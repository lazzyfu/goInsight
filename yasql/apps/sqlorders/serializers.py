# -*- coding:utf-8 -*-
# edit by fuzongfei
import json
import os
import re
import subprocess
import uuid
from datetime import datetime

import sqlparse
from django.core.cache import cache
from django.db.models import F
from django.utils import timezone
from django.utils.html import linebreaks
from rest_framework import serializers
from rest_framework.exceptions import PermissionDenied

from config import REOMOTE_USER
from sqlorders import models, utils, libs, tasks
from sqlorders.api.goInceptionApi import InceptionApi
from sqlorders.libs import check_export_column_unique
from users.models import UserAccounts


class ReleaseVersionsSerializer(serializers.ModelSerializer):
    class Meta:
        model = models.ReleaseVersions
        fields = ['username', 'version', 'expire_time', 'created_at']

    def to_representation(self, instance):
        ret = super(ReleaseVersionsSerializer, self).to_representation(instance)
        ret["created_at"] = datetime.strftime(instance.created_at, "%Y-%m-%d %H:%M:%S")
        return ret


class DbEnvironmentSerializer(serializers.ModelSerializer):
    class Meta:
        model = models.DbEnvironment
        fields = ['id', 'name']


class DbSchemasSerializer(serializers.Serializer):
    env_id = serializers.CharField()
    use_type = serializers.IntegerField()
    rds_category = serializers.ChoiceField(choices=utils.rdsCategory)

    @property
    def query(self):
        super(DbSchemasSerializer, self)
        vdata = self.validated_data
        queryset = models.DbSchemas.objects.filter(
            cid__env_id=vdata['env_id'],
            cid__use_type=vdata['use_type'],
            cid__rds_category=vdata['rds_category']
        ).annotate(host=F('cid__host'),
                   port=F('cid__port'),
                   comment=F('cid__comment')
                   ).values('id', 'cid', 'host', 'port', 'schema', 'comment')
        return queryset


class IncepSyntaxCheckSerializer(serializers.Serializer):
    rds_category = serializers.ChoiceField(choices=utils.rdsCategory,
                                           error_messages={
                                               'required': '请先选择左侧的【DB类别】',
                                               'blank': '请先选择左侧的【DB类别】'
                                           })
    database = serializers.CharField(
        error_messages={
            'required': '请先选择左侧的【环境】和【库名】',
            'blank': '请先选择左侧的【环境】和【库名】'
        }
    )
    sqls = serializers.CharField(error_messages={
        'required': '审核内容不能为空',
        'blank': '审核内容不能为空'
    })
    sql_type = serializers.ChoiceField(choices=utils.sqlTypeChoice)

    def validate(self, attrs):
        # ddl工单提交ddl语句，dml工单提交dml语句
        status, msg = libs.verify_sql_type(sqls=attrs['sqls'], sql_type=attrs['sql_type'])
        if not status:
            raise serializers.ValidationError(msg)

        # EXPORT类型的工单不需要检查语法
        if attrs['sql_type'] == 'EXPORT' and attrs['rds_category'] != 3:
            cid, database = attrs['database'].split('__')
            obj = models.DbConfig.objects.get(pk=cid)
            cfg = {
                'host': obj.host,
                'port': obj.port,
                'database': database
            }
            cfg.update(REOMOTE_USER)

            # 判断导出工单提交的列是否重复
            status, msg = check_export_column_unique(config=cfg, sqls=attrs['sqls'])
            if not status:
                raise serializers.ValidationError(msg)

            raise serializers.ValidationError('EXPORT类型的工单不需要语法检查，请直接提交')

        # TiDB的ALTER语句需要单独的处理
        if attrs['rds_category'] == 2:
            if attrs['sql_type'] == 'DDL':
                sc = re.compile(r'^ALTER([\s\S]+)([,])(\s+)(ADD|CHANGE|RENAME|MODIFY|DROP|CONVERT)([\s\S]+)', re.I)
                for sql in sqlparse.split(attrs['sqls']):
                    try:
                        m = sc.match(sql)
                        if m.groups():
                            raise serializers.ValidationError('TiDB下的同一个表的多个操作，请拆分成多个ALTER语句')
                    except AttributeError:
                        pass

        # 检查goInception是否可以访问
        status, msg = InceptionApi().check_cnx()
        if not status:
            raise serializers.ValidationError(msg)

        return super(IncepSyntaxCheckSerializer, self).validate(attrs)

    def check(self):
        vdata = self.validated_data
        if vdata['rds_category'] == 3:
            return True, None
        cid, database = vdata['database'].split('__')
        obj = models.DbConfig.objects.get(pk=cid)
        cfg = {
            'host': obj.host,
            'port': obj.port,
            'database': database
        }
        cfg.update(REOMOTE_USER)
        api = InceptionApi(cfg=cfg, sqls=vdata['sqls'], rds_category=vdata['rds_category'])
        context = api.run_check()
        status = api.is_check_pass()
        return status, context


