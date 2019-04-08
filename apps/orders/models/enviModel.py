# -*- coding:utf-8 -*-
# edit by fuzongfei
from django.db import models


class SysEnvironment(models.Model):
    envi_id = models.AutoField(primary_key=True, null=False, verbose_name=u'环境ID')
    envi_name = models.CharField(max_length=30, default='', null=False, verbose_name=u'环境名')
    updated_at = models.DateTimeField(auto_now=True, verbose_name=u'更新时间')
    created_at = models.DateTimeField(auto_now_add=True, verbose_name=u'创建时间')

    def __str__(self):
        return self.envi_name

    class Meta:
        verbose_name = u'系统环境'
        verbose_name_plural = verbose_name

        default_permissions = ()
        app_label = 'orders'
        db_table = 'auditsql_sys_environment'

