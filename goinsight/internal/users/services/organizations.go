package services

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/lazzyfu/goinsight/internal/global"

	"github.com/lazzyfu/goinsight/pkg/pagination"

	"github.com/lazzyfu/goinsight/internal/users/forms"
	"github.com/lazzyfu/goinsight/internal/users/models"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type GetOrganizationsServices struct {
	C *gin.Context
}

func (s *GetOrganizationsServices) getChildOrganizations(key string, level uint64) []map[string]interface{} {
	// 迭代子节点，获取所有递归的子节点
	var childNodes []models.InsightOrganizations
	global.App.DB.Table("insight_organizations").Where("`key` like ? and level=?", key+"-%", level).Scan(&childNodes)
	if len(childNodes) == 0 {
		return nil
	}
	var data []map[string]interface{} = []map[string]interface{}{}
	for _, row := range childNodes {
		var childNode map[string]interface{} = map[string]interface{}{
			"title":       row.Name,
			"key":         row.Key,
			"slots":       map[string]string{"icon": "dir"},
			"scopedSlots": map[string]string{"title": "custom"},
		}
		childData := s.getChildOrganizations(row.Key, row.Level+1)
		if childData != nil {
			childNode["children"] = childData
		} else {
			childNode["is-leaf"] = true
		}
		data = append(data, childNode)
	}
	return data
}

func (s *GetOrganizationsServices) Run() (responseData interface{}) {
	// 获取ROOT组织
	var rootNodes []models.InsightOrganizations
	global.App.DB.Table("insight_organizations").
		Where("parent_id=0 and level=1").
		Scan(&rootNodes)
	if len(rootNodes) == 0 {
		return
	}
	var data []map[string]interface{} = []map[string]interface{}{}
	for _, row := range rootNodes {
		// 迭代父节点
		var rootNode map[string]interface{} = map[string]interface{}{
			"title":       row.Name,
			"key":         row.Key,
			"slots":       map[string]string{"icon": "dir"},
			"scopedSlots": map[string]string{"title": "custom"},
		}
		// 迭代子节点
		childData := s.getChildOrganizations(row.Key, row.Level+1)
		if childData != nil {
			rootNode["children"] = childData
		}
		data = append(data, rootNode)
	}
	return data
}

type CreateRootOrganizationsService struct {
	*forms.CreateRootOrganizationsForm
	C        *gin.Context
	Username string
}

func (s *CreateRootOrganizationsService) Run() error {
	tx := global.App.DB.Model(&models.InsightOrganizations{})
	organization := models.InsightOrganizations{Name: s.Name, ParentID: 0, Creator: s.Username, Updater: s.Username, Level: 1}
	result := tx.Create(&organization)
	if result.Error != nil {
		mysqlErr := result.Error.(*mysql.MySQLError)
		switch mysqlErr.Number {
		case 1062:
			return fmt.Errorf("记录`%s`已存在", s.Name)
		}
		return result.Error
	}
	// 更新key
	key := fmt.Sprintf("0-%d", organization.ID)
	global.App.DB.Model(&models.InsightOrganizations{}).Where("id=?", organization.ID).Update("key", key)
	return nil
}

type CreateChildOrganizationsService struct {
	*forms.CreateChildOrganizationsForm
	C        *gin.Context
	Username string
}

func (s *CreateChildOrganizationsService) Run() error {
	// 判断父节点是否存在
	var parentOrganization models.InsightOrganizations
	parentResult := global.App.DB.Table("insight_organizations").Where("`key`=? and `name`=?", s.ParentNodeKey, s.ParentNodeName).First(&parentOrganization)
	if parentResult.RowsAffected == 0 {
		return fmt.Errorf("父节点%s不存在", s.ParentNodeName)
	}
	// 将path json数据转换为数组
	var pathJsonD []string
	if len(parentOrganization.Path) != 0 {
		err := json.Unmarshal([]byte(parentOrganization.Path), &pathJsonD)
		if err != nil {
			return err
		}
	}
	pathJsonD = append(pathJsonD, parentOrganization.Key)
	// 将path数据转换为json
	pathJson, err := json.Marshal(pathJsonD)
	if err != nil {
		return err
	}

	tx := global.App.DB.Model(&models.InsightOrganizations{})
	levelLength := parentOrganization.Level + 1
	organization := models.InsightOrganizations{
		Name:     s.Name,
		ParentID: parentOrganization.ID,
		Path:     datatypes.JSON(pathJson),
		Creator:  s.Username,
		Updater:  s.Username,
		Level:    levelLength,
	}
	result := tx.Create(&organization)
	if result.Error != nil {
		mysqlErr := result.Error.(*mysql.MySQLError)
		switch mysqlErr.Number {
		case 1062:
			return fmt.Errorf("记录`%s`已存在", s.Name)
		}
		return result.Error
	}
	// 更新key
	key := fmt.Sprintf("%s-%d", parentOrganization.Key, organization.ID)
	global.App.DB.Model(&models.InsightOrganizations{}).Where("id=?", organization.ID).Update("key", key)
	return nil
}

type UpdateOrganizationsService struct {
	*forms.UpdateOrganizationsForm
	C        *gin.Context
	Username string
}

