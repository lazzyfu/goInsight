# -*- coding:utf-8 -*-
# edit by fuzongfei

from django.contrib import admin
from django_celery_results.models import TaskResult

from sqlorders.models import MysqlConfig, SysConfig, SqlOrdersEnvironment, SqlOrdersContents
from sqlorders.models import SqlOrdersExecTasks


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


class SqlOrdersContentsAdmin(admin.ModelAdmin):
    list_display = ('title', 'sql_type', 'envi_id', 'proposer', 'progress', 'created_at')
    list_display_links = ('title',)
    search_fields = ('title', 'proposer')
    fieldsets = (
        ('详情',
         {'fields': ['title', 'description', 'envi_id', 'progress', 'remark', 'task_version', 'proposer', 'auditor',
                     'email_cc', 'sql_type', 'host', 'port',
                     'database']}
         ),
        ('内容',
         {'fields': ['contents']}
         )
    )
    readonly_fields = ('host', 'port', 'database', 'proposer', 'auditor', 'sql_type', 'email_cc')


class MysqlConfigAdmin(admin.ModelAdmin):
    list_display = ('host', 'port', 'user', 'envi_id', 'is_type', 'comment', 'updated_at')
    ordering = ('-created_at',)
    list_display_links = ('host',)

    def get_readonly_fields(self, request, obj=None):
        if obj:
            return ["host", ]
        else:
            return []


class SysConfigAdmin(admin.ModelAdmin):
    list_display = ('name', 'value', 'is_enabled')
    list_display_links = ('name',)
    readonly_fields = ('key', 'name')
    fields = ('name', 'value', 'is_enabled')


class SqlOrdersEnvironmentAdmin(admin.ModelAdmin):
    list_display = ('envi_id', 'envi_name')
    list_display_links = ('envi_name',)


admin.site.unregister(TaskResult)
admin.site.register(SqlOrdersEnvironment, SqlOrdersEnvironmentAdmin)
admin.site.register(MysqlConfig, MysqlConfigAdmin)
admin.site.register(SysConfig, SysConfigAdmin)
admin.site.register(SqlOrdersContents, SqlOrdersContentsAdmin)
admin.site.register(SqlOrdersExecTasks, SqlOrdersExecTasksAdmin)
