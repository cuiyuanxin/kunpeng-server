package middleware

import (
	"context"
	"time"

	kperrors "github.com/cuiyuanxin/kunpeng/pkg/errors"
	"github.com/cuiyuanxin/kunpeng/pkg/response"
	"github.com/gin-gonic/gin"
)

// TimeoutConfig 超时配置
type TimeoutConfig struct {
	Timeout time.Duration // 超时时间
}

// DefaultTimeoutConfig 默认超时配置
var DefaultTimeoutConfig = TimeoutConfig{
	Timeout: 30 * time.Second, // 默认30秒超时
}

// Timeout 中间件，用于控制请求超时
func Timeout(config TimeoutConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 创建带超时的上下文
		ctx, cancel := context.WithTimeout(c.Request.Context(), config.Timeout)
		defer cancel()

		// 将超时上下文设置到请求中
		c.Request = c.Request.WithContext(ctx)

		// 使用通道来检测请求是否完成
		done := make(chan struct{})

		go func() {
			defer close(done)
			c.Next()
		}()

		select {
		case <-done:
			// 请求正常完成
			return
		case <-ctx.Done():
			// 请求超时
			if ctx.Err() == context.DeadlineExceeded {
				// 检查响应是否已经写入
				if !c.Writer.Written() {
					c.Abort()
					response.FailWithCode(c, kperrors.ErrTimeout)
				}
			}
			return
		}
	}
}

// TimeoutWithDefault 使用默认配置的超时中间件
func TimeoutWithDefault() gin.HandlerFunc {
	return Timeout(DefaultTimeoutConfig)
}
