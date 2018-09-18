# -*- coding:utf-8 -*-
# edit by fuzongfei
from django import forms
from django.core.exceptions import PermissionDenied
from django.db.models import Min, F

from sqlorders.models import envi_choice, SqlOrdersEnvironment, MysqlSchemas
from sqlorders.utils import GetTableInfo
from sqlquery.models import MysqlSchemasGrant
from sqlquery.sqlQueryApi import MySQLQuery


class GetSchemasGrantForm(forms.Form):
    envi_id = forms.ChoiceField(required=True, choices=envi_choice, label=u'环境id')

    def query(self, request):
        cdata = self.cleaned_data
        envi_id = cdata.get('envi_id')
        is_master = 1

        # 判断是否为只读环境，查找parent_id最小的，即认为为只读环境（一般为生产环境的从库）
        parent_id_min = SqlOrdersEnvironment.objects.all().aggregate(Min('parent_id'))['parent_id__min']
        if int(envi_id) == SqlOrdersEnvironment.objects.get(parent_id=parent_id_min).envi_id:
            is_master = 0

        query = f"select b.id, b.host, b.port, b.schema from sqlaudit_schemas_grant a " \
                f"join sqlaudit_mysql_schemas b on a.schema_id = b.schema_join join sqlaudit_user_accounts c  " \
                f"on c.uid = a.user_id where c.uid={request.user.uid} " \
                f"and b.envi_id={envi_id} and b.is_master={is_master}"

        context = []
        for row in MysqlSchemas.objects.raw(query):
            data = GetTableInfo(row.host, row.port, row.schema).get_online_tables()
            show_schema = '_'.join((row.comment, row.schema))
            context.append({
                'id': '___'.join((row.host, str(row.port), row.schema)),
                'text': show_schema,
                'children': data
            })
        return context


class GetStruInfoForm(forms.Form):
    schema = forms.CharField(max_length=1024, required=True)

    def query(self):
        cdata = self.cleaned_data
        host, port, schema = cdata.get('schema').split('___')
        if len(schema.split('.')) == 2:
            data = GetTableInfo(host, port, schema).get_stru_info()
            context = {'status': 0, 'data': data}
        else:
            context = {'status': 2, 'msg': ''}
        return context


class ExecSqlQueryForm(forms.Form):
    envi_id = forms.ChoiceField(required=True, choices=envi_choice, label=u'环境id')
    schema = forms.CharField(max_length=1024)
    contents = forms.CharField(widget=forms.Textarea, label=u'sql')

    def execute(self, request):
        cdata = self.cleaned_data
        envi_id = cdata.get('envi_id')
        schema = cdata.get('schema')
        contents = cdata.get('contents')

        host, port, schema = schema.split('___')
        if len(schema.split('.')) == 2:
            schema = schema.split('.')[0]

        # 判断主机所属的envi_id是否和接收的envi_id相等
        if int(envi_id) == MysqlSchemas.objects.filter(host=host).first().envi_id:
            # 验证传入的host是否属于该用户的授权主机
            schemas = MysqlSchemasGrant.objects.filter(user__uid=request.user.uid).annotate(
                schemas=F('schema__schema')).values_list('schemas', flat=True)
            if schema in schemas:
                # 判断是否是只读
                # 判断依据：parent_id最小的
                parent_id_min = SqlOrdersEnvironment.objects.all().aggregate(Min('parent_id'))['parent_id__min']
                if int(envi_id) == SqlOrdersEnvironment.objects.get(parent_id=parent_id_min).envi_id:
                    mysql_query = MySQLQuery(querys=contents, host=host, port=port, schema=schema, rw='r', envi=0)
                else:
                    mysql_query = MySQLQuery(querys=contents, host=host, port=port, schema=schema, rw='rw', envi=0)
                result = mysql_query.query(request)
            else:
                raise PermissionDenied
        else:
            raise PermissionDenied

        return result
