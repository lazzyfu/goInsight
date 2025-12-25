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
	OrderID    string `form:"order_id" json:"order_id" binding:"required,uuid"`
	NewClaimer string `form:"new_claimer" json:"new_claimer" binding:"required,max=32"`
	Msg        string `form:"msg" json:"msg" binding:"max=256"`
}
type RevokeOrderForm struct {
	OrderID string `form:"order_id" json:"order_id" binding:"required,uuid"`
	Msg     string `form:"msg" json:"msg" binding:"max=256"`
}

type CompleteOrderForm struct {
	OrderID string `form:"order_id" json:"order_id" binding:"required,uuid"`
	Msg     string `form:"msg" json:"msg" binding:"max=256"`
}

type FailOrderForm struct {
	OrderID string `form:"order_id" json:"order_id" binding:"required,uuid"`
	Msg     string `form:"msg" json:"msg" binding:"max=256"`
}

type ReviewOrderForm struct {
	OrderID string `form:"order_id" json:"order_id" binding:"required,uuid"`
	Msg     string `form:"msg" json:"msg" binding:"max=256"`
}
