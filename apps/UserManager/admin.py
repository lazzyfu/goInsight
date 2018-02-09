from django.contrib import admin
from django.contrib.auth.hashers import make_password
from django.contrib.auth.models import Group

from .models import UserAccount, Groups, GroupsDetail, Contacts, Roles, RolesDetail, ContactsDetail

# Register your models here.

admin.site.site_title = '后台'
admin.site.site_header = '数据库审核系统'

# 不注册系统的Group
admin.site.unregister(Group)


# @admin.register(UserAccount)
class UserAccountAdmin(admin.ModelAdmin):
    list_display = ('uid', 'username', 'displayname', 'is_superuser', 'email', 'avatar_file', 'date_joined')
    list_display_links = ('username',)
    search_fields = ('username',)

    actions = ['reset_password']

    # 重置密码
    def reset_password(self, request, queryset):
        rows_updated = queryset.update(password=make_password('123.com'))
        if rows_updated == 1:
            message_bit = "1 user was"
        else:
            message_bit = "%s users were" % rows_updated
        self.message_user(request, "%s successfully reset password, password is: 123.com" % message_bit)

@admin.register(Groups)
class GroupsAdmin(admin.ModelAdmin):
    list_display = ('group_id', 'group_name', 'created_at', 'updated_at')
    ordering = ('-created_at',)


@admin.register(GroupsDetail)
class GroupsDetailAdmin(admin.ModelAdmin):
    list_display = ('id', 'user', 'group', 'created_at', 'updated_at')
    ordering = ('-created_at',)


@admin.register(Roles)
class RolesAdmin(admin.ModelAdmin):
    list_display = ('role_id', 'role_name', 'created_at', 'updated_at')
    ordering = ('-created_at',)


@admin.register(RolesDetail)
class RolesDetailAdmin(admin.ModelAdmin):
    list_display = ('user', 'role', 'created_at', 'updated_at')
    ordering = ('-created_at',)


@admin.register(Contacts)
class ContactsAdmin(admin.ModelAdmin):
    list_display = ('contact_id', 'contact_name', 'contact_email', 'created_at', 'updated_at')
    ordering = ('-created_at',)


@admin.register(ContactsDetail)
class ContactsDetailAdmin(admin.ModelAdmin):
    list_display = ('id', 'contact', 'group', 'bcc', 'created_at', 'updated_at')


admin.site.register(UserAccount, UserAccountAdmin)
