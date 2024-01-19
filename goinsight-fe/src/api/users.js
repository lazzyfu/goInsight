import axios from '@/utils/request'

// 用户
export const getUsersApi = (params) =>
  axios.request({
    url: '/api/v1/admin/users',
    method: 'get',
    params: params,
  })

export const createUsersApi = (data) =>
  axios.request({
    url: '/api/v1/admin/users',
    method: 'post',
    data: data,
  })

export const updateUsersApi = (data) =>
  axios.request({
    url: `/api/v1/admin/users/${data.uid}`,
    method: 'put',
    data: data,
  })

export const deleteUsersApi = (data) =>
  axios.request({
    url: `/api/v1/admin/users/${data}`,
    method: 'delete',
    data: data,
  })

export const changeUsersPassApi = (params) =>
  axios.request({
    url: '/api/v1/admin/users/change/password',
    method: 'post',
    data: params,
  })

// 角色
export const getRolesApi = (params) =>
  axios.request({
    url: '/api/v1/admin/roles',
    method: 'get',
    params: params,
  })

export const createRolesApi = (data) =>
  axios.request({
    url: '/api/v1/admin/roles',
    method: 'post',
    data: data,
  })

export const updateRolesApi = (data) =>
  axios.request({
    url: `/api/v1/admin/roles/${data.id}`,
    method: 'put',
    data: data,
  })

export const deleteRolesApi = (data) =>
  axios.request({
    url: `/api/v1/admin/roles/${data}`,
    method: 'delete',
    data: data,
  })

// 组织
export const getOrganizationsApi = (params) =>
  axios.request({
    url: '/api/v1/admin/organizations',
    method: 'get',
    params: params,
  })

export const createRootOrganizationsApi = (data) =>
  axios.request({
    url: '/api/v1/admin/organizations/root-node',
    method: 'post',
    data: data,
  })

export const createChildOrganizationsApi = (data) =>
  axios.request({
    url: '/api/v1/admin/organizations/child-node',
    method: 'post',
    data: data,
  })

export const updateOrganizationsApi = (data) =>
  axios.request({
    url: `/api/v1/admin/organizations`,
    method: 'put',
    data: data,
  })

export const deleteOrganizationsApi = (data) =>
  axios.request({
    url: `/api/v1/admin/organizations`,
    method: 'delete',
    data: data,
  })

export const getOrganizationsUsersApi = (params) =>
  axios.request({
    url: '/api/v1/admin/organizations/users',
    method: 'get',
    params: params,
  })

export const bindOrganizationsUsersApi = (data) =>
  axios.request({
    url: '/api/v1/admin/organizations/users',
    method: 'post',
    data: data,
  })

export const deleteOrganizationsUsersApi = (data) =>
  axios.request({
    url: `/api/v1/admin/organizations/users`,
    method: 'delete',
    data: data,
  })

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
