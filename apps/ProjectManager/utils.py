# -*- coding:utf-8 -*-
# edit by fuzongfei
import json

from ProjectManager.models import IncepMakeExecTask


def update_tasks_status(**kwargs):
    """
    更新任务进度
    更新备份信息
    """
    id = kwargs.get('id')
    taskid = kwargs.get('taskid')
    exec_result = kwargs.get('exec_result')

    data = IncepMakeExecTask.objects.get(id=id, taskid=taskid)
    errlevel = [x['errlevel'] for x in exec_result]
    if 1 in errlevel or 2 in errlevel:
        # 报错，设置为：未执行
        data.exec_status = 0
        data.save()
    else:
        # 未报错
        # 更新进度为：已完成
        data.exec_status = 1

        # 更新备份信息
        data.sequence = exec_result[1]['sequence']
        data.backup_dbname = exec_result[1]['backup_dbname']
        data.exec_log = exec_result
        data.save()