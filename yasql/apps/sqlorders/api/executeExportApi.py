# -*- coding:utf-8 -*-
# edit by fuzongfei
import base64
import csv
import os
import subprocess
import time
from decimal import Decimal

import pymysql
from asgiref.sync import async_to_sync
from celery.utils.log import get_task_logger
from channels.layers import get_channel_layer
from django.core.files import File
from django.utils import timezone
from django.utils.crypto import get_random_string
from openpyxl import Workbook
from openpyxl.cell.cell import ILLEGAL_CHARACTERS_RE

from sqlorders import models, tasks

channel_layer = get_channel_layer()
logger = get_task_logger('celery.logger')


class PullMsg(object):
    def __init__(self, task_id):
        self.task_id = task_id

    def pull(self, msg=None):
        """
        msg = {
        'type': 'execute',   // 取值：execute/processlist/ghost
        'data': str(err)
        }
        """
        async_to_sync(channel_layer.group_send)(
            self.task_id,
            {
                "type": "user.message",
                'text': {
                    'type': 'export',
                    'data': msg
                }
            }
        )


class ExecuteExport(object):
    def __init__(self, config):
        """
        {'host': '127.0.0.1', 'port': 3306, 'charset': 'utf8', 'rds_type': 3, 'database': 'test', 'file_format': 'csv',
        'user': 'yasql_rw', 'password': '123.com', 'task_id': '1e0695520bb640e2ab9dcb8258aeb937'，
        'id': 11, 'username': ''}
        """
        self.config = config
        self.sql = None

        # 文件编码
        self.encoding = 'gbk'

        # 实例化消息推送
        self.pm = PullMsg(self.config['task_id'])

        # 文件名, export_file_${task_id}_${time}
        self.title = str(int(time.time() * 10000))
        self.tmp_file = f"export_file_{config['task_id']}_{self.title}.{config['file_format']}"
        self.tmp_zip_file = f"{self.tmp_file}.zip"

        self.execute_log = []
        self.result = {'status': 'success',
                       'rollbacksql': '',
                       'consuming_time': 0.000,
                       'execute_log': '',
                       'affected_rows': 0
                       }

    def _cnx(self):
        # 新建连接
        cfg = self.config.copy()
        del cfg['id']
        del cfg['username']
        del cfg['task_id']
        del cfg['rds_type']
        del cfg['file_format']
        cnx = pymysql.connect(**cfg)
        with cnx.cursor() as cursor:
            cursor.execute("set session lock_wait_timeout = 30")
        with cnx.cursor() as cursor:
            cursor.execute("set session net_read_timeout=3600")
        with cnx.cursor() as cursor:
            cursor.execute("set session net_write_timeout=3600")
        return cnx

    def correct_int_row(self, x):
        # 解决数字类型变科学计数法以及精度丢失
        if isinstance(x, int) and len(str(x)) > 8:
            return str(x)

        if isinstance(x, Decimal):
            return '{0:f}'.format(x)

        return x

    def _export_to_xlsx(self, cnx):
        # 导出xlsx格式的文件，使用wirte_only能够有效降低内存的使用
        wb = Workbook(write_only=True)
        wb.encoding = self.encoding
        ws = wb.create_sheet()
        ws.title = self.title

        # 推送消息
        msg = f'正在执行导出SQL: {self.sql} \n'
        self.pm.pull(msg=msg)
        self.execute_log.append(msg)

        # 使用游标读取数据，避免数据量过大产生OOM
        cnx.cursorclass = pymysql.cursors.SSCursor
        with cnx.cursor() as cursor:
            cursor.execute(self.sql)
            # 推送消息
            msg = f'正在处理并生成XLSX数据 \n'
            self.pm.pull(msg=msg)
            self.execute_log.append(msg)
            # 标题
            ws.append([x[0] for x in cursor.description])
            # 返回行数
            self.result['affected_rows'] = cursor.rownumber
            # 操作数据
            while True:
                row = cursor.fetchone()
                if not row:
                    break

                # 过滤掉特殊字符
                filter_illegal_characters_row = list(
                    map(
                        (lambda x: ILLEGAL_CHARACTERS_RE.sub(r'', x) if isinstance(x, str) else x), row
                    )
                )
                # 处理科学计数法
                _row = [self.correct_int_row(x) for x in filter_illegal_characters_row]
                ws.append(_row)

        # 保存到文件
        wb.save(self.tmp_file)

    def _exprt_to_csv(self, cnx):
        # 导出csv的文件
        # 推送消息
        msg = f'正在执行导出SQL: {self.sql} \n'
        self.pm.pull(msg=msg)
        self.execute_log.append(msg)

        # 打开csv文件
        with open(self.tmp_file, 'w', newline='', encoding='utf-8') as csvfile:
            # 使用游标读取数据，避免数据量过大产生OOM
            cnx.cursorclass = pymysql.cursors.SSDictCursor
            with cnx.cursor() as cursor:
                cursor.execute(self.sql)
                # 推送消息
                msg = f'正在处理并生成CSV数据 \n'
                self.pm.pull(msg=msg)
                self.execute_log.append(msg)
                # 标题
                fieldnames = ([x[0] for x in cursor.description])
                # 返回行数
                self.result['affected_rows'] = cursor.rownumber
                # 实例化csv
                writer = csv.DictWriter(csvfile, fieldnames=fieldnames)
                writer.writeheader()
                # 操作数据
                while True:
                    row = cursor.fetchone()
                    if not row:
                        break
                    # 过滤掉特殊字符
                    for k, v in row.items():
                        filter_illegal_characters_value = ILLEGAL_CHARACTERS_RE.sub(r'', v) \
                            if isinstance(v, str) else v
                        row[k] = filter_illegal_characters_value
                    writer.writerow(row)

    def compress_file(self):
        # 推送消息
        msg = f'正在压缩文件: {self.tmp_file} -> {self.tmp_zip_file} \n'
        self.pm.pull(msg=msg)
        self.execute_log.append(msg)

        # 压缩并加密，随机生成24位长度的字符串
        # 需要安装p7zip 7za a -tzip -p123.com b_stru.sql.zip b_stru.sql
        salt = get_random_string(24)
        status, output = subprocess.getstatusoutput(f"7za a -tzip -p{salt} {self.tmp_zip_file} {self.tmp_file}")
        self.pm.pull(msg=output)
        self.execute_log.append(output)

        if status == 0:
            # 存储文件
            with open(self.tmp_zip_file, 'rb') as f:
                myfile = File(f)
                obj = models.DbExportFiles.objects.create(
                    task_id=self.config['id'],
                    file_name=self.tmp_zip_file,
                    file_size=os.path.getsize(self.tmp_zip_file),
                    files=myfile,
                    content_type=self.config['file_format'],
                    encryption_key=salt,
                    created_at=timezone.now(),
                )
        else:
            self.result['status'] = 'fail'

        # 删除临时文件
        for f in [self.tmp_file, self.tmp_zip_file]:
            msg = f'删除临时文件: {f} \n'
            self.pm.pull(msg=msg)
            self.execute_log.append(msg)
            if os.path.exists(f):
                os.remove(f)

        # 发送消息
        tasks.msg_notice.delay(
            pk=models.DbOrdersExecuteTasks.objects.get(pk=self.config['id']).order_id,
            op='_export',
            username=self.config['username'],
            export_file_name=base64.b64encode(obj.file_name.encode()).decode(),
            export_file_encryption_key=obj.encryption_key,
            export_sql=self.sql
        )

    def run(self, sql):
        try:
            cnx = self._cnx()
            self.sql = sql
            start_time = time.time()
            if self.config['file_format'].upper() == 'XLSX':
                self._export_to_xlsx(cnx)

            if self.config['file_format'].upper() == 'CSV':
                self._exprt_to_csv(cnx)
            end_time = time.time()
            self.result['consuming_time'] = round(float(end_time - start_time), 3)  # 简单计算耗时

            self.compress_file()

            msg = f"执行耗时: {self.result['consuming_time']} \n"
            self.pm.pull(msg=msg)
            self.execute_log.append(msg)
        except Exception as err:
            logger.error(err)
            self.pm.pull(msg=err)
            self.execute_log.append(err)
            self.result['status'] = 'fail'
        finally:
            self.result['execute_log'] = '\n'.join(self.execute_log)
            return self.result
