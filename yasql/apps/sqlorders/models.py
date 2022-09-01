# -*- coding:utf-8 -*-
# edit by xff
import humanfriendly as humanfriendly
from django.db import models
from django_jsonfield_backport.models import JSONField

# Create your models here.
from sqlorders import utils


class DbEnvironment(models.Model):
    """
    DB所在的环境，也可以叫做业务线等，可根据自己的业务线理解，此处姑且叫做DB环境
    比如：测试环境、预发布环境、压测环境、生产环境，每个环境均有N套DB集群
    """
    id = models.AutoField(primary_key=True, verbose_name=u'主键ID')
    name = models.CharField(max_length=128, null=False, unique=True, verbose_name='名称')
    created_at = models.DateTimeField(auto_now_add=True, verbose_name=u'创建时间')
    updated_at = models.DateTimeField(auto_now=True, verbose_name=u'更新时间')

    def __str__(self):
        return self.name

    class Meta:
        verbose_name = u'DB环境'
        verbose_name_plural = verbose_name

        default_permissions = ()
        app_label = 'sqlorders'
        db_table = 'yasql_db_environment'


class ReleaseVersions(models.Model):
    """
    工单发布版本号
    用于每周上线时，设定的发布版本号，一般和代码发布版本号相关联
    """
    id = models.AutoField(primary_key=True, verbose_name=u'主键id')
    username = models.CharField(default='', null=False, max_length=128, verbose_name=u'创建用户')
    version = models.CharField(default='', null=False, max_length=128, unique=True, verbose_name=u'版本号')
    expire_time = models.DateField(verbose_name=u'截止上线日期')
    created_at = models.DateTimeField(auto_now_add=True, verbose_name=u'创建时间')
    updated_at = models.DateTimeField(auto_now=True, verbose_name=u'更新时间')

    def __str__(self):
        return self.version

    class Meta:
        verbose_name = u'工单发布版本号'
        verbose_name_plural = verbose_name

        default_permissions = ()
        app_label = 'sqlorders'
        db_table = 'yasql_release_versions'


class DbConfig(models.Model):
    """
    mysql远程主机或实例配置表
    """
    id = models.AutoField(primary_key=True, verbose_name=u'主键ID')
    host = models.CharField(max_length=128, null=False, verbose_name=u'数据库地址')
    port = models.IntegerField(null=False, default=3306, verbose_name=u'数据库端口')
    character = models.CharField(max_length=32, null=False, choices=utils.characterChoice,
                                 default='utf8', verbose_name=u'表字符集')
    custom_audit_parameters = JSONField(default=dict, blank=True, null=True, verbose_name='自定义审核参数')
    env = models.ForeignKey(DbEnvironment, blank=True, null=True, on_delete=models.SET_NULL, verbose_name='环境')
    use_type = models.SmallIntegerField(choices=utils.useTypeChoice, default=0, verbose_name=u'用途')
    rds_type = models.SmallIntegerField(choices=utils.rdsTypeChoice, default=2, verbose_name=u'数据库的类型')
    rds_category = models.SmallIntegerField(choices=utils.rdsCategory, default=1, verbose_name=u'数据库类别')
    comment = models.CharField(max_length=128, null=True, verbose_name=u'主机描述')
    created_at = models.DateTimeField(auto_now_add=True, verbose_name=u'创建时间')
    updated_at = models.DateTimeField(auto_now=True, verbose_name=u'更新时间')

    def __str__(self):
        return self.comment

    class Meta:
        verbose_name = u'DB主机配置'
        verbose_name_plural = verbose_name

        # 允许一个实例存在2中用途,即：查询和审核
        unique_together = (('host', 'port', 'use_type'),)

        default_permissions = ()
        app_label = 'sqlorders'
        db_table = 'yasql_dbconfig'


class DbSchemas(models.Model):
    """
    存储远程数据库的库名
    """
    id = models.AutoField(primary_key=True, verbose_name=u'主键ID')
    cid = models.ForeignKey(DbConfig, blank=True, null=True, on_delete=models.SET_NULL, verbose_name='数据库')
    schema = models.CharField(null=False, max_length=64, default='', verbose_name=u'库名')
    created_at = models.DateTimeField(auto_now_add=True, verbose_name=u'创建时间')
    updated_at = models.DateTimeField(auto_now=True, verbose_name=u'更新时间')

    class Meta:
        verbose_name = u'DB库'
        verbose_name_plural = verbose_name

        unique_together = (('cid', 'schema'),)

        default_permissions = ()
        app_label = 'sqlorders'
        db_table = 'yasql_dbschemas'


