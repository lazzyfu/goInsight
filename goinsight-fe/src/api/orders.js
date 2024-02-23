import axios from '@/utils/request'

// 获取环境
export const getEnvironmentsApi = (params) =>
  axios.request({
    url: '/api/v1/orders/environments',
    method: 'get',
    params: params,
  })

// 获取指定环境的Instances
export const getInstancesApi = (params) =>
  axios.request({
    url: '/api/v1/orders/instances',
    method: 'get',
    params: params,
  })

// 获取指定实例的schemas
export const getSchemasApi = (params) =>
  axios.request({
    url: '/api/v1/orders/schemas',
    method: 'get',
    params: params,
  })

// 获取审核/复核/抄送人
export const getUsersApi = (params) =>
  axios.request({
    url: '/api/v1/orders/users',
    method: 'get',
    params: params,
  })

// 语法检查
export const syntaxCheckApi = (data) =>
  axios.request({
    url: `/api/v1/orders/syntax-inspect`,
    method: 'post',
    data: data,
  })

// 提交工单
export const createOrdersApi = (data) =>
  axios.request({
    url: '/api/v1/orders/commit',
    method: 'post',
    data: data,
  })

// 获取工单记录
export const getListApi = (params) =>
  axios.request({
    url: '/api/v1/orders/list',
    method: 'get',
    params: params,
  })

// 获取工单详情
export const getOrdersDetailApi = (params) =>
  axios.request({
    url: `/api/v1/orders/detail/${params}`,
    method: 'get',
    params: params,
  })

// 获取操作日志
export const getOpLogsApi = (params) =>
  axios.request({
    url: `/api/v1/orders/detail/oplogs`,
    method: 'get',
    params: params,
  })

// 审核
export const approveOrdersApi = (data) =>
  axios.request({
    url: `/api/v1/orders/operate/approve`,
    method: 'put',
    data: data,
  })

// 反馈
export const feedbackOrdersApi = (data) =>
  axios.request({
    url: `/api/v1/orders/operate/feedback`,
    method: 'put',
    data: data,
  })

// 复核
export const reviewOrdersApi = (data) =>
  axios.request({
    url: `/api/v1/orders/operate/review`,
    method: 'put',
    data: data,
  })

// 关闭
export const closeOrdersApi = (data) =>
  axios.request({
    url: `/api/v1/orders/operate/close`,
    method: 'put',
    data: data,
  })

// hook
export const hookOrdersApi = (data) =>
  axios.request({
    url: `/api/v1/orders/hook`,
    method: 'post',
    data: data,
  })

// generate orders tasks
export const generateTasksApi = (data) =>
  axios.request({
    url: `/api/v1/orders/generate-tasks`,
    method: 'post',
    data: data,
  })

// 获取tasks
export const getTasksApi = (params) =>
  axios.request({
    url: `/api/v1/orders/tasks/${params.order_id}`,
    method: 'get',
    params: params,
  })

// 预览任务
export const previewTasksApi = (params) =>
  axios.request({
    url: `/api/v1/orders/tasks/preview`,
    method: 'get',
    params: params,
  })

// 执行单个任务
export const executeSingleTaskApi = (data) =>
  axios.request({
    url: `/api/v1/orders/tasks/execute-single`,
    method: 'post',
    data: data,
  })

// 批量执行任务
export const executeAllTaskApi = (data) =>
  axios.request({
    url: `/api/v1/orders/tasks/execute-all`,
    method: 'post',
    data: data,
  })

// 下载导出文件
export const downloadExportFile = (params) =>
  axios.request({
    url: `/api/v1/orders/download/exportfile/${params}`,
    method: 'get',
    responseType: 'blob',
  })
