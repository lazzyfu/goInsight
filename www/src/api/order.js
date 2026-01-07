import { get, getBlob, post, put } from '@/utils/request'

export const getOrderEnvironmentsApi = (params) => get('/api/v1/orders/environments', params)
export const getOrderInstancesApi = (params) => get('/api/v1/orders/instances', params)
export const getOrderSchemasApi = (params) => get('/api/v1/orders/schemas', params)
export const getOrderUsersApi = (params) => get('/api/v1/orders/users', params)

export const inspectOrderSyntaxApi = (data) => post('/api/v1/orders/inspect-syntax', data)

export const createOrderApi = (data) => post('/api/v1/orders', data)
export const getOrderListApi = (params) => get('/api/v1/orders', params)
export const getOrderDetailApi = (params) => get(`/api/v1/orders/${params.order_id}`)

export const getOrderApprovalStatusApi = (params) =>
  get(`/api/v1/orders/${params.order_id}/approvals`)
export const getOrderLogsApi = (params) => get(`/api/v1/orders/${params.order_id}/logs`)

export const approvalOrderApi = (data) => put(`/api/v1/orders/${data.order_id}/actions/approval`, data)
export const claimOrderApi = (data) => put(`/api/v1/orders/${data.order_id}/actions/claim`, data)
export const revokeOrderApi = (data) => put(`/api/v1/orders/${data.order_id}/actions/revoke`, data)
export const transferOrderApi = (data) => put(`/api/v1/orders/${data.order_id}/actions/transfer`, data)
export const completeOrderApi = (data) => put(`/api/v1/orders/${data.order_id}/actions/complete`, data)
export const failOrderApi = (data) => put(`/api/v1/orders/${data.order_id}/actions/fail`, data)
export const reviewOrderApi = (data) => put(`/api/v1/orders/${data.order_id}/actions/review`, data)

// 生成tasks
export const generateOrderTasksApi = (data) => post(`/api/v1/orders/${data.order_id}/tasks`, data)

// 获取tasks
export const getOrderTasksApi = (params) => get(`/api/v1/orders/${params.order_id}/tasks`, params)
export const executeTaskApi = (data) => post(`/api/v1/orders/${data.order_id}/tasks/execute`, data)
export const executebatchTasksApi = (data) =>
  post(`/api/v1/orders/${data.order_id}/tasks/execute-batch`, data)

// 下载导出文件
export const downloadExportFileApi = (params) =>
  getBlob(`/api/v1/orders/${params.order_id}/tasks/exports/${params}`)
