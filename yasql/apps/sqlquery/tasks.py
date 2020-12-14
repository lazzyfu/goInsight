# -*- coding:utf-8 -*-
# edit by fuzongfei
import pymysql
from celery import shared_task
from celery.utils.log import get_task_logger

from config import QUERY_USER
from sqlorders.models import DbConfig
from sqlquery.models import DbQuerySchemas, DbQueryTables

logger = get_task_logger('noahCelery')


@shared_task(queue='dbtask')
def sqlquery_sync_schemas_tables():
    """
    同步远程的schema信息到表noah_sqlquery_meta，查询授权使用
    """
    ignored_schemas = ('PERFORMANCE_SCHEMA', 'INFORMATION_SCHEMA', 'PERCONA', 'MYSQL', 'SYS',
                       'DM_META', 'DM_HEARTBEAT', 'DBMS_MONITOR', 'METRICS_SCHEMA', 'TIDB_BINLOG', 'TIDB_LOADER')
    query = f"select TABLE_SCHEMA,GROUP_CONCAT(TABLE_NAME) as TABLE_NAME " \
            f"FROM INFORMATION_SCHEMA.TABLES " \
            f"WHERE  TABLE_SCHEMA NOT IN {ignored_schemas} GROUP BY TABLE_SCHEMA"
    # 同步DB的schema信息
    for row in DbConfig.objects.filter(use_type=1):
        try:
            config = {
                'host': row.host,
                'port': row.port,
                'read_timeout': 10,  # socket.timeout: timed out，比如阿里的rds就很操蛋，没开白名单会hang住
                'cursorclass': pymysql.cursors.DictCursor
            }
            config.update(QUERY_USER)
            cnx = pymysql.connect(**config)
            with cnx.cursor() as cursor:
                cursor.execute('set session group_concat_max_len = 1073741824')
                cursor.execute(query)
                for i in cursor.fetchall():
                    obj, _ = DbQuerySchemas.objects.update_or_create(
                        cid_id=row.id,
                        schema=i['TABLE_SCHEMA'],
                        defaults={'schema': i['TABLE_SCHEMA']}
                    )
                    for t in i['TABLE_NAME'].split(','):
                        DbQueryTables.objects.update_or_create(
                            schema_id=obj.pk,
                            table=t,
                            defaults={'table': t}
                        )
            cnx.close()
        except Exception as err:
            logger.error(f"异常主机：{row.host}")
            logger.error(err)
            continue
    logger.info(f'dbms_sync_dbschems任务结束...')
