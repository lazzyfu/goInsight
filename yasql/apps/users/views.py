# Create your views here.
from uuid import uuid4

from django.utils.decorators import method_decorator
from django.views.decorators.csrf import csrf_exempt, ensure_csrf_cookie
from rest_framework import status
from rest_framework.generics import ListAPIView, UpdateAPIView
from rest_framework.permissions import AllowAny
from rest_framework.views import APIView

from libs.response import JsonResponseV1
from users import serializers, models


class Login(APIView):
    """登录"""
    permission_classes = [AllowAny]
    serializer_class = serializers.LoginSerializer

    def perform_authentication(self, request):
        # 不检查JWT，兼容/admin
        pass

    @method_decorator(csrf_exempt)
    def dispatch(self, *args, **kwargs):
        return super().dispatch(*args, **kwargs)

    @method_decorator(ensure_csrf_cookie)
    def post(self, request):
        serializer = serializers.LoginSerializer(data=request.data)

        if serializer.is_valid():
            code, data = serializer.login()
            if code:
                return JsonResponseV1(data={'token': data})
            return JsonResponseV1(message=data, code='0001')
        return JsonResponseV1(serializer.errors, code='0001')


class Logout(APIView):
    """退出登录"""

    def post(self, request, *args, **kwargs):
        user = request.user
        # 每次退出登录时，修改user_secret的值，即可实现token失效
        user.user_secret = uuid4()
        user.save()
        return JsonResponseV1(status=status.HTTP_200_OK)


class UserInfo(APIView):
    """获取登录用户信息"""

    def get(self, request):
        user = serializers.UserInfoSerializer(request.user)
        return JsonResponseV1(data=user.data)


class GetUsers(ListAPIView):
    queryset = models.UserAccounts.objects.all()
    serializer_class = serializers.UsersListSerializer

    # 获取审核、复核、抄送的用户列表
    def get(self, request, *args, **kwargs):
        queryset = self.get_queryset()
        serializer = self.get_serializer(queryset, many=True)
        return JsonResponseV1(data=serializer.data)


class UpdateUserInfo(UpdateAPIView):
    """更新用户信息"""

    def put(self, request, *args, **kwargs):
        serializer = serializers.UpdateUserInfoSerializer(
            instance=models.UserAccounts.objects.get(username=kwargs['username']),
            data=request.data
        )

        if serializer.is_valid():
            serializer.save()
            return JsonResponseV1(message="更新成功")
        return JsonResponseV1(code='0001', message=serializer.errors)


class ChangePassword(APIView):
    """修改密码"""

    def post(self, request):
        serializer = serializers.ChangePasswordSerializer(data=request.data)
        if serializer.is_valid():
            code, msg = serializer.change(request)
            if code:
                return JsonResponseV1(message=msg)
            return JsonResponseV1(message=msg, code='0001')
        return JsonResponseV1(message=serializer.errors, code='0001', flat=True)


class ChangeAvatar(APIView):
    """修改头像"""

    def post(self, request):
        photo_instance = models.UserAccounts.objects.get(uid=request.user.uid)
        fileobj = request.FILES.get('avatar_file')
        photo_instance.avatar_file = fileobj
        photo_instance.save()
        return JsonResponseV1(message="上传成功")
