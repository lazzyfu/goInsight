from django.contrib import admin
from django.contrib.auth.hashers import make_password
from django.contrib.auth.models import Group

from .models import UserAccount, Groups, GroupsDetail, Contacts, Roles, RolesDetail, ContactsDetail, PermissionDetail

# Register your models here.

admin.site.site_title = '后台'
admin.site.site_header = '数据库审核系统'

# 不注册系统的Group
admin.site.unregister(Group)


class RolesDetailInline(admin.StackedInline):
    model = RolesDetail
    max_num = 1


class GroupsDetailInline(admin.StackedInline):
    model = GroupsDetail
    extra = 1


class PermissionDetailInline(admin.StackedInline):
    model = PermissionDetail
    extra = 1


class ContactsDetailInline(admin.StackedInline):
    model = ContactsDetail
    extra = 1


class UserAccountAdmin(admin.ModelAdmin):
    list_display = (
        'uid', 'username', 'displayname', 'is_superuser', 'is_active', 'email', 'avatar_file', 'user_role',
        'user_group',
        'date_joined')
    list_display_links = ('username',)
    search_fields = ('username', 'email', 'displayname', 'user_group')
    fieldsets = (
        ('个人信息', {'fields': ['username', 'displayname', 'email', 'is_superuser', 'is_active', 'avatar_file']}),
    )
    inlines = [RolesDetailInline, GroupsDetailInline]

    actions = ['reset_password']

    # 重置密码
    def reset_password(self, request, queryset):
        rows_updated = queryset.update(password=make_password('123.com'))
        if rows_updated == 1:
            message_bit = "1 user was"
        else:
            message_bit = "%s users were" % rows_updated
        self.message_user(request, "%s successfully reset password, password is: 123.com" % message_bit)

    reset_password.short_description = u'重置用户密码为：123.com'


class RolesAdmin(admin.ModelAdmin):
    list_display = ('role_id', 'role_name', 'permission', 'created_at', 'updated_at')
    ordering = ('-created_at',)
    list_display_links = ('role_name',)
    inlines = [PermissionDetailInline, ]


class ContactsAdmin(admin.ModelAdmin):
    list_display = ('contact_id', 'contact_name', 'contact_email', 'contact_group', 'created_at', 'updated_at')
    ordering = ('-created_at',)
    list_display_links = ('contact_email',)

    inlines = [ContactsDetailInline, ]


@admin.register(Groups)
class GroupsAdmin(admin.ModelAdmin):
    list_display = ('group_id', 'group_name', 'created_at', 'updated_at')
    ordering = ('-created_at',)


# 注册
admin.site.register(UserAccount, UserAccountAdmin)
admin.site.register(Roles, RolesAdmin)
admin.site.register(Contacts, ContactsAdmin)
