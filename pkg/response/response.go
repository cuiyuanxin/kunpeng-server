package response

import (
	"net/http"

	kperrors "github.com/cuiyuanxin/kunpeng/pkg/errors"
	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int    `json:"code"`    // 错误码
	Message string `json:"message"` // 错误消息
	Data    any    `json:"data"`    // 响应数据
	TraceID string `json:"traceId"` // 链路追踪ID
}

// Success 成功响应
func Success(c *gin.Context, data any) {
	resp := Response{
		Code:    0,
		Message: "Success",
		Data:    data,
		TraceID: getTraceID(c),
	}
	c.JSON(http.StatusOK, resp)
}

// SuccessWithMessage 带消息的成功响应
func SuccessWithMessage(c *gin.Context, message string, data any) {
	resp := Response{
		Code:    0,
		Message: message,
		Data:    data,
		TraceID: getTraceID(c),
	}
	c.JSON(http.StatusOK, resp)
}

// Fail 失败响应
func Fail(c *gin.Context, err error) {
	e := kperrors.FromError(err)
	resp := Response{
		Code:    e.Code,
		Message: e.Message,
		Data:    e.Data,
		TraceID: getTraceID(c),
	}
	c.JSON(e.StatusCode(), resp)
}

// FailWithCode 带错误码的失败响应
func FailWithCode(c *gin.Context, code int) {
	resp := Response{
		Code:    code,
		Message: kperrors.GetMessage(code),
		Data:    nil,
		TraceID: getTraceID(c),
	}
	e := kperrors.New(code, nil)
	c.JSON(e.StatusCode(), resp)
}

// FailWithMessage 带消息的失败响应
func FailWithMessage(c *gin.Context, code int, message string) {
	resp := Response{
		Code:    code,
		Message: message,
		Data:    nil,
		TraceID: getTraceID(c),
	}
	e := kperrors.New(code, nil)
	c.JSON(e.StatusCode(), resp)
}

// FailWithData 带数据的失败响应
func FailWithData(c *gin.Context, code int, data any) {
	resp := Response{
		Code:    code,
		Message: kperrors.GetMessage(code),
		Data:    data,
		TraceID: getTraceID(c),
	}
	e := kperrors.New(code, nil)
	c.JSON(e.StatusCode(), resp)
}

// FailWithMessageAndData 带消息和数据的失败响应
func FailWithMessageAndData(c *gin.Context, code int, message string, data any) {
	resp := Response{
		Code:    code,
		Message: message,
		Data:    data,
		TraceID: getTraceID(c),
	}
	e := kperrors.New(code, nil)
	c.JSON(e.StatusCode(), resp)
}

// FailWithError 带错误的失败响应
func FailWithError(c *gin.Context, err error) {
	Fail(c, err)
}

// Ok 简单成功响应，无数据
func Ok(c *gin.Context) {
	Success(c, nil)
}

// OkWithData 带数据的成功响应
func OkWithData(c *gin.Context, data any) {
	Success(c, data)
}

// OkWithMessage 带消息的成功响应
func OkWithMessage(c *gin.Context, message string) {
	SuccessWithMessage(c, message, nil)
}

// OkWithMessageAndData 带消息和数据的成功响应
func OkWithMessageAndData(c *gin.Context, message string, data any) {
	SuccessWithMessage(c, message, data)
}

// FailWithValidation 参数验证失败响应
func FailWithValidation(c *gin.Context, data any) {
	resp := Response{
		Code:    kperrors.ErrValidation,
		Message: kperrors.GetMessage(kperrors.ErrValidation),
		Data:    data,
		TraceID: getTraceID(c),
	}
	e := kperrors.New(kperrors.ErrValidation, nil)
	c.JSON(e.StatusCode(), resp)
}

// getTraceID 获取链路追踪ID
func getTraceID(c *gin.Context) string {
	traceID, exists := c.Get("trace_id")
	if !exists {
		return ""
	}
	if id, ok := traceID.(string); ok {
		return id
	}
	return ""
}
