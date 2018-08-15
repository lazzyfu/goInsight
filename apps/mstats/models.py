from django.db import models

# Create your models here.
from user_manager.models import UserAccount


class MySQLQueryLog(models.Model):
    id = models.AutoField(primary_key=True, verbose_name=u'主键id')
    user = models.CharField(max_length=30, null=False, verbose_name=u'用户名')
    host = models.CharField(max_length=32, null=False, verbose_name=u'目标数据库地址')
    database = models.CharField(max_length=32, null=False, verbose_name=u'目标数据库')
    envi = models.SmallIntegerField(null=False, default=1, verbose_name=u'0：线上环境， 1：线下环境')
    query_sql = models.TextField(null=False, default='', verbose_name=u'查询SQL')
    query_time = models.CharField(null=False, default='', max_length=128, verbose_name=u'查询时间，单位s')
    query_status = models.CharField(max_length=2048, default='', verbose_name=u'查询状态，成功或失败的原因')
    affect_rows = models.IntegerField(default=0, null=False, verbose_name=u'影响行数')
    created_at = models.DateTimeField(auto_now_add=True, verbose_name=u'查询时间')

    class Meta:
        verbose_name = u'mysql查询记录日志'
        verbose_name_plural = verbose_name

        default_permissions = ()
        db_table = 'auditsql_mysql_query_log'


class MysqlSchemaInfo(models.Model):
    id = models.AutoField(primary_key=True, verbose_name=u'主键id')
    user = models.CharField(max_length=30, null=False, verbose_name=u'用户名')
    password = models.CharField(max_length=30, null=False, verbose_name=u'密码')
    host = models.CharField(max_length=32, null=False, verbose_name=u'地址')
    port = models.IntegerField(null=False, default=3306, verbose_name=u'端口')
    schema = models.CharField(null=False, max_length=64, default='', verbose_name=u'schema信息')
    envi = models.SmallIntegerField(null=False, default=1, verbose_name=u'0: 线下其他环境，1：test环境，2：Staging环境 3：生产环境')
    is_master = models.SmallIntegerField(null=False, default=1, verbose_name=u'0：从库（用于查询）， 1：主库（用于审核）')
    schema_join = models.CharField(null=False, default='', max_length=128, unique=True,
                                   verbose_name=u'host、port、schema的组合, 关联MysqlSchemaGrant的schema')
    comment = models.CharField(max_length=128, null=True, verbose_name=u'主机描述')
    created_at = models.DateTimeField(auto_now_add=True, verbose_name=u'创建时间')
    updated_at = models.DateTimeField(auto_now=True, verbose_name=u'更新时间')

    def __str__(self):
        return self.schema

    class Meta:
        verbose_name = u'所有MySQL集群节点的schema汇总'
        verbose_name_plural = verbose_name

        default_permissions = ()
        db_table = 'auditsql_mysql_schema'
        unique_together = (('host', 'port', 'schema'),)


class MysqlSchemaGrant(models.Model):
    id = models.AutoField(primary_key=True, verbose_name=u'主键id')
    user = models.ForeignKey(UserAccount, null=True, blank=True, on_delete=models.SET_NULL)
    schema = models.ForeignKey(MysqlSchemaInfo, to_field='schema_join', on_delete=models.CASCADE)
    created_at = models.DateTimeField(auto_now_add=True, verbose_name=u'创建时间')
    updated_at = models.DateTimeField(auto_now=True, verbose_name=u'更新时间')

    class Meta:
        verbose_name = u'生产环境数据库授权'
        verbose_name_plural = verbose_name

        default_permissions = ()
        db_table = 'auditsql_mysql_schema_grant'
        unique_together = (('user', 'schema'),)


