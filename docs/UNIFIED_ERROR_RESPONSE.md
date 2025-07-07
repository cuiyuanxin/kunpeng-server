# 统一错误处理与响应系统

## 概述

本文档介绍了如何使用融合后的 `response` 包和 `errors` 包，实现统一的错误处理和响应机制。这个系统提供了预设的错误代码、错误描述和详细的错误信息处理能力。

## 核心特性

### 1. 统一错误码体系
- **分类管理**: 按业务模块对错误进行分类（1000-8999）
- **预设错误**: 提供常用的业务错误定义
- **自定义错误**: 支持创建新的业务错误
- **详细信息**: 支持添加具体的错误详情

### 2. 自动错误转换
- **GORM错误转换**: 自动将数据库错误转换为业务错误
- **系统错误处理**: 统一处理panic和系统级错误
- **HTTP状态码映射**: 自动映射业务错误到合适的HTTP状态码

### 3. 便捷响应函数
- **成功响应**: 统一的成功数据返回格式
- **错误响应**: 自动识别和处理业务错误
- **分页响应**: 内置分页数据响应支持
- **条件响应**: 根据错误情况自动选择响应类型

## 错误码分类

```go
// 通用错误 (1000-1999)
ErrCodeInternalError     = 1000  // 服务器内部错误
ErrCodeInvalidParams     = 1001  // 请求参数无效
ErrCodeValidationFailed  = 1002  // 数据验证失败
ErrCodeResourceNotFound  = 1003  // 资源不存在
ErrCodeResourceExists    = 1004  // 资源已存在
ErrCodePermissionDenied  = 1005  // 权限不足

// 认证相关错误 (2000-2999)
ErrCodeUnauthorized      = 2000  // 未授权访问
ErrCodeTokenExpired      = 2001  // 令牌已过期
ErrCodeTokenInvalid      = 2002  // 令牌无效
ErrCodeLoginFailed       = 2003  // 登录失败

// 用户相关错误 (3000-3999)
ErrCodeUserNotFound      = 3000  // 用户不存在
ErrCodeUserExists        = 3001  // 用户已存在
ErrCodeUsernameTaken     = 3002  // 用户名已被占用
ErrCodeEmailTaken        = 3003  // 邮箱已被占用

// 数据库相关错误 (6000-6999)
ErrCodeDatabaseError     = 6000  // 数据库操作失败
ErrCodeRecordNotFound    = 6001  // 记录不存在
ErrCodeDuplicateKey      = 6002  // 数据重复
```

## 基本使用方法

### 1. 控制器层使用

#### 简单错误处理
```go
func (ctrl *UserController) GetUser(c *gin.Context) {
    userID := c.Param("id")
    
    // 调用服务层，自动处理错误
    user, err := ctrl.userService.GetUserByID(userID)
    response.SuccessOrError(c, user, err)
}
```

#### 参数验证错误
```go
func (ctrl *UserController) CreateUser(c *gin.Context) {
    var req CreateUserRequest
    
    if err := c.ShouldBindJSON(&req); err != nil {
        response.ValidationError(c, err.Error())
        return
    }
    
    user, err := ctrl.userService.CreateUser(req.Username, req.Email)
    response.SuccessOrError(c, user, err)
}
```

#### 分页响应
```go
func (ctrl *UserController) GetUsers(c *gin.Context) {
    page, size := getPageParams(c) // 获取分页参数
    
    users, total, err := ctrl.userService.GetUsers(page, size)
    response.SuccessPageOrError(c, users, total, page, size, err)
}
```

### 2. 服务层使用

#### 预定义错误使用
```go
func (s *UserService) GetUserByID(userID string) (*User, error) {
    if userID == "" {
        return nil, errors.ErrInvalidParams.WithDetails("用户ID不能为空")
    }
    
    var user User
    err := s.db.First(&user, "id = ?", userID).Error
    if err != nil {
        return nil, errors.ConvertGormError(err) // 自动转换GORM错误
    }
    
    return &user, nil
}
```

#### 自定义错误详情
```go
func (s *UserService) CreateUser(username, email string) (*User, error) {
    // 检查用户名是否已存在
    var count int64
    err := s.db.Model(&User{}).Where("username = ?", username).Count(&count).Error
    if err != nil {
        return nil, errors.ConvertGormError(err)
    }
    
    if count > 0 {
        return nil, errors.ErrUsernameTaken.WithDetailsf("用户名 '%s' 已被占用", username)
    }
    
    // ... 创建用户逻辑
}
```

### 3. 中间件使用

#### 认证中间件
```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            errors.AbortWithError(c, errors.ErrUnauthorized.WithDetails("缺少认证令牌"))
            return
        }
        
        if !isValidToken(token) {
            errors.AbortWithError(c, errors.ErrTokenInvalid)
            return
        }
        
        c.Next()
    }
}
```

#### 全局错误处理中间件
```go
func setupRouter() *gin.Engine {
    r := gin.New()
    
    // 注册全局错误处理中间件
    r.Use(errors.ErrorHandlerMiddleware())
    
    // ... 其他中间件和路由
    
    return r
}
```

