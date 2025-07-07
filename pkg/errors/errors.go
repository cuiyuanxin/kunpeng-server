package errors

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/cuiyuanxin/kunpeng/pkg/response"
)

// ErrorCode 错误码类型
type ErrorCode int

// BusinessError 业务错误结构
type BusinessError struct {
	Code       ErrorCode `json:"code"`
	Message    string    `json:"message"`
	HTTPStatus int       `json:"-"`
	Details    string    `json:"details,omitempty"`
}

// Error 实现 error 接口
func (e *BusinessError) Error() string {
	return e.Message
}

// WithDetails 添加详细信息
func (e *BusinessError) WithDetails(details string) *BusinessError {
	return &BusinessError{
		Code:       e.Code,
		Message:    e.Message,
		HTTPStatus: e.HTTPStatus,
		Details:    details,
	}
}

// WithDetailsf 添加格式化的详细信息
func (e *BusinessError) WithDetailsf(format string, args ...interface{}) *BusinessError {
	return e.WithDetails(fmt.Sprintf(format, args...))
}

// 错误码定义
const (
	// 通用错误 (1000-1999)
	ErrCodeSuccess           ErrorCode = 0
	ErrCodeInternalError     ErrorCode = 1000
	ErrCodeInvalidParams     ErrorCode = 1001
	ErrCodeValidationFailed  ErrorCode = 1002
	ErrCodeResourceNotFound  ErrorCode = 1003
	ErrCodeResourceExists    ErrorCode = 1004
	ErrCodePermissionDenied  ErrorCode = 1005
	ErrCodeRateLimitExceeded ErrorCode = 1006
	ErrCodeServiceUnavailable ErrorCode = 1007

	// 认证相关错误 (2000-2999)
	ErrCodeUnauthorized      ErrorCode = 2000
	ErrCodeTokenExpired      ErrorCode = 2001
	ErrCodeTokenInvalid      ErrorCode = 2002
	ErrCodeLoginFailed       ErrorCode = 2003
	ErrCodePasswordIncorrect ErrorCode = 2004
	ErrCodeAccountLocked     ErrorCode = 2005
	ErrCodeAccountDisabled   ErrorCode = 2006

	// 用户相关错误 (3000-3999)
	ErrCodeUserNotFound      ErrorCode = 3000
	ErrCodeUserExists        ErrorCode = 3001
	ErrCodeUsernameTaken     ErrorCode = 3002
	ErrCodeEmailTaken        ErrorCode = 3003
	ErrCodeInvalidEmail      ErrorCode = 3004
	ErrCodeWeakPassword      ErrorCode = 3005

	// 管理员相关错误 (4000-4999)
	ErrCodeAdminNotFound     ErrorCode = 4000
	ErrCodeAdminExists       ErrorCode = 4001
	ErrCodeInsufficientPrivileges ErrorCode = 4002

	// 角色权限相关错误 (5000-5999)
	ErrCodeRoleNotFound      ErrorCode = 5000
	ErrCodeRoleExists        ErrorCode = 5001
	ErrCodePermissionNotFound ErrorCode = 5002
	ErrCodeRoleInUse         ErrorCode = 5003

	// 数据库相关错误 (6000-6999)
	ErrCodeDatabaseError     ErrorCode = 6000
	ErrCodeRecordNotFound    ErrorCode = 6001
	ErrCodeDuplicateKey      ErrorCode = 6002
	ErrCodeConstraintViolation ErrorCode = 6003

	// 文件相关错误 (7000-7999)
	ErrCodeFileNotFound      ErrorCode = 7000
	ErrCodeFileUploadFailed  ErrorCode = 7001
	ErrCodeFileTypeNotAllowed ErrorCode = 7002
	ErrCodeFileSizeExceeded  ErrorCode = 7003

	// 第三方服务错误 (8000-8999)
	ErrCodeExternalServiceError ErrorCode = 8000
	ErrCodeAPICallFailed     ErrorCode = 8001
	ErrCodeNetworkError      ErrorCode = 8002
)

