# -*- coding:utf-8 -*-
# edit by fuzongfei
import json
import socket
import threading

import paramiko
from gevent.socket import wait_read
from asgiref.sync import async_to_sync
from channels.generic.websocket import WebsocketConsumer
from channels.layers import get_channel_layer
# from .settings import WEBSHELL_HOST, WEBSHELL_PORT, WEBSHELL_USER, WEBSHELL_PASSWORD

# from mstats.models import WebShellOpLog
from webshell.models import WebShellOpLog

channel_layer = get_channel_layer()


class EchoConsumer(WebsocketConsumer):
    def connect(self):
        # 创建channels group， 命名为：用户名，并使用channel_layer写入到redis
        async_to_sync(self.channel_layer.group_add)(self.scope['user'].username, self.channel_name)

        # 返回给receive方法处理
        self.accept()

    def receive(self, text_data=None, bytes_data=None):
        async_to_sync(self.channel_layer.group_send)(
            self.scope['user'].username,
            {
                "type": "user.message",
                "text": text_data,
            },
        )

    def user_message(self, event):
        # 消费
        self.send(text_data=event["text"])

    def disconnect(self, close_code):
        async_to_sync(self.channel_layer.group_discard)(self.scope['user'].username, self.channel_name)


# 定义SESSION池
PARAMIKO_SESSION_POOL = {}
# 定义连接，如果连接存在，跳过创建
PARAMIKO_CONN = {}


class ParamikoConn(object):
    def __init__(self, user):
        self.user = user

    def init_conn(self):
        """初始化连接"""
        conn = paramiko.Transport('127.0.0.1', 22)
        conn.start_client()
        conn.auth_password(username='webshell', password='oWd5CRnELlg1ZBi9')
        PARAMIKO_CONN['conn'] = conn

    def create_session(self):
        """创建会话"""
        if not PARAMIKO_CONN.get('conn') or not PARAMIKO_CONN.get('conn').is_active():
            self.init_conn()
        channel = PARAMIKO_CONN.get('conn').open_session()
        channel.get_pty(term='xterm')
        channel.invoke_shell()
        channel.setblocking(False)
        channel.settimeout(0.0)
        return channel

    def session_pool(self):
        """
        用户会话池
        一个用户一个会话
        如果该用户websocket未断开，使用已存在的session
        如果该用户的webscoket重连，为该用户新建session
        """
        if not PARAMIKO_SESSION_POOL.get(self.user):
            PARAMIKO_SESSION_POOL[self.user] = self.create_session()


class MyThread(threading.Thread):
    def __init__(self, channel, send):
        threading.Thread.__init__(self)
        self.channel = channel
        self.send = send

    def run(self):
        # 当用户会话没有退出时
        while not self.channel.exit_status_ready():
            # 等待读取
            wait_read(self.channel.fileno())
            try:
                data = self.channel.recv(65535).decode('utf8')
                self.send(text_data=data)
            except socket.timeout as err:
                pass
            except ValueError as err:
                # 客户端关闭了websocket，无法发送结果
                pass
        return False


class SSHTerminalConsumer(WebsocketConsumer):
    cmd_string = ''
    channel = None

    def connect(self):
        # 为当前登陆用户创建一个session
        paramiko_conn = ParamikoConn(self.scope['user'].username)
        paramiko_conn.session_pool()

        # 获取当前用户连接的的channel
        self.channel = PARAMIKO_SESSION_POOL.get(self.scope['user'].username)

        # 记录用户登录
        self.logging_op('Login.')

        # 启动一个线程，负责读取数据
        t1 = MyThread(self.channel, self.send)
        t1.setDaemon(True)
        t1.start()

        self.accept()

    def receive(self, text_data=None, bytes_data=None):
        data = json.loads(text_data)
        resize = data.get('resize')
        if resize and len(resize) == 2:
            # 设置paramiko的窗口尺寸
            self.channel.resize_pty(width=resize[0], height=resize[1])
        else:
            self.handler_write(data.get('data'))

    def handler_write(self, data):
        """
        写入数据
        from websocket to paramiko
        """
        try:
            # 记录用户操作
            if data.endswith('\r'):
                if len(data.strip()) > 1:
                    self.cmd_string = data

                if self.cmd_string.strip():
                    self.logging_op(self.cmd_string)
                self.cmd_string = ''
            else:
                self.cmd_string += data
            # 发送数据到paramiko
            self.channel.send(data)
        except socket.error as err:
            self.logging_op(f'Logout. {str(err)}')

    def logging_op(self, cmd_string):
        # 记录用户的操作命令
        WebShellOpLog.objects.create(
            user=self.scope['user'].username,
            session_id=self.channel.get_id(),
            op_cmd=cmd_string,
        )

    def disconnect(self, close_code):
        # websocket断开时，关闭用户channel，并从pool中移除session
        # 记录用户登出
        self.logging_op('Logout.')
        self.channel.close()
        if PARAMIKO_SESSION_POOL.get(self.scope['user'].username):
            del PARAMIKO_SESSION_POOL[self.scope['user'].username]
