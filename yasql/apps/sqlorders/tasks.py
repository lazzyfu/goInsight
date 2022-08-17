# -*- coding:utf-8 -*-
# edit by xff
import json
import os
import time
from datetime import datetime

import pymysql
from clickhouse_driver import Client
from celery import shared_task
from celery.utils.log import get_task_logger
from django.db.models import F
from django.utils import timezone

from config import REOMOTE_USER
from sqlorders import models
from sqlorders.api.executeExportApi import ExecuteExport
from sqlorders.api.executeSqlApi import ExecuteSQL
from sqlorders.libs import remove_sql_comment
from sqlorders.notice import MsgNotice

logger = get_task_logger('celery.logger')


def update_dborders_progress_to_processing(id, username):
    # 更新父工单的进度为：处理中
    obj = models.DbOrders.objects.get(pk=id)
    executor = json.loads(obj.executor)
    for i in executor:
        i['user'] = username
        i['time'] = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
    obj.executor = json.dumps(executor)
    obj.progress = 3
    obj.updated_at = timezone.now()
    obj.save()


def update_dborders_progress_to_finish(task_id, username=None):
    # 更新父工单的进度为：已完成
    obj = models.DbOrdersExecuteTasks.objects.filter(task_id=task_id)
    progress = obj.values_list('progress', flat=True)
    order_id = obj.first().order_id

    if order_id:
        if all([False for i in list(progress) if i != 1]):
            data = models.DbOrders.objects.get(pk=order_id)
            if data.progress != 4:
                data.progress = 4
                data.updated_at = timezone.now()
                data.save()
            # 推送消息
            msg_notice.delay(
                pk=order_id,
                op='_feedback',
                username=username
            )


def save_rbsql_as_file(rollbacksql):
    """当备份的数据太大时,数据库由于max_allowed_packet问题无法保存,此时保存到文件"""
    if not os.path.exists(r'media/rbsql'):
        os.makedirs('media/rbsql')
    filename = f"media/rbsql/rbsql_{datetime.now().strftime('%Y%m%d%H%M%S%f')}.sql"
    with open(filename, 'w') as f:
        f.write(rollbacksql)
    return filename


def save_ghost_log_as_file(log):
    """当日志太大时,数据库由于max_allowed_packet问题无法保存,此时保存到文件"""
    if not os.path.exists(r'media/ghost_log'):
        os.makedirs('media/ghost_log')
    filename = f"media/ghost_log/ghost_{datetime.now().strftime('%Y%m%d%H%M%S%f')}.sql"
    with open(filename, 'w') as f:
        f.write(log)
    return filename


def update_task_progress_to_finish(id=None, result=None):
    """
    更新当前任务的进度
    result:
    - {'status': 'fail', 'execute_log': ''}
    - {'status': 'success', 'execute_log': '', 'rollbacksql': '', 'consuming_time': 0.000, 'affected_rows': 0}
    """
    obj = models.DbOrdersExecuteTasks.objects.get(pk=id)
    if result['status'] == 'fail':
        obj.progress = 3
        obj.execute_log = result.get('execute_log')
        obj.save()
    if result['status'] == 'success':
        rbsql = result.get('rollbacksql')
        execute_log = result.get('execute_log')
        if obj.progress == 2:
            try:
                obj.rollback_sql = rbsql
                obj.save()
            except Exception as err:
                filename = save_rbsql_as_file(rbsql)
                obj.rollback_sql = '\n'.join([
                    '回滚数据超出max_allowed_packet,写入到数据库失败',
                    '备份数据已经以文本的形式进行了保存',
                    '存储路径：',
                    filename
                ])
            try:
                obj.execute_log = execute_log
            except Exception as err:
                filename = save_ghost_log_as_file(execute_log)
                obj.rollback_sql = '\n'.join([
                    '日志数据超出max_allowed_packet,写入到数据库失败',
                    '日志数据已经以文本的形式进行了保存',
                    '存储路径：',
                    filename
                ])
            finally:
                obj.consuming_time = result.get('consuming_time')
                obj.progress = 1
                obj.affected_rows = result.get('affected_rows')
                obj.save()


@shared_task(queue='dbtask')
def async_execute_sql(id=None, sql=None):
    """异步执行SQL,获取当前任务的db信息"""
    config = models.DbOrdersExecuteTasks.objects.filter(pk=id).annotate(
        host=F('order__cid__host'),
        port=F('order__cid__port'),
        charset=F('order__cid__character'),
        rds_type=F('order__cid__rds_type'),
        rds_category=F('order__rds_category')
    ).values('host', 'port', 'charset', 'rds_type', 'sql_type', 'rds_category').first()

    # 获取当前任务对应的数据库
    database = models.DbOrdersExecuteTasks.objects.filter(pk=id).annotate(
        database=F('order__database')
    ).values('database', 'task_id').first()
    config.update(database)
    config.update(REOMOTE_USER)

    # 执行SQL
    execute_sql = ExecuteSQL(config=config)
    result = execute_sql.run_by_sql(remove_sql_comment(sql))
    # 更新当前任务的进度
    update_task_progress_to_finish(id=id, result=result)


