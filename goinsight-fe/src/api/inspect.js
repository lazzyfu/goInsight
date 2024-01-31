import axios from '@/utils/request'

// dbconfig
export const adminGetInspectParamsApi = (params) =>
  axios.request({
    url: '/api/v1/admin/inspect/params',
    method: 'get',
    params: params,
  })

export const adminUpdateInspectParamsApi = (data) =>
  axios.request({
    url: `/api/v1/admin/inspect/params/${data.id}`,
    method: 'put',
    data: data,
  })