class SqlOrdersCommitSerializer(serializers.ModelSerializer):
    env_id = serializers.IntegerField()

    class Meta:
        model = models.DbOrders
        exclude = ['applicant', 'auditor', 'reviewer', 'email_cc', 'cid', 'order_id']

    def convert_to_dict(self, data):
        """
        auditor = reviewer = [
            {'user': 'zhangsan', 'is_superuser': 0, 'status': 0, 'msg': '', 'time': ''},
            ...
            ]
        status:
            0：未审核或未复核
            1：已审核或已复核
            2: 已驳回
        is_superuser:
            0：不是一键审核人或复核人
            1：是一键审核人或复核人
        msg: 附加的信息
        time: 操作时间
        """
        r = []
        for i in data:
            r.append({'user': i, 'is_superuser': 0, 'status': 0, 'msg': '', 'time': ''})
        return json.dumps(r)

    def check_number(self, data):
        """单次最大支持2048条DML和DDL语句提交"""
        sql_list = [sql for sql in sqlparse.split(data)]
        if len(sql_list) > 2048:
            return False, len(sql_list)
        return True, None

    def validate_sql_type(self, data):
        if data in ('DDL', 'DML'):
            # 检查语句条数
            status, length = self.check_number(self.initial_data.get('contents'))
            if not status:
                raise serializers.ValidationError(f'单次最大支持一次提交2048条SQL语句，当前条数: {length}')
        return data

    def validate(self, attrs):
        # 判断提交的SQL是否符合SQL类型
        # 如：DML只能提交DML语句，DDL工单只能提交DDL语句
        status, msg = libs.verify_sql_type(sqls=attrs['contents'], sql_type=attrs['sql_type'])
        if not status:
            raise serializers.ValidationError(msg)

        # 当语法未检查通过时，拦截提交
        cid, database = attrs['database'].split('__')
        obj = models.DbConfig.objects.get(pk=cid)
        cfg = {
            'host': obj.host,
            'port': obj.port,
            'database': database
        }
        cfg.update(REOMOTE_USER)
        api = InceptionApi(cfg=cfg, sqls=attrs['contents'], rds_category=attrs['rds_category'])

        # 检查goInception是否可以访问
        status, msg = api.check_cnx()
        if not status:
            raise serializers.ValidationError(msg)

        # # 禁止提交insert into ... select 语句
        # 需要屏蔽insert into ... select 语句，取消注释
        # if attrs['sql_type'] == 'DML':
        #     if not api.check_insert_select():
        #         raise serializers.ValidationError('禁止提交INSERT INTO ... SELECT ...语句')

        # 判断导出工单提交的列是否重复
        if attrs['sql_type'] == 'EXPORT' and attrs['rds_category'] != 3:
            status, msg = check_export_column_unique(config=cfg, sqls=attrs['contents'])
            if not status:
                raise serializers.ValidationError(msg)

        # 导出工单不检查语法，仅检测是否以SELECT开头
        if attrs['sql_type'] != 'EXPORT':
            if not api.is_check_pass():
                raise serializers.ValidationError('SQL语法存在异常，提交失败，请先检查SQL语法，请点击：[语法检查]')

        # TiDB的ALTER语句需要单独的处理
        if attrs['rds_category'] == 2:
            if attrs['sql_type'] == 'DDL':
                sc = re.compile(r'^ALTER([\s\S]+)([,])(\s+)(ADD|CHANGE|RENAME|MODIFY|DROP|CONVERT)([\s\S]+)', re.I)
                for sql in sqlparse.split(attrs['contents']):
                    try:
                        m = sc.match(sql)
                        if m.groups():
                            raise serializers.ValidationError('TiDB下的同一个表的多个操作，请拆分成多个ALTER语句')
                    except AttributeError:
                        pass

        # 提交工单
        request = self.context['request']
        attrs["applicant"] = request.user.username
        attrs["auditor"] = self.convert_to_dict(self.initial_data.get('auditor'))
        attrs["reviewer"] = self.convert_to_dict(self.initial_data.get('reviewer'))
        attrs["executor"] = self.convert_to_dict(['None'])
        attrs["closer"] = self.convert_to_dict(['None'])
        attrs["email_cc"] = ','.join(self.initial_data.get('email_cc'))
        attrs['title'] = attrs['title'] + '_[' + datetime.now().strftime("%Y%m%d%H%M%S") + ']'
        attrs['order_id'] = ''.join(str(uuid.uuid4()).split('-'))  # 基于UUID生成工单id
        attrs['cid_id'], attrs['database'] = attrs['database'].split('__')
        attrs["department"] = request.user.department

        return super(SqlOrdersCommitSerializer, self).validate(attrs)

    def save(self, **kwargs):
        super(SqlOrdersCommitSerializer, self).save()

        # 推送消息
        tasks.msg_notice.delay(
            pk=self.instance.pk,
            op='_commit',
            username=self.context['request'].user.username
        )


