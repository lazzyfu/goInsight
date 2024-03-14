import axios from '@/utils/request'

// 获取环境
export const getEnvironmentsApi = (params) =>
  axios.request({
    url: '/api/v1/das/environments',
    method: 'get',
    params: params,
  })

// 获取授权的schemas
export const getSchemasApi = (params) =>
  axios.request({
    url: '/api/v1/das/schemas',
    method: 'get',
    params: params,
  })

// 获取授权的tables
export const getTablesApi = (params) =>
  axios.request({
    url: '/api/v1/das/tables',
    method: 'get',
    params: params,
  })

// 执行MySQL/TiDB查询接口
export const executeMySQLQueryApi = (data) =>
  axios.request({
    url: '/api/v1/das/execute/query/mysql',
    method: 'post',
    data: data,
  })

// 执行Clickhouse查询接口
export const executeClickHouseQueryApi = (data) =>
  axios.request({
    url: '/api/v1/das/execute/query/clickhouse',
    method: 'post',
    data: data,
  })

// 获取指定库有权限的表
export const getUserGrantsApi = (params) =>
  axios.request({
    url: '/api/v1/das/user/grants',
    method: 'get',
    params: params,
  })

// 获取数据字典
export const getDBDictApi = (params) =>
  axios.request({
    url: '/api/v1/das/dbdict',
    method: 'get',
    params: params,
  })

// 获取历史SQL
export const getHistory = (params) =>
  axios.request({
    url: '/api/v1/das/history',
    method: 'get',
    params: params,
  })

// 获取收藏的SQL
export const getFavoritesApi = (params) =>
  axios.request({
    url: '/api/v1/das/favorites',
    method: 'get',
    params: params,
  })

// 新增收藏的SQL
export const createFavoritesApi = (data) =>
  axios.request({
    url: '/api/v1/das/favorites',
    method: 'post',
    data: data,
  })

// 更新收藏的SQL
export const updateFavoritesApi = (data) =>
  axios.request({
    url: `/api/v1/das/favorites/${data.id}`,
    method: 'put',
    data: data,
  })

// 删除收藏的SQL
export const deleteFavoritesApi = (data) =>
  axios.request({
    url: `/api/v1/das/favorites/${data}`,
    method: 'delete',
    data: data,
  })

// 获取表结果或表元信息
export const getTableInfo = (params) =>
  axios.request({
    url: '/api/v1/das/table-info',
    method: 'get',
    params: params,
  })

// admin
export const adminGetSchemasListGrantApi = (params) =>
  axios.request({
    url: '/api/v1/das/admin/schemas/grant',
    method: 'get',
    params: params,
  })

export const adminCreateSchemasGrantApi = (data) =>
  axios.request({
    url: '/api/v1/das/admin/schemas/grant',
    method: 'post',
    data: data,
  })

export const adminDeleteSchemasGrantApi = (data) =>
  axios.request({
    url: `/api/v1/das/admin/schemas/grant/${data}`,
    method: 'delete',
    data: data,
  })

export const adminGetTablesGrantApi = (params) =>
  axios.request({
    url: '/api/v1/das/admin/tables/grant',
    method: 'get',
    params: params,
  })

export const adminCreateTablesGrantApi = (data) =>
  axios.request({
    url: '/api/v1/das/admin/tables/grant',
    method: 'post',
    data: data,
  })

export const adminDeleteTablesGrantApi = (data) =>
  axios.request({
    url: `/api/v1/das/admin/tables/grant/${data}`,
    method: 'delete',
    data: data,
  })

export const adminGetInstancesListApi = (params) =>
  axios.request({
    url: '/api/v1/das/admin/instances/list',
    method: 'get',
    params: params,
  })

export const adminGetSchemasListApi = (params) =>
  axios.request({
    url: '/api/v1/das/admin/schemas/list',
    method: 'get',
    params: params,
  })

export const adminGetTablesListApi = (params) =>
  axios.request({
    url: '/api/v1/das/admin/tables/list',
    method: 'get',
    params: params,
  })
