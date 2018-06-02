# -*- coding:utf-8 -*-
# edit by fuzongfei
from datetime import datetime

from django import forms

from project_manager.inception.inception_api import IncepSqlCheck
from project_manager.models import IncepMakeExecTask


class OflineAuditForm(forms.Form):
    host = forms.CharField(required=True)
    database = forms.CharField(required=True, max_length=64)
    operate_type = forms.CharField(required=True)
    group_id = forms.IntegerField(required=True)
    contents = forms.CharField(widget=forms.Textarea)

    def save(self, request):
        cleaned_data = super(OflineAuditForm, self).clean()

        host = cleaned_data['host']
        database = cleaned_data['database']
        operate_type = cleaned_data.get('operate_type')
        group_id = cleaned_data['group_id']
        sql_content = cleaned_data['contents']

        # 实例化
        of_audit = IncepSqlCheck(sql_content, host, database, request.user.username)

        # 生成执行任务
        check_result = of_audit.is_check_pass()
        if check_result['status'] == 2:
            context = check_result
        else:
            # 对OSC执行的SQL生成sqlsha1
            result = of_audit.make_sqlsha1()
            taskid = datetime.now().strftime("%Y%m%d%H%M%S%f")
            # 生成执行任务记录
            for row in result:
                IncepMakeExecTask.objects.create(
                    uid=request.user.uid,
                    user=request.user.username,
                    taskid=taskid,
                    dst_host=host,
                    group_id=group_id,
                    dst_database=database,
                    sql_content=row['SQL'],
                    sqlsha1=row['sqlsha1'],
                    affected_row=row['Affected_rows'],
                    type=operate_type
                )
            context = {'status': 0, 'msg': '',
                       'jump_url': f'/projects/pt/perform_records/perform_details/{taskid}'}
        return context
