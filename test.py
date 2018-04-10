# -*- coding:utf-8 -*-
# edit by fuzongfei

import pymysql

conn = pymysql.connect(host='10.72.63.128', user='inception_user', password='inception@123com', db='test',
                       charset='utf8', cursorclass=pymysql.cursors.DictCursor)

with conn.cursor() as cursor:
    sql = "select * from information_schema.processlist"

    cursor.execute(sql)

    result = []
    for row in cursor.fetchall():
        print(row)
        pass