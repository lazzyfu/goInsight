import { get, post } from "@/utils/request"

export const getUsersApi = (params) => get('/api/v1/admin/users', params)
export const changePasswordApi = (params) => post('/api/v1/admin/users/change-password', params)
