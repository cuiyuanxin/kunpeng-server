# Form Validation Package

这是一个专门用于处理表单验证的包，将参数绑定和验证逻辑从控制器中分离出来，提供更灵活和可维护的验证方案。

## 特性

- **分离关注点**: 将 `c.ShouldBind` 绑定逻辑和验证逻辑分离
- **多语言支持**: 支持中文和英文错误消息
- **灵活绑定**: 支持 JSON、表单、查询参数、URI、Header 等多种绑定方式
- **自定义验证**: 支持注册自定义验证规则和翻译
- **友好错误**: 结构化的错误信息，便于前端处理
- **便捷方法**: 提供 Must 系列方法，自动处理错误响应

## 核心概念

### 分离绑定和验证

传统方式将绑定和验证耦合在一起，新的设计将它们分离：

```go
// 传统方式 - 绑定和验证耦合
if err := c.ShouldBindJSON(&req); err != nil {
    // 处理绑定错误和验证错误
}

// 新方式 - 分离绑定和验证
// 1. 仅绑定
if err := formValidator.ShouldBindJSON(c, &req); err != nil {
    // 处理绑定错误（JSON格式错误等）
}

// 2. 业务逻辑处理（设置默认值等）
if req.Page == 0 {
    req.Page = 1
}

// 3. 验证
errs := formValidator.Validate(&req)
if len(errs) > 0 {
    // 处理验证错误
}
```

## 快速开始

### 1. 创建验证器实例

```go
import "github.com/cuiyuanxin/kunpeng/pkg/app"

formValidator := app.NewFormValidator()
```

### 2. 在控制器中使用

```go
type UserRequest struct {
    Username string `json:"username" validate:"required,min=3,max=50"`
    Email    string `json:"email" validate:"required,email"`
    Age      int    `json:"age" validate:"required,min=1,max=120"`
}

// 方式1: 使用 Must 方法（推荐）
func CreateUser(c *gin.Context) {
    var req UserRequest
    
    // 自动处理错误响应
    if !formValidator.MustBindJSONAndValidate(c, &req) {
        return // 验证失败，已自动返回错误响应
    }
    
    // 验证成功，处理业务逻辑
    response.Success(c, req)
}

// 方式2: 手动处理错误
func UpdateUser(c *gin.Context) {
    var req UserRequest
    
    errs := formValidator.BindJSONAndValidate(c, &req)
    if len(errs) > 0 {
        formValidator.HandleValidationErrors(c, errs)
        return
    }
    
    // 验证成功，处理业务逻辑
    response.Success(c, req)
}

// 方式3: 分离绑定和验证
func CreateProduct(c *gin.Context) {
    var req ProductRequest
    
    // 先绑定
    if err := formValidator.ShouldBindJSON(c, &req); err != nil {
        response.BadRequest(c, "JSON格式错误")
        return
    }
    
    // 业务逻辑处理
    if req.Description == "" {
        req.Description = "暂无描述"
    }
    
    // 再验证
    errs := formValidator.Validate(&req)
    if len(errs) > 0 {
        formValidator.HandleValidationErrors(c, errs)
        return
    }
    
    response.Success(c, req)
}
```

## API 文档

### FormValidator 结构体

#### 创建实例

```go
func NewFormValidator() *FormValidator
```

#### 语言设置

```go
// 设置语言（zh/en）
func (f *FormValidator) SetLanguage(lang string) error

// 获取当前语言
func (f *FormValidator) GetLanguage() string
```

#### 仅绑定方法（不验证）

```go
// 绑定表单参数
func (f *FormValidator) ShouldBind(c *gin.Context, obj interface{}) error

// 绑定JSON参数
func (f *FormValidator) ShouldBindJSON(c *gin.Context, obj interface{}) error

// 绑定查询参数
func (f *FormValidator) ShouldBindQuery(c *gin.Context, obj interface{}) error

// 绑定URI参数
func (f *FormValidator) ShouldBindUri(c *gin.Context, obj interface{}) error

// 绑定Header参数
func (f *FormValidator) ShouldBindHeader(c *gin.Context, obj interface{}) error
```

#### 仅验证方法（不绑定）

```go
// 验证结构体
func (f *FormValidator) Validate(obj interface{}) ValidErrors

// 验证单个变量
func (f *FormValidator) ValidateVar(field interface{}, tag string) ValidErrors
```

#### 绑定并验证方法

```go
// 绑定表单并验证
func (f *FormValidator) BindAndValidate(c *gin.Context, obj interface{}) ValidErrors

// 绑定JSON并验证
func (f *FormValidator) BindJSONAndValidate(c *gin.Context, obj interface{}) ValidErrors

// 绑定查询参数并验证
func (f *FormValidator) BindQueryAndValidate(c *gin.Context, obj interface{}) ValidErrors
```

#### Must 系列方法（自动处理错误响应）

```go
// 绑定表单并验证，失败时自动返回错误响应
func (f *FormValidator) MustBindAndValidate(c *gin.Context, obj interface{}) bool

// 绑定JSON并验证，失败时自动返回错误响应
func (f *FormValidator) MustBindJSONAndValidate(c *gin.Context, obj interface{}) bool

// 绑定查询参数并验证，失败时自动返回错误响应
func (f *FormValidator) MustBindQueryAndValidate(c *gin.Context, obj interface{}) bool
```

