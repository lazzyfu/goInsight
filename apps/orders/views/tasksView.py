# -*- coding:utf-8 -*-
# edit by fuzongfei
from rest_framework import status
from rest_framework.renderers import TemplateHTMLRenderer
from rest_framework.response import Response
from rest_framework.views import APIView

from orders.permissions import CanExecutePermission, anyof, CanCommitPermission, CanAuditPermission
from orders.serializers.taskSerializers import GenerateSubtasksSerializer, SubtasksDetailSerializer, \
    FullExecuteSerializer, SingleExecuteSerializer, GetTasksLogSerializer, StopExecuteSerializer


class GenerateSubtasksView(APIView):
    """将SQL工单拆解成子任务"""
    # 有执行权限的用户可以执行
    permission_classes = (CanExecutePermission,)

    def post(self, request):
        serializer = GenerateSubtasksSerializer(data=request.data)

        if serializer.is_valid():
            s, data = serializer.save(request)
            code = 0 if s else 2
            return Response(data={'code': code, 'data': data}, status=status.HTTP_200_OK)
        else:
            errors = [str(v[0]) for k, v in serializer.errors.items()]
            data = {'code': 2, 'data': '\n'.join(errors)}
            return Response(data=data, status=status.HTTP_400_BAD_REQUEST)


class RenderSubtasksView(APIView):
    """有执行权限的用户可以访问"""
    permission_classes = (anyof(CanCommitPermission, CanExecutePermission, CanAuditPermission),)

    renderer_classes = [TemplateHTMLRenderer]
    template_name = 'orders/subtasks.html'

    def get(self, request, taskid):
        return Response(data={'taskid': taskid}, status=status.HTTP_200_OK)


class SubTasksDetailView(APIView):
    """
    获取子任务数据
    有执行权限的用户可以访问
    """

    # permission_classes = (CanExecutePermission,)

    def post(self, request):
        serializer = SubtasksDetailSerializer(data=request.data)

        if serializer.is_valid():
            data = serializer.query()
            return Response(data=data, status=status.HTTP_200_OK)
        else:
            errors = [str(v[0]) for k, v in serializer.errors.items()]
            data = {'code': 2, 'data': '\n'.join(errors)}
            return Response(data=data, status=status.HTTP_400_BAD_REQUEST)


class FullExecuteView(APIView):
    """工单一键全部执行"""
    # 有执行权限的用户可以访问
    permission_classes = (CanExecutePermission,)

    def post(self, request):
        serializer = FullExecuteSerializer(data=request.data)

        if serializer.is_valid():
            s, data = serializer.execute(request)
            code = 0 if s else 2
            return Response(data={'code': code, 'data': data}, status=status.HTTP_200_OK)
        else:
            errors = [str(v[0]) for k, v in serializer.errors.items()]
            data = {'code': 2, 'data': '\n'.join(errors)}
            return Response(data=data, status=status.HTTP_400_BAD_REQUEST)


class SingleExecuteView(APIView):
    """工单单条执行"""
    # 有执行权限的用户可以访问
    permission_classes = (CanExecutePermission,)

    def post(self, request):
        serializer = SingleExecuteSerializer(data=request.data)

        if serializer.is_valid():
            s, data = serializer.execute(request)
            code = 0 if s else 2
            return Response(data={'code': code, 'data': data}, status=status.HTTP_200_OK)
        else:
            errors = [str(v[0]) for k, v in serializer.errors.items()]
            data = {'code': 2, 'data': '\n'.join(errors)}
            return Response(data=data, status=status.HTTP_400_BAD_REQUEST)


class StopExecuteView(APIView):
    """
    执行任务-操作，支持：暂停、恢复、终止
    只支持停止修改表结构的操作
    """
    # 有执行权限的用户可以访问
    permission_classes = (CanExecutePermission,)

    def post(self, request):
        serializer = StopExecuteSerializer(data=request.data)

        if serializer.is_valid():
            s, data = serializer.op()
            code = 0 if s else 2
            return Response(data={'code': code, 'data': data}, status=status.HTTP_200_OK)
        else:
            errors = [str(v[0]) for k, v in serializer.errors.items()]
            data = {'code': 2, 'data': '\n'.join(errors)}
            return Response(data=data, status=status.HTTP_400_BAD_REQUEST)


class GetTasksLogView(APIView):
    """获取任务的执行日志和回滚SQL"""

    def post(self, request):
        serializer = GetTasksLogSerializer(data=request.data)

        if serializer.is_valid():
            s, data = serializer.query()
            code = 0 if s else 2
            return Response(data={'code': code, 'data': data}, status=status.HTTP_200_OK)
        else:
            errors = [str(v[0]) for k, v in serializer.errors.items()]
            data = {'code': 2, 'data': '\n'.join(errors)}
            return Response(data=data, status=status.HTTP_400_BAD_REQUEST)
