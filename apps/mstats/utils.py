# -*- coding:utf-8 -*-
# edit by fuzongfei
import configparser
import re

import os
import paramiko
import pymysql
from django.http import Http404
from django.shortcuts import get_object_or_404

from project_manager.models import InceptionHostConfig


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


class MySQLuser_manager(object):
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


class ParamikoOutput(object):
    def __init__(self, ssh_host=None, ssh_port=None, ssh_user=None, ssh_password=None):
        self.ssh_host = ssh_host
        self.ssh_port = ssh_port
        self.ssh_user = ssh_user
        self.ssh_password = ssh_password

    def check_connection(self):
        s = paramiko.SSHClient()
        s.set_missing_host_key_policy(paramiko.AutoAddPolicy())
        try:
            s.connect(hostname=self.ssh_host, port=self.ssh_port, username=self.ssh_user, password=self.ssh_password,
                      timeout=1)
            return True
        except Exception as err:
            return err

    def run(self, cmd):
        try:
            s = paramiko.SSHClient()
            s.set_missing_host_key_policy(paramiko.AutoAddPolicy())
            s.connect(hostname=self.ssh_host, port=self.ssh_port, username=self.ssh_user, password=self.ssh_password,
                      timeout=1)

            msg = [stdin, stdout, stderr] = s.exec_command(cmd)
            out = []
            for item in msg:
                try:
                    for line in item:
                        out.append(line.strip('\n'))
                except Exception as err:
                    pass
            s.close()
            return {'status': 0, 'data': out}
        except Exception as err:
            return {'status': 2, 'msg': '连接超时，请检查SSH或网络'}


class GeneralCParser(object):
    """从传入的字符串中读取配置并生成变量"""

    def __init__(self, backup_dir=None, parser_string=None):
        if isinstance(parser_string, str):
            self.config = configparser.ConfigParser(allow_no_value=True)
            self.config.read_string(parser_string)
            self.backup_dir = backup_dir
            self.check_obj = []

            # 获取mysql backup user
            mysql = self.config['mysql']
            self.mysql_cmd = mysql.get('mysql_tool')
            self.mysql_user = mysql.get('user')
            self.mysql_host = mysql.get('host')
            self.mysql_password = mysql.get('password')
            self.mysql_port = mysql.get('port')

            self.check_obj.append(self.mysql_cmd)
            self.check_obj.append(self.backup_dir)

            # 是否启用compress
            try:
                xb_compress = self.config['compress']
                try:
                    if 'compress' in xb_compress:
                        self.compress = xb_compress.get('compress')
                    if 'compress_chunk_size' in xb_compress:
                        self.compress_chunk_size = xb_compress.get('compress_chunk_size')
                    if 'compress_threads' in xb_compress:
                        self.compress_threads = xb_compress.get('compress_threads')
                    if self.compress and self.compress_chunk_size and self.compress_threads:
                        self.compress_args = f"--compress={self.compress} " \
                                             f"--compress-chunk-size={self.compress_chunk_size} " \
                                             f"--compress-threads={self.compress_threads}"
                except AttributeError as err:
                    self.compress_args = ''
            except KeyError as err:
                self.compress_args = ''

            # 是否启用encrypt
            try:
                xb_encrypt = self.config['encrypt']
                try:
                    if 'encrypt' in xb_encrypt:
                        self.encrypt = xb_encrypt.get('encrypt')
                    if 'encrypt_key' in xb_encrypt:
                        self.encrypt_key = xb_encrypt.get('encrypt_key')
                    if 'encrypt_threads' in xb_encrypt:
                        self.encrypt_threads = xb_encrypt.get('encrypt_threads')
                    if 'encrypt_chunk_size' in xb_encrypt:
                        self.encrypt_chunk_size = xb_encrypt.get('encrypt_chunk_size')
                    if self.encrypt and self.encrypt_key and self.encrypt_threads and self.encrypt_chunk_size:
                        self.encrypt_args = f"--encrypt={self.encrypt} --encrypt-key={self.encrypt_key} " \
                                            f"--encrypt-threads={self.encrypt_threads} " \
                                            f"--encrypt-chunk-size={self.encrypt_chunk_size}"
                except AttributeError as err:
                    self.encrypt_args = ''
            except KeyError as err:
                self.encrypt_args = ''

            # 获取mysqldump
            try:
                mysqldump = self.config['mysqldump']
                self.mysqldump_cmd = mysqldump.get('backup_tool')
                self.mysqldump_backupdir = os.path.join(self.backup_dir, 'mysqldump')
                self.backup_dbs = mysqldump.get('backup_dbs').split(',')
                self.single_table = mysqldump.get('single_table')
                self.dump_options = mysqldump.get('dump_options')

                self.check_obj.append(self.mysqldump_cmd)

            except KeyError as err:
                pass

            # 获取xtrabackup
            try:
                xtrabackup = self.config['xtrabackup']
                self.xtrabackup_cmd = xtrabackup.get('backup_tool')
                self.defaults_file = xtrabackup.get('defaults-file')
                self.xtrabackup_backupdir = os.path.join(self.backup_dir, 'xtrabackup', "`date +%F_%T`")
                self.xtra_options = xtrabackup.get('xtra_options')

                self.check_obj.append(self.xtrabackup_cmd)

            except KeyError as err:
                pass


