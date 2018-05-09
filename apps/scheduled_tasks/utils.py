# -*- coding:utf-8 -*-
# edit by fuzongfei
from djcelery.models import PeriodicTask, PeriodicTasks


# 更新任务状态
def update_periodic_tasks():
    dummy_periodic_task = PeriodicTask()
    dummy_periodic_task.no_changes = False
    PeriodicTasks.changed(dummy_periodic_task)

