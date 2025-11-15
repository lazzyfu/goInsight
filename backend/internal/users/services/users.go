package services

import (
	"errors"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/lazzyfu/goinsight/internal/global"

	"github.com/lazzyfu/goinsight/pkg/pagination"

	"github.com/lazzyfu/goinsight/internal/users/forms"
	"github.com/lazzyfu/goinsight/internal/users/models"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type GetUsersServices struct {
	*forms.GetUsersForm
	C *gin.Context
}

func (s *GetUsersServices) Run() (responseData interface{}, total int64, err error) {
	type user struct {
		models.InsightUsers
		Organization    string `json:"organization,omitempty"`
		OrganizationKey string `json:"organization_key,omitempty"`
		Role            string `json:"role,omitempty"`
		RoleID          uint64 `json:"role_id,omitempty"`
	}
	var users []user
	tx := global.App.DB.Table("insight_users a").
		Select(`a.*, c.key as organization_key, ifnull(d.name, '-/-') as role,ifnull(
			concat(
				(
					SELECT
						GROUP_CONCAT(
							ia.name
							ORDER BY
								ia.name ASC SEPARATOR '/'
						) AS concatenated_names
					FROM
						insight_organizations ia
					WHERE
						EXISTS (
							SELECT
								1
							FROM
								insight_organizations
							WHERE
								JSON_CONTAINS(c.path, CONCAT('\"', ia.key, '\"'))
						)
				),
				'/',
				c.name
			),
			c.name
		) as organization`).
		Joins("left join insight_organizations_users b on a.uid=b.uid").
		Joins("left join insight_organizations c on b.organization_key=c.key").
		Joins("left join insight_roles d on d.id=a.role_id")
	// 搜索
	if s.Search != "" {
		tx = tx.Where("`username` like ? or `nick_name` like ? or `email` like ? or `mobile` like ?", "%"+s.Search+"%", "%"+s.Search+"%", "%"+s.Search+"%", "%"+s.Search+"%")
	}
	if s.OrganizationKey != "" {
		tx = tx.Where("c.key like ?", s.OrganizationKey+"%")
	}
	if s.RoleID > 0 {
		tx = tx.Where("a.role_id=?", s.RoleID)
	}
	total = pagination.Pager(&s.PaginationQ, tx, &users)
	return &users, total, nil
}

type CreateUsersService struct {
	*forms.CreateUsersForm
	C *gin.Context
}

func (s *CreateUsersService) Run() error {
	// 加密密码
	hashedPassword := models.BcryptPW(s.Password)
	user := models.InsightUsers{
		Username:    s.Username,
		Password:    hashedPassword,
		Email:       s.Email,
		NickName:    s.NickName,
		Mobile:      s.Mobile,
		RoleID:      s.RoleID,
		IsTwoFA:     s.IsTwoFA,
		IsSuperuser: s.IsSuperuser,
		IsActive:    s.IsActive,
	}
	return global.App.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.InsightUsers{}).Create(&user).Error; err != nil {
			mysqlErr := err.(*mysql.MySQLError)
			switch mysqlErr.Number {
			case 1062:
				return fmt.Errorf("用户`%s`已存在", s.Username)
			}
			global.App.Log.Error(err)
			return err
		}
		return nil
	})

}

type UpdateUsersService struct {
	*forms.UpdateUsersForm
	C   *gin.Context
	UID uint64
}

func (s *UpdateUsersService) Run() error {
	return global.App.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.InsightUsers{}).Where("uid=?", s.UID).Updates(map[string]interface{}{
			"username":     s.Username,
			"email":        s.Email,
			"nick_name":    s.NickName,
			"mobile":       s.Mobile,
			"role_id":      s.RoleID,
			"is_two_fa":    s.IsTwoFA,
			"is_active":    s.IsActive,
			"is_superuser": s.IsSuperuser,
		}).Error; err != nil {
			mysqlErr := err.(*mysql.MySQLError)
			switch mysqlErr.Number {
			case 1062:
				return fmt.Errorf("用户`%s`已存在", s.Username)
			}
			global.App.Log.Error(err)
			return err
		}
		return nil
	})

}

type DeleteUsersService struct {
	C   *gin.Context
	UID uint64
}

func (s *DeleteUsersService) Run() error {
	tx := global.App.DB.Where("uid=?", s.UID).Delete(&models.InsightUsers{})
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

type ResetUsersPasswordService struct {
	*forms.ResetUsersPasswordForm
	C *gin.Context
}

func (s *ResetUsersPasswordService) Run() error {
	if s.Password != s.VerifyPassword {
		return errors.New("您两次输入的密码不一致")
	}
	// 加密密码
	hashedPassword := models.BcryptPW(s.Password)
	global.App.DB.Model(&models.InsightUsers{}).Where("uid=?", s.UID).Update("password", hashedPassword)
	return nil
}

type ChangeUserAvatarService struct {
	C        *gin.Context
	Username string
	File     *multipart.FileHeader
}

func (s *ChangeUserAvatarService) Run() error {
	// 保存图片文件
	fileName := fmt.Sprintf("%s_%d.jpg", s.Username, time.Now().Unix())
	err := s.C.SaveUploadedFile(s.File, "./media/avatars/"+fileName)
	if err != nil {
		return err
	}

	// 调用业务逻辑
	global.App.DB.Model(&models.InsightUsers{}).Where("username=?", s.Username).Update("avatar_file", "/media/avatars/"+fileName)
	if err != nil {
		return err
	}
	return nil
}
