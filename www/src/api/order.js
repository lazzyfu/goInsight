import { get, post } from "@/utils/request"

export const getEnvironmentsApi = (params) => get('/api/v1/orders/environments', params)
export const getInstancesApi = (params) => get('/api/v1/orders/instances', params)
export const getSchemasApi = (params) => get('/api/v1/orders/schemas', params)
export const getUsersApi = (params) => get('/api/v1/orders/users', params)
export const createOrdersApi = (data) => post('/api/v1/orders/commit', data)
export const inspectSQLSyntaxApi = (data) => post('/api/v1/orders/inspect-sql-syntax', data)
export const getOrdersListApi = (params) => get('/api/v1/orders/list', params)
