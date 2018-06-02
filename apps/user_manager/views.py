import json

from PIL import Image
from django.contrib.auth import logout, login, authenticate
from django.contrib.auth.hashers import make_password
from django.db.models import F
from django.urls import reverse, reverse_lazy
from django.http import HttpResponseRedirect, HttpResponse
from django.shortcuts import render
from django.views import View
from django.views.generic import FormView, RedirectView

from user_manager.models import UserAccount, PermissionDetail, RolesDetail
from .forms import ChangePasswordForm, LoginForm


# Create your views here.

class LoginView(FormView):
    """用户登录视图， success_url登陆成功后访问的页面"""
    template_name = 'login.html'
    form_class = LoginForm
    success_url = reverse_lazy('p_ol_records')

    def form_valid(self, form):
        cleaned_data = form.cleaned_data
        username = cleaned_data.get('username')
        password = cleaned_data.get('password')

        try:
            user = authenticate(username=username, password=password)

            obj = UserAccount.objects.get(username=username)
            if make_password(password) != obj.password:
                result = {'msg': '用户名或密码错误'}

            if not obj.is_active:
                result = {'msg': '用户未激活，请联系管理员'}
            else:
                if user is not None:
                    login(self.request, user)
                    # 将用户所属的组id写入到session中
                    groups = list(UserAccount.objects.get(uid=user.uid).groupsdetail_set.all().values_list(
                        'group__group_id', flat=True))
                    self.request.session['groups'] = groups
                    try:
                        user_role = self.request.user.user_role()
                        # 将用户权限写入到session
                        perm_list = list(PermissionDetail.objects.annotate(
                            permission_name=F('permission__permission_name')).filter(role__role_name=user_role).values_list(
                            'permission_name', flat=True))
                        self.request.session['perm_list'] = perm_list
                        # 判断用户是否属于任何一个组
                        if len(list(groups)) > 0:
                            return super(LoginView, self).form_valid(form)
                        else:
                            result = {'msg': '用户未被分配项目组，请联系管理员'}
                    except RolesDetail.DoesNotExist:
                        result = {'msg': '用户未被分配角色，请联系管理员'}
        except UserAccount.DoesNotExist:
            result = {'msg': '用户不存在，请联系管理员'}
        return render(self.request, self.template_name, result)


class LogoutView(RedirectView):
    """用户登出视图"""
    permanent = False
    url = reverse_lazy('p_login')

    def get(self, request, *args, **kwargs):
        logout(self.request)
        return super(LogoutView, self).get(request, *args, **kwargs)


class IndexView(View):
    """访问首页，重定向的页面"""

    def get(self, request):
        return HttpResponseRedirect(reverse('p_ol_records'))


class UserProfileView(View):
    def get(self, request):
        return render(request, 'profile.html')


class ChangePasswordView(View):
    def post(self, request):
        form = ChangePasswordForm(request.POST)
        if form.is_valid():
            cleaned_data = form.cleaned_data
            old_password = cleaned_data['old_password']
            new_password = cleaned_data['new_password']
            verify_password = cleaned_data['verify_password']

            user = UserAccount.objects.get(uid=request.user.uid)
            if new_password == verify_password:
                if user.check_password(old_password):
                    if old_password != new_password:
                        user.password = make_password(new_password)
                        user.save()
                        context = {'status': 0, 'msg': '密码修改成功'}
                    else:
                        context = {'status': 2, 'msg': '新密码等于旧密码，请重新输入'}
                else:
                    context = {'status': 2, 'msg': '旧密码错误，请重新输入'}
            else:
                context = {'status': 2, 'msg': '密码不匹配，请重新输入'}
        else:
            error = form.errors.as_text()
            context = {'status': 2, 'msg': error}
        return HttpResponse(json.dumps(context))


class ChangePicView(View):
    def get(self, request):
        return render(request, 'userpic.html')

    def post(self, request):
        avatar_data = eval(request.POST.get('avatar_data'))

        # 保存图片到upload_to位置，并将路径写入到字段avatar_file
        photo = request.FILES.get('avatar_file')
        photo_instance = UserAccount.objects.get(uid=request.user.uid)
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
