# -*- coding:utf-8 -*-
# edit by fuzongfei

from django.db import models

# 工单进度
from users.models import UserAccounts

progressChoice = (
    ('0', u'待批准'),
    ('1', u'未批准'),
    ('2', u'已批准'),
    ('3', u'处理中'),
    ('4', u'已完成'),
    ('5', u'已关闭'),
    ('6', u'已复核'),
    ('7', u'已勾住')
)

# 操作类型选择
# OPS为运维工单
sqlTypeChoice = (
    ('DDL', u'DDL工单'),
    ('DML', u'DML工单'),
    ('OPS', u'运维工单'),
    ('EXPORT', u'数据导出工单')
)

# 导出工单支持的文件格式
fileFormatChoice = (
    ('xlsx', 'xlsx格式'),
    ('csv', 'csv格式')
)

# 支持多人审核
# 支持多人复核
"""
status: 0表示未审核，1表示已审核
is_superuser：0表示不是一键审核人，1表示是一键审核人，即该用户为超级管理员时，可一键通过审核
auditor = reviewer = [
            {'user': 'zhangsan', 'is_superuser': 0, 'status': 0, 'time': '', 'msg': ''},
            ...,
            {'user': 'lisi', 'is_superuser': 1, 'status': 0, 'time': '', 'msg': ''}
            ]
close_info = {'user': 'zs', 'msg': '关闭原因', 'time': '关闭时间'}
"""


class Orders(models.Model):
    id = models.AutoField(primary_key=True, verbose_name=u'主键ID')
    title = models.CharField(max_length=100, verbose_name=u'工单标题')
    description = models.CharField(max_length=2048, default='', null=False, verbose_name=u'工单需求描述')
    remark = models.CharField(max_length=128, default='', null=False, verbose_name=u'工单备注')
    sql_type = models.CharField(max_length=30, default='DML', choices=sqlTypeChoice, verbose_name=u'工单类型')
    file_format = models.CharField(max_length=30, choices=fileFormatChoice, default='xlsx', verbose_name=u'文件格式')
    envi = models.ForeignKey('orders.SysEnvironment', null=False, to_field='envi_id', on_delete=models.CASCADE,
                             verbose_name=u'环境')
    applicant = models.CharField(max_length=30, default='', verbose_name=u'工单申请人')
    auditor = models.CharField(max_length=4096, default='', verbose_name=u'工单审核人')
    reviewer = models.CharField(max_length=4096, default='', verbose_name=u'工单复核人')
    email_cc = models.CharField(max_length=4096, default='', verbose_name=u'抄送联系人')
    host = models.CharField(null=False, default='', max_length=128, verbose_name=u'主机')
    port = models.IntegerField(null=False, default=3306, verbose_name=u'端口')
    database = models.CharField(null=False, default='', max_length=30, verbose_name=u'库名')
    progress = models.CharField(max_length=30, default='0', choices=progressChoice, verbose_name=u'进度')
    version = models.ForeignKey('OnlineVersion', null=True, blank=True, on_delete=models.CASCADE,
                                verbose_name=u'上线版本号')
    close_info = models.CharField(max_length=1024, default='', verbose_name=u'关闭记录的详情')
    contents = models.TextField(default='', verbose_name=u'工单内容')
    created_at = models.DateTimeField(auto_now_add=True, verbose_name=u'创建时间')
    updated_at = models.DateTimeField(auto_now=True, verbose_name=u'更新时间')

    def __str__(self):
        return self.title

    def applicant_avatar_file(self):
        # 申请人头像
        return UserAccounts.objects.get(username=self.applicant).avatar_file

    class Meta:
        verbose_name = u'工单记录'
        verbose_name_plural = verbose_name

        default_permissions = ()
        app_label = 'orders'
        db_table = 'auditsql_orders'


# 执行任务的进度
TaskProgress = (
    ('0', u'未执行'),
    ('1', u'已完成'),
    ('2', u'处理中'),
    ('3', u'失败'),
    ('4', u'异常')
)


