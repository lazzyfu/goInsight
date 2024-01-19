import axios from '@/utils/request'

export const UpdateUserInfo = (params) =>
  axios.request({
    url: `/api/v1/user/${params.uid}`,
    method: 'put',
    data: params,
  })

export const ChangePassword = (params) =>
  axios.request({
    url: '/api/v1/user/change/password',
    method: 'post',
    data: params,
  })

export const ChangeAvatarApi = (params) =>
  axios.request({
    url: '/api/v1/user/change/avatar',
    method: 'post',
    data: params,
    headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
    contentType: false,
    processData: false,
  })

export const GetOTPAuthURLApi = (params) =>
  axios.request({
    url: '/api/v1/user/otp-auth-url',
    method: 'get',
    params: params,
  })

export const GetOTPCallbackApi = (params) =>
  axios.request({
    url: '/api/v1/user/otp-auth-callback',
    method: 'get',
    params: params,
  })