@shared_task(queue='dbtask')
def async_execute_export(id=None, username=None, sql=None):
    """异步执行导出, 获取当前任务的db信息"""
    config = models.DbOrdersExecuteTasks.objects.filter(pk=id).annotate(
        host=F('order__cid__host'),
        port=F('order__cid__port'),
        charset=F('order__cid__character'),
        rds_type=F('order__cid__rds_type'),
        rds_category=F('order__rds_category')
    ).values('host', 'port', 'charset', 'rds_type', 'rds_category').first()

    # 获取当前任务对应的数据库
    database = models.DbOrdersExecuteTasks.objects.filter(pk=id).annotate(
        database=F('order__database')
    ).values('database', 'task_id', 'file_format').first()
    config['id'] = id
    config['username'] = username
    config.update(database)
    config.update(REOMOTE_USER)

    # 执行SQL
    execute_export = ExecuteExport(config=config)
    result = execute_export.run(remove_sql_comment(sql))
    # 更新当前任务的进度
    update_task_progress_to_finish(id=id, result=result)


@shared_task(queue='dbtask')
def async_execute_single(id=None, username=None):
    """执行单条SQL"""
    obj = models.DbOrdersExecuteTasks.objects.get(pk=id)

    if obj.sql_type in ['DML', 'DDL']:
        async_execute_sql.delay(
            id=id,
            sql=obj.sql
        )
    if obj.sql_type in ['EXPORT']:
        async_execute_export.delay(
            id=id,
            username=obj.applicant,  # 收件人为工单的申请人
            sql=obj.sql
        )
    # 监控当前任务是否执行完成,当执行完成后,执行修改父工单的状态
    # 否则方法update_dborders_progress_to_finish将先于任务执行完成
    # 此处的while不能引用obj,obj属于查询结果的对象,而不会每次查询数据,否则死的很惨...
    while models.DbOrdersExecuteTasks.objects.get(pk=id).progress == 2:
        time.sleep(0.01)
        continue
    update_dborders_progress_to_finish(task_id=obj.task_id, username=username)


@shared_task(queue='dbtask')
def async_execute_multi(task_id=None, username=None):
    # 匹配未执行和失败的工单
    for row in models.DbOrdersExecuteTasks.objects.filter(task_id=task_id, progress__in=(0, 3)):
        # 将当前任务进度设置为：处理中
        obj = models.DbOrdersExecuteTasks.objects.get(pk=row.id)
        obj.executor = username
        obj.execute_time = timezone.now()
        obj.progress = 2
        obj.save()
        # 执行
        async_execute_single.delay(id=row.id, username=username)
        # 监控当前任务是否执行完成,当执行完成后,继续执行下一个任务
        while models.DbOrdersExecuteTasks.objects.get(id=row.id).progress == 2:
            time.sleep(0.01)
            continue
    update_dborders_progress_to_finish(task_id=task_id, username=username)


def dbms_sync_clickhouse_schema(row):
    ignored_schemas = ('_temporary_and_external_tables', 'system', 'default')
    query = f"select name from system.databases where name not in {ignored_schemas}"
    config = {
        'host': row.host,
        'port': row.port,
        'database': 'default',
        'connect_timeout': 5,
        'send_receive_timeout': 5,
    }
    # 请在clickhouse创建好用户
    config.update(REOMOTE_USER)
    cnx = Client(**config)
    result = cnx.execute(query)
    for i in result:
        schema = i[0]
        models.DbSchemas.objects.update_or_create(
            cid_id=row.id,
            schema=schema,
            defaults={'schema': schema}
        )
    cnx.disconnect()


def dbms_sync_mysql_schema(row):
    ignored_schemas = ('PERFORMANCE_SCHEMA', 'INFORMATION_SCHEMA', 'PERCONA', 'MYSQL', 'SYS',
                       'DM_META', 'DM_HEARTBEAT', 'DBMS_MONITOR', 'METRICS_SCHEMA', 'TIDB_BINLOG', 'TIDB_LOADER')
    query = f"SELECT SCHEMA_NAME FROM INFORMATION_SCHEMA.SCHEMATA"
    config = {
        'host': row.host,
        'port': row.port,
        'read_timeout': 10,  # socket.timeout: timed out,比如阿里的rds就很操蛋,没开白名单会hang住
        'cursorclass': pymysql.cursors.DictCursor
    }
    config.update(REOMOTE_USER)
    cnx = pymysql.connect(**config)
    with cnx.cursor() as cursor:
        cursor.execute(query)
        for i in cursor.fetchall():
            if i['SCHEMA_NAME'].upper() not in ignored_schemas:
                models.DbSchemas.objects.update_or_create(
                    cid_id=row.id,
                    schema=i['SCHEMA_NAME'],
                    defaults={'schema': i['SCHEMA_NAME']}
                )
    cnx.close()


@shared_task(queue='dbtask')
def dbms_sync_dbschems():
    """
    同步远程的schema信息到表yasql_dbms_sql_schemas
    用于提交工单使用
    """
    # 同步DB信息
    for row in models.DbConfig.objects.filter(use_type=0):
        try:
            category = row.rds_category
            if category == 3:  # clickhouse
                dbms_sync_clickhouse_schema(row)
            elif category in (1, 2):
                dbms_sync_mysql_schema(row)
        except Exception as err:
            logger.error(f"异常主机：{row.host}")
            logger.error(err)
            continue
    logger.info(f'dbms_sync_dbschems任务结束...')


@shared_task()
def msg_notice(**kwargs):
    MsgNotice(**kwargs).run()
