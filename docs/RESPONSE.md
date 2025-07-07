# 统一响应格式文档

## 概述

本项目实现了统一的 HTTP 响应格式，确保 API 接口返回的数据结构一致性和规范性。响应格式支持成功响应、错误响应、分页响应等多种场景，并提供了丰富的状态码和消息定义。

## 响应结构

### 基础响应结构

```go
type Response struct {
    Code    int         `json:"code"`    // 业务状态码
    Message string      `json:"message"` // 响应消息
    Data    interface{} `json:"data"`    // 响应数据
}
```

### 分页响应结构

```go
type PageResponse struct {
    Code    int         `json:"code"`     // 业务状态码
    Message string      `json:"message"`  // 响应消息
    Data    interface{} `json:"data"`     // 响应数据
    Total   int64       `json:"total"`    // 总记录数
    Page    int         `json:"page"`     // 当前页码
    Size    int         `json:"size"`     // 每页大小
}
```

## 状态码定义

### 成功状态码

```go
const (
    CodeSuccess = 200  // 操作成功
)
```

### 客户端错误状态码 (4xx)

```go
const (
    CodeBadRequest          = 400  // 请求参数错误
    CodeUnauthorized        = 401  // 未授权
    CodeForbidden          = 403  // 禁止访问
    CodeNotFound           = 404  // 资源不存在
    CodeMethodNotAllowed   = 405  // 方法不允许
    CodeConflict           = 409  // 资源冲突
    CodeValidationFailed   = 422  // 数据验证失败
    CodeTooManyRequests    = 429  // 请求过于频繁
)
```

### 服务器错误状态码 (5xx)

```go
const (
    CodeServerError        = 500  // 服务器内部错误
    CodeNotImplemented     = 501  // 功能未实现
    CodeBadGateway         = 502  // 网关错误
    CodeServiceUnavailable = 503  // 服务不可用
    CodeGatewayTimeout     = 504  // 网关超时
)
```

## 响应消息定义

### 成功消息

```go
const (
    MsgSuccess = "操作成功"
)
```

### 错误消息

```go
const (
    MsgBadRequest          = "请求参数错误"
    MsgUnauthorized        = "未授权访问"
    MsgForbidden          = "禁止访问"
    MsgNotFound           = "资源不存在"
    MsgMethodNotAllowed   = "请求方法不允许"
    MsgConflict           = "资源冲突"
    MsgValidationFailed   = "数据验证失败"
    MsgTooManyRequests    = "请求过于频繁"
    MsgServerError        = "服务器内部错误"
    MsgNotImplemented     = "功能未实现"
    MsgBadGateway         = "网关错误"
    MsgServiceUnavailable = "服务不可用"
    MsgGatewayTimeout     = "网关超时"
)
```

## 响应函数

### 成功响应

#### Success - 标准成功响应

```go
func Success(c *gin.Context, data interface{}) {
    c.JSON(http.StatusOK, Response{
        Code:    CodeSuccess,
        Message: MsgSuccess,
        Data:    data,
    })
}
```

**使用示例**:
```go
// 返回用户信息
user := &User{ID: 1, Name: "张三", Email: "zhangsan@example.com"}
response.Success(c, user)
```

**响应示例**:
```json
{
  "code": 200,
  "message": "操作成功",
  "data": {
    "id": 1,
    "name": "张三",
    "email": "zhangsan@example.com"
  }
}
```

#### SuccessWithMessage - 自定义消息成功响应

```go
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
    c.JSON(http.StatusOK, Response{
        Code:    CodeSuccess,
        Message: message,
        Data:    data,
    })
}
```

**使用示例**:
```go
// 用户注册成功
response.SuccessWithMessage(c, "用户注册成功", gin.H{"user_id": 123})
```

#### SuccessPage - 分页成功响应

