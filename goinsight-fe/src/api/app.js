import axios from '@/utils/request'

export const getAppTitleApi = (params) =>
  axios.request({
    url: '/api/v1/app/title',
    method: 'get',
    params: params,
  })
