import { get, post, put } from "@/utils/request"

export const getOrderEnvironmentsApi = (params) => get('/api/v1/orders/environments', params)
export const getOrderInstancesApi = (params) => get('/api/v1/orders/instances', params)
export const getOrderSchemasApi = (params) => get('/api/v1/orders/schemas', params)
export const getOrderUsersApi = (params) => get('/api/v1/orders/users', params)
export const inspectOrderSyntaxApi = (data) => post('/api/v1/orders/inspect-syntax', data)
export const createOrderApi = (data) => post('/api/v1/orders', data)
export const getOrderListApi = (params) => get('/api/v1/orders', params)
export const getOrderDetailApi = (params) => get(`/api/v1/orders/${params.order_id}`)
export const getOrderApprovalStatusApi = (params) => get(`/api/v1/orders/approval/${params.order_id}`)
export const approvalOrderApi = (data) =>  put(`/api/v1/orders/approval`, data)
export const claimOrderApi = (data) => put(`/api/v1/orders/claim`, data)
export const closeOrderApi = (data) => put(`/api/v1/orders/close`, data)
