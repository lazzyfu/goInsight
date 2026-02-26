package middleware

import (
	"net/http/httptest"
	"testing"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

func TestGetUserNameFromJWT(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("missing id claim", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Set("JWT_PAYLOAD", jwt.MapClaims{"nick_name": "n"})

		username, ok := GetUserNameFromJWT(c)
		if ok {
			t.Fatalf("expected missing claim to be invalid, got username=%q", username)
		}
	})

	t.Run("invalid id claim type", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Set("JWT_PAYLOAD", jwt.MapClaims{"id": 123})

		username, ok := GetUserNameFromJWT(c)
		if ok {
			t.Fatalf("expected non-string claim to be invalid, got username=%q", username)
		}
	})

	t.Run("valid id claim", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Set("JWT_PAYLOAD", jwt.MapClaims{"id": "alice"})

		username, ok := GetUserNameFromJWT(c)
		if !ok || username != "alice" {
			t.Fatalf("expected valid username alice, got ok=%v username=%q", ok, username)
		}
	})
}
