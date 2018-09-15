# -*- coding:utf-8 -*-
# edit by fuzongfei
from channels.auth import AuthMiddlewareStack
from channels.routing import URLRouter, ProtocolTypeRouter
from django.urls import path

from .consumers import EchoConsumer, SSHTerminalConsumer

application = ProtocolTypeRouter({
    "websocket": AuthMiddlewareStack(
        URLRouter([
            path(r"ws/", EchoConsumer),
            path(r"ssh/", SSHTerminalConsumer),
        ])
    )
})
