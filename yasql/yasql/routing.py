# -*- coding:utf-8 -*-
# edit by fuzongfei

from channels.auth import AuthMiddlewareStack
from channels.routing import URLRouter, ProtocolTypeRouter
from django.urls import path

from sqlorders.consumers import SqlConsumer

application = ProtocolTypeRouter({
    "websocket": AuthMiddlewareStack(
        URLRouter([
            path("ws/sql/<str:flag>/", SqlConsumer),  # flag为task_id
        ])
    )
})