class SqlOrdersListSerializer(serializers.ModelSerializer):
    applicant = serializers.SerializerMethodField()
    auditor = serializers.SerializerMethodField()
    reviewer = serializers.SerializerMethodField()
    # 获取choices字段
    progress = serializers.SerializerMethodField()
    # 获取外键的字段
    version = serializers.ReadOnlyField(source='version.version')
    host = serializers.ReadOnlyField(source='cid.host')
    port = serializers.ReadOnlyField(source='cid.port')
    # 环境
    env_name = serializers.ReadOnlyField(source='env.name')

    class Meta:
        model = models.DbOrders
        fields = ['id', 'title', 'progress', 'remark', 'env_name', 'sql_type', 'file_format', 'department', 'is_hide',
                  'applicant', 'order_id', 'auditor', 'reviewer', 'version', 'cid', 'host', 'port', 'database',
                  'created_at']

    def get_applicant(self, obj):
        try:
            obj.applicant = UserAccounts.objects.get(username=obj.applicant).displayname
        except UserAccounts.DoesNotExist:
            pass
        return obj.applicant

    def get_auditor(self, obj):
        data = json.loads(obj.auditor)
        for row in data:
            try:
                row['display_name'] = UserAccounts.objects.get(username=row['user']).displayname
            except UserAccounts.DoesNotExist:
                pass
        return json.dumps(data)

    def get_reviewer(self, obj):
        data = json.loads(obj.reviewer)
        for row in data:
            try:
                row['display_name'] = UserAccounts.objects.get(username=row['user']).displayname
            except UserAccounts.DoesNotExist:
                pass
        return json.dumps(data)

    def get_progress(self, obj):
        obj.progress = dict(utils.sqlProgressChoice).get(obj.progress)
        return obj.progress

    def to_representation(self, instance):
        ret = super(SqlOrdersListSerializer, self).to_representation(instance)
        ret["escape_title"] = instance.title
        ret["created_at"] = datetime.strftime(instance.created_at, "%Y-%m-%d %H:%M:%S")
        return ret


class SqlOrderDetailSerializer(serializers.ModelSerializer):
    # 获取choices字段
    env_id = serializers.SerializerMethodField()
    progress = serializers.SerializerMethodField()
    host = serializers.ReadOnlyField(source='cid.host')
    port = serializers.ReadOnlyField(source='cid.port')

    class Meta:
        model = models.DbOrders
        fields = "__all__"

    def get_env_id(self, obj):
        obj.env_id = models.DbEnvironment.objects.get(pk=obj.env_id).name
        return obj.env_id

    def get_progress(self, obj):
        obj.progress = dict(utils.sqlProgressChoice).get(obj.progress)
        return obj.progress

    def to_representation(self, instance):
        ret = super(SqlOrderDetailSerializer, self).to_representation(instance)
        ret["created_at"] = datetime.strftime(instance.created_at, "%Y-%m-%d %H:%M:%S")
        ret["display_rds_category"] = dict(utils.rdsCategory).get(instance.rds_category)

        # 数据隐藏按钮打开了
        # 仅允许申请人、审核人、复核人和超权用户查看数据
        if instance.is_hide == 'ON' and not self.context['request'].user.is_superuser:
            allowed_view_users = [instance.applicant]
            allowed_view_users.extend([x['user'] for x in json.loads(instance.auditor)])
            allowed_view_users.extend([x['user'] for x in json.loads(instance.reviewer)])
            if self.context['request'].user.username not in allowed_view_users:
                raise PermissionDenied(detail='您没有权限查看该工单的数据，5s后，自动跳转到工单列表页面')
        return ret


