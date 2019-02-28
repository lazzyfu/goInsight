# -*- coding:utf-8 -*-
# edit by fuzongfei
import os
import subprocess

import pymysql
import sqlparse
from django.utils.crypto import get_random_string

from sqlaudit.settings import SOAR_CONFIG


class SoarAnalyze(object):
    def __init__(self, user, password, host, port, schema, type, contents):
        # 远程数据库模拟生产环境
        self.user = user
        self.password = password
        self.host = host
        self.port = port if isinstance(port, int) else int(port)
        self.schema = schema
        self.type = type
        self.contents = contents

        self.soar_user = SOAR_CONFIG.get('testenv').get('SOAR_USER')
        self.soar_password = SOAR_CONFIG.get('testenv').get('SOAR_PASSWORD')
        self.soar_host = SOAR_CONFIG.get('testenv').get('SOAR_HOST')
        self.soar_port = SOAR_CONFIG.get('testenv').get('SOAR_PORT')
        self.soar_arguments = ' '.join(SOAR_CONFIG.get('arguments'))

    def check_conn(self):
        try:
            conn = pymysql.connect(user=self.user,
                                   host=self.host,
                                   password=self.password,
                                   port=self.port,
                                   connect_timeout=2)

            if conn:
                return True, None
            conn.close()
        except pymysql.Error as err:
            return False, str(err)

    def execute(self, cmd):
        status, output = subprocess.getstatusoutput(cmd)
        return status, output

    def advisor(self):
        # 返回list格式
        advisor_result = []
        for sql in sqlparse.split(self.contents):
            # 创建临时文件
            tmp_file = f'media/tmp/soar/{get_random_string(32)}.sql'
            with open(tmp_file, encoding='utf-8', mode='w') as file:
                file.write(sql)

            cmd = f'soar -query {tmp_file} -online-dsn ' \
                f'"{self.user}:{self.password}@{self.host}:{self.port}/{self.schema}" ' \
                f'-test-dsn "{self.soar_user}:{self.soar_password}@{self.soar_host}:{self.soar_port}" '

            soar_cmd = ' '.join([cmd, self.soar_arguments])
            p = subprocess.Popen(soar_cmd,
                                 shell=True,
                                 stdin=subprocess.PIPE,
                                 stdout=subprocess.PIPE,
                                 stderr=subprocess.PIPE,
                                 encoding='utf-8')
            out, err = p.communicate()
            # 删除临时文件
            os.remove(tmp_file)
            advisor_result.append(' '.join([out, err]))
        return advisor_result

    def run(self):
        context = None
        status, msg = self.check_conn()
        # 如果tmp临时目录不存在，创建
        if not os.path.exists('media/tmp/soar'):
            os.makedirs('media/tmp/soar')
        if status:
            if self.type == 'advisor':
                context = {'status': 0, 'data': self.advisor()}

        else:
            context = {'status': 2, 'msg': msg}

        return context