# -*- coding:utf-8 -*-
# edit by xff
import pymysql
import sqlparse

from collections import Counter


def remove_sql_comment(sql):
    # 执行前,删除语句开头的注释
    for stmt in sqlparse.split(sql):
        statement = sqlparse.parse(stmt)[0]
        comment = statement.token_first()
        if isinstance(comment, sqlparse.sql.Comment):
            return statement.value.replace(comment.value, '')
        return statement.value


def verify_sql_type(sqls=None, sql_type=None):
    # 提交工单/语法检查时，判断SQL的类型是DDL还是DML
    # 保证分开提交，ddl工单提交ddl语句，dml工单提交dml语句

    result = []

    for sql in sqlparse.split(sqls):
        """解析SQL类型，返回是DML还是DDL"""
        res = sqlparse.parse(remove_sql_comment(sql))
        syntax_type = res[0].token_first().ttype.__str__()
        if syntax_type == 'Token.Keyword.DDL':
            result.append('DDL')
        if syntax_type == 'Token.Keyword.DML':
            result.append('DML')

    if not all([i == sql_type for i in result]):
        if sql_type == 'DDL':
            return False, 'DDL模式下, 不支持SELECT|UPDATE|DELETE|INSERT等语句'
        if sql_type == 'DML':
            return False, 'DML模式下, 不支持ALTER|CREATE|TRUNCATE|DROP等语句'
    return True, None


def check_export_column_unique(config, sqls):
    # ERROR 1060 (42S21): Duplicate column name 'xxx'
    # 判断提交导出工单的列是否重复，重复不允许提交。dict数据类型不支持key重复

    conn = pymysql.connect(**config)
    for sql in sqlparse.split(sqls):
        explain_query = f"explain select count(*) as count from ({sql.strip(';')}) as subquery"
        try:
            with conn.cursor() as cursor:
                cursor.execute(explain_query)
        except Exception as err:
            if err.args[0] == 1060:
                return False, f"SELECT列名重复, 请使用AS更换别名，错误信息：{err.args[1]}"
            else:
                return False, err
    conn.close()
    return True, None


def handle_duplicate_column(column):
    column_count = Counter(column)
    _column = []
    _ = 0
    for col in column:
        if column_count[col] > 1:
            col = f"{col}_{str(_)}"
            _ += 1
        _column.append(col)
    return _column
