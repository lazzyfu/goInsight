# -*- coding:utf-8 -*-
# edit by fuzongfei
import sqlparse
from django.db.models import F
from rest_framework import serializers

from orders.models import MysqlSchemas, SysEnvironment, Orders, OnlineVersion
from orders.utils.inceptionApi import InceptionSqlApi
from orders.utils.tools import split_sqltype
from orders.utils.validatorsCheck import envi_validator, sql_type_validator, online_version_validator


class BeautifySQLSerializer(serializers.Serializer):
    """
    注释格式必须符合规范即可
    格式：# 这是注释 中间要有空格
    """
    contents = serializers.CharField(allow_blank=False, error_messages={'blank': '格式化的SQL不能为空'})

    def beautify(self):
        sdata = self.validated_data
        contents = sdata.get('contents')

        split_sqls = []
        for stmt in sqlparse.split(contents):
            sql = sqlparse.parse(stmt)[0]
            sql_comment = sql.token_first()
            if isinstance(sql_comment, sqlparse.sql.Comment):
                split_sqls.append({'comment': sql_comment.value, 'sql': sql.value.replace(sql_comment.value, '')})
            else:
                split_sqls.append({'comment': '', 'sql': sql.value})

        beautify_sqls = []
        for row in split_sqls:
            comment = row['comment']
            sql = row['sql']
            res = sqlparse.parse(sql)
            syntax_type = res[0].token_first().ttype.__str__()
            if syntax_type == 'Token.Keyword.DDL':
                sql_format = sqlparse.format(sql)
                beautify_sqls.append(comment + sql_format)
            elif syntax_type == 'Token.Keyword.DML':
                sql_format = sqlparse.format(sql, strip_whitespace=True, reindent=True)
                beautify_sqls.append(comment + sql_format)
            else:
                beautify_sqls.append(comment + sql)
        return '\n\n'.join(beautify_sqls)


class SyntaxCheckSerializer(serializers.Serializer):
    schema = serializers.CharField(required=True, max_length=256, error_messages={'required': 'schema不能为空'})
    sql_type = serializers.CharField(validators=[sql_type_validator], required=True,
                                     error_messages={'required': 'SQL类型不能为空'})
    contents = serializers.CharField(allow_blank=False, error_messages={'required': '语法检查输入的SQL不能为空'})

    def check(self, request):
        sdata = self.validated_data
        host, port, database = sdata.get('schema').split(',')
        sql_type = sdata.get('sql_type')
        contents = sdata.get('contents')

        # 对检测的SQL类型进行区分
        status, msg = split_sqltype(contents, sql_type)

        # 实例化
        inception_check = InceptionSqlApi(host, port, database, contents, request.user.username)

        if not status:
            return False, msg

        if status:
            # SQL语法检查
            context = inception_check.run_check()
            return True, context['data']


class GetSchemasSerializer(serializers.Serializer):
    envi_id = serializers.CharField(required=True, validators=[envi_validator],
                                    error_messages={'required': 'envi_id字段不能为空'})

    def query(self):
        sdata = self.validated_data
        envi_id = sdata.get('envi_id')
        queryset = MysqlSchemas.objects.filter(envi__envi_id=envi_id, cid__type=0).annotate(
            host=F('cid__host'),
            port=F('cid__port'),
            comment=F('cid__comment')
        ).values('host', 'port', 'schema', 'comment')
        return queryset


class OnlineVersionDetailSerializer(serializers.Serializer):
    version = serializers.CharField(required=True, validators=[online_version_validator],
                                    error_messages={'required': '上线版本号不能为空'})

    def query(self):
        sdata = self.validated_data
        version = sdata.get('version')

        queryset = SysEnvironment.objects.values('envi_id', 'envi_name').order_by('-envi_id')
        dynamic_columns_join = ''
        for row in queryset:
            dynamic_columns_join += f"max(if(envi_id={row['envi_id']}, progress, -1)) as {row['envi_name']},"

        # 获取任务下所有工单分别在各个环境中的状态
        # 此处的环境为动态环境
        version_id = OnlineVersion.objects.get(version=version).id
        query = f"select " + dynamic_columns_join + \
                "id, title, applicant " \
                    f"from auditsql_orders where version_id={version_id} " \
                    f"group by title,applicant order by id desc"
        result = []
        print(query)

        data = Orders.objects.raw(query)
        print(data.columns)
        dynamic_columns = list(data.columns)[:-3]

        # 获取列名并进行拼接
        columns_definition = [{'field': 'id', 'title': 'ID', 'visible': False},
                              {'field': 'title', 'title': '标题'},
                              {'field': 'applicant', 'title': '申请人'},
                              {'field': 'version', 'title': '上线版本号'},
                              ]

        dynamic_columns_definition = [{'field': x, 'title': x, 'formatter': 'render_onlinetasks_status'} for x in
                                      dynamic_columns]

        # 获取列名对应的数据
        for row in data:
            columns = {
                'id': row.id,
                'title': row.title,
                'applicant': row.applicant,
                'auditor': row.auditor,
                'version': version,
            }
            for i in dynamic_columns:
                columns[i] = getattr(row, i)
            result.append(columns)

        print(result)
        return {'columns': columns_definition + dynamic_columns_definition, 'data': result}
