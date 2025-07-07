package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/cuiyuanxin/kunpeng/pkg/auth"
	"github.com/cuiyuanxin/kunpeng/pkg/response"
)

// AuthMiddleware 认证中间件
type AuthMiddleware struct {
	jwtManager *auth.JWTManager
	db         *gorm.DB
}

// NewAuthMiddleware 创建认证中间件
func NewAuthMiddleware(jwtManager *auth.JWTManager, db *gorm.DB) *AuthMiddleware {
	return &AuthMiddleware{
		jwtManager: jwtManager,
		db:         db,
	}
}

// JWTAuth JWT认证中间件
func (m *AuthMiddleware) JWTAuth() gin.HandlerFunc {
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
		claims, err := m.jwtManager.ParseToken(tokenString)
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

// AdminAuth 管理员权限中间件
func (m *AuthMiddleware) AdminAuth() gin.HandlerFunc {
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

// OptionalAuth 可选认证中间件（不强制要求认证）
func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// 没有认证头，继续执行但不设置用户信息
			c.Next()
			return
		}

		// 检查Bearer前缀
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			// 格式错误，继续执行但不设置用户信息
			c.Next()
			return
		}

		tokenString := parts[1]
		claims, err := m.jwtManager.ParseToken(tokenString)
		if err != nil {
			// Token无效，继续执行但不设置用户信息
			c.Next()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// RequireAuth 要求认证的中间件（简化版）
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return m.JWTAuth()
}