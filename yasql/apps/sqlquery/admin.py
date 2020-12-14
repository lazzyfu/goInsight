from django import forms
from django.contrib import admin
from django.contrib.admin.widgets import AutocompleteSelect

from sqlquery import models


# Register your models here.


class TableSearchForm(forms.ModelForm):
    class Meta:
        fields = ('tables',)
        widgets = {
            'tables': AutocompleteSelect(
                models.DbQueryUserAllowedTables._meta.get_field('tables').remote_field,
                admin.site,
                attrs={'data-dropdown-auto-width': 'true'}
            ),
        }


class DbQueryUserAllowedTablesInline(admin.StackedInline):
    model = models.DbQueryUserAllowedTables
    autocomplete_fields = ('tables',)
    extra = 1

    fieldsets = (
        (None, {
            'classes': ('collapse',),
            'fields': ('tables',),
        }),
    )

    form = TableSearchForm


class DbQueryUserDenyTablesInline(admin.StackedInline):
    model = models.DbQueryUserDenyTables
    autocomplete_fields = ('tables',)
    extra = 1

    fieldsets = (
        (None, {
            'classes': ('collapse',),
            'fields': ('tables',),
        }),
    )

    form = TableSearchForm


class DbQueryGroupAllowedTablesInline(admin.StackedInline):
    model = models.DbQueryGroupAllowedTables
    autocomplete_fields = ('tables',)
    extra = 1

    fieldsets = (
        (None, {
            'classes': ('collapse',),
            'fields': ('tables',),
        }),
    )

    form = TableSearchForm


class DbQueryGroupDenyTablesInline(admin.StackedInline):
    model = models.DbQueryGroupDenyTables
    autocomplete_fields = ('tables',)
    extra = 1

    fieldsets = (
        (None, {
            'classes': ('collapse',),
            'fields': ('tables',),
        }),
    )

    form = TableSearchForm


@admin.register(models.DbQueryTables)
class DbQueryTablesAdmin(admin.ModelAdmin):
    list_display = ('display_comment', 'display_schema', 'table')
    search_fields = ('table', 'schema__schema', 'schema__cid__comment')
    list_per_page = 30

    def has_add_permission(self, request):
        return False

    def has_change_permission(self, request, obj=None):
        return False


@admin.register(models.DbQueryUserPrivs)
class DbQueryUserPrivsAdmin(admin.ModelAdmin):
    list_display = ('user',)
    list_display_links = ('user',)
    list_per_page = 20
    autocomplete_fields = ('user',)
    filter_horizontal = ('schemas',)
    inlines = [DbQueryUserAllowedTablesInline, DbQueryUserDenyTablesInline]

    search_fields = ('user__username', 'schemas__schema')


@admin.register(models.DbQueryGroupPrivs)
class DbQueryGroupPrivsAdmin(admin.ModelAdmin):
    list_display = ('group',)
    list_display_links = ('group',)
    list_per_page = 20
    filter_horizontal = ('user', 'schemas')
    inlines = [DbQueryGroupAllowedTablesInline, DbQueryGroupDenyTablesInline]

    search_fields = ('group', 'user__username', 'schemas__schema')
