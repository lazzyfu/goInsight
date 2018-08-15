# -*- coding:utf-8 -*-
# edit by fuzongfei
import datetime
import re
import hashlib

import requests

# 定义阈值参数
# 指定时间范围内执行次数大于100
from dingtalkchatbot.chatbot import DingtalkChatbot

CNT = 200

# 指定时间范围内SQL的平均响应时间，单位：ms
AVG = 100

# 忽略的库或参数
IGNORE_VALUES = ['information_schema', 'performance_schema', 'THROTTLE']

# WEBHOOK
WEBHOOK = "https://oapi.dingtalk.com/robot/send?access_token=8469850b137b308" \
          "b24151fe136dc128ab471bc782e8f2e735abab6eac175ad10"

# 初始化requests会话
s = requests.Session()
s.auth = ('pmm_user', 'pmm@yunniao_monitor')

uuid_url = "http://pmm.ops.xunhuji.me:8888/qan-api/instances?deleted=no"
r_uuid = s.get(uuid_url)

# 使用UTC时间，统计区间为1小时
local_now = datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S')
local_before_one_hours = (datetime.datetime.now() - datetime.timedelta(hours=1)).strftime('%Y-%m-%d %H:%M:%S')

utc_now = datetime.datetime.utcnow().strftime('%Y-%m-%dT%H:%M:%S')
utc_before_one_hours = (datetime.datetime.utcnow() - datetime.timedelta(hours=1)).strftime('%Y-%m-%dT%H:%M:%S')


def check_rules(abstract, rule):
    if abstract == '':
        return False
    else:
        if not re.search(rule, abstract, re.I):
            return True
        else:
            return False


result = []

for i in r_uuid.json():
    if i.get('Subsystem') == 'mysql':
        UUID = i.get('UUID')
        query_url = f"http://pmm.ops.xunhuji.me:8888/qan-api/qan/profile" \
                    f"/{UUID}?begin={utc_before_one_hours}&end={utc_now}&offset=0"
        r_query = s.get(query_url)
        query = r_query.json().get('Query')
        if isinstance(query, list):
            for line in query:
                TimeRange = f"{local_before_one_hours} -> {local_now}"
                Name = i.get('Name')
                Qps = round(line.get('QPS'), 2)
                Cnt = line.get('Stats').get('Cnt')
                Avg = round(line.get('Stats').get('Avg') * 1000, 2)
                Abstract = line.get('Abstract')
                Fingerprint = line.get('Fingerprint')
                Md5Sum = hashlib.md5(Fingerprint.encode('utf8')).hexdigest()
                Version = '-'.join(([i.get('Subsystem'), i.get('Version')]))
                if all(check_rules(line.get('Abstract'), rule) for rule in IGNORE_VALUES):
                    if Cnt > CNT and Avg > AVG:
                        result.append({
                            'TimeRange': TimeRange,
                            'Name': Name,
                            'Version': Version,
                            'Abstract': Abstract,
                            'Fingerprint': Fingerprint,
                            'QPS': Qps,
                            'Cnt': Cnt,
                            'Avg': Avg
                        })

                        # MysqlSlowLog.objects.create(
                        #     hostname=Name,
                        #     version=Version,
                        #     qps=Qps,
                        #     cnt=Cnt,
                        #     avg=Avg,
                        #     fingerprint=Fingerprint,
                        #     md5sum=Md5Sum,
                        #     timerange=TimeRange
                        # )

for i in result:
    print(i)

xiaoding = DingtalkChatbot(WEBHOOK)

title = f"#### 探测到新的慢查询\n#### 开始时间：{local_before_one_hours}\n#### 结束时间：{local_now}\n"
text = []
for p in result:
    text.append(f"- 主机名：{p.get('Name')}\n"
                f"- MySQL版本：{p.get('Version')}\n"
                f"- 区间内每秒查询数(QPS)：{p.get('QPS')}\n"
                f"- 区间内总查询次数：{p.get('Cnt')}\n"
                f"- 区间内平均响应时间：{p.get('Avg')}ms\n"
                f"- SQL指纹：\n"
                f"> {p.get('Fingerprint')}\n")

xiaoding.send_markdown(title='慢查询', text=title + ''.join(text), is_at_all=True)
# text = f"慢查询统计时间区间：{local_before_one_hours} -> {local_now}\n" \
#        f"主机名：{p.get('Name')}\n" \
#        f"MySQL版本：{p.get('Version')}\n" \
#        f"区间内每秒查询数(QPS)：{p.get('QPS')}\n" \
#        f"区间内总查询次数：{p.get('Cnt')}\n" \
#        f"区间内平均响应时间：{p.get('Avg')}ms\n" \
#        f"SQL指纹：{p.get('Fingerprint')}\n"
# xiaoding.send_text(msg=text, is_at_all=True)