class CheckCParserValid(GeneralCParser):
    """检测备份配置文件在目标主机的有效性"""

    def __init__(self, ssh_user=None, ssh_password=None, ssh_host=None, ssh_port=None, backup_dir=None,
                 parser_string=None):
        # 获取ssh
        self.ssh_user = ssh_user
        self.ssh_password = ssh_password
        self.ssh_host = ssh_host
        self.ssh_port = ssh_port
        self.parser_string = parser_string
        GeneralCParser.__init__(self, backup_dir, self.parser_string)
        self.result = {}

        self.paramiko_conn = ParamikoOutput(ssh_user=self.ssh_user,
                                            ssh_password=self.ssh_password,
                                            ssh_host=self.ssh_host,
                                            ssh_port=self.ssh_port)

    def check_ssh_conn(self):
        result = self.paramiko_conn.check_connection()
        if result is not True:
            self.result = {'status': 2, 'msg': 'SSH Authentication Failed'}
            return False
        else:
            return True

    def check_obj_exisit(self):
        """检测指定的命令文件和目录是否存在是否存在"""
        is_true = []
        for i in self.check_obj:
            output = self.paramiko_conn.run(f"ls {i} && echo $?")['data']
            print(output)
            if output[-1] != '0':
                self.result = {'status': 2, 'msg': str(output[0]) + ', 请确认文件或目录存在'}
                is_true.append(False)
                break
            else:
                is_true.append(True)
        return is_true

    def check_schema_exisit(self):
        """检测使用mysqldump时，备份的库是否存在"""
        for db in self.backup_dbs:
            cmd = f"{self.mysql_cmd} --user={self.mysql_user} --password='{self.mysql_password}' " \
                  f"--host={self.mysql_host} --port={self.mysql_port} " \
                  f"-e \"select count(*) from information_schema.SCHEMATA where SCHEMA_NAME='{db}'\""

            output = self.paramiko_conn.run(cmd)['data']
            if output[1] == '0':
                self.result = {'status': 2, 'msg': f'mysqldump指定备份的库不存在：{db}'}
                break

    def mkdir_backup_dir(self):
        """创建备份目录，自动创建self.backup_dir/{mysqldump,xtrabackup}"""
        cmd1 = f"if [ ! -d {self.backup_dir}/mysqldump ]; then mkdir {self.backup_dir}/mysqldump;fi "
        cmd2 = f"if [ ! -d {self.backup_dir}/xtrabackup ]; then mkdir {self.backup_dir}/xtrabackup;fi "

        self.paramiko_conn.run(cmd1)
        self.paramiko_conn.run(cmd2)

    def check_mysql_conn(self):
        """检测mysql备份用户是否可以连接到数据库"""
        cmd = f"{self.mysql_cmd} --user={self.mysql_user} --password='{self.mysql_password}' " \
              f"--host={self.mysql_host} --port={self.mysql_port} -e 'select 1'"
        output = self.paramiko_conn.run(cmd)['data']
        if output[0] != '1':
            self.result = {'status': 2, 'msg': str(output[-1])}
            return False
        else:
            return True

    def run(self):
        if self.check_ssh_conn():
            if all(self.check_obj_exisit()):
                if self.check_mysql_conn():
                    if self.config.has_section('mysqldump'):
                        self.check_schema_exisit()
            self.mkdir_backup_dir()
        return self.result if self.result else True


