from django.db import models

# Create your models here.
from sqlorders.models import MysqlSchemas
from users.models import UserAccounts


class MysqlSchemasGrant(models.Model):
    id = models.AutoField(primary_key=True, verbose_name=u'主键id')
    user = models.ForeignKey(UserAccounts, null=True, blank=True, on_delete=models.SET_NULL)
    schema = models.ForeignKey(MysqlSchemas, to_field='schema_join', on_delete=models.CASCADE)
    created_at = models.DateTimeField(auto_now_add=True, verbose_name=u'创建时间')
    updated_at = models.DateTimeField(auto_now=True, verbose_name=u'更新时间')

    class Meta:
        verbose_name = u'SQL授权表'
        verbose_name_plural = verbose_name

        default_permissions = ()
        db_table = 'sqlaudit_schemas_grant'
        unique_together = (('user', 'schema'),)


class MySQLQueryLog(models.Model):
    id = models.AutoField(primary_key=True, verbose_name=u'主键id')
    user = models.CharField(max_length=30, null=False, verbose_name=u'用户名')
    host = models.CharField(max_length=128, null=False, verbose_name=u'目标数据库地址')
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
        db_table = 'sqlaudit_sql_query_log'

