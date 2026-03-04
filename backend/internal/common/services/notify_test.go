package services

import (
	"testing"

	"github.com/lazzyfu/goinsight/internal/common/forms"
)

func buildValidNotifyForm() *forms.AdminUpdateNotifyConfigForm {
	return &forms.AdminUpdateNotifyConfigForm{
		NoticeURL: "https://goinsight.example.com",
		Wechat: forms.NotifyWechatForm{
			Enable:  false,
			Webhook: "",
		},
		DingTalk: forms.NotifyDingTalkForm{
			Enable:   false,
			Webhook:  "",
			Keywords: "",
		},
		Mail: forms.NotifyMailForm{
			Enable:   false,
			Username: "",
			Password: "",
			Host:     "",
			Port:     0,
		},
	}
}

func TestNormalizeNoticeURL(t *testing.T) {
	got := normalizeNoticeURL("https://goinsight.example.com///")
	if got != "https://goinsight.example.com" {
		t.Fatalf("unexpected normalize result, got=%s", got)
	}
}

func TestValidateNotifyConfigForSave(t *testing.T) {
	t.Run("valid basic config", func(t *testing.T) {
		form := buildValidNotifyForm()
		if err := validateNotifyConfigForSave(form, false); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("invalid notice_url", func(t *testing.T) {
		form := buildValidNotifyForm()
		form.NoticeURL = "ftp://invalid.example.com"
		if err := validateNotifyConfigForSave(form, false); err == nil {
			t.Fatal("expected validation error")
		}
	})

	t.Run("wechat enabled without webhook", func(t *testing.T) {
		form := buildValidNotifyForm()
		form.Wechat.Enable = true
		if err := validateNotifyConfigForSave(form, false); err == nil {
			t.Fatal("expected validation error")
		}
	})

	t.Run("dingtalk enabled without webhook", func(t *testing.T) {
		form := buildValidNotifyForm()
		form.DingTalk.Enable = true
		if err := validateNotifyConfigForSave(form, false); err == nil {
			t.Fatal("expected validation error")
		}
	})

	t.Run("dingtalk enabled without keywords", func(t *testing.T) {
		form := buildValidNotifyForm()
		form.DingTalk.Enable = true
		form.DingTalk.Webhook = "https://oapi.dingtalk.com/robot/send?access_token=xxx"
		form.DingTalk.Keywords = ""
		if err := validateNotifyConfigForSave(form, false); err == nil {
			t.Fatal("expected validation error")
		}
	})

	t.Run("dingtalk enabled with keywords", func(t *testing.T) {
		form := buildValidNotifyForm()
		form.DingTalk.Enable = true
		form.DingTalk.Webhook = "https://oapi.dingtalk.com/robot/send?access_token=xxx"
		form.DingTalk.Keywords = "GoInsight"
		if err := validateNotifyConfigForSave(form, false); err != nil {
			t.Fatalf("unexpected validation error: %v", err)
		}
	})

	t.Run("mail enabled without required fields", func(t *testing.T) {
		form := buildValidNotifyForm()
		form.Mail.Enable = true
		if err := validateNotifyConfigForSave(form, false); err == nil {
			t.Fatal("expected validation error")
		}
	})

	t.Run("mail enabled with invalid port", func(t *testing.T) {
		form := buildValidNotifyForm()
		form.Mail.Enable = true
		form.Mail.Username = "ops@example.com"
		form.Mail.Host = "smtp.example.com"
		form.Mail.Password = "secret"
		form.Mail.Port = 65536
		if err := validateNotifyConfigForSave(form, false); err == nil {
			t.Fatal("expected validation error")
		}
	})

	t.Run("mail enabled without password should fail even with existing password", func(t *testing.T) {
		form := buildValidNotifyForm()
		form.Mail.Enable = true
		form.Mail.Username = "ops@example.com"
		form.Mail.Host = "smtp.example.com"
		form.Mail.Port = 465
		form.Mail.Password = ""
		if err := validateNotifyConfigForSave(form, true); err == nil {
			t.Fatal("expected validation error")
		}
	})

	t.Run("mail enabled without new password and no existing password", func(t *testing.T) {
		form := buildValidNotifyForm()
		form.Mail.Enable = true
		form.Mail.Username = "ops@example.com"
		form.Mail.Host = "smtp.example.com"
		form.Mail.Port = 465
		form.Mail.Password = ""
		if err := validateNotifyConfigForSave(form, false); err == nil {
			t.Fatal("expected validation error")
		}
	})
}
