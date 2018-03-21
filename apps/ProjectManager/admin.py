from django.contrib import admin, messages

# Register your models here.
from ProjectManager.models import InceptionHostConfig, Remark, InceptionHostConfigDetail
from ProjectManager.utils import check_mysql_conn


class InceptionHostConfigDetailInline(admin.StackedInline):
    model = InceptionHostConfigDetail
    extra = 1


class InceptionHostConfigAdmin(admin.ModelAdmin):
    list_display = ('id', 'user', 'password', 'host', 'port', 'type', 'is_enable', 'group_name', 'comment')
    list_display_links = ('user',)
    search_fields = ('user',)
    ordering = ('id',)

    inlines = [InceptionHostConfigDetailInline, ]

    actions = ['check_connection_status']

    def check_connection_status(self, request, queryset):
        for row in queryset.filter():
            user = row.user
            host = row.host
            password = row.password
            port = row.port
            result = check_mysql_conn(user, host, password, port)
            status = result['status']
            message_bit = result['msg']
            if status == 'INFO':
                self.message_user(request, f'{host}: {message_bit}', level=messages.INFO)
            elif status == 'ERROR':
                self.message_user(request, f'{host}: {message_bit}', level=messages.ERROR)

    check_connection_status.short_description = u'测试账号到数据库的连接'

@admin.register(Remark)
class RemarkAdmin(admin.ModelAdmin):
    list_display = ('id', 'remark', 'created_at', 'updated_at')


admin.site.register(InceptionHostConfig, InceptionHostConfigAdmin)
