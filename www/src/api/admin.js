import { del, get, post, put } from '@/utils/request'

export const getUsersApi = (params) => get('/api/v1/admin/users', params)
export const addUsersApi = (data) => post('/api/v1/admin/users', data)
export const updateUsersApi = (data) => put(`/api/v1/admin/users/${data.uid}`, data)
export const deleteUsersApi = (data) => del(`/api/v1/admin/users/${data}`, data)
export const ResetPasswordApi = (params) => post('/api/v1/admin/users/reset-password', params)
export const getRolesApi = (params) => get('/api/v1/admin/roles', params)
export const createRolesApi = (data) => post('/api/v1/admin/roles', data)
export const updateRolesApi = (data) => put(`/api/v1/admin/roles/${data.id}`, data)
export const deleteRolesApi = (data) => del(`/api/v1/admin/roles/${data}`, data)
export const getOrganizationsApi = (data) => get('/api/v1/admin/organizations', data)
export const createRootOrganizationsApi = (data) =>
  post('/api/v1/admin/organizations/root-node', data)
export const getOrganizationsUsersApi = (params) => get('/api/v1/admin/organizations/users', params)
export const bindOrganizationsUsersApi = (data) => post('/api/v1/admin/organizations/users', data)
export const deleteOrganizationsUsersApi = (data) => del(`/api/v1/admin/organizations/users`, data)
export const createChildOrganizationsApi = (data) =>
  post('/api/v1/admin/organizations/child-node', data)
export const updateOrganizationsApi = (data) => put(`/api/v1/admin/organizations`, data)
export const deleteOrganizationsApi = (data) => del(`/api/v1/admin/organizations`, data)
// 实例
export const getInstancesApi = (params) => get('/api/v1/admin/instances', params)
export const createInstancesApi = (data) => post('/api/v1/admin/instances', data)
export const updateInstancesApi = (data) => put(`/api/v1/admin/instances/${data.id}`, data)
export const deleteInstancesApi = (data) => del(`/api/v1/admin/instances/${data}`, data)
export const getInstanceInspectParamsApi = (params) =>
  get('/api/v1/admin/instances/inspect/params', params)
export const createInstanceInspectParamsApi = (data) =>
  post('/api/v1/admin/instances/inspect/params', data)
export const updateInstanceInspectParamsApi = (data) =>
  put(`/api/v1/admin/instances/inspect/params/${data.id}`, data)
export const deleteInstanceInspectParamsApi = (data) =>
  del(`/api/v1/admin/instances/inspect/params/${data}`, data)
// 环境
export const getEnvironmentsApi = (params) => get('/api/v1/admin/environments', params)
export const createEnvironmentsApi = (data) => post('/api/v1/admin/environment', data)
export const updateEnvironmentsApi = (data) => put(`/api/v1/admin/environment/${data.id}`, data)
export const deleteEnvironmentsApi = (data) => del(`/api/v1/admin/environment/${data}`, data)
// inspect
export const getInspectParamsApi = (params) => get('/api/v1/admin/inspect/params', params)
export const updateInspectParamsApi = (data) => put(`/api/v1/admin/inspect/params/${data.id}`, data)
// das
export const getDasSchemasListGrantApi = (params) => get('/api/v1/admin/das/schemas/grant', params)
export const createDasSchemasGrantApi = (data) => post('/api/v1/admin/das/schemas/grant', data)
export const deleteDasSchemasGrantApi = (data) =>
  del(`/api/v1/admin/das/schemas/grant/${data}`, data)
export const getDasTablesGrantApi = (params) => get('/api/v1/admin/das/tables/grant', params)
export const createDasTablesGrantApi = (data) => post('/api/v1/admin/das/tables/grant', data)
export const deleteDasTablesGrantApi = (data) => del(`/api/v1/admin/das/tables/grant/${data}`, data)
export const getDasInstancesListApi = (params) => get('/api/v1/admin/das/instances/list', params)
export const getDasSchemasListApi = (params) => get('/api/v1/admin/das/schemas/list', params)
export const getDasTablesListApi = (params) => get('/api/v1/admin/das/tables/list', params)
// 审批流
export const getApprovalFlowUnboundUsersApi = (params) =>
  get('/api/v1/admin/approval-flows/unbound-users', params)
export const getApprovalFlowsApi = (params) => get('/api/v1/admin/approval-flows', params)
export const createApprovalFlowsApi = (data) => post('/api/v1/admin/approval-flows', data)
export const updateApprovalFlowsApi = (data) => put(`/api/v1/admin/approval-flows/${data.id}`, data)
export const deleteApprovalFlowsApi = (data) => del(`/api/v1/admin/approval-flows/${data}`, data)
export const bindUsersToApprovalFlowApi = (data) =>
  post(`/api/v1/admin/approval-flows/bind-users`, data)
export const getApprovalFlowUsersApi = (data) =>
  get('/api/v1/admin/approval-flows/bind-users', data)
export const deleteUsersFromApprovalFlowApi = (data) =>
  del(`/api/v1/admin/approval-flows/bind-users/${data.id}`, data)
