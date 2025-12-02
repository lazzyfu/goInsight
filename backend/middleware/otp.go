package middleware

import (
	"errors"
	"net/http"

	"github.com/lazzyfu/goinsight/internal/global"

	userModels "github.com/lazzyfu/goinsight/internal/users/models"

	"github.com/lazzyfu/goinsight/pkg/response"

	"github.com/gin-contrib/requestid"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// login form
type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	OtpCode  string `form:"otp_code" json:"otp_code"`
}

// user info
type userInfo struct {
	username string
	password string
	user     *userModels.InsightUsers
}

func (u *userInfo) checkIfUserExist() (err error) {
	result := global.App.DB.Table("insight_users u").
		Where("u.username=?", u.username).
		Scan(&u.user)
	if result.RowsAffected == 0 {
		return errors.New("用户名或密码错误")
	}
	return nil
}

func (u *userInfo) checkUserValid() error {
	// 用户是否激活
	if !u.user.IsActive {
		return errors.New("用户未激活")
	}
	// 验证密码
	err := bcrypt.CompareHashAndPassword([]byte(u.user.Password), []byte(u.password))
	if err != nil {
		return errors.New("用户名或密码错误")
	}
	return nil
}

func (u *userInfo) checkIfUserNeedsOTP() bool {
	// 用户是否开启2FA认证
	return u.user.IsTwoFA
}

func (u *userInfo) checkIfOtpSecretNullString() bool {
	// 检查用户的otp秘钥是否为空
	return len(u.user.OtpSecret) == 0
}

func OTPMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// login json bind to loginVals
		var loginVals login
		if err := c.ShouldBind(&loginVals); err != nil {
			response.Fail(c, "missing Username or Password")
			// stop the handler chain so the jwt.LoginHandler won't run after we've written a response
			c.Abort()
			return
		}
		username := loginVals.Username
		password := loginVals.Password

		check := userInfo{username: username, password: password}
		// 检查用户是否存在
		if err := check.checkIfUserExist(); err != nil {
			global.App.Log.WithFields(logrus.Fields{"request_id": requestid.Get(c), "username": username, "error": err}).Error(err.Error())
			c.AbortWithStatusJSON(400, gin.H{"message": err.Error(), "data": nil, "request_id": requestid.Get(c)})
			return
		}
		// 验证用户有效性
		if err := check.checkUserValid(); err != nil {
			global.App.Log.WithFields(logrus.Fields{"request_id": requestid.Get(c), "username": username, "error": err}).Error(err.Error())
			c.AbortWithStatusJSON(400, gin.H{"message": err.Error(), "data": nil, "request_id": requestid.Get(c)})
			return
		}
		// 用户是否开启2FA认证
		otpCode := loginVals.OtpCode
		needsOTP := check.checkIfUserNeedsOTP()

		if needsOTP && otpCode == "" {
			var (
				code    string
				message string
				data    any
			)

			if check.checkIfOtpSecretNullString() {
				code = "4001"
				message = "需要重新绑定 OTP"
				data = gin.H{"action": "bind_otp"}
			} else {
				code = "4002"
				message = "需要输入 OTP"
				data = gin.H{"action": "input_otp"}
			}

			c.JSON(http.StatusOK, response.Response{
				RequestID: requestid.Get(c),
				Code:      code,
				Message:   message,
				Data:      data,
			})
			// stop the handler chain so the jwt.LoginHandler won't run after we've written a response
			c.Abort()
			return
		}

		c.Set("loginUserName", username)
		c.Set("loginOtpCode", loginVals.OtpCode)
		if needsOTP {
			c.Set("loginNeedsOTP", "YES")
		} else {
			c.Set("loginNeedsOTP", "NO")
		}

		c.Next()
	}
}
