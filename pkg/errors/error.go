package errors

import (
	"fmt"
	"net/http"

	"github.com/cuiyuanxin/kunpeng/pkg/i18n"
)

// Error 自定义错误结构
type Error struct {
	Code    int    // 错误码
	Message string // 错误消息
	Err     error  // 原始错误
	Data    any    // 附加数据
}

// New 创建新的错误（支持i18n）
func New(code int, err error) *Error {
	// 尝试获取国际化消息
	messageKey := fmt.Sprintf("error.%d", code)
	localizedMessage := i18n.T(messageKey)

	// 如果没有找到对应的翻译，使用默认消息
	message := GetMessage(code)
	if localizedMessage != messageKey {
		message = localizedMessage
	}

	return &Error{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// NewWithMessage 创建带自定义消息的错误（支持i18n）
func NewWithMessage(code int, message string, err error) *Error {
	// 如果提供了自定义消息，优先使用自定义消息
	// 否则尝试获取国际化消息
	finalMessage := message
	if message == "" {
		messageKey := fmt.Sprintf("error.%d", code)
		localizedMessage := i18n.T(messageKey)

		if localizedMessage != messageKey {
			finalMessage = localizedMessage
		} else {
			finalMessage = GetMessage(code)
		}
	}

	return &Error{
		Code:    code,
		Message: finalMessage,
		Err:     err,
	}
}

// NewWithData 创建带附加数据的错误（支持i18n）
func NewWithData(code int, err error, data any) *Error {
	// 尝试获取国际化消息
	messageKey := fmt.Sprintf("error.%d", code)
	localizedMessage := i18n.T(messageKey)

	// 如果没有找到对应的翻译，使用默认消息
	message := GetMessage(code)
	if localizedMessage != messageKey {
		message = localizedMessage
	}

	return &Error{
		Code:    code,
		Message: message,
		Err:     err,
		Data:    data,
	}
}

// Error 实现error接口（支持i18n）
func (e *Error) Error() string {
	// 尝试获取当前语言的国际化消息
	messageKey := fmt.Sprintf("error.%d", e.Code)
	localizedMessage := i18n.T(messageKey)

	// 如果找到翻译，使用翻译消息，否则使用原始消息
	message := e.Message
	if localizedMessage != messageKey {
		message = localizedMessage
	}

	if e.Err != nil {
		return fmt.Sprintf("错误码: %d, 错误信息: %s, 原始错误: %s", e.Code, message, e.Err.Error())
	}
	return fmt.Sprintf("错误码: %d, 错误信息: %s", e.Code, message)
}

// Unwrap 实现errors.Unwrap接口
func (e *Error) Unwrap() error {
	return e.Err
}

// WithMessage 设置错误消息
func (e *Error) WithMessage(message string) *Error {
	e.Message = message
	return e
}

// WithData 设置附加数据
func (e *Error) WithData(data any) *Error {
	e.Data = data
	return e
}

// StatusCode 获取对应的HTTP状态码
func (e *Error) StatusCode() int {
	switch {
	case e.Code == ErrUnauthorized || e.Code == ErrInvalidToken || e.Code == ErrTokenExpired:
		return http.StatusUnauthorized
	case e.Code == ErrForbidden || e.Code == ErrPermDenied:
		return http.StatusForbidden
	case e.Code == ErrNotFound:
		return http.StatusNotFound
	case e.Code == ErrMethodNotAllowed:
		return http.StatusMethodNotAllowed
	case e.Code == ErrParam:
		return http.StatusBadRequest
	case e.Code == ErrTooManyRequests:
		return http.StatusTooManyRequests
	case e.Code == ErrTimeout:
		return http.StatusRequestTimeout
	case IsSystemError(e.Code):
		return http.StatusInternalServerError
	default:
		return http.StatusBadRequest
	}
}

// IsCode 判断错误是否为指定错误码
func IsCode(err error, code int) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(*Error); ok {
		return e.Code == code
	}
	return false
}

// FromError 从标准错误转换为自定义错误
func FromError(err error) *Error {
	if err == nil {
		return nil
	}
	if e, ok := err.(*Error); ok {
		return e
	}
	return New(ErrSystem, err)
}
