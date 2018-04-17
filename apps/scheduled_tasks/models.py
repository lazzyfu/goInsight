from django.db import models


# Create your models here.


class DeadLockRecords(models.Model):
    ts = models.DateTimeField(null=False)
    thread = models.IntegerField(null=False)
    txn_time = models.SmallIntegerField(null=False)
    user = models.CharField(max_length=16, null=False)
    ip = models.CharField(max_length=15, null=False)
    db = models.CharField(max_length=64, null=False)
    tbl = models.CharField(max_length=64, null=False)
    idx = models.CharField(max_length=64, null=False)
    lock_type = models.CharField(max_length=16, null=False)
    lock_mode = models.CharField(max_length=1, null=False)
    wait_hold = models.CharField(max_length=1, null=False)
    victim = models.SmallIntegerField(null=False)
    query = models.TextField()
    is_sign = models.CharField(choices=(('0', '未推送'), ('1', '已推送')), default='0', max_length=1)
    src_host = models.CharField(max_length=30, default='')

    class Meta:
        verbose_name = u'记录死锁表'
        verbose_name_plural = verbose_name
        permissions = ()
        db_table = "dbaudit_deadlocks_records"
