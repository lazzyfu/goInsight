# -*- coding:utf-8 -*-
# edit by fuzongfei
import re

import pymysql
from django.http import Http404
from django.shortcuts import get_object_or_404

from ProjectManager.models import InceptionHostConfig


def get_mysql_user_info(host):
    data = InceptionHostConfig.objects.get(host=host)
    conn = pymysql.connect(host=host,
                           user=data.user,
                           password=data.password,
                           charset='utf8',
                           port=data.port,
                           cursorclass=pymysql.cursors.DictCursor, )

    # 此方法可能受中间件影响，查询结果不准确
    # mysql_version = int(conn.server_version.split('.', 2)[1])
    version_query = 'select version();'

    try:
        with conn.cursor() as cursor:
            cursor.execute(version_query)
            mysql_version = int(cursor.fetchone().get('version()').split('.', 2)[1])

            id = 1

            user_query = "select user from mysql.user"
            cursor.execute(user_query)
            user_list = []

            for row in cursor.fetchall():
                user_list.append(row.get('user'))

            user_set = list(set(user_list))

            user_dict = []
            for i in user_set:
                user_dict.append({'id': id, 'pid': 0, 'privileges': '', 'schema': '', 'user': i})
                id += 1

            if mysql_version > 6:
                user_info_query = "select concat(\"'\",user,\"'\",'@',\"'\",host,\"'\") as username, " \
                                  "concat(left(authentication_string,5),'...',right(authentication_string,2)) " \
                                  "as password, password_expired from mysql.user"
            else:
                user_info_query = "select concat(\"'\",user,\"'\",'@',\"'\",host,\"'\") as username," \
                                  "concat(left(password,5),'...',right(password,2)) as password, " \
                                  "password_expired from mysql.user"
            cursor.execute(user_info_query)

            privileges_dict = []

            for row in cursor.fetchall():
                user = row.get('username')
                password = row.get('password')
                password_expired = row.get('password_expired')
                username = user.split('@')[0].replace("'", '').strip()
                user_host = user.split('@')[1].replace("'", '').strip()

                privileges_query = f"show grants for {user}"
                cursor.execute(privileges_query)
                pid = 0
                for i in cursor.fetchall():
                    for k, v in i.items():
                        data = re.split(r' TO ', v.replace('GRANT', '').strip())[0].split(r' ON ')
                        privileges = data[0]
                        schema = data[1]
                        for j in user_dict:
                            if username == j.get('user'):
                                pid = j.get('id')
                        privileges_dict.append({
                            'id': id,
                            'pid': pid,
                            'user': username,
                            'host': user_host,
                            'password': password,
                            'privileges': privileges,
                            'schema': schema,
                            'password_expired': password_expired
                        })
                        id += 1
            data = user_dict + privileges_dict
            return data
    finally:
        conn.close()


class MySQLUserManager(object):
    def __init__(self, kwargs):
        self.db_host = kwargs.get('db_host')
        # username = user@host
        self.username = kwargs.get('username')
        self.schema = kwargs.get('schema')
        self.password = kwargs.get('password')
        self.privileges = kwargs.get('privileges')

        data = InceptionHostConfig.objects.get(host=self.db_host)
        self.dst_user = '@'.join((data.user, data.host))
        self.conn = pymysql.connect(host=data.host,
                                    user=data.user,
                                    password=data.password,
                                    port=data.port,
                                    charset='utf8')

    def flush_privileges(self):
        """刷新权限"""
        with self.conn.cursor() as cursor:
            flush_query = f"flush privileges"
            cursor.execute(flush_query)
            cursor.close()

    def rollback_user(self):
        """移除执行错误后，已经生成的用户"""
        with self.conn.cursor() as cursor:
            rollback_user_query = f"drop user {self.username}"
            cursor.execute(rollback_user_query)
            cursor.close()

    def priv_modify(self):
        try:
            with self.conn.cursor() as cursor:
                # 先移除all privileges
                revoke_query = f"revoke all ON {self.schema} from {self.username}"
                cursor.execute(revoke_query)

                # 赋予新权限
                modify_query = f"grant {self.privileges} on {self.schema} to {self.username}"
                cursor.execute(modify_query)
                return {'status': 0, 'msg': '修改成功'}
        except self.conn.OperationalError as error:
            return {'status': 2, 'msg': f'授权失败，请检查{self.dst_user}是否有with grant option权限'}
        except self.conn.ProgrammingError as error:
            return {'status': 2, 'msg': str(error)}
        finally:
            self.flush_privileges()
            self.conn.close()

    def new_host(self):
        try:
            with self.conn.cursor() as cursor:
                # 创建用户主机
                new_host_query = f"create user {self.username} identified by \"{self.password}\""
                cursor.execute(new_host_query)

                # 赋予新权限
                modify_query = f"grant {self.privileges} on {self.schema} to {self.username}"
                cursor.execute(modify_query)
                return {'status': 0, 'msg': '创建成功'}
        except self.conn.OperationalError as error:
            self.rollback_user()
            return {'status': 2, 'msg': f'授权失败，请检查{self.dst_user}是否有with grant option权限'}
        except self.conn.ProgrammingError as error:
            self.rollback_user()
            return {'status': 2, 'msg': f'语法错误，{str(error)}'}
        except self.conn.InternalError as error:
            return {'status': 2, 'msg': f'创建用户失败，请检查主机是否存在, {str(error)}'}
        finally:
            self.flush_privileges()
            self.conn.close()

    def delete_host(self):
        try:
            with self.conn.cursor() as cursor:
                delete_host_query = f"drop user {self.username}"
                cursor.execute(delete_host_query)
                cursor.close()
                return {'status': 0, 'msg': '删除成功'}
        except Exception as error:
            return {'status': 2, 'msg': str(error)}
        finally:
            self.flush_privileges()
            self.conn.close()


def check_mysql_conn_status(fun):
    def wapper(request, *args, **kwargs):
        host = request.GET.get('host')
        data = get_object_or_404(InceptionHostConfig, host=host)

        try:
            conn = pymysql.connect(user=data.user,
                                   host=host,
                                   password=data.password,
                                   port=data.port,
                                   use_unicode=True,
                                   connect_timeout=1)

            if conn:
                return fun(request, *args, **kwargs)
            conn.close()
        except pymysql.Error as err:
            raise Http404

    return wapper
