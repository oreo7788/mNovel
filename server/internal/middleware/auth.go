package middleware

import (
	"strings"

	"yidaiku-server/pkg/auth"
	"yidaiku-server/pkg/response"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware 认证中间件
func AuthMiddleware(jwtManager *auth.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "缺少认证信息")
			c.Abort()
			return
		}

		// 解析Bearer Token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "认证格式错误")
			c.Abort()
			return
		}

		tokenString := parts[1]
		if tokenString == "" {
			response.Unauthorized(c, "缺少认证信息")
			c.Abort()
			return
		}

		// 验证Token（不向客户端暴露具体错误，避免泄露 "bad token" 等库信息）
		claims, err := jwtManager.ValidateAccessToken(tokenString)
		if err != nil {
			response.Unauthorized(c, "认证失败，请重新登录")
			c.Abort()
			return
		}

		// 将用户ID存入上下文
		c.Set("user_id", claims.UserID)
		c.Next()
	}
}

// OptionalAuthMiddleware 可选认证中间件（不强制要求登录）
func OptionalAuthMiddleware(jwtManager *auth.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		tokenString := parts[1]
		claims, err := jwtManager.ValidateAccessToken(tokenString)
		if err == nil {
			c.Set("user_id", claims.UserID)
		}

		c.Next()
	}
}

// GetUserID 从上下文获取用户ID
func GetUserID(c *gin.Context) string {
	userID, exists := c.Get("user_id")
	if !exists {
		return ""
	}
	return userID.(string)
}

// MustGetUserID 从上下文获取用户ID（必须存在）
func MustGetUserID(c *gin.Context) (string, bool) {
	userID := GetUserID(c)
	if userID == "" {
		response.Unauthorized(c, "请先登录")
		return "", false
	}
	return userID, true
}