#### 错误处理方法

```go
// 处理验证错误
func (f *FormValidator) HandleValidationErrors(c *gin.Context, errs ValidErrors)

// 处理验证错误并返回自定义消息
func (f *FormValidator) HandleValidationErrorsWithMessage(c *gin.Context, errs ValidErrors, message string)
```

#### 自定义验证

```go
// 注册自定义验证规则
func (f *FormValidator) RegisterValidation(tag string, fn val.Func) error

// 注册自定义翻译
func (f *FormValidator) RegisterTranslation(tag string, text string) error
```

### ValidError 和 ValidErrors

```go
type ValidError struct {
    Key     string `json:"key"`     // 字段名
    Message string `json:"message"` // 错误消息
}

type ValidErrors []*ValidError

// 实现 error 接口
func (v *ValidError) Error() string
func (v ValidErrors) Error() string

// 获取所有错误消息
func (v ValidErrors) Errors() []string
```

## 在 BaseController 中使用

`BaseController` 已经集成了 `FormValidator`，可以直接使用：

```go
type MyController struct {
    *controller.BaseController
}

func (mc *MyController) CreateUser(c *gin.Context) {
    var req UserRequest
    
    // 使用 BaseController 的方法
    if !mc.MustBindJSONAndValidate(c, &req) {
        return
    }
    
    // 业务逻辑
    response.Success(c, req)
}

// 或者分离绑定和验证
func (mc *MyController) UpdateUser(c *gin.Context) {
    var req UserRequest
    
    // 先绑定
    if err := mc.ShouldBindJSON(c, &req); err != nil {
        response.BadRequest(c, "JSON格式错误")
        return
    }
    
    // 业务逻辑处理
    // ...
    
    // 再验证
    errs := mc.Validate(&req)
    if len(errs) > 0 {
        mc.HandleValidationErrors(c, errs)
        return
    }
    
    response.Success(c, req)
}
```

## 多语言支持

### 设置语言

```go
// 通过代码设置
formValidator.SetLanguage("en") // 英文
formValidator.SetLanguage("zh") // 中文

// 通过HTTP头部设置
// Accept-Language: en
// Accept-Language: zh

// 通过查询参数设置
// ?lang=en
// ?lang=zh
```

### 中间件自动切换语言

```go
r.Use(func(c *gin.Context) {
    lang := c.GetHeader("Accept-Language")
    if lang == "" {
        lang = c.Query("lang")
    }
    if lang != "" {
        formValidator.SetLanguage(lang)
    }
    c.Next()
})
```

## 自定义验证

### 注册自定义验证规则

```go
// 注册验证规则
formValidator.RegisterValidation("custom_username", func(fl app.FieldLevel) bool {
    username := fl.Field().String()
    // 自定义验证逻辑
    return isValidUsername(username)
})

// 注册翻译
formValidator.RegisterTranslation("custom_username", "用户名格式不正确")

// 使用自定义验证
type UserRequest struct {
    Username string `json:"username" validate:"required,custom_username"`
}
```

## 常用验证标签

- `required` - 必填
- `min=3` - 最小长度/值
- `max=50` - 最大长度/值
- `len=11` - 固定长度
- `email` - 邮箱格式
- `url` - URL格式
- `numeric` - 数字
- `alpha` - 字母
- `alphanum` - 字母数字
- `oneof=red green blue` - 枚举值
- `gt=0` - 大于
- `gte=0` - 大于等于
- `lt=100` - 小于
- `lte=100` - 小于等于
- `omitempty` - 可选字段

## 错误响应格式

```json
{
    "code": 400,
    "message": "参数验证失败",
    "data": [
        {
            "key": "Username",
            "message": "Username长度必须至少为3个字符"
        },
        {
            "key": "Email",
            "message": "Email必须是一个有效的邮箱"
        }
    ]
}
```

## 最佳实践

1. **使用 Must 方法**: 对于简单的验证场景，使用 `MustBindJSONAndValidate` 等方法
2. **分离复杂逻辑**: 对于需要业务逻辑处理的场景，分离绑定和验证
3. **统一错误处理**: 使用 `HandleValidationErrors` 方法统一处理错误响应
4. **自定义验证**: 为业务特定的验证需求注册自定义验证规则
5. **多语言支持**: 根据用户偏好设置合适的语言

## 与旧版本的区别

| 特性 | 旧版本 (pkg/validator) | 新版本 (pkg/app/form) |
|------|----------------------|----------------------|
| 绑定和验证 | 耦合在一起 | 可以分离 |
| 错误处理 | 手动处理 | 提供自动处理方法 |
| 方法命名 | BindAndValid | BindAndValidate |
| 返回值 | (bool, ValidErrors) | ValidErrors |
| Must 方法 | 无 | 提供 Must 系列方法 |
| 单字段验证 | 无 | 支持 ValidateVar |
| 多种绑定 | 有限支持 | 全面支持 |

## 示例项目

查看 `examples/form_example.go` 和 `cmd/form_example/main.go` 获取完整的使用示例。

运行示例：

```bash
go run cmd/form_example/main.go
```

然后访问 `http://localhost:8082` 查看所有可用的API接口。