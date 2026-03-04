package notifier

import (
	"fmt"
	"strings"
	"testing"

	commonModels "github.com/lazzyfu/goinsight/internal/common/models"
	"github.com/lazzyfu/goinsight/internal/global"
	"github.com/lazzyfu/goinsight/pkg/utils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupNotifySettingsTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", strings.ReplaceAll(t.Name(), "/", "_"))
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite db failed: %v", err)
	}
	if err := db.AutoMigrate(&commonModels.InsightNotifySettings{}); err != nil {
		t.Fatalf("migrate notify settings table failed: %v", err)
	}

	oldDB := global.App.DB
	oldSecretKey := global.App.Config.App.SECRET_KEY
	global.App.DB = db
	global.App.Config.App.SECRET_KEY = "12345678901234567890123456789012"
	t.Cleanup(func() {
		global.App.DB = oldDB
		global.App.Config.App.SECRET_KEY = oldSecretKey
	})

	return db
}

func TestLoadRuntimeNotifyConfig(t *testing.T) {
	db := setupNotifySettingsTestDB(t)

	encryptedPassword, err := utils.Encrypt("mail-secret")
	if err != nil {
		t.Fatalf("encrypt password failed: %v", err)
	}

	record := commonModels.InsightNotifySettings{
		ConfigKey:        runtimeNotifyConfigKeyDefault,
		NoticeURL:        "https://goinsight.example.com///",
		WechatEnable:     true,
		WechatWebhook:    " https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=demo ",
		DingTalkEnable:   true,
		DingTalkWebhook:  " https://oapi.dingtalk.com/robot/send?access_token=demo ",
		DingTalkKeywords: " goinsight ",
		MailEnable:       true,
		MailUsername:     " ops@example.com ",
		MailPassword:     encryptedPassword,
		MailHost:         " smtp.example.com ",
		MailPort:         465,
	}
	if err := db.Create(&record).Error; err != nil {
		t.Fatalf("insert notify settings failed: %v", err)
	}

	cfg, err := LoadRuntimeNotifyConfig()
	if err != nil {
		t.Fatalf("load runtime notify config failed: %v", err)
	}

	if cfg.NoticeURL != "https://goinsight.example.com" {
		t.Fatalf("unexpected notice url: %s", cfg.NoticeURL)
	}
	if !cfg.Wechat.Enable || !strings.Contains(cfg.Wechat.Webhook, "qyapi.weixin.qq.com") {
		t.Fatalf("unexpected wechat config: %#v", cfg.Wechat)
	}
	if !cfg.DingTalk.Enable || cfg.DingTalk.Keywords != "goinsight" {
		t.Fatalf("unexpected dingtalk config: %#v", cfg.DingTalk)
	}
	if !cfg.Mail.Enable || cfg.Mail.Password != "mail-secret" {
		t.Fatalf("unexpected mail config: %#v", cfg.Mail)
	}
}

func TestLoadRuntimeNotifyConfigNotFound(t *testing.T) {
	setupNotifySettingsTestDB(t)

	_, err := LoadRuntimeNotifyConfig()
	if err == nil {
		t.Fatal("expected error when notify settings not found")
	}
}

func TestLoadNoticeURL(t *testing.T) {
	db := setupNotifySettingsTestDB(t)
	record := commonModels.InsightNotifySettings{
		ConfigKey: runtimeNotifyConfigKeyDefault,
		NoticeURL: "https://example.com/",
	}
	if err := db.Create(&record).Error; err != nil {
		t.Fatalf("insert notify settings failed: %v", err)
	}

	noticeURL, err := LoadNoticeURL()
	if err != nil {
		t.Fatalf("load notice url failed: %v", err)
	}
	if noticeURL != "https://example.com" {
		t.Fatalf("unexpected notice url: %s", noticeURL)
	}
}
