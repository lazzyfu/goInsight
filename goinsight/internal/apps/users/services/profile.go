package services

import (
	"fmt"
	"goInsight/global"
	"goInsight/internal/apps/users/forms"
	"goInsight/internal/apps/users/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
)

type GetUserInfoServices struct {
	C        *gin.Context
	Username string
}

func (s *GetUserInfoServices) Run() (responseData interface{}, err error) {
	var user map[string]interface{}
	tx := global.App.DB.Table("insight_users a").
		Select("a.*, ifnull(c.name, '-/-') as organization, ifnull(d.name, '-/-') as role").
		Joins("left join insight_organizations_users b on a.uid=b.uid").
		Joins("left join insight_organizations c on b.organization_key=c.key").
		Joins("left join insight_roles d on d.id=a.role_id").
		Where("a.username=?", s.Username).
		Scan(&user)
	if tx.RowsAffected == 0 {
		return user, fmt.Errorf("用户`%s`不存在", s.Username)
	}
	return user, err
}

type UpdateUserInfoService struct {
	*forms.UpdateUserInfoForm
	C   *gin.Context
	UID uint32
}

func (s *UpdateUserInfoService) Run() error {
	tx := global.App.DB.Model(&models.InsightUsers{}).Where("uid=?", s.UID)
	data := make(map[string]interface{})
	if s.NickName != "" {
		data["NickName"] = s.NickName
	}
	if s.Mobile != "" {
		data["Mobile"] = s.Mobile
	}
	if s.Email != "" {
		data["Email"] = s.Email
	}
	tx.Updates(data)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

// 用户更改密码
type ChangeUserPasswordService struct {
	*forms.ChangeUserPasswordForm
	C        *gin.Context
	Username string
}

func (s *ChangeUserPasswordService) Run() error {
	var user models.InsightUsers
	global.App.DB.Model(&models.InsightUsers{}).Where("username=?", s.Username).Scan(&user)
	// 验证老密码是否正确
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(s.CurrentPassword))
	if err != nil {
		return fmt.Errorf("密码更改失败，旧密码输入不正确")
	}
	// 验证新老密码是否一致
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(s.NewPassword))
	if err == nil {
		return fmt.Errorf("密码更改失败，新旧密码一致")
	}
	// 加密密码
	hashedPassword := models.BcryptPW(s.NewPassword)
	global.App.DB.Model(&models.InsightUsers{}).Where("username=?", s.Username).Update("password", hashedPassword)
	return nil
}

type GetOTPAuthURLService struct {
	*forms.GetOTPAuthURLForm
	C *gin.Context
}

func (s *GetOTPAuthURLService) Run() (data interface{}, err error) {
	var user models.InsightUsers
	result := global.App.DB.Model(&models.InsightUsers{}).Where("username=?", s.Username).Scan(&user)
	if result.RowsAffected == 0 {
		return data, fmt.Errorf("用户名或密码不正确")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(s.Password))
	if err != nil {
		return data, fmt.Errorf("用户名或密码不正确")
	}
	if user.OtpSecret != "" {
		return data, fmt.Errorf("不能重复绑定，如需重新绑定，请联系系统管理员")
	}
	secret, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "goInsight",
		AccountName: s.Username, // You can customize the account name
	})
	if err != nil {
		return data, fmt.Errorf("failed to generate OTP secret")
	}
	uuid := uuid.New().String()
	global.App.Redis.Set(s.C.Request.Context(), uuid, secret.Secret(), 10*time.Minute)
	data = map[string]string{"otpAuthUrl": secret.URL(), "callback": uuid}
	return data, nil
}

type GetOTPAuthCallbackService struct {
	*forms.GetOTPAuthCallbackForm
	C *gin.Context
}

func (s *GetOTPAuthCallbackService) Run() error {
	var user models.InsightUsers
	result := global.App.DB.Model(&models.InsightUsers{}).Where("username=?", s.Username).Scan(&user)
	if result.RowsAffected == 0 {
		return fmt.Errorf("用户名或密码不正确")
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(s.Password))
	if err != nil {
		return fmt.Errorf("用户名或密码不正确")
	}
	secret, err := global.App.Redis.Get(s.C.Request.Context(), s.Callback).Result()
	if err != nil {
		return err
	}
	global.App.DB.Model(&models.InsightUsers{}).Where("username=?", s.Username).Update("otp_secret", secret)
	return nil
}