func (s *UpdateOrganizationsService) Run() error {
	tx := global.App.DB.Table("insight_organizations").Where("`key`=?", s.Key)
	result := tx.Updates(map[string]interface{}{
		"name":    s.Name,
		"Updater": s.Username,
	})
	if result.Error != nil {
		mysqlErr := result.Error.(*mysql.MySQLError)
		switch mysqlErr.Number {
		case 1062:
			return fmt.Errorf("记录`%s`已存在", s.Username)
		}
		return result.Error
	}
	return nil
}

type DeleteOrganizationsService struct {
	*forms.DeleteOrganizationsForm
	C *gin.Context
}

func (s *DeleteOrganizationsService) Run() error {
	var organization models.InsightOrganizations
	result := global.App.DB.Table("insight_organizations").Where("`key`=? and `name`=?", s.Key, s.Name).First(&organization)
	if result.RowsAffected == 0 {
		return fmt.Errorf("节点`%s`不存在", s.Name)
	}

	return global.App.DB.Transaction(func(tx *gorm.DB) error {
		// 删除当前节点
		if err := tx.Where("`key`=? and `name`=?", s.Key, s.Name).
			Delete(&models.InsightOrganizations{}).Error; err != nil {
			global.App.Log.Error(err)
			return err
		}
		// 删除所有的子节点
		if err := tx.Where("`key` like ? and `parent_id`=?", s.Key+"-%", organization.ID).
			Delete(&models.InsightOrganizations{}).Error; err != nil {
			global.App.Log.Error(err)
			return err
		}
		return nil
	})
}

type GetOrganizationsUsersServices struct {
	*forms.GetOrganizationsUsersForm
	C *gin.Context
}

func (s *GetOrganizationsUsersServices) Run() (responseData interface{}, total int64, err error) {
	type user struct {
		models.InsightUsers
		OrganizationKey  string `json:"organization_key"`
		OrganizationName string `json:"organization_name"`
	}

	var users []user

	tx := global.App.DB.Table("insight_users a").
		Select(`
            a.uid,
            a.username,
            ifnull(
                concat(
                    (
                        SELECT GROUP_CONCAT(a.name ORDER BY a.name ASC SEPARATOR '/') AS concatenated_names
                        FROM insight_organizations a
                        WHERE EXISTS (
                            SELECT 1 FROM insight_organizations b
                            WHERE JSON_CONTAINS(c.path, CONCAT('\"', a.key, '\"'))
                        )
                    ),
                    '/',
                    c.name
                ),
                c.name
            ) as organization_name,
            b.organization_key as organization_key
        `).
		Joins("JOIN insight_organizations_users b ON a.uid=b.uid").
		Joins("JOIN insight_organizations c ON c.key=b.organization_key").
		Where("b.`organization_key` LIKE ?", s.Key+"%")

	// 搜索
	if s.Search != "" {
		tx = tx.Where("a.`username` like ? ", "%"+s.Search+"%")
	}
	total = pagination.Pager(&s.PaginationQ, tx, &users)
	return &users, total, nil
}

type BindOrganizationsUsersService struct {
	*forms.BindOrganizationsUsersForm
	C *gin.Context
}

func (s *BindOrganizationsUsersService) Run() error {
	// 判断节点是否存在
	var organization models.InsightOrganizations
	tx := global.App.DB.Table("insight_organizations").Where("`key`=?", s.Key).First(&organization)
	if tx.RowsAffected == 0 {
		return fmt.Errorf("节点[%s]不存在", s.Key)
	}
	// 创建记录
	for _, uid := range s.Users {
		result := global.App.DB.Create(&models.InsightOrganizationsUsers{Uid: uid, OrganizationKey: s.Key})
		if result.Error != nil {
			type user struct {
				Username         string
				OrganizationName string
			}
			var record user
			global.App.DB.Table("insight_organizations_users a").
				Select(`b.username,ifnull(
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
				) as organization_name`).
				Joins("join insight_users b on a.uid = b.uid").
				Joins("join insight_organizations c on a.organization_key = c.key").
				Where("b.uid=?", uid).Scan(&record)
			mysqlErr := result.Error.(*mysql.MySQLError)
			switch mysqlErr.Number {
			case 1062:
				return fmt.Errorf("绑定失败，当前用户%s已绑定到组织[%s]，不能重复绑定", record.Username, record.OrganizationName)
			}
			return result.Error
		}
	}
	return nil
}

type DeleteOrganizationsUsersService struct {
	*forms.DeleteOrganizationsUsersForm
	C *gin.Context
}

func (s *DeleteOrganizationsUsersService) Run() error {
	var organizationUsers models.InsightOrganizationsUsers
	result := global.App.DB.Table("insight_organizations_users").Where("`organization_key`=? and `uid`=?", s.Key, s.Uid).First(&organizationUsers)
	if result.RowsAffected == 0 {
		return errors.New("记录`%s`不存在")
	}

	return global.App.DB.Transaction(func(tx *gorm.DB) error {
		// 删除当前节点
		if err := tx.Where("`organization_key`=? and `uid`=?", s.Key, s.Uid).
			Delete(&models.InsightOrganizationsUsers{}).Error; err != nil {
			global.App.Log.Error(err)
			return err
		}
		return nil
	})
}