// 预定义错误实例
var (
	// 通用错误
	ErrInternalError     = &BusinessError{ErrCodeInternalError, "服务器内部错误", http.StatusInternalServerError, ""}
	ErrInvalidParams     = &BusinessError{ErrCodeInvalidParams, "请求参数无效", http.StatusBadRequest, ""}
	ErrValidationFailed  = &BusinessError{ErrCodeValidationFailed, "数据验证失败", http.StatusBadRequest, ""}
	ErrResourceNotFound  = &BusinessError{ErrCodeResourceNotFound, "资源不存在", http.StatusNotFound, ""}
	ErrResourceExists    = &BusinessError{ErrCodeResourceExists, "资源已存在", http.StatusConflict, ""}
	ErrPermissionDenied  = &BusinessError{ErrCodePermissionDenied, "权限不足", http.StatusForbidden, ""}
	ErrRateLimitExceeded = &BusinessError{ErrCodeRateLimitExceeded, "请求频率超限", http.StatusTooManyRequests, ""}
	ErrServiceUnavailable = &BusinessError{ErrCodeServiceUnavailable, "服务不可用", http.StatusServiceUnavailable, ""}

	// 认证相关错误
	ErrUnauthorized      = &BusinessError{ErrCodeUnauthorized, "未授权访问", http.StatusUnauthorized, ""}
	ErrTokenExpired      = &BusinessError{ErrCodeTokenExpired, "令牌已过期", http.StatusUnauthorized, ""}
	ErrTokenInvalid      = &BusinessError{ErrCodeTokenInvalid, "令牌无效", http.StatusUnauthorized, ""}
	ErrLoginFailed       = &BusinessError{ErrCodeLoginFailed, "登录失败", http.StatusUnauthorized, ""}
	ErrPasswordIncorrect = &BusinessError{ErrCodePasswordIncorrect, "密码错误", http.StatusUnauthorized, ""}
	ErrAccountLocked     = &BusinessError{ErrCodeAccountLocked, "账户已锁定", http.StatusForbidden, ""}
	ErrAccountDisabled   = &BusinessError{ErrCodeAccountDisabled, "账户已禁用", http.StatusForbidden, ""}

	// 用户相关错误
	ErrUserNotFound      = &BusinessError{ErrCodeUserNotFound, "用户不存在", http.StatusNotFound, ""}
	ErrUserExists        = &BusinessError{ErrCodeUserExists, "用户已存在", http.StatusConflict, ""}
	ErrUsernameTaken     = &BusinessError{ErrCodeUsernameTaken, "用户名已被占用", http.StatusConflict, ""}
	ErrEmailTaken        = &BusinessError{ErrCodeEmailTaken, "邮箱已被占用", http.StatusConflict, ""}
	ErrInvalidEmail      = &BusinessError{ErrCodeInvalidEmail, "邮箱格式无效", http.StatusBadRequest, ""}
	ErrWeakPassword      = &BusinessError{ErrCodeWeakPassword, "密码强度不足", http.StatusBadRequest, ""}

	// 管理员相关错误
	ErrAdminNotFound     = &BusinessError{ErrCodeAdminNotFound, "管理员不存在", http.StatusNotFound, ""}
	ErrAdminExists       = &BusinessError{ErrCodeAdminExists, "管理员已存在", http.StatusConflict, ""}
	ErrInsufficientPrivileges = &BusinessError{ErrCodeInsufficientPrivileges, "权限不足", http.StatusForbidden, ""}

	// 角色权限相关错误
	ErrRoleNotFound      = &BusinessError{ErrCodeRoleNotFound, "角色不存在", http.StatusNotFound, ""}
	ErrRoleExists        = &BusinessError{ErrCodeRoleExists, "角色已存在", http.StatusConflict, ""}
	ErrPermissionNotFound = &BusinessError{ErrCodePermissionNotFound, "权限不存在", http.StatusNotFound, ""}
	ErrRoleInUse         = &BusinessError{ErrCodeRoleInUse, "角色正在使用中", http.StatusConflict, ""}

	// 数据库相关错误
	ErrDatabaseError     = &BusinessError{ErrCodeDatabaseError, "数据库操作失败", http.StatusInternalServerError, ""}
	ErrRecordNotFound    = &BusinessError{ErrCodeRecordNotFound, "记录不存在", http.StatusNotFound, ""}
	ErrDuplicateKey      = &BusinessError{ErrCodeDuplicateKey, "数据重复", http.StatusConflict, ""}
	ErrConstraintViolation = &BusinessError{ErrCodeConstraintViolation, "数据约束违反", http.StatusBadRequest, ""}

	// 文件相关错误
	ErrFileNotFound      = &BusinessError{ErrCodeFileNotFound, "文件不存在", http.StatusNotFound, ""}
	ErrFileUploadFailed  = &BusinessError{ErrCodeFileUploadFailed, "文件上传失败", http.StatusInternalServerError, ""}
	ErrFileTypeNotAllowed = &BusinessError{ErrCodeFileTypeNotAllowed, "文件类型不允许", http.StatusBadRequest, ""}
	ErrFileSizeExceeded  = &BusinessError{ErrCodeFileSizeExceeded, "文件大小超限", http.StatusBadRequest, ""}

	// 第三方服务错误
	ErrExternalServiceError = &BusinessError{ErrCodeExternalServiceError, "外部服务错误", http.StatusBadGateway, ""}
	ErrAPICallFailed     = &BusinessError{ErrCodeAPICallFailed, "API调用失败", http.StatusBadGateway, ""}
	ErrNetworkError      = &BusinessError{ErrCodeNetworkError, "网络错误", http.StatusBadGateway, ""}
)

// HandleError 统一错误处理函数
func HandleError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	// 如果是业务错误，直接返回
	if bizErr, ok := err.(*BusinessError); ok {
		responseData := gin.H{
			"code":    bizErr.Code,
			"message": bizErr.Message,
		}
		if bizErr.Details != "" {
			responseData["details"] = bizErr.Details
		}
		c.JSON(bizErr.HTTPStatus, responseData)
		return
	}

	// 其他错误统一处理为内部服务器错误
	response.ServerError(c, "服务器内部错误")
}

// NewBusinessError 创建新的业务错误
func NewBusinessError(code ErrorCode, message string, httpStatus int) *BusinessError {
	return &BusinessError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
	}
}

// IsBusinessError 判断是否为业务错误
func IsBusinessError(err error) bool {
	_, ok := err.(*BusinessError)
	return ok
}

// GetErrorCode 获取错误码
func GetErrorCode(err error) ErrorCode {
	if bizErr, ok := err.(*BusinessError); ok {
		return bizErr.Code
	}
	return ErrCodeInternalError
}