# -*- coding:utf-8 -*-
# edit by fuzongfei

import json

from asgiref.sync import async_to_sync
from channels.generic.websocket import WebsocketConsumer


class SqlConsumer(WebsocketConsumer):
    received_value = None
    group_name = None

    def connect(self):
        # 创建channels group， 命名为：用户名，并使用channel_layer写入到redis
        # 以flag为组，每个用户只能消费自己组内的消息
        # flag为前端传入的当前任务的taskid，后台pull消息时，以taskid为组进行写入消息到redis
        self.group_name = self.scope['url_route']['kwargs']['flag']
        async_to_sync(self.channel_layer.group_add)(self.group_name, self.channel_name)

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
        # websocket消费
        # flag用于前端页面的消息定位
        self.send(text_data=json.dumps({'flag': self.group_name, 'data': event["text"]}))

    def disconnect(self, close_code):
        async_to_sync(self.channel_layer.group_discard)(self.group_name, self.channel_name)