class DbOrders(models.Model):
    id = models.AutoField(primary_key=True, verbose_name=u'主键ID')
    title = models.CharField(max_length=64, null=False, verbose_name=u'工单标题')
    order_id = models.CharField(max_length=128, null=False, blank=True, default='', verbose_name=u'工单号')
    demand = models.CharField(max_length=256, null=False, blank=True, verbose_name=u'需求描述')
    is_hide = models.CharField(max_length=10, null=False, blank=True, default='OFF', verbose_name=u'是否隐藏')
    remark = models.CharField(max_length=16, choices=utils.orderRemark, verbose_name=u'工单备注')
    rds_category = models.SmallIntegerField(choices=utils.rdsCategory, default=1, verbose_name=u'数据库类别')
    sql_type = models.CharField(max_length=30, null=False, default='DML', choices=utils.sqlTypeChoice,
                                verbose_name=u'工单类型')
    file_format = models.CharField(max_length=30, choices=utils.fileFormatChoice, null=False, blank=True,
                                   verbose_name=u'导出工单的文件格式')
    env = models.ForeignKey(DbEnvironment, blank=True, null=True, on_delete=models.SET_NULL, verbose_name='环境')
    applicant = models.CharField(max_length=30, null=False, verbose_name=u'工单申请人')
    department = models.CharField(max_length=128, null=False, blank=True, verbose_name=u'申请人所在的部门')
    auditor = models.CharField(max_length=512, null=False, verbose_name=u'工单审核人')
    executor = models.CharField(max_length=512, null=True, blank=True, verbose_name=u'工单执行人')
    closer = models.CharField(max_length=512, null=True, blank=True, verbose_name=u'工单关闭人')
    reviewer = models.CharField(max_length=512, null=False, verbose_name=u'工单复核人')
    email_cc = models.CharField(max_length=2048, null=True, blank=True, verbose_name=u'抄送人')
    cid = models.ForeignKey(DbConfig, blank=True, null=True, on_delete=models.SET_NULL, verbose_name='数据库')
    database = models.CharField(max_length=32, null=False, verbose_name=u'库名')
    progress = models.SmallIntegerField(default=0, choices=utils.sqlProgressChoice, verbose_name=u'进度')
    version = models.ForeignKey('ReleaseVersions', blank=True, null=True, on_delete=models.SET_NULL,
                                verbose_name=u'上线版本号')
    contents = models.TextField(default='', verbose_name=u'工单内容')
    created_at = models.DateTimeField(auto_now_add=True, verbose_name=u'创建时间')
    updated_at = models.DateTimeField(auto_now=True, verbose_name=u'更新时间')

    def __str__(self):
        return self.title

    class Meta:
        verbose_name = u'DB工单记录'
        verbose_name_plural = verbose_name

        default_permissions = ()
        app_label = 'sqlorders'
        db_table = 'yasql_dborders'

        unique_together = (('order_id',),)
        index_together = (('title',), ('remark',), ('applicant',), ('database',), ('progress',))


class DbOrdersExecuteTasks(models.Model):
    id = models.AutoField(primary_key=True, verbose_name='主键id')
    order = models.ForeignKey(DbOrders, blank=True, null=True, on_delete=models.SET_NULL, verbose_name='关联工单id')
    applicant = models.CharField(max_length=30, null=True, blank=True, verbose_name=u'工单申请人')
    task_id = models.CharField(null=False, max_length=128, verbose_name='任务ID')
    sql_type = models.CharField(max_length=30, null=False, default='DML', choices=utils.sqlTypeChoice,
                                verbose_name=u'工单类型')
    executor = models.CharField(max_length=30, null=False, default='', verbose_name='工单执行人')
    sql = models.TextField(verbose_name='执行的SQL', default='')
    progress = models.SmallIntegerField(default=0, choices=utils.taskProgressChoice, verbose_name=u'执行进度')
    affected_rows = models.IntegerField(default=0, verbose_name=u'影响行数')
    consuming_time = models.DecimalField(default=0.000, max_digits=10, decimal_places=3, verbose_name='耗时')
    execute_log = models.TextField(verbose_name=u'执行的日志', default='')
    rollback_sql = models.TextField(verbose_name=u'回滚的SQL', default='')
    file_format = models.CharField(max_length=30, choices=(('xlsx', 'xlsx'), ('csv', 'csv'), ('txt', 'txt')),
                                   default='xlsx', verbose_name=u'导出的文件格式')
    execute_time = models.DateTimeField(auto_now=True, verbose_name='工单执行时间')
    created_time = models.DateTimeField(auto_now_add=True, verbose_name=u'创建时间')

    class Meta:
        verbose_name = u'DB工单执行任务'
        verbose_name_plural = verbose_name

        index_together = (('task_id',), ('progress',),)

        default_permissions = ()
        app_label = 'sqlorders'
        db_table = 'yasql_dborders_execute_tasks'


class DbExportFiles(models.Model):
    id = models.AutoField(primary_key=True, verbose_name=u'主键id')
    task = models.ForeignKey(DbOrdersExecuteTasks, blank=True, null=True, on_delete=models.SET_NULL,
                             verbose_name=u'关联执行任务的主键id')
    file_name = models.CharField(max_length=200, default='', unique=True, verbose_name=u'文件名')
    file_size = models.IntegerField(default=0, verbose_name=u'文件大小，单位B')
    files = models.FileField(upload_to='export/')
    encryption_key = models.CharField(max_length=128, default='', verbose_name='加密密钥')
    content_type = models.CharField(max_length=100, default='', verbose_name=u'文件的类型')
    created_at = models.DateTimeField(auto_now_add=True, verbose_name=u'创建时间')

    def size(self):
        return humanfriendly.format_size(self.file_size, binary=True)

    class Meta:
        verbose_name = u'DB导出文件'
        verbose_name_plural = verbose_name

        default_permissions = ()
        db_table = 'yasql_db_export_files'
