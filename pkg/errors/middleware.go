package errors

import (
	"fmt"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// ErrorHandlerMiddleware 全局错误处理中间件
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				// 记录panic日志
				fmt.Printf("Panic recovered: %v, stack: %s, path: %s, method: %s\n",
					r,
					string(debug.Stack()),
					c.Request.URL.Path,
					c.Request.Method,
				)
				
				// 返回服务器错误
				HandleError(c, ErrInternalError.WithDetails(fmt.Sprintf("服务器内部错误: %v", r)))
				return
			}
			
			// 检查是否有错误
			if len(c.Errors) > 0 {
				lastError := c.Errors.Last()
				
				// 检查是否是业务错误
				if bizErr, ok := lastError.Err.(*BusinessError); ok {
					HandleError(c, bizErr)
					return
				}
				
				// 记录错误日志
				fmt.Printf("Request error: %s, type: %d, path: %s, method: %s\n",
					lastError.Error(),
					int(lastError.Type),
					c.Request.URL.Path,
					c.Request.Method,
				)
				
				// 返回通用错误
				HandleError(c, ErrInternalError.WithDetails(lastError.Error()))
				return
			}
		}()
		
		c.Next()
	})
}

// SetError 设置错误到Gin上下文
func SetError(c *gin.Context, err error) {
	c.Error(err)
}

// AbortWithError 中断请求并返回错误
func AbortWithError(c *gin.Context, err error) {
	HandleError(c, err)
	c.Abort()
}

// AbortWithBusinessError 中断请求并返回业务错误
func AbortWithBusinessError(c *gin.Context, bizErr *BusinessError) {
	HandleError(c, bizErr)
	c.Abort()
}