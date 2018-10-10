import logging

from django.db import models

# Create your models here.

# sql工单环境定义
from users.models import UserAccounts

logger = logging.getLogger('django')


class SqlOrdersEnvironment(models.Model):
    id = models.AutoField(primary_key=True, verbose_name=u'主键id')
    envi_id = models.IntegerField(null=False, default=1, verbose_name=u'ID，起始值：1')
    parent_id = models.IntegerField(null=False, default=0, verbose_name=u'父ID，起始值：0')
    envi_name = models.CharField(max_length=30, default='', null=False, verbose_name=u'环境')
    created_at = models.DateTimeField(auto_now_add=True, verbose_name=u'创建时间')
    updated_at = models.DateTimeField(auto_now=True, verbose_name=u'更新时间')

    class Meta:
        verbose_name = u'环境配置'
        verbose_name_plural = verbose_name

        db_table = 'sqlaudit_sql_order_environment'

        # 建立唯一索引
        unique_together = (('envi_id',), ('parent_id',))


envi_choice = [(x, y) for x, y in list(SqlOrdersEnvironment.objects.all().values_list('envi_id', 'envi_name'))]
# envi_choice = ((0, '1'), (1, '1'))
type_choice = ((0, '数据查询'), (1, 'SQL审核'))


class MysqlConfig(models.Model):
    id = models.AutoField(primary_key=True, verbose_name=u'主键id')
    host = models.CharField(max_length=128, null=False, verbose_name=u'地址')
    port = models.IntegerField(null=False, default=3306, verbose_name=u'端口')
    user = models.CharField(max_length=32, default='', null=False, verbose_name=u'用户名')
    password = models.CharField(max_length=64, default='', null=False, verbose_name=u'密码')
    envi_id = models.IntegerField(choices=envi_choice, verbose_name=u'环境')
    is_master = models.SmallIntegerField(choices=type_choice, verbose_name=u'用途')
    comment = models.CharField(max_length=128, null=True, verbose_name=u'主机描述')
    created_at = models.DateTimeField(auto_now_add=True, verbose_name=u'创建时间')
    updated_at = models.DateTimeField(auto_now=True, verbose_name=u'更新时间')

    def __str__(self):
        return self.host

    class Meta:
        verbose_name = u'数据库主机配置'
        verbose_name_plural = verbose_name

        default_permissions = ()
        db_table = 'sqlaudit_mysql_config'
        unique_together = (('host', 'port'),)


class MysqlSchemas(models.Model):
    id = models.AutoField(primary_key=True, verbose_name=u'主键id')
    user = models.CharField(max_length=30, null=False, verbose_name=u'用户名')
    password = models.CharField(max_length=30, null=False, verbose_name=u'密码')
    host = models.CharField(max_length=128, null=False, verbose_name=u'地址')
    port = models.IntegerField(null=False, default=3306, verbose_name=u'端口')
    schema = models.CharField(null=False, max_length=64, default='', verbose_name=u'schema信息')
    envi_id = models.IntegerField(null=False, verbose_name=u'环境')
    is_master = models.SmallIntegerField(null=False, verbose_name=u'用途, 主库：1， 从库：0')
    schema_join = models.CharField(null=False, max_length=128,
                                   verbose_name=u'host、port、schema的组合')
    comment = models.CharField(max_length=128, null=True, verbose_name=u'主机描述')
    created_at = models.DateTimeField(auto_now_add=True, verbose_name=u'创建时间')
    updated_at = models.DateTimeField(auto_now=True, verbose_name=u'更新时间')

    def __str__(self):
        try:
            envi_name = SqlOrdersEnvironment.objects.get(envi_id=self.envi_id).envi_name
            return '_'.join((envi_name, self.comment, self.schema))
        except Exception as err:
            logger.error(err)
            logger.error('请先配置环境')
            return ''

    class Meta:
        verbose_name = u'MySQL集群汇总库'
        verbose_name_plural = verbose_name

        default_permissions = ()
        db_table = 'sqlaudit_mysql_schemas'
        unique_together = (('host', 'port', 'schema'), ('schema_join',))


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
sql_type_choice = (
    ('DDL', u'DDL'),
    ('DML', u'DML')
)


class SqlOrdersContents(models.Model):
    id = models.AutoField(primary_key=True, verbose_name=u'主键id')
    title = models.CharField(max_length=100, verbose_name=u'标题')
    description = models.CharField(max_length=2048, default='', null=False, verbose_name=u'需求url或者描述性文字')
    sql_type = models.CharField(max_length=5, default='DML', choices=sql_type_choice,
                                verbose_name=u'SQL类型: DDL or DML')
    envi_id = models.IntegerField(choices=envi_choice, verbose_name=u'环境')
    proposer = models.CharField(max_length=30, default='', verbose_name=u'申请人， 一般为开发或者产品，存储username')
    auditor = models.CharField(max_length=30, default='', verbose_name=u'审核人')
    operate_time = models.DateTimeField(auto_now_add=True, verbose_name=u'审核时间')
    email_cc = models.CharField(max_length=4096, default='', verbose_name=u'抄送联系人')
    host = models.CharField(null=False, default='', max_length=128, verbose_name=u'mysql主机')
    port = models.IntegerField(null=False, default=3306, verbose_name=u'端口')
    database = models.CharField(null=False, default='', max_length=80, verbose_name=u'库名')
    progress = models.CharField(max_length=10, default='0', choices=progress_choices, verbose_name=u'任务进度')
    remark = models.CharField(max_length=32, default='', null=False, verbose_name=u'工单备注')
    task_version = models.CharField(max_length=256, default='', verbose_name=u'上线任务版本')
    close_user = models.CharField(max_length=30, default='', verbose_name=u'关闭记录的用户')
    close_reason = models.CharField(max_length=1024, default='', verbose_name=u'关闭原因')
    close_time = models.DateTimeField(auto_now_add=True, verbose_name=u'关闭时间')
    contents = models.TextField(default='', verbose_name=u'提交的内容')
    created_at = models.DateTimeField(auto_now_add=True, verbose_name=u'创建时间')
    updated_at = models.DateTimeField(auto_now=True, verbose_name=u'更新时间')

    def __str__(self):
        return self.title

    def proposer_avatar_file(self):
        return UserAccounts.objects.get(username=self.proposer).avatar_file

    class Meta:
        verbose_name = u'工单记录'
        verbose_name_plural = verbose_name

        db_table = 'sqlaudit_sql_orders_contents'


