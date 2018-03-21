# -*- coding:utf-8 -*-
# edit by fuzongfei

import datetime
import difflib
import hashlib
import time

import mysql.connector as mdb
from celery import shared_task
from channels.layers import get_channel_layer
from django.core.cache import cache
from django.core.mail import EmailMessage
from django.db.models import F
from django.template.loader import render_to_string

from AuditSQL.settings import EMAIL_FROM
from ProjectManager.inception.inception_api import IncepSqlCheck
from ProjectManager.models import OnlineAuditContents, OnlineAuditContentsReply, MonitorSchema
from ProjectManager.utils import update_tasks_status
from UserManager.models import ContactsDetail, UserAccount, Contacts

channel_layer = get_channel_layer()


class GetUserInfo(object):
    def __init__(self, latest_id):
        self.latest_id = latest_id

    def get_user_email(self):
        obj = OnlineAuditContents.objects.get(id=self.latest_id)
        user_list = [obj.proposer, obj.verifier, obj.operate_dba]
        user_email = list(UserAccount.objects.filter(username__in=user_list).values_list('email', flat=True))
        return user_email

    def get_contact_email(self):
        cc = list(OnlineAuditContents.objects.get(pk=self.latest_id).email_cc.split(','))
        contact_email = list(Contacts.objects.filter(contact_id__in=cc).values_list('contact_email', flat=True))
        return contact_email

    # 获取项目组密送成员的邮箱
    def get_bcc_email(self):
        group_id = OnlineAuditContents.objects.get(pk=self.latest_id).group_id
        bcc_email = ContactsDetail.objects.filter(group__group_id=group_id).filter(bcc='1').annotate(
            contact_email=F('contact__contact_email')
        ).values_list('contact_email', flat=True)
        return list(bcc_email)


@shared_task
def send_commit_mail(**kwargs):
    latest_id = kwargs['latest_id']
    userinfo = GetUserInfo(latest_id)

    receiver = userinfo.get_user_email()
    cc = userinfo.get_contact_email()
    bcc = userinfo.get_bcc_email()

    # 向_commit_mail.html渲染data数据
    record = OnlineAuditContents.objects.annotate(group_name=F('group__group_name')).get(pk=latest_id)
    email_html_body = render_to_string('_send_commit_mail.html', {'data': record})

    # 发送邮件
    msg = EmailMessage(subject=record.title,
                       body=email_html_body,
                       from_email=EMAIL_FROM,
                       to=receiver,
                       cc=cc,
                       bcc=bcc,
                       )
    msg.content_subtype = "html"

    # 如果存在上传文件，作为附件发送
    # attachments = UploadFiles.objects.filter(content_id=latest_id).filter(type='0')
    # if attachments:
    #     for i in attachments:
    #         msg.attach_file(r'media/{}'.format(i.files))
    msg.send()


@shared_task
def send_verify_mail(**kwargs):
    latest_id = kwargs['latest_id']
    userinfo = GetUserInfo(latest_id)

    receiver = userinfo.get_user_email()
    cc = userinfo.get_contact_email()
    bcc = userinfo.get_bcc_email()

    # 向mail_template.html渲染data数据
    record = OnlineAuditContents.objects.get(pk=latest_id)
    email_html_body = render_to_string('_send_verify_mail.html', {
        'data': record,
        'user_role': kwargs['user_role'],
        'username': kwargs['username'],
        'addition_info': kwargs['addition_info']
    })

    # 发送邮件
    headers = {'Reply: ': receiver}
    title = 'Re: ' + record.title
    msg = EmailMessage(subject=title,
                       body=email_html_body,
                       from_email=EMAIL_FROM,
                       to=receiver,
                       cc=cc,
                       bcc=bcc,
                       headers=headers)
    msg.content_subtype = "html"
    msg.send()


@shared_task
def send_reply_mail(**kwargs):
    latest_id = kwargs['latest_id']
    reply_id = kwargs['reply_id']
    userinfo = GetUserInfo(latest_id)

    receiver = userinfo.get_user_email()
    cc = userinfo.get_contact_email()
    bcc = userinfo.get_bcc_email()

    title = OnlineAuditContents.objects.get(pk=reply_id).title
    reply_record = OnlineAuditContentsReply.objects.get(pk=latest_id)

    # 向mail_template.html渲染data数据
    email_html_body = render_to_string('_send_reply_mail.html', {
        'reply_record': reply_record,
        'username': kwargs['username'],
    })

    # 发送邮件
    headers = {'Reply: ': receiver}
    title = 'Re: ' + title
    msg = EmailMessage(subject=title,
                       body=email_html_body,
                       from_email=EMAIL_FROM,
                       to=receiver,
                       cc=cc,
                       bcc=bcc,
                       headers=headers)
    msg.content_subtype = "html"
    msg.send()


