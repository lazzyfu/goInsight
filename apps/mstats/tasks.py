# -*- coding:utf-8 -*-
# edit by fuzongfei
import datetime
import hashlib
import logging
import re
import subprocess

import pymysql
import requests
from celery import shared_task
from dingtalkchatbot.chatbot import DingtalkChatbot

from AuditSQL.settings import REMOTE_MYSQL_USER, REMOTE_MYSQL_PASSWORD
from mstats.models import MysqlSchemaInfo, MysqlSlowLog, DeadlockCommand, DeadlockRecord, MySQLConfigSource
from project_manager.models import Webhook

logger = logging.getLogger(__name__)


@shared_task
def sync_schemas():
    user = REMOTE_MYSQL_USER
    password = REMOTE_MYSQL_PASSWORD

    ignored_params = ('information_schema', 'mysql', 'percona', 'performance_schema', 'sys', 'test')
    schema_filter_query = "select schema_name from information_schema.schemata " \
                          "where SCHEMA_NAME not in {0}".format(ignored_params)

    collect_from_host = []
    for row in MySQLConfigSource.objects.all():
        collect_from_host.append({
            'user': REMOTE_MYSQL_USER,
            'password': REMOTE_MYSQL_PASSWORD,
            'db_host': row.host,
            'db_port': row.port,
            'envi': row.envi,
            'is_master': row.is_master,
            'comment': row.comment
        })

    # 连接到目标数据库，统计schema
    for row in collect_from_host:
        try:
            cnx = pymysql.connect(user=row['user'],
                                  password=row['password'],
                                  host=row['db_host'],
                                  port=row['db_port'],
                                  charset='utf8mb4',
                                  cursorclass=pymysql.cursors.DictCursor)
            try:
                with cnx.cursor() as cursor:
                    cursor.execute(schema_filter_query)
                    for i in cursor.fetchall():
                        MysqlSchemaInfo.objects.update_or_create(
                            user=user,
                            password=password,
                            host=row['db_host'],
                            port=row['db_port'],
                            schema=i['schema_name'],
                            envi=row['envi'],
                            is_master=row['is_master'],
                            schema_join='_'.join(
                                ([row['db_host'], str(row['db_port']), i['schema_name']])),
                            comment=row['comment'],
                        )
            finally:
                cnx.close()
        except Exception as err:
            logger.error(err)
            continue


def check_rules(abstract, rule):
    if abstract == '':
        return False
    else:
        if not re.search(rule, abstract, re.I):
            return True
        else:
            return False


@shared_task
def check_mysql_slowlog():
    # 定义阈值参数
    # 指定时间范围内执行次数大于100
    cnt_value = 200

    # 指定时间范围内SQL的平均响应时间，单位：ms
    avg_value = 100

    # 忽略的库或参数
    ignore_values = ['information_schema', 'performance_schema', 'THROTTLE']

    # webhook
    webhook = "https://oapi.dingtalk.com/robot/send?access_token=8469850b137b308" \
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

    result = []

    for i in r_uuid.json():
        if i.get('Subsystem') == 'mysql':
            uuid = i.get('UUID')
            query_url = f"http://pmm.ops.xunhuji.me:8888/qan-api/qan/profile" \
                        f"/{uuid}?begin={utc_before_one_hours}&end={utc_now}&offset=0"
            r_query = s.get(query_url)
            query = r_query.json().get('Query')
            if isinstance(query, list):
                for line in query:
                    timerange = f"{local_before_one_hours} -> {local_now}"
                    name = i.get('Name')
                    qps = round(line.get('QPS'), 2)
                    cnt = line.get('Stats').get('Cnt')
                    avg = round(line.get('Stats').get('Avg') * 1000, 2)
                    abstract = line.get('Abstract')
                    fingerprint = line.get('Fingerprint')
                    md5sum = hashlib.md5(fingerprint.encode('utf8')).hexdigest()
                    version = '-'.join(([i.get('Subsystem'), i.get('Version')]))
                    if all(check_rules(abstract, rule) for rule in ignore_values):
                        if cnt > cnt_value and avg > avg_value:
                            result.append({
                                'timerange': timerange,
                                'name': name,
                                'version': version,
                                'abstract': abstract,
                                'fingerprint': fingerprint,
                                'qps': qps,
                                'cnt': cnt,
                                'avg': avg,
                                'md5sum': md5sum
                            })

                            MysqlSlowLog.objects.get_or_create(
                                md5sum=md5sum,
                                defaults={'hostname': name,
                                          'version': version,
                                          'qps': qps,
                                          'cnt': cnt,
                                          'avg': avg,
                                          'fingerprint': fingerprint,
                                          'timerange': timerange
                                          }
                            )

    xiaoding = DingtalkChatbot(webhook)
    text = []
    if result:
        for p in result:
            obj = MysqlSlowLog.objects.filter(md5sum=p.get('md5sum')).first()
            is_pull = obj.is_pull if obj else 0
            if is_pull == 0:
                text.append(f"- 主机名：{p.get('name')}\n"
                            f"- MySQL版本：{p.get('version')}\n"
                            f"- 区间内每秒查询数(QPS)：{p.get('qps')}\n"
                            f"- 区间内总查询次数：{p.get('cnt')}\n"
                            f"- 区间内平均响应时间：{p.get('avg')}ms\n"
                            f"- MD5：{p.get('md5sum')}\n"
                            f"- SQL摘要：{p.get('abstract')}\n"
                            f"- SQL指纹：\n"
                            f"> {p.get('fingerprint')}\n\n\n")
        title = f"#### 探测到新的慢查询，总数：{len(text)}条\n#### 开始时间：{local_before_one_hours}\n#### 结束时间：{local_now}\n"
        xiaoding.send_markdown(title='SQL慢查询', text=title + ' '.join(text))


