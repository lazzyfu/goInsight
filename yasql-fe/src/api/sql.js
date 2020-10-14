import axios from '@/utils/request'

// 获取用户
export const getUsers = axios.request({
  url: '/users/list',
  method: 'get'
})

// 获取工单环境
export const getDbEnvironment = axios.request({
  url: '/sqlorders/envs',
  method: 'get'
})

// 获取上线版本
export const getReleaseVersions = axios.request({
  url: '/sqlorders/versions/get',
  method: 'get'
})

// 获取指定环境的schemas
export const getDbSchemas = params =>
  axios.request({
    url: '/sqlorders/schemas',
    method: 'get',
    params: params
  })

// 语法检查
export const incepSyntaxCheck = params =>
  axios.request({
    url: '/sqlorders/incep/syntaxcheck',
    method: 'post',
    data: params
  })

// 提交sql工单
export const commitSqlOrders = params =>
  axios.request({
    url: '/sqlorders/commit',
    method: 'post',
    data: params
  })

// 获取指定环境的工单列表
export const getSqlOrdersList = params =>
  axios.request({
    url: '/sqlorders/list',
    method: 'get',
    params: params
  })

// 获取工单详情
export const getSqlOrdersDetail = params =>
  axios.request({
    url: `/sqlorders/detail/${params}`,
    method: 'get'
  })

// 更新SQL工单状态，如：审核，反馈，关闭等
export const opSqlOrders = params =>
  axios.request({
    url: `/sqlorders/op/${params.action}/${params.pk}`,
    method: 'put',
    data: params
  })

// 生成工单执行任务
export const generateSqlOrdersExecuteTasks = params =>
  axios.request({
    url: '/sqlorders/tasks/generate',
    method: 'post',
    data: params
  })

// 获取工单的执行预览任务列表
export const getSqlOrdersTasksPreviewList = params =>
  axios.request({
    url: `/sqlorders/tasks/preview/${params.task_id}`,
    method: 'get',
    params: params
  })

// 根据orderid获取taskid
export const getSqlOrdersTaskId = params =>
  axios.request({
    url: `/sqlorders/tasks/get/${params.order_id}`,
    method: 'get',
    params: params
  })

// 获取工单的执行任务列表
export const getSqlOrdersTasksList = params =>
  axios.request({
    url: `/sqlorders/tasks/list/${params.task_id}`,
    method: 'get',
    params: params
  })

// 执行单个任务
export const executeSingleTask = params =>
  axios.request({
    url: '/sqlorders/tasks/execute/single',
    method: 'post',
    data: params
  })

// 执行全部任务
export const executeMultiTask = params =>
  axios.request({
    url: '/sqlorders/tasks/execute/multi',
    method: 'post',
    data: params
  })

// 获取任务的结果
export const getTasksResult = id =>
  axios.request({
    url: `/sqlorders/tasks/result/${id}`,
    method: 'get'
  })

// 节流任务
export const TaskThrottle = params =>
  axios.request({
    url: '/sqlorders/tasks/throttle',
    method: 'post',
    data: params
  })

// hook
export const HookSqlOrders = params =>
  axios.request({
    url: '/sqlorders/hook',
    method: 'post',
    data: params
  })

// 下载导出文件
export const downloadExportFiles = params =>
  axios.request({
    url: `/sqlorders/export/download/${params}`,
    method: 'get',
    responseType: 'blob'
  })

// 获取上线版本列表 - 管理使用
export const listReleaseVersions = params =>
  axios.request({
    url: '/sqlorders/versions/list',
    method: 'get',
    params: params
  })

// 新建版本
export const createReleaseVersions = params =>
  axios.request({
    url: '/sqlorders/versions/create',
    method: 'post',
    data: params
  })

// 更新上线版本信息
export const updateReleaseVersions = params =>
  axios.request({
    url: `/sqlorders/versions/update/${params.key}`,
    method: 'put',
    data: params
  })

// 删除版本
export const deleteReleaseVersions = params =>
  axios.request({
    url: `/sqlorders/versions/delete/${params}`,
    method: 'delete',
    data: params
  })

// 获取指定版本内工单在所有环境的进度
export const viewReleaseVersions = params =>
  axios.request({
    url: `/sqlorders/versions/view/${params}`,
    method: 'get',
    params: params
  })
