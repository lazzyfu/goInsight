from django.contrib import admin

# Register your models here.

from djcelery.models import (
    TaskState, WorkerState, PeriodicTask,
    IntervalSchedule, CrontabSchedule)

admin.site.unregister(TaskState)
admin.site.unregister(WorkerState)
admin.site.unregister(IntervalSchedule)
admin.site.unregister(CrontabSchedule)
admin.site.unregister(PeriodicTask)

