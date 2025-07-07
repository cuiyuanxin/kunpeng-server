# Validator 表单验证包

这是一个基于 `github.com/go-playground/validator/v10` 的 Gin 表单验证包，提供了更友好的错误处理和多语言支持。

## 特性

- 🌍 **多语言支持**: 支持中文和英文错误信息
- 🔧 **灵活绑定**: 支持 JSON、表单、查询参数等多种绑定方式
- 📝 **友好错误**: 结构化的错误信息，便于前端处理
- 🎯 **自定义验证**: 支持注册自定义验证规则和翻译
- 🚀 **易于集成**: 与现有的 Gin 项目无缝集成

## 快速开始

### 1. 创建验证器实例

```go
import "github.com/cuiyuanxin/kunpeng/pkg/validator"

// 创建验证器
v := validator.New()

// 设置语言（可选，默认为中文）
v.SetLanguage("zh") // 或 "en"
```

### 2. 定义请求结构体

```go
type UserRequest struct {
    Username string `json:"username" validate:"required,min=3,max=50" label:"用户名"`
    Email    string `json:"email" validate:"required,email" label:"邮箱"`
    Password string `json:"password" validate:"required,min=6" label:"密码"`
    Age      int    `json:"age" validate:"required,min=1,max=120" label:"年龄"`
    Phone    string `json:"phone" validate:"omitempty,len=11" label:"手机号"`
}
```

### 3. 在控制器中使用

```go
func RegisterHandler(c *gin.Context) {
    var req UserRequest
    
    // 绑定并验证 JSON 参数
    valid, errs := v.BindJSONAndValid(c, &req)
    if !valid {
        response.Error(c, http.StatusBadRequest, "参数验证失败", errs)
        return
    }
    
    // 验证通过，处理业务逻辑
    // ...
}
```

## API 文档

### 核心方法

#### `BindJSONAndValid(c *gin.Context, obj interface{}) (bool, ValidErrors)`
绑定 JSON 参数并验证

#### `BindAndValid(c *gin.Context, obj interface{}) (bool, ValidErrors)`
绑定表单/查询参数并验证

#### `Validate(obj interface{}) (bool, ValidErrors)`
直接验证结构体

#### `SetLanguage(lang string) error`
设置验证错误信息的语言
- `"zh"`: 中文
- `"en"`: 英文

### 自定义验证

#### 注册自定义验证规则

```go
// 注册验证规则
v.RegisterValidation("not_admin", func(fl validator.FieldLevel) bool {
    return fl.Field().String() != "admin"
})

// 注册翻译
v.RegisterTranslation("not_admin", "{0}不能为admin")
```

#### 使用自定义验证

```go
type CustomRequest struct {
    Username string `json:"username" validate:"required,not_admin" label:"用户名"`
}
```

## 错误处理

### 错误结构

```go
type ValidError struct {
    Key     string `json:"key"`     // 字段名
    Message string `json:"message"` // 错误信息
}

type ValidErrors []*ValidError
```

### 错误示例

```json
[
    {
        "key": "UserRequest.Username",
        "message": "用户名长度必须至少为3个字符"
    },
    {
        "key": "UserRequest.Email",
        "message": "邮箱必须是一个有效的邮箱地址"
    }
]
```

## 常用验证标签

| 标签 | 说明 | 示例 |
|------|------|------|
| `required` | 必填字段 | `validate:"required"` |
| `min` | 最小长度/值 | `validate:"min=3"` |
| `max` | 最大长度/值 | `validate:"max=50"` |
| `len` | 固定长度 | `validate:"len=11"` |
| `email` | 邮箱格式 | `validate:"email"` |
| `oneof` | 枚举值 | `validate:"oneof=admin user guest"` |
| `gt` | 大于 | `validate:"gt=0"` |
| `gte` | 大于等于 | `validate:"gte=18"` |
| `lt` | 小于 | `validate:"lt=100"` |
| `lte` | 小于等于 | `validate:"lte=120"` |
| `omitempty` | 可选字段 | `validate:"omitempty,email"` |

## 与 BaseController 集成

项目中的 `BaseController` 已经集成了验证器，可以直接使用：

```go
type UserController struct {
    *BaseController
}

func (uc *UserController) Register(c *gin.Context) {
    var req UserRequest
    
    // 使用 BaseController 的验证方法
    if err := uc.BindAndValidate(c, &req); err != nil {
        return // 错误已经在方法内部处理
    }
    
    // 处理业务逻辑
    // ...
}
```

## 完整示例

查看 `examples/validator_example.go` 文件获取完整的使用示例，包括：

- JSON 参数验证
- 表单参数验证
- 查询参数验证
- 自定义验证规则
- 多语言支持
- 错误处理

## 运行示例

```bash
# 运行验证器示例
go run examples/validator_example.go

# 测试接口
curl -X POST http://localhost:8081/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test","email":"invalid-email","password":"123"}'
```

## 注意事项

1. 结构体字段需要添加 `validate` 标签来定义验证规则
2. 使用 `label` 标签可以自定义字段在错误信息中的显示名称
3. 验证器实例是线程安全的，可以在多个 goroutine 中使用
4. 自定义验证规则需要在使用前注册
5. 语言设置会影响所有后续的验证错误信息