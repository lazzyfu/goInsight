import axios from '@/utils/request'

// 获取系统仪表盘
export const getSysDash = params =>
  axios.request({
    url: '/users/dashboard/sys',
    method: 'get',
    params: params
  })

// 获取用户仪表盘
export const getSelfDash = params =>
  axios.request({
    url: '/users/dashboard/self',
    method: 'get',
    params: params
  })
