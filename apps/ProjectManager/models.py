from django.db import models

# Create your models here.
class InceptionHostConfig(models.Model):
    id = models.AutoField(primary_key=True, verbose_name=u'主键id')
    user = models.CharField(max_length=30, null=False, verbose_name=u'用户名')
    password = models.CharField(max_length=128, null=False, verbose_name=u'密码')
    host = models.CharField(max_length=32, null=False, verbose_name=u'ip地址')
    port = models.IntegerField(null=False, verbose_name=u'端口')
    comment = models.CharField(max_length=128, verbose_name=u'主机描述')

    class Meta:
        verbose_name = u'inception连接数据库配置'
        verbose_name_plural = verbose_name

        default_permissions = ()
        db_table = 'sqlaudit_inception_hostconfig'


class InceptionSqlOperateRecord(models.Model):
    id = models.AutoField(primary_key=True, verbose_name=u'主键id')
    op_uid = models.IntegerField(null=False, default=0, verbose_name=u'操作用户uid')
    op_user = models.CharField(max_length=30, null=False, verbose_name=u'操作用户')
    step_id = models.IntegerField(null=False, default=1, verbose_name=u'inception输出顺序ID')
    workid = models.CharField(null=False, max_length=128, verbose_name=u'工单号')
    dst_host = models.CharField(null=False, max_length=30, verbose_name=u'操作目标数据库主机')
    dst_database = models.CharField(null=False, max_length=80, verbose_name=u'操作目标数据库')
    stage = models.CharField(null=False, max_length=128, default='', verbose_name=u'inception检测阶段')
    stagestatus = models.CharField(null=False, default='', max_length=1024, verbose_name=u'inception检测阶段状态')
    errlevel = models.IntegerField(null=False, default=0, verbose_name=u'inception输出错误级别')
    errormessage = models.TextField(null=False, default='', verbose_name=u'inception输出错误信息')
    op_sql = models.TextField(verbose_name=u'执行的SQL', default='')
    affected_rows = models.IntegerField(null=False, default=0, verbose_name=u'影响的行数')
    sequence = models.CharField(null=False, default='', max_length=1024, verbose_name=u'备份记录id，inception生成的sequence')
    backup_dbname = models.CharField(null=False, max_length=1024, default='', verbose_name=u'inception生成的备份的库名')
    execute_time = models.CharField(null=False, default='0.000', max_length=128, verbose_name=u'inception执行时间')
    op_time = models.DateTimeField(auto_now_add=True, verbose_name=u'用户操作时间')

    class Meta:
        verbose_name = u'inception操作结果记录表'
        verbose_name_plural = verbose_name

        default_permissions = ()
        db_table = 'sqlaudit_inception_sql_operate_record'