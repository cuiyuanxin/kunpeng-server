package errors

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// WrapError 包装错误，添加上下文信息
func WrapError(err error, message string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", message, err)
}

// WrapErrorf 包装错误，添加格式化的上下文信息
func WrapErrorf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf(format+": %w", append(args, err)...)
}

// ConvertGormError 转换GORM错误为业务错误
func ConvertGormError(err error) *BusinessError {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return ErrRecordNotFound
	case errors.Is(err, gorm.ErrInvalidTransaction):
		return ErrDatabaseError.WithDetails("无效的事务操作")
	case errors.Is(err, gorm.ErrNotImplemented):
		return ErrDatabaseError.WithDetails("功能未实现")
	case errors.Is(err, gorm.ErrMissingWhereClause):
		return ErrDatabaseError.WithDetails("缺少WHERE条件")
	case errors.Is(err, gorm.ErrUnsupportedRelation):
		return ErrDatabaseError.WithDetails("不支持的关联关系")
	case errors.Is(err, gorm.ErrPrimaryKeyRequired):
		return ErrDatabaseError.WithDetails("缺少主键")
	default:
		errorMsg := err.Error()
		// 检查常见的数据库错误模式
		switch {
		case strings.Contains(errorMsg, "duplicate") || strings.Contains(errorMsg, "unique"):
			return ErrDuplicateKey
		case strings.Contains(errorMsg, "foreign key constraint"):
			return ErrConstraintViolation.WithDetails("外键约束违反")
		case strings.Contains(errorMsg, "check constraint"):
			return ErrConstraintViolation.WithDetails("检查约束违反")
		case strings.Contains(errorMsg, "not null constraint"):
			return ErrConstraintViolation.WithDetails("非空约束违反")
		default:
			return ErrDatabaseError.WithDetails(errorMsg)
		}
	}
}

// ConvertValidationError 转换验证错误为业务错误
func ConvertValidationError(err error) *BusinessError {
	if err == nil {
		return nil
	}

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		var messages []string
		for _, fieldError := range validationErrors {
			message := getValidationErrorMessage(fieldError)
			messages = append(messages, message)
		}
		return ErrValidationFailed.WithDetails(strings.Join(messages, "; "))
	}

	return ErrValidationFailed.WithDetails(err.Error())
}

// getValidationErrorMessage 获取验证错误的友好消息
func getValidationErrorMessage(fieldError validator.FieldError) string {
	fieldName := fieldError.Field()
	tag := fieldError.Tag()
	param := fieldError.Param()

	switch tag {
	case "required":
		return fmt.Sprintf("%s 是必填字段", fieldName)
	case "email":
		return fmt.Sprintf("%s 必须是有效的邮箱地址", fieldName)
	case "min":
		return fmt.Sprintf("%s 长度不能少于 %s 个字符", fieldName, param)
	case "max":
		return fmt.Sprintf("%s 长度不能超过 %s 个字符", fieldName, param)
	case "len":
		return fmt.Sprintf("%s 长度必须为 %s 个字符", fieldName, param)
	case "numeric":
		return fmt.Sprintf("%s 必须是数字", fieldName)
	case "alpha":
		return fmt.Sprintf("%s 只能包含字母", fieldName)
	case "alphanum":
		return fmt.Sprintf("%s 只能包含字母和数字", fieldName)
	case "url":
		return fmt.Sprintf("%s 必须是有效的URL", fieldName)
	case "uri":
		return fmt.Sprintf("%s 必须是有效的URI", fieldName)
	case "gte":
		return fmt.Sprintf("%s 必须大于或等于 %s", fieldName, param)
	case "lte":
		return fmt.Sprintf("%s 必须小于或等于 %s", fieldName, param)
	case "gt":
		return fmt.Sprintf("%s 必须大于 %s", fieldName, param)
	case "lt":
		return fmt.Sprintf("%s 必须小于 %s", fieldName, param)
	case "oneof":
		return fmt.Sprintf("%s 必须是以下值之一: %s", fieldName, param)
	case "unique":
		return fmt.Sprintf("%s 已存在", fieldName)
	default:
		return fmt.Sprintf("%s 验证失败", fieldName)
	}
}

// SafeExecute 安全执行函数，捕获panic并转换为错误
func SafeExecute(fn func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = ErrInternalError.WithDetailsf("panic recovered: %v", r)
		}
	}()
	return fn()
}

// SafeExecuteWithResult 安全执行函数并返回结果，捕获panic并转换为错误
func SafeExecuteWithResult[T any](fn func() (T, error)) (result T, err error) {
	defer func() {
		if r := recover(); r != nil {
			var zero T
			result = zero
			err = ErrInternalError.WithDetailsf("panic recovered: %v", r)
		}
	}()
	return fn()
}

// Chain 链式错误处理，如果任何一个函数返回错误，则停止执行
func Chain(fns ...func() error) error {
	for _, fn := range fns {
		if err := fn(); err != nil {
			return err
		}
	}
	return nil
}

// Parallel 并行执行多个函数，收集所有错误
func Parallel(fns ...func() error) []error {
	ch := make(chan error, len(fns))
	
	for _, fn := range fns {
		go func(f func() error) {
			ch <- f()
		}(fn)
	}
	
	var errs []error
	for i := 0; i < len(fns); i++ {
		if err := <-ch; err != nil {
			errs = append(errs, err)
		}
	}
	
	return errs
}

// FirstError 返回第一个非nil错误
func FirstError(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

// CombineErrors 合并多个错误为一个错误
func CombineErrors(errs ...error) error {
	var validErrs []error
	for _, err := range errs {
		if err != nil {
			validErrs = append(validErrs, err)
		}
	}
	
	if len(validErrs) == 0 {
		return nil
	}
	
	if len(validErrs) == 1 {
		return validErrs[0]
	}
	
	var messages []string
	for _, err := range validErrs {
		messages = append(messages, err.Error())
	}
	
	return ErrInternalError.WithDetails(strings.Join(messages, "; "))
}

// RetryOnError 在错误时重试执行函数
func RetryOnError(fn func() error, maxRetries int) error {
	var lastErr error
	for i := 0; i <= maxRetries; i++ {
		if err := fn(); err != nil {
			lastErr = err
			if i == maxRetries {
				return WrapErrorf(err, "failed after %d retries", maxRetries)
			}
			continue
		}
		return nil
	}
	return lastErr
}