from django.db import models

# Create your models here.
from UserManager.models import Groups, Contacts, UserAccount


class InceptionHostConfig(models.Model):
    id = models.AutoField(primary_key=True, verbose_name=u'主键id')
    user = models.CharField(max_length=30, null=False, verbose_name=u'用户名')
    password = models.CharField(max_length=128, null=False, verbose_name=u'密码')
    host = models.CharField(max_length=32, null=False, verbose_name=u'ip地址')
    port = models.IntegerField(null=False, verbose_name=u'端口')
    comment = models.CharField(max_length=128, verbose_name=u'主机描述')

    class Meta:
        verbose_name = u'inception连接数据库配置'
        verbose_name_plural = verbose_name

        default_permissions = ()
        db_table = 'sqlaudit_inception_hostconfig'


class InceptionSqlOperateRecord(models.Model):
    id = models.AutoField(primary_key=True, verbose_name=u'主键id')
    op_uid = models.IntegerField(null=False, default=0, verbose_name=u'操作用户uid')
    op_user = models.CharField(max_length=30, null=False, verbose_name=u'操作用户')
    step_id = models.IntegerField(null=False, default=1, verbose_name=u'inception输出顺序ID')
    workid = models.CharField(null=False, max_length=128, verbose_name=u'工单号')
    dst_host = models.CharField(null=False, max_length=30, verbose_name=u'操作目标数据库主机')
    dst_database = models.CharField(null=False, max_length=80, verbose_name=u'操作目标数据库')
    stage = models.CharField(null=False, max_length=128, default='', verbose_name=u'inception检测阶段')
    stagestatus = models.CharField(null=False, default='', max_length=1024, verbose_name=u'inception检测阶段状态')
    errlevel = models.IntegerField(null=False, default=0, verbose_name=u'inception输出错误级别')
    errormessage = models.TextField(null=False, default='', verbose_name=u'inception输出错误信息')
    op_sql = models.TextField(verbose_name=u'执行的SQL', default='')
    affected_rows = models.IntegerField(null=False, default=0, verbose_name=u'影响的行数')
    sequence = models.CharField(null=False, default='', max_length=1024, verbose_name=u'备份记录id，inception生成的sequence')
    backup_dbname = models.CharField(null=False, max_length=1024, default='', verbose_name=u'inception生成的备份的库名')
    execute_time = models.CharField(null=False, default='0.000', max_length=128, verbose_name=u'inception执行时间')
    op_time = models.DateTimeField(auto_now_add=True, verbose_name=u'用户操作时间')

    class Meta:
        verbose_name = u'inception操作结果记录表'
        verbose_name_plural = verbose_name

        default_permissions = ()
        db_table = 'sqlaudit_inception_sql_operate_record'


progress_status_choices = (
    ('0', u'待批准'),
    ('1', u'未批准'),
    ('2', u'已批准'),
    ('3', u'处理中'),
    ('4', u'已完成'),
    ('5', u'已关闭')
)


class OnlineAuditContents(models.Model):
    id = models.AutoField(primary_key=True, verbose_name=u'主键id')
    group = models.ForeignKey(Groups, on_delete=models.CASCADE, verbose_name=u'项目组id')
    title = models.CharField(max_length=100, verbose_name=u'标题')
    remark = models.CharField(default='', max_length=30, verbose_name=u'备注')
    proposer = models.CharField(max_length=30, default='', verbose_name=u'申请人')
    verifier = models.CharField(max_length=30, default='', verbose_name=u'批准人')
    operate_dba = models.CharField(max_length=30, default='', verbose_name=u'执行dba')
    email_cc = models.CharField(max_length=1024, default='', verbose_name=u'抄送人')
    progress_status = models.CharField(max_length=10, default='0', choices=progress_status_choices, verbose_name=u'进度')
    contents = models.TextField(default='', verbose_name=u'提交的内容')
    created_at = models.DateTimeField(auto_now_add=True, verbose_name=u'创建时间')
    updated_at = models.DateTimeField(auto_now=True, verbose_name=u'更新时间')

    def __str__(self):
        return self.title

    def email_cc_list(self):
        return '\n'.join(
            Contacts.objects.filter(contact_id__in=self.email_cc.split(',')).values_list('email', flat=True))

    def items_id(self):
        """
        实例化items，并返回，否则无法操作
        :return: items_id
        """
        return self.items.items_id

    def remark_list(self):
        return self.remark.split(',')

    def proposer_info(self):
        return UserAccount.objects.get(username=self.proposer)

    class Meta:
        verbose_name = u'线上操作审计表'
        verbose_name_plural = verbose_name

        default_permissions = ()
        db_table = 'sqlaudit_audit_contents'
        unique_together = ('title',)

        permissions = (
            ('leader_verify', u'批准权限'),
        )

class Remark(models.Model):
    id = models.AutoField(primary_key=True, verbose_name=u'主键')
    remark = models.CharField(default='', max_length=30, verbose_name=u'备注')
    created_at = models.DateTimeField(auto_now_add=True, verbose_name=u'创建时间')
    updated_at = models.DateTimeField(auto_now=True, verbose_name=u'更新时间')

    class Meta:
        verbose_name = u'线上审计备注表'
        verbose_name_plural = verbose_name

        default_permissions = ()
        db_table = 'sqlaudit_audit_remark'