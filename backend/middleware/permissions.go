package middleware

import (
	"github.com/lazzyfu/goinsight/internal/global"

	userModels "github.com/lazzyfu/goinsight/internal/users/models"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// This function checks if the user has admin permissions
func HasAdminPermission() gin.HandlerFunc {
	// This anonymous function is used to handle the request
	return func(c *gin.Context) {
		// Extract the username from the JWT claims
		username, ok := GetUserNameFromJWT(c)
		if !ok {
			c.AbortWithStatusJSON(401, gin.H{"code": 401, "msg": "认证信息无效", "data": nil, "request_id": requestid.Get(c)})
			return
		}
		// Declare a variable to store the user
		var user userModels.InsightUsers

		// Query the database for the user
		global.App.DB.Table("insight_users u").
			Where("u.username=?", username).
			Scan(&user)
		// Check if the user is a superuser
		if !user.IsSuperuser {
			global.App.Log.WithFields(logrus.Fields{"request_id": requestid.Get(c), "username": username}).Error("您不是超级管理员.")
			c.AbortWithStatusJSON(403, gin.H{"code": 403, "msg": "无权限", "data": nil, "request_id": requestid.Get(c)})
		}
		// Call the next handler in the chain
		c.Next()
	}
}

// 下面方法暂时用不到，如果您有定制需求，可以取消注释，绑定到路由即可
// HasPermission 返回一个 gin.HandlerFunc，用于检查用户是否具有指定的权限
// 使用方法：v1.POST("generate-tasks", middleware.HasPermission("ExecuteOrders"), views.GenerateTasksView)
// func HasPermission(permission string) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		// 从JWT声明中提取用户名
// 		username := jwt.ExtractClaims(c)["id"].(string)
// 		// 获取用户绑定的organization
// 		type organization struct {
// 			Key  string
// 			Path *datatypes.JSON // 此处使用指针支持返回空
// 		}
// 		var resultOrganization organization
// 		if err := global.App.DB.Select("io.key, io.path").
// 			Table("insight_org_users ou").
// 			Joins("join insight_users iu on ou.uid=iu.uid").
// 			Joins("join insight_orgs io on io.key=ou.organization_key").
// 			Where("iu.username=?", username).
// 			Scan(&resultOrganization).Error; err != nil {
// 			// 处理数据库查询错误
// 			global.App.Log.WithFields(logrus.Fields{"request_id": requestid.Get(c), "username": username, "error": err}).Error("无法获取用户的组织信息")
// 			c.AbortWithStatusJSON(500, gin.H{"code": 500, "msg": "无法获取用户的组织信息", "data": nil, "request_id": requestid.Get(c)})
// 			return
// 		}
// 		// 解析JSON数组
// 		var okeys []string
// 		if resultOrganization.Path != nil {
// 			if err := json.Unmarshal(*resultOrganization.Path, &okeys); err != nil {
// 				// 处理JSON解析错误
// 				global.App.Log.WithFields(logrus.Fields{"request_id": requestid.Get(c), "username": username, "error": err}).Error("无法解析用户的组织信息")
// 				c.AbortWithStatusJSON(500, gin.H{"code": 500, "msg": "无法解析用户的组织信息", "data": nil, "request_id": requestid.Get(c)})
// 				return
// 			}
// 		}

// 		okeys = append(okeys, resultOrganization.Key)
// 		// 查看用户的权限
// 		type count struct {
// 			Count int
// 		}
// 		var resultCount count
// 		if err := global.App.DB.Select("count(*) as count").
// 			Table("insight_permissions_organizations a").
// 			Joins("join insight_permissions b on a.permission_id=b.id").
// 			Joins("join insight_orgs c on a.organization_key=c.key").
// 			Where("c.key in ? and b.name=?", okeys, permission).Scan(&resultCount).Error; err != nil {
// 			// 处理数据库查询错误
// 			global.App.Log.WithFields(logrus.Fields{"request_id": requestid.Get(c), "username": username, "error": err}).Error("无法查询用户的权限信息")
// 			c.AbortWithStatusJSON(500, gin.H{"code": 500, "msg": "无法查询用户的权限信息", "data": nil, "request_id": requestid.Get(c)})
// 			return
// 		}

// 		if resultCount.Count == 0 {
// 			global.App.Log.WithFields(logrus.Fields{"request_id": requestid.Get(c), "username": username}).Errorf("您没有%s权限", permission)
// 			c.AbortWithStatusJSON(403, gin.H{"code": 403, "msg": fmt.Sprintf("您没有%s权限", permission), "data": nil, "request_id": requestid.Get(c)})
// 		}
// 		// 用户有权限，调用链中的下一个处理程序
// 		c.Next()
// 	}
// }