```go
func SuccessPage(c *gin.Context, data interface{}, total int64, page, size int) {
    c.JSON(http.StatusOK, PageResponse{
        Code:    CodeSuccess,
        Message: MsgSuccess,
        Data:    data,
        Total:   total,
        Page:    page,
        Size:    size,
    })
}
```

**使用示例**:
```go
// 返回用户列表
users := []User{{ID: 1, Name: "张三"}, {ID: 2, Name: "李四"}}
response.SuccessPage(c, users, 100, 1, 10)
```

**响应示例**:
```json
{
  "code": 200,
  "message": "操作成功",
  "data": [
    {"id": 1, "name": "张三"},
    {"id": 2, "name": "李四"}
  ],
  "total": 100,
  "page": 1,
  "size": 10
}
```

### 错误响应

#### BadRequest - 请求参数错误

```go
func BadRequest(c *gin.Context, message string) {
    c.JSON(http.StatusBadRequest, Response{
        Code:    CodeBadRequest,
        Message: message,
        Data:    nil,
    })
}
```

**使用示例**:
```go
// 参数验证失败
response.BadRequest(c, "用户名不能为空")
```

#### Unauthorized - 未授权访问

```go
func Unauthorized(c *gin.Context, message string) {
    c.JSON(http.StatusUnauthorized, Response{
        Code:    CodeUnauthorized,
        Message: message,
        Data:    nil,
    })
}
```

**使用示例**:
```go
// JWT 令牌无效
response.Unauthorized(c, "令牌已过期，请重新登录")
```

#### Forbidden - 禁止访问

```go
func Forbidden(c *gin.Context, message string) {
    c.JSON(http.StatusForbidden, Response{
        Code:    CodeForbidden,
        Message: message,
        Data:    nil,
    })
}
```

**使用示例**:
```go
// 权限不足
response.Forbidden(c, "您没有权限访问此资源")
```

#### NotFound - 资源不存在

```go
func NotFound(c *gin.Context, message string) {
    c.JSON(http.StatusNotFound, Response{
        Code:    CodeNotFound,
        Message: message,
        Data:    nil,
    })
}
```

**使用示例**:
```go
// 用户不存在
response.NotFound(c, "用户不存在")
```

#### ServerError - 服务器内部错误

```go
func ServerError(c *gin.Context, message string) {
    c.JSON(http.StatusInternalServerError, Response{
        Code:    CodeServerError,
        Message: message,
        Data:    nil,
    })
}
```

**使用示例**:
```go
// 数据库连接失败
response.ServerError(c, "数据库连接失败")
```

#### Custom - 自定义响应

```go
func Custom(c *gin.Context, httpCode, businessCode int, message string, data interface{}) {
    c.JSON(httpCode, Response{
        Code:    businessCode,
        Message: message,
        Data:    data,
    })
}
```

**使用示例**:
```go
// 自定义业务错误
response.Custom(c, http.StatusOK, 1001, "用户已存在", nil)
```

## 使用场景

### 1. 用户认证场景

```go
// 登录成功
func Login(c *gin.Context) {
    // ... 验证逻辑
    
    token, err := jwtManager.GenerateToken(user.ID, user.Username, user.Role)
    if err != nil {
        response.ServerError(c, "生成令牌失败")
        return
    }
    
    response.Success(c, gin.H{
        "token": token,
        "user": gin.H{
            "id":       user.ID,
            "username": user.Username,
            "role":     user.Role,
        },
    })
}

// 登录失败
func LoginFailed(c *gin.Context) {
    response.Unauthorized(c, "用户名或密码错误")
}
```

### 2. 数据验证场景

