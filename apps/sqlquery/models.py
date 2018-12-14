# -*- coding:utf-8 -*-
# edit by fuzongfei

from django.db import models

# Create your models here.
from sqlorders.models import MysqlSchemas, MysqlConfig
from users.models import UserAccounts


class MysqlRulesChain(models.Model):
    id = models.AutoField(primary_key=True, verbose_name=u'主键id')
    cid = models.ForeignKey(MysqlConfig, to_field='id', on_delete=models.CASCADE, verbose_name=u'主机')
    action = models.CharField(max_length=32, choices=(('allow', u'允许'),), verbose_name=u'允许规则')
    schema = models.CharField(max_length=128, null=False, default='', verbose_name=u'库名')
    table = models.CharField(max_length=128, null=False, default='*', verbose_name=u'表名')
    comment = models.CharField(max_length=255, null=False, default='', verbose_name=u'规则')
    created_at = models.DateTimeField(auto_now_add=True, verbose_name=u'创建时间')
    updated_at = models.DateTimeField(auto_now=True, verbose_name=u'更新时间')

    def __str__(self):
        return self.comment

    class Meta:
        verbose_name = u'MySQL授权规则链'
        verbose_name_plural = verbose_name

        default_permissions = ()
        db_table = 'sqlaudit_mysql_rules_chain'

        unique_together = (('schema', 'table'),)


class MysqlRulesGroup(models.Model):
    id = models.AutoField(primary_key=True, verbose_name=u'主键')
    name = models.CharField(max_length=30, null=False, default='', verbose_name=u'组名')
    rule = models.ManyToManyField(MysqlRulesChain)
    user = models.ManyToManyField(UserAccounts)
    created_at = models.DateTimeField(auto_now_add=True, verbose_name=u'创建时间')
    updated_at = models.DateTimeField(auto_now=True, verbose_name=u'更新时间')

    def __str__(self):
        return self.name

    class Meta:
        verbose_name = u'MySQL授权规则组'
        verbose_name_plural = verbose_name

        default_permissions = ()
        db_table = 'sqlaudit_mysql_rules_group'


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
        verbose_name = u'MySQL查询日志'
        verbose_name_plural = verbose_name

        default_permissions = ()
        db_table = 'sqlaudit_sql_query_log'
