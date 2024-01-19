import axios from '@/utils/request'

// dbconfig
export const adminGetDBConfigApi = (params) =>
  axios.request({
    url: '/api/v1/admin/dbconfig',
    method: 'get',
    params: params,
  })

export const adminCreateDBConfigApi = (data) =>
  axios.request({
    url: '/api/v1/admin/dbconfig',
    method: 'post',
    data: data,
  })

export const adminUpdateDBConfigApi = (data) =>
  axios.request({
    url: `/api/v1/admin/dbconfig/${data.id}`,
    method: 'put',
    data: data,
  })

export const adminDeleteDBConfigApi = (data) =>
  axios.request({
    url: `/api/v1/admin/dbconfig/${data}`,
    method: 'delete',
    data: data,
  })

// environment
export const adminGetEnvironmentsApi = (params) =>
  axios.request({
    url: '/api/v1/admin/environment',
    method: 'get',
    params: params,
  })

export const adminCreateEnvironmentsApi = (data) =>
  axios.request({
    url: '/api/v1/admin/environment',
    method: 'post',
    data: data,
  })

export const adminUpdateEnvironmentsApi = (data) =>
  axios.request({
    url: `/api/v1/admin/environment/${data.id}`,
    method: 'put',
    data: data,
  })

export const adminDeleteEnvironmentsApi = (data) =>
  axios.request({
    url: `/api/v1/admin/environment/${data}`,
    method: 'delete',
    data: data,
  })
