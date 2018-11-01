from django.http import JsonResponse
from django.shortcuts import render

# Create your views here.
from django.views import View

from webshell.models import WebShellInfo


class WebSSHView(View):
    def get(self, request):
        return render(request, 'webshell/webssh.html')


class GetWebSSHCmdView(View):
    def get(self, request):
        query = f"select a.id, a.command, a.comment from sqlaudit_web_shell a " \
                f"join sqlaudit_web_shell_grant b on a.id = b.shell_id " \
                f"join sqlaudit_user_accounts c on b.user_id = c.uid where c.uid = {request.user.uid}"
        data = []
        for row in WebShellInfo.objects.raw(query):
            data.append({
                'command': row.command,
                'comment': row.comment
            })
        return JsonResponse(list(data), safe=False)
