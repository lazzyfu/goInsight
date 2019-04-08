# -*- coding:utf-8 -*-
# edit by fuzongfei
import json
from datetime import datetime

from django.utils import timezone
from rest_framework import serializers

from orders.models import Orders
from orders.utils.msgNotice import MsgPush


class OrderApproveSerializer(serializers.Serializer):
    id = serializers.IntegerField(required=True, error_messages={'required': 'id字段不能为空'})
    status = serializers.CharField(max_length=10, required=True, error_messages={
        'required': '状态字段不能为空',
        'max_length': '最大长度不能超过10个字符'
    })
    msg = serializers.CharField(required=False)

    def op(self, request):
        sdata = self.validated_data
        id = sdata.get('id')
        status = sdata.get('status')
        msg = sdata.get('msg')

        data = Orders.objects.get(pk=id)
        auditor = json.loads(data.auditor)
        allowed_auditor = [x.get('user') for x in auditor]

        if status == u'通过':
            # 超级管理员可以一键审核通过
            if request.user.is_superuser:
                auditor.append({
                    'user': request.user.username,
                    'is_superuser': 1,
                    'status': 1,
                    'time': datetime.now().strftime("%Y-%m-%d %H:%M:%S"),
                    'msg': msg
                })
                data.progress = '2'
                data.auditor = json.dumps(auditor)
                data.save()
                # 推送消息
                msg_push = MsgPush(id=id, user=request.user.username, type='approve', msg=msg)
                msg_push.send()
                return True, '操作成功'

            if request.user.username in allowed_auditor:
                # 当用户在审核的人员列表中时，允许审核
                for i in auditor:
                    if request.user.username == i['user']:
                        if i['status'] == 1:
                            return False, '您已审核过，请不要重复执行'
                        else:
                            i['user'] = request.user.username
                            i['status'] = 1
                            i['time'] = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
                            i['msg'] = msg
                            data.auditor = json.dumps(auditor)
                            data.save()
                            # 当所有非超级管理员的审核人员批准后，修改工单状态为已完成
                            f_status = [x['status'] for x in auditor if x['is_superuser'] == 0]
                            if 1 in f_status and all(f_status):
                                data.progress = '2'
                                data.save()
                                # 推送消息
                                msg_push = MsgPush(id=id, user=request.user.username, type='approve', msg=msg)
                                msg_push.send()
                            return True, '操作成功'
            else:
                return False, '您不在审核人列表中，没有权限操作'
        if status == u'不通过':
            if request.user.is_superuser:
                auditor.append({
                    'user': request.user.username,
                    'is_superuser': 1,
                    'status': 1,
                    'time': datetime.now().strftime("%Y-%m-%d %H:%M:%S"),
                    'msg': msg
                })
                data.progress = '1'
                data.auditor = json.dumps(auditor)
                data.save()
                # 推送消息
                msg_push = MsgPush(id=id, user=request.user.username, type='approve', msg=msg)
                msg_push.send()
                return True, '操作成功'

            if request.user.username in allowed_auditor:
                for i in auditor:
                    if request.user.username == i['user']:
                        if i['status'] == 1:
                            return False, '您已审核过，请不要重复执行'
                        else:
                            i['user'] = request.user.username
                            i['status'] = 1
                            i['time'] = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
                            i['msg'] = msg
                            data.progress = '1'
                            data.auditor = json.dumps(auditor)
                            data.save()
                            # 推送消息
                            msg_push = MsgPush(id=id, user=request.user.username, type='approve', msg=msg)
                            msg_push.send()
                            return True, '操作成功'
            else:
                return False, '您不在审核人列表中，没有权限操作'


