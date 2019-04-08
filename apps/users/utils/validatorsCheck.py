# -*- coding:utf-8 -*-
# edit by fuzongfei
import re

from rest_framework import serializers


def password_validate(value):
    """
    检查密码复杂度
    """
    pattern = re.compile('(?=.*[0-9])(?=.*[a-zA-Z])(?=.*[^a-zA-Z0-9]).{8,30}')
    if not pattern.findall(value):
        raise serializers.ValidationError('新密码复杂度不够(大小写字符,数字,符号组合)')


def audit_permission_validator(value):
    allowed_permissons = ['can_audit']
    if value not in allowed_permissons:
        raise serializers.ValidationError('权限校验不匹配')
