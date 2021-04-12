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
    # 实例信息
    path(r'<int:pk>/info', views.RedisInfo.as_view(), name='v1.redis.info'),
    # 性能分析
    path(r'<int:pk>/analysis', views.RedisAnalysis.as_view(), name='v1.redis.analysis'),
]

urlpatterns = [
    re_path(r'^v1/redis/', include(v1_patterns))
]