from django.db import models


# Create your models here.

class MySQLQueryLog(models.Model):
    id = models.AutoField(primary_key=True, verbose_name=u'主键id')
    user = models.CharField(max_length=30, null=False, verbose_name=u'用户名')
    host = models.CharField(max_length=32, null=False, verbose_name=u'目标数据库地址')
    database = models.CharField(max_length=32, null=False, verbose_name=u'目标数据库')
    query_sql = models.TextField(null=False, default='', verbose_name=u'查询SQL')
    query_status = models.CharField(max_length=2048, default='', verbose_name=u'查询状态，成功或失败的原因')
    affect_rows = models.IntegerField(default=0, null=False, verbose_name=u'影响行数')
    created_at = models.DateTimeField(auto_now_add=True, verbose_name=u'查询时间')

    class Meta:
        verbose_name = u'mysql查询记录日志'
        verbose_name_plural = verbose_name

        default_permissions = ()
        db_table = 'auditsql_mysql_query_log'

