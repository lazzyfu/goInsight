from django.contrib import admin
from django.contrib.auth.models import Group

from .models import UserAccount, Groups, GroupsDetail, Contacts, Roles, RolesDetail, ContactsDetail

# Register your models here.

admin.site.site_title = '后台'
admin.site.site_header = '数据库审核系统'

# 不注册系统的Group
admin.site.unregister(Group)


@admin.register(UserAccount)
class UserAccountAdmin(admin.ModelAdmin):
    list_display = ('uid', 'username', 'is_superuser', 'email', 'avatar_file', 'date_joined')
    list_display_links = ('username',)
    search_fields = ('username',)


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
    list_display = ('id', 'contact', 'group', 'created_at', 'updated_at')
