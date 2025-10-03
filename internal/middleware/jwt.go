package middleware

import (
	"github.com/cuiyuanxin/kunpeng/internal/service"
	kperrors "github.com/cuiyuanxin/kunpeng/pkg/errors"
	"github.com/cuiyuanxin/kunpeng/pkg/jwt"
	"github.com/cuiyuanxin/kunpeng/pkg/response"
	"github.com/gin-gonic/gin"
)

// JWT 中间件，用于验证JWT令牌
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.FailWithCode(c, kperrors.ErrUnauthorized)
			c.Abort()
			return
		}

		// 提取token
		token := jwt.ExtractTokenFromHeader(authHeader)
		if token == "" {
			response.FailWithCode(c, kperrors.ErrInvalidToken)
			c.Abort()
			return
		}

		// 解析令牌
		claims, err := jwt.ParseToken(token)
		if err != nil {
			response.Fail(c, err)
			c.Abort()
			return
		}

		// 验证token类型必须是access token
		if claims.TokenType != jwt.AccessTokenType {
			response.FailWithCode(c, kperrors.ErrInvalidToken)
			c.Abort()
			return
		}

		// 检查token是否在黑名单中
		tokenBlacklistService := service.GetTokenBlacklistService()
		isBlacklisted, err := tokenBlacklistService.IsTokenBlacklisted(token)
		if err != nil {
			response.FailWithCode(c, kperrors.ErrSystem)
			c.Abort()
			return
		}
		if isBlacklisted {
			response.FailWithCode(c, kperrors.ErrInvalidToken)
			c.Abort()
			return
		}

		// 将用户信息存储到上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role_id", claims.RoleID)
		c.Set("app_key", claims.AppKey)

		c.Next()
	}
}
