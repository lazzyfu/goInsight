package forms

import "github.com/lazzyfu/goinsight/pkg/pagination"

type CreateRootOrganizationsForm struct {
	Name string `form:"name"  json:"name" binding:"required,min=2,max=32"`
}

type CreateChildOrganizationsForm struct {
	ParentNodeName string `form:"parent_node_name"  json:"parent_node_name" binding:"required,min=2,max=32"`
	ParentNodeKey  string `form:"parent_node_key"  json:"parent_node_key" binding:"required,min=3,max=256"`
	Name           string `form:"name"  json:"name" binding:"required,min=2,max=32"`
}

type UpdateOrganizationsForm struct {
	Key  string `form:"key"  json:"key" binding:"required,min=3,max=256"`
	Name string `form:"name"  json:"name" binding:"required,min=2,max=32"`
}

type DeleteOrganizationsForm struct {
	Key  string `form:"key"  json:"key" binding:"required,min=3,max=256"`
	Name string `form:"name"  json:"name" binding:"required,min=2,max=32"`
}

type GetOrganizationsUsersForm struct {
	PaginationQ pagination.Pagination
	Key         string `form:"key"  json:"key" binding:"required,min=3,max=256"`
	Search      string `form:"search"`
}

type BindOrganizationsUsersForm struct {
	Key   string   `form:"key"  json:"key" binding:"required,min=3,max=256"`
	Users []uint64 `form:"users"  json:"users" binding:"required"`
}

type DeleteOrganizationsUsersForm struct {
	Key string `form:"key"  json:"key" binding:"required,min=3,max=256"`
	Uid uint64 `form:"uid"  json:"uid" binding:"required"`
}
