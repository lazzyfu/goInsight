import { del, get, post, put } from "@/utils/request"

export const getUsersApi = (params) => get('/api/v1/admin/users', params)
export const addUsersApi = (data) => post('/api/v1/admin/users', data)
export const updateUsersApi = (data) => put(`/api/v1/admin/users/${data.uid}`, data)
export const deleteUsersApi = (data) => del(`/api/v1/admin/users/${data}`, data)
export const changePasswordApi = (params) => post('/api/v1/admin/users/change-password', params)
export const getRolesApi = (params) => get('/api/v1/admin/roles', params)
