from django.db import models


# Create your models here.
from sqlorders.models import envi_choice
from users.models import UserAccounts


class WebShellInfo(models.Model):
    id = models.AutoField(primary_key=True, verbose_name=u'主键id')
    command = models.CharField(max_length=1024, null=False, default='', verbose_name=u'命令')
    comment = models.CharField(max_length=128, null=True, verbose_name=u'描述')
    envi_id = models.IntegerField(choices=envi_choice, verbose_name=u'环境')
    created_at = models.DateTimeField(auto_now_add=True, verbose_name=u'创建时间')
    updated_at = models.DateTimeField(auto_now=True, verbose_name=u'更新时间')

    def __str__(self):
        return self.comment

    class Meta:
        verbose_name = u'webshell'
        verbose_name_plural = verbose_name

        default_permissions = ()
        db_table = 'sqlaudit_web_shell'
        unique_together = (('comment',),)


class WebShellGrant(models.Model):
    id = models.AutoField(primary_key=True, verbose_name=u'主键id')
    user = models.ForeignKey(UserAccounts, null=True, blank=True, on_delete=models.SET_NULL)
    shell = models.ForeignKey(WebShellInfo, on_delete=models.CASCADE)
    created_at = models.DateTimeField(auto_now_add=True, verbose_name=u'创建时间')
    updated_at = models.DateTimeField(auto_now=True, verbose_name=u'更新时间')

    class Meta:
        verbose_name = u'webshell授权'
        verbose_name_plural = verbose_name

        default_permissions = ()
        db_table = 'sqlaudit_web_shell_grant'
        unique_together = (('user', 'shell'),)


class WebShellOpLog(models.Model):
    id = models.AutoField(primary_key=True, verbose_name=u'主键id')
    user = models.CharField(max_length=128, null=False, default='', verbose_name=u'操作用户')
    session_id = models.CharField(max_length=128, null=False, default='', verbose_name=u'会话id')
    op_cmd = models.CharField(max_length=4096, null=True, default='', verbose_name=u'操作的命令')
    created_at = models.DateTimeField(auto_now_add=True, verbose_name=u'创建时间')
    updated_at = models.DateTimeField(auto_now=True, verbose_name=u'更新时间')

    class Meta:
        verbose_name = u'web shell操作记录表'
        verbose_name_plural = verbose_name

        default_permissions = ()
        db_table = 'sqlaudit_web_shell_oplog'


class DeadlockCommand(models.Model):
    id = models.AutoField(primary_key=True, verbose_name=u'主键id')
    command = models.CharField(max_length=1024, null=False, default='', verbose_name=u'pt-deadlock-logger命令')
    schema_id = models.CharField(null=False, default='', max_length=128, verbose_name=u'关联MysqlSchemaInfo的id')
    created_at = models.DateTimeField(auto_now_add=True, verbose_name=u'创建时间')
    updated_at = models.DateTimeField(auto_now=True, verbose_name=u'更新时间')

    class Meta:
        verbose_name = u'死锁命令表'
        verbose_name_plural = verbose_name

        default_permissions = ()
        db_table = 'sqlaudit_deadlock_command'


class DeadlockRecord(models.Model):
    id = models.AutoField(primary_key=True, verbose_name=u'主键id')
    server = models.CharField(null=False, max_length=20)
    ts = models.DateTimeField(null=False, auto_now=True)
    thread = models.IntegerField(null=False)
    txn_id = models.CharField(max_length=1024, null=False)
    txn_time = models.PositiveSmallIntegerField(null=False)
    user = models.CharField(max_length=16, null=False)
    hostname = models.CharField(max_length=20, null=False)
    ip = models.CharField(max_length=15, null=False)
    db = models.CharField(max_length=64, null=False)
    tbl = models.CharField(max_length=64, null=False)
    idx = models.CharField(max_length=64, null=False)
    lock_type = models.CharField(max_length=16, null=False)
    lock_mode = models.CharField(max_length=1, null=False)
    wait_hold = models.CharField(max_length=1, null=False)
    victim = models.SmallIntegerField(null=False)
    query = models.TextField()
    is_pull = models.SmallIntegerField(null=False, default=0, verbose_name=u'是否已推送，0：未推送，1：已推送')

    class Meta:
        verbose_name = u'死锁采集记录表,pt-deadlock-logger工具采集'
        verbose_name_plural = verbose_name

        unique_together = ('server', 'ts', 'thread')

        permissions = ()
        db_table = "sqlaudit_deadlocks_records"
