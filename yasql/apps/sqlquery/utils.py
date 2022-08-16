# -*- coding:utf-8 -*-
# edit by xff

import sqlparse


# 执行前,删除语句开头的注释
def remove_sql_comment(sql):
    for stmt in sqlparse.split(sql.rstrip(';')):
        statement = sqlparse.parse(stmt)[0]
        comment = statement.token_first()
        if isinstance(comment, sqlparse.sql.Comment):
            return statement.value.replace(comment.value, '')
        return statement.value
