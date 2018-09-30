from django import forms
from django.contrib import admin

# Register your models here.
from django_celery_results.models import TaskResult

from sqlorders.models import MysqlConfig, SqlOrdersExecTasks, SysConfig, SqlOrdersEnvironment


class MysqlConfigAdmin(admin.ModelAdmin):
    list_display = ('host', 'port', 'user', 'envi_id', 'is_master', 'comment', 'updated_at')
    ordering = ('-created_at',)
    list_display_links = ('host',)


class SqlOrdersExecTasksAdmin(admin.ModelAdmin):
    list_display = ('user', 'taskid', 'envi_id', 'host', 'sql')
    ordering = ('-created_time',)
    list_display_links = ('taskid', 'user')
    search_fields = ('taskid', 'user', 'host')
    fieldsets = (
        ('任务详情',
         {'fields': ['user', 'taskid', 'envi_id', 'host', 'port', 'database', 'exec_status', 'sql_type', 'sql']}),
    )
    readonly_fields = ('user', 'taskid', 'envi_id', 'host', 'port', 'sql_type')


class SysConfigAdmin(admin.ModelAdmin):
    list_display = ('name', 'value', 'is_enabled')
    list_display_links = ('name',)
    readonly_fields = ('key', 'name')
    fields = ('name', 'value', 'is_enabled')


class SqlOrdersEnvironmentAdmin(admin.ModelAdmin):
    list_display = ('envi_id', 'parent_id', 'envi_name')
    list_display_links = ('envi_name',)


admin.site.register(SqlOrdersEnvironment, SqlOrdersEnvironmentAdmin)
admin.site.register(MysqlConfig, MysqlConfigAdmin)
admin.site.register(SysConfig, SysConfigAdmin)
admin.site.register(SqlOrdersExecTasks, SqlOrdersExecTasksAdmin)
admin.site.unregister(TaskResult)
