# -*- coding:utf-8 -*-
# edit by fuzongfei
import datetime

from rest_framework import serializers

from orders.models import Orders


class GetSqlCountSerializer(serializers.Serializer):
    starttime = serializers.CharField(required=True, min_length=1, error_messages={
        'required': '开始时间不能为空',
    })
    endtime = serializers.CharField(required=True, min_length=1, error_messages={
        'required': '结束时间不能为空',
    })

    def query(self):
        sdata = self.validated_data
        starttime = sdata.get('starttime')
        endtime = sdata.get('endtime')

        if starttime == 'None' and endtime == 'None':
            tim = datetime.datetime.now()
            endtime = tim.strftime('%Y-%m-%d')
            starttime = (tim + datetime.timedelta(days=-30)).strftime('%Y-%m-%d')

        result = []

        query = f"select s.id, e.envi_name as enviname, `database` as db, sql_type as sqltype, count(*) as num " \
            f"from auditsql_orders s join auditsql_sys_environment e on s.envi_id = e.envi_id " \
            f"where s.created_at >= '{starttime}' and s.created_at <= '{endtime}' " \
            f"group by s.envi_id, s.`database`, s.sql_type"

        for row in Orders.objects.raw(query):
            result.append({
                'enviname': row.enviname,
                'db': row.db,
                'sqltype': row.sqltype,
                'num': row.num
            })
        return result
