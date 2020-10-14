# -*- coding:utf-8 -*-
# edit by fuzongfei
import re
from operator import itemgetter

import pymysql

from config import INCEPTION


class InceptionApi(object):
    """Inpcetion语法检查接口"""

    def __init__(self, cfg=None, sqls=None, rds_category=None):
        # 目标数据库连接串配置
        self.cfg = cfg
        # 传入的SQL
        self.sqls = sqls
        # DB类别，mysql还是tidb
        self.rds_category = rds_category
        # inception连接串配置
        self.inception_cfg = {
            'host': INCEPTION['host'],
            'port': INCEPTION['port'],
            'user': '',
            'password': '',
            'db': '',
            'charset': 'utf8mb4',
            'cursorclass': pymysql.cursors.DictCursor
        }

    def set_dynamic_variables(self):
        """
        当选择tidb时，merge_alter_table = false;
        当选择mysql时，merge_alter_table = true;
        """
        variables_cmd = [
            'inception set merge_alter_table = true',
        ]
        # 当为tidb时
        if self.rds_category == 2:
            variables_cmd = [
                'inception set merge_alter_table = false',
            ]

        cnx = pymysql.connect(**self.inception_cfg)
        cursor = cnx.cursor()
        for cmd in variables_cmd:
            cursor.execute(cmd)
        cursor.close()
        cnx.close()

    def check_cnx(self):
        """检查inception的连接状态"""
        try:
            self.inception_cfg['read_timeout'] = 0.5
            pymysql.connect(**self.inception_cfg)
            return True, None
        except pymysql.err.OperationalError as err:
            return False, f"不能访问goInception服务，请联系DBA，错误信息: {err}"

    def conn_incep(self, magic_sqls=None):
        """
        连接到inception，执行语法检查
        此处不能是with cursor，不知道为啥
        """
        self.set_dynamic_variables()
        cnx = pymysql.connect(**self.inception_cfg)
        cursor = cnx.cursor()
        cursor.execute(magic_sqls)
        result = cursor.fetchall()
        cursor.close()
        cnx.close()
        return result

    def run_check(self):
        # fingerprint 开启sql指纹功能。dml语句相似时，可以根据相同的指纹ID复用explain结果，减少远端数据库explain操作，以提高审核速度
        magic_sqls = f"/*--user={self.cfg['user']};--password={self.cfg['password']};" \
                     f"--host={self.cfg['host']};--port={self.cfg['port']};" \
                     f"--check=1;fingerprint=true;*/" \
                     f"inception_magic_start;" \
                     f"\nuse {self.cfg['database']};" \
                     f"\n{self.sqls}" \
                     f"\n;inception_magic_commit;"
        return self.conn_incep(magic_sqls)

    def check_insert_select(self):
        """检查语句中是否包含insert into ... select 语句"""
        rs = re.compile(r'insert([\s\S]+)into([\s\S]+)select(.*)', re.I)
        data = self.run_check()
        for row in data:
            # 匹配到，返回False
            if rs.search(row.get('sql')):
                return False
        return True

    def is_check_pass(self):
        """判断语法检查是否通过"""
        data = self.run_check()
        keys = ['error_level']
        error_level = [itemgetter(*keys)(row) for row in data]
        if all([i == 0 for i in error_level]):
            return True
        return False