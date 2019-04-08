from django.contrib import admin

# Register your models here.
from query.models import MysqlPrivBlacklist, MysqlUserGroupMap, MySQLQueryLog, QueryBusinessGroup


@admin.register(MySQLQueryLog)
class MySQLQueryLogAdmin(admin.ModelAdmin):
    list_display = ('user', 'host', 'database', 'query_sql', 'query_time', 'affect_rows', 'created_at')
    list_display_links = ('query_sql',)
    readonly_fields = ('user', 'host', 'query_sql', 'query_time', 'database', 'affect_rows')
    search_fields = ('query_sql', 'user', 'database')
    list_filter = ('user', 'created_at')


@admin.register(MysqlPrivBlacklist)
class MysqlPrivBlacklistAdmin(admin.ModelAdmin):
    list_display = ('comment', 'schema', 'table', 'columns', 'updated_at', 'created_at')
    list_display_links = ('comment',)
    ordering = ('-updated_at',)
    search_fields = ('schema', 'table', 'columns', 'comment')
    list_per_page = 20
    list_filter = ('schema',)


@admin.register(MysqlUserGroupMap)
class MysqlUserGroupMapAdmin(admin.ModelAdmin):
    list_display = ('comment', 'schema', 'group_alias', 'group', 'updated_at', 'created_at')
    list_display_links = ('comment',)
    ordering = ('id',)
    search_fields = ('group', 'user__username', 'schema', 'comment__comment')
    filter_horizontal = ('user',)
    list_per_page = 20
    list_filter = ('schema', 'comment',)

    readonly_fields = ('comment', 'group', 'schema')

    fieldsets = (
        ('',
         {'fields': ['comment', 'group', 'schema']}
         ),
        ('',
         {'fields': ['user']}
         )
    )

    def has_add_permission(self, request):
        return


@admin.register(QueryBusinessGroup)
class QueryBusinessGroupAdmin(admin.ModelAdmin):
    list_display = ('group', 'config', 'schema', 'tables', 'updated_at', 'created_at')
    list_display_links = ('group',)
    ordering = ('id',)
    search_fields = ('group', 'schema', 'tables', 'user__username')
    filter_horizontal = ('user',)
    list_per_page = 20
    list_filter = ('schema', 'group',)

    fieldsets = (
        ('',
         {'fields': ['group', 'config', 'schema', 'tables', 'map_mysqluser']}
         ),
        ('',
         {'fields': ['user']}
         )
    )
