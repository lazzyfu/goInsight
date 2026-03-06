package services

import (
	"errors"
	"fmt"
	netmail "net/mail"
	"net/url"
	"strings"
	"time"

	"github.com/lazzyfu/goinsight/internal/common/forms"
	commonModels "github.com/lazzyfu/goinsight/internal/common/models"
	"github.com/lazzyfu/goinsight/internal/global"
	userModels "github.com/lazzyfu/goinsight/internal/users/models"
	"github.com/lazzyfu/goinsight/middleware"
	"github.com/lazzyfu/goinsight/pkg/notifier"
	"github.com/lazzyfu/goinsight/pkg/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const notifyConfigKeyDefault = "default"

type AdminGetNotifyConfigService struct {
	C *gin.Context
}

func (s *AdminGetNotifyConfigService) Run() (responseData any, err error) {
	record, err := loadOrCreateNotifySettings(global.App.DB)
	if err != nil {
		return nil, err
	}
	return buildNotifyConfigResponse(record), nil
}

type AdminUpdateNotifyConfigService struct {
	*forms.AdminUpdateNotifyConfigForm
	C *gin.Context
}

func (s *AdminUpdateNotifyConfigService) Run() error {
	record, err := loadOrCreateNotifySettings(global.App.DB)
	if err != nil {
		return err
	}

	hasStoredMailPassword := strings.TrimSpace(record.MailPassword) != ""
	if err := validateNotifyConfigForSave(s.AdminUpdateNotifyConfigForm, hasStoredMailPassword); err != nil {
		return err
	}

	noticeURL := normalizeNoticeURL(s.NoticeURL)
	wechatWebhook := strings.TrimSpace(s.Wechat.Webhook)
	dingtalkWebhook := strings.TrimSpace(s.DingTalk.Webhook)
	dingtalkKeywords := strings.TrimSpace(s.DingTalk.Keywords)
	mailUsername := strings.TrimSpace(s.Mail.Username)
	mailHost := strings.TrimSpace(s.Mail.Host)

	encryptedMailPassword := record.MailPassword
	inputMailPassword := strings.TrimSpace(s.Mail.Password)
	if inputMailPassword != "" {
		encryptedMailPassword, err = utils.Encrypt(inputMailPassword)
		if err != nil {
			return fmt.Errorf("邮件密码加密失败: %w", err)
		}
	}

	updates := map[string]any{
		"notice_url":         noticeURL,
		"wechat_enable":      s.Wechat.Enable,
		"wechat_webhook":     wechatWebhook,
		"ding_talk_enable":   s.DingTalk.Enable,
		"ding_talk_webhook":  dingtalkWebhook,
		"ding_talk_keywords": dingtalkKeywords,
		"mail_enable":        s.Mail.Enable,
		"mail_username":      mailUsername,
		"mail_password":      encryptedMailPassword,
		"mail_host":          mailHost,
		"mail_port":          s.Mail.Port,
	}

	if err := global.App.DB.Model(&commonModels.InsightNotifySettings{}).Where("id = ?", record.ID).Updates(updates).Error; err != nil {
		return err
	}

	record.NoticeURL = noticeURL
	record.WechatEnable = s.Wechat.Enable
	record.WechatWebhook = wechatWebhook
	record.DingTalkEnable = s.DingTalk.Enable
	record.DingTalkWebhook = dingtalkWebhook
	record.DingTalkKeywords = dingtalkKeywords
	record.MailEnable = s.Mail.Enable
	record.MailUsername = mailUsername
	record.MailPassword = encryptedMailPassword
	record.MailHost = mailHost
	record.MailPort = s.Mail.Port

	return nil
}

type AdminTestNotifyConfigService struct {
	*forms.AdminTestNotifyConfigForm
	C *gin.Context
}

func (s *AdminTestNotifyConfigService) Run() error {
	username, ok := middleware.GetUserNameFromJWT(s.C)
	if !ok {
		return errors.New("获取当前登录用户失败")
	}

	runtimeCfg, err := notifier.LoadRuntimeNotifyConfig()
	if err != nil {
		return fmt.Errorf("加载通知配置失败: %w", err)
	}
	if err := validateRuntimeChannelConfig(s.Channel, runtimeCfg); err != nil {
		return err
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	baseMsg := fmt.Sprintf("您好，这是一条 GoInsight 消息通知测试。\n时间：%s", now)

	switch s.Channel {
	case "wechat":
		wechatNotifier := notifier.NewWechatNotifier(notifier.WechatNotifierConfig{
			Webhook: runtimeCfg.Wechat.Webhook,
		})
		if err := wechatNotifier.SendMessage("", []string{username}, baseMsg); err != nil {
			return fmt.Errorf("企业微信测试发送失败: %w", err)
		}
		return nil

	case "dingtalk":
		dingtalkNotifier := notifier.NewDingTalkNotifier(notifier.DingTalkNotifierConfig{
			Webhook: runtimeCfg.DingTalk.Webhook,
		})
		messageWithKeywords := baseMsg
		if strings.TrimSpace(runtimeCfg.DingTalk.Keywords) != "" {
			messageWithKeywords = fmt.Sprintf("%s\nKeywords：%s", baseMsg, runtimeCfg.DingTalk.Keywords)
		}
		if err := dingtalkNotifier.SendMessage("", []string{username}, messageWithKeywords); err != nil {
			return fmt.Errorf("钉钉测试发送失败: %w", err)
		}
		return nil

	case "mail":
		var user userModels.InsightUsers
		if err := global.App.DB.Table("insight_users").Where("username = ?", username).Take(&user).Error; err != nil {
			return fmt.Errorf("邮件测试发送失败: 查询用户邮箱失败(%w)", err)
		}
		if strings.TrimSpace(user.Email) == "" {
			return errors.New("邮件测试发送失败: 当前用户未配置邮箱")
		}

		emailNotifier := notifier.NewEmailNotifier(notifier.EmailNotifierConfig{
			Host:     runtimeCfg.Mail.Host,
			Port:     runtimeCfg.Mail.Port,
			Username: runtimeCfg.Mail.Username,
			Password: runtimeCfg.Mail.Password,
		})
		if err := emailNotifier.SendMessage("【GoInsight】消息通知测试", []string{user.Email}, baseMsg); err != nil {
			return fmt.Errorf("邮件测试发送失败: %w", err)
		}
		return nil

	default:
		return errors.New("不支持的测试渠道")
	}
}

func validateRuntimeChannelConfig(channel string, runtimeCfg *notifier.RuntimeNotifyConfig) error {
	if runtimeCfg == nil {
		return errors.New("通知配置不存在")
	}

	switch channel {
	case "wechat":
		if !runtimeCfg.Wechat.Enable {
			return errors.New("请先启用企业微信通知后再测试")
		}
		if !isValidHTTPURL(runtimeCfg.Wechat.Webhook) {
			return errors.New("企业微信 Webhook 配置不正确，请先保存合法地址")
		}
		return nil

	case "dingtalk":
		if !runtimeCfg.DingTalk.Enable {
			return errors.New("请先启用钉钉通知后再测试")
		}
		if !isValidHTTPURL(runtimeCfg.DingTalk.Webhook) {
			return errors.New("钉钉 Webhook 配置不正确，请先保存合法地址")
		}
		if strings.TrimSpace(runtimeCfg.DingTalk.Keywords) == "" {
			return errors.New("钉钉关键字不能为空，请先完善后再测试")
		}
		return nil

	case "mail":
		if !runtimeCfg.Mail.Enable {
			return errors.New("请先启用邮件通知后再测试")
		}
		mailUsername := strings.TrimSpace(runtimeCfg.Mail.Username)
		if mailUsername == "" {
			return errors.New("邮件发件账号不能为空，请先完善后再测试")
		}
		if _, err := netmail.ParseAddress(mailUsername); err != nil {
			return errors.New("邮件发件账号格式不正确，请先完善后再测试")
		}
		if strings.TrimSpace(runtimeCfg.Mail.Host) == "" {
			return errors.New("邮件 SMTP 主机不能为空，请先完善后再测试")
		}
		if runtimeCfg.Mail.Port < 1 || runtimeCfg.Mail.Port > 65535 {
			return errors.New("邮件 SMTP 端口范围必须为 1-65535")
		}
		if strings.TrimSpace(runtimeCfg.Mail.Password) == "" {
			return errors.New("邮件 SMTP 密码不能为空，请先完善后再测试")
		}
		return nil

	default:
		return errors.New("不支持的测试渠道")
	}
}

func normalizeNoticeURL(raw string) string {
	return strings.TrimRight(strings.TrimSpace(raw), "/")
}

func validateNotifyConfigForSave(form *forms.AdminUpdateNotifyConfigForm, hasStoredMailPassword bool) error {
	if form == nil {
		return errors.New("通知配置不能为空")
	}

	form.NoticeURL = normalizeNoticeURL(form.NoticeURL)
	if form.NoticeURL == "" {
		return errors.New("通知地址不能为空")
	}
	if !isValidHTTPURL(form.NoticeURL) {
		return errors.New("通知地址格式不正确，仅支持 http/https")
	}

	form.Wechat.Webhook = strings.TrimSpace(form.Wechat.Webhook)
	form.DingTalk.Webhook = strings.TrimSpace(form.DingTalk.Webhook)
	form.DingTalk.Keywords = strings.TrimSpace(form.DingTalk.Keywords)
	form.Mail.Username = strings.TrimSpace(form.Mail.Username)
	form.Mail.Password = strings.TrimSpace(form.Mail.Password)
	form.Mail.Host = strings.TrimSpace(form.Mail.Host)

	if form.Wechat.Enable && form.Wechat.Webhook == "" {
		return errors.New("启用企业微信通知时，Webhook 不能为空")
	}
	if form.Wechat.Enable && !isValidHTTPURL(form.Wechat.Webhook) {
		return errors.New("启用企业微信通知时，Webhook 地址格式不正确")
	}

	if form.DingTalk.Enable && form.DingTalk.Webhook == "" {
		return errors.New("启用钉钉通知时，Webhook 不能为空")
	}
	if form.DingTalk.Enable && !isValidHTTPURL(form.DingTalk.Webhook) {
		return errors.New("启用钉钉通知时，Webhook 地址格式不正确")
	}
	if form.DingTalk.Enable && form.DingTalk.Keywords == "" {
		return errors.New("启用钉钉通知时，关键字不能为空")
	}

	if form.Mail.Enable {
		if form.Mail.Username == "" {
			return errors.New("启用邮件通知时，用户名不能为空")
		}
		if _, err := netmail.ParseAddress(form.Mail.Username); err != nil {
			return errors.New("启用邮件通知时，用户名必须是合法邮箱")
		}
		if form.Mail.Host == "" {
			return errors.New("启用邮件通知时，SMTP 主机不能为空")
		}
		if form.Mail.Port < 1 || form.Mail.Port > 65535 {
			return errors.New("启用邮件通知时，SMTP 端口范围必须为 1-65535")
		}
		if form.Mail.Password == "" && !hasStoredMailPassword {
			return errors.New("启用邮件通知时，SMTP 密码不能为空")
		}
	}

	return nil
}

func isValidHTTPURL(raw string) bool {
	u, err := url.Parse(raw)
	if err != nil {
		return false
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}
	return strings.TrimSpace(u.Host) != ""
}

func buildNotifyConfigResponse(record *commonModels.InsightNotifySettings) *forms.AdminNotifyConfigResponse {
	return &forms.AdminNotifyConfigResponse{
		NoticeURL: normalizeNoticeURL(record.NoticeURL),
		Wechat: forms.NotifyWechatForm{
			Enable:  record.WechatEnable,
			Webhook: record.WechatWebhook,
		},
		DingTalk: forms.NotifyDingTalkForm{
			Enable:   record.DingTalkEnable,
			Webhook:  record.DingTalkWebhook,
			Keywords: record.DingTalkKeywords,
		},
		Mail: forms.NotifyMailForm{
			Enable:      record.MailEnable,
			Username:    record.MailUsername,
			Password:    "",
			Host:        record.MailHost,
			Port:        record.MailPort,
			HasPassword: strings.TrimSpace(record.MailPassword) != "",
		},
	}
}

func loadOrCreateNotifySettings(db *gorm.DB) (*commonModels.InsightNotifySettings, error) {
	var record commonModels.InsightNotifySettings
	err := db.Where("config_key = ?", notifyConfigKeyDefault).First(&record).Error
	if err == nil {
		return &record, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	seed := buildDefaultNotifySettings()
	if err := db.Create(seed).Error; err != nil {
		return nil, err
	}
	return seed, nil
}

func buildDefaultNotifySettings() *commonModels.InsightNotifySettings {
	return &commonModels.InsightNotifySettings{
		ConfigKey:        notifyConfigKeyDefault,
		NoticeURL:        defaultNoticeURLFromAppConfig(),
		WechatEnable:     false,
		WechatWebhook:    "",
		DingTalkEnable:   false,
		DingTalkWebhook:  "",
		DingTalkKeywords: "",
		MailEnable:       false,
		MailUsername:     "",
		MailPassword:     "",
		MailHost:         "",
		MailPort:         465,
	}
}

func defaultNoticeURLFromAppConfig() string {
	address := strings.TrimSpace(global.App.Config.App.ListenAddress)
	if address == "" {
		return ""
	}
	if strings.HasPrefix(address, "http://") || strings.HasPrefix(address, "https://") {
		return strings.TrimRight(address, "/")
	}
	return "http://" + strings.TrimRight(address, "/")
}
