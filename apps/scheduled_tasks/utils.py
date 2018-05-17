# -*- coding:utf-8 -*-
# edit by fuzongfei
# from djcelery.models import PeriodicTask, PeriodicTasks


from django_celery_beat.models import PeriodicTask, PeriodicTasks


# 当定时任务变更后，刷新定时任务
def refresh_periodic_tasks():
    dummy_periodic_task = PeriodicTask()
    dummy_periodic_task.no_changes = False
    PeriodicTasks.changed(dummy_periodic_task)
