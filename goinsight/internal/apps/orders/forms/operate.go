package forms

type ApproveForm struct {
	OrderID string `form:"order_id" json:"order_id" binding:"required,uuid"`
	Msg     string `form:"msg" json:"msg" binding:"max=256"`
	Status  string `form:"status" json:"status" binding:"required,oneof=pass reject"`
}

type FeedbackForm struct {
	OrderID  string `form:"order_id" json:"order_id" binding:"required,uuid"`
	Msg      string `form:"msg" json:"msg" binding:"max=256"`
	Progress string `form:"progress" json:"progress" binding:"required,oneof=执行中 已完成"`
}

type ReviewForm struct {
	OrderID string `form:"order_id" json:"order_id" binding:"required,uuid"`
	Msg     string `form:"msg" json:"msg" binding:"max=256"`
}

type CloseForm struct {
	OrderID string `form:"order_id" json:"order_id" binding:"required,uuid"`
	Msg     string `form:"msg" json:"msg" binding:"max=256"`
}
