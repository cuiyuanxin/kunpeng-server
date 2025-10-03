# Go-i18n/v2 与 Validator 融合指南

本指南介绍如何在现有的 Kunpeng 项目中集成 go-i18n/v2 来实现更强大的国际化支持。

## 概述

我们已经成功将 go-i18n/v2 与现有的 validator 系统进行了融合，提供了两种验证方案：

1. **原有方案**：使用 `constants` 包存储错误消息
2. **go-i18n 方案**：使用 YAML 文件和 go-i18n/v2 库

## 项目结构

```
pkg/
├── i18n/                    # 新增的 i18n 包
│   ├── i18n.go             # i18n 核心功能
│   └── locales/            # 翻译文件目录
│       ├── zh.yaml         # 中文翻译
│       └── en.yaml         # 英文翻译
├── validator/
│   ├── validator.go        # 原有验证器
│   └── validator_i18n.go   # 新增的 i18n 验证器
└── constants/
    └── regex.go            # 原有的常量和错误消息
```

## 核心功能

### 1. i18n 包 (`pkg/i18n/i18n.go`)

提供核心的国际化功能：

```go
// 初始化 i18n
i18n.Init()

// 设置语言
i18n.SetLanguage("zh") // 或 "en"

// 翻译消息
message := i18n.T("validator.username")

// 带字段名的翻译
message := i18n.TWithField("validator.username", "用户名")

// 从 Accept-Language 头解析语言
lang := i18n.GetLanguageFromAcceptLanguage("zh-CN,zh;q=0.9")
```

### 2. 增强的验证器 (`pkg/validator/validator_i18n.go`)

提供与 go-i18n 集成的验证功能：

```go
// 初始化 i18n 验证器
validator.InitI18n()

// 切换语言
validator.SwitchLanguageI18n("en")

// 翻译验证错误
errorMsg := validator.TranslateI18n(err)

// 绑定和验证（自动语言检测）
err := validator.BindAndValidateJSONI18n(c, &req)
```

### 3. 自定义验证方法增强

在 DTO 中添加了 go-i18n 支持：

```go
// 使用 go-i18n 的验证方法
func (req *UserLoginReq) ValidateWithI18n() error {
    return req.validateWithI18n()
}

// 带语言参数的验证
func (req *UserLoginReq) ValidateWithLanguage(language string) error {
    i18n.SetLanguage(language)
    return req.validateWithI18n()
}
```

## 翻译文件格式

### 中文翻译 (`pkg/i18n/locales/zh.yaml`)

```yaml
# 验证器错误消息
validator.required: "{{.Field}}是必填字段"
validator.username: "{{.Field}}格式不正确，应为5-20位的字母、数字或下划线"
validator.password: "{{.Field}}格式不正确，应包含大小写字母、数字和特殊字符，长度6-25位"
validator.mobile: "{{.Field}}格式不正确，请输入正确的手机号码"

# 正则验证错误消息
regex:
  username: "用户名必须是5-20位的字母、数字或下划线组合"
  password: "密码必须包含大小写字母、数字和特殊字符，长度6-25位"
  mobile: "手机号格式不正确"
  captcha: "验证码必须是6位数字"
  email: "邮箱格式不正确"
  idcard: "身份证号格式不正确"
  ipv4: "IP地址格式不正确"
  url: "URL格式不正确"
  required: "该字段不能为空"
  invalid_format: "格式不正确"

# 业务相关错误消息
user.username_required: "用户名是必填字段"
user.password_required: "密码是必填字段"
```

### 英文翻译 (`pkg/i18n/locales/en.yaml`)

```yaml
# Validator error messages
validator.required: "{{.Field}} is required"
validator.username: "{{.Field}} format is incorrect, should be 5-20 characters of letters, numbers or underscores"
validator.password: "{{.Field}} format is incorrect, should contain uppercase and lowercase letters, numbers and special characters, 6-25 characters long"
validator.mobile: "{{.Field}} format is incorrect, please enter a valid mobile number"

# Business related error messages
user.username_required: "Username is required"
user.password_required: "Password is required"
```

