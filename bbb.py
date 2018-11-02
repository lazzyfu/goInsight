# -*- coding:utf-8 -*-
# edit by fuzongfei
# import pymysql
#
# connection = pymysql.connect(host='10.10.1.202',
#                              port=3306,
#                              user='yops',
#                              password='Fuzongfei_1991',
#                              charset='utf8mb4',
#                              database='aa',
#                              )
#
# try:
#     with connection.cursor() as cursor:
#         # Create a new record
#         sql = "INSERT INTO `a1a` (`I_REF`, `name`) VALUES (2 , 'gggg...')"
#         cursor.execute(sql)
#         print(cursor.description)
#
#     # connection is not autocommit by default. So you must commit to save
#     # your changes.
#     connection.commit()
#
#     # with connection.cursor() as cursor:
#     #     # Read a single record
#     #     sql = "select * from a1a"
#     #     cursor.execute(sql)
#     #     result = cursor.fetchall()
#     #     for i in result:
#     #         print(i)
# except Exception as err:
#     print(err)
# finally:
#     connection.close()


table_name = extract_tables("update aa,bb set aa.id=bb.id where aa.xx = bb.xx;")
for i in table_name:
    print(i.name)