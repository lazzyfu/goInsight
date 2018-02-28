from django.contrib import admin

# Register your models here.
from ProjectManager.models import InceptionHostConfig, Remark, InceptionHostConfigDetail


class InceptionHostConfigDetailInline(admin.StackedInline):
    model = InceptionHostConfigDetail
    extra = 1


class InceptionHostConfigAdmin(admin.ModelAdmin):
    list_display = ('id', 'user', 'password', 'host', 'port', 'type', 'is_enable', 'group_name', 'comment')
    list_display_links = ('user',)
    search_fields = ('user',)
    ordering = ('id',)

    inlines = [InceptionHostConfigDetailInline, ]


@admin.register(Remark)
class RemarkAdmin(admin.ModelAdmin):
    list_display = ('id', 'remark', 'created_at', 'updated_at')


admin.site.register(InceptionHostConfig, InceptionHostConfigAdmin)
