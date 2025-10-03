package validator

import (
	"reflect"
	"regexp"
	"strings"
	"sync"

	"github.com/cuiyuanxin/kunpeng/pkg/constants"
	kperrors "github.com/cuiyuanxin/kunpeng/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

var (
	trans    ut.Translator
	validate *validator.Validate
	once     sync.Once
)

// Init 初始化验证器
func Init() {
	once.Do(func() {
		// 获取gin的验证器
		validate = binding.Validator.Engine().(*validator.Validate)

		// 注册一个获取json tag的自定义方法
		validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		// 初始化翻译器
		zhT := zh.New()
		enT := en.New()
		uni := ut.New(enT, zhT, enT)

		// 默认使用中文
		trans, _ = uni.GetTranslator("zh")

		// 注册翻译器
		_ = zh_translations.RegisterDefaultTranslations(validate, trans)

		// 注册自定义验证规则
		registerCustomValidations()
	})
}

// registerCustomValidations 注册自定义验证规则
func registerCustomValidations() {
	// 注册用户名验证规则
	_ = RegisterCustomValidation("username", validateUsername, "{0}格式不正确")
	// 注册密码验证规则
	_ = RegisterCustomValidation("password", validatePassword, "{0}格式不正确")
	// 注册手机号验证规则
	_ = RegisterCustomValidation("mobile", validateMobile, "{0}格式不正确")
}

// validateUsername 验证用户名格式
func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	// 验证用户名格式：5-20位的字母、数字或下划线
	matched, _ := regexp.MatchString(constants.UsernameRegex, username)
	return matched
}

// validatePassword 验证密码格式
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	// 验证密码格式：包含大小写字母、数字和特殊字符，长度6-25位
	matched, _ := regexp.MatchString(constants.PasswordRegex, password)
	return matched
}

// validateMobile 验证手机号格式
func validateMobile(fl validator.FieldLevel) bool {
	mobile := fl.Field().String()
	// 验证中国大陆手机号格式：1开头，第二位为3-9，总共11位数字
	matched, _ := regexp.MatchString(constants.MobileRegex, mobile)
	return matched
}

// Translate 翻译错误信息
func Translate(err error) string {
	if err == nil {
		return ""
	}

	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return err.Error()
	}

	var errMsgs []string
	for _, e := range errs {
		errMsgs = append(errMsgs, e.Translate(trans))
	}

	return strings.Join(errMsgs, ", ")
}

// RegisterCustomValidation 注册自定义验证规则
func RegisterCustomValidation(tag string, fn validator.Func, errMsg string) error {
	if err := validate.RegisterValidation(tag, fn); err != nil {
		return err
	}
	return RegisterCustomTranslation(tag, errMsg)
}

// RegisterCustomTranslation 注册自定义翻译
func RegisterCustomTranslation(tag string, errMsg string) error {
	return validate.RegisterTranslation(tag, trans, func(ut ut.Translator) error {
		return ut.Add(tag, errMsg, true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(tag, fe.Field())
		return t
	})
}

// SwitchLanguage 切换语言
func SwitchLanguage(language string) {
	if validate == nil {
		Init()
	}

	uni := ut.New(en.New(), zh.New())
	var ok bool
	trans, ok = uni.GetTranslator(language)
	if !ok {
		trans, _ = uni.GetTranslator("zh")
	}

	switch language {
	case "en":
		_ = en_translations.RegisterDefaultTranslations(validate, trans)
	default:
		_ = zh_translations.RegisterDefaultTranslations(validate, trans)
	}
}

// GetCurrentLanguage 获取当前语言
func GetCurrentLanguage() string {
	if trans == nil {
		return "zh"
	}
	// 通过translator的locale获取当前语言
	locale := trans.Locale()
	if locale == "en" {
		return "en"
	}
	return "zh"
}

// GetLanguageFromContext 从gin上下文获取语言设置
func GetLanguageFromContext(c *gin.Context) string {
	lang := c.GetHeader("Accept-Language")
	if lang == "" {
		return "zh"
	}
	if strings.Contains(strings.ToLower(lang), "en") {
		return "en"
	}
	return "zh"
}

// BindAndValidate 绑定并验证请求参数
func BindAndValidate(c *gin.Context, obj interface{}) error {
	// 根据请求头切换语言
	lang := c.GetHeader("Accept-Language")
	if lang != "" {
		SwitchLanguage(lang)
	}

	// 绑定参数
	if err := c.ShouldBind(obj); err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			return kperrors.NewWithMessage(kperrors.ErrParam, Translate(err), err)
		}
		return kperrors.New(kperrors.ErrParam, err)
	}

	return nil
}

// BindAndValidateUri 绑定并验证URI参数
func BindAndValidateUri(c *gin.Context, obj interface{}) error {
	// 根据请求头切换语言
	lang := c.GetHeader("Accept-Language")
	if lang != "" {
		SwitchLanguage(lang)
	}

	// 绑定参数
	if err := c.ShouldBindUri(obj); err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			return kperrors.NewWithMessage(kperrors.ErrParam, Translate(err), err)
		}
		return kperrors.New(kperrors.ErrParam, err)
	}

	return nil
}

// BindAndValidateQuery 绑定并验证Query参数
func BindAndValidateQuery(c *gin.Context, obj interface{}) error {
	// 根据请求头切换语言
	lang := c.GetHeader("Accept-Language")
	if lang != "" {
		SwitchLanguage(lang)
	}

	// 绑定参数
	if err := c.ShouldBindQuery(obj); err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			return kperrors.NewWithMessage(kperrors.ErrParam, Translate(err), err)
		}
		return kperrors.New(kperrors.ErrParam, err)
	}

	return nil
}

// BindAndValidateJSON 绑定并验证JSON参数
func BindAndValidateJSON(c *gin.Context, obj interface{}) error {
	// 根据请求头切换语言
	lang := c.GetHeader("Accept-Language")
	if lang != "" {
		SwitchLanguage(lang)
	}

	// 绑定参数
	if err := c.ShouldBindJSON(obj); err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			return kperrors.NewWithMessage(kperrors.ErrParam, Translate(err), err)
		}
		return kperrors.New(kperrors.ErrParam, err)
	}

	return nil
}