class OpSqlOrderSerializer(serializers.Serializer):
    action = serializers.CharField(required=False, max_length=32, allow_null=True, allow_blank=True)
    msg = serializers.CharField(required=False, max_length=128, allow_null=True, allow_blank=True)
    btn = serializers.CharField(required=True, max_length=12, allow_null=True, allow_blank=True)

    def validate(self, attrs):
        f = getattr(self, self.context["handler"])
        data = f(self.context["request"], attrs)
        return data

    def update(self, instance, validated_data):
        for k, v in validated_data.items():
            setattr(instance, k, v)
        instance.save()
        return instance

    def save(self, **kwargs):
        super(OpSqlOrderSerializer, self).save()
        # 推送消息
        tasks.msg_notice.delay(
            pk=self.instance.pk,
            op=self.context['handler'],
            username=self.context['request'].user.username
        )

    def _approve(self, request, attr):
        # 状态不为：待审核，raise error
        if self.instance.progress not in (0,):
            raise serializers.ValidationError({'status': '非可审核状态，禁止操作'}, code='status')

        data = {}
        auditor = json.loads(self.instance.auditor)
        allowed_auditor = [x.get('user') for x in auditor]

        # 超级审核人
        if request.user.is_superuser:
            if attr['btn'] == 'ok':
                data['progress'] = 2
                status = 1
            else:
                data['progress'] = 1  # 当点击驳回按钮时，将工单状态设置为已驳回
                status = 2
            auditor.append({
                'user': request.user.username,
                'is_superuser': 1,
                'status': status,
                'time': datetime.now().strftime("%Y-%m-%d %H:%M:%S"),
                'msg': attr['msg']
            })
        # 普通审核人
        if not request.user.is_superuser:
            if request.user.username in allowed_auditor:
                for i in auditor:
                    if request.user.username == i['user']:
                        if i['status'] == 1:
                            raise serializers.ValidationError({'status': '您已审核过，请不要重复执行'}, code='status')
                        else:
                            i['user'] = request.user.username
                            i['time'] = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
                            i['msg'] = attr['msg']
                            if attr['btn'] == 'ok':
                                i['status'] = 1
                                # 当所有非超级管理员的审核人员批准后，修改工单状态为已完成
                                f_status = [x['status'] for x in auditor if x['is_superuser'] == 0]
                                if 1 in f_status and all(f_status):
                                    data['progress'] = 2
                            else:
                                i['status'] = 2
                                data['progress'] = 1  # 当点击驳回按钮时，将工单状态设置为已驳回
            else:
                raise serializers.ValidationError({'status': '权限拒绝，您没有当前工单的审核权限'}, code='status')
        data['auditor'] = json.dumps(auditor)
        return data  # update方法处理

    def _feedback(self, request, attr):
        # 状态不为：待审核/已批准/处理中，raise error
        if self.instance.progress not in (0, 2, 3):
            raise serializers.ValidationError({'status': '非可反馈状态，禁止操作'}, code='status')
        # 更新执行人的信息
        executor = json.loads(self.instance.executor)
        executor[0]['user'] = request.user.username
        executor[0]['time'] = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
        executor[0]['msg'] = attr['msg']
        # 当点击处理中，修改工单状态为处理中
        # 当点击执行完成按钮时，修改工单状态为已完成
        data = {'progress': 3 if attr['btn'] == 'ok' else 4, 'executor': json.dumps(executor)}
        return data

    def _close(self, request, attr):
        # 状态不为：待审核/已批准/处理中，raise error
        if self.instance.progress not in (0, 2, 3):
            raise serializers.ValidationError({'status': '非可关闭状态，禁止操作'}, code='status')
        # 更新关闭人的信息
        closer = json.loads(self.instance.closer)
        closer[0]['user'] = request.user.username
        closer[0]['time'] = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
        closer[0]['msg'] = attr['msg']
        # 当点击执行中，修改工单状态为执行中
        # 当点击执行完成按钮时，修改工单状态为已完成
        if attr['btn'] == 'ok':
            return {'progress': 5, 'closer': json.dumps(closer)}

    def _review(self, request, attr):
        # 状态不为：已完成，raise error
        if self.instance.progress not in (4,):
            raise serializers.ValidationError({'status': '工单未完成，不能复核'}, code='status')
        # 更新关闭人的信息
        data = {}
        reviewer = json.loads(self.instance.reviewer)
        allowed_reviewer = [x.get('user') for x in reviewer]
        # 超级审核人无效，必须每个复核人进行复核
        if attr['btn'] == 'ok':
            if request.user.username in allowed_reviewer:
                for i in reviewer:
                    if request.user.username == i['user']:
                        if i['status'] == 1:
                            raise serializers.ValidationError({'status': '您已复核过，请不要重复执行'}, code='status')
                        else:
                            i['user'] = request.user.username
                            i['time'] = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
                            i['msg'] = attr['msg']
                            i['status'] = 1
                            # 当所有非超级管理员的审核人员批准后，修改工单状态为已完成
                            f_status = [x['status'] for x in reviewer if x['is_superuser'] == 0]
                            if 1 in f_status and all(f_status):
                                data['progress'] = 6
                data['reviewer'] = json.dumps(reviewer)
                return data  # update方法处理
            else:
                raise serializers.ValidationError({'status': '权限拒绝，您没有当前工单的审核权限'}, code='status')
        return data


