# -*- coding:utf-8 -*-
# edit by fuzongfei
from django.db.models import F

from user_manager.models import UserAccount, Contacts, ContactsDetail


class GetEmailAddr(object):
    def __init__(self, obj, latest_id):
        self.obj = obj
        self.latest_id = latest_id

    def get_user_email(self, *args):
        """
        args
        传入用户参数列表，可传入：proposer, verifier, operator
        返回传入用户的邮箱地址
        """
        obj = self.obj.objects.get(pk=self.latest_id)
        user_list = []
        if 'proposer' in args:
            user_list.append(obj.proposer)
        if 'verifier' in args:
            user_list.append(obj.verifier)
        if 'operator' in args:
            user_list.append(obj.operator)

        user_email = list(UserAccount.objects.filter(username__in=user_list).values_list('email', flat=True))
        return user_email

    def get_contact_email(self):
        cc = list(self.obj.objects.get(pk=self.latest_id).email_cc.split(','))
        contact_email = list(Contacts.objects.filter(contact_id__in=cc).values_list('contact_email', flat=True))
        return contact_email

    # 获取项目组密送成员的邮箱
    def get_bcc_email(self):
        group_id = self.obj.objects.get(pk=self.latest_id).group_id
        bcc_email = ContactsDetail.objects.filter(group__group_id=group_id).filter(bcc='1').annotate(
            contact_email=F('contact__contact_email')
        ).values_list('contact_email', flat=True)
        return list(bcc_email)
