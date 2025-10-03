package middleware

import (
	"github.com/cuiyuanxin/kunpeng/pkg/tracer"
	"github.com/gin-gonic/gin"
)

// Tracer 中间件，用于实现链路追踪
func Tracer() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取追踪器
		t := tracer.GetTracer()
		if !t.IsEnabled() {
			c.Next()
			return
		}

		// 创建追踪上下文
		ctx := t.NewContext(c.Request.Context())

		// 获取追踪ID
		traceID := tracer.GetTraceID(ctx)
		if traceID != "" {
			// 设置追踪ID到请求头
			c.Header("X-Trace-ID", traceID)
			// 设置追踪ID到上下文
			c.Set("trace_id", traceID)
		}

		// 设置追踪上下文到请求
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
