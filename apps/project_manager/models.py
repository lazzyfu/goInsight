from django.db import models

from user_manager.models import UserAccount

# Create your models here.

# 审核进度选择
progress_choices = (
    ('0', u'待批准'),
    ('1', u'未批准'),
    ('2', u'已批准'),
    ('3', u'处理中'),
    ('4', u'已完成'),
    ('5', u'已关闭'),
    ('6', u'已勾住')
)

# 操作类型选择
operate_type_choice = (
    ('DDL', u'数据库定义语言'),
    ('DML', u'数据库操纵语言')
)


class AuditContents(models.Model):
    id = models.AutoField(primary_key=True, verbose_name=u'主键id')
    title = models.CharField(max_length=100, verbose_name=u'标题')
    url = models.CharField(max_length=1024, default='', null=False, verbose_name=u'上线的Confluence链接')
    operate_type = models.CharField(max_length=5, default='DML', choices=operate_type_choice,
                                    verbose_name=u'操作类型: DDL or DML')
    envi_desc = models.SmallIntegerField(null=False, default=0, verbose_name=u'环境：0：测试、1：staging、2：生产、3：线下其他环境')
    proposer = models.CharField(max_length=30, default='', verbose_name=u'申请人， 一般为开发或者产品，存储username')
    operator = models.CharField(max_length=30, default='', verbose_name=u'审核DBA')
    operate_time = models.DateTimeField(auto_now_add=True, verbose_name=u'审核时间')
    host = models.CharField(null=False, default='', max_length=30, verbose_name=u'操作数据库主机')
    database = models.CharField(null=False, default='', max_length=80, verbose_name=u'操作数据库')
    port = models.IntegerField(null=False, default=3306, verbose_name=u'端口')
    progress = models.CharField(max_length=10, default='0', choices=progress_choices, verbose_name=u'任务进度')
    remark = models.SmallIntegerField(default=0, null=False, verbose_name=u'工单备注, 0: 周三上线，1：紧急上线，2：数据修复')
    tasks = models.CharField(max_length=256, default='', verbose_name=u'部署步骤版本')
    close_user = models.CharField(max_length=30, default='', verbose_name=u'关闭记录的用户')
    close_reason = models.CharField(max_length=1024, default='', verbose_name=u'关闭原因')
    close_time = models.DateTimeField(auto_now_add=True, verbose_name=u'关闭时间')
    contents = models.TextField(default='', verbose_name=u'提交的内容')
    created_at = models.DateTimeField(auto_now_add=True, verbose_name=u'创建时间')
    updated_at = models.DateTimeField(auto_now=True, verbose_name=u'更新时间')

    def __str__(self):
        return self.title

    def proposer_avatar_file(self):
        return UserAccount.objects.get(username=self.proposer).avatar_file

    class Meta:
        verbose_name = u'审核内容'
        verbose_name_plural = verbose_name

        db_table = 'auditsql_work_order'


class AuditTasks(models.Model):
    id = models.AutoField(primary_key=True, verbose_name=u'主键id')
    username = models.CharField(default='', null=False, max_length=128, verbose_name=u'创建用户')
    tasks = models.CharField(default='', null=False, max_length=128, verbose_name=u'任务名')
    expire_time = models.DateTimeField(default='2000-11-01 01:01:01', verbose_name=u'任务截止上线日期')
    created_at = models.DateTimeField(auto_now_add=True, verbose_name=u'创建时间')

    def __str__(self):
        return self.tasks

    class Meta:
        verbose_name = u'审核任务'
        verbose_name_plural = verbose_name

        db_table = 'auditsql_audit_tasks'
        unique_together = ('tasks',)


export_progress_choices = (
    ('0', u'未执行'),
    ('1', u'导出中'),
    ('2', u'已生成')
)


