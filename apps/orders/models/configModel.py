# -*- coding:utf-8 -*-
# edit by fuzongfei

from django.db import models
from django.utils.html import format_html

typeChoice = (
    (0, 'SQL审核'),
    (1, 'SQL查询'),
)
characterChoice = (
    ('utf8', 'utf8'),
    ('utf8mb4', 'utf8mb4')
)
rdsTypeChoice = (
    (0, '非阿里云RDS'),
    (1, '阿里云RDS')
)


class MysqlConfig(models.Model):
    """
    mysql远程主机配置表
    """
    id = models.AutoField(primary_key=True, verbose_name=u'主键ID')
    host = models.CharField(max_length=128, null=False, verbose_name=u'地址')
    port = models.IntegerField(null=False, default=3306, verbose_name=u'端口')
    user = models.CharField(max_length=32, null=False, verbose_name=u'用户名')
    password = models.CharField(max_length=64, null=False, verbose_name=u'密码')
    character = models.CharField(max_length=32, null=False, choices=characterChoice,
                                 default='utf8', verbose_name=u'库表字符集')
    envi = models.ForeignKey('orders.SysEnvironment', to_field='envi_id', null=False,
                             on_delete=models.CASCADE, verbose_name=u'环境')
    type = models.SmallIntegerField(choices=typeChoice, default=0, verbose_name=u'用途')
    rds_type = models.SmallIntegerField(choices=rdsTypeChoice, default=0, verbose_name=u'数据库的类型')
    comment = models.CharField(max_length=128, null=True, verbose_name=u'主机描述')
    created_at = models.DateTimeField(auto_now_add=True, verbose_name=u'创建时间')
    updated_at = models.DateTimeField(auto_now=True, verbose_name=u'更新时间')

    def __str__(self):
        return self.comment

    def colored_type(self):
        if self.type == 0:
            color_code = 'red'
        if self.type == 1:
            color_code = 'blue'
        return format_html(
            '<span style="color: {};">{}</span>', color_code, self.get_type_display()
        )

    class Meta:
        verbose_name = u'MySQL主机配置'
        verbose_name_plural = verbose_name

        unique_together = (('host', 'port'),)

        default_permissions = ()
        app_label = 'orders'
        db_table = 'auditsql_mysql_config'


class MysqlSchemas(models.Model):
    """
    mysql远程主机元数据采集表，存储远程主机的库信息
    """
    id = models.AutoField(primary_key=True, verbose_name=u'主键ID')
    cid = models.ForeignKey(MysqlConfig, null=False, to_field='id', on_delete=models.CASCADE, verbose_name=u'主机')
    schema = models.CharField(null=False, max_length=64, default='', verbose_name=u'库名')
    envi = models.ForeignKey('orders.SysEnvironment', null=False, to_field='envi_id', on_delete=models.CASCADE)
    created_at = models.DateTimeField(auto_now_add=True, verbose_name=u'创建时间')
    updated_at = models.DateTimeField(auto_now=True, verbose_name=u'更新时间')

    class Meta:
        verbose_name = u'MySQL库信息'
        verbose_name_plural = verbose_name

        unique_together = (('cid', 'schema'),)

        default_permissions = ()
        app_label = 'orders'
        db_table = 'auditsql_mysql_schemas'