def connect_db(**kwargs):
    config = {
        'user': kwargs['user'],
        'password': kwargs['password'],
        'host': kwargs['host'],
        'port': kwargs['port'],
        'database': kwargs['database'],
        'raw': 'True'
    }

    return mdb.connect(**config)


@shared_task
def schema_modify_monitor(**kwargs):
    check_time = datetime.datetime.now().strftime("%Y-%m-%d %H:%M:%S")
    conn = connect_db(**kwargs)
    cursor = conn.cursor(dictionary=True)

    query_info = "select table_schema,table_name,group_concat(COLUMN_NAME) as column_name," \
                 "group_concat(COLUMN_DEFAULT) as column_default,group_concat(IS_NULLABLE) as is_nullable," \
                 "group_concat(DATA_TYPE) as data_type,group_concat(CHARACTER_MAXIMUM_LENGTH) as char_length," \
                 "group_concat(COLUMN_TYPE) as column_type,group_concat(COLUMN_COMMENT) as column_comment " \
                 "from columns where table_schema='{schema}' " \
                 "group by table_schema,table_name".format(schema=kwargs['schema'])

    cursor.execute(query_info)

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
    if MonitorSchema.objects.filter(table_schema=kwargs['schema']).first() is None:
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
            MonitorSchema.objects.filter(table_schema=kwargs['schema']).values_list('table_name', flat=True))
        new_data = table_list

        # 找出已删除的表，并处理
        table_remove = list(set(old_data).difference(set(new_data)))
        if table_remove:
            table_change_data.append({'remove': table_remove})
            # 从本地库中删除该表的记录
            MonitorSchema.objects.filter(table_schema=kwargs['schema']).filter(table_name__in=table_remove).delete()

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
            html_data = diff_data.make_file(old_table_stru, new_table_stru, '旧表-表结构', '新表-表结构', context=False,
                                            numlines=5)

        email_html_body = render_to_string('_monitor_table.html',
                                           {'html_data': html_data, 'table_change_data': table_change_data})
        title = '{db}库表变更[来自:{host},检测时间:{check_time}]'.format(db=kwargs['schema'], host=kwargs['describle'],
                                                               check_time=check_time)
        msg = EmailMessage(subject=title,
                           body=email_html_body,
                           from_email=EMAIL_FROM,
                           to=kwargs['receiver'].split(','),
                           )
        msg.content_subtype = "html"
        msg.send()
    cursor.close()
    conn.close()


"""
status = 0: 推送执行结果
statu = 1: 推送执行进度
"""

@shared_task
def get_osc_percent(user, sqlsha1, redis_key, host, database):
    while True:
        if redis_key in cache:
            data = cache.get(redis_key)
            if data == 'start':
                sql = f"inception get osc_percent '{sqlsha1}'"
                incep_sql_check = IncepSqlCheck(sql, host, database, user)

                # 执行SQL
                incep_sql_check.run_status(1)

                # 每2s获取一次
                time.sleep(2)

            elif data == 'end':
                # 删除key
                cache.delete_pattern(redis_key)
                break
        else:
            break


@shared_task
def incep_async_tasks(user, sql, host, database, redis_key, id, taskid):
    incep_sql_check = IncepSqlCheck(sql, host, database, user)

    # 执行SQL
    exec_result = incep_sql_check.run_exec(0)

    # 告诉获取进度的线程退出
    cache.set(redis_key, 'end')

    # 更新任务进度
    update_tasks_status(id=id, taskid=taskid, exec_result=exec_result)


@shared_task
def stop_incep_osc(user, sqlsha1, redis_key, host, database):
    sql = f"inception stop alter '{sqlsha1}'"

    # 执行SQL
    incep_sql_check = IncepSqlCheck(sql, host, database, user)
    incep_sql_check.run_status(0)

    # 告诉获取进度的线程退出
    cache.set(redis_key, 'end')
