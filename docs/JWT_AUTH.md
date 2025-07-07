# JWT 认证系统文档

## 概述

本项目采用 JWT (JSON Web Token) 作为用户认证和授权的核心机制。JWT 认证系统提供了安全、无状态的用户身份验证解决方案，支持令牌生成、解析、验证和刷新等完整功能。

## 功能特性

### 1. 令牌管理
- **令牌生成**: 用户登录成功后生成 JWT 令牌
- **令牌解析**: 解析请求中的 JWT 令牌获取用户信息
- **令牌验证**: 验证令牌的有效性和完整性
- **令牌刷新**: 在令牌即将过期时提供刷新机制

### 2. 用户信息
- **用户ID**: 唯一标识用户身份
- **用户名**: 用户登录名
- **角色权限**: 用户角色信息，支持基于角色的访问控制
- **过期时间**: 令牌有效期管理

### 3. 安全特性
- **HMAC-SHA256 签名**: 使用安全的签名算法
- **过期时间控制**: 可配置的令牌有效期
- **签发者验证**: 验证令牌签发者信息
- **时间窗口验证**: 检查令牌的生效时间和过期时间

## 配置说明

### JWT 配置结构

```yaml
jwt:
  secret: "your-jwt-secret-key"    # JWT 签名密钥（生产环境必须修改）
  expire_time: 24h                  # 令牌过期时间
  issuer: "kunpeng"                 # 令牌签发者
```

### 配置参数说明

- **secret**: JWT 签名密钥，用于令牌的签名和验证
  - 生产环境必须使用强密码
  - 建议长度至少 32 字符
  - 支持通过环境变量配置

- **expire_time**: 令牌有效期
  - 支持时间单位：h(小时)、m(分钟)、s(秒)
  - 建议值：1h-24h
  - 过短影响用户体验，过长存在安全风险

- **issuer**: 令牌签发者标识
  - 用于验证令牌来源
  - 通常设置为应用名称

## API 使用说明

### 初始化 JWT 管理器

```go
import (
    "github.com/cuiyuanxin/kunpeng/internal/auth"
    "github.com/cuiyuanxin/kunpeng/internal/config"
)

// 创建 JWT 管理器
jwtManager := auth.NewJWTManager(&cfg.JWT)
```

### 生成令牌

```go
// 用户登录成功后生成令牌
token, err := jwtManager.GenerateToken(userID, username, role)
if err != nil {
    // 处理错误
    return err
}

// 返回令牌给客户端
response.Success(c, gin.H{
    "token": token,
    "user": userInfo,
})
```

### 解析令牌

```go
// 从请求头获取令牌
token := c.GetHeader("Authorization")
// 移除 "Bearer " 前缀
token = strings.TrimPrefix(token, "Bearer ")

// 解析令牌
claims, err := jwtManager.ParseToken(token)
if err != nil {
    // 令牌无效
    response.Unauthorized(c, "Invalid token")
    return
}

// 获取用户信息
userID := claims.UserID
username := claims.Username
role := claims.Role
```

### 验证令牌

```go
// 简单验证令牌是否有效
isValid := jwtManager.ValidateToken(tokenString)
if !isValid {
    // 令牌无效
    response.Unauthorized(c, "Token validation failed")
    return
}
```

### 刷新令牌

```go
// 刷新即将过期的令牌
newToken, err := jwtManager.RefreshToken(oldToken)
if err != nil {
    // 令牌不符合刷新条件
    response.BadRequest(c, "Token refresh failed")
    return
}

// 返回新令牌
response.Success(c, gin.H{
    "token": newToken,
})
```

## 中间件集成

### JWT 认证中间件

```go
// 在路由中使用 JWT 认证中间件
protected := router.Group("/api/v1/user")
protected.Use(middleware.JWTAuth(jwtManager))
{
    protected.GET("/profile", userHandler.GetProfile)
    protected.PUT("/profile", userHandler.UpdateProfile)
}
```

### 角色权限中间件

```go
// 需要管理员权限的路由
admin := router.Group("/api/v1/admin")
admin.Use(middleware.JWTAuth(jwtManager))
admin.Use(middleware.RequireRole("admin"))
{
    admin.GET("/users", userHandler.GetUsers)
    admin.DELETE("/users/:id", userHandler.DeleteUser)
}
```

## 令牌结构

### Claims 结构

```go
type Claims struct {
    UserID   uint   `json:"user_id"`    // 用户ID
    Username string `json:"username"`   // 用户名
    Role     string `json:"role"`       // 用户角色
    jwt.RegisteredClaims                // 标准声明
}
```

### 标准声明字段

