# -*- coding:utf-8 -*-
# edit by fuzongfei

import difflib
import hashlib
from datetime import datetime

import pymysql
from celery import shared_task
from django.core.mail import EmailMessage
from django.template.loader import render_to_string

from AuditSQL.settings import EMAIL_FROM
from ProjectManager.models import InceptionHostConfig, MonitorSchema


@shared_task
def monitor_schema_modify(**kwargs):
    """监控数据库表结构变更并提供邮件通知"""
    host = kwargs.get('host')
    schema = kwargs.get('schema')
    receiver = kwargs.get('receiver')
    check_time = datetime.now().strftime("%Y-%m-%d %H:%M:%S")

    config = InceptionHostConfig.objects.get(host=host, is_enable=0)

    cnx = pymysql.connect(host=config.host,
                          user=config.user,
                          password=config.password,
                          port=config.port,
                          charset="utf8",
                          cursorclass=pymysql.cursors.DictCursor)

    try:
        with cnx.cursor() as cursor:
            query = f"select table_schema,table_name,group_concat(COLUMN_NAME) as column_name," \
                    f"group_concat(COLUMN_DEFAULT) as column_default,group_concat(IS_NULLABLE) as is_nullable," \
                    f"group_concat(DATA_TYPE) as data_type,group_concat(CHARACTER_MAXIMUM_LENGTH) as char_length," \
                    f"group_concat(COLUMN_TYPE) as column_type,group_concat(COLUMN_COMMENT) as column_comment " \
                    f"from information_schema.columns where table_schema='{schema}' " \
                    f"group by table_schema,table_name"
            cursor.execute(query)

            source_info = []
            table_list = []
            diff_old_data = ''
            diff_new_data = ''
            table_change_data = []

            for row in cursor.fetchall():
                table_schema = row['table_schema']
                table_name = row['table_name']

                md5_source = ''.join(str(row.values()))
                md5_sum = hashlib.md5(md5_source.encode('utf8')).hexdigest()
                source_info.append({'table_schema': table_schema, 'table_name': table_name, 'md5_sum': md5_sum})
                table_list.append(table_name)

            # 如果当前库没有记录，则进行初始化全量同步
            if MonitorSchema.objects.filter(table_schema=schema).first() is None:
                for row in source_info:
                    table_schema = row['table_schema']
                    table_name = row['table_name']

                    query_table_stru = "show create table {}".format('.'.join((table_schema, table_name)))
                    cursor.execute(query_table_stru)
                    for i in cursor:
                        table_stru = i['Create Table']
                        row['table_stru'] = str(table_stru)
                        MonitorSchema.objects.create(**row)
            else:
                # 如果存在，开始核验数据
                old_data = list(
                    MonitorSchema.objects.filter(table_schema=schema).values_list('table_name', flat=True))
                new_data = table_list

                # 找出已删除的表，并处理
                table_remove = list(set(old_data).difference(set(new_data)))
                if table_remove:
                    table_change_data.append({'remove': table_remove})
                    # 从本地库中删除该表的记录
                    MonitorSchema.objects.filter(table_schema=schema).filter(
                        table_name__in=table_remove).delete()

                # 找出新增的表，并处理
                table_add = list(set(new_data).difference(set(old_data)))
                if table_add:
                    for i in table_add:
                        for j in source_info:
                            if i in j.values():
                                table_change_data.append({'add': j})
                                table_schema = j['table_schema']
                                table_name = j['table_name']
                                query_table_stru = "show create table {}".format('.'.join((table_schema, table_name)))
                                cursor.execute(query_table_stru)
                                for x in cursor:
                                    table_stru = x['Create Table']
                                    j['table_stru'] = str(table_stru)
                                    MonitorSchema.objects.create(**j)

                # 找出相同的表，并核验表结构
                table_intersection = list(set(old_data).intersection(set(new_data)))
                for row in source_info:
                    table_schema = row['table_schema']
                    table_name = row['table_name']
                    new_md5_sum = row['md5_sum']

                    if table_name in table_intersection:
                        old_table = MonitorSchema.objects.get(table_schema=table_schema, table_name=table_name)
                        if new_md5_sum != old_table.md5_sum:
                            query_table_stru = "show create table {}".format('.'.join((table_schema, table_name)))
                            cursor.execute(query_table_stru)
                            for i in cursor:
                                table_stru = i['Create Table']
                                diff_old_data += old_table.table_stru + '\n' * 3
                                diff_new_data += table_stru + '\n' * 3
                                # 更新新表表结构到本地
                                MonitorSchema.objects.update_or_create(table_schema=table_schema, table_name=table_name,
                                                                       defaults={'table_stru': table_stru,
                                                                                 'md5_sum': new_md5_sum})

                if (diff_old_data and diff_new_data) or table_change_data:
                    html_data = ''
                    if diff_old_data and diff_new_data:
                        diff_data = difflib.HtmlDiff(tabsize=2)
                        old_table_stru = list(diff_old_data.split('\n'))
                        new_table_stru = list(diff_new_data.split('\n'))
                        html_data = diff_data.make_file(old_table_stru, new_table_stru, '旧表-表结构', '新表-表结构',
                                                        context=False,
                                                        numlines=5)

                    email_html_body = render_to_string('_monitor_table.html',
                                                       {'html_data': html_data,
                                                        'table_change_data': table_change_data,
                                                        'host': host,
                                                        'schema': schema
                                                        })
                    title = f'表结构检测[{host}-{schema}_{check_time}]'
                    msg = EmailMessage(subject=title,
                                       body=email_html_body,
                                       from_email=EMAIL_FROM,
                                       to=receiver.split(','),
                                       )
                    msg.content_subtype = "html"
                    msg.send()
                cursor.close()
    finally:
        cnx.close()
