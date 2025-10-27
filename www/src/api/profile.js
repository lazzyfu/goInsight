import { get, post, put } from "@/utils/request";

export const UpdateUserInfoApi = (params) => put(`/api/v1/user/${params.uid}`, params)
export const ChangePasswordApi = (params) => post('/api/v1/user/change/password', params)

export const ChangeAvatarApi = (params) => post('/api/v1/user/change/avatar', params, {
  headers: { 'Content-Type': 'multipart/form-data' },
  contentType: false,
  processData: false,
});

export const GetOTPAuthURLApi = (params) => get('/api/v1/user/otp-auth-url', params)
export const GetOTPCallbackApi = (params) => get('/api/v1/user/otp-auth-callback', params)
