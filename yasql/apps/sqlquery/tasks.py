# -*- coding:utf-8 -*-
# edit by xff
import clickhouse_driver
import pymysql
from celery import shared_task
from celery.utils.log import get_task_logger

from config import QUERY_USER
from sqlorders.models import DbConfig
from sqlquery.models import DbQuerySchemas, DbQueryTables

logger = get_task_logger('celery.logger')


@shared_task(queue='dbtask')
def sqlquery_sync_schemas_tables():
    """
        同步远程的schema信息到表yasql_sqlquery_schemas/yasql_sqlquery_tables，查询授权使用
        """
    ignored_schemas = (
        'PERFORMANCE_SCHEMA', 'INFORMATION_SCHEMA', 'PERCONA', 'MYSQL', 'SYS',
        'DM_META', 'DM_HEARTBEAT', 'DBMS_MONITOR', 'METRICS_SCHEMA',
        'TIDB_BINLOG', 'TIDB_LOADER',
        '_TEMPORARY_AND_EXTERNAL_TABLES', 'DEFAULT', 'MYSQL_MONITOR', 'SYSTEM'
    )

    mysql_query = f"SELECT " \
                  f"TABLE_SCHEMA," \
                  f"GROUP_CONCAT(TABLE_NAME) AS TABLE_NAME " \
                  f"FROM INFORMATION_SCHEMA.TABLES " \
                  f"WHERE  TABLE_SCHEMA NOT IN {ignored_schemas} " \
                  f"GROUP BY TABLE_SCHEMA"

    clickhouse_query = f"SELECT " \
                       f"`database` AS TABLE_SCHEMA, " \
                       f"`tables`  AS TABLE_NAME " \
                       f"FROM " \
                       f"(" \
                       f"SELECT " \
                       f"`database`, " \
                       f"groupArray(`name`) AS `name_array`, " \
                       f"arrayStringConcat(name_array, ',') AS `tables` " \
                       f"FROM system.tables " \
                       f"WHERE `database` NOT IN {ignored_schemas} GROUP BY `database`" \
                       f")"

    for row in DbConfig.objects.filter(use_type=1):
        try:
            result = []
            if row.rds_category in [1, 2]:
                # mysql/tidb
                config = {
                    'host': row.host,
                    'port': row.port,
                    'read_timeout': 3,  # socket.timeout: timed out，比如阿里的rds就很操蛋，没开白名单会hang住
                }
                config.update(QUERY_USER)
                cnx = pymysql.connect(**config)
                cursor = cnx.cursor()
                cursor.execute('set session group_concat_max_len = 1073741824')
                cursor.execute(mysql_query)
                result = cursor.fetchall()
                cursor.close()
                cnx.close()
            if row.rds_category in [3]:
                # clickhouse
                config = {
                    'host': row.host,
                    'port': row.port,
                    'connect_timeout': 3,
                    'database': 'system',
                }
                config.update(QUERY_USER)
                cnx = clickhouse_driver.connect(**config)
                cursor = cnx.cursor()
                cursor.execute(clickhouse_query)
                result = cursor.fetchall()
                cursor.close()
                cnx.close()

            for i in result:
                obj, _ = DbQuerySchemas.objects.update_or_create(
                    cid_id=row.id,
                    schema=i[0],
                    defaults={'schema': i[0]}
                )
                for t in i[1].split(','):
                    DbQueryTables.objects.update_or_create(
                        schema_id=obj.pk,
                        table=t,
                        defaults={'table': t}
                    )
        except Exception as err:
            logger.error(f"异常主机：{row.host}")
            logger.error(err)
            continue
    logger.info(f'dbms_sync_dbschems任务结束...')
