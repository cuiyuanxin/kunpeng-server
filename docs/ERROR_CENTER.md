# 错误中心文档

## 概述

错误中心是一个独立集中的错误处理系统，用于统一管理应用程序中的所有错误定义、错误处理逻辑和错误响应格式。

## 核心特性

- **统一错误码**: 所有错误都有唯一的错误码标识
- **分类管理**: 按业务模块对错误进行分类
- **链式调用**: 支持错误信息的链式添加和修改
- **自动转换**: 自动将系统错误转换为业务错误
- **中间件支持**: 提供全局错误处理中间件
- **详细日志**: 自动记录错误日志和堆栈信息

## 错误码分类

### 通用错误 (1000-1999)
- `1000`: 服务器内部错误
- `1001`: 请求参数无效
- `1002`: 数据验证失败
- `1003`: 资源不存在
- `1004`: 资源已存在
- `1005`: 权限不足
- `1006`: 请求频率超限
- `1007`: 服务不可用

### 认证相关错误 (2000-2999)
- `2000`: 未授权访问
- `2001`: 令牌已过期
- `2002`: 令牌无效
- `2003`: 登录失败
- `2004`: 密码错误
- `2005`: 账户已锁定
- `2006`: 账户已禁用

### 用户相关错误 (3000-3999)
- `3000`: 用户不存在
- `3001`: 用户已存在
- `3002`: 用户名已被占用
- `3003`: 邮箱已被占用
- `3004`: 邮箱格式无效
- `3005`: 密码强度不足

### 管理员相关错误 (4000-4999)
- `4000`: 管理员不存在
- `4001`: 管理员已存在
- `4002`: 权限不足

### 角色权限相关错误 (5000-5999)
- `5000`: 角色不存在
- `5001`: 角色已存在
- `5002`: 权限不存在
- `5003`: 角色正在使用中

### 数据库相关错误 (6000-6999)
- `6000`: 数据库操作失败
- `6001`: 记录不存在
- `6002`: 数据重复
- `6003`: 数据约束违反

### 文件相关错误 (7000-7999)
- `7000`: 文件不存在
- `7001`: 文件上传失败
- `7002`: 文件类型不允许
- `7003`: 文件大小超限

### 第三方服务错误 (8000-8999)
- `8000`: 外部服务错误
- `8001`: API调用失败
- `8002`: 网络错误

## 基本使用

### 1. 导入错误包

```go
import (
    "github.com/cuiyuanxin/kunpeng/pkg/errors"
)
```

### 2. 在处理器中使用预定义错误

```go
func (h *UserHandler) GetUser(c *gin.Context) {
    userID := c.Param("id")
    
    user, err := h.userService.GetByID(userID)
    if err != nil {
        // 直接使用预定义错误
        errors.AbortWithError(c, errors.ErrUserNotFound)
        return
    }
    
    response.Success(c, user)
}
```

### 3. 添加详细信息

```go
func (h *UserHandler) CreateUser(c *gin.Context) {
    var req model.UserCreateRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        // 添加详细错误信息
        bizErr := errors.ErrValidationFailed.WithDetails(err.Error())
        errors.AbortWithError(c, bizErr)
        return
    }
    
    // 业务逻辑...
}
```

### 4. 格式化详细信息

```go
func (s *UserService) GetByID(id string) (*model.User, error) {
    var user model.User
    err := s.db.First(&user, "id = ?", id).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.ErrUserNotFound.WithDetailsf("用户ID: %s", id)
        }
        return nil, errors.ConvertGormError(err)
    }
    return &user, nil
}
```

## 中间件使用

### 1. 注册全局错误处理中间件

```go
// 在路由设置中添加
func (r *Router) Setup() {
    // 添加全局错误处理中间件
    r.engine.Use(errors.ErrorHandlerMiddleware())
    
    // 其他中间件...
    r.setupAPIRoutes()
}
```

### 2. 在处理器中设置错误

```go
func (h *UserHandler) UpdateUser(c *gin.Context) {
    // 方式1: 直接设置错误，由中间件处理
    if someCondition {
        errors.SetError(c, errors.ErrPermissionDenied)
        return
    }
    
    // 方式2: 设置验证错误
    if validationErr != nil {
        errors.SetValidationError(c, validationErr)
        return
    }
    
    // 方式3: 设置数据库错误
    if dbErr != nil {
        errors.SetDatabaseError(c, dbErr)
        return
    }
}
```

## 服务层错误处理

### 1. 数据库错误转换

```go
func (s *UserService) Create(req *model.UserCreateRequest) (*model.User, error) {
    user := &model.User{
        Username: req.Username,
        Email:    req.Email,
    }
    
    err := s.db.Create(user).Error
    if err != nil {
        // 自动转换GORM错误为业务错误
        return nil, errors.ConvertGormError(err)
    }
    
    return user, nil
}
```

### 2. 验证错误转换

```go
func (s *UserService) ValidateUser(req *model.UserCreateRequest) error {
    if err := s.validator.Struct(req); err != nil {
        // 转换验证错误为业务错误
        return errors.ConvertValidationError(err)
    }
    return nil
}
```

### 3. 错误包装

```go
func (s *UserService) ComplexOperation(userID string) error {
    user, err := s.GetByID(userID)
    if err != nil {
        return errors.WrapErrorf(err, "获取用户失败, ID: %s", userID)
    }
    
    err = s.updateUserStatus(user)
    if err != nil {
        return errors.WrapError(err, "更新用户状态失败")
    }
    
    return nil
}
```

