package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"goInsight/global"
	"goInsight/internal/users/models"
	"goInsight/pkg/utils"
	"net/http"

	"gopkg.in/gomail.v2"
)

// WechatNotifierConfig holds configuration for the WeChat notifier
type WechatNotifierConfig struct {
	Webhook string
}

// DingTalkNotifierConfig holds configuration for the DingTalk notifier
type DingTalkNotifierConfig struct {
	Webhook string
}

// EmailNotifierConfig holds configuration for the email notifier
type EmailNotifierConfig struct {
	Username string
	Host     string
	Port     int
	Password string
}

// Notifier interface defines methods for sending messages
type Notifier interface {
	SendMessage(subject string, users []string, msg string) error
}

// WechatNotifier implements Notifier for sending WeChat messages
type WechatNotifier struct {
	Config WechatNotifierConfig
}

// SendMessage 发送企业微信消息
func (w *WechatNotifier) SendMessage(subject string, users []string, msg string) error {
	payload := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]interface{}{
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

// DingTalkNotifier implements Notifier for sending DingTalk messages
type DingTalkNotifier struct {
	Config DingTalkNotifierConfig
}

// SendMessage sends a DingTalk message
func (d *DingTalkNotifier) SendMessage(subject string, users []string, msg string) error {
	// Construct the DingTalk message payload
	payload := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": msg,
		},
		"at": map[string]interface{}{
			"atMobiles": users,
			"isAtAll":   false,
		},
	}

	// Convert the payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// Send the HTTP POST request to the DingTalk webhook
	resp, err := http.Post(d.Config.Webhook, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode == http.StatusOK {
		global.App.Log.Info("DingTalk message sent successfully to users:", users)
	} else {
		return fmt.Errorf("failed to send DingTalk message to users: %v, Status code: %d", users, resp.StatusCode)
	}

	return nil
}

// EmailNotifier implements Notifier for sending emails
type EmailNotifier struct {
	Config EmailNotifierConfig
}

func (e *EmailNotifier) SendMessage(subject string, emails []string, msg string) error {
	// Create a new message
	m := gomail.NewMessage()

	// Set the sender (you should replace "your-email@example.com" with the actual sender's email address)
	m.SetHeader("From", e.Config.Username)

	// Set the recipients
	m.SetHeader("To", emails...)

	// Set the email subject
	m.SetHeader("Subject", subject)

	// Set the email body (text/plain)
	m.SetBody("text/plain", msg)

	// Create a new dialer (you should replace the SMTP server and credentials with your own)
	d := gomail.NewDialer(e.Config.Host, e.Config.Port, e.Config.Username, e.Config.Password)

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

// NewWechatNotifier creates a new WeChat notifier instance
func NewWechatNotifier(config WechatNotifierConfig) Notifier {
	return &WechatNotifier{Config: config}
}

// NewDingTalkNotifier creates a new DingTalk notifier instance
func NewDingTalkNotifier(config DingTalkNotifierConfig) Notifier {
	return &DingTalkNotifier{Config: config}
}

// NewEmailNotifier creates a new email notifier instance
func NewEmailNotifier(config EmailNotifierConfig) Notifier {
	return &EmailNotifier{Config: config}
}

// SendMessage sends a notification message to users via WeChat and Email
func SendMessage(subject, orderID string, users []string, msg string) {
	// Add the order link to the notification message
	noticeURL := fmt.Sprintf("%s/orders/detail/%s", global.App.Config.Notify.NoticeURL, orderID)
	msg = fmt.Sprintf("%s\n\n工单地址：%s", msg, noticeURL)

	// Deduplicate user list
	newUsers := utils.RemoveDuplicate(users)

	// Send WeChat message
	if global.App.Config.Notify.Wechat.Enable {
		wechatConfig := WechatNotifierConfig{Webhook: global.App.Config.Notify.Wechat.Webhook}
		wechatNotifier := NewWechatNotifier(wechatConfig)
		if err := wechatNotifier.SendMessage("", newUsers, msg); err != nil {
			global.App.Log.Error("Error sending WeChat message:", err)
		}
	}

	// Send DingTalk message
	if global.App.Config.Notify.DingTalk.Enable {
		dingTalkConfig := DingTalkNotifierConfig{Webhook: global.App.Config.Notify.DingTalk.Webhook}
		withKeywordsMsg := fmt.Sprintf("%s\nKeywords：%s", msg, global.App.Config.Notify.DingTalk.Keywords)
		dingTalkNotifier := NewDingTalkNotifier(dingTalkConfig)
		if err := dingTalkNotifier.SendMessage("", newUsers, withKeywordsMsg); err != nil {
			global.App.Log.Error("Error sending DingTalk message:", err)
		}
	}

	// Send Email message
	if global.App.Config.Notify.Mail.Enable {
		emailConfig := EmailNotifierConfig{
			Host:     global.App.Config.Notify.Mail.Host,
			Port:     global.App.Config.Notify.Mail.Port,
			Username: global.App.Config.Notify.Mail.Username,
			Password: global.App.Config.Notify.Mail.Password,
		}
		emailNotifier := NewEmailNotifier(emailConfig)

		// Retrieve user emails from the database
		var usersWithEmail []string
		var usersFromDB []models.InsightUsers
		global.App.DB.Table("insight_users").Where("username in ?", newUsers).Scan(&usersFromDB)

		for _, user := range usersFromDB {
			usersWithEmail = append(usersWithEmail, user.Email)
		}

		// Send email to each user
		for _, userEmail := range usersWithEmail {
			if err := emailNotifier.SendMessage("【工单】"+subject, []string{userEmail}, msg); err != nil {
				global.App.Log.Error("Error sending email:", err)
			}
		}
	}
}
