from django.contrib import admin

from webshell.models import WebShellInfo


# Register your models here.


@admin.register(WebShellInfo)
class WebShellInfoAdmin(admin.ModelAdmin):
    list_display = ('id', 'comment', 'command', 'envi_id', 'created_at')
    list_display_links = ('comment',)
    ordering = ('-created_at',)

