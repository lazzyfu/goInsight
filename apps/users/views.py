import io
import json

from PIL import Image
from django.contrib.auth import logout
from django.core.serializers.json import DjangoJSONEncoder
from django.http import HttpResponse, HttpResponseRedirect
from django.shortcuts import render
from django.urls import reverse_lazy, reverse
from django.views import View
from django.views.generic import FormView, RedirectView

from users import verifyCode
from users.forms import LoginForm, ChangePasswordForm, ChangeMobileForm
from users.models import UserAccounts


class LoginView(FormView):
    """用户登录， success_url登陆成功后跳转的页面"""
    template_name = 'login.html'
    form_class = LoginForm
    success_url = reverse_lazy('p_user_profile')

    def form_valid(self, form):
        result = form.authentication(self.request)
        if result['status'] is True:
            return super(LoginView, self).form_valid(form)
        else:
            return render(self.request, self.template_name, {'msg': result['msg']})

    def form_invalid(self, form):
        error = form.errors.as_text()
        return render(self.request, self.template_name, {'msg': error})


class LogoutView(RedirectView):
    """用户登出"""
    permanent = False
    url = reverse_lazy('p_login')

    def get(self, request, *args, **kwargs):
        logout(self.request)
        return super(LogoutView, self).get(request, *args, **kwargs)


class VerifyView(View):
    def get(self, request):
        """
        生成随机验证码
        """
        stream = io.BytesIO()
        img, code = verifyCode.create_validate_code()
        img.save(stream, 'png')
        request.session['verifycode'] = code
        return HttpResponse(stream.getvalue())


class IndexView(View):
    """访问首页，重定向的页面"""

    def get(self, request):
        return HttpResponseRedirect(reverse('p_user_profile'))


class UserProfileView(View):
    """用户profile"""

    def get(self, request):
        return render(request, 'profile.html')


class ChangePasswordView(View):
    """用户修改密码"""

    def post(self, request):
        form = ChangePasswordForm(request.POST)
        if form.is_valid():
            context = form.change_pass(request)
        else:
            error = form.errors.as_text()
            context = {'status': 2, 'msg': error}
        return HttpResponse(json.dumps(context))


class ChangeMobileView(View):
    """用户修改手机号"""

    def post(self, request):
        form = ChangeMobileForm(request.POST)
        if form.is_valid():
            context = form.change_mobile(request)
        else:
            error = form.errors.as_text()
            context = {'status': 2, 'msg': error}
        return HttpResponse(json.dumps(context))


class ChangePicView(View):
    """用户头像修改"""

    def get(self, request):
        return render(request, 'userpicture.html')

    def post(self, request):
        avatar_data = eval(request.POST.get('avatar_data'))

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

        result = {'state': 200}

        return HttpResponse(json.dumps(result))


class GetUserMailView(View):
    def get(self, request):
        queryset = UserAccounts.objects.all().values('username', 'email')
        serialize_data = json.dumps(list(queryset), cls=DjangoJSONEncoder)
        return HttpResponse(serialize_data)
