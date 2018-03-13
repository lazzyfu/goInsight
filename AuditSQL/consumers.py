# -*- coding:utf-8 -*-
# edit by fuzongfei

from asgiref.sync import async_to_sync
from channels.generic.websocket import WebsocketConsumer

from channels.layers import get_channel_layer

channel_layer = get_channel_layer()


class EchoConsumer(WebsocketConsumer):
    def connect(self):
        # 创建channels group， 命名为：用户名，并使用channel_layer写入到redis
        async_to_sync(self.channel_layer.group_add)(self.scope['user'].username, self.channel_name)

        # 返回给receive方法处理
        self.accept()

    def receive(self, text_data):
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

    # def disconnect(self):
    #     async_to_sync(self.channel_layer.group_discard)(self.scope['user'].username, self.channel_name)
