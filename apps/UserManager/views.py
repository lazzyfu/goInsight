from django.contrib.auth import logout, login
from django.urls import reverse, reverse_lazy
from django.http import HttpResponseRedirect
from django.shortcuts import render
from django.views import View
from django.views.generic import FormView, RedirectView

from UserManager.models import UserAccount
from .forms import LoginForm


# Create your views here.

class LoginView(FormView):
    """用户登录视图"""
    template_name = 'login.html'
    form_class = LoginForm
    success_url = reverse_lazy('p_project')

    def form_valid(self, form):
        user = form.is_verify()
        if user is not None:
            login(self.request, user)
            # 将用户所属的组id写入到session中
            groups = UserAccount.objects.get(uid=self.request.user.uid).groupsdetail_set.all().values_list(
                'group__group_id', flat=True)
            self.request.session['groups'] = list(groups)
            return super(LoginView, self).form_valid(form)
        else:
            return render(self.request, self.template_name, {'msg': '用户名或密码错误'})


class LogoutView(RedirectView):
    """用户登出视图"""
    permanent = False
    url = reverse_lazy('p_login')

    def get(self, request, *args, **kwargs):
        logout(self.request)
        return super(LogoutView, self).get(request, *args, **kwargs)


class IndexView(View):
    """用户登录后，重定向的页面"""

    def get(self, request):
        return HttpResponseRedirect(reverse('p_project'))
