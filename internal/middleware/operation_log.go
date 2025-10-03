package middleware

import (
	"bytes"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/cuiyuanxin/kunpeng/internal/service"
	"github.com/gin-gonic/gin"
)

// responseWriter 包装gin.ResponseWriter以捕获响应内容
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// OperationLog 操作日志中间件
func OperationLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 跳过不需要记录的路径
		if shouldSkipPath(c.Request.URL.Path) {
			c.Next()
			return
		}

		// 记录开始时间
		startTime := time.Now()

		// 读取请求体
		reqBody := ""
		if c.Request.Body != nil {
			bodyBytes, _ := io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			reqBody = string(bodyBytes)
		}

		// 包装ResponseWriter以捕获响应
		w := &responseWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		c.Writer = w

		// 处理请求
		c.Next()

		// 计算耗时
		duration := int(time.Since(startTime).Milliseconds())

		// 获取用户信息
		userID, username := getUserInfo(c)

		// 获取响应内容
		respBody := w.body.String()
		if len(respBody) > 1000 {
			respBody = respBody[:1000] + "..."
		}

		// 获取错误信息
		errorMessage := ""
		if len(c.Errors) > 0 {
			errorMessage = c.Errors.String()
		}

		// 获取状态
		status := int8(1) // 成功
		if c.Writer.Status() >= 400 {
			status = 0 // 失败
		}

		// 异步记录操作日志
		go func() {
			operationLogService := service.GetOperationLogService()
			operationLogService.RecordOperationLog(
				userID,
				username,
				getModule(c.Request.URL.Path),
				getAction(c.Request.Method, c.Request.URL.Path),
				c.Request.Method,
				c.Request.URL.Path,
				c.ClientIP(),
				c.Request.UserAgent(),
				reqBody,
				respBody,
				status,
				errorMessage,
				duration,
			)
		}()
	}
}

// shouldSkipPath 判断是否跳过记录的路径
func shouldSkipPath(path string) bool {
	skipPaths := []string{
		"/health",
		"/swagger",
		"/api/v1/user/info", // 跳过获取用户信息接口，避免频繁记录
	}

	for _, skipPath := range skipPaths {
		if strings.HasPrefix(path, skipPath) {
			return true
		}
	}
	return false
}

// getUserInfo 获取用户信息
func getUserInfo(c *gin.Context) (uint, string) {
	userID := uint(0)
	username := "anonymous"

	// 从JWT中获取用户信息
	if userIDStr, exists := c.Get("userID"); exists {
		if id, ok := userIDStr.(string); ok {
			if parsedID, err := strconv.ParseUint(id, 10, 32); err == nil {
				userID = uint(parsedID)
			}
		}
	}

	if usernameStr, exists := c.Get("username"); exists {
		if name, ok := usernameStr.(string); ok {
			username = name
		}
	}

	return userID, username
}

// getModule 根据路径获取模块名
func getModule(path string) string {
	if strings.Contains(path, "/users") {
		return "用户管理"
	}
	if strings.Contains(path, "/roles") {
		return "角色管理"
	}
	if strings.Contains(path, "/menus") {
		return "菜单管理"
	}
	if strings.Contains(path, "/apis") {
		return "API管理"
	}
	if strings.Contains(path, "/depts") {
		return "部门管理"
	}
	if strings.Contains(path, "/posts") {
		return "岗位管理"
	}
	if strings.Contains(path, "/dict") {
		return "字典管理"
	}
	if strings.Contains(path, "/user-logs") {
		return "日志管理"
	}
	return "系统管理"
}

// getAction 根据请求方法和路径获取操作名
func getAction(method, path string) string {
	switch method {
	case "GET":
		if strings.Contains(path, "/all") {
			return "获取全部"
		}
		if strings.Contains(path, "/tree") {
			return "获取树形结构"
		}
		if strings.HasSuffix(path, "s") || !strings.Contains(path, "/") {
			return "查询列表"
		}
		return "查询详情"
	case "POST":
		if strings.Contains(path, "/clean") {
			return "清理数据"
		}
		return "新增"
	case "PUT":
		if strings.Contains(path, "/status") {
			return "修改状态"
		}
		if strings.Contains(path, "/password") {
			return "修改密码"
		}
		return "修改"
	case "DELETE":
		if strings.Contains(path, "/batch") {
			return "批量删除"
		}
		return "删除"
	default:
		return "未知操作"
	}
}
