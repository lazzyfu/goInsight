# -*- coding:utf-8 -*-
# edit by fuzongfei
import logging
import time

import pymysql
from celery import shared_task
from channels.layers import get_channel_layer
from django.core.cache import cache
from django.utils import timezone

from orders.api.executeStatementApi import ExecuteSql
from orders.models import OrdersTasks, MysqlConfig, Orders, MysqlSchemas
from orders.utils.task import update_orders_progress, update_task_progress, ExportToFiles

channel_layer = get_channel_layer()
logger = logging.getLogger('django')


@shared_task
def async_full_execute(username, query, key):
    """
    批量串行执行SQL
    """
    taskid = key
    for row in OrdersTasks.objects.raw(query):
        id = row.id
        host = row.host
        port = row.port
        database = row.database
        sql = row.sql

        obj = OrdersTasks.objects.get(id=id)
        if obj.task_progress not in ('1', '2'):
            # 将任务进度设置为: 处理中
            obj.executor = username
            obj.execition_time = timezone.now()
            obj.task_progress = '2'
            obj.save()

            # 执行SQL
            # 获取工单的备注
            sql_type = Orders.objects.get(pk=obj.order_id).sql_type
            # 如果type为EXPORT
            if sql_type == 'EXPORT':
                async_export_tasks.delay(username=username,
                                         id=id,
                                         sql=sql,
                                         host=obj.host,
                                         port=obj.port,
                                         database=obj.database)
            # 如果type为DML和DDL
            if sql_type in ['DML', 'DDL']:
                async_execute_sql.delay(
                    username=username,
                    id=id,
                    sql=sql,
                    host=host,
                    port=port,
                    database=database,
                    task_progress='2')
        while OrdersTasks.objects.get(id=id).task_progress == '2':
            time.sleep(0.2)
            continue
    cache.delete(key)
    # 更新父任务进度
    update_orders_progress(username, taskid)


@shared_task
def async_execute_sql(id=None, username=None, sql=None, host=None, port=None, database=None, task_progress=None):
    """执行SQL"""
    obj = MysqlConfig.objects.get(host=host, port=port)

    execute_sql = ExecuteSql(host=host,
                             port=port,
                             user=obj.user,
                             password=obj.password,
                             charset=obj.character,
                             database=database,
                             database_type=obj.rds_type,
                             username=username)
    result = execute_sql.run_by_sql(sql)

    # 更新任务进度
    update_task_progress(id=id, exec_result=result, task_progress=task_progress)


@shared_task
def async_export_tasks(username=None, id=None, sql=None, host=None, port=None, database=None):
    export_to_excel = ExportToFiles(id, username, sql, host, port, database)
    export_to_excel.run()


@shared_task()
def periodic_sync_schemas():
    """
    定时任务，后台admin可设置为10分钟同步一次
    同步系统表auditsql_mysql_config里面配置的数据库的库元数据到本地表auditsql_mysql_schemas
    """
    ignored_schemas = ('information_schema', 'mysql', 'percona', 'performance_schema', 'sys', 'test')
    schema_filter_query = "select schema_name from information_schema.schemata " \
                          "where schema_name not in {}".format(ignored_schemas)

    for row in MysqlConfig.objects.all():
        try:
            cnx = pymysql.connect(user=row.user,
                                  password=row.password,
                                  host=row.host,
                                  port=row.port,
                                  charset=row.character,
                                  cursorclass=pymysql.cursors.DictCursor)

            with cnx.cursor() as cursor:
                cursor.execute(schema_filter_query)
                for i in cursor.fetchall():
                    MysqlSchemas.objects.update_or_create(
                        cid_id=row.id, schema=i['schema_name'],
                        defaults={
                            'cid_id': row.id,
                            'schema': i['schema_name'],
                            'envi_id': row.envi_id
                        })
        except Exception as err:
            logger.error(f'定时任务失败: periodic_sync_schemas, 主机: {row.host}')
            logger.error(err.args[1])
            continue
