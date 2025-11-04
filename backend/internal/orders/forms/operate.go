package forms

type ApprovalOrderForm struct {
	OrderID string `form:"order_id" json:"order_id" binding:"required,uuid"`
	Msg     string `form:"msg" json:"msg" binding:"max=256"`
	Status  string `form:"status" json:"status" binding:"required,oneof=APPROVED REJECTED"`
}

type ClaimOrderForm struct {
	OrderID string `form:"order_id" json:"order_id" binding:"required,uuid"`
	Msg     string `form:"msg" json:"msg" binding:"max=256"`
}
type TransferOrderForm struct {
	OrderID     string `form:"order_id" json:"order_id" binding:"required,uuid"`
	NewExecutor string `form:"new_executor" json:"new_executor" binding:"required,max=32"`
	TransferMsg string `form:"transfer_msg" json:"transfer_msg" binding:"max=256"`
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

type CloseOrderForm struct {
	OrderID string `form:"order_id" json:"order_id" binding:"required,uuid"`
	Msg     string `form:"msg" json:"msg" binding:"max=256"`
}
