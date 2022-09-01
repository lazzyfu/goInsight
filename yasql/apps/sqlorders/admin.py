from django.contrib import admin

# Register your models here.
from sqlorders import models
from django_json_widget.widgets import JSONEditorWidget


@admin.register(models.DbConfig)
class DbConfigAdmin(admin.ModelAdmin):
    formfield_overrides = {
        # fields.JSONField: {'widget': JSONEditorWidget}, # if django < 3.1
        models.JSONField: {'widget': JSONEditorWidget},
    }
    list_display = ('host', 'port', 'character', 'custom_audit_parameters',
                    'env', 'use_type', 'rds_type', 'rds_category', 'comment')
    list_display_links = ('host',)
    list_filter = ('use_type', 'rds_type', 'rds_category')
    search_fields = ('host', 'port', 'env__name', 'comment')
    list_per_page = 10


@admin.register(models.DbSchemas)
class DbSchemasAdmin(admin.ModelAdmin):
    list_display = ('cid', 'schema', 'created_at')
    list_display_links = ('cid',)
    search_fields = ('cid__host', 'schema')
    list_per_page = 10


@admin.register(models.DbEnvironment)
class DbEnvironmentAdmin(admin.ModelAdmin):
    list_display = ('id', 'name', 'created_at', 'updated_at')
    list_display_links = ('name',)
    search_fields = ('name',)
    list_per_page = 10


@admin.register(models.DbOrders)
class DbOrdersAdmin(admin.ModelAdmin):
    list_display = ('title', 'order_id', 'applicant', 'progress', 'rds_category', 'sql_type', 'env', 'created_at')
    list_display_links = ('title',)
    list_filter = ('progress', 'env')
    search_fields = ('title', 'order_id', 'applicant')
    list_per_page = 20
    exclude = ['executor', 'auditor', 'closer', 'reviewer', 'email_cc']
    readonly_fields = ['is_hide', 'demand', 'remark', 'rds_category', 'sql_type', 'env',
                       'applicant', 'cid', 'database', 'version', 'order_id', 'title']


@admin.register(models.DbOrdersExecuteTasks)
class DbOrdersExecuteTasksAdmin(admin.ModelAdmin):
    list_display = ('task_id', 'applicant', 'order', 'sql_type', 'progress', 'created_time')
    list_display_links = ('task_id',)
    search_fields = ('taskid', 'order__order_id', 'applicant')
    list_per_page = 20
    exclude = ['rollback_sql', 'execute_log', 'executor']
    readonly_fields = ['consuming_time', 'affected_rows', 'file_format', 'task_id', 'applicant', 'order']
