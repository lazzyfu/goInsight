# -*- coding:utf-8 -*-
# by pandonglin
from redisms import models, serializers
from redisms.redisApi import REDIS_CMDS, RedisApi
from libs.response import JsonResponseV1
from rest_framework.views import APIView
from libs import permissions


class RedisLCmds(APIView):
    def get(self, request, *args, **kwargs):
        return JsonResponseV1(REDIS_CMDS)


class MyRedisList(APIView):
    permission_classes = (permissions.CanExecRedisPermission,)

    def get(self, request, *args, **kwargs):
        my_queryset = models.RedisGroup.objects.filter(rg_group__user=request.user)
        serializer = serializers.RedisGroupSerializers(my_queryset, many=True)
        return JsonResponseV1(serializer.data)


class ExecRedisCmd(APIView):
    permission_classes = (permissions.CanExecRedisPermission,)

    def post(self, request):
        serializer = serializers.RedisSerializer(data=request.data, context={"request": request})
        if serializer.is_valid():
            data = serializer.save()
            return JsonResponseV1(data)
        return JsonResponseV1(message=serializer.errors, code="0001")


class RedisMetrics(APIView):
    permission_classes = (permissions.CanExecRedisPermission,)

    def get(self, request, pk):
        data = {}
        try:
            db = request.query_params.get("db")
            redis_ins = models.RedisConfig.objects.get(pk=pk, group__rg_group__user=request.user)
            redis_api = RedisApi(redis_ins.host, redis_ins.port, password=redis_ins.decrypt_password)
            data = redis_api.get_metrics("db%s" % db)
        except Exception as err:
            pass
        return JsonResponseV1(data)
