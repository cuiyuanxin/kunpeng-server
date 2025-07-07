package app

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	val "github.com/go-playground/validator/v10"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/cuiyuanxin/kunpeng/pkg/response"
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

// FormValidator 表单验证器
type FormValidator struct {
	validator *val.Validate
	trans     ut.Translator
	language  string
}

// NewFormValidator 创建新的表单验证器实例
func NewFormValidator() *FormValidator {
	validator := val.New()
	
	// 创建翻译器
	uni := ut.New(en.New(), zh.New())
	trans, _ := uni.GetTranslator("zh")
	
	// 注册翻译器
	zh_translations.RegisterDefaultTranslations(validator, trans)
	
	return &FormValidator{
		validator: validator,
		trans:     trans,
		language:  "zh",
	}
}

// SetLanguage 设置语言
func (f *FormValidator) SetLanguage(lang string) error {
	uni := ut.New(en.New(), zh.New())
	trans, found := uni.GetTranslator(lang)
	if !found {
		return nil // 使用默认语言
	}
	
	f.trans = trans
	f.language = lang
	
	// 根据语言注册翻译
	switch lang {
	case "en":
		en_translations.RegisterDefaultTranslations(f.validator, trans)
	case "zh":
		zh_translations.RegisterDefaultTranslations(f.validator, trans)
	}
	
	return nil
}

// GetLanguage 获取当前语言
func (f *FormValidator) GetLanguage() string {
	return f.language
}

// ShouldBind 绑定请求参数（不进行验证）
func (f *FormValidator) ShouldBind(c *gin.Context, obj interface{}) error {
	return c.ShouldBind(obj)
}

// ShouldBindJSON 绑定JSON参数（不进行验证）
func (f *FormValidator) ShouldBindJSON(c *gin.Context, obj interface{}) error {
	return c.ShouldBindJSON(obj)
}

// ShouldBindQuery 绑定查询参数（不进行验证）
func (f *FormValidator) ShouldBindQuery(c *gin.Context, obj interface{}) error {
	return c.ShouldBindQuery(obj)
}

// ShouldBindUri 绑定URI参数（不进行验证）
func (f *FormValidator) ShouldBindUri(c *gin.Context, obj interface{}) error {
	return c.ShouldBindUri(obj)
}

// ShouldBindHeader 绑定Header参数（不进行验证）
func (f *FormValidator) ShouldBindHeader(c *gin.Context, obj interface{}) error {
	return c.ShouldBindHeader(obj)
}

// Validate 验证结构体
func (f *FormValidator) Validate(obj interface{}) ValidErrors {
	err := f.validator.Struct(obj)
	if err != nil {
		return f.translateErrors(err)
	}
	return nil
}

// ValidateVar 验证单个变量
func (f *FormValidator) ValidateVar(field interface{}, tag string) ValidErrors {
	err := f.validator.Var(field, tag)
	if err != nil {
		return f.translateErrors(err)
	}
	return nil
}

// BindAndValidate 绑定并验证请求参数
func (f *FormValidator) BindAndValidate(c *gin.Context, obj interface{}) ValidErrors {
	// 先绑定参数
	if err := c.ShouldBind(obj); err != nil {
		return f.translateErrors(err)
	}
	
	// 再验证参数
	return f.Validate(obj)
}

// BindJSONAndValidate 绑定JSON并验证请求参数
func (f *FormValidator) BindJSONAndValidate(c *gin.Context, obj interface{}) ValidErrors {
	// 先绑定JSON参数
	if err := c.ShouldBindJSON(obj); err != nil {
		return f.translateErrors(err)
	}
	
	// 再验证参数
	return f.Validate(obj)
}

// BindQueryAndValidate 绑定查询参数并验证
func (f *FormValidator) BindQueryAndValidate(c *gin.Context, obj interface{}) ValidErrors {
	// 先绑定查询参数
	if err := c.ShouldBindQuery(obj); err != nil {
		return f.translateErrors(err)
	}
	
	// 再验证参数
	return f.Validate(obj)
}

// translateErrors 翻译错误信息
func (f *FormValidator) translateErrors(err error) ValidErrors {
	var errs ValidErrors
	
	// 检查是否为验证错误
	verrs, ok := err.(val.ValidationErrors)
	if !ok {
		// 非验证错误，直接返回
		errs = append(errs, &ValidError{
			Key:     "bind_error",
			Message: err.Error(),
		})
		return errs
	}
	
	// 翻译验证错误
	for key, value := range verrs.Translate(f.trans) {
		errs = append(errs, &ValidError{
			Key:     key,
			Message: value,
		})
	}
	
	return errs
}

// RegisterValidation 注册自定义验证规则
func (f *FormValidator) RegisterValidation(tag string, fn val.Func) error {
	return f.validator.RegisterValidation(tag, fn)
}

// RegisterTranslation 注册自定义翻译
func (f *FormValidator) RegisterTranslation(tag string, text string) error {
	return f.validator.RegisterTranslation(tag, f.trans, func(ut ut.Translator) error {
		return ut.Add(tag, text, true)
	}, func(ut ut.Translator, fe val.FieldError) string {
		t, _ := ut.T(tag, fe.Field())
		return t
	})
}

// GetValidator 获取底层验证器实例
func (f *FormValidator) GetValidator() *val.Validate {
	return f.validator
}

// GetTranslator 获取翻译器实例
func (f *FormValidator) GetTranslator() ut.Translator {
	return f.trans
}

// FieldLevel 类型别名，方便外部使用
type FieldLevel = val.FieldLevel

// HandleValidationErrors 处理验证错误并返回响应
func (f *FormValidator) HandleValidationErrors(c *gin.Context, errs ValidErrors) {
	if len(errs) > 0 {
		response.Error(c, http.StatusBadRequest, "参数验证失败", errs)
	}
}

// HandleValidationErrorsWithMessage 处理验证错误并返回自定义消息
func (f *FormValidator) HandleValidationErrorsWithMessage(c *gin.Context, errs ValidErrors, message string) {
	if len(errs) > 0 {
		response.Error(c, http.StatusBadRequest, message, errs)
	}
}

// MustBindAndValidate 绑定并验证，如果失败则自动返回错误响应
func (f *FormValidator) MustBindAndValidate(c *gin.Context, obj interface{}) bool {
	errs := f.BindAndValidate(c, obj)
	if len(errs) > 0 {
		f.HandleValidationErrors(c, errs)
		return false
	}
	return true
}

// MustBindJSONAndValidate 绑定JSON并验证，如果失败则自动返回错误响应
func (f *FormValidator) MustBindJSONAndValidate(c *gin.Context, obj interface{}) bool {
	errs := f.BindJSONAndValidate(c, obj)
	if len(errs) > 0 {
		f.HandleValidationErrors(c, errs)
		return false
	}
	return true
}

// MustBindQueryAndValidate 绑定查询参数并验证，如果失败则自动返回错误响应
func (f *FormValidator) MustBindQueryAndValidate(c *gin.Context, obj interface{}) bool {
	errs := f.BindQueryAndValidate(c, obj)
	if len(errs) > 0 {
		f.HandleValidationErrors(c, errs)
		return false
	}
	return true
}