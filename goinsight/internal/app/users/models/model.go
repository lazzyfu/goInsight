/*
@Time    :   2023/08/14 16:51:22
@Author  :   lazzyfu
*/

package models

import (
	"goInsight/internal/app/common/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/datatypes"
)

func BcryptPW(password string) string {
	cost := 10
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), cost)
	return string(hashedPassword)
}

// 用户表
type InsightUsers struct {
	Uid         uint64           `gorm:"type:bigint;primaryKey;autoIncrement;comment:用户ID" json:"uid"`
	Username    string           `gorm:"type:varchar(32);not null;uniqueIndex:uniq_username;comment:用户名" json:"username"`
	Password    string           `gorm:"type:varchar(128);not null;comment:密码" json:"password"`
	Email       string           `gorm:"type:varchar(254);not null;comment:email" json:"email"`
	NickName    string           `gorm:"type:varchar(32);not null;comment:显示名" json:"nick_name"`
	Mobile      string           `gorm:"type:varchar(11);not null;comment:手机号" json:"mobile"`
	AvatarFile  string           `gorm:"type:varchar(254);not null;comment:头像文件地址" json:"avatar_file"`
	RoleID      uint64           `gorm:"type:bigint;null;comment:角色ID" json:"role_id"`
	IsSuperuser bool             `gorm:"type:boolean;default:false;comment:是否为超级管理员" json:"is_superuser"`
	IsActive    bool             `gorm:"type:boolean;default:true;comment:是否激活" json:"is_active"`
	IsStaff     bool             `gorm:"type:boolean;default:false;comment:是否为员工" json:"is_staff"`
	IsTwoFA     bool             `gorm:"type:boolean;default:false;comment:是否启用2FA认证" json:"is_two_fa"`
	OtpSecret   string           `gorm:"type:varchar(128);not null;comment:otp_secret" json:"otp_secret"`
	LastLogin   models.LocalTime `gorm:"autoUpdateTime;comment:最后一次登录" json:"last_login"`
	DateJoined  models.LocalTime `gorm:"index:date_joined;autoCreateTime;comment:加入时间" json:"date_joined"`
	UpdatedAt   models.LocalTime `gorm:"index:idx_updated_at;autoUpdateTime;comment:更新时间" json:"updated_at"`
}

func (InsightUsers) TableName() string {
	return "insight_users"
}

// 角色表
type InsightRoles struct {
	*models.Model
	Name string `gorm:"type:varchar(32);not null;uniqueIndex:uniq_name;comment:角色名" json:"name"`
}

func (InsightRoles) TableName() string {
	return "insight_roles"
}

// key全局唯一
// ParentID+Name组成唯一索引,表示同一级别下节点名不能重复
type InsightOrganizations struct {
	ID        uint64           `gorm:"primaryKey" json:"id"`
	Name      string           `gorm:"type:varchar(32);not null;uniqueIndex:uniq_name;comment:节点名" json:"name"`
	ParentID  uint64           `gorm:"not null;default 0;uniqueIndex:uniq_name;comment:父节点ID,0值表示父节点" json:"parent_id"`
	Key       string           `gorm:"type:varchar(256);default null;uniqueIndex:uniq_key;comment:搜索路径" json:"key"`
	Level     uint64           `gorm:"not null;default 1;comment:当前节点到根节点的距离或者层级,父节点起始值为1" json:"level"`
	Path      datatypes.JSON   `gorm:"type:json;null;default:null;comment:绝对路径" json:"path"`
	Creator   string           `gorm:"type:varchar(64);default null;comment:创建人" json:"creator"`
	Updater   string           `gorm:"type:varchar(64);default null;comment:更新人" json:"updater"`
	CreatedAt models.LocalTime `gorm:"index:idx_created_at;autoCreateTime;comment:创建时间" json:"created_at"`
	UpdatedAt models.LocalTime `gorm:"index:idx_updated_at;autoUpdateTime;comment:更新时间" json:"updated_at"`
}

func (InsightOrganizations) TableName() string {
	return "insight_organizations"
}

type InsightOrganizationsUsers struct {
	*models.Model
	Uid             uint64 `gorm:"type:bigint;not null;uniqueIndex:uniq_uid;comment:用户ID" json:"uid"`
	OrganizationKey string `gorm:"type:varchar(256);not null;index:organization_key;comment:搜索路径" json:"organization_key"`
}

func (InsightOrganizationsUsers) TableName() string {
	return "insight_organizations_users"
}
