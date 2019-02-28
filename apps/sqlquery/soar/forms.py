# -*- coding:utf-8 -*-
# edit by fuzongfei

from django import forms

from sqlorders.models import MysqlConfig
from sqlquery.soar.soarApi import SoarAnalyze


class SoarAnalyzeForm(forms.Form):
    type = forms.ChoiceField(required=True, choices=(('advisor', 'advisor'),),
                             label=u'soar操作类型')
    contents = forms.CharField(widget=forms.Textarea, required=True, label=u'操作的内容')
    database = forms.CharField(min_length=1, max_length=256, required=True, strip=True, label=u'库名')

    def analyze(self):
        data = self.cleaned_data
        type = data.get('type')
        host, port, schema = data.get('database').split(',')
        contents = data.get('contents')

        obj = MysqlConfig.objects.get(host=host, port=port)
        analyze = SoarAnalyze(obj.user, obj.password, obj.host, obj.port, schema, type, contents)
        return analyze.run()