package notifier

import (
	"errors"
	"fmt"
	"strings"

	commonModels "github.com/lazzyfu/goinsight/internal/common/models"
	"github.com/lazzyfu/goinsight/internal/global"
	"github.com/lazzyfu/goinsight/pkg/utils"
	"gorm.io/gorm"
)

const runtimeNotifyConfigKeyDefault = "default"

type RuntimeNotifyConfig struct {
	NoticeURL string
	Wechat    struct {
		Enable  bool
		Webhook string
	}
	DingTalk struct {
		Enable   bool
		Webhook  string
		Keywords string
	}
	Mail struct {
		Enable   bool
		Username string
		Password string
		Host     string
		Port     int
	}
}

func LoadRuntimeNotifyConfig() (*RuntimeNotifyConfig, error) {
	if global.App.DB == nil {
		return nil, errors.New("database not initialized")
	}

	var settings commonModels.InsightNotifySettings
	if err := global.App.DB.Where("config_key = ?", runtimeNotifyConfigKeyDefault).Take(&settings).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("notify settings (%s) not found", runtimeNotifyConfigKeyDefault)
		}
		return nil, err
	}

	decryptedMailPassword := ""
	if strings.TrimSpace(settings.MailPassword) != "" {
		password, err := utils.Decrypt(settings.MailPassword)
		if err != nil {
			return nil, fmt.Errorf("decrypt notify mail password failed: %w", err)
		}
		decryptedMailPassword = password
	}

	cfg := &RuntimeNotifyConfig{
		NoticeURL: normalizeNoticeURL(settings.NoticeURL),
	}
	cfg.Wechat.Enable = settings.WechatEnable
	cfg.Wechat.Webhook = strings.TrimSpace(settings.WechatWebhook)
	cfg.DingTalk.Enable = settings.DingTalkEnable
	cfg.DingTalk.Webhook = strings.TrimSpace(settings.DingTalkWebhook)
	cfg.DingTalk.Keywords = strings.TrimSpace(settings.DingTalkKeywords)
	cfg.Mail.Enable = settings.MailEnable
	cfg.Mail.Username = strings.TrimSpace(settings.MailUsername)
	cfg.Mail.Password = decryptedMailPassword
	cfg.Mail.Host = strings.TrimSpace(settings.MailHost)
	cfg.Mail.Port = settings.MailPort

	return cfg, nil
}

func LoadNoticeURL() (string, error) {
	cfg, err := LoadRuntimeNotifyConfig()
	if err != nil {
		return "", err
	}
	return cfg.NoticeURL, nil
}

func normalizeNoticeURL(raw string) string {
	return strings.TrimRight(strings.TrimSpace(raw), "/")
}
