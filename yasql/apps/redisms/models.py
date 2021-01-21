# -*- coding:utf-8 -*-
# __Author__: PanDongLin
from django.db import models
from users.models import UserAccounts
from libs.SecurityTools import KMS


class RedisEnv(models.Model):
    """
    redis所在的环境，也可以叫做业务线等，可根据自己的业务线理解
    比如：测试环境、预发布环境、压测环境、生产环境，每个环境均有N套Redis集群
    """
    id = models.AutoField(primary_key=True, verbose_name=u'主键ID')
    name = models.CharField(max_length=128, unique=True, verbose_name='名称')
    updated_at = models.DateTimeField(auto_now=True, verbose_name=u'更新时间')

    def __str__(self):
        return self.name

    class Meta:
        verbose_name = u'redis环境'
        verbose_name_plural = verbose_name
        db_table = 'yasql_redis_env'
        default_permissions = ()


class RedisGroup(models.Model):
    """授权是以redis组为基础，用户加入组"""
    id = models.AutoField(primary_key=True, verbose_name='主键id')
    name = models.CharField(max_length=32, verbose_name="redis组名")
    env = models.ForeignKey("RedisEnv", related_name="rg_env", on_delete=models.CASCADE, verbose_name="环境")
    updated_at = models.DateTimeField(auto_now=True, verbose_name='更新时间')

    class Meta:
        verbose_name = u'redis组'
        verbose_name_plural = verbose_name
        db_table = 'yasql_redis_group'
        default_permissions = ()

    @property
    def display_full_name(self):
        return "%s-%s" % (self.name, self.env.name)

    def __str__(self):
        return self.display_full_name


class RedisGrant(models.Model):
    """把用户和组关联起来"""
    PERMISSION_TYPE = (
        ('read_only', 'read_only'),  # 只读
        ('read_write', 'read_write'), # 读写
    )
    id = models.AutoField(primary_key=True, verbose_name='主键id')
    user = models.ForeignKey(UserAccounts, related_name="rg_user", on_delete=models.CASCADE, verbose_name="用户")
    group = models.ForeignKey("RedisGroup", related_name='rg_group', on_delete=models.CASCADE, verbose_name="redis组")
    permission_code = models.CharField(choices=PERMISSION_TYPE, max_length=32, verbose_name="读写权限")
    updated_at = models.DateTimeField(auto_now=True, verbose_name='更新时间')

    class Meta:
        unique_together = ["user", "group", "permission_code"]
        verbose_name = u'redis授权'
        verbose_name_plural = verbose_name
        db_table = 'yasql_redis_grant'
        default_permissions = ()

    def __str__(self):
        return "%s:%s-%s" %(self.user, self.group, self.permission_code)


class RedisConfig(models.Model):
    """redis实例配置"""
    id = models.AutoField(primary_key=True, verbose_name='主键id')
    group = models.ForeignKey("RedisGroup", related_name='rg', on_delete=models.CASCADE, verbose_name="redis组")
    host = models.CharField(max_length=128, verbose_name=u'地址')
    port = models.IntegerField(default=6379, verbose_name=u'端口')
    password = models.CharField(max_length=256, null=True, blank=True, verbose_name=u'密码')
    updated_at = models.DateTimeField(auto_now=True, verbose_name=u'更新时间')

    class Meta:
        verbose_name = u'redis配置'
        verbose_name_plural = verbose_name
        db_table = 'yasql_redis_config'
        default_permissions = ()

    def __str__(self):
        return self.full_dsn

    @property
    def full_dsn(self):
        return "%s:%s" % (self.host, self.port)

    @property
    def decrypt_password(self):
        if self.password:
            password = KMS().decrypt(self.password)
            return password
        else:
            return self.password


class RedisLog(models.Model):
    id = models.AutoField(primary_key=True, verbose_name=u'主键id')
    user = models.CharField(max_length=32, verbose_name=u'用户名')
    host_dsn = models.CharField(max_length=255, verbose_name=u'redis地址')
    exec_cmd = models.TextField(null=False, blank=True, verbose_name=u'执行命令')
    created_at = models.DateTimeField(auto_now_add=True, verbose_name=u'执行时间')

    class Meta:
        verbose_name = u'redis执行日志'
        verbose_name_plural = verbose_name
        db_table = 'yasql_redis_log'
        default_permissions = ()

    def __str__(self):
        return "%s > %s" %(self.host_dsn, self.exec_cmd)