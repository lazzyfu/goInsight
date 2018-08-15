# -*- coding:utf-8 -*-
# edit by fuzongfei

# -*- coding:utf-8 -*-
# edit by fuzongfei

import pymysql.cursors

user = 'dba'
password = 'yn_dbapass'

IGNORED_PARAMS = ('information_schema', 'mysql', 'percona', 'test', 'performance_schema', 'sys')
SCHEMA_INFO = []
SCHEMA_QUERY = "select schema_name from information_schema.schemata " \
               "where SCHEMA_NAME not in {0}".format(IGNORED_PARAMS)

CONN_INFO = [
    {'user': 'incep_user', 'password': 'incep123com', 'db_host': '192.168.203.16', 'db_port': 3306, 'envi': 1,
     'comment': 'dev-mysql-192.168.203.16'},
    {'user': 'incep_user', 'password': 'incep123com', 'db_host': '192.168.200.236 ', 'db_port': 3306, 'envi': 1,
     'comment': 'f1-mysql-192.168.200.236'},
    {'user': 'incep_user', 'password': 'incep123com', 'db_host': '192.168.203.36', 'db_port': 3306, 'envi': 1,
     'comment': 'f2-mysql-192.168.203.36'},
    {'user': 'incep_user', 'password': 'incep123com', 'db_host': '192.168.200.159', 'db_port': 3306, 'envi': 1,
     'comment': 'f3-mysql-192.168.200.159'},
    {'user': 'incep_user', 'password': 'incep123com', 'db_host': '192.168.203.56', 'db_port': 3306, 'envi': 1,
     'comment': 'm1-mysql-192.168.203.56'}
]

for row in CONN_INFO:
    cnx = pymysql.connect(user=row['user'],
                          password=row['password'],
                          host=row['db_host'],
                          port=row['db_port'],
                          charset='utf8mb4',
                          cursorclass=pymysql.cursors.DictCursor)
    try:
        with cnx.cursor() as cursor:
            cursor.execute(SCHEMA_QUERY)
            schema = []
            for i in cursor.fetchall():
                SCHEMA_INFO.append({
                    'user': 'sqlaudit_query',
                    'password': 'sqlaudit_query123com',
                    'db_host': row['db_host'],
                    'db_port': row['db_port'],
                    'envi': row['envi'],
                    'comment': row['comment'],
                    'schema': i['schema_name']})
    finally:
        cnx.close()
for i in SCHEMA_INFO:
    cmd = "insert into auditsql_mysql_schema(user, password, `host`, `port`, `schema`, envi, comment, created_at, updated_at)" \
          " values('{0}', '{1}', '{2}', {3}, '{4}', {5}, '{6}', now(),now());".format(
        i['user'], i['password'], i['db_host'], i['db_port'], i['schema'], i['envi'], i['comment'])
    print(cmd)
