from django.contrib import admin

# Register your models here.
from django.contrib.auth.hashers import make_password
from django.contrib.auth.models import Group

from users import models


class UserRolesInline(admin.TabularInline):
    model = models.UserRoles.user.through
    verbose_name = u"用户角色"
    verbose_name_plural = u"用户角色"
    max_num = 1


class RolePermissionInline(admin.TabularInline):
    model = models.RolePermissions.role.through
    verbose_name = u'用户权限'
    verbose_name_plural = u'用户权限'
    extra = 1


class UserAccountsAdmin(admin.ModelAdmin):
    list_display = (
        'username', 'user_role', 'displayname', 'email', 'mobile', 'is_active', 'date_joined'
    )
    list_display_links = ('username',)
    search_fields = ('username', 'email', 'displayname')
    fieldsets = (
        ('个人信息',
         {'fields': [
             'username', 'displayname', 'email', 'mobile', 'password', 'is_active', 'avatar_file', 'is_superuser'
         ]}
         ),
    )
    inlines = [UserRolesInline, ]
    exclude = ('users',)

    # 支持密码修改
    # 直接在密码输入框输入明文保存即可
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

    def has_delete_permission(self, request, obj=None):
        return False

    def has_add_permission(self, request):
        return False

    def has_change_permission(self, request, obj=None):
        return False


admin.site.register(models.UserAccounts, UserAccountsAdmin)
admin.site.register(models.UserRoles, UserRolesAdmin)
admin.site.register(models.RolePermissions, RolePermissionAdmin)
admin.site.unregister(Group)