exec_progress = (
    ('0', u'未执行'),
    ('1', u'已完成'),
    ('2', u'处理中'),
    ('3', u'回滚中'),
    ('4', u'已回滚'),
    ('5', u'失败'),
    ('6', u'异常')
)


class SqlOrdersExecTasks(models.Model):
    id = models.AutoField(primary_key=True, verbose_name=u'主键id')
    uid = models.IntegerField(null=False, default=0, verbose_name=u'操作用户uid')
    user = models.CharField(max_length=30, null=False, verbose_name=u'操作用户')
    taskid = models.CharField(null=False, max_length=128, verbose_name=u'任务号')
    related_id = models.IntegerField(null=False, default=0, verbose_name=u'关联SqlOrdersContents的主键id')
    envi_id = models.IntegerField(choices=envi_choice, verbose_name=u'环境')
    host = models.CharField(null=False, max_length=128, verbose_name=u'操作目标数据库主机')
    database = models.CharField(null=False, max_length=80, verbose_name=u'操作目标数据库')
    port = models.IntegerField(null=False, default=3306, verbose_name=u'端口')
    sql = models.TextField(verbose_name=u'执行的SQL', default='')
    sql_type = models.CharField(max_length=5, default='DML', choices=sql_type_choice,
                                verbose_name=u'SQL类型')
    is_ghost = models.IntegerField(choices=((0, '否'), (1, '是')), default=0, verbose_name=u'是否启用ghost改表')
    ghost_pid = models.IntegerField(null=False, default=0, verbose_name=u'ghost进程pid')
    sqlsha1 = models.CharField(null=False, max_length=120, default='', verbose_name=u'sqlsha1')
    rollback_sqlsha1 = models.CharField(null=False, max_length=120, default='', verbose_name=u'rollback sqlsha1')
    celery_task_id = models.CharField(null=False, max_length=256, default='', verbose_name=u'celery执行任务ID')
    exec_status = models.CharField(max_length=10, default='0', choices=exec_progress, verbose_name=u'执行进度')
    sequence = models.CharField(null=False, default='', max_length=1024, verbose_name=u'备份记录id，inception生成的sequence')
    affected_row = models.IntegerField(null=False, default=0, verbose_name=u'预计影响行数')
    backup_dbname = models.CharField(null=False, max_length=1024, default='', verbose_name=u'inception生成的备份的库名')
    exec_log = models.TextField(verbose_name=u'执行的记录', default='')
    created_time = models.DateTimeField(auto_now_add=True, verbose_name=u'生成时间')

    class Meta:
        verbose_name = u'工单执行任务'
        verbose_name_plural = verbose_name

        default_permissions = ()
        db_table = 'sqlaudit_sql_orders_execute_tasks'


class SqlOrdersTasksVersions(models.Model):
    id = models.AutoField(primary_key=True, verbose_name=u'主键id')
    username = models.CharField(default='', null=False, max_length=128, verbose_name=u'创建用户')
    tasks_version = models.CharField(default='', null=False, max_length=128, verbose_name=u'任务版本')
    expire_time = models.DateTimeField(default='2000-11-01 01:01:01', verbose_name=u'任务截止上线日期')
    created_at = models.DateTimeField(auto_now_add=True, verbose_name=u'创建时间')

    def __str__(self):
        return self.tasks_version

    class Meta:
        verbose_name = u'SQL工单上线任务版本'
        verbose_name_plural = verbose_name

        db_table = 'sqlaudit_sql_orders_tasks_versions'
        unique_together = ('tasks_version',)


class SysConfig(models.Model):
    id = models.AutoField(primary_key=True, verbose_name=u'主键id')
    name = models.CharField(max_length=256, default='', null=False, verbose_name=u'名称')
    key = models.CharField(max_length=256, default='', null=False, verbose_name=u'key')
    value = models.TextField(max_length=256, null=True, blank=True, verbose_name=u'值')
    is_enabled = models.CharField(max_length=2, choices=(('0', '启用'), ('1', '禁用')), default='1', verbose_name=u'是否启用')
    created_at = models.DateTimeField(auto_now_add=True, verbose_name=u'创建时间')

    class Meta:
        verbose_name = u'系统配置'
        verbose_name_plural = verbose_name

        default_permissions = ()
        db_table = 'sqlaudit_sys_config'


class SqlOrderReply(models.Model):
    id = models.AutoField(primary_key=True, verbose_name=u'主键')
    reply = models.ForeignKey(SqlOrdersContents, on_delete=models.CASCADE, null=False, default='')
    user = models.ForeignKey(UserAccounts, on_delete=models.CASCADE, null=False, default='')
    reply_contents = models.TextField(default='', verbose_name=u'回复内容')
    created_at = models.DateTimeField(auto_now_add=True, verbose_name=u'回复时间')

    class Meta:
        verbose_name = u'工单回复'
        verbose_name_plural = verbose_name

        default_permissions = ()
        db_table = 'sqlaudit_sql_order_reply'

    def reply_id(self):
        return self.reply.id

    def user_id(self):
        return self.user.uid
