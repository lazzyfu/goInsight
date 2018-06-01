from django.contrib import admin, messages

# Register your models here.
from project_manager.models import InceptionHostConfig, InceptionHostConfigDetail, DomainName


class InceptionHostConfigDetailInline(admin.StackedInline):
    model = InceptionHostConfigDetail
    extra = 1


class InceptionHostConfigAdmin(admin.ModelAdmin):
    list_display = ('id', 'user', 'password', 'host', 'port', 'type', 'purpose', 'is_enable', 'group_name', 'comment')
    list_display_links = ('user',)
    search_fields = ('user',)
    ordering = ('id',)

    inlines = [InceptionHostConfigDetailInline, ]

admin.site.register(InceptionHostConfig, InceptionHostConfigAdmin)
