# -*- coding:utf-8 -*-
# edit by fuzongfei
import datetime

from rest_framework import status
from rest_framework.response import Response
from rest_framework.views import APIView

from orders.models import SysEnvironment, OnlineVersion
from orders.permissions import anyof, CanCommitPermission, CanExecutePermission, CanAuditPermission
from orders.serializers import BeautifySQLSerializer, GetSchemasSerializer, SyntaxCheckSerializer, \
    OnlineVersionDetailSerializer
from orders.serializers.commitSerializers import OnlineVersionListSerializer


class BeautifySQLView(APIView):
    """
    美化SQL
    判断SQL类型（DML还是DDL），并分别进行美化
    最后合并返回
    """

    def post(self, request):
        serializer = BeautifySQLSerializer(data=request.data)
        if serializer.is_valid():
            data = serializer.beautify()
            return Response(data={'code': 0, 'data': data}, status=status.HTTP_200_OK)
        else:
            errors = [str(v[0]) for k, v in serializer.errors.items()]
            data = {'code': 2, 'data': '\n'.join(errors)}
            return Response(data=data, status=status.HTTP_400_BAD_REQUEST)


class SyntaxCheckView(APIView):
    """语法检查"""

    def post(self, request):
        serializer = SyntaxCheckSerializer(data=request.data)
        if serializer.is_valid():
            s, data = serializer.check(request)
            code = 0 if s else 2
            return Response(data={'code': code, 'data': data}, status=status.HTTP_200_OK)
        else:
            errors = [str(v[0]) for k, v in serializer.errors.items()]
            data = {'code': 2, 'data': '\n'.join(errors)}
            return Response(data=data, status=status.HTTP_400_BAD_REQUEST)


class GetSchemasView(APIView):
    """获取指定环境的schema"""

    def post(self, request):
        serializer = GetSchemasSerializer(data=request.data)
        if serializer.is_valid():
            data = serializer.query()
            return Response(data={'code': 0, 'data': data}, status=status.HTTP_200_OK)
        else:
            errors = [str(v[0]) for k, v in serializer.errors.items()]
            data = {'code': 2, 'data': '\n'.join(errors)}
            return Response(data=data, status=status.HTTP_400_BAD_REQUEST)


class GetSysEnviView(APIView):
    """获取系统环境"""

    def get(self, request):
        queryset = SysEnvironment.objects.all().values('envi_id', 'envi_name')
        return Response(data=queryset, status=status.HTTP_200_OK)


class OnlineVersionNoExpireView(APIView):
    """获取上线版本"""

    def get(self, request):
        """
        如果当前任务的提交时间大于任务设置的过期时间，不允许选择该任务
        is_disable：是否禁用，0：否，1：是
        """
        before_30_days = (datetime.datetime.now() - datetime.timedelta(days=30)).strftime("%Y-%m-%d %H:%M:%S")
        query = f"select id, version, if(now() > date_add(expire_time, interval 8 hour ),1,0) as is_disable " \
            f"from auditsql_online_version " \
            f"where is_deleted='0' and created_at >= '{before_30_days}' order by created_at desc"
        data = []
        for row in OnlineVersion.objects.raw(query):
            data.append({'version': row.version, 'id': row.id, 'is_disable': row.is_disable})
        return Response(data=data, status=status.HTTP_200_OK)


class OnlineVersionListView(APIView):
    """创建或返回上线版本列表"""
    # 拥有工单提交/执行/审核权限的用户可以操作
    permission_classes = (anyof(CanCommitPermission, CanExecutePermission, CanAuditPermission),)

    def get(self, request):
        data = OnlineVersion.objects.filter(is_deleted='0'). \
            values('id', 'version', 'username', 'expire_time', 'updated_at', 'created_at').order_by('-created_at')
        return Response(data=data, status=status.HTTP_200_OK)

    def post(self, request):
        serializer = OnlineVersionListSerializer(data=request.data)
        if serializer.is_valid():
            s, data = serializer.op(request)
            code = 0 if s else 2
            return Response(data={'code': code, 'data': data}, status=status.HTTP_200_OK)
        else:
            errors = [str(v[0]) for k, v in serializer.errors.items()]
            data = {'code': 2, 'data': '\n'.join(errors)}
            return Response(data=data, status=status.HTTP_400_BAD_REQUEST)


class OnlineVersionDetailView(APIView):
    """获取指定上线版本的所有工单数据"""

    def post(self, request):
        serializer = OnlineVersionDetailSerializer(data=request.data)
        if serializer.is_valid():
            data = serializer.query()
            return Response(data=data, status=status.HTTP_200_OK)
        else:
            errors = [str(v[0]) for k, v in serializer.errors.items()]
            data = {'code': 2, 'data': '\n'.join(errors)}
            return Response(data=data, status=status.HTTP_400_BAD_REQUEST)
