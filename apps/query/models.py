# -*- coding:utf-8 -*-
# edit by fuzongfei

from django.db import models

# Create your models here.
from orders.models import MysqlConfig
from users.models import UserAccounts


class MysqlPrivBlacklist(models.Model):
    id = models.AutoField(primary_key=True, verbose_name=u'主键id')
    schema = models.CharField(max_length=128, null=False, default='', verbose_name=u'库名')
    table = models.CharField(max_length=128, null=False, default='*', verbose_name=u'表名')
    columns = models.CharField(max_length=4096, null=False, default='*', verbose_name=u'列名')
    comment = models.CharField(max_length=255, null=False, default='', verbose_name=u'描述')
    created_at = models.DateTimeField(auto_now_add=True, verbose_name=u'创建时间')
    updated_at = models.DateTimeField(auto_now=True, verbose_name=u'更新时间')

    def __str__(self):
        return self.comment

    class Meta:
        verbose_name = u'库表黑名单'
        verbose_name_plural = verbose_name

        default_permissions = ()
        db_table = 'auditsql_mysql_priv_blacklist'

        unique_together = (('schema', 'table'),)


class MysqlUserGroupMap(models.Model):
    id = models.AutoField(primary_key=True, verbose_name=u'主键id')
    user = models.ManyToManyField(UserAccounts, blank=True, verbose_name='用户')
    group = models.CharField(max_length=128, null=False, default='', verbose_name='MySQL用户组名')
    schema = models.CharField(max_length=128, null=False, default='', verbose_name='库')
    comment = models.ForeignKey(MysqlConfig, to_field='id', on_delete=models.CASCADE, verbose_name=u'主机描述')
    created_at = models.DateTimeField(auto_now_add=True, verbose_name=u'创建时间')
    updated_at = models.DateTimeField(auto_now=True, verbose_name=u'更新时间')

    def group_alias(self):
        if self.group.startswith('n_'):
            return '普通用户权限'
        if self.group.startswith('s_'):
            return '管理员权限'

    def __str__(self):
        user = '普通用户权限' if self.group.startswith('n_') else '管理员权限'
        return '-'.join([MysqlConfig.objects.get(pk=self.comment_id).comment, self.schema, user])

    class Meta:
        verbose_name = u'库表权限组'
        verbose_name_plural = verbose_name

        default_permissions = ()
        db_table = 'auditsql_mysql_usergroup_map'

        unique_together = (('comment', 'group'),)


class QueryBusinessGroup(models.Model):
    id = models.AutoField(primary_key=True, verbose_name=u'主键id')
    user = models.ManyToManyField(UserAccounts, blank=True, verbose_name='用户')
    group = models.CharField(max_length=128, null=False, default='', verbose_name='业务组名')
    config = models.ForeignKey(MysqlConfig, default='', on_delete=models.CASCADE, verbose_name=u'关联主机')
    schema = models.CharField(max_length=128, null=False, default='', verbose_name='关联库')
    tables = models.TextField(verbose_name=u'关联表')
    map_mysqluser = models.ForeignKey(MysqlUserGroupMap, on_delete=models.CASCADE,
                                      verbose_name=u'映射的mysql用户')
    created_at = models.DateTimeField(auto_now_add=True, verbose_name=u'创建时间')
    updated_at = models.DateTimeField(auto_now=True, verbose_name=u'更新时间')

    def __str__(self):
        return self.group

    class Meta:
        verbose_name = u'库表业务组'
        verbose_name_plural = verbose_name

        default_permissions = ()
        db_table = 'yops_query_business_group'


class MySQLQueryLog(models.Model):
    id = models.AutoField(primary_key=True, verbose_name=u'主键id')
    user = models.CharField(max_length=30, null=False, verbose_name=u'用户名')
    host = models.CharField(max_length=128, null=False, verbose_name=u'目标数据库地址')
    database = models.CharField(max_length=32, null=False, verbose_name=u'目标数据库')
    envi = models.SmallIntegerField(null=False, default=1, verbose_name=u'环境')
    query_sql = models.TextField(null=False, default='', verbose_name=u'查询SQL')
    query_time = models.CharField(null=False, default='', max_length=128, verbose_name=u'查询时间，单位s')
    query_status = models.CharField(max_length=2048, default='', verbose_name=u'查询状态，成功或失败的原因')
    affect_rows = models.IntegerField(default=0, null=False, verbose_name=u'影响行数')
    created_at = models.DateTimeField(auto_now_add=True, verbose_name=u'查询时间')

    class Meta:
        verbose_name = u'查询日志'
        verbose_name_plural = verbose_name

        default_permissions = ()
        db_table = 'auditsql_sql_query_log'
