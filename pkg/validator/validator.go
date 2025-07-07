package validator

import (
	"strings"

	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	val "github.com/go-playground/validator/v10"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

// ValidError 验证错误结构体
type ValidError struct {
	Key     string `json:"key"`
	Message string `json:"message"`
}

// ValidErrors 验证错误集合
type ValidErrors []*ValidError

// Error 实现error接口
func (v *ValidError) Error() string {
	return v.Message
}

// Error 实现error接口
func (v ValidErrors) Error() string {
	return strings.Join(v.Errors(), ",")
}

// Errors 获取所有错误信息
func (v ValidErrors) Errors() []string {
	var errs []string
	for _, err := range v {
		errs = append(errs, err.Error())
	}
	return errs
}

// Validator 验证器结构体
type Validator struct {
	validator *val.Validate
	trans     ut.Translator
}

// New 创建新的验证器实例
func New() *Validator {
	validator := val.New()
	
	// 创建翻译器
	uni := ut.New(en.New(), zh.New())
	trans, _ := uni.GetTranslator("zh")
	
	// 注册翻译器
	zh_translations.RegisterDefaultTranslations(validator, trans)
	
	return &Validator{
		validator: validator,
		trans:     trans,
	}
}

// SetLanguage 设置语言
func (v *Validator) SetLanguage(lang string) error {
	uni := ut.New(en.New(), zh.New())
	trans, found := uni.GetTranslator(lang)
	if !found {
		return nil // 使用默认语言
	}
	
	v.trans = trans
	
	// 根据语言注册翻译
	switch lang {
	case "en":
		en_translations.RegisterDefaultTranslations(v.validator, trans)
	case "zh":
		zh_translations.RegisterDefaultTranslations(v.validator, trans)
	}
	
	return nil
}

// BindAndValid 绑定并验证请求参数
func (v *Validator) BindAndValid(c *gin.Context, obj interface{}) (bool, ValidErrors) {
	var errs ValidErrors
	
	// 绑定参数
	err := c.ShouldBind(obj)
	if err != nil {
		// 检查是否为验证错误
		verrs, ok := err.(val.ValidationErrors)
		if !ok {
			// 非验证错误，直接返回
			errs = append(errs, &ValidError{
				Key:     "bind_error",
				Message: err.Error(),
			})
			return false, errs
		}
		
		// 翻译验证错误
		for key, value := range verrs.Translate(v.trans) {
			errs = append(errs, &ValidError{
				Key:     key,
				Message: value,
			})
		}
		
		return false, errs
	}
	
	return true, nil
}

// BindJSONAndValid 绑定JSON并验证请求参数
func (v *Validator) BindJSONAndValid(c *gin.Context, obj interface{}) (bool, ValidErrors) {
	var errs ValidErrors
	
	// 绑定JSON参数
	err := c.ShouldBindJSON(obj)
	if err != nil {
		// 检查是否为验证错误
		verrs, ok := err.(val.ValidationErrors)
		if !ok {
			// 非验证错误，直接返回
			errs = append(errs, &ValidError{
				Key:     "bind_error",
				Message: err.Error(),
			})
			return false, errs
		}
		
		// 翻译验证错误
		for key, value := range verrs.Translate(v.trans) {
			errs = append(errs, &ValidError{
				Key:     key,
				Message: value,
			})
		}
		
		return false, errs
	}
	
	return true, nil
}

// Validate 直接验证结构体
func (v *Validator) Validate(obj interface{}) (bool, ValidErrors) {
	var errs ValidErrors
	
	err := v.validator.Struct(obj)
	if err != nil {
		// 检查是否为验证错误
		verrs, ok := err.(val.ValidationErrors)
		if !ok {
			// 非验证错误，直接返回
			errs = append(errs, &ValidError{
				Key:     "validation_error",
				Message: err.Error(),
			})
			return false, errs
		}
		
		// 翻译验证错误
		for key, value := range verrs.Translate(v.trans) {
			errs = append(errs, &ValidError{
				Key:     key,
				Message: value,
			})
		}
		
		return false, errs
	}
	
	return true, nil
}

// RegisterValidation 注册自定义验证规则
func (v *Validator) RegisterValidation(tag string, fn val.Func) error {
	return v.validator.RegisterValidation(tag, fn)
}

// FieldLevel 类型别名，方便外部使用
type FieldLevel = val.FieldLevel

// RegisterTranslation 注册自定义翻译
func (v *Validator) RegisterTranslation(tag string, text string) error {
	return v.validator.RegisterTranslation(tag, v.trans, func(ut ut.Translator) error {
		return ut.Add(tag, text, true)
	}, func(ut ut.Translator, fe val.FieldError) string {
		t, _ := ut.T(tag, fe.Field())
		return t
	})
}

// GetValidator 获取底层验证器实例
func (v *Validator) GetValidator() *val.Validate {
	return v.validator
}

// GetTranslator 获取翻译器实例
func (v *Validator) GetTranslator() ut.Translator {
	return v.trans
}