# -*- coding:utf-8 -*-
# edit by fuzongfei

import pymysql.cursors

user = 'incep_user'
password = 'incep123com'

IGNORED_PARAMS = ('information_schema', 'mysql', 'percona', 'performance_schema', 'sys')
SCHEMA_INFO = []
SCHEMA_QUERY = "select schema_name from information_schema.schemata " \
               "where SCHEMA_NAME not in {0}".format(IGNORED_PARAMS)

OFFLINE_CONN_INFO = [
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

ONLINE_CONN_INFO = []

# 连接到192.168.0.67获取线上数据
conn = pymysql.connect(host='192.168.0.67',
                       user=user,
                       password=password,
                       db='yndbops',
                       charset='utf8mb4',
                       cursorclass=pymysql.cursors.DictCursor)

try:
    with conn.cursor() as cursor:
        sql = 'select b.db_host,b.db_port from mysql_cluster a join mysql_node b ' \
              'on a.id = b.cluster_id where b.is_master=1 group by a.id,b.db_host,b.db_port'
        cursor.execute(sql)

        for i in cursor:
            db_host = i['db_host']
            db_port = i['db_port']

            ONLINE_CONN_INFO.append({
                'user': user,
                'password': password,
                'db_host': db_host,
                'db_port': db_port,
                'envi': 0,
                'comment': None
            })
finally:
    conn.close()

ALL_CONN_INFO = OFFLINE_CONN_INFO + ONLINE_CONN_INFO

# 连接到目标数据库，统计schema
for row in ALL_CONN_INFO:
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
                    'user': user,
                    'password': password,
                    'db_host': row['db_host'],
                    'db_port': row['db_port'],
                    'schema': i['schema_name'],
                    'envi': row['envi'],
                    'comment': row['comment'],
                })
    finally:
        cnx.close()

for i in SCHEMA_INFO:
    cmd = "insert into auditsql_mysql_schema(user,password,`host`,`port`,`schema`,created_at,updated_at) values('{0}', '{1}', '{2}', {3}, '{4}', now(),now());".format(i['user'], i['password'], i['db_host'], i['db_port'], i['schema'])
    print(cmd)