@shared_task
def detect_deadlock():
    # webhook
    if Webhook.objects.filter().first():
        webhook_addr = Webhook.objects.get().webhook_addr
        xiaoding = DingtalkChatbot(webhook_addr)
        # 检查实例，并生生成实例死锁记录的命令
        # 使用本机的数据库作为死锁记录
        # 库名：auditsql，表名：dbaudit_deadlocks_records
        command = "/usr/bin/pt-deadlock-logger --user={user} --password={password} --host={host} --port={port} " \
                  "--no-version-check --create-dest-table " \
                  "--dest h=localhost,u=root,p=123.com,D=auditsql,t=dbaudit_deadlocks_records --iterations 1"

        query = "SELECT id, `user`, `password`, `host`, `port` FROM auditsql_mysql_schema " \
                "WHERE auditsql_mysql_schema.is_master = 1 group by host,port"

        for row in MysqlSchemaInfo.objects.raw(query):
            format_command = command.format(user=row.user, password=row.password, host=row.host, port=row.port)
            if not DeadlockCommand.objects.filter(schema_id=row.id):
                print('www....')
                DeadlockCommand.objects.create(schema_id=row.id, command=format_command)

        # 轮询探测死锁
        for row in DeadlockCommand.objects.all():
            process = subprocess.Popen(row.command, shell=True)
            process.wait()
        # 检查进程是否启动，若没有，则启动进程
        # for row in DeadlockCommand.objects.filter(pid__gte=0):
        #     # 检查进程是否运行
        #     # 如果不存在，启动该进程
        #     if not psutil.pid_exists(row.pid):
        #         process = psutil.Popen(['/usr/bin/pt-deadlock-logger', row.command])
        #         DeadlockCommand.objects.filter(id=row.id).update(pid=process.pid, is_process_run=1)

        # 检查死锁，并发送报告
        i = 0
        step = 2
        result = []
        data = list(DeadlockRecord.objects.filter(is_pull=0).values())
        while i <= (len(data) - step):
            result.append({'data': [data[i], data[i + 1]]})
            i += step

        format_deadlock_data = ''
        j = 1
        for row in result:
            double_data = ''
            for i in row['data']:
                text = f"主机：{i['server']}\n" \
                       f"时间: {i['ts']}\n" \
                       f"线程ID: {i['thread']}\n" \
                       f"事务ID: {i['txn_id']}\n" \
                       f"事务激活时间: {i['txn_time']}\n" \
                       f"用户名: {i['user']}\n" \
                       f"主机名: {i['hostname']}\n" \
                       f"IP: {i['ip']}\n" \
                       f"库名: {i['db']}\n" \
                       f"表名: {i['tbl']} \n" \
                       f"发生死锁的索引: {i['idx']}\n" \
                       f"锁类型: {i['lock_type']}\n" \
                       f"锁模式: {i['lock_mode']}\n" \
                       f"请求锁: {i['wait_hold']}\n" \
                       f"是否回滚: {'否' if i['victim'] == 0 else '是'}\n" \
                       f"查询: {i['query']}\n\n"
                double_data += text
                DeadlockRecord.objects.filter(id=i['id']).update(is_pull=1)

            format_deadlock_data += ''.join((f'## 死锁记录{j} ##:\n', double_data))
            j += 1

        if result:
            check_time = datetime.datetime.now().strftime("%Y-%m-%d %H:%M:%S")
            xiaoding.send_text(msg='\n'.join((f'【警告 ◕﹏◕，探测到新的死锁记录，探测时间：{check_time}】\n', format_deadlock_data)),
                               is_at_all=True)
