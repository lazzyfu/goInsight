# 连接到inception
import json

import pymysql
import time

from asgiref.sync import async_to_sync

from AuditSQL.consumers import channel_layer

# while True:
#     conn = pymysql.connect(host='10.72.63.197', user='root', password='', db='',
#                            port=6033, use_unicode=True, charset="utf8")
#     cur = conn.cursor()
#     cur.execute("inception get osc_percent '*C54B086283FCD804F1A71DF1A53B2D9BB8C1297A'")
#     result = cur.fetchall()
#     if result:
#         field_names = [i[0] for i in cur.description]
#         result_all = []
#         for row in result:
#             result_all.append(dict(map(lambda x, y: [x, y], field_names, row)))
#         cur.close()
#         conn.close()
#         print(result_all)
#         time.sleep(2)
#         async_to_sync(channel_layer.group_send)('fuzongfei', {"type": "user.message",
#                                                               'text': json.dumps(result_all)})

aa = "[{'ID': 1, 'stage': 'RERUN', 'errlevel': 0, 'stagestatus': 'Execute Successfully', 'errormessage': 'None', 'SQL': 'use test', 'Affected_rows': 0, 'sequence': \"'1520844509_15056_0'\", 'backup_dbname': 'None', 'execute_time': '0.000', 'sqlsha1': ''}, {'ID': 2, 'stage': 'EXECUTED', 'errlevel': 0, 'stagestatus': 'Execute Successfully\\nBackup successfully', 'errormessage': 'None', 'SQL': \"ALTER TABLE xboss_account MODIFY `display_money` tinyint(3) NOT NULL DEFAULT '1' COMMENT '\u662f\u5426\u663e\u793a\u91d1\u989d\uff1a1\u663e\u793a\uff0c2\u9690\u85cf'\", 'Affected_rows': 1, 'sequence': \"'1520844539_15056_1'\", 'backup_dbname': '10_72_63_128_3306_test', 'execute_time': '30.360', 'sqlsha1': '*F8470F1EA082A00E0821BE60DC949B736A0CCDE6'}]"
print(aa)

for row in eval(aa):
    print(row)