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
  get(`/api/v1/orders/approval/${params.order_id}`)
export const getOrderLogsApi = (params) => get(`/api/v1/orders/logs/${params.order_id}`)
export const approvalOrderApi = (data) => put(`/api/v1/orders/approval`, data)
export const claimOrderApi = (data) => put(`/api/v1/orders/claim`, data)
export const revokeOrderApi = (data) => put(`/api/v1/orders/revoke`, data)
export const transferOrderApi = (data) => put(`/api/v1/orders/transfer`, data)
export const completeOrderApi = (data) => put(`/api/v1/orders/complete`, data)
export const failOrderApi = (data) => put(`/api/v1/orders/fail`, data)
export const reviewOrderApi = (data) => put(`/api/v1/orders/review`, data)
export const generateOrderTasksApi = (data) => post(`/api/v1/orders/generate-tasks`, data)
// 获取tasks
export const getOrderTasksApi = (params) => get(`/api/v1/orders/tasks/${params.order_id}`, params)
export const executeTaskApi = (data) => post(`/api/v1/orders/tasks/execute`, data)
export const executebatchTasksApi = (data) => post(`/api/v1/orders/tasks/execute-batch`, data)
// 下载导出文件
export const downloadExportFileApi = (params) =>
  getBlob(`/api/v1/orders/tasks/download/exportfile/${params}`)
