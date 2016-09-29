package main

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// The result of an auth/auth check
type AuthResult struct {
	Authorized bool
	Error      string
}

// Checks authentication/authorization
func checkAuth(c *gin.Context) AuthResult {
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		return AuthResult{
			Authorized: false,
			Error:      "Missing Authorization header",
		}
	}
	if !strings.HasPrefix(authHeader, "Token ") {
		return AuthResult{
			Authorized: false,
			Error:      "Malformatted Authorization header",
		}
	}
	if authHeader != "Token "+Config.Token {
		return AuthResult{
			Authorized: false,
			Error:      "Token in Authorization header is incorrect",
		}
	}
	return AuthResult{
		Authorized: true,
	}
}
