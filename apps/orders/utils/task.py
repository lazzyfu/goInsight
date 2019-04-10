# -*- coding:utf-8 -*-
# edit by fuzongfei
import csv
import json
import logging
import os
import smtplib
import subprocess
import time
from datetime import datetime

import pymysql
from asgiref.sync import async_to_sync
from channels.layers import get_channel_layer
from django.core.files import File
from django.core.mail import EmailMessage
from django.template.loader import render_to_string
from django.utils import timezone
from django.utils.crypto import get_random_string
from openpyxl import Workbook
from openpyxl.cell.cell import ILLEGAL_CHARACTERS_RE

from opsql.settings import EMAIL_FROM, DOMAIN_NAME
from orders.models import OrdersTasks, Orders, MysqlConfig, UserAccounts, ExportFiles
from orders.utils.msgNotice import MsgPush

logger = logging.getLogger('django')
channel_layer = get_channel_layer()


def update_orders_progress(username, taskid):
    # 当点击全部执行时有效
    # 检查子任务是否都执行完成，如果执行完成，将当前子任务所在的工单进度设置为已完成
    obj = OrdersTasks.objects.filter(taskid=taskid)
    task_progress = obj.values_list('task_progress', flat=True)
    order_id = obj.first().order_id

    if order_id:
        if all([False for i in list(task_progress) if i != '1']):
            data = Orders.objects.get(id=order_id)
            if data.progress != '4':
                data.progress = '4'
                data.updated_at = timezone.now()
                data.save()
                # 推送消息
                msg_push = MsgPush(id=order_id, user=username, type='feedback')
                msg_push.send()


def update_task_progress(id=None, exec_result=None, task_progress=None):
    """更新当前任务的进度"""
    # exec_result的数据格式
    # {'status': 'success', 'rollbacksql': [sql,], 'affected_rows': 1, 'runtime': '1.000s', 'exec_log': ''}
    # 或
    # {'status': 'fail', 'exec_log': ''}
    data = OrdersTasks.objects.get(id=id)
    if exec_result['status'] in ['fail', 'warn']:
        # 标记为失败
        data.task_progress = '3'
        data.task_execlog = exec_result.get('exec_log')
        data.save()
    elif exec_result['status'] == 'success':
        # 执行状态为处理中时，状态变为已完成
        if task_progress == '2':
            rbsql = exec_result.get('rollbacksql')
            affected_rows = int(exec_result.get('affected_rows'))
            consume_time = exec_result.get('runtime')
            exec_log = exec_result.get('exec_log')
            try:
                data.rollback_sql = rbsql
                data.save()
            except Exception as err:
                filename = save_rbsql_as_file(rbsql)
                data.rollback_sql = '\n'.join([
                    '数据超出max_allowed_packet，写入到数据库失败',
                    '备份数据已经以文本的形式进行了保存',
                    '存储路径：',
                    filename
                ])
            finally:
                data.consume_time = consume_time
                data.task_execlog = exec_log
                data.task_progress = '1'
                data.affected_row = affected_rows
                data.save()


def save_rbsql_as_file(rollbacksql):
    """当备份的数据太大时，数据库由于max_allowed_packet问题无法保存，此时保存到文件"""
    if not os.path.exists(r'media/rbsql'):
        os.makedirs('media/rbsql')

    filename = f"media/rbsql/rbsql_{datetime.now().strftime('%Y%m%d%H%M%S%f')}.sql"
    with open(filename, 'w') as f:
        f.write(rollbacksql)
    return filename


