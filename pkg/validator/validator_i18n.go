package validator

import (
	"reflect"
	"regexp"
	"strings"
	"sync"

	"github.com/cuiyuanxin/kunpeng/pkg/constants"
	kperrors "github.com/cuiyuanxin/kunpeng/pkg/errors"
	"github.com/cuiyuanxin/kunpeng/pkg/i18n"
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
	i18nTrans    ut.Translator
	i18nValidate *validator.Validate
	i18nOnce     sync.Once
)

// InitI18n 初始化带 go-i18n 支持的验证器
func InitI18n() {
	i18nOnce.Do(func() {
		// 初始化 i18n
		i18n.Init()

		// 获取gin的验证器
		i18nValidate = binding.Validator.Engine().(*validator.Validate)

		// 注册一个获取json tag的自定义方法
		i18nValidate.RegisterTagNameFunc(func(fld reflect.StructField) string {
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
		i18nTrans, _ = uni.GetTranslator("zh")

		// 注册翻译器
		_ = zh_translations.RegisterDefaultTranslations(i18nValidate, i18nTrans)

		// 注册自定义验证规则
		registerI18nCustomValidations()
	})
}

// registerI18nCustomValidations 注册自定义验证规则（使用 go-i18n）
func registerI18nCustomValidations() {
	// 注册用户名验证规则
	_ = RegisterI18nCustomValidation("username", validateUsername, "validator.username")
	// 注册密码验证规则
	_ = RegisterI18nCustomValidation("password", validatePassword, "validator.password")
	// 注册手机号验证规则
	_ = RegisterI18nCustomValidation("mobile", validateMobile, "validator.mobile")
	// 注册验证码验证规则
	_ = RegisterI18nCustomValidation("captcha", validateCaptcha, "validator.captcha")
}

// validateCaptcha 验证验证码格式
func validateCaptcha(fl validator.FieldLevel) bool {
	captcha := fl.Field().String()
	// 验证验证码格式：4位数字或字母
	matched, _ := regexp.MatchString(constants.CaptchaRegex, captcha)
	return matched
}

// TranslateI18n 使用 go-i18n 翻译错误信息
func TranslateI18n(err error) string {
	if err == nil {
		return ""
	}

	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return err.Error()
	}

	var errMsgs []string
	for _, e := range errs {
		// 首先尝试使用 go-i18n 翻译
		messageID := "validator." + e.Tag()
		translated := i18n.TWithField(messageID, e.Field())

		// 如果 go-i18n 翻译失败（返回的是 messageID），则使用原有的翻译器
		if translated == messageID {
			translated = e.Translate(i18nTrans)
		}

		errMsgs = append(errMsgs, translated)
	}

	return strings.Join(errMsgs, ", ")
}

// RegisterI18nCustomValidation 注册自定义验证规则（使用 go-i18n）
func RegisterI18nCustomValidation(tag string, fn validator.Func, messageID string) error {
	if err := i18nValidate.RegisterValidation(tag, fn); err != nil {
		return err
	}
	return RegisterI18nCustomTranslation(tag, messageID)
}

// RegisterI18nCustomTranslation 注册自定义翻译（使用 go-i18n）
func RegisterI18nCustomTranslation(tag string, messageID string) error {
	return i18nValidate.RegisterTranslation(tag, i18nTrans, func(ut ut.Translator) error {
		// 使用 go-i18n 获取翻译
		translated := i18n.T(messageID)
		return ut.Add(tag, translated, true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		// 使用 go-i18n 翻译，传递字段名
		return i18n.TWithField(messageID, fe.Field())
	})
}

// SwitchLanguageI18n 切换语言（使用 go-i18n）
func SwitchLanguageI18n(language string) {
	if i18nValidate == nil {
		InitI18n()
	}

	// 设置 i18n 语言
	i18n.SetLanguage(language)

	// 同时更新 validator 的翻译器
	uni := ut.New(en.New(), zh.New())
	var ok bool
	i18nTrans, ok = uni.GetTranslator(language)
	if !ok {
		i18nTrans, _ = uni.GetTranslator("zh")
	}

	switch language {
	case "en":
		_ = en_translations.RegisterDefaultTranslations(i18nValidate, i18nTrans)
	default:
		_ = zh_translations.RegisterDefaultTranslations(i18nValidate, i18nTrans)
	}

	// 重新注册自定义验证规则的翻译
	registerI18nCustomValidations()
}

// GetCurrentLanguageI18n 获取当前语言（使用 go-i18n）
func GetCurrentLanguageI18n() string {
	return i18n.GetCurrentLanguage()
}

// GetLanguageFromContextI18n 从gin上下文获取语言设置（使用 go-i18n）
func GetLanguageFromContextI18n(c *gin.Context) string {
	lang := c.GetHeader("Accept-Language")
	return i18n.GetLanguageFromAcceptLanguage(lang)
}

// BindAndValidateI18n 绑定并验证请求参数（使用 go-i18n）
func BindAndValidateI18n(c *gin.Context, obj interface{}) error {
	// 根据请求头切换语言
	lang := GetLanguageFromContextI18n(c)
	SwitchLanguageI18n(lang)

	// 绑定参数
	if err := c.ShouldBind(obj); err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			return kperrors.NewWithMessage(kperrors.ErrParam, TranslateI18n(err), err)
		}
		return kperrors.New(kperrors.ErrParam, err)
	}

	return nil
}

// BindAndValidateUriI18n 绑定并验证URI参数（使用 go-i18n）
func BindAndValidateUriI18n(c *gin.Context, obj interface{}) error {
	// 根据请求头切换语言
	lang := GetLanguageFromContextI18n(c)
	SwitchLanguageI18n(lang)

	// 绑定参数
	if err := c.ShouldBindUri(obj); err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			return kperrors.NewWithMessage(kperrors.ErrParam, TranslateI18n(err), err)
		}
		return kperrors.New(kperrors.ErrParam, err)
	}

	return nil
}

// BindAndValidateQueryI18n 绑定并验证Query参数（使用 go-i18n）
func BindAndValidateQueryI18n(c *gin.Context, obj interface{}) error {
	// 根据请求头切换语言
	lang := GetLanguageFromContextI18n(c)
	SwitchLanguageI18n(lang)

	// 绑定参数
	if err := c.ShouldBindQuery(obj); err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			return kperrors.NewWithMessage(kperrors.ErrParam, TranslateI18n(err), err)
		}
		return kperrors.New(kperrors.ErrParam, err)
	}

	return nil
}

// BindAndValidateJSONI18n 绑定并验证JSON参数（使用 go-i18n）
func BindAndValidateJSONI18n(c *gin.Context, obj interface{}) error {
	// 根据请求头切换语言
	lang := GetLanguageFromContextI18n(c)
	SwitchLanguageI18n(lang)

	// 绑定参数
	if err := c.ShouldBindJSON(obj); err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			return kperrors.NewWithMessage(kperrors.ErrParam, TranslateI18n(err), err)
		}
		return kperrors.New(kperrors.ErrParam, err)
	}

	return nil
}
