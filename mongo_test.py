# -*- coding:utf-8 -*-
# edit by fuzongfei

# from pymongo import MongoClient
# import urllib.parse
#
# username = urllib.parse.quote_plus('dba')
# password = urllib.parse.quote_plus('123.com')
#
# client = MongoClient('mongodb://%s:%s@172.17.101.41:27017' % (username, password))
#
# # 获取所有的库
# # 返回库名 和大小
# schema = client.list_databases()
# for i in schema:
#     schema = i['name']
#
#     db = client[schema]
#     for j in db.list_collections():
#         print(j)
#
#
#
# db = client['test']
#
# userinfo = db['userinfo']
#
# print()
# for i in userinfo.find().limit(20):
#     print(i)
#
# result = userinfo.find({"name": "welcome you", 'age': {'$gt': 20}})
#


import subprocess

cmd = "mongo --host=172.17.101.41 -udba -p'123.com' --authenticationDatabase admin --eval 'db.userinfo.find().skip(20)'"

status, output = subprocess.getstatusoutput(cmd)

print(status)
print(output)