class ExportToFiles(object):
    """
    传入参数：user, sql, host, port, id
    """

    def __init__(self, id=None, user=None, sql=None, host=None, port=None, database=None):
        self.user = user
        self.sql = sql
        self.id = id
        self.encoding = 'gbk'
        self.type = 'excel'
        obj = MysqlConfig.objects.get(host=host, port=port)
        self.conn = pymysql.connect(host=obj.host,
                                    user=obj.user,
                                    password=obj.password,
                                    port=obj.port,
                                    db=database,
                                    max_allowed_packet=1024 * 1024 * 1024,
                                    charset='utf8')

        self.execute_log = []
        self.affected_row = 0

        # 文件名
        file_format = OrdersTasks.objects.get(id=self.id).file_format
        num = datetime.now().strftime("%Y%m%d%H%M%S")
        self.title = f'result_{num}'
        self.file = self.title + f'.{file_format}'
        self.zip_file = self.file + '.zip'

    def pull_msg(self, msg):
        # 推送消息
        msg = f"{msg} \n"
        pull_msg = {'status': 3, 'data': msg}
        async_to_sync(channel_layer.group_send)(self.user, {"type": "user.message",
                                                            'text': json.dumps(pull_msg)})

    def set_session_timeout(self):
        with self.conn.cursor() as cursor:
            cursor.execute("set session net_read_timeout=3600")

        with self.conn.cursor() as cursor:
            cursor.execute("set session net_write_timeout=3600")

    def get_count(self):
        # 查询当前SQL的返回的查询数量，返回分页的SQL列表
        status = True
        msg = None
        try:
            count_query = f"select count(*) as count from ({self.sql}) as subquery"
            self.conn.cursorclass = pymysql.cursors.DictCursor
            with self.conn.cursor() as cursor:
                cursor.execute(count_query)
                count = cursor.fetchone()
                self.affected_row = count['count']
            msg = f"SQL导出记录总数：{count['count']}"
        except Exception as err:
            msg = f"导出失败，发现错误：{str(err.args[1])}"
            status = False
        finally:
            self.execute_log.append(msg)
            self.pull_msg(msg)
            return status

    def send_attachments(self, file, encryption_key):
        """发送导出的文件到用户的邮箱"""
        try:
            queryset = OrdersTasks.objects.get(id=self.id)
            if UserAccounts.objects.filter(username=queryset.applicant).exists():
                user_email = [UserAccounts.objects.get(username=queryset.applicant).email]
                self.pull_msg(f'发送附件到用户的邮箱：{user_email[0]}')
                self.execute_log.append(f'发送附件到用户的邮箱：{user_email[0]}')

                # 发送通知里面的域名提示
                domain_name_tips = DOMAIN_NAME['value']

                # 向mail_template.html渲染data数据
                email_html_body = render_to_string('mail/export_file_mail.html',
                                                   {'encryption_key': encryption_key,
                                                    'file': file,
                                                    'domain_name_tips': domain_name_tips
                                                    }
                                                   )

                # 发送邮件
                title = Orders.objects.get(id=queryset.order_id).title
                headers = {'Reply: ': user_email}
                title = 'Re: ' + title
                msg = EmailMessage(subject=title,
                                   body=email_html_body,
                                   from_email=EMAIL_FROM,
                                   to=user_email,
                                   headers=headers
                                   )
                msg.content_subtype = "html"
                msg.send()
                return True
            else:
                self.pull_msg(f'用户邮箱错误，发送失败')
                self.execute_log.append(f'用户邮箱错误，发送失败')
                return False
        except smtplib.SMTPAuthenticationError:
            self.pull_msg('ERROR：Authentication error when sending Email')
            self.execute_log.append('ERROR：Authentication error when sending Email')
            return False

    def compress_file(self):
        # 压缩文件
        msg = f'正在压缩文件：{self.file} ---> {self.zip_file}'
        self.pull_msg(msg)
        self.execute_log.append(msg)

        # 压缩并加密，随机生成18位长度的字符串
        # 7za a -tzip -p123.com b_stru.sql.zip b_stru.sql
        # 需要安装p7zip
        salt = get_random_string(18)
        status, output = subprocess.getstatusoutput(f"7za a -tzip -p{salt} {self.zip_file} {self.file}")
        self.pull_msg(output)
        self.execute_log.append(output)

        msg = f'生成加密密钥：{salt}'
        self.pull_msg(msg)

        # 存储文件
        with open(self.zip_file, 'rb') as f:
            myfile = File(f)
            obj = ExportFiles.objects.create(
                task_id=self.id,
                file_name=self.zip_file,
                file_size=os.path.getsize(self.zip_file),
                files=myfile,
                content_type='xlsx',
                encryption_key=salt
            )
            obj.save()

        # 删除临时文件
        msg = f'删除源文件：{self.file}'
        self.pull_msg(msg)
        self.execute_log.append(msg)
        os.remove(self.file) if os.path.exists(self.file) else None
        os.remove(self.zip_file) if os.path.exists(self.zip_file) else None

        # 发送邮件附件
        status = self.send_attachments(file=obj.files, encryption_key=salt)
        return status

    def export_csv(self):
        # 导出成csv格式
        # 获取列名作为标题
        self.conn.cursorclass = pymysql.cursors.DictCursor
        with self.conn.cursor() as cursor:
            cursor.execute(self.sql)
            title = []
            for column_name in cursor.fetchone():
                title.append(column_name)

        with open(r'%s' % self.file, 'w', newline='', encoding='utf-8') as csvfile:
            fieldnames = title
            writer = csv.DictWriter(csvfile, fieldnames=fieldnames)
            writer.writeheader()

            # 获取数据，并写入到表格
            if self.affected_row <= 100000:
                # 当导出数据量小于10W时，使用fetchall直接读取到内存中
                self.conn.cursorclass = pymysql.cursors.DictCursor
                with self.conn.cursor() as cursor:
                    msg = f'正在导出SQL：{self.sql}'
                    self.pull_msg(msg)
                    self.execute_log.append(msg)

                    cursor.execute(self.sql)
                    rows = cursor.fetchall()

                    msg = f'正在处理数据\n编码为：UTF-8'
                    self.pull_msg(msg)
                    self.execute_log.append(msg)

                    for row in rows:
                        # 过滤掉特殊字符
                        for k, v in row.items():
                            filter_illegal_characters_value = ILLEGAL_CHARACTERS_RE.sub(r'', v) if isinstance(v,
                                                                                                              str) else v
                            row[k] = filter_illegal_characters_value
                        writer.writerow(row)
            elif self.affected_row > 100000:
                # 当导出数据量大于10W时，使用SSCursor进行迭代读取
                self.conn.cursorclass = pymysql.cursors.SSDictCursor
                with self.conn.cursor() as cursor:
                    msg = f'正在导出SQL：{self.sql}'
                    self.pull_msg(msg)
                    self.execute_log.append(msg)

                    cursor.execute(self.sql)
                    while True:
                        row = cursor.fetchone()
                        if row:
                            # 过滤掉特殊字符
                            for k, v in row.items():
                                filter_illegal_characters_value = \
                                    ILLEGAL_CHARACTERS_RE.sub(r'', v) if isinstance(v, str) else v
                                row[k] = filter_illegal_characters_value
                            writer.writerow(row)
                        else:
                            break

    def export_xlsx(self):
        # 导出成xlsx格式
        # num：保存文件的结尾_num标识，为str类型
        # 使用write_only能够有效降低内存的使用
        wb = Workbook(write_only=True)
        wb.encoding = f'{self.encoding}'
        ws = wb.create_sheet()
        ws.title = self.title

        # 获取列名作为标题
        self.conn.cursorclass = pymysql.cursors.DictCursor
        with self.conn.cursor() as cursor:
            cursor.execute(self.sql)
            title = []
            for column_name in cursor.fetchone():
                title.append(column_name)
        ws.append(title)

        # 获取数据，并写入到表格
        if self.affected_row <= 100000:
            # 当导出数据量小于10W时，使用fetchall直接读取到内存中
            self.conn.cursorclass = pymysql.cursors.Cursor
            with self.conn.cursor() as cursor:
                msg = f'正在导出SQL：{self.sql}'
                self.pull_msg(msg)
                self.execute_log.append(msg)

                cursor.execute(self.sql)
                rows = cursor.fetchall()

                msg = f'正在处理数据'
                self.pull_msg(msg)
                self.execute_log.append(msg)

                for row in rows:
                    # 过滤掉特殊字符
                    filter_illegal_characters_row = list(
                        map(
                            (lambda x: ILLEGAL_CHARACTERS_RE.sub(r'', x) if isinstance(x, str) else x), row
                        )
                    )
                    ws.append(filter_illegal_characters_row)
            wb.save(self.file)
        elif self.affected_row > 100000:
            # 当导出数据量大于10W时，使用SSCursor进行迭代读取
            self.conn.cursorclass = pymysql.cursors.SSCursor
            with self.conn.cursor() as cursor:
                msg = f'正在导出SQL：{self.sql}'
                self.pull_msg(msg)
                self.execute_log.append(msg)

                cursor.execute(self.sql)
                while True:
                    row = cursor.fetchone()
                    if row:
                        # 过滤掉特殊字符
                        filter_illegal_characters_row = list(
                            map(
                                (lambda x: ILLEGAL_CHARACTERS_RE.sub(r'', x) if isinstance(x, str) else x), row
                            )
                        )
                        ws.append(filter_illegal_characters_row)
                    else:
                        break
            wb.save(self.file)

    def run(self):
        queryset = OrdersTasks.objects.get(id=self.id)
        status = self.get_count()
        if status:
            if self.affected_row == 0:
                queryset.task_progress = '1'
                queryset.save()
            else:
                start_time = time.time()
                self.set_session_timeout()
                if queryset.file_format == 'xlsx':
                    self.export_xlsx()
                if queryset.file_format == 'csv':
                    self.export_csv()
                if self.compress_file() is True:
                    queryset.task_progress = '1'
                else:
                    queryset.task_progress = '4'
                end_time = time.time()
                consume_time = ''.join((str(round(end_time - start_time, 2)), 's'))
                msg = f'执行耗时：{consume_time}'
                self.execute_log.append(msg)
                self.pull_msg(msg)
                queryset.consume_time = consume_time
        else:
            queryset.task_progress = '3'
        queryset.task_execlog = '\n'.join(self.execute_log)
        queryset.save()