class OlAuditContentsReply(models.Model):
    id = models.AutoField(primary_key=True, verbose_name=u'主键')
    reply = models.ForeignKey(AuditContents, on_delete=models.CASCADE, null=False, default='')
    user = models.ForeignKey(UserAccount, on_delete=models.CASCADE, null=False, default='')
    reply_contents = models.TextField(default='', verbose_name=u'回复内容')
    created_at = models.DateTimeField(auto_now_add=True, verbose_name=u'回复时间')

    class Meta:
        verbose_name = u'线上审核回复表'
        verbose_name_plural = verbose_name

        default_permissions = ()
        db_table = 'auditsql_work_order_reply'

    def reply_id(self):
        return self.reply.id

    def user_id(self):
        return self.user.uid


class MonitorSchema(models.Model):
    table_schema = models.CharField(null=False, max_length=512)
    table_name = models.CharField(null=False, max_length=512)
    table_stru = models.TextField(null=False, default='')
    md5_sum = models.CharField(null=False, max_length=256)

    class Meta:
        verbose_name = u'监控表结构变更表'
        verbose_name_plural = verbose_name
        permissions = ()
        db_table = "auditsql_monitor_schema"


exec_progress = (
    ('0', u'未执行'),
    ('1', u'已完成'),
    ('2', u'处理中'),
    ('3', u'回滚中'),
    ('4', u'已回滚'),
    ('5', u'失败'),
    ('6', u'异常')
)


class IncepMakeExecTask(models.Model):
    id = models.AutoField(primary_key=True, verbose_name=u'主键id')
    uid = models.IntegerField(null=False, default=0, verbose_name=u'操作用户uid')
    user = models.CharField(max_length=30, null=False, verbose_name=u'操作用户')
    taskid = models.CharField(null=False, max_length=128, verbose_name=u'任务号')
    related_id = models.IntegerField(null=False, default=0, verbose_name=u'关联AuditContents的主键id')
    envi_desc = models.SmallIntegerField(null=False, default=0, verbose_name=u'环境：0：测试、1：staging、2：生产、3：线下其他环境')
    dst_host = models.CharField(null=False, max_length=30, verbose_name=u'操作目标数据库主机')
    dst_database = models.CharField(null=False, max_length=80, verbose_name=u'操作目标数据库')
    dst_port = models.IntegerField(null=False, default=3306, verbose_name=u'端口')
    sql_content = models.TextField(verbose_name=u'执行的SQL', default='')
    type = models.CharField(max_length=5, default='', choices=(('DDL', u'数据库定义语言'), ('DML', u'数据库操纵语言')))
    sqlsha1 = models.CharField(null=False, max_length=120, default='', verbose_name=u'sqlsha1')
    rollback_sqlsha1 = models.CharField(null=False, max_length=120, default='', verbose_name=u'rollback sqlsha1')
    celery_task_id = models.CharField(null=False, max_length=256, default='', verbose_name=u'celery执行任务ID')
    exec_status = models.CharField(max_length=10, default='0', choices=exec_progress, verbose_name=u'执行进度')
    sequence = models.CharField(null=False, default='', max_length=1024, verbose_name=u'备份记录id，inception生成的sequence')
    affected_row = models.IntegerField(null=False, default=0, verbose_name=u'预计影响行数')
    backup_dbname = models.CharField(null=False, max_length=1024, default='', verbose_name=u'inception生成的备份的库名')
    exec_log = models.TextField(verbose_name=u'执行成功的记录', default='')
    make_time = models.DateTimeField(auto_now_add=True, verbose_name=u'生成时间')

    class Meta:
        verbose_name = u'生成Inception执行任务'
        verbose_name_plural = verbose_name

        default_permissions = ()
        db_table = 'auditsql_incep_tasks'


class DomainName(models.Model):
    id = models.AutoField(primary_key=True, verbose_name=u'主键id')
    domain_name = models.CharField(max_length=256, default='', null=False, verbose_name=u'域名地址')

    class Meta:
        verbose_name = u'域名地址'
        verbose_name_plural = verbose_name

        default_permissions = ()
        db_table = 'auditsql_domain_name'


class Webhook(models.Model):
    id = models.AutoField(primary_key=True, verbose_name=u'主键id')
    webhook_addr = models.CharField(max_length=256, default='', null=False, verbose_name=u'webhook地址')

    class Meta:
        verbose_name = u'钉钉机器人'
        verbose_name_plural = verbose_name

        default_permissions = ()
        db_table = 'auditsql_webhook'
