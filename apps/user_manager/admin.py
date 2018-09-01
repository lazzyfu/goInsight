from django.contrib import admin
from django.contrib.auth.hashers import make_password, check_password
from django.contrib.auth.models import Group

from mstats.models import MysqlSchemaGrant, MysqlSchemaInfo, WebShellGrant
from .models import UserAccount, Roles, RolesDetail, PermissionDetail, \
    RolePermission, SystemMsg

# Register your models here.

admin.site.site_title = '后台'
admin.site.site_header = '数据库审核系统'

# 不注册系统的Group
admin.site.unregister(Group)


class RolesDetailInline(admin.TabularInline):
    model = RolesDetail
    max_num = 1


class PermissionDetailInline(admin.StackedInline):
    model = PermissionDetail
    extra = 1


class MysqlSchemaGrantInline(admin.TabularInline):
    model = MysqlSchemaGrant
    extra = 1

    verbose_name = u'Database授权'
    verbose_name_plural = u'Database授权库授权'

    # 此处过滤指定的数据
    def formfield_for_foreignkey(self, db_field, request, **kwargs):
        kwargs['queryset'] = MysqlSchemaInfo.objects.filter(envi=3, is_master=0)
        return super(MysqlSchemaGrantInline, self).formfield_for_foreignkey(db_field, request, **kwargs)


class WebShellGrantInline(admin.TabularInline):
    model = WebShellGrant
    extra = 1

    verbose_name = u'WebShell授权'
    verbose_name_plural = u'WebShell授权'


class UserAccountAdmin(admin.ModelAdmin):
    list_display = (
        'fullname', 'mobile', 'user_role', 'user_schema', 'user_shell',
        'date_joined')
    list_display_links = ('fullname',)
    search_fields = ('username', 'email', 'displayname')
    fieldsets = (
        ('个人信息',
         {'fields': ['username', 'displayname', 'password', 'mobile', 'is_active', 'avatar_file']}),
    )
    inlines = [RolesDetailInline, MysqlSchemaGrantInline, WebShellGrantInline]

    actions = ['reset_password']

    # 重置密码功能
    def reset_password(self, request, queryset):
        rows_updated = queryset.update(password=make_password('123.com'))
        if rows_updated == 1:
            message_bit = "1 user was"
        else:
            message_bit = "%s users were" % rows_updated
        self.message_user(request, "%s successfully reset password, password is: 123.com" % message_bit)

    reset_password.short_description = u'重置用户密码为：123.com'

    # 支持密码修改
    def save_model(self, request, obj, form, change):
        obj.user = request.user
        data = form.clean()
        password = data.get('password')
        if 'password' in form.changed_data:
            obj.password = make_password(password)

        return super().save_model(request, obj, form, change)


class RolesAdmin(admin.ModelAdmin):
    list_display = ('role_id', 'role_name', 'permission', 'created_at', 'updated_at')
    ordering = ('-created_at',)
    list_display_links = ('role_name',)
    inlines = [PermissionDetailInline, ]


@admin.register(RolePermission)
class PermissionAdmin(admin.ModelAdmin):
    list_display = ('id', 'permission_name', 'permission_desc', 'created_at', 'updated_at')
    ordering = ('-created_at',)
    readonly_fields = ['permission_name', 'permission_desc']


@admin.register(SystemMsg)
class SystemMsgAdmin(admin.ModelAdmin):
    list_display = ('id', 'title', 'content', 'created_at')
    ordering = ('-created_at',)


# 注册
admin.site.register(UserAccount, UserAccountAdmin)
admin.site.register(Roles, RolesAdmin)
