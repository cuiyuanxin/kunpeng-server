package middleware

import (
	"context"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CORS 跨域中间件
func CORS() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "authorization, origin, content-type, accept")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})
}

// Logger 日志中间件
func Logger(logger *zap.Logger) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// 处理请求
		c.Next()

		// 记录日志
		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		if raw != "" {
			path = path + "?" + raw
		}

		logger.Info("HTTP Request",
			zap.String("method", method),
			zap.String("path", path),
			zap.String("client_ip", clientIP),
			zap.Int("status_code", statusCode),
			zap.Duration("latency", latency),
			zap.String("user_agent", c.Request.UserAgent()),
		)
	})
}

// CustomLogger 自定义日志中间件（生产环境使用）
// 参考《Go程序编程之旅》的实现方式
func CustomLogger(logger *zap.Logger) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		requestID := c.GetString("request_id")
		if requestID == "" {
			requestID = generateRequestID()
			c.Set("request_id", requestID)
		}

		// 记录请求开始
		logger.Info("Request started",
			zap.String("request_id", requestID),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("client_ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.String("referer", c.Request.Referer()),
		)

		// 处理请求
		c.Next()

		// 记录请求结束
		latency := time.Since(start)
		statusCode := c.Writer.Status()
		bodySize := c.Writer.Size()

		if raw != "" {
			path = path + "?" + raw
		}

		// 根据状态码选择日志级别
		logFields := []zap.Field{
			zap.String("request_id", requestID),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("client_ip", c.ClientIP()),
			zap.Int("status_code", statusCode),
			zap.Duration("latency", latency),
			zap.Int("body_size", bodySize),
			zap.String("user_agent", c.Request.UserAgent()),
		}

		if statusCode >= 500 {
			logger.Error("Request completed with server error", logFields...)
		} else if statusCode >= 400 {
			logger.Warn("Request completed with client error", logFields...)
		} else {
			logger.Info("Request completed successfully", logFields...)
		}
	})
}

// Recovery 恢复中间件
func Recovery(logger *zap.Logger, responseFunc func(*gin.Context, string)) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		logger.Error("Panic recovered",
			zap.Any("error", recovered),
			zap.String("path", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
		)

		responseFunc(c, "Internal server error")
		c.Abort()
	})
}

// CustomRecovery 自定义恢复中间件（生产环境使用）
// 参考《Go程序编程之旅》的实现方式
func CustomRecovery(logger *zap.Logger, responseFunc func(*gin.Context, string)) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		requestID := c.GetString("request_id")
		if requestID == "" {
			requestID = generateRequestID()
		}

		// 获取堆栈信息
		stack := make([]byte, 4096)
		length := runtime.Stack(stack, false)
		stackTrace := string(stack[:length])

		// 记录详细的panic信息
		logger.Error("Panic recovered in custom recovery",
			zap.String("request_id", requestID),
			zap.Any("panic", recovered),
			zap.String("path", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
			zap.String("client_ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.String("stack_trace", stackTrace),
			zap.Time("timestamp", time.Now()),
		)

		// 发送错误响应
		responseFunc(c, "Internal server error")
		c.Abort()
	})
}

// RequestID 请求ID中间件
func RequestID() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			// 生成新的请求ID
			requestID = generateRequestID()
		}

		c.Header("X-Request-ID", requestID)
		c.Set("request_id", requestID)
		c.Next()
	})
}

// generateRequestID 生成请求ID
func generateRequestID() string {
	// 简单的时间戳+随机数实现
	// 实际项目中可以使用UUID或其他更复杂的算法
	return time.Now().Format("20060102150405") + "-" + "random"
}

// Timeout 超时中间件
func Timeout(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// RateLimiter 限流中间件（简单实现）
func RateLimiter() gin.HandlerFunc {
	// 这里可以集成更复杂的限流算法，如令牌桶、滑动窗口等
	return gin.HandlerFunc(func(c *gin.Context) {
		// 简单的IP限流示例
		// 实际项目中建议使用Redis实现分布式限流
		c.Next()
	})
}