## 使用方式

### 方式一：在控制器中使用

```go
func (uc *UserController) Login(c *gin.Context) {
    var req dto.UserLoginReq
    
    // 使用 go-i18n 验证器（自动语言检测）
    if err := validator.BindAndValidateJSONI18n(c, &req); err != nil {
        response.Error(c, err)
        return
    }
    
    // 或者手动调用自定义验证
    lang := validator.GetLanguageFromContextI18n(c)
    if err := req.ValidateWithLanguage(lang); err != nil {
        response.Error(c, kperrors.NewWithMessage(kperrors.ErrParam, err.Error(), err))
        return
    }
    
    // 业务逻辑...
}
```

### 方式二：直接使用 i18n 包

```go
// 设置语言
i18n.SetLanguage("en")

// 获取翻译消息
errorMsg := i18n.TWithField("validator.username", "username")
// 输出: "username format is incorrect, should be 5-20 characters of letters, numbers or underscores"
```

### 方式三：正则验证错误消息

```go
// 使用 go-i18n 获取正则验证错误消息
import "github.com/cuiyuanxin/kunpeng/pkg/constants"

// 获取当前语言的错误消息
message := constants.GetErrorMessageI18n("username")

// 获取指定语言的错误消息
message := constants.GetErrorMessageI18nWithLanguage("username", "en")

// 支持的错误类型:
// username, password, mobile, captcha, email, idcard, ipv4, url, required, invalid_format
```

## 迁移策略

### 1. 渐进式迁移

- 保留原有的 `validator.go` 和 `constants` 包
- 新功能使用 go-i18n 方案
- 逐步将现有功能迁移到 go-i18n

### 2. 并行使用

```go
// 原有方案
err := validator.BindAndValidateJSON(c, &req)

// go-i18n 方案
err := validator.BindAndValidateJSONI18n(c, &req)
```

### 3. 完全替换

- 将所有验证逻辑迁移到 go-i18n
- 移除 `constants` 包中的错误消息
- 统一使用 YAML 翻译文件

## 优势对比

| 特性 | 原有方案 | go-i18n 方案 |
|------|----------|-------------|
| 翻译文件格式 | Go 常量 | YAML 文件 |
| 模板变量支持 | 简单字符串替换 | 完整模板支持 |
| 复数规则 | 不支持 | 支持 |
| 文件嵌入 | 编译时嵌入 | embed.FS 嵌入 |
| 维护性 | 中等 | 高 |
| 扩展性 | 有限 | 强 |
| 学习成本 | 低 | 中等 |
| 适用场景 | 小型项目 | 大型项目 |

## 最佳实践

### 1. 翻译文件组织

```yaml
# 按功能模块组织
validator.username: "用户名格式错误"
validator.password: "密码格式错误"

user.login_failed: "登录失败"
user.permission_denied: "权限不足"

order.not_found: "订单不存在"
order.status_invalid: "订单状态无效"
```

### 2. 错误消息设计

```yaml
# 使用模板变量
validator.min_length: "{{.Field}}长度不能少于{{.Min}}个字符"
validator.max_length: "{{.Field}}长度不能超过{{.Max}}个字符"

# 支持复数规则
item.count:
  one: "有{{.Count}}个项目"
  other: "有{{.Count}}个项目"
```

### 3. 性能优化

- 在应用启动时初始化 i18n
- 缓存翻译结果
- 使用嵌入式文件系统减少 I/O

## 示例代码

完整的示例代码请参考 `example_go_i18n_integration.go` 文件，其中包含：

1. 自定义验证示例
2. 结构体验证示例
3. HTTP 请求验证示例
4. 方案对比说明

## 总结

go-i18n/v2 与 validator 的融合为项目提供了更强大和灵活的国际化支持。通过渐进式迁移策略，可以在不影响现有功能的前提下，逐步享受 go-i18n 带来的优势。

建议：
- 新项目直接使用 go-i18n 方案
- 现有项目可以并行使用两种方案
- 根据项目规模和需求选择合适的迁移策略