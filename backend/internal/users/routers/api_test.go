package routers

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRegisterApiRoutes_ProfileUpdatePath(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	group := r.Group("/api/v1/profile")
	RegisterApiRoutes(group)

	found := false
	for _, route := range r.Routes() {
		if route.Method == "PUT" && route.Path == "/api/v1/profile/:uid" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected PUT route /api/v1/profile/:uid, got routes: %#v", r.Routes())
	}
}
