# -*- coding:utf-8 -*-
# edit by fuzongfei
import logging
from datetime import datetime

import pymysql
import sqlparse
from django_redis import get_redis_connection
from rest_framework import serializers

from config import QUERY_USER
from sqlorders.models import DbConfig
from sqlquery.api.getTablesMeta import GetTablesMeta
from sqlquery.api.sqlQuery import SqlQuery
from sqlquery.api.verifyUserPrivs import VerifyUserPrivs
from sqlquery.models import DbQueryUserPrivs, DbQueryGroupPrivs, DbQueryLog
from sqlquery.utils import remove_sql_comment

logger = logging.getLogger('main')


class ExecuteQuerySerializer(serializers.Serializer):
    key = serializers.CharField(max_length=256, error_messages={'blank': '请先选中左侧的一个库名'})
    sql = serializers.CharField()
    character = serializers.ChoiceField(
        choices=(
            ('utf8', 'utf8'),
            ('utf8mb4', 'utf8mb4'),
            ('latin1', 'latin1'),
            ('gbk', 'gbk'))
    )
    query_hash = serializers.CharField(max_length=64, min_length=64)

    def validate_sql(self, value):
        if len(sqlparse.split(value)) >= 2:
            raise serializers.ValidationError("每次请执行一条SQL语句")

        # 移除sql注释和尾部的分号
        return remove_sql_comment(value)

    def validate(self, attrs):
        """检查用户是否有库表查询权限"""
        username = self.context['request'].user.username
        key = attrs.get('key')
        sql = attrs.get('sql')
        status, _ = VerifyUserPrivs(username=username, key=key, sql=sql).run
        if not status:
            raise serializers.ValidationError(_)
        attrs['tables'] = _
        return super(ExecuteQuerySerializer, self).validate(attrs)

    def execute(self):
        vdata = self.validated_data
        key = vdata.get('key')
        sql = vdata.get('sql')
        character = vdata.get('character')
        query_hash = vdata.get('query_hash')
        username = self.context['request'].user.username

        cid, schema = key.split('___')
        obj = DbConfig.objects.get(pk=cid)
        config = {
            'host': obj.host,
            'port': obj.port,
            'charset': character
        }
        config.update(QUERY_USER)

        kwargs = {
            'username': username,
            'sql': sql,
            'query_hash': query_hash,
            'schema': schema,
            'tables': vdata.get('tables'),
            'rds_category': obj.rds_category,
            'config': config
        }
        return SqlQuery(kwargs=kwargs).execute()


class GetTreeSerializer(serializers.Serializer):
    key = serializers.CharField(max_length=256)

    def query(self, request):
        vdata = self.validated_data
        key = vdata.get('key')

        tree_data = []

        # root节点，获取库名
        if key == 'root':
            obj_user = DbQueryUserPrivs.objects.filter(
                user__username=request.user.username
            ).values_list('schemas__cid',
                          'schemas__cid__comment',
                          'schemas__schema',
                          'schemas__cid__host',
                          'schemas__cid__port'
                          )

            obj_group = DbQueryGroupPrivs.objects.filter(
                user__username=request.user.username
            ).values_list('schemas__cid',
                          'schemas__cid__comment',
                          'schemas__schema',
                          'schemas__cid__host',
                          'schemas__cid__port'
                          )
            for row in set(list(obj_user) + list(obj_group)):
                row = {
                    'title': '___'.join(row[1:3]),
                    'key': '___'.join([str(row[0]), row[2]]),
                    'hostname': ':'.join([row[3], str(row[4])])
                }
                tree_data.append(row)
            return True, tree_data

        # 二级节点，获取表名和列名
        if len(key.split('___')) == 2:
            cid, schema = key.split('___')
            obj = DbConfig.objects.get(pk=cid)
            config = {
                'host': obj.host,
                'port': obj.port
            }
            config.update(QUERY_USER)

            return GetTablesMeta(config=config, cid=cid).get_tables_treedata(schema=schema)


class DeleteQueryHashSerializer(serializers.Serializer):
    query_hash = serializers.CharField(max_length=64, min_length=64)

    def del_session(self):
        username = self.context['request'].user.username
        query_hash = self.validated_data.get('query_hash')
        cnx_redis = get_redis_connection('default')
        data = cnx_redis.smembers(query_hash)
        if data:
            thread_id, host, port = data.pop().split(':')
            config = {
                'host': host,
                'port': int(port),
                'cursorclass': pymysql.cursors.DictCursor
            }
            config.update(QUERY_USER)
            cnx = pymysql.connect(**config)
            try:
                with cnx.cursor() as cursor:
                    cursor.execute("select version() as version")
                    server_version = cursor.fetchone()['version']
                    # 如果是tidb，执行kill tidb thread_id
                    if 'TIDB' in server_version.upper():
                        cursor.execute(f'KILL TIDB {int(thread_id)}')
                        logger.info(f"KILL TIDB QUERY,Thread ID:{thread_id},"
                                    f"HOST:{host}:{port}, USER:{username}")
                    else:
                        cursor.execute(f'KILL CONNECTION {int(thread_id)}')
                        logger.info(f"KILL MySQL QUERY,Thread ID:{thread_id},"
                                    f"HOST:{host}:{port}, USER:{username}")
                    cursor.close()
                cnx_redis.srem(query_hash, data)
            except Exception as err:
                logger.error(err)


class GetTableInfoSerializer(serializers.Serializer):
    key = serializers.CharField(max_length=256)
    type = serializers.ChoiceField(choices=(('table_structure', 'table_structure'), ('table_base', 'table_base')))

    def query(self):
        vdata = self.validated_data
        key = vdata.get('key')
        type = vdata.get('type')
        cid, schema, table = key.split('___')
        obj = DbConfig.objects.get(pk=cid)
        config = {
            'host': obj.host,
            'port': obj.port,
        }
        config.update(QUERY_USER)
        if type == 'table_structure':
            return GetTablesMeta(config=config, cid=cid).get_table_structure(schema=schema, table=table)

        if type == 'table_base':
            return GetTablesMeta(config=config, cid=cid).get_table_base(schema=schema, table=table)


class GetHistorySQLSerializer(serializers.ModelSerializer):
    class Meta:
        fields = '__all__'
        model = DbQueryLog

    def to_representation(self, instance):
        ret = super(GetHistorySQLSerializer, self).to_representation(instance)
        ret["created_at"] = datetime.strftime(instance.created_at, "%Y-%m-%d %H:%M")
        return ret


class GetDBDictSerializer(serializers.Serializer):
    key = serializers.CharField(max_length=256)

    def query(self):
        vdata = self.validated_data
        key = vdata.get('key')
        cid, schema = key.split('___')
        obj = DbConfig.objects.get(pk=cid)
        config = {
            'host': obj.host,
            'port': obj.port,
            'db': 'INFORMATION_SCHEMA'
        }
        config.update(QUERY_USER)
        return GetTablesMeta(config=config, cid=cid).get_db_dict(schema=schema)
