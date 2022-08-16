# -*- coding:utf-8 -*-
# edit by xff
import re

from rest_framework import serializers


def jwt_get_user_secret(user):
    return user.user_secret


def validate_password(value):
    """
    检查密码复杂度
    """
    pattern = re.compile('(?=.*[0-9])(?=.*[a-zA-Z])(?=.*[^a-zA-Z0-9]).{8,30}')
    if not pattern.findall(value):
        raise serializers.ValidationError('新密码复杂度不够(大小写字符,数字,符号组合)')