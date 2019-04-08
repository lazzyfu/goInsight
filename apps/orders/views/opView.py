# -*- coding:utf-8 -*-
# edit by fuzongfei
from rest_framework import status
from rest_framework.response import Response
from rest_framework.views import APIView

from orders.permissions import CanExecutePermission, CanCommitPermission, CanAuditPermission, anyof
from orders.serializers.opSerializers import OrderApproveSerializer, OrderFeedbackSerializer, OrderReviewSerializer, \
    OrderCloseSerializer


class OrderApproveView(APIView):
    """工单审批"""

    # 拥有审核权限的用户可以操作
    permission_classes = (CanAuditPermission,)

    def post(self, request):
        serializer = OrderApproveSerializer(data=request.data)
        if serializer.is_valid():
            s, data = serializer.op(request)
            code = 0 if s else 2
            return Response(data={'code': code, 'data': data}, status=status.HTTP_200_OK)
        else:
            errors = [str(v[0]) for k, v in serializer.errors.items()]
            data = {'code': 2, 'data': '\n'.join(errors)}
            return Response(data=data, status=status.HTTP_400_BAD_REQUEST)


class OrderFeedbackView(APIView):
    """
    工单的反馈
    当工单不需要实际执行或运维工单或执行内容包含了失败的情况，工单不会自动变为已完成状态，此时需要人工点击反馈进行变更工单的状态
    """

    # 拥有执行权限的用户可以操作
    permission_classes = (CanExecutePermission,)

    def post(self, request):
        serializer = OrderFeedbackSerializer(data=request.data)
        if serializer.is_valid():
            s, data = serializer.op(request)
            code = 0 if s else 2
            return Response(data={'code': code, 'data': data}, status=status.HTTP_200_OK)
        else:
            errors = [str(v[0]) for k, v in serializer.errors.items()]
            data = {'code': 2, 'data': '\n'.join(errors)}
            return Response(data=data, status=status.HTTP_400_BAD_REQUEST)


class OrderReviewView(APIView):
    """
    工单的复核
    工单完成后，需要相关人员进行核对，核对完成后，需要点击复核按钮更新工单状态为已核对
    """
    # 拥有工单提交/执行/审核权限的用户可以操作
    permission_classes = (anyof(CanCommitPermission, CanExecutePermission, CanAuditPermission),)

    def post(self, request):
        serializer = OrderReviewSerializer(data=request.data)
        if serializer.is_valid():
            s, data = serializer.op(request)
            code = 0 if s else 2
            return Response(data={'code': code, 'data': data}, status=status.HTTP_200_OK)
        else:
            errors = [str(v[0]) for k, v in serializer.errors.items()]
            data = {'code': 2, 'data': '\n'.join(errors)}
            return Response(data=data, status=status.HTTP_400_BAD_REQUEST)


class OrderCloseView(APIView):
    # 拥有工单提交/执行/审核权限的用户可以操作
    permission_classes = (anyof(CanCommitPermission, CanExecutePermission, CanAuditPermission),)

    def post(self, request):
        serializer = OrderCloseSerializer(data=request.data)
        if serializer.is_valid():
            s, data = serializer.op(request)
            code = 0 if s else 2
            return Response(data={'code': code, 'data': data}, status=status.HTTP_200_OK)
        else:
            print(serializer.errors.items())
            errors = [str(v[0]) for k, v in serializer.errors.items()]
            data = {'code': 2, 'data': '\n'.join(errors)}
            return Response(data=data, status=status.HTTP_400_BAD_REQUEST)