class OrderFeedbackSerializer(serializers.Serializer):
    id = serializers.IntegerField(required=True, error_messages={'required': 'id字段不能为空'})
    status = serializers.CharField(max_length=10, required=True, error_messages={
        'required': '状态字段不能为空',
        'max_length': '最大长度不能超过10个字符'
    })
    msg = serializers.CharField(required=False)

    def op(self, request):
        sdata = self.validated_data
        id = sdata.get('id')
        status = sdata.get('status')
        msg = sdata.get('msg')

        data = Orders.objects.get(pk=id)
        # 当用户点击的是处理中, 状态变为：处理中
        if status == u'处理中':
            data.progress = '3'
            data.updated_at = timezone.now()
            data.save()
            # 推送消息
            msg_push = MsgPush(id=id, user=request.user.username, type='feedback', msg=msg)
            msg_push.send()
            return True, '操作成功'

        # 当用户点击的是已完成, 状态变为：已完成
        elif status == u'已完成':
            data.progress = '4'
            data.updated_at = timezone.now()
            data.save()
            # 推送消息
            msg_push = MsgPush(id=id, user=request.user.username, type='feedback', msg=msg)
            msg_push.send()
            return True, '操作成功'


class OrderReviewSerializer(serializers.Serializer):
    id = serializers.IntegerField(required=True, error_messages={'required': 'id字段不能为空'})
    status = serializers.CharField(max_length=10, required=True, error_messages={
        'required': '状态字段不能为空',
        'max_length': '最大长度不能超过10个字符'
    })
    msg = serializers.CharField(required=False)

    def op(self, request):
        sdata = self.validated_data
        id = sdata.get('id')
        status = sdata.get('status')
        msg = sdata.get('msg')

        data = Orders.objects.get(pk=id)

        reviewer = json.loads(data.reviewer)
        allowed_reviewer = [x.get('user') for x in reviewer]

        # 只有选定的复核人可以复核，不在复核列表里面的人员均不可复核（超级管理员也不行）
        if status == u'已核对':
            if request.user.username in allowed_reviewer:
                for i in reviewer:
                    if request.user.username == i['user']:
                        if i['status'] == 1:
                            return False, '您已复核过，请不要重复执行'
                        else:
                            i['user'] = request.user.username
                            i['status'] = 1
                            i['time'] = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
                            i['msg'] = msg
                            data.reviewer = json.dumps(reviewer)
                            data.save()
                            # 检查所有的复核人是否都通过复核，通过将工单设置为已复核状态
                            f_status = [x['status'] for x in reviewer if x['is_superuser'] == 0]
                            if 1 in f_status and all(f_status):
                                data.progress = '6'
                                data.save()
                            # 推送消息
                            msg_push = MsgPush(id=id, user=request.user.username, type='review', msg=msg)
                            msg_push.send()
                            return True, '操作成功'
            else:
                return False, '您不在复核人列表中，没有权限操作'
        if status == u'关闭窗口':
            return True, '窗口关闭，操作结束'


class OrderCloseSerializer(serializers.Serializer):
    id = serializers.IntegerField(required=True, error_messages={'required': 'id字段不能为空'})
    status = serializers.CharField(max_length=10, required=True, error_messages={
        'required': '状态字段不能为空',
        'max_length': '最大长度不能超过10个字符'
    })
    msg = serializers.CharField(required=False, allow_blank=True)

    def op(self, request):
        sdata = self.validated_data
        id = sdata.get('id')
        status = sdata.get('status')
        msg = sdata.get('msg')

        data = Orders.objects.get(pk=id)
        close_info = {}

        if data.progress == '5':
            return False, '该记录已被关闭、请不要重复执行'
        if status == u'提交':
            if len(msg) < 5:
                return False, '<关闭原因>请至少输入5个字符'
            if data.progress in ['3', '4', '6', '7']:
                return False, '工单正在处理中或已完成，无法关闭'
            else:
                data.progress = '5'
                close_info['user'] = request.user.username
                close_info['msg'] = msg
                close_info['time'] = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
                data.close_info = json.dumps(close_info)
                data.save()
                # 推送消息
                msg_push = MsgPush(id=id, user=request.user.username, type='close', msg=msg)
                msg_push.send()
                return True, '操作成功'
        if status == u'关闭窗口':
            return True, '窗口关闭，操作结束'
