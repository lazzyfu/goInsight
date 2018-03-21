import pymysql

# try:
# 连接到inception
try:
    conn = pymysql.connect(host='10.72.63.128', user='inception_user', password='inception@123com',
                           port=3306, use_unicode=True, charset="utf8", connect_timeout=1)

    if conn:
        print('ok')
    conn.close()
except pymysql.Error as err:
    print(err)