## 高级功能

### 1. 安全执行

```go
func (s *UserService) SafeOperation() error {
    return errors.SafeExecute(func() error {
        // 可能会panic的代码
        riskyOperation()
        return nil
    })
}
```

### 2. 链式错误处理

```go
func (s *UserService) ChainOperations() error {
    return errors.Chain(
        func() error { return s.validateInput() },
        func() error { return s.checkPermissions() },
        func() error { return s.performOperation() },
        func() error { return s.updateCache() },
    )
}
```

### 3. 重试机制

```go
func (s *UserService) RetryableOperation() error {
    return errors.RetryOnError(func() error {
        return s.callExternalAPI()
    }, 3) // 最多重试3次
}
```

### 4. 并行错误收集

```go
func (s *UserService) ParallelOperations() error {
    errs := errors.Parallel(
        func() error { return s.operation1() },
        func() error { return s.operation2() },
        func() error { return s.operation3() },
    )
    
    if len(errs) > 0 {
        return errors.CombineErrors(errs...)
    }
    
    return nil
}
```

## 自定义错误

### 1. 创建新的业务错误

```go
// 定义新的错误码
const (
    ErrCodeCustomOperation errors.ErrorCode = 9001
)

// 创建自定义错误
var ErrCustomOperation = errors.NewBusinessError(
    ErrCodeCustomOperation,
    "自定义操作失败",
    http.StatusBadRequest,
)
```

### 2. 在运行时创建错误

```go
func (s *UserService) CustomValidation(data string) error {
    if len(data) > 100 {
        return errors.NewBusinessError(
            errors.ErrCodeValidationFailed,
            "数据长度超过限制",
            http.StatusBadRequest,
        ).WithDetailsf("当前长度: %d, 最大长度: 100", len(data))
    }
    return nil
}
```

## 错误响应格式

### 标准错误响应

```json
{
    "code": 3000,
    "message": "用户不存在"
}
```

### 带详细信息的错误响应

```json
{
    "code": 3000,
    "message": "用户不存在",
    "details": "用户ID: 12345"
}
```

### 验证错误响应

```json
{
    "code": 1002,
    "message": "数据验证失败",
    "details": "username 是必填字段; email 必须是有效的邮箱地址"
}
```

## 最佳实践

### 1. 错误处理原则

- **早期返回**: 一旦发现错误，立即返回，避免继续执行
- **错误包装**: 在错误传递过程中添加上下文信息
- **统一处理**: 使用中间件统一处理错误响应
- **详细日志**: 记录足够的错误信息用于调试

### 2. 服务层错误处理

```go
func (s *UserService) GetUser(id string) (*model.User, error) {
    // 1. 参数验证
    if id == "" {
        return nil, errors.ErrInvalidParams.WithDetails("用户ID不能为空")
    }
    
    // 2. 数据库操作
    var user model.User
    err := s.db.First(&user, "id = ?", id).Error
    if err != nil {
        return nil, errors.ConvertGormError(err)
    }
    
    // 3. 业务逻辑验证
    if user.Status == "deleted" {
        return nil, errors.ErrUserNotFound.WithDetails("用户已被删除")
    }
    
    return &user, nil
}
```

### 3. 处理器层错误处理

```go
func (h *UserHandler) GetUser(c *gin.Context) {
    userID := c.Param("id")
    
    user, err := h.userService.GetUser(userID)
    if err != nil {
        // 统一错误处理，无需判断错误类型
        errors.AbortWithError(c, err)
        return
    }
    
    response.Success(c, user)
}
```

### 4. 中间件错误处理

```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            errors.AbortWithError(c, errors.ErrUnauthorized)
            return
        }
        
        claims, err := validateToken(token)
        if err != nil {
            errors.AbortWithError(c, errors.ErrTokenInvalid.WithDetails(err.Error()))
            return
        }
        
        c.Set("user_id", claims.UserID)
        c.Next()
    }
}
```

## 迁移指南

### 从旧的错误处理迁移

1. **替换直接的HTTP响应**:
   ```go
   // 旧方式
   c.JSON(http.StatusBadRequest, gin.H{"error": "用户不存在"})
   
   // 新方式
   errors.AbortWithError(c, errors.ErrUserNotFound)
   ```

2. **替换response包的错误函数**:
   ```go
   // 旧方式
   response.BadRequest(c, "参数错误")
   
   // 新方式
   errors.AbortWithError(c, errors.ErrInvalidParams.WithDetails("参数错误"))
   ```

3. **统一错误码**:
   ```go
   // 旧方式
   const (
       USER_NOT_FOUND = 1001
       INVALID_PARAMS = 1002
   )
   
   // 新方式
   // 使用预定义的错误码
   errors.ErrUserNotFound
   errors.ErrInvalidParams
   ```

## 相关文档

- [统一响应格式文档](RESPONSE.md)
- [中间件系统文档](MIDDLEWARE.md)
- [日志系统文档](LOGGING.md)
- [API文档生成指南](../README.md#api-文档)

---

**注意**: 错误中心是应用程序错误处理的核心，所有新的错误处理都应该使用这个系统。旧的错误处理方式应该逐步迁移到新系统。