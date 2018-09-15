from django.contrib import admin
# Register your models here.
from django.contrib.auth.hashers import make_password
from django.contrib.auth.models import Group
from django.db.models import Min, Q

from sqlorders.models import SqlOrdersEnvironment, MysqlSchemas
from sqlquery.models import MysqlSchemasGrant
from users.models import UserAccounts, UserRoles, RolePermission
from webshell.models import WebShellGrant

admin.site.site_title = '后台'
admin.site.site_header = '数据库审核系统'


class UserRolesInline(admin.TabularInline):
    model = UserRoles.user.through
    verbose_name = u"用户角色"
    verbose_name_plural = u"用户角色"
    max_num = 1


class RolePermissionInline(admin.TabularInline):
    model = RolePermission.role.through
    verbose_name = u'用户权限'
    verbose_name_plural = u'用户权限'
    extra = 1


class WebShellGrantInline(admin.TabularInline):
    model = WebShellGrant
    extra = 1

    verbose_name = u'WebShell授权'
    verbose_name_plural = u'WebShell授权'


class MysqlSchemasGrantInline(admin.TabularInline):
    model = MysqlSchemasGrant
    extra = 1

    verbose_name = u'授权库'
    verbose_name_plural = u'授权库'

    # 此处过滤指定的数据
    # 将生产环境的主库过滤掉
    def formfield_for_foreignkey(self, db_field, request, **kwargs):
        try:
            parent_id_min = SqlOrdersEnvironment.objects.all().aggregate(Min('parent_id'))['parent_id__min']
            envi_id = SqlOrdersEnvironment.objects.get(parent_id=parent_id_min).envi_id
            kwargs['queryset'] = MysqlSchemas.objects.filter(Q(envi_id=envi_id, is_master=0) | ~Q(envi_id=envi_id))
        except Exception as err:
            pass
        finally:
            return super(MysqlSchemasGrantInline, self).formfield_for_foreignkey(db_field, request, **kwargs)


class UserAccountsAdmin(admin.ModelAdmin):
    list_display = (
        'username', 'displayname', 'email', 'mobile', 'is_active', 'user_role', 'user_schema',
        'date_joined')
    list_display_links = ('username',)
    search_fields = ('username', 'email', 'displayname')
    fieldsets = (
        ('个人信息',
         {'fields': ['username', 'displayname', 'email', 'mobile', 'password', 'is_active', 'avatar_file']}),
    )
    inlines = [UserRolesInline, MysqlSchemasGrantInline, WebShellGrantInline]

    exclude = ('users',)
    actions = ['reset_password']

    # 支持密码修改
    def save_model(self, request, obj, form, change):
        obj.user = request.user
        data = form.clean()
        password = data.get('password')
        if 'password' in form.changed_data:
            obj.password = make_password(password)

        return super().save_model(request, obj, form, change)


class UserRolesAdmin(admin.ModelAdmin):
    list_display = ('role_name', 'permission', 'created_at', 'updated_at')
    ordering = ('-created_at',)
    fieldsets = (
        ('角色信息',
         {'fields': ['role_name']}),
    )
    list_display_links = ('role_name',)
    inlines = [RolePermissionInline, ]


class RolePermissionAdmin(admin.ModelAdmin):
    list_display = ('permission_name', 'permission_desc', 'created_at', 'updated_at')
    ordering = ('-created_at',)
    fieldsets = (
        ('权限信息',
         {'fields': ['permission_name', 'permission_desc']}),
    )
    list_display_links = ('permission_desc',)
    readonly_fields = ('permission_name', 'permission_desc')


admin.site.register(UserAccounts, UserAccountsAdmin)
admin.site.register(UserRoles, UserRolesAdmin)
admin.site.register(RolePermission, RolePermissionAdmin)
admin.site.unregister(Group)
