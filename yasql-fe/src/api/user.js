import axios from '@/utils/request'

// 更新用户的信息
export const UpdateUserInfo = params =>
  axios.request({
    url: `/users/update/${params.username}`,
    method: 'put',
    data: params
  })

// 更改密码
export const ChangePassword = params =>
  axios.request({
    url: '/users/change/password',
    method: 'post',
    data: params
  })

// 更改密码
export const ChangeAvatar = params =>
  axios.request({
    url: '/users/change/avatar',
    method: 'post',
    data: params,
    headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
    contentType: false, 
    processData: false,
  })

