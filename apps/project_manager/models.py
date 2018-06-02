from django.db import models
# Create your models here.
from django.db.models import F

from user_manager.models import Groups, Contacts, UserAccount, GroupsDetail

# 审核进度选择
progress_choices = (
    ('0', u'待批准'),
    ('1', u'未批准'),
    ('2', u'已批准'),
    ('3', u'处理中'),
    ('4', u'已完成'),
    ('5', u'已关闭')
)

# 操作类型选择
operate_type_choice = (
    ('DDL', u'数据库定义语言'),
    ('DML', u'数据库操纵语言')
)


class AuditContents(models.Model):
    id = models.AutoField(primary_key=True, verbose_name=u'主键id')
    title = models.CharField(max_length=100, verbose_name=u'标题')
    group = models.ForeignKey(Groups, on_delete=models.CASCADE, verbose_name=u'关联项目组表id')
    operate_type = models.CharField(max_length=5, default='DML', choices=operate_type_choice,
                                    verbose_name=u'操作类型: DDL or DML')
    proposer = models.CharField(max_length=30, default='', verbose_name=u'申请人， 一般为开发或者产品，存储username')
    verifier = models.CharField(max_length=30, default='', verbose_name=u'批准人，一般为项目经理或Leader， 存储username')
    operator = models.CharField(max_length=30, default='', verbose_name=u'执行人，一般为DBA， 存储username')
    email_cc = models.CharField(max_length=1024, default='', verbose_name=u'抄送人， 存储contact_id，以逗号分隔')
    host = models.CharField(null=False, default='', max_length=30, verbose_name=u'操作数据库主机')
    database = models.CharField(null=False, default='', max_length=80, verbose_name=u'操作数据库')
    progress = models.CharField(max_length=10, default='0', choices=progress_choices, verbose_name=u'任务进度')
    verifier_time = models.DateTimeField(auto_now_add=True, verbose_name=u'审批时间')
    operate_time = models.DateTimeField(auto_now_add=True, verbose_name=u'执行的时间')
    close_user = models.CharField(max_length=30, default='', verbose_name=u'关闭记录的用户')
    close_reason = models.CharField(max_length=1024, default='', verbose_name=u'关闭原因')
    close_time = models.DateTimeField(auto_now_add=True, verbose_name=u'关闭时间')
    created_at = models.DateTimeField(auto_now_add=True, verbose_name=u'创建时间')
    updated_at = models.DateTimeField(auto_now=True, verbose_name=u'更新时间')

    def __str__(self):
        return self.title

    def proposer_avatar_file(self):
        return UserAccount.objects.get(username=self.proposer).avatar_file

    def email_cc_list(self):
        return '\n'.join(
            Contacts.objects.filter(contact_id__in=self.email_cc.split(',')).values_list('contact_email', flat=True))

    class Meta:
        verbose_name = u'审核内容'
        verbose_name_plural = verbose_name

        db_table = 'auditsql_audit_contents'
        unique_together = ('title',)


class OlAuditDetail(models.Model):
    id = models.AutoField(primary_key=True, verbose_name=u'主键id')
    ol = models.ForeignKey(AuditContents, on_delete=models.CASCADE, verbose_name=u'关联审核内容表id')
    contents = models.TextField(default='', verbose_name=u'提交的内容')
    created_at = models.DateTimeField(auto_now_add=True, verbose_name=u'创建时间')
    updated_at = models.DateTimeField(auto_now=True, verbose_name=u'更新时间')

    class Meta:
        verbose_name = u'线上审核内容详情表'
        verbose_name_plural = verbose_name

        db_table = 'auditsql_ol_audit_detail'


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
        db_table = 'auditsql_ol_audit_reply'

    def reply_id(self):
        return self.reply.id

    def user_id(self):
        return self.user.uid


class InceptionHostConfig(models.Model):
    id = models.AutoField(primary_key=True, verbose_name=u'主键id')
    user = models.CharField(max_length=30, null=False, verbose_name=u'用户名')
    password = models.CharField(max_length=128, null=False, verbose_name=u'密码')
    host = models.CharField(max_length=32, null=False, verbose_name=u'ip地址')
    port = models.IntegerField(null=False, verbose_name=u'端口')
    type = models.IntegerField(null=False, default=0, verbose_name=u'0:线下数据库，1:线上数据库')
    purpose = models.CharField(default='0', max_length=2, choices=(('0', u'审核'), ('1', u'查询')), verbose_name=u'用途')
    is_enable = models.IntegerField(null=False, default=0, verbose_name=u'0:启用，1：禁用')
    protection_user = models.TextField(default='root', null=False, verbose_name=u'被保护的数据库账号， 以逗号分隔')
    comment = models.CharField(max_length=128, verbose_name=u'主机描述')

    def group_name(self):
        group = InceptionHostConfigDetail.objects.annotate(group_name=F('group__group_name')).filter(
            config__id=self.id).values_list(
            'group_name',
            flat=True)
        return ', '.join(group)

    class Meta:
        verbose_name = u'数据库连接账号'
        verbose_name_plural = verbose_name

        unique_together = [('host', 'port'), ('is_enable', 'comment')]

        default_permissions = ()
        db_table = 'auditsql_inception_hostconfig'


class InceptionHostConfigDetail(models.Model):
    """
    inception主机分组
    """
    id = models.AutoField(primary_key=True, verbose_name=u'主键')
    config = models.ForeignKey(InceptionHostConfig, on_delete=models.CASCADE, null=False)
    group = models.ForeignKey(Groups, on_delete=models.CASCADE, null=False)
    created_at = models.DateTimeField(auto_now_add=True, verbose_name=u'创建时间')
    updated_at = models.DateTimeField(auto_now=True, verbose_name=u'更新时间')

    class Meta:
        verbose_name = u'inception主机分组'
        verbose_name_plural = verbose_name

        default_permissions = ()
        db_table = 'auditsql_inception_hostconfig_detail'


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
)


class IncepMakeExecTask(models.Model):
    id = models.AutoField(primary_key=True, verbose_name=u'主键id')
    uid = models.IntegerField(null=False, default=0, verbose_name=u'操作用户uid')
    user = models.CharField(max_length=30, null=False, verbose_name=u'操作用户')
    group = models.ForeignKey(Groups, on_delete=models.CASCADE, verbose_name=u'项目组id')
    taskid = models.CharField(null=False, max_length=128, verbose_name=u'任务号')
    related_id = models.IntegerField(null=False, default=0, verbose_name=u'关联AuditContents的主键id')
    category = models.CharField(null=False, max_length=2, default='0', choices=(('0', u'线下任务'), ('1', u'线上任务')),
                                verbose_name=u'任务分类')
    dst_host = models.CharField(null=False, max_length=30, verbose_name=u'操作目标数据库主机')
    dst_database = models.CharField(null=False, max_length=80, verbose_name=u'操作目标数据库')
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
