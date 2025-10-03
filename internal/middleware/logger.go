package middleware

import (
	"bytes"
	"io"
	"time"

	"github.com/cuiyuanxin/kunpeng/pkg/config"
	"github.com/cuiyuanxin/kunpeng/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Logger 中间件，用于记录访问日志
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 如果不是生产环境，则跳过
		if !config.IsProduction() {
			c.Next()
			return
		}

		// 开始时间
		start := time.Now()

		// 获取请求体
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 创建自定义响应写入器
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// 处理请求
		c.Next()

		// 结束时间
		end := time.Now()
		latency := end.Sub(start)

		// 获取用户信息
		userID, _ := c.Get("user_id")
		username, _ := c.Get("username")

		// 记录日志
		logger.Info("访问日志",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("query", c.Request.URL.RawQuery),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.Any("user_id", userID),
			zap.Any("username", username),
			zap.Int("status", c.Writer.Status()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("latency", latency),
			zap.String("request", string(requestBody)),
			zap.String("response", blw.body.String()),
		)
	}
}

// bodyLogWriter 自定义响应写入器，用于记录响应体
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write 重写Write方法，用于记录响应体
func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