```go
func CreateUser(c *gin.Context) {
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.BadRequest(c, "请求参数格式错误")
        return
    }
    
    // 验证必填字段
    if req.Username == "" {
        response.BadRequest(c, "用户名不能为空")
        return
    }
    
    if req.Email == "" {
        response.BadRequest(c, "邮箱不能为空")
        return
    }
    
    // 检查用户是否已存在
    if userService.ExistsByUsername(req.Username) {
        response.Custom(c, http.StatusOK, 1001, "用户名已存在", nil)
        return
    }
    
    // 创建用户
    user, err := userService.Create(&req)
    if err != nil {
        response.ServerError(c, "创建用户失败")
        return
    }
    
    response.SuccessWithMessage(c, "用户创建成功", user)
}
```

### 3. 分页查询场景

```go
func GetUsers(c *gin.Context) {
    // 获取分页参数
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
    
    // 参数验证
    if page < 1 {
        response.BadRequest(c, "页码必须大于0")
        return
    }
    
    if size < 1 || size > 100 {
        response.BadRequest(c, "每页大小必须在1-100之间")
        return
    }
    
    // 查询数据
    users, total, err := userService.GetList(page, size)
    if err != nil {
        response.ServerError(c, "查询用户列表失败")
        return
    }
    
    response.SuccessPage(c, users, total, page, size)
}
```

### 4. 权限控制场景

```go
func DeleteUser(c *gin.Context) {
    userID, _ := strconv.ParseUint(c.Param("id"), 10, 32)
    currentUserID := c.GetUint("user_id")
    currentUserRole := c.GetString("role")
    
    // 检查是否为管理员或本人
    if currentUserRole != "admin" && currentUserID != uint(userID) {
        response.Forbidden(c, "您没有权限删除此用户")
        return
    }
    
    // 检查用户是否存在
    user, err := userService.GetByID(uint(userID))
    if err != nil {
        response.NotFound(c, "用户不存在")
        return
    }
    
    // 删除用户
    if err := userService.Delete(uint(userID)); err != nil {
        response.ServerError(c, "删除用户失败")
        return
    }
    
    response.SuccessWithMessage(c, "用户删除成功", nil)
}
```

## 错误处理最佳实践

### 1. 统一错误处理

```go
// 定义业务错误类型
type BusinessError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
}

func (e *BusinessError) Error() string {
    return e.Message
}

// 业务错误常量
var (
    ErrUserNotFound     = &BusinessError{Code: 1001, Message: "用户不存在"}
    ErrUserExists       = &BusinessError{Code: 1002, Message: "用户已存在"}
    ErrInvalidPassword  = &BusinessError{Code: 1003, Message: "密码错误"}
    ErrTokenExpired     = &BusinessError{Code: 1004, Message: "令牌已过期"}
)

// 统一错误处理函数
func HandleError(c *gin.Context, err error) {
    switch e := err.(type) {
    case *BusinessError:
        response.Custom(c, http.StatusOK, e.Code, e.Message, nil)
    case *gorm.ErrRecordNotFound:
        response.NotFound(c, "记录不存在")
    default:
        klogger.Error("Unexpected error", zap.Error(err))
        response.ServerError(c, "系统错误")
    }
}
```

### 2. 参数验证错误处理

```go
// 验证错误处理
func HandleValidationError(c *gin.Context, err error) {
    if validationErrors, ok := err.(validator.ValidationErrors); ok {
        var messages []string
        for _, fieldError := range validationErrors {
            switch fieldError.Tag() {
            case "required":
                messages = append(messages, fmt.Sprintf("%s 是必填字段", fieldError.Field()))
            case "email":
                messages = append(messages, fmt.Sprintf("%s 必须是有效的邮箱地址", fieldError.Field()))
            case "min":
                messages = append(messages, fmt.Sprintf("%s 长度不能少于 %s 个字符", fieldError.Field(), fieldError.Param()))
            case "max":
                messages = append(messages, fmt.Sprintf("%s 长度不能超过 %s 个字符", fieldError.Field(), fieldError.Param()))
            default:
                messages = append(messages, fmt.Sprintf("%s 验证失败", fieldError.Field()))
            }
        }
        response.BadRequest(c, strings.Join(messages, "; "))
    } else {
        response.BadRequest(c, "请求参数格式错误")
    }
}
```

