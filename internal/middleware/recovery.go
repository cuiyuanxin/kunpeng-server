package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/cuiyuanxin/kunpeng/pkg/config"
	kperrors "github.com/cuiyuanxin/kunpeng/pkg/errors"
	"github.com/cuiyuanxin/kunpeng/pkg/logger"
	"github.com/cuiyuanxin/kunpeng/pkg/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Recovery 中间件，用于捕获异常
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 获取堆栈信息
				stack := string(debug.Stack())

				// 记录错误日志
				logger.Error("系统异常", fmt.Errorf("%v", err), zap.String("stack", stack))

				// 如果是生产环境，则发送报警
				if config.IsProduction() {
					go sendAlert(fmt.Sprintf("%v", err), stack, c.Request.URL.Path)
				}

				// 返回错误响应
				response.FailWithCode(c, kperrors.ErrSystem)
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}

// sendAlert 发送报警
func sendAlert(errMsg, stack, path string) {
	// TODO: 实现报警通知，可以通过邮件、飞书、钉钉、企业微信、电话、短信等方式
	// 这里只是一个占位函数，实际实现需要根据具体需求来开发
	fmt.Println(stack)
	logger.Info("发送报警通知",
		zap.String("error", errMsg),
		zap.String("path", path),
	)
}
