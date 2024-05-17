/*
@Time    :   2023/08/14 17:53:11
@Author  :   xff
*/

package middleware

import (
	"errors"
	"goInsight/global"
	userModels "goInsight/internal/apps/users/models"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp/totp"
)

var identityKey = "id"

// payload
type User struct {
	UserName string
	NickName string
}

func VerifyOTPCode(username string, otp_code string) bool {
	var user userModels.InsightUsers
	global.App.DB.Model(&userModels.InsightUsers{}).Where("username=?", username).Take(&user)

	if user.OtpSecret == "" {
		return false
	}

	valid := totp.Validate(otp_code, user.OtpSecret)
	if valid {
		return true
	} else {
		return false
	}
}

func InitAuthMiddleware() (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "goInsight",
		Key:         []byte(global.App.Config.App.SECRET_KEY),
		Timeout:     24 * time.Hour,
		MaxRefresh:  24 * time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*User); ok {
				return jwt.MapClaims{
					identityKey: v.UserName,
					"nick_name": v.NickName,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &User{
				UserName: claims[identityKey].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			// middleware otp->jwt，因此这里不需要再次校验用户信息
			username := c.GetString("loginUserName")
			if len(username) == 0 {
				return nil, jwt.ErrMissingLoginValues
			}
			// 获取用户信息
			var user userModels.InsightUsers
			global.App.DB.Table("insight_users u").
				Where("u.username=?", username).
				Scan(&user)
			// 验证otp_code
			loginNeedsOTP := c.GetString("loginNeedsOTP")
			if loginNeedsOTP == "YES" {
				otpCode := c.GetString("loginOtpCode")
				if len(otpCode) == 0 {
					return nil, errors.New("otp code is empty")
				}
				if ok := VerifyOTPCode(username, otpCode); !ok {
					return nil, errors.New("otp code error")
				}
			}
			// 更新登录时间
			global.App.DB.Table("insight_users u").
				Where("u.username=?", username).
				Update("last_login", time.Now().Format("2006-01-02 15:04:05"))
			// payload
			return &User{
				UserName: user.Username,
				NickName: user.NickName,
			}, nil
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "JWT",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})

}