- **iss** (Issuer): 令牌签发者
- **iat** (Issued At): 令牌签发时间
- **exp** (Expiration Time): 令牌过期时间
- **nbf** (Not Before): 令牌生效时间

## 安全最佳实践

### 1. 密钥管理

```bash
# 使用环境变量存储密钥
export JWT_SECRET="your-very-strong-secret-key-here"
```

```yaml
# 配置文件中引用环境变量
jwt:
  secret: "${JWT_SECRET}"
```

### 2. 令牌传输

```javascript
// 客户端请求示例
fetch('/api/v1/user/profile', {
    headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
    }
});
```

### 3. 令牌存储

- **推荐**: 存储在内存中或 httpOnly Cookie
- **避免**: 存储在 localStorage 或 sessionStorage
- **HTTPS**: 生产环境必须使用 HTTPS 传输

### 4. 过期处理

```go
// 客户端应处理令牌过期
if response.Code == 401 {
    // 令牌过期，重新登录或刷新令牌
    redirectToLogin()
}
```

## 错误处理

### 常见错误类型

1. **令牌格式错误**
   ```
   Error: "Invalid authorization header format"
   Solution: 确保使用 "Bearer <token>" 格式
   ```

2. **令牌签名无效**
   ```
   Error: "Invalid token"
   Solution: 检查密钥配置是否正确
   ```

3. **令牌过期**
   ```
   Error: "Token expired"
   Solution: 使用刷新令牌或重新登录
   ```

4. **权限不足**
   ```
   Error: "Insufficient permissions"
   Solution: 检查用户角色权限
   ```

### 错误处理示例

```go
func handleJWTError(c *gin.Context, err error) {
    switch {
    case strings.Contains(err.Error(), "expired"):
        response.Unauthorized(c, "Token expired")
    case strings.Contains(err.Error(), "invalid"):
        response.Unauthorized(c, "Invalid token")
    default:
        response.Unauthorized(c, "Authentication failed")
    }
}
```

## 性能优化

### 1. 令牌缓存

```go
// 使用 Redis 缓存用户信息，减少数据库查询
func cacheUserInfo(userID uint, userInfo *UserInfo) {
    key := fmt.Sprintf("user:%d", userID)
    redis.Set(context.Background(), key, userInfo, time.Hour)
}
```

### 2. 批量验证

```go
// 对于高频 API，可以考虑批量验证令牌
func batchValidateTokens(tokens []string) map[string]bool {
    results := make(map[string]bool)
    for _, token := range tokens {
        results[token] = jwtManager.ValidateToken(token)
    }
    return results
}
```

## 监控和日志

### 认证事件日志

```go
// 记录认证相关事件
klogger.Info("User login",
    zap.Uint("user_id", userID),
    zap.String("username", username),
    zap.String("ip", clientIP),
)

klogger.Warn("Invalid token attempt",
    zap.String("token", maskedToken),
    zap.String("ip", clientIP),
    zap.Error(err),
)
```

### 安全监控

- 监控异常的令牌使用模式
- 记录失败的认证尝试
- 跟踪令牌刷新频率
- 监控权限提升尝试

## 测试示例

### 单元测试

```go
func TestJWTManager_GenerateToken(t *testing.T) {
    cfg := &config.JWT{
        Secret:     "test-secret",
        ExpireTime: time.Hour,
        Issuer:     "test",
    }
    
    jwtManager := auth.NewJWTManager(cfg)
    
    token, err := jwtManager.GenerateToken(1, "testuser", "user")
    assert.NoError(t, err)
    assert.NotEmpty(t, token)
    
    // 验证生成的令牌
    claims, err := jwtManager.ParseToken(token)
    assert.NoError(t, err)
    assert.Equal(t, uint(1), claims.UserID)
    assert.Equal(t, "testuser", claims.Username)
    assert.Equal(t, "user", claims.Role)
}
```

### 集成测试

```go
func TestJWTMiddleware(t *testing.T) {
    // 创建测试路由
    router := gin.New()
    router.Use(middleware.JWTAuth(jwtManager))
    router.GET("/protected", func(c *gin.Context) {
        response.Success(c, "access granted")
    })
    
    // 测试有效令牌
    token, _ := jwtManager.GenerateToken(1, "testuser", "user")
    req := httptest.NewRequest("GET", "/protected", nil)
    req.Header.Set("Authorization", "Bearer "+token)
    
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusOK, w.Code)
}
```

## 相关文档

- [中间件系统文档](MIDDLEWARE.md)
- [用户认证 API 文档](../api/auth.md)
- [安全配置指南](SECURITY.md)
- [JWT 官方规范](https://tools.ietf.org/html/rfc7519)

---

**注意**: 在生产环境中，请确保使用强密钥、启用 HTTPS，并定期轮换密钥以保证系统安全。