class GenerateSqlOrdersTasksSerializer(serializers.Serializer):
    id = serializers.IntegerField()

    def validate_id(self, data):
        obj = models.DbOrders.objects.get(pk=data)
        if obj.progress not in (2, 3, 4, 5, 6, 7):
            raise serializers.ValidationError('工单状态为：未通过/已驳回，禁止操作')
        return data

    def save(self, request):
        id = self.validated_data['id']
        obj = models.DbOrders.objects.get(pk=id)
        # 如果记录已经存在，直接跳转
        if models.DbOrdersExecuteTasks.objects.filter(order__id=id).exists():
            return models.DbOrdersExecuteTasks.objects.filter(order__id=id).first().task_id
        # 生成记录
        splitsqls = [sql.strip(';') for sql in sqlparse.split(obj.contents, encoding='utf8')]
        task_id = ''.join(str(uuid.uuid4()).split('-'))  # 基于UUID生成任务ID
        for sql in splitsqls:
            models.DbOrdersExecuteTasks.objects.create(
                applicant=obj.applicant,
                task_id=task_id,
                sql=sql.strip(';'),
                sql_type=obj.sql_type,
                file_format=obj.file_format,
                order_id=id
            )
        return task_id


class SqlOrdersTasksListSerializer(serializers.ModelSerializer):
    progress = serializers.SerializerMethodField()

    class Meta:
        model = models.DbOrdersExecuteTasks
        exclude = ['execute_log', 'rollback_sql']

    def get_progress(self, obj):
        # 格式化进度
        obj.progress = dict(utils.taskProgressChoice).get(obj.progress)
        return obj.progress


class ExecuteSingleTaskSerializer(serializers.Serializer):
    id = serializers.IntegerField()

    def check_task_locked(self, id):
        key = f'execute_single_task_id_{str(id)}'
        if cache.get(key):
            raise serializers.ValidationError('当前任务已被触发，10s内请不要重复点击')
        # 锁定10s
        cache.set(key, 1, 10)

    def validate_id(self, data):
        # 检查任务是否被锁定，避免连击
        self.check_task_locked(data)

        try:
            obj = models.DbOrdersExecuteTasks.objects.get(pk=data)
        except models.DbOrdersExecuteTasks.DoesNotExist as err:
            raise serializers.ValidationError('查询的结果不存在')

        # 当父工单的状态不为已批准或处理中的时候，禁止执行
        if models.DbOrders.objects.get(pk=obj.order_id).progress not in (2, 3):
            raise serializers.ValidationError('工单状态为：已驳回/已关闭/已完成，禁止操作')

        progress_distinct = models.DbOrdersExecuteTasks.objects.filter(
            task_id=obj.task_id
        ).distinct().values_list('progress', flat=True)
        if obj.progress == 1:
            # 避免重复点击已完成的任务
            raise serializers.ValidationError('当前任务已经完成，请不要重复执行')
        if 2 in progress_distinct:
            # 每次只能有一条SQL在执行，避免数据不一致或数据库压力
            raise serializers.ValidationError('当前有任务正在执行中，请不要重复执行')
        return data

    def execute(self, request):
        id = self.validated_data.get('id')
        obj = models.DbOrdersExecuteTasks.objects.get(pk=id)
        # 更新父工单的状态为：处理中
        tasks.update_dborders_progress_to_processing(obj.order_id, request.user.username)
        # 将当前任务进度设置为：处理中
        obj.executor = request.user.username
        obj.execute_time = timezone.now()
        obj.progress = 2
        obj.save()
        # 执行
        tasks.async_execute_single.delay(id=id, username=request.user.username)

        # 推送消息
        tasks.msg_notice.delay(
            pk=obj.order_id,
            op='_feedback',
            username=request.user.username
        )


