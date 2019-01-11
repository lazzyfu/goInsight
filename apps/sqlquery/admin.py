from django.contrib import admin
from django.contrib import messages

# Register your models here.
from sqlquery.models import MysqlRulesGroup, MysqlRulesChain
from sqlquery.tasks import mysql_privileges_operate
from . import models


class MySQLQueryLogAdmin(admin.ModelAdmin):
    list_display = ('user', 'host', 'database', 'query_sql', 'query_time', 'affect_rows', 'created_at')
    list_display_links = ('query_sql',)
    readonly_fields = ('user', 'host', 'query_sql', 'query_time', 'database', 'affect_rows')
    search_fields = ('query_sql', 'user', 'database')
    list_filter = ('user', 'created_at')


class MysqlRulesChainAdmin(admin.ModelAdmin):
    list_display = ('comment', 'cid', 'action', 'schema', 'table', 'created_at')
    list_display_links = ('comment',)
    ordering = ('id',)
    search_fields = ('action', 'schema', 'table', 'comment')

    readonly_fields = ('comment', 'action', 'cid', 'schema', 'table', 'created_at')

    def get_actions(self, request):
        actions = super().get_actions(request)
        if 'delete_selected' in actions:
            del actions['delete_selected']
        return actions

    def has_add_permission(self, request):
        return False

    def has_delete_permission(self, request, obj=None):
        return False

    def has_change_permission(self, request, obj=None):
        return False


class MysqlRulesGroupAdmin(admin.ModelAdmin):
    list_display = ('name', 'created_at')
    list_display_links = ('name',)
    ordering = ('id',)
    filter_horizontal = ('user', 'rule')
    search_fields = ('user__username', 'rule__comment')

    def save_model(self, request, obj, form, change):
        data = form.clean()
        before_rule_id = []
        if form.initial.get('rule'):
            before_rule_id = [x.id for x in form.initial.get('rule')]

        before_users = []
        if form.initial.get('user'):
            before_users = [x.username for x in form.initial.get('user')]

        after_rule_id = []
        if data.get('rule'):
            after_rule_id = [x.id for x in data.get('rule')]

        after_users = []
        if data.get('user'):
            after_users = [x.username for x in data.get('user')]

        mysql_privileges_operate.apply_async((before_rule_id,
                                              after_rule_id,
                                              before_users,
                                              after_users),
                                             countdown=1)
        messages.add_message(request, messages.SUCCESS, f'任务已提交到后台处理，请稍后几分钟检测权限是否变更')
        return super().save_model(request, obj, form, change)


admin.site.register(models.MySQLQueryLog, MySQLQueryLogAdmin)
admin.site.register(MysqlRulesChain, MysqlRulesChainAdmin)
admin.site.register(MysqlRulesGroup, MysqlRulesGroupAdmin)
