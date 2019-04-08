# -*- coding:utf-8 -*-
# edit by fuzongfei
import json
from datetime import datetime

import sqlparse
from django.utils import timezone
from rest_framework import serializers

from orders.models import Orders, OrderReply, OnlineVersion
from orders.utils.inceptionApi import InceptionSqlApi
from orders.utils.msgNotice import MsgPush
from orders.utils.tools import split_sqltype
from orders.utils.validatorsCheck import sql_type_validator, envi_validator, file_foramt_validator


class OrdersCommitSerializer(serializers.Serializer):
    title = serializers.CharField(max_length=100, required=True)
    description = serializers.CharField(max_length=1024, required=False)
    version = serializers.IntegerField(required=False)
    auditor = serializers.CharField(required=True, error_messages={'required': '工单审核人不能为空'})
    reviewer = serializers.CharField(required=True, error_messages={'required': '工单复核人不能为空'})
    email_cc = serializers.CharField(required=False)
    schema = serializers.CharField(required=False, max_length=256)
    remark = serializers.CharField(required=True, max_length=256, min_length=1, error_messages={'required': '备注不能为空'})
    sql_type = serializers.CharField(validators=[sql_type_validator], default='')
    file_format = serializers.CharField(validators=[file_foramt_validator], default='xlsx')
    envi_id = serializers.CharField(required=False, validators=[envi_validator],
                                    error_messages={'required': 'envi_id字段不能为空'})
    contents = serializers.CharField(allow_blank=False, error_messages={'required': '提交的内容不能为空'})

    def convert_to_dict(self, data):
        """
        auditor = reviewer = [
            {'user': 'zhangsan', 'is_superuser': 0, 'status': 0, 'time': '', 'msg': ''},
            ...,
            {'user': 'lisi', 'is_superuser': 1, 'status': 0, 'time': '', 'msg': ''}
            ]
        status:
            0：未审核或未复核
            1：已审核或已复核
        is_superuser:
            0：不是一键审核人或复核人
            1：是一键审核人或复核人
        time：操作的时间
        msg：附加的消息
        """
        r = []
        for i in data:
            r.append({'user': i, 'is_superuser': 0, 'status': 0, 'time': '', 'msg': ''})
        return json.dumps(r)

    def check_sqlnumber(self, data):
        """单次最大支持1000条DML和DDL语句提交"""
        sql_list = [sql for sql in sqlparse.split(data)]
        if len(sql_list) > 1000:
            return False, len(sql_list)
        return True, None

    def createobj(self, data):
        obj = Orders.objects.create(**data)
        return obj

    def save(self, request):
        sdata = self.validated_data
        host, port, database = sdata.get('schema').split(',') if sdata.get('schema') else [0, 0, '']
        sql_type = sdata.get('sql_type')
        contents = sdata.get('contents')
        version_id = sdata.get('version') if sdata.get('version') else -1
        data = {
            'title': sdata.get('title') + '_[' + datetime.now().strftime("%Y%m%d%H%M%S") + ']',
            'description': sdata.get('description'),
            'version_id': version_id,
            'applicant': request.user.username,
            'auditor': self.convert_to_dict(self.initial_data.getlist('auditor')),
            'reviewer': self.convert_to_dict(self.initial_data.getlist('reviewer')),
            'email_cc': ','.join(self.initial_data.getlist('email_cc')),
            'remark': sdata.get('remark'),
            'host': host,
            'port': port if isinstance(port, int) else int(port),
            'database': database,
            'sql_type': sql_type,
            'contents': contents,
            'file_format': sdata.get('file_format'),
            'envi_id': sdata.get('envi_id'),
        }

        # DML和DDL语句需要检查语法
        if sql_type in ('DML', 'DDL'):
            # 检查语句条数
            status, length = self.check_sqlnumber(contents)
            if not status:
                msg = f'最大支持一次提交1000条SQL语句，当前条数: {length}'
                return False, msg

            # 检查提交的语句类型是否匹配
            status, msg = split_sqltype(sql=contents, sql_type=sql_type)
            if not status:
                return False, msg

            # Inception语法检查
            result = InceptionSqlApi(host, port, database, contents, request.user.username).is_check_pass()
            if result.get('status') == 2:
                return False, result['msg']

        obj = self.createobj(data)
        # 推送消息
        msg_push = MsgPush(id=obj.id, type='commit')
        msg_push.send()
        # 跳转到工单记录页面
        jump_url = f'/orders/envi/{obj.envi_id}/'
        return True, jump_url


class OrderReplySerializer(serializers.Serializer):
    reply_id = serializers.IntegerField(required=True, error_messages={'required': '回复工单id不能为空'})
    reply_contents = serializers.CharField(required=True, min_length=2, max_length=2048, error_messages={
        'required': '回复内容不能为空',
        'min_length': '回复至少输入2个字符',
        'max_length': '回复最多2048个字符'
    })

    def save(self, request):
        sdata = self.validated_data
        reply_id = sdata.get('reply_id')
        reply_contents = sdata.get('reply_contents')
        obj = OrderReply.objects.create(
            reply_id=reply_id,
            user_id=request.user.uid,
            reply_contents=reply_contents)

        # 推送消息
        msg_push = MsgPush(id=obj.id, user=request.user.username, type='reply')
        msg_push.send()
        return True, '回复成功'


class HookOrdersSerializer(serializers.Serializer):
    id = serializers.IntegerField(required=True, error_messages={'required': '工单id不能为空'})
    schema = serializers.CharField(required=True, max_length=256, error_messages={'required': '库不能为空'})
    envi_id = serializers.CharField(required=True, validators=[envi_validator],
                                    error_messages={'required': 'envi_id字段不能为空'})

    def save(self, request):
        sdata = self.validated_data
        id = sdata.get('id')
        envi_id = sdata.get('envi_id')
        jump_url = f'/orders/envi/{envi_id}'

        data = Orders.objects.get(pk=id)
        if data.progress == '7':
            return False, '当前工单已被勾住， 请不要重复执行'
        else:
            # OPS运维工单默认0
            host, port, database = [0, 0, '']
            if data.sql_type in ['DML', 'DDL']:
                host, port, database = sdata['schema'].split(',')

            # 已复核的工单才能执行hook
            if data.progress in ['6']:
                # 重置复核信息
                reviewer = json.loads(data.reviewer)
                for i in reviewer:
                    i['status'] = 0
                    i['msg'] = ''
                    i['time'] = ''
                obj = Orders.objects.create(
                    title=data.title,
                    description=data.description,
                    version_id=data.version_id,
                    sql_type=data.sql_type,
                    host=host,
                    database=database,
                    port=port,
                    envi_id=envi_id,
                    progress='2',
                    remark=data.remark,
                    applicant=data.applicant,
                    auditor=data.auditor,
                    reviewer=json.dumps(reviewer),
                    file_format=data.file_format,
                    contents=data.contents,
                    updated_at=timezone.now()
                )

                # 更新状态为：已勾住
                data.progress = '7'
                data.save()

                # 推送消息
                msg_push = MsgPush(id=obj.id, user=request.user.username, type='hook')
                msg_push.send()

                # 跳转到工单记录页面
                return True, jump_url
            else:
                return False, '当前工单未完成，无法执行钩子操作'


class OnlineVersionListSerializer(serializers.Serializer):
    id = serializers.CharField(required=False)
    version = serializers.CharField(required=False)
    expire_time = serializers.CharField(required=False)
    action = serializers.ChoiceField(choices=(('create', 'create'), ('delete', 'delete')))

    def op(self, request):
        sdata = self.validated_data
        action = sdata.get('action')
        version = sdata.get('version')
        expire_time = sdata.get('expire_time')

        if action == 'create':
            if OnlineVersion.objects.filter(version=version).exists():
                return False, '记录已存在，不能重复创建'
            else:
                OnlineVersion.objects.create(version=version, expire_time=expire_time,
                                             username=request.user.displayname)
                return True, '记录创建成功'
        elif action == 'delete':
            id = self.initial_data.getlist('id')[0]
            for i in id.split(','):
                OnlineVersion.objects.filter(pk=i).update(is_deleted='1')
            return True, '记录删除成功'