class ExecuteMultiTasksSerializer(serializers.Serializer):
    task_id = serializers.CharField()

    def check_task_locked(self, taskid):
        key = f'execute_multi_tasks_taskid_{taskid}'
        if cache.get(key):
            raise serializers.ValidationError('当前任务已被触发，10s内请不要重复点击')
        # 锁定10s
        cache.set(key, 1, 10)

    def validate_task_id(self, data):
        # 检查任务是否被锁定，避免连击
        self.check_task_locked(data)
        
        if not models.DbOrdersExecuteTasks.objects.filter(task_id=data).exists():
            raise serializers.ValidationError('查询的结果不存在')

        # 当父工单的状态不为已批准或处理中的时候，禁止执行
        if models.DbOrdersExecuteTasks.objects.filter(task_id=data).exists():
            obj = models.DbOrdersExecuteTasks.objects.filter(task_id=data).first()
            if models.DbOrders.objects.get(pk=obj.order_id).progress not in (2, 3):
                raise serializers.ValidationError('工单状态为：已驳回/已关闭/已完成，禁止操作')

        progress_distinct = models.DbOrdersExecuteTasks.objects.filter(
            task_id=data
        ).distinct().values_list('progress', flat=True)
        if 2 in progress_distinct:
            # 每次只能有一条SQL在执行，避免数据不一致或数据库压力
            raise serializers.ValidationError('当前有任务正在执行中，请不要重复执行')
        return data

    def execute(self, request):
        task_id = self.validated_data.get('task_id')
        # 更新父工单的状态为：处理中
        id = models.DbOrdersExecuteTasks.objects.filter(task_id=task_id).first().order_id
        tasks.update_dborders_progress_to_processing(id, request.user.username)
        # 执行
        tasks.async_execute_multi.delay(task_id=task_id, username=request.user.username)

        # 推送消息
        tasks.msg_notice.delay(
            pk=id,
            op='_feedback',
            username=request.user.username
        )


class ThrottleTaskSerializer(serializers.Serializer):
    id = serializers.IntegerField()
    op = serializers.ChoiceField(choices=(('pause', 'pause'), ('recovery', 'recovery')))

    def validate_id(self, data):
        try:
            obj = models.DbOrdersExecuteTasks.objects.get(pk=data)
        except models.DbOrdersExecuteTasks.DoesNotExist as err:
            raise serializers.ValidationError('查询的结果不存在')

        if obj.progress in (0, 1, 3):
            raise serializers.ValidationError('非执行中的任务，禁止操作')

        return data

    def execute(self, request):
        op = self.validated_data['op']
        obj = models.DbOrdersExecuteTasks.objects.get(pk=self.validated_data['id'])
        database = models.DbOrders.objects.get(pk=obj.order_id).database
        sql = libs.remove_sql_comment(obj.sql)
        try:
            syntaxcompile = re.compile(r'^ALTER(\s+)TABLE(\s+)([\S]*)(.*)', re.I)
            syntaxmatch = syntaxcompile.match(sql)
            # 由于gh-ost不支持反引号，会被解析成命令，因此此处替换掉
            table = syntaxmatch.group(3).replace('`', '')
            # 将schema.table进行处理，这种情况gh-ost不识别，只保留table
            if len(table.split('.')) > 1:
                table = table.split('.')[1]
            sock = f"/tmp/gh-ost.{database}.{table}.sock"
            if os.path.exists(sock):
                """
                echo throttle | nc -U /tmp/gh-ost.test.sample_data_0.sock
                echo no-throttle | nc -U /tmp/gh-ost.test.sample_data_0.sock
                """
                if op == 'pause':
                    pause_cmd = f"echo throttle | nc -U {sock}"
                    p = subprocess.Popen(pause_cmd, shell=True)
                    p.wait()
                    obj.progress = 4
                    obj.execute_time = timezone.now()
                    obj.save()
                    return '暂停操作已执行，请查看输出'
                if op == 'recovery':
                    recovery_cmd = f"echo no-throttle | nc -U {sock}"
                    p = subprocess.Popen(recovery_cmd, shell=True)
                    p.wait()
                    obj.progress = 3
                    obj.execute_time = timezone.now()
                    obj.save()
                    return '恢复操作已执行，请查看输出'
            else:
                return f'操作失败：不能发现sock文件[{sock}]'
        except AttributeError as err:
            return '非有效的ALTER语句，禁止操作'


