package middleware

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

const identityClaimKey = "id"

func GetUserNameFromJWT(c *gin.Context) (string, bool) {
	claims := jwt.ExtractClaims(c)
	raw, ok := claims[identityClaimKey]
	if !ok {
		return "", false
	}
	username, ok := raw.(string)
	if !ok || username == "" {
		return "", false
	}
	return username, true
}
