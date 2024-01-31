package forms

type UpdateUserInfoForm struct {
	NickName string `form:"nick_name"  json:"nick_name"`
	Mobile   string `form:"mobile" json:"mobile"`
	Email    string `form:"email" json:"email"`
}

type ChangeUserPasswordForm struct {
	CurrentPassword string `form:"current_password" json:"current_password" binding:"required,min=7,max=32"`
	NewPassword     string `form:"new_password" json:"new_password" binding:"required,min=7,max=32"`
	VerifyPassword  string `form:"verify_password" json:"verify_password" binding:"required,min=7,max=32"`
}

type GetOTPAuthURLForm struct {
	Username string `form:"username"  json:"username" binding:"required,min=3,max=32"`
	Password string `form:"password" json:"password" binding:"required,min=3,max=128"`
}

type GetOTPAuthCallbackForm struct {
	Username string `form:"username"  json:"username" binding:"required,min=3,max=32"`
	Password string `form:"password" json:"password" binding:"required,min=3,max=128"`
	Callback string `form:"callback" json:"callback" binding:"required"`
}
