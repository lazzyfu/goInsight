package models

// 系统消息通知配置（单例）
type InsightNotifySettings struct {
	*Model
	ConfigKey string `gorm:"type:varchar(32);not null;default:default;uniqueIndex:uniq_notify_config_key;comment:配置唯一键" json:"config_key"`

	NoticeURL string `gorm:"type:varchar(512);not null;default:'';comment:通知链接前缀" json:"notice_url"`

	WechatEnable  bool   `gorm:"type:boolean;not null;default:false;comment:是否启用企业微信通知" json:"wechat_enable"`
	WechatWebhook string `gorm:"type:varchar(512);not null;default:'';comment:企业微信Webhook" json:"wechat_webhook"`

	DingTalkEnable   bool   `gorm:"type:boolean;not null;default:false;comment:是否启用钉钉通知" json:"dingtalk_enable"`
	DingTalkWebhook  string `gorm:"type:varchar(512);not null;default:'';comment:钉钉Webhook" json:"dingtalk_webhook"`
	DingTalkKeywords string `gorm:"type:varchar(256);not null;default:'';comment:钉钉关键字" json:"dingtalk_keywords"`

	MailEnable   bool   `gorm:"type:boolean;not null;default:false;comment:是否启用邮件通知" json:"mail_enable"`
	MailUsername string `gorm:"type:varchar(254);not null;default:'';comment:邮件用户名" json:"mail_username"`
	MailPassword string `gorm:"type:varchar(1024);not null;default:'';comment:邮件密码(加密)" json:"mail_password"`
	MailHost     string `gorm:"type:varchar(255);not null;default:'';comment:邮件主机" json:"mail_host"`
	MailPort     int    `gorm:"type:int;not null;default:0;comment:邮件端口" json:"mail_port"`
}

func (InsightNotifySettings) TableName() string {
	return "insight_notify_settings"
}
