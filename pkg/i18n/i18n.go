package i18n

import (
	"embed"
	"fmt"
	"strings"
	"sync"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

//go:embed locales/*.yaml
var localeFS embed.FS

var (
	bundle      *i18n.Bundle
	localizer   *i18n.Localizer
	once        sync.Once
	currentLang string = "zh"
	mu          sync.RWMutex
)

// Init 初始化国际化
func Init() {
	once.Do(func() {
		bundle = i18n.NewBundle(language.Chinese)
		bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)

		// 加载语言文件
		_, err := bundle.LoadMessageFileFS(localeFS, "locales/zh.yaml")
		if err != nil {
			fmt.Printf("Failed to load zh.yaml: %v\n", err)
		}

		_, err = bundle.LoadMessageFileFS(localeFS, "locales/en.yaml")
		if err != nil {
			fmt.Printf("Failed to load en.yaml: %v\n", err)
		}

		// 默认使用中文
		localizer = i18n.NewLocalizer(bundle, "zh")
	})
}

// SetLanguage 设置当前语言
func SetLanguage(lang string) {
	mu.Lock()
	defer mu.Unlock()

	if bundle == nil {
		Init()
	}

	// 标准化语言代码
	switch {
	case strings.Contains(strings.ToLower(lang), "en"):
		lang = "en"
	default:
		lang = "zh"
	}

	currentLang = lang
	localizer = i18n.NewLocalizer(bundle, lang)
}

// GetCurrentLanguage 获取当前语言
func GetCurrentLanguage() string {
	mu.RLock()
	defer mu.RUnlock()
	return currentLang
}

// T 翻译消息
func T(messageID string, templateData ...map[string]interface{}) string {
	mu.RLock()
	defer mu.RUnlock()

	if localizer == nil {
		Init()
	}

	cfg := &i18n.LocalizeConfig{
		MessageID: messageID,
	}

	if len(templateData) > 0 {
		cfg.TemplateData = templateData[0]
	}

	result, err := localizer.Localize(cfg)
	if err != nil {
		// 如果翻译失败，返回消息ID
		return messageID
	}

	return result
}

// TWithField 翻译带字段名的消息
func TWithField(messageID string, fieldName string) string {
	return T(messageID, map[string]interface{}{
		"Field": fieldName,
	})
}

// GetLanguageFromAcceptLanguage 从 Accept-Language 头解析语言
func GetLanguageFromAcceptLanguage(acceptLang string) string {
	if acceptLang == "" {
		return "zh"
	}
	if strings.Contains(strings.ToLower(acceptLang), "en") {
		return "en"
	}
	return "zh"
}
