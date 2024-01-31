package forms

import (
	"goInsight/internal/pkg/pagination"
)

type GetUsersForm struct {
	PaginationQ     pagination.Pagination
	OrganizationKey string `form:"organization_key"`
	RoleID          uint64 `form:"role_id"`
	Search          string `form:"search"`
}

type CreateUsersForm struct {
	Username    string `form:"username"  json:"username" binding:"required,min=2,max=32"`
	Password    string `form:"password"  json:"password" binding:"required,min=7,max=128"`
	Email       string `form:"email" json:"email" binding:"required,min=3,max=254"`
	NickName    string `form:"nick_name" json:"nick_name" binding:"required,min=1,max=32"`
	Mobile      string `form:"mobile" json:"mobile"`
	RoleID      uint64 `form:"role_id" json:"role_id"`
	IsTwoFA     bool   `form:"is_two_fa" json:"is_two_fa" validate:"boolean"`
	IsSuperuser bool   `form:"is_superuser" json:"is_superuser" validate:"boolean"`
	IsActive    bool   `form:"is_active" json:"is_active" validate:"boolean" `
}

type UpdateUsersForm struct {
	Username    string `form:"username"  json:"username" binding:"required,min=2,max=32"`
	Email       string `form:"email" json:"email" binding:"required,min=3,max=254"`
	NickName    string `form:"nick_name" json:"nick_name" binding:"required,min=1,max=32"`
	Mobile      string `form:"mobile" json:"mobile"`
	RoleID      uint64 `form:"role_id" json:"role_id"`
	IsTwoFA     bool   `form:"is_two_fa" json:"is_two_fa" validate:"boolean"`
	IsSuperuser bool   `form:"is_superuser" json:"is_superuser" validate:"boolean"`
	IsActive    bool   `form:"is_active" json:"is_active" validate:"boolean"`
}

type ChangeUsersPasswordForm struct {
	UID            uint64 `form:"uid" json:"uid" binding:"required"`
	Password       string `form:"password" json:"password" binding:"required,min=7,max=32"`
	VerifyPassword string `form:"verify_password" json:"verify_password" binding:"required,min=7,max=32"`
}
