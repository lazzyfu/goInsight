# Create your views here.
import io

from PIL import Image
from django.contrib.auth import logout
from django.http import HttpResponseRedirect, HttpResponse
from django.views import View
from django.views.generic import RedirectView
from rest_framework import status
from rest_framework.permissions import IsAuthenticated
from rest_framework.renderers import TemplateHTMLRenderer
from rest_framework.response import Response
from rest_framework.reverse import reverse_lazy
from rest_framework.views import APIView

from users.models import UserAccounts
from users.serializers import LoginSerializer, ChangePasswordSerializer, ChangeMobileSerializer, GetAuditorSerializer
from users.utils import verifyCode


class LoginView(APIView):
    """用户登录"""

    renderer_classes = [TemplateHTMLRenderer]
    template_name = 'login.html'
    # 此处解析url使用rest_framework.reverse.reverse_lazy方法
    success_url = reverse_lazy('p_profile')

    def get(self, request):
        return Response()

    def post(self, request):
        serializer = LoginSerializer(data=request.data)

        if serializer.is_valid():
            s, msg = serializer.authentication(request)
            if s:
                # 重定向到用户详情页面
                return HttpResponseRedirect(self.success_url)
            else:
                data = {'code': 2, 'data': msg}
                return Response(data=data, status=status.HTTP_400_BAD_REQUEST)
        else:
            errors = [str(v[0]) for k, v in serializer.errors.items()]
            data = {'code': 2, 'data': '\n'.join(errors)}
            return Response(data=data, status=status.HTTP_400_BAD_REQUEST)


class LogoutView(RedirectView):
    """用户登出"""

    permanent = False
    url = reverse_lazy('p_login')

    def get(self, request, *args, **kwargs):
        logout(self.request)
        return super(LogoutView, self).get(request, *args, **kwargs)


class VerifyCodeView(View):
    """验证码"""

    def get(self, request):
        stream = io.BytesIO()
        img, code = verifyCode.create_validate_code()
        img.save(stream, 'png')
        request.session['verifycode'] = code
        return HttpResponse(stream.getvalue())


class UserProfileView(APIView):
    """用户profile"""

    renderer_classes = [TemplateHTMLRenderer]
    template_name = 'profile.html'
    permission_classes = (IsAuthenticated,)

    def get(self, request):
        return Response()


class ChangePasswordView(APIView):
    """修改密码"""

    def post(self, request):
        serializer = ChangePasswordSerializer(data=request.data)
        if serializer.is_valid():
            s, msg = serializer.change(request)
            code = 0 if s else 2
            data = {'code': code, 'data': msg}
            return Response(data=data, status=status.HTTP_200_OK)
        else:
            errors = [str(v[0]) for k, v in serializer.errors.items()]
            data = {'code': 2, 'data': '\n'.join(errors)}
            return Response(data=data, status=status.HTTP_400_BAD_REQUEST)


class ChangeMobileView(APIView):
    """修改手机号"""

    def post(self, request):
        serializer = ChangeMobileSerializer(data=request.data)
        if serializer.is_valid():
            s, msg = serializer.change(request)
            if s:
                data = {'code': 0, 'data': msg}
                return Response(data=data, status=status.HTTP_200_OK)
        else:
            errors = [str(v[0]) for k, v in serializer.errors.items()]
            data = {'code': 2, 'data': '\n'.join(errors)}
            return Response(data=data, status=status.HTTP_400_BAD_REQUEST)


class ChangePicView(APIView):
    """用户头像修改"""

    def post(self, request):
        avatar_data = eval(request.data.get('avatar_data'))

        # 保存图片到upload_to位置，并将路径写入到字段avatar_file
        photo = request.FILES.get('avatar_file')
        photo_instance = UserAccounts.objects.get(uid=request.user.uid)
        photo_instance.avatar_file = photo
        photo_instance.save()

        # 获取截取图片的坐标
        x = avatar_data['x']
        y = avatar_data['y']
        w = avatar_data['width']
        h = avatar_data['height']

        # 裁剪图片
        # photo_instance.avatar_file：获取上面存储到数据库中的原始的图片（绝对路径）
        # photo_instance.avatar_file.path：获取原始图片的存储位置
        img = Image.open(photo_instance.avatar_file)
        # 按照前端传递来的坐标进行裁剪
        cropped_image = img.crop((x, y, w + x, h + y))
        # 对裁剪后的图片进行尺寸重新格式化
        resized_image = cropped_image.resize((305, 304), Image.ANTIALIAS)
        # 将裁剪后的图片替换掉原始图片，生成新的图片
        resized_image.save(photo_instance.avatar_file.path, 'PNG')

        return Response(data={'state': 200}, status=status.HTTP_200_OK)


class GetEmailCcView(APIView):
    """获取抄送的用户和邮箱"""

    def get(self, request):
        queryset = UserAccounts.objects.all().values('username', 'email')
        return Response(queryset)


class GetAuditorView(APIView):
    """获取有审核权限的用户"""

    def post(self, request):
        serializer = GetAuditorSerializer(data=request.data)
        if serializer.is_valid():
            data = serializer.query()
            return Response(data=data)
        else:
            errors = [str(v[0]) for k, v in serializer.errors.items()]
            data = {'code': 2, 'data': '\n'.join(errors)}
            return Response(data=data, status=status.HTTP_400_BAD_REQUEST)
