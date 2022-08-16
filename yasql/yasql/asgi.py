# -*- coding:utf-8 -*-
# edit by xff

"""
ASGI entrypoint. Configures Django and then runs the application
defined in the ASGI_APPLICATION setting.
"""
import os
import django
from channels.routing import get_default_application

os.environ.setdefault("DJANGO_SETTINGS_MODULE", "yasql.settings")

django.setup()

application = get_default_application()
