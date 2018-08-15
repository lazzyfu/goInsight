from django.contrib import admin

# Register your models here.
from mstats.models import MysqlSlowLog, WebShellInfo, MySQLConfigSource


@admin.register(MysqlSlowLog)
class MysqlSlowLogAdmin(admin.ModelAdmin):
    list_display = ('id', 'hostname', 'qps', 'cnt', 'avg', 'md5sum', 'is_pull', 'fingerprint', 'created_at')
    list_display_links = ('hostname', 'md5sum',)
    ordering = ('-created_at',)


@admin.register(WebShellInfo)
class WebShellInfoAdmin(admin.ModelAdmin):
    list_display = ('id', 'comment', 'command', 'envi', 'created_at')
    list_display_links = ('comment',)
    ordering = ('-created_at',)


@admin.register(MySQLConfigSource)
class MySQLConfigSourceAdmin(admin.ModelAdmin):
    list_display = ('id', 'host', 'port', 'envi', 'is_master', 'comment', 'created_at', 'updated_at')
    list_display_links = ('host', 'comment',)
    ordering = ('-created_at',)
