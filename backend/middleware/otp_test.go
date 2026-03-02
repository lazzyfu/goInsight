package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/lazzyfu/goinsight/internal/global"
	userModels "github.com/lazzyfu/goinsight/internal/users/models"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupOTPMiddlewareTestDB(t *testing.T) {
	t.Helper()

	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", strings.ReplaceAll(t.Name(), "/", "_"))
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite failed: %v", err)
	}
	if err := db.Exec(`
		CREATE TABLE insight_users (
			uid INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			is_active BOOLEAN NOT NULL DEFAULT 1,
			is_two_fa BOOLEAN NOT NULL DEFAULT 0,
			otp_secret TEXT
		);
	`).Error; err != nil {
		t.Fatalf("create table failed: %v", err)
	}

	oldDB := global.App.DB
	oldLog := global.App.Log
	global.App.DB = db
	global.App.Log = logrus.New()
	t.Cleanup(func() {
		global.App.DB = oldDB
		global.App.Log = oldLog
	})
}

func TestOTPMiddleware_WhenTwoFAEnabledAndSecretMissingAlwaysReturnsBindHint(t *testing.T) {
	setupOTPMiddlewareTestDB(t)
	gin.SetMode(gin.TestMode)

	if err := global.App.DB.Table("insight_users").Create(map[string]any{
		"username":   "alice",
		"password":   userModels.BcryptPW("secret123"),
		"is_active":  true,
		"is_two_fa":  true,
		"otp_secret": "",
	}).Error; err != nil {
		t.Fatalf("create user failed: %v", err)
	}

	r := gin.New()
	r.POST("/login", OTPMiddleware(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	body, _ := json.Marshal(gin.H{
		"username": "alice",
		"password": "secret123",
		"otp_code": "123456",
	})
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var resp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response failed: %v", err)
	}
	if resp["code"] != "4001" {
		t.Fatalf("expected code 4001, got %#v", resp["code"])
	}
}

func TestOTPMiddleware_WhenTwoFAEnabledAndNoCodeReturnsInputHint(t *testing.T) {
	setupOTPMiddlewareTestDB(t)
	gin.SetMode(gin.TestMode)

	if err := global.App.DB.Table("insight_users").Create(map[string]any{
		"username":   "bob",
		"password":   userModels.BcryptPW("secret123"),
		"is_active":  true,
		"is_two_fa":  true,
		"otp_secret": "JBSWY3DPEHPK3PXP",
	}).Error; err != nil {
		t.Fatalf("create user failed: %v", err)
	}

	r := gin.New()
	r.POST("/login", OTPMiddleware(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	body, _ := json.Marshal(gin.H{
		"username": "bob",
		"password": "secret123",
	})
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var resp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response failed: %v", err)
	}
	if resp["code"] != "4002" {
		t.Fatalf("expected code 4002, got %#v", resp["code"])
	}
}
