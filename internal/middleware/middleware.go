// Package middleware 包含项目特定的中间件
// 通用中间件已迁移到 pkg/middleware
package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	gorm "gorm.io/gorm"

	"github.com/cuiyuanxin/kunpeng/pkg/auth"
	"github.com/cuiyuanxin/kunpeng/pkg/response"
	"github.com/cuiyuanxin/kunpeng/pkg/middleware"
)

// 重新导出通用中间件
var (
	CORS           = middleware.CORS
	Logger         = middleware.Logger
	CustomLogger   = middleware.CustomLogger
	Recovery       = middleware.Recovery
	CustomRecovery = middleware.CustomRecovery
	RequestID      = middleware.RequestID
	Timeout        = middleware.Timeout
	RateLimiter    = middleware.RateLimiter
)

// JWTAuth JWT认证中间件（项目特定）
func JWTAuth(jwtManager *auth.JWTManager, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "Missing authorization header")
			c.Abort()
			return
		}

		// 检查Bearer前缀
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "Invalid authorization header format")
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := jwtManager.ParseToken(tokenString)
		if err != nil {
			response.Unauthorized(c, "Invalid token: "+err.Error())
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// AdminAuth 管理员权限中间件（项目特定）
func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			response.Unauthorized(c, "Missing role information")
			c.Abort()
			return
		}

		if role != "admin" {
			response.Forbidden(c, "Admin access required")
			c.Abort()
			return
		}

		c.Next()
	}
}

// PermissionAuth 权限验证中间件（项目特定）
func PermissionAuth(requiredPermission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			response.Unauthorized(c, "Missing user information")
			c.Abort()
			return
		}

		// TODO: 实现权限检查逻辑
		// 这里应该查询数据库检查用户是否具有所需权限
		_ = userID
		_ = requiredPermission

		c.Next()
	}
}
