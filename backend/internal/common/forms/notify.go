package forms

type NotifyWechatForm struct {
	Enable  bool   `json:"enable"`
	Webhook string `json:"webhook"`
}

type NotifyDingTalkForm struct {
	Enable   bool   `json:"enable"`
	Webhook  string `json:"webhook"`
	Keywords string `json:"keywords"`
}

type NotifyMailForm struct {
	Enable      bool   `json:"enable"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Host        string `json:"host"`
	Port        int    `json:"port"`
	HasPassword bool   `json:"has_password"`
}

type AdminUpdateNotifyConfigForm struct {
	NoticeURL string             `json:"notice_url" binding:"required"`
	Wechat    NotifyWechatForm   `json:"wechat"`
	DingTalk  NotifyDingTalkForm `json:"dingtalk"`
	Mail      NotifyMailForm     `json:"mail"`
}

type AdminTestNotifyConfigForm struct {
	Channel string `json:"channel" binding:"required,oneof=wechat dingtalk mail"`
}

type AdminNotifyConfigResponse struct {
	NoticeURL string             `json:"notice_url"`
	Wechat    NotifyWechatForm   `json:"wechat"`
	DingTalk  NotifyDingTalkForm `json:"dingtalk"`
	Mail      NotifyMailForm     `json:"mail"`
}
