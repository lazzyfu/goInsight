import re

from django.test import TestCase

# Create your tests here.
import pymysql

conn = pymysql.connect(host='10.72.63.128', user='inception_user', password='inception@123com', db='test',
                       charset='utf8', cursorclass=pymysql.cursors.DictCursor)

# 此方法可能受中间件影响，查询结果不准确
# mysql_version = int(conn.server_version.split('.', 2)[1])

VERSION_QUERY = 'select version();'

try:
    with conn.cursor() as cursor:
        cursor.execute(VERSION_QUERY)
        mysql_version = int(cursor.fetchone().get('version()').split('.', 2)[1])

        id = 1

        user_query = "select user from mysql.user"
        cursor.execute(user_query)
        user_list = []
        for row in cursor.fetchall():
            user_list.append(row.get('user'))

        user_set = list(set(user_list))

        user_dict = []
        for i in user_set:
            user_dict.append({'id': id, 'pid': 0, 'privileges': '', 'schema': '', 'user': i})
            id += 1

        if mysql_version > 6:
            user_info_query = "select concat(\"'\",user,\"'\",'@',\"'\",host,\"'\") as username, " \
                              "concat(left(authentication_string,5),'...',right(authentication_string,2)) " \
                              "as password from mysql.user"
        else:
            user_info_query = "select concat(\"'\",user,\"'\",'@',\"'\",host,\"'\") as username," \
                              "concat(left(password,5),'...',right(password,2)) as password from mysql.user"
        cursor.execute(user_info_query)

        privileges_dict = []

        for row in cursor.fetchall():
            user = row.get('username')
            password = row.get('password')
            print(user)
            username = user.split('@')[0].replace("'", '')

            privileges_query = f"show grants for {user}"
            cursor.execute(privileges_query)
            pid = 0
            for i in cursor.fetchall():
                for k, v in i.items():
                    data = re.split(r' TO ', v.replace('GRANT', '').strip())[0].split(r' ON ')
                    privileges = data[0]
                    host = data[1]
                    for j in user_dict:
                        if username == j.get('user'):
                            pid = j.get('id')
                    privileges_dict.append({
                        'id': id,
                        'pid': pid,
                        'user': user,
                        'password': password,
                        'privileges': privileges,
                        'schema': host
                    })
                    id += 1
        data = user_dict + privileges_dict
        for i in data:
            print(i)

finally:
    conn.close()
