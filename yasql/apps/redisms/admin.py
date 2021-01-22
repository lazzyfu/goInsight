from django.contrib import admin
from redisms import models
from django import forms
from libs.SecurityTools import KMS


class RedisConfigForm(forms.ModelForm):
    raw_password = forms.CharField(
        label='输入明文密码(转成秘文保存)',
        widget=forms.TextInput(attrs={'size': 65}),
        required=False
    )

    class Meta:
        model = models.RedisConfig
        fields = '__all__'


class RedisGroupAdmin(admin.ModelAdmin):
    list_display = ('id', 'name', 'env',)
    list_display_links = ('id',)
    list_filter = ('name', )
    readonly_fields = ()
    search_fields = ('name',)


class RedisConfigAdmin(admin.ModelAdmin):
    form = RedisConfigForm
    list_display = ('id', 'host', 'port')
    list_display_links = ('id',)
    list_filter = ('host', )
    readonly_fields = ('password',)
    search_fields = ('host',)

    def save_model(self, request, obj, form, change):
        raw_password = form.cleaned_data['raw_password']
        if raw_password:
            obj.password = KMS().encrypt(raw_password)
        super().save_model(request, obj, form, change)


class RedisGrantAdmin(admin.ModelAdmin):
    list_display = ('id', 'user', 'group', 'permission_code',)
    list_display_links = ()
    list_filter = ('id',)
    readonly_fields = ()
    search_fields = ('group',)


class RedisLogAdmin(admin.ModelAdmin):
    list_display = ('id', 'user', 'host_dsn', 'exec_cmd', 'created_at')
    list_display_links = ('id',)
    list_filter = ('user', 'created_at',)
    readonly_fields = ()
    search_fields = ('exec_cmd',)


admin.site.register(models.RedisEnv)
admin.site.register(models.RedisGroup, RedisGroupAdmin)
admin.site.register(models.RedisConfig, RedisConfigAdmin)
admin.site.register(models.RedisGrant, RedisGrantAdmin)
admin.site.register(models.RedisLog, RedisLogAdmin)