## 高级功能

### 1. 安全执行

防止panic导致程序崩溃：

```go
func (s *UserService) RiskyOperation() error {
    return errors.SafeExecute(func() error {
        // 可能会panic的代码
        result := someRiskyFunction()
        return processResult(result)
    })
}

func (s *UserService) RiskyOperationWithResult() (interface{}, error) {
    return errors.SafeExecuteWithResult(func() (interface{}, error) {
        // 可能会panic的代码
        result := someRiskyFunction()
        return result, nil
    })
}
```

### 2. 链式错误处理

按顺序执行多个操作，遇到错误立即返回：

```go
func (s *UserService) ComplexOperation(userID string) error {
    return errors.Chain(
        func() error {
            return s.validateUser(userID)
        },
        func() error {
            return s.checkPermissions(userID)
        },
        func() error {
            return s.performOperation(userID)
        },
        func() error {
            return s.updateCache(userID)
        },
    )
}
```

### 3. 重试机制

对于可能临时失败的操作，提供重试支持：

```go
func (s *ExternalService) CallAPI() error {
    return errors.RetryOnError(func() error {
        // 调用外部API
        return callExternalAPI()
    }, 3) // 最多重试3次
}
```

### 4. GORM错误自动转换

系统会自动将常见的GORM错误转换为对应的业务错误：

```go
// 自动转换示例
err := db.First(&user, "id = ?", userID).Error
if err != nil {
    return errors.ConvertGormError(err)
    // "record not found" -> ErrRecordNotFound
    // "duplicate key" -> ErrDuplicateKey
    // "constraint violation" -> ErrConstraintViolation
}
```

## 响应格式

### 成功响应
```json
{
    "code": 200,
    "message": "success",
    "data": {
        "id": 1,
        "username": "john",
        "email": "john@example.com"
    }
}
```

### 错误响应
```json
{
    "code": 3000,
    "message": "用户不存在",
    "data": {
        "details": "用户ID: 123 不存在于系统中"
    }
}
```

### 分页响应
```json
{
    "code": 200,
    "message": "success",
    "data": [
        {"id": 1, "username": "user1"},
        {"id": 2, "username": "user2"}
    ],
    "total": 100,
    "page": 1,
    "size": 10
}
```

## 最佳实践

### 1. 错误处理原则
- **早期返回**: 一旦发现错误，立即返回，避免继续执行
- **错误包装**: 在错误传递过程中添加上下文信息
- **统一处理**: 使用中间件统一处理错误响应
- **详细日志**: 记录足够的错误信息用于调试

### 2. 控制器层
- 使用 `response.SuccessOrError()` 简化错误处理
- 参数验证失败使用 `response.ValidationError()`
- 认证失败使用 `response.AuthError()`
- 权限不足使用 `response.PermissionError()`

### 3. 服务层
- 返回具体的业务错误，不要返回通用错误
- 使用 `WithDetails()` 或 `WithDetailsf()` 添加错误详情
- 对于数据库操作，使用 `errors.ConvertGormError()` 转换错误
- 对于可能panic的代码，使用 `errors.SafeExecute()`

### 4. 中间件层
- 使用 `errors.AbortWithError()` 中断请求并返回错误
- 注册全局错误处理中间件 `errors.ErrorHandlerMiddleware()`
- 对于设置错误但不立即返回的情况，使用 `errors.SetError()`

## 扩展自定义错误

### 添加新的错误码
```go
// 在 errors.go 中添加新的错误码
const (
    // 订单相关错误 (9000-9999)
    ErrCodeOrderNotFound    ErrorCode = 9000
    ErrCodeOrderCancelled   ErrorCode = 9001
    ErrCodeOrderExpired     ErrorCode = 9002
)

// 添加预定义错误实例
var (
    ErrOrderNotFound  = &BusinessError{ErrCodeOrderNotFound, "订单不存在", http.StatusNotFound, ""}
    ErrOrderCancelled = &BusinessError{ErrCodeOrderCancelled, "订单已取消", http.StatusConflict, ""}
    ErrOrderExpired   = &BusinessError{ErrCodeOrderExpired, "订单已过期", http.StatusGone, ""}
)
```

### 创建业务特定的响应函数
```go
// 在 response.go 中添加业务特定的响应函数
func OrderError(c *gin.Context, details string) {
    HandleBusinessError(c, errors.ErrOrderNotFound.WithDetails(details))
}

func OrderCancelledError(c *gin.Context, orderID string) {
    HandleBusinessError(c, errors.ErrOrderCancelled.WithDetailsf("订单 %s 已被取消", orderID))
}
```

## 相关文档

- [错误中心文档](ERROR_CENTER.md)
- [统一响应格式文档](RESPONSE.md)
- [中间件系统文档](MIDDLEWARE.md)
- [处理器与服务层文档](HANDLER_SERVICE.md)

---

**注意**: 这个统一的错误处理和响应系统是项目的核心组件，所有新的功能开发都应该使用这个系统来处理错误和响应。