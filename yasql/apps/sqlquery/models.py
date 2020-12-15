from django.db import models

# Create your models here.
from sqlorders.models import DbConfig
from users.models import UserAccounts


class DbQuerySchemas(models.Model):
    """
    存储远程查询数据库的库名
    """
    id = models.AutoField(primary_key=True, verbose_name=u'主键ID')
    cid = models.ForeignKey(DbConfig, blank=True, null=True, on_delete=models.SET_NULL, verbose_name='数据库')
    schema = models.CharField(null=False, max_length=64, default='', verbose_name=u'库名')
    created_at = models.DateTimeField(auto_now_add=True, verbose_name=u'创建时间')
    updated_at = models.DateTimeField(auto_now=True, verbose_name=u'更新时间')

    def __str__(self):
        return '.'.join([self.cid.comment, self.schema])

    def display_comment(self):
        return self.cid.comment

    display_comment.short_description = '主机'

    class Meta:
        verbose_name = u'DB查询库'
        verbose_name_plural = verbose_name

        unique_together = (('cid', 'schema'),)

        default_permissions = ()
        app_label = 'sqlquery'
        db_table = 'yasql_sqlquery_schemas'


class DbQueryTables(models.Model):
    """
    存储远程查询数据库的表名
    """
    id = models.AutoField(primary_key=True, verbose_name=u'主键ID')
    schema = models.ForeignKey(DbQuerySchemas, blank=True, null=True, on_delete=models.SET_NULL, verbose_name='库名')
    table = models.CharField(null=False, max_length=128, default='', verbose_name=u'表名')
    created_at = models.DateTimeField(auto_now_add=True, verbose_name=u'创建时间')
    updated_at = models.DateTimeField(auto_now=True, verbose_name=u'更新时间')

    def __str__(self):
        return '.'.join([self.schema.cid.comment, self.schema.schema, self.table])

    def display_comment(self):
        return self.schema.cid.comment

    def display_schema(self):
        return self.schema.schema

    display_comment.short_description = '主机'
    display_schema.short_description = '库名'

    class Meta:
        verbose_name = u'DB查询表'
        verbose_name_plural = verbose_name

        default_permissions = ()
        app_label = 'sqlquery'
        db_table = 'yasql_sqlquery_tables'


class DbQueryUserPrivs(models.Model):
    id = models.AutoField(primary_key=True, verbose_name=u'主键ID')
    user = models.ForeignKey(UserAccounts, blank=True, null=True, on_delete=models.SET_NULL, verbose_name='用户')
    schemas = models.ManyToManyField(DbQuerySchemas, verbose_name='允许访问的库')
    created_at = models.DateTimeField(auto_now_add=True, verbose_name=u'创建时间')
    updated_at = models.DateTimeField(auto_now=True, verbose_name=u'更新时间')

    def __str__(self):
        return self.user.username

    class Meta:
        verbose_name = u'DB查询用户权限'
        verbose_name_plural = verbose_name

        default_permissions = ()
        app_label = 'sqlquery'
        db_table = 'yasql_sqlquery_user_privileges'


class DbQueryGroupPrivs(models.Model):
    id = models.AutoField(primary_key=True, verbose_name=u'主键ID')
    group = models.CharField(null=False, max_length=128, default='', verbose_name=u'组名')
    user = models.ManyToManyField(UserAccounts, verbose_name='用户')
    schemas = models.ManyToManyField(DbQuerySchemas, verbose_name='允许访问的库')
    created_at = models.DateTimeField(auto_now_add=True, verbose_name=u'创建时间')
    updated_at = models.DateTimeField(auto_now=True, verbose_name=u'更新时间')

    def __str__(self):
        return self.group

    class Meta:
        verbose_name = u'DB查询组权限'
        verbose_name_plural = verbose_name

        default_permissions = ()
        app_label = 'sqlquery'
        db_table = 'yasql_sqlquery_group_privileges'


class DbQueryUserAllowedTables(models.Model):
    id = models.AutoField(primary_key=True, verbose_name=u'主键ID')
    tables = models.ForeignKey(DbQueryTables, blank=True, null=True, on_delete=models.SET_NULL, verbose_name='表')
    user_privs = models.ForeignKey(DbQueryUserPrivs, blank=True, null=True, on_delete=models.SET_NULL,
                                   verbose_name='权限')
    created_at = models.DateTimeField(auto_now_add=True, verbose_name=u'创建时间')
    updated_at = models.DateTimeField(auto_now=True, verbose_name=u'更新时间')

    def __str__(self):
        return '.'.join([self.tables.schema.cid.comment, self.tables.schema.schema, self.tables.table])

    class Meta:
        verbose_name = u'允许用户访问的表'
        verbose_name_plural = verbose_name

        default_permissions = ()
        app_label = 'sqlquery'
        db_table = 'yasql_sqlquery_user_allowed_tables'


