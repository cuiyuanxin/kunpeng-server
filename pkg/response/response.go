package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/cuiyuanxin/kunpeng/pkg/errors"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// PageResponse 分页响应结构
type PageResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Total   int64       `json:"total"`
	Page    int         `json:"page"`
	Size    int         `json:"size"`
}

// 响应状态码常量
const (
	CodeSuccess      = 200
	CodeBadRequest   = 400
	CodeUnauthorized = 401
	CodeForbidden    = 403
	CodeNotFound     = 404
	CodeServerError  = 500
)

// 响应消息常量
const (
	MsgSuccess      = "success"
	MsgBadRequest   = "bad request"
	MsgUnauthorized = "unauthorized"
	MsgForbidden    = "forbidden"
	MsgNotFound     = "not found"
	MsgServerError  = "internal server error"
)

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: MsgSuccess,
		Data:    data,
	})
}

// SuccessWithMessage 带消息的成功响应
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: message,
		Data:    data,
	})
}

// SuccessPage 分页成功响应
func SuccessPage(c *gin.Context, data interface{}, total int64, page, size int) {
	c.JSON(http.StatusOK, PageResponse{
		Code:    CodeSuccess,
		Message: MsgSuccess,
		Data:    data,
		Total:   total,
		Page:    page,
		Size:    size,
	})
}

// BadRequest 错误请求响应
func BadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, Response{
		Code:    CodeBadRequest,
		Message: message,
	})
}

// BadRequestWithData 带数据的错误请求响应
func BadRequestWithData(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusBadRequest, Response{
		Code:    CodeBadRequest,
		Message: message,
		Data:    data,
	})
}

// Unauthorized 未授权响应
func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, Response{
		Code:    CodeUnauthorized,
		Message: message,
	})
}

// Forbidden 禁止访问响应
func Forbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, Response{
		Code:    CodeForbidden,
		Message: message,
	})
}

// NotFound 未找到响应
func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, Response{
		Code:    CodeNotFound,
		Message: message,
	})
}

// ServerError 服务器错误响应
func ServerError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, Response{
		Code:    CodeServerError,
		Message: message,
	})
}

// ServerErrorWithData 带数据的服务器错误响应
func ServerErrorWithData(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusInternalServerError, Response{
		Code:    CodeServerError,
		Message: message,
		Data:    data,
	})
}

// Error 错误响应（通用）
func Error(c *gin.Context, httpCode int, message string, data interface{}) {
	c.JSON(httpCode, Response{
		Code:    httpCode,
		Message: message,
		Data:    data,
	})
}

// Custom 自定义响应
func Custom(c *gin.Context, httpCode, code int, message string, data interface{}) {
	c.JSON(httpCode, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

// HandleBusinessError 处理业务错误的统一响应
func HandleBusinessError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	// 如果是业务错误，使用预设的错误码和描述
	if bizErr, ok := err.(*errors.BusinessError); ok {
		responseData := Response{
			Code:    int(bizErr.Code),
			Message: bizErr.Message,
		}
		
		// 如果有详细信息，添加到响应中
		if bizErr.Details != "" {
			responseData.Data = gin.H{"details": bizErr.Details}
		}
		
		c.JSON(bizErr.HTTPStatus, responseData)
		return
	}

	// 其他错误统一处理为内部服务器错误
	ServerError(c, "服务器内部错误")
}

// SuccessOrError 根据错误情况返回成功或错误响应
func SuccessOrError(c *gin.Context, data interface{}, err error) {
	if err != nil {
		HandleBusinessError(c, err)
		return
	}
	Success(c, data)
}

// SuccessPageOrError 根据错误情况返回分页成功或错误响应
func SuccessPageOrError(c *gin.Context, data interface{}, total int64, page, size int, err error) {
	if err != nil {
		HandleBusinessError(c, err)
		return
	}
	SuccessPage(c, data, total, page, size)
}

// BusinessError 直接返回业务错误
func BusinessError(c *gin.Context, bizErr *errors.BusinessError) {
	HandleBusinessError(c, bizErr)
}

// BusinessErrorWithDetails 返回带详细信息的业务错误
func BusinessErrorWithDetails(c *gin.Context, bizErr *errors.BusinessError, details string) {
	HandleBusinessError(c, bizErr.WithDetails(details))
}

// ValidationError 返回验证错误
func ValidationError(c *gin.Context, details string) {
	HandleBusinessError(c, errors.ErrValidationFailed.WithDetails(details))
}

// AuthError 返回认证错误
func AuthError(c *gin.Context, details string) {
	HandleBusinessError(c, errors.ErrUnauthorized.WithDetails(details))
}

// PermissionError 返回权限错误
func PermissionError(c *gin.Context, details string) {
	HandleBusinessError(c, errors.ErrPermissionDenied.WithDetails(details))
}

// DatabaseError 返回数据库错误
func DatabaseError(c *gin.Context, details string) {
	HandleBusinessError(c, errors.ErrDatabaseError.WithDetails(details))
}