class MysqlSlowLog(models.Model):
    id = models.AutoField(primary_key=True, verbose_name=u'主键id')
    hostname = models.CharField(null=False, default='', max_length=128, verbose_name=u'慢查询主机')
    version = models.CharField(null=False, default='', max_length=128, verbose_name=u'数据库版本')
    qps = models.CharField(null=False, default='0.00', max_length=128, verbose_name=u'区间内每秒平均查询数')
    cnt = models.IntegerField(null=False, default=0, verbose_name=u'区间内执行总数')
    avg = models.CharField(null=False, default='', max_length=128, verbose_name=u'区间内查询的平均响应时间')
    fingerprint = models.TextField(verbose_name=u'SQL指纹')
    md5sum = models.CharField(null=False, max_length=32, default='', verbose_name=u'MD5值')
    is_pull = models.SmallIntegerField(null=False, default=0, verbose_name=u'是否钉钉推送：0：推送，1：不推送')
    timerange = models.CharField(null=False, default='', max_length=128, verbose_name=u'时间范围')
    created_at = models.DateTimeField(auto_now_add=True, verbose_name=u'创建时间')

    def __str__(self):
        return self.hostname

    class Meta:
        verbose_name = u'mysql慢查询统计表'
        verbose_name_plural = verbose_name

        default_permissions = ()
        db_table = 'auditsql_mysql_slow_log'
        index_together = ('md5sum',)


class WebShellInfo(models.Model):
    id = models.AutoField(primary_key=True, verbose_name=u'主键id')
    command = models.CharField(max_length=1024, null=False, default='', verbose_name=u'用户命令')
    envi = models.SmallIntegerField(null=False, default=1, verbose_name=u'0: 线下环境，1: 生产环境')
    comment = models.CharField(max_length=128, null=True, verbose_name=u'命令描述')
    created_at = models.DateTimeField(auto_now_add=True, verbose_name=u'创建时间')
    updated_at = models.DateTimeField(auto_now=True, verbose_name=u'更新时间')

    def __str__(self):
        return self.comment

    class Meta:
        verbose_name = u'mongodb和redis'
        verbose_name_plural = verbose_name

        default_permissions = ()
        db_table = 'auditsql_web_shell'
        unique_together = (('comment',),)


class WebShellGrant(models.Model):
    id = models.AutoField(primary_key=True, verbose_name=u'主键id')
    user = models.ForeignKey(UserAccount, null=True, blank=True, on_delete=models.SET_NULL)
    shell = models.ForeignKey(WebShellInfo, on_delete=models.CASCADE)
    created_at = models.DateTimeField(auto_now_add=True, verbose_name=u'创建时间')
    updated_at = models.DateTimeField(auto_now=True, verbose_name=u'更新时间')

    class Meta:
        verbose_name = u'web shell授权表'
        verbose_name_plural = verbose_name

        default_permissions = ()
        db_table = 'auditsql_web_shell_grant'
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
        db_table = 'auditsql_web_shell_oplog'


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
        db_table = 'auditsql_deadlock_command'


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
        db_table = "dbaudit_deadlocks_records"


class MySQLConfigSource(models.Model):
    id = models.AutoField(primary_key=True, verbose_name=u'主键id')
    host = models.CharField(max_length=32, null=False, verbose_name=u'地址')
    port = models.IntegerField(null=False, default=3306, verbose_name=u'端口')
    envi = models.SmallIntegerField(null=False, default=1, verbose_name=u'1：test环境，2：Staging环境 3：生产环境')
    is_master = models.SmallIntegerField(null=False, default=1, verbose_name=u'0：从库（用于查询）， 1：主库（用于审核）')
    comment = models.CharField(max_length=128, null=True, verbose_name=u'主机描述')
    created_at = models.DateTimeField(auto_now_add=True, verbose_name=u'创建时间')
    updated_at = models.DateTimeField(auto_now=True, verbose_name=u'更新时间')

    def __str__(self):
        return self.host

    class Meta:
        verbose_name = u'数据库主机配置'
        verbose_name_plural = verbose_name

        default_permissions = ()
        db_table = 'auditsql_mysql_config'
        unique_together = (('host', 'port'),)
