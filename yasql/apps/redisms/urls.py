# -*- coding:utf-8 -*-
# by pandonglin
from django.urls import path, re_path, include
from redisms import views


v1_patterns = [
    # 获取可执行的命令列表
    path(r'cmds', views.RedisLCmds.as_view(), name='v1.redis.cmds'),
    # 获取指定redis授权信息
    path(r'list', views.MyRedisList.as_view(), name='v1.redis.list'),
    # 执行查询
    path(r'exec_cmd', views.ExecRedisCmd.as_view(), name='v1.redis.cmd'),
    # 健康检测
    path(r'<int:pk>/check', views.RedisHealthCheck.as_view(), name='v1.redis.check'),
    # 性能指标
    path(r'<int:pk>/metrics', views.RedisMetrics.as_view(), name='v1.redis.metrics'),
]

urlpatterns = [
    re_path(r'^v1/redis/', include(v1_patterns))
]