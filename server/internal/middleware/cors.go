package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware 跨域中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Requested-With")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// RateLimitMiddleware 限流中间件（简单实现）
func RateLimitMiddleware(maxRequests int, duration time.Duration) gin.HandlerFunc {
	// 简单的滑动窗口限流
	// 生产环境建议使用Redis实现分布式限流
	return func(c *gin.Context) {
		// TODO: 实现限流逻辑
		c.Next()
	}
}

// RequestLogMiddleware 请求日志中间件
func RequestLogMiddleware() gin.HandlerFunc {
	return gin.Logger()
}

// RecoveryMiddleware 恢复中间件
func RecoveryMiddleware() gin.HandlerFunc {
	return gin.Recovery()
}