class GetTasksResultSerializer(serializers.ModelSerializer):
    class Meta:
        model = models.DbOrdersExecuteTasks
        fields = ['rollback_sql', 'execute_log']

    def get_execute_log(self, obj):
        return linebreaks(obj.execute_log)


class HookSqlOrdersSerializer(serializers.ModelSerializer):
    id = serializers.IntegerField()
    reset = serializers.ChoiceField(choices=(('ON', 'ON'), ('OFF', 'OFF')))
    env_id = serializers.IntegerField()

    class Meta:
        model = models.DbOrders
        fields = ['id', 'reset', 'database', 'remark', 'env_id']

    def validate_id(self, data):
        try:
            obj = models.DbOrders.objects.get(pk=data)
        except models.DbOrders.DoesNotExist:
            raise serializers.ValidationError('查询的结果不存在')
        if obj.progress not in (6,):
            raise serializers.ValidationError('仅状态为【已复核】的工单，才能使用钩子')
        return data

    def reset_json_value(self, data):
        json_data = json.loads(data)
        for row in json_data:
            if row['is_superuser'] == 0:
                row['status'] = 0
                row['msg'] = ''
                row['time'] = ''
            else:
                json_data.remove(row)
        return json.dumps(json_data)

    def validate(self, attrs):
        obj = models.DbOrders.objects.get(pk=attrs['id'])
        # 判断工单是否存在
        if models.DbOrders.objects.filter(title=obj.title, env=attrs['env_id']).exists():
            raise serializers.ValidationError(f"目标【{models.DbEnvironment.objects.get(pk=attrs['env_id'])}】工单已存在")
        return super(HookSqlOrdersSerializer, self).validate(attrs)

    def save(self):
        value = [{"user": "None", "is_superuser": 0, "status": 0, "msg": "", "time": ""}]
        vdata = self.validated_data
        obj = models.DbOrders.objects.get(pk=vdata['id'])

        # 新建hook工单
        cid_id, database = vdata['database'].split('__')
        # 重置审核状态，需要Leader重新审核
        auditor = obj.auditor
        progress = 2
        if vdata['reset'] == 'ON':
            auditor = self.reset_json_value(obj.auditor)
            progress = 0
        models.DbOrders.objects.create(
            title=obj.title,
            order_id=''.join(str(uuid.uuid4()).split('-')),
            demand=obj.demand,
            remark=vdata['remark'],
            env_id=vdata['env_id'],
            database=database,
            rds_category=obj.rds_category,
            sql_type=obj.sql_type,
            file_format=obj.file_format,
            applicant=obj.applicant,
            auditor=auditor,
            department=obj.department,
            reviewer=self.reset_json_value(obj.reviewer),
            executor=json.dumps(value),
            closer=json.dumps(value),
            email_cc=obj.email_cc,
            cid_id=cid_id,
            progress=progress,
            version=obj.version,
            contents=obj.contents,
            created_at=timezone.now(),
            updated_at=timezone.now()
        )
        # 更新工单状态为已勾住
        obj.progress = 7
        obj.save()


class ReleaseVersionsListSerializer(serializers.ModelSerializer):
    class Meta:
        model = models.ReleaseVersions
        fields = "__all__"

    def to_representation(self, instance):
        ret = super(ReleaseVersionsListSerializer, self).to_representation(instance)
        ret["key"] = instance.pk
        ret["created_at"] = datetime.strftime(instance.created_at, "%Y-%m-%d %H:%M:%S")
        ret["expire_time"] = datetime.strftime(instance.expire_time, "%Y-%m-%d")
        return ret


class ReleaseVersionsCreateSerializer(serializers.ModelSerializer):
    class Meta:
        model = models.ReleaseVersions
        fields = ['username', 'version', 'expire_time']
