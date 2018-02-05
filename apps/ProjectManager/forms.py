# -*- coding:utf-8 -*-
# edit by fuzongfei
from ast import literal_eval
from datetime import datetime

from django import forms

from ProjectManager.models import InceptionSqlOperateRecord


class InceptionSqlOperateForm(forms.Form):
    host = forms.CharField(required=True)
    database = forms.CharField(required=True, max_length=64)
    op_action = forms.CharField(required=True)
    op_type = forms.CharField(required=True)
    sql_content = forms.CharField(widget=forms.Textarea)

    def is_save(self, request, result):
        cleaned_data = super(InceptionSqlOperateForm, self).clean()
        host = cleaned_data.get('host')
        database = cleaned_data.get('database')
        op_user = request.user.username
        op_uid = request.user.uid
        workid = datetime.now().strftime("%Y%m%d%H%M%S%f")

        for line in result:
            op_sql = line['SQL']
            step_id = line['ID']
            stage = line['stage']
            stagestatus = line['stagestatus']
            errlevel = line['errlevel']
            errormessage = line['errormessage']
            Affected_rows = line['Affected_rows']
            sequence = line['sequence']
            backup_dbname = line['backup_dbname']
            execute_time = line['execute_time']

            InceptionSqlOperateRecord.objects.create(
                op_user=op_user,
                op_uid=op_uid,
                workid=workid,
                dst_host=host,
                dst_database=database,
                stagestatus=stagestatus,
                op_sql=op_sql,
                step_id=step_id,
                stage=stage,
                errlevel=errlevel,
                errormessage=errormessage,
                affected_rows=Affected_rows,
                sequence=literal_eval(sequence),
                backup_dbname=backup_dbname,
                execute_time=execute_time
            )
        context = {'data': result, 'errMsg': '执行完成', 'errCode': 200}
        return context
