# -*- coding:utf-8 -*-
# edit by xff

from operator import itemgetter
import requests
import json

import logging

from config import GAUDIT_API

logger = logging.getLogger('main')


class GAuditApi(object):
    """语法检查接口"""

    def __init__(self, cfg=None, sqls=None, rds_category=None):
        # 目标数据库连接串配置
        self.cfg = cfg
        # 传入的SQL
        self.sqls = sqls
        # DB类别，mysql还是tidb
        self.rds_category = rds_category

    def format(self):
        jsond = {
            "db_user": self.cfg['user'],
            "db_password": self.cfg['password'],
            "db_host": self.cfg['host'],
            "db_port": self.cfg['port'],
            "db": self.cfg['database'],
            "timeout": 3000,
            "custom_audit_parameters": self.cfg['custom_audit_parameters']
        }

        jsond["sqltext"] = self.sqls
        return jsond

    def request(self):
        header = {
            "Content-Type": "application/json"
        }

        try:
            resp = requests.post(
                GAUDIT_API,
                data=json.dumps(self.format(), ensure_ascii=False).encode('utf-8'),
                headers=header
            )
        except requests.exceptions.ConnectionError as err:
            logger.error(err)
            return 500, f"请求审核服务器gAudit异常,请联系DBA,错误信息:{err.args[0]}"

        if resp.status_code == 200:
            return resp.status_code, resp.json()
        return resp.status_code, f"请求审核服务器gAudit异常,请联系DBA;Code:{resp.status_code} Reason:{resp.reason}"

    def check(self):
        """判断语法检查是否通过
        返回值：status, data, msg
        """
        if self.rds_category == 3:
            # clickhouse
            return True, "clickhouse跳过语法审核"
        status_code, data = self.request()
        if status_code == 200:
            if data["code"] != "0000":
                return False, None, data["message"]
            keys = ['level']
            levels = [itemgetter(*keys)(row) for row in data["data"]]
            if all([i == "INFO" for i in levels]):
                return True, data["data"], None
            return False, data["data"], None
        return False, None, data

    def check_dml_scan_limits(self):
        """分析
        分析单条SQL的影响行数
        分析所有SQL的影响行数
        """
        if self.rds_category == 3:
            return True
        limits = db.DML_SCAN_LIMITS
        status, data, _ = self.check()
        if status is False:
            return False

        keys = ['affected_rows']
        affected_rows = [itemgetter(*keys)(row) for row in data]

        # 如果此处需要分别对DML语句中的insert、update、delete做限制，在这里解析data分别处理即可
        single_max_affected_rows = max(affected_rows)
        sum_affected_rows = sum(affected_rows)

        if single_max_affected_rows <= limits.get('single') and sum_affected_rows <= limits.get('sum'):
            return True
        return False
