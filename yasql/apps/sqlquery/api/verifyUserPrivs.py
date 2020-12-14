# -*- coding:utf-8 -*-
# edit by fuzongfei
from django.db.models import F

from sqlquery import models
from sqlquery.api.extactTableName import get_query_tables


class VerifyUserPrivs(object):
    """检查用户是否有库或表的权限"""

    def __init__(self, username, key, sql):
        self.username = username
        self.cid, self.schema = key.split('___')
        self.sql = sql

    @property
    def _extract_tables_name(self):
        """提取表名，返回list"""
        return get_query_tables(self.sql)

    @property
    def _verify_db(self):
        """判断用户是否有库权限"""
        obj_user = models.DbQueryUserPrivs.objects.filter(
            user__username=self.username,
            schemas__cid__id=self.cid,
            schemas__schema=self.schema
        )
        obj_group = models.DbQueryGroupPrivs.objects.filter(
            user__username=self.username,
            schemas__cid__id=self.cid,
            schemas__schema=self.schema
        )
        if obj_user.count() >= 1 or obj_group.count() >= 1:
            return True
        return False

    @property
    def _get_user_allowed_tables(self):
        """获取当前库允许用户访问的表名"""
        return list(models.DbQueryUserAllowedTables.objects.filter(
            user_privs__user__username=self.username,
            tables__schema__cid=self.cid,
            tables__schema__schema=self.schema
        ).annotate(allowed_tables=F('tables__table')).values_list('allowed_tables', flat=True))

    @property
    def _get_user_deny_tables(self):
        """获取当前库禁止用户访问的表名"""
        return list(models.DbQueryUserDenyTables.objects.filter(
            user_privs__user__username=self.username,
            tables__schema__cid=self.cid,
            tables__schema__schema=self.schema
        ).annotate(deny_tables=F('tables__table')).values_list('deny_tables', flat=True))

    @property
    def _get_group_allowed_tables(self):
        """获取当前库允许用户所在组访问的表名"""
        return list(models.DbQueryGroupAllowedTables.objects.filter(
            group_privs__user__username=self.username,
            tables__schema__cid=self.cid,
            tables__schema__schema=self.schema
        ).annotate(allowed_tables=F('tables__table')).values_list('allowed_tables', flat=True))

    @property
    def _get_group_deny_tables(self):
        """获取当前库禁止用户所在组访问的表名"""
        return list(models.DbQueryGroupDenyTables.objects.filter(
            group_privs__user__username=self.username,
            tables__schema__cid=self.cid,
            tables__schema__schema=self.schema
        ).annotate(deny_tables=F('tables__table')).values_list('deny_tables', flat=True))

    def _rule_engine(self, allowed_tables, deny_tables, input_tables):
        """
        获取当前库允许用户访问的表名
        获取当前库禁止用户访问的表名
        获取当前库允许用户所在组访问的表名
        获取当前库禁止用户所在组访问的表名
        """
        # 如果SQL中的表名在禁止列表，不允许访问
        deny_t = list(set(input_tables).intersection(set(deny_tables)))
        if deny_t:
            return False, f"您没有表{','.join(deny_t)}的访问权限."

        # 如果允许列表为空，允许访问
        # 如果允许列表不为空，如果SQL中的表名全部在列表中，允许访问
        if allowed_tables:
            if not set(input_tables).issubset(set(allowed_tables)):
                return False, f"您没有表{','.join(list(set(input_tables).difference(set(allowed_tables))))}的访问权限."

        return True, input_tables

    @property
    def run(self):
        """判断用户是否有表权限"""
        if not self._verify_db:
            # 没有db的权限
            return False, f"您没有当前库：{self.schema}的查询权限"

        # 有db权限
        allowed_tables = self._get_user_allowed_tables + self._get_group_allowed_tables
        deny_tables = self._get_user_deny_tables + self._get_group_deny_tables
        extract_tables_name = self._extract_tables_name
        return self._rule_engine(
            allowed_tables=allowed_tables,
            deny_tables=deny_tables,
            input_tables=extract_tables_name
        )
