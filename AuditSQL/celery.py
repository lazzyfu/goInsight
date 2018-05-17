# -*- coding:utf-8 -*-
# edit by fuzongfei

from __future__ import absolute_import, unicode_literals
import os
from celery import Celery
from AuditSQL import settings


class MyCelery(Celery):
    def now(self):
        """Return the current time and date as a datetime."""
        from datetime import datetime
        import pytz
        # 由于celery-4.1.0存在时区bug， 此处写死时区，不受系统时区影响
        tz = pytz.timezone('Asia/Shanghai')
        return datetime.now(tz)


# set the default Django settings module for the 'celery' program.
os.environ.setdefault('DJANGO_SETTINGS_MODULE', 'AuditSQL.settings')

app = MyCelery('AuditSQL')

# Using a string here means the worker don't have to serialize
# the configuration object to child processes.
# - namespace='CELERY' means all celery-related configuration keys
#   should have a `CELERY_` prefix.
app.config_from_object('django.conf:settings', namespace='CELERY')

# Load task modules from all registered Django app configs.
app.autodiscover_tasks(lambda: settings.INSTALLED_APPS)


@app.task(bind=True)
def debug_task(self):
    print('Request: {0!r}'.format(self.request))
