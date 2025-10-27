import { get, post } from "@/utils/request"

export const Login = (params) => post('/api/v1/user/login', params)
export const GetUserProfileApi = (params) => get('/api/v1/user/profile', params)
