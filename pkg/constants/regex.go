package constants

import (
	"github.com/cuiyuanxin/kunpeng/pkg/i18n"
	"regexp"
)

// 正则表达式常量
const (
	// 用户名正则：5-20位的字母、数字或下划线
	UsernameRegex = `^[a-zA-Z0-9_]{5,20}$`

	// 密码正则：基本字符集和长度验证（6-25位，只包含字母、数字和指定特殊字符）
	PasswordRegex = `^[a-zA-Z\d~!@#$%^&_-]{6,25}$`

	// 手机号正则：中国大陆手机号格式，1开头，第二位为3-9，总共11位数字
	MobileRegex = `^1[3-9]\d{9}$`

	// 验证码正则：6位数字
	CaptchaRegex = `^\d{6}$`

	// 邮箱正则：标准邮箱格式
	EmailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	// 身份证号正则：18位身份证号
	IDCardRegex = `^[1-9]\d{5}(18|19|20)\d{2}((0[1-9])|(1[0-2]))(([0-2][1-9])|10|20|30|31)\d{3}[0-9Xx]$`

	// IP地址正则：IPv4格式
	IPv4Regex = `^((25[0-5]|2[0-4]\d|[01]?\d\d?)\.){3}(25[0-5]|2[0-4]\d|[01]?\d\d?)$`

	// URL正则：HTTP/HTTPS URL格式
	URLRegex = `^https?://[\w\-]+(\.[\w\-]+)+([\w\-\.,@?^=%&:/~\+#]*[\w\-\@?^=%&/~\+#])?$`
)

// GetErrorMessageI18n 使用 go-i18n 获取错误消息
func GetErrorMessageI18n(msgType string) string {
	return i18n.T("regex." + msgType)
}

// GetErrorMessageI18nWithLanguage 使用 go-i18n 根据指定语言获取错误消息
func GetErrorMessageI18nWithLanguage(msgType, language string) string {
	i18n.SetLanguage(language)
	return i18n.T("regex." + msgType)
}

// ValidatePassword 验证密码复杂性（Go 兼容版本）
// 要求：包含大小写字母、数字和特殊字符，长度6-25位
func ValidatePassword(password string) bool {
	if len(password) < 6 || len(password) > 25 {
		return false
	}

	// 基本字符集验证
	basicRegex := regexp.MustCompile(PasswordRegex)
	if !basicRegex.MatchString(password) {
		return false
	}

	// 手动检查复杂性要求
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`\d`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[~!@#$%^&_-]`).MatchString(password)

	return hasLower && hasUpper && hasDigit && hasSpecial
}
