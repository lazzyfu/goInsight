from django.contrib import admin

# Register your models here.

from django_celery_beat.models import PeriodicTask, IntervalSchedule, CrontabSchedule, SolarSchedule
from django_celery_results.models import TaskResult

admin.site.unregister(IntervalSchedule)
admin.site.unregister(CrontabSchedule)
admin.site.unregister(PeriodicTask)
admin.site.unregister(SolarSchedule)
admin.site.unregister(TaskResult)
