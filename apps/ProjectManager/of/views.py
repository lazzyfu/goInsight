# -*- coding:utf-8 -*-
# edit by fuzongfei
import json
from datetime import datetime

from django.http import HttpResponse
from django.shortcuts import render
from django.utils.decorators import method_decorator
from django.views import View

from ProjectManager.forms import InceptionSqlCheckForm
from ProjectManager.inception.inception_api import sql_filter, IncepSqlCheck
from ProjectManager.models import IncepMakeExecTask
from ProjectManager.utils import check_incep_alive


class IncepOfAuditView(View):
    def get(self, request):
        return render(request, 'incep_perform_audit.html')

    @method_decorator(check_incep_alive)
    def post(self, request):
        form = InceptionSqlCheckForm(request.POST)
        if form.is_valid():
            cleaned_data = form.cleaned_data
            host = cleaned_data['host']
            database = cleaned_data['database']
            op_action = cleaned_data.get('op_action')
            op_type = cleaned_data['op_type']
            group_id = cleaned_data['group_id']
            sql_content = cleaned_data['sql_content']

            # 对检测的SQL类型进行区分
            filter_result = sql_filter(sql_content, op_action)

            # 实例化
            incep_of_audit = IncepSqlCheck(sql_content, host, database, self.request.user.username)

            if filter_result['errCode'] == 400:
                context = filter_result
            else:
                # SQL语法检查
                if op_type == 'check':
                    context = incep_of_audit.run_check()

                # 生成执行任务
                elif op_type == 'make':
                    # 生成执行任务之前，检测是否审核通过
                    check_result = incep_of_audit.is_check_pass()
                    if check_result['errCode'] == 400:
                        context = check_result
                    else:
                        # 对OSC执行的SQL生成sqlsha1
                        result = incep_of_audit.make_sqlsha1()
                        taskid = datetime.now().strftime("%Y%m%d%H%M%S%f")
                        # 生成执行任务记录
                        for row in result:
                            IncepMakeExecTask.objects.create(
                                uid=self.request.user.uid,
                                user=self.request.user.username,
                                taskid=taskid,
                                dst_host=host,
                                group_id=group_id,
                                dst_database=database,
                                sql_content=row['SQL'],
                                sqlsha1=row['sqlsha1'],
                                type=filter_result['type']
                            )
                        context = {'errCode': 201,
                                   'dst_url': f'/projects/incep_perform_records/incep_perform_details/{taskid}'}
            return HttpResponse(json.dumps(context))
        else:
            error = "请选择主机、库名和项目组"
            context = {'errCode': 400, 'errMsg': error}

            return HttpResponse(json.dumps(context))
