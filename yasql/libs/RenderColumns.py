# -*- coding:utf-8 -*-
# edit by xff


def render_dynamic_columns(render_columns):
    try:
        columns = []
        for x in render_columns:
            col = {'title': x['value'],
                   'dataIndex': x['key'],
                   'key': x['key'],
                   'scopedSlots': {'customRender': x['key']}
                   }
            if x.get('width'):
                col['width'] = x['width']
            if x.get('fixed'):
                col['fixed'] = x['fixed']
            if x.get('ellipsis'):
                col['ellipsis'] = True
            columns.append(col)
    except IndexError as err:
        columns = []
    return columns
