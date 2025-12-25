import { del, get, post, put } from "@/utils/request"

export const GetSchemasApi = (params) => get('/api/v1/das/schemas', params)
export const GetSchemaTablesApi = (params) => get('/api/v1/das/schema/tables', params)
export const GetHistoryApi = (params) => get('/api/v1/das/history', params)
export const GetPermittedTablesBySchemaApi = (params) => get('/api/v1/das/schema/grants', params)
export const ExecuteMySQLQueryApi = (params) => post('/api/v1/das/query/mysql', params)
export const ExecuteClickHouseQueryApi = (params) => post('/api/v1/das/query/clickhouse', params)
export const GetTableInfoApi = (params) => get('/api/v1/das/table-info', params)
export const GetDBDictApi = (params) => get('/api/v1/das/dbdict', params)
export const CreateFavoritesApi = (data) => post('/api/v1/das/favorites', data)
export const GetFavoritesApi = (params) => get('/api/v1/das/favorites', params)
export const UpdateFavoritesApi = (data) => put(`/api/v1/das/favorites/${data.id}`, data)
export const DeleteFavoritesApi = (data) => del(`/api/v1/das/favorites/${data.id}`, data)
