import axios from '@/utils/request'

// 获取tree
export const getQueryTree = params =>
  axios.request({
    url: '/api/sqlquery/tree',
    method: 'post',
    data: params
  })

// 执行查询
export const ExecuteQuery = params =>
  axios.request({
    url: '/api/sqlquery/execute-query',
    method: 'post',
    data: params
  })

// 删除会话
export const deleteQuerySession = params =>
  axios.request({
    url: '/api/sqlquery/delete-query-hash',
    method: 'post',
    data: params
  })

// 获取表信息
export const getTableInfo = params =>
  axios.request({
    url: '/api/sqlquery/get/tableinfo',
    method: 'post',
    data: params
  })

// 获取历史SQL
export const getHistorySql = params =>
  axios.request({
    url: '/api/sqlquery/get/history/sql',
    method: 'get',
    params: params
  })

// 获取数据字典
export const getDBDict = params =>
  axios.request({
    url: '/api/sqlquery/get/dbdict',
    method: 'post',
    data: params
  })