## 响应格式扩展

### 1. 带时间戳的响应

```go
type TimestampResponse struct {
    Code      int         `json:"code"`
    Message   string      `json:"message"`
    Data      interface{} `json:"data"`
    Timestamp int64       `json:"timestamp"`
}

func SuccessWithTimestamp(c *gin.Context, data interface{}) {
    c.JSON(http.StatusOK, TimestampResponse{
        Code:      CodeSuccess,
        Message:   MsgSuccess,
        Data:      data,
        Timestamp: time.Now().Unix(),
    })
}
```

### 2. 带请求ID的响应

```go
type TrackedResponse struct {
    Code      int         `json:"code"`
    Message   string      `json:"message"`
    Data      interface{} `json:"data"`
    RequestID string      `json:"request_id"`
}

func SuccessWithRequestID(c *gin.Context, data interface{}) {
    requestID := c.GetString("request_id")
    c.JSON(http.StatusOK, TrackedResponse{
        Code:      CodeSuccess,
        Message:   MsgSuccess,
        Data:      data,
        RequestID: requestID,
    })
}
```

### 3. 多语言响应

```go
// 消息国际化
var messages = map[string]map[string]string{
    "zh": {
        "success":      "操作成功",
        "bad_request":  "请求参数错误",
        "unauthorized": "未授权访问",
        "forbidden":    "禁止访问",
        "not_found":    "资源不存在",
        "server_error": "服务器内部错误",
    },
    "en": {
        "success":      "Success",
        "bad_request":  "Bad Request",
        "unauthorized": "Unauthorized",
        "forbidden":    "Forbidden",
        "not_found":    "Not Found",
        "server_error": "Internal Server Error",
    },
}

func getMessage(lang, key string) string {
    if langMessages, ok := messages[lang]; ok {
        if message, ok := langMessages[key]; ok {
            return message
        }
    }
    return messages["zh"][key] // 默认中文
}

func SuccessI18n(c *gin.Context, data interface{}) {
    lang := c.GetHeader("Accept-Language")
    if lang == "" {
        lang = "zh"
    }
    
    c.JSON(http.StatusOK, Response{
        Code:    CodeSuccess,
        Message: getMessage(lang, "success"),
        Data:    data,
    })
}
```

## 测试

### 响应格式测试

```go
func TestSuccessResponse(t *testing.T) {
    router := gin.New()
    router.GET("/test", func(c *gin.Context) {
        response.Success(c, gin.H{"message": "test"})
    })
    
    req := httptest.NewRequest("GET", "/test", nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusOK, w.Code)
    
    var resp response.Response
    err := json.Unmarshal(w.Body.Bytes(), &resp)
    assert.NoError(t, err)
    assert.Equal(t, response.CodeSuccess, resp.Code)
    assert.Equal(t, response.MsgSuccess, resp.Message)
}

func TestErrorResponse(t *testing.T) {
    router := gin.New()
    router.GET("/test", func(c *gin.Context) {
        response.BadRequest(c, "参数错误")
    })
    
    req := httptest.NewRequest("GET", "/test", nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusBadRequest, w.Code)
    
    var resp response.Response
    err := json.Unmarshal(w.Body.Bytes(), &resp)
    assert.NoError(t, err)
    assert.Equal(t, response.CodeBadRequest, resp.Code)
    assert.Equal(t, "参数错误", resp.Message)
}
```

## 相关文档

- [中间件系统文档](MIDDLEWARE.md)
- [JWT 认证系统文档](JWT_AUTH.md)
- [路由管理文档](ROUTER.md)
- [API 文档生成指南](../README.md#api-文档)

---

**最佳实践**: 始终使用统一的响应格式；为不同的错误场景提供明确的状态码和消息；在开发环境中提供详细的错误信息，在生产环境中避免暴露敏感信息；使用业务状态码区分不同的业务错误。