class OrdersTasks(models.Model):
    id = models.AutoField(primary_key=True, verbose_name=u'主键ID')
    taskid = models.CharField(null=False, max_length=128, verbose_name=u'任务ID')
    applicant = models.CharField(max_length=30, null=False, verbose_name=u'工单申请人')
    executor = models.CharField(max_length=30, null=False, default='', verbose_name=u'工单执行人')
    order = models.ForeignKey(Orders, to_field='id', null=False, on_delete=models.CASCADE, verbose_name=u'工单标题')
    envi = models.ForeignKey('orders.SysEnvironment', to_field='envi_id', null=False, on_delete=models.CASCADE,
                             verbose_name=u'环境')
    host = models.CharField(null=False, max_length=128, verbose_name=u'主机')
    database = models.CharField(null=False, max_length=80, verbose_name=u'库名')
    port = models.IntegerField(null=False, default=3306, verbose_name=u'端口')
    sql = models.TextField(default='', verbose_name=u'SQL')
    sql_type = models.CharField(max_length=10, default='DML', choices=sqlTypeChoice, verbose_name=u'SQL类型')
    is_ghost = models.IntegerField(choices=((0, '否'), (1, '是')), default=0, verbose_name=u'是否启用ghost改表')
    ghost_pid = models.IntegerField(null=False, default=0, verbose_name=u'ghost进程pid')
    task_progress = models.CharField(max_length=10, default='0', choices=TaskProgress, verbose_name=u'任务进度')
    affected_row = models.IntegerField(null=False, default=0, verbose_name=u'影响行数')
    task_execlog = models.TextField(default='', verbose_name=u'任务执行的日志')
    rollback_sql = models.TextField(default='', verbose_name=u'回滚的SQL')
    file_format = models.CharField(max_length=30, choices=fileFormatChoice, verbose_name=u'导出的文件格式')
    consume_time = models.CharField(max_length=1024, null=False, default='0.00', verbose_name=u'任务耗时，单位s')
    execition_time = models.DateTimeField(auto_now=True, verbose_name=u'工单执行时间')
    created_time = models.DateTimeField(auto_now_add=True, verbose_name=u'生成时间')

    def export_file_format(self):
        if self.sql_type in ['DML', 'DDL', 'OPS']:
            return None
        else:
            return self.file_format

    def sql_content(self):
        return self.sql[:32]

    export_file_format.short_description = '文件格式'
    sql_content.short_description = 'SQL'

    class Meta:
        verbose_name = u'工单子任务'
        verbose_name_plural = verbose_name

        default_permissions = ()
        app_label = 'orders'
        db_table = 'auditsql_orders_tasks'


class OnlineVersion(models.Model):
    id = models.AutoField(primary_key=True, verbose_name=u'主键id')
    username = models.CharField(default='', null=False, max_length=128, verbose_name=u'创建用户')
    version = models.CharField(default='', null=False, max_length=128, verbose_name=u'版本号')
    expire_time = models.DateTimeField(default='2000-11-01 01:01:01', verbose_name=u'截止上线日期')
    is_deleted = models.CharField(max_length=2, choices=(('0', '未删除'), ('1', '已删除')), default='0',
                                  verbose_name=u'标记为是否删除')
    created_at = models.DateTimeField(auto_now_add=True, verbose_name=u'创建时间')
    updated_at = models.DateTimeField(auto_now=True, verbose_name=u'更新时间')

    def __str__(self):
        return self.version

    class Meta:
        verbose_name = u'上线版本号'
        verbose_name_plural = verbose_name

        default_permissions = ()
        app_label = 'orders'
        db_table = 'auditsql_online_version'
        unique_together = ('version',)


class OrderReply(models.Model):
    id = models.AutoField(primary_key=True, verbose_name=u'主键')
    reply = models.ForeignKey(Orders, on_delete=models.CASCADE, null=False, default='')
    user = models.ForeignKey(UserAccounts, on_delete=models.CASCADE, null=False, default='')
    reply_contents = models.TextField(default='', verbose_name=u'回复内容')
    created_at = models.DateTimeField(auto_now_add=True, verbose_name=u'回复时间')

    class Meta:
        verbose_name = u'工单回复'
        verbose_name_plural = verbose_name

        default_permissions = ()
        app_label = 'orders'
        db_table = 'auditsql_reply'

    def reply_id(self):
        return self.reply.id

    def user_id(self):
        return self.user.uid


class ExportFiles(models.Model):
    id = models.AutoField(primary_key=True, verbose_name=u'主键id')
    task = models.ForeignKey(OrdersTasks, on_delete=models.CASCADE, null=False, blank=True,
                             verbose_name=u'关联子任务的ID')
    file_name = models.CharField(max_length=256, default='', verbose_name=u'文件名')
    file_size = models.IntegerField(default=0, verbose_name=u'文件大小，单位B')
    files = models.FileField(upload_to='files/%Y/%m/%d/')
    encryption_key = models.CharField(max_length=128, default='', verbose_name='加密密钥')
    content_type = models.CharField(max_length=100, default='', verbose_name=u'文件的类型')

    def size(self):
        return ''.join((str(round(self.file_size / 1024 / 1024, 2)), 'MB')) if self.file_size > 1048576 else ''.join(
            (str(round(self.file_size / 1024, 2)), 'KB'))

    class Meta:
        verbose_name = u'sql导出excel表'
        verbose_name_plural = verbose_name

        default_permissions = ()
        db_table = 'auditsql_export_files'
