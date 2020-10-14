# -*- coding:utf-8 -*-
# edit by fuzongfei

"""
http返回封装
"""
from rest_framework.response import Response


class JsonResponseV1(Response):
    """
    封装by fuzongfei
    """

    def __init__(self, data={}, code='0000', message='', status=None, flat=False):
        """
        code = '0000' 为success
        code = '0001' 为error
        status = status=status.HTTP_200_OK ... 返回的http状态
        flat = True， 不返还KEY {'id': [ErrorDetail(string='审核未通过或工单已关闭，操作失败', code='invalid')]} 即不返回id
        """
        super().__init__(None, status=status)
        if isinstance(message, dict):
            if flat:
                fmt_msg = ' '.join(
                    [': '.join([str(v[0])]) for _, v in message.items()]
                )
            else:
                fmt_msg = ' '.join(
                    [': '.join([k, str(v[0])]) for k, v in message.items()]
                ) if isinstance(message, dict) else message
        else:
            fmt_msg = message
        self.data = {
            'data': data,
            'code': code,
            'message': fmt_msg
        }
