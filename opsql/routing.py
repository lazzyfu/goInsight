# -*- coding:utf-8 -*-
# edit by fuzongfei
from channels.auth import AuthMiddlewareStack
from channels.routing import URLRouter, ProtocolTypeRouter
from django.urls import re_path, path

from .consumers import OrderMsgConsumer

application = ProtocolTypeRouter({
    "websocket": AuthMiddlewareStack(
        URLRouter([
            re_path(r"ws/(?P<stream>\w+)/", OrderMsgConsumer),
        ])
    )
})