class DbQueryUserDenyTables(models.Model):
    id = models.AutoField(primary_key=True, verbose_name=u'主键ID')
    tables = models.ForeignKey(DbQueryTables, blank=True, null=True, on_delete=models.SET_NULL, verbose_name='表')
    user_privs = models.ForeignKey(DbQueryUserPrivs, blank=True, null=True, on_delete=models.SET_NULL,
                                   verbose_name='权限')
    created_at = models.DateTimeField(auto_now_add=True, verbose_name=u'创建时间')
    updated_at = models.DateTimeField(auto_now=True, verbose_name=u'更新时间')

    def __str__(self):
        return '.'.join([self.tables.schema.cid.comment, self.tables.schema.schema, self.tables.table])

    class Meta:
        verbose_name = u'禁止用户访问的表'
        verbose_name_plural = verbose_name

        default_permissions = ()
        app_label = 'sqlquery'
        db_table = 'yasql_sqlquery_user_deny_tables'


class DbQueryGroupAllowedTables(models.Model):
    id = models.AutoField(primary_key=True, verbose_name=u'主键ID')
    tables = models.ForeignKey(DbQueryTables, blank=True, null=True, on_delete=models.SET_NULL, verbose_name='表')
    group_privs = models.ForeignKey(DbQueryGroupPrivs, blank=True, null=True, on_delete=models.SET_NULL,
                                    verbose_name='权限')
    created_at = models.DateTimeField(auto_now_add=True, verbose_name=u'创建时间')
    updated_at = models.DateTimeField(auto_now=True, verbose_name=u'更新时间')

    def __str__(self):
        return '.'.join([self.tables.schema.cid.comment, self.tables.schema.schema, self.tables.table])

    class Meta:
        verbose_name = u'允许组访问的表'
        verbose_name_plural = verbose_name

        default_permissions = ()
        app_label = 'sqlquery'
        db_table = 'yasql_sqlquery_group_allowed_tables'


class DbQueryGroupDenyTables(models.Model):
    id = models.AutoField(primary_key=True, verbose_name=u'主键ID')
    tables = models.ForeignKey(DbQueryTables, blank=True, null=True, on_delete=models.SET_NULL, verbose_name='表')
    group_privs = models.ForeignKey(DbQueryGroupPrivs, blank=True, null=True, on_delete=models.SET_NULL,
                                    verbose_name='权限')
    created_at = models.DateTimeField(auto_now_add=True, verbose_name=u'创建时间')
    updated_at = models.DateTimeField(auto_now=True, verbose_name=u'更新时间')

    def __str__(self):
        return '.'.join([self.tables.schema.cid.comment, self.tables.schema.schema, self.tables.table])

    class Meta:
        verbose_name = u'禁止组访问的表'
        verbose_name_plural = verbose_name

        default_permissions = ()
        app_label = 'sqlquery'
        db_table = 'yasql_sqlquery_group_deny_tables'


class DbQueryLog(models.Model):
    id = models.AutoField(primary_key=True, verbose_name=u'主键id')
    username = models.CharField(max_length=64, null=False, verbose_name=u'用户名')
    host = models.CharField(max_length=256, null=False, verbose_name=u'目标数据库地址')
    schema = models.CharField(null=False, max_length=128, default='', verbose_name=u'目标数据库')
    tables = models.CharField(null=False, max_length=200, default='', verbose_name=u'目标表名')
    query_sql = models.TextField(null=False, default='', verbose_name=u'查询SQL')
    query_consume_time = models.FloatField(null=False, default=0.000, verbose_name=u'查询耗时，单位s')
    query_status = models.CharField(max_length=2048, default='', verbose_name=u'查询是否成功或失败的原因')
    affected_rows = models.IntegerField(default=0, null=False, verbose_name=u'影响影响行数')
    created_at = models.DateTimeField(auto_now_add=True, verbose_name=u'查询时间')

    class Meta:
        verbose_name = u'DB查询日志'
        verbose_name_plural = verbose_name

        index_together = (('username',), ('schema',), ('tables',),)

        default_permissions = ()
        app_label = 'sqlquery'
        db_table = 'yasql_sqlquery_log'
