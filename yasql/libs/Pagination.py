# -*- coding:utf-8 -*-
# edit by xff

from rest_framework.pagination import PageNumberPagination


class Pagination(PageNumberPagination):
    page_size = 10  # 默认每页显示条数
    max_page_size = 30  # 每页最大显示条数
    page_query_param = 'page'  # 当前页面值的参数
    page_size_query_param = "page_size"  # 当前页面大小参数

    def get_page_size(self, request):
        try:
            page_size = int(request.query_params.get('page_size'))
        except:
            page_size = self.page_size
        if page_size > 0:
            return min(page_size, self.max_page_size)
        else:
            return 10000  # 为非正整数时，限制为10000
