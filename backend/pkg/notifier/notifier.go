package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/lazzyfu/goinsight/internal/global"

	"github.com/lazzyfu/goinsight/pkg/utils"

	"github.com/lazzyfu/goinsight/internal/users/models"

	"gopkg.in/gomail.v2"
)

// WechatNotifierConfig 保存企业微信通知器的配置
type WechatNotifierConfig struct {
	Webhook string
}

// DingTalkNotifierConfig 保存钉钉通知器的配置
type DingTalkNotifierConfig struct {
	Webhook string
}

// EmailNotifierConfig 保存邮件通知器的配置
type EmailNotifierConfig struct {
	Username string
	Host     string
	Port     int
	Password string
}

// Notifier 接口定义发送消息的方法
type Notifier interface {
	SendMessage(subject string, users []string, msg string) error
}

// WechatNotifier 实现了 Notifier 用于发送企业微信消息
type WechatNotifier struct {
	Config WechatNotifierConfig
}

// SendMessage 发送企业微信消息
func (w *WechatNotifier) SendMessage(subject string, users []string, msg string) error {
	payload := map[string]any{
		"msgtype": "markdown",
		"markdown": map[string]any{
			"content":        msg,
			"mentioned_list": users,
		},
	}
	messageJSON, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	resp, err := http.Post(w.Config.Webhook, "application/json", bytes.NewBuffer(messageJSON))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		global.App.Log.Info("Message sent successfully to users:", users)
	} else {
		return fmt.Errorf("failed to send message to users: %v, Status code: %d", users, resp.StatusCode)
	}

	return nil
}

// DingTalkNotifier 实现了 Notifier 用于发送钉钉消息
type DingTalkNotifier struct {
	Config DingTalkNotifierConfig
}

// SendMessage 发送钉钉消息
func (d *DingTalkNotifier) SendMessage(subject string, users []string, msg string) error {
	// 构造钉钉消息负载
	payload := map[string]any{
		"msgtype": "text",
		"text": map[string]string{
			"content": msg,
		},
		"at": map[string]any{
			"atMobiles": users,
			"isAtAll":   false,
		},
	}

	// 将负载转换为 JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// 发送 HTTP POST 请求到钉钉 webhook
	resp, err := http.Post(d.Config.Webhook, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode == http.StatusOK {
		global.App.Log.Info("DingTalk message sent successfully to users:", users)
	} else {
		return fmt.Errorf("failed to send DingTalk message to users: %v, Status code: %d", users, resp.StatusCode)
	}

	return nil
}

// EmailNotifier 实现了 Notifier 用于发送邮件
type EmailNotifier struct {
	Config EmailNotifierConfig
}

func (e *EmailNotifier) SendMessage(subject string, emails []string, msg string) error {
	// 创建新的邮件消息
	m := gomail.NewMessage()

	// 设置发件人（请替换为实际的发件人邮箱）
	m.SetHeader("From", e.Config.Username)

	// 设置收件人
	m.SetHeader("To", emails...)

	// 设置邮件主题
	m.SetHeader("Subject", subject)

	// 设置邮件正文（文本）
	m.SetBody("text/plain", msg)

	// 创建新的拨号器（请使用实际的 SMTP 服务器和凭据）
	d := gomail.NewDialer(e.Config.Host, e.Config.Port, e.Config.Username, e.Config.Password)

	// 发送邮件
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

// NewWechatNotifier 创建新的企业微信通知器实例
func NewWechatNotifier(config WechatNotifierConfig) Notifier {
	return &WechatNotifier{Config: config}
}

// NewDingTalkNotifier 创建新的钉钉通知器实例
func NewDingTalkNotifier(config DingTalkNotifierConfig) Notifier {
	return &DingTalkNotifier{Config: config}
}

// NewEmailNotifier 创建新的邮件通知器实例
func NewEmailNotifier(config EmailNotifierConfig) Notifier {
	return &EmailNotifier{Config: config}
}

// SendMessage 通过企业微信/钉钉/邮件发送通知消息给用户
func SendMessage(subject, orderID string, users []string, msg string) {
	// 在通知消息中添加工单链接
	noticeURL := fmt.Sprintf("%s/orders/%s", global.App.Config.Notify.NoticeURL, orderID)
	msg = fmt.Sprintf("%s\n\n工单地址：%s", msg, noticeURL)

	// 去重用户列表
	newUsers := utils.RemoveDuplicate(users)

	// 发送企业微信消息
	if global.App.Config.Notify.Wechat.Enable {
		wechatConfig := WechatNotifierConfig{Webhook: global.App.Config.Notify.Wechat.Webhook}
		wechatNotifier := NewWechatNotifier(wechatConfig)
		if err := wechatNotifier.SendMessage("", newUsers, msg); err != nil {
			global.App.Log.Error("Error sending WeChat message:", err)
		}
	}

	// 发送钉钉消息
	if global.App.Config.Notify.DingTalk.Enable {
		dingTalkConfig := DingTalkNotifierConfig{Webhook: global.App.Config.Notify.DingTalk.Webhook}
		withKeywordsMsg := fmt.Sprintf("%s\nKeywords：%s", msg, global.App.Config.Notify.DingTalk.Keywords)
		dingTalkNotifier := NewDingTalkNotifier(dingTalkConfig)
		if err := dingTalkNotifier.SendMessage("", newUsers, withKeywordsMsg); err != nil {
			global.App.Log.Error("Error sending DingTalk message:", err)
		}
	}

	// 发送邮件消息
	if global.App.Config.Notify.Mail.Enable {
		emailConfig := EmailNotifierConfig{
			Host:     global.App.Config.Notify.Mail.Host,
			Port:     global.App.Config.Notify.Mail.Port,
			Username: global.App.Config.Notify.Mail.Username,
			Password: global.App.Config.Notify.Mail.Password,
		}
		emailNotifier := NewEmailNotifier(emailConfig)

		// 从数据库中获取用户邮箱
		var usersWithEmail []string
		var usersFromDB []models.InsightUsers
		global.App.DB.Table("insight_users").Where("username in ?", newUsers).Scan(&usersFromDB)

		for _, user := range usersFromDB {
			usersWithEmail = append(usersWithEmail, user.Email)
		}

		// 给每个用户发送邮件
		for _, userEmail := range usersWithEmail {
			if err := emailNotifier.SendMessage("【工单】"+subject, []string{userEmail}, msg); err != nil {
				global.App.Log.Error("Error sending email:", err)
			}
		}
	}
}