class GeneralBackupCmd(GeneralCParser):
    """生成备份命令并返回"""

    def __init__(self, ssh_user=None, ssh_password=None, ssh_host=None, ssh_port=None, backup_dir=None,
                 parser_string=None):
        # 获取ssh
        self.ssh_user = ssh_user
        self.ssh_password = ssh_password
        self.ssh_host = ssh_host
        self.ssh_port = ssh_port
        self.parser_string = parser_string
        GeneralCParser.__init__(self, backup_dir, self.parser_string)

        self.mysql_args = f"--user={self.mysql_user} --password='{self.mysql_password}' " \
                          f"--host={self.mysql_host} --port={self.mysql_port}"

        self.backup_cmd = {}

    def general_xtrabackup_cmd(self):
        """生成xtrabackup备份命令"""
        cmd = ' '.join((f"{self.xtrabackup_cmd} --defaults-file={self.defaults_file} "
                        f"--backup --target-dir={self.xtrabackup_backupdir} {self.xtra_options}",
                        self.mysql_args, self.compress_args, self.encrypt_args))
        self.backup_cmd['xtrabackup_cmd'] = cmd

    def general_mysqldump_cmd(self):
        """生成mysqldump备份命令"""
        cmd = []
        for db in self.backup_dbs:
            cmd.append(' '.join((self.mysqldump_cmd,
                                 self.mysql_args,
                                 self.dump_options,
                                 db,
                                 f'| gzip > {self.mysqldump_backupdir}/full_{db}_`date +%F_%T`.sql.gz '
                                 f'2>/dev/null')))
        self.backup_cmd['mysqldump_full_cmd'] = cmd

    def run(self):
        if self.config.has_section('mysqldump'):
            self.general_mysqldump_cmd()

        if self.config.has_section('xtrabackup'):
            self.general_xtrabackup_cmd()

        return self.backup_cmd

#         if single_table == 'enable':
#             dbs = tuple(backup_dbs) if len(backup_dbs) > 1 else (repr(backup_dbs))
#             query_cmd = " ".join(('mysql',
#                                   self.get_mysql(),
#                                   "-e",
#                                   f"\"select table_schema, table_name from information_schema.tables "
#                                   f"where table_schema in {dbs}\" 2>/dev/null"))
#             paramiko_output = ParamikoOutput(self.ssh_host, self.ssh_port, self.ssh_user, self.ssh_password)
#             query_output = paramiko_output.run(query_cmd)
#             del query_output[0]
#             mysqldump_single_cmd = []
#             for i in query_output:
#                 schema = i.split('\t')[0]
#                 table = i.split('\t')[1]
#                 mysqldump_single_cmd.append(' '.join((backup_tool,
#                                                       self.get_mysql(),
#                                                       schema,
#                                                       table,
#                                                       dump_options,
#                                                       f'| gzip > {backupdir}/singlebackup_'
#                                                       f'{schema}_{table}_`date +%F_%T`.sql.gz 2>/dev/null')))
#             self.backup_cmd['mysqldump_single_cmd'] = mysqldump_single_cmd
#     except KeyError as err:
#         return ''
#
