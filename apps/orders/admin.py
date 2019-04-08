from django.contrib import admin

# Register your models here.
from django_celery_beat.models import SolarSchedule
from django_celery_results.models import TaskResult
from django_celery_beat.apps import BeatConfig

from orders.models import SysEnvironment, MysqlConfig, Orders, OrdersTasks

BeatConfig.verbose_name = '定时任务'


@admin.register(SysEnvironment)
class SysEnvironmentAdmin(admin.ModelAdmin):
    list_display = ('envi_name', 'updated_at', 'created_at')
    list_display_links = ('envi_name',)


@admin.register(MysqlConfig)
class MysqlConfigAdmin(admin.ModelAdmin):
    list_display = (
        'host', 'port', 'user', 'envi', 'character', 'colored_type', 'rds_type', 'comment', 'updated_at', 'created_at',
    )
    ordering = ('-created_at',)
    search_fields = ('host', 'user', 'envi',)
    list_per_page = 20
    list_display_links = ('host',)
    list_filter = ('envi', 'type', 'rds_type')

    def get_readonly_fields(self, request, obj=None):
        if obj:
            return ["host", ]
        else:
            return []


@admin.register(Orders)
class OrdersAdmin(admin.ModelAdmin):
    list_display = ('title', 'sql_type', 'envi', 'applicant', 'progress', 'created_at')
    list_display_links = ('title',)
    search_fields = ('title', 'applicant')
    fieldsets = (
        ('详情',
         {'fields': ['title', 'description', 'envi', 'progress', 'remark',
                     'version', 'applicant',
                     'sql_type', 'host', 'port', 'database']}
         ),
        ('内容',
         {'fields': ['contents']}
         )
    )
    readonly_fields = ('title', 'description', 'envi', 'remark', 'host', 'port', 'database',
                       'applicant', 'sql_type', 'version')

    def has_add_permission(self, request):
        return False

    def has_delete_permission(self, request, obj=None):
        return False


@admin.register(OrdersTasks)
class OrdersTasksAdmin(admin.ModelAdmin):
    list_display = ('applicant', 'taskid', 'order', 'envi', 'host', 'sql_content')
    ordering = ('-created_time',)
    list_display_links = ('taskid', 'applicant')
    search_fields = ('taskid', 'applicant', 'host')
    list_per_page = 20
    list_filter = ('applicant', 'host', 'envi', 'created_time', 'sql_type')
    fieldsets = (
        ('任务详情',
         {'fields': ['applicant', 'taskid', 'envi', 'host', 'port', 'database', 'task_progress',
                     'sql_type', 'export_file_format', 'sql']}),
    )
    readonly_fields = ('applicant', 'taskid', 'order', 'envi', 'export_file_format',
                       'host', 'port', 'database', 'sql_type')

    def has_delete_permission(self, request, obj=None):
        return False

    def has_add_permission(self, request):
        return False


admin.site.unregister(TaskResult)
admin.site.unregister(SolarSchedule)
