import { post, put } from '@/utils/request'

export const UpdateUserInfoApi = (params) => put(`/api/v1/profile/${params.uid}`, params)
export const ChangePasswordApi = (params) => post('/api/v1/profile/change/password', params)

export const ChangeAvatarApi = (params) =>
  post('/api/v1/profile/change/avatar', params, {
    headers: { 'Content-Type': 'multipart/form-data' },
    contentType: false,
    processData: false,
  })

export const GetOTPAuthURLApi = (params) => post('/api/v1/user/otp-auth-url', params)
export const GetOTPCallbackApi = (params) => post('/api/v1/user/otp-auth-callback', params)
