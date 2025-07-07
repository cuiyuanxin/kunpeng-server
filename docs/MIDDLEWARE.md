# 中间件系统文档

## 概述

本项目基于 Gin 框架构建了完整的中间件系统，提供了跨域处理、日志记录、错误恢复、JWT 认证、角色权限、请求追踪、限流控制和超时管理等功能。中间件采用洋葱模型设计，支持链式调用和灵活配置。

## 中间件架构

### 执行顺序

```
请求 → Recovery → Logger → CORS → RequestID → RateLimiter → JWTAuth → RequireRole → 业务处理器 → 响应
```

### 洋葱模型

```
┌─────────────────────────────────────────┐
│                Recovery                 │
│  ┌───────────────────────────────────┐  │
│  │             Logger                │  │
│  │  ┌─────────────────────────────┐  │  │
│  │  │           CORS              │  │  │
│  │  │  ┌───────────────────────┐  │  │  │
│  │  │  │      RequestID        │  │  │  │
│  │  │  │  ┌─────────────────┐  │  │  │  │
│  │  │  │  │   RateLimiter   │  │  │  │  │
│  │  │  │  │  ┌───────────┐  │  │  │  │  │
│  │  │  │  │  │  JWTAuth  │  │  │  │  │  │
│  │  │  │  │  │  Handler  │  │  │  │  │  │
│  │  │  │  │  └───────────┘  │  │  │  │  │
│  │  │  │  └─────────────────┘  │  │  │  │
│  │  │  └───────────────────────┘  │  │  │
│  │  └─────────────────────────────┘  │  │
│  └───────────────────────────────────┘  │
└─────────────────────────────────────────┘
```

## 核心中间件

### 1. Recovery 中间件

**功能**: 捕获 panic 异常，防止服务器崩溃

```go
// 使用方式
router.Use(middleware.Recovery())
```

**特性**:
- 自动捕获 panic 异常
- 记录详细的错误日志
- 返回统一的错误响应
- 防止服务器进程终止

**日志记录**:
```go
klogger.Error("Panic recovered",
    zap.Any("error", recovered),
    zap.String("path", c.Request.URL.Path),
    zap.String("method", c.Request.Method),
)
```

### 2. Logger 中间件

**功能**: 记录 HTTP 请求和响应信息

```go
// 使用方式
router.Use(middleware.Logger())
```

**记录信息**:
- 请求方法和路径
- 客户端 IP 地址
- 响应状态码
- 请求处理时间
- User-Agent 信息

**日志格式**:
```json
{
  "level": "info",
  "msg": "HTTP Request",
  "method": "GET",
  "path": "/api/v1/users",
  "client_ip": "192.168.1.100",
  "status_code": 200,
  "latency": "15.2ms",
  "user_agent": "Mozilla/5.0..."
}
```

### 3. CORS 中间件

**功能**: 处理跨域资源共享

```go
// 使用方式
router.Use(middleware.CORS())
```

**配置的响应头**:
```
Access-Control-Allow-Origin: *
Access-Control-Allow-Methods: GET,POST,PUT,PATCH,DELETE,OPTIONS
Access-Control-Allow-Headers: authorization, origin, content-type, accept
Access-Control-Allow-Credentials: true
Access-Control-Max-Age: 86400
```

**OPTIONS 请求处理**:
- 自动响应预检请求
- 返回 204 No Content 状态码
- 提供完整的 CORS 头信息

### 4. RequestID 中间件

**功能**: 为每个请求生成唯一标识符

```go
// 使用方式
router.Use(middleware.RequestID())
```

**特性**:
- 自动生成请求 ID
- 支持客户端传入的 X-Request-ID
- 添加到响应头中
- 存储在上下文中供后续使用

**使用示例**:
```go
// 在处理器中获取请求 ID
requestID := c.GetString("request_id")
klogger.Info("Processing request", zap.String("request_id", requestID))
```

### 5. RateLimiter 中间件

**功能**: 请求限流控制

```go
// 使用方式
router.Use(middleware.RateLimiter())
```

**当前实现**:
- 基础框架已就绪
- 可扩展为令牌桶、滑动窗口等算法
- 支持基于 IP、用户、API 的限流

**扩展示例**:
```go
// Redis 分布式限流实现
func RedisRateLimiter(limit int, window time.Duration) gin.HandlerFunc {
    return func(c *gin.Context) {
        key := "rate_limit:" + c.ClientIP()
        count, _ := redis.Incr(context.Background(), key)
        if count == 1 {
            redis.Expire(context.Background(), key, window)
        }
        if count > int64(limit) {
            response.Custom(c, http.StatusTooManyRequests, 429, "Rate limit exceeded", nil)
            c.Abort()
            return
        }
        c.Next()
    }
}
```

### 6. JWTAuth 中间件

**功能**: JWT 令牌认证

```go
// 使用方式
protected := router.Group("/api/v1/user")
protected.Use(middleware.JWTAuth(jwtManager))
```

**认证流程**:
1. 从 Authorization 头获取令牌
2. 验证 Bearer 格式
3. 解析和验证 JWT 令牌
4. 提取用户信息到上下文

**上下文信息**:
```go
// 在后续处理器中可获取
userID := c.GetUint("user_id")
username := c.GetString("username")
role := c.GetString("role")
```

### 7. RequireRole 中间件

**功能**: 基于角色的访问控制

```go
// 使用方式
admin := router.Group("/api/v1/admin")
admin.Use(middleware.JWTAuth(jwtManager))
admin.Use(middleware.RequireRole("admin", "superadmin"))
```

**权限检查**:
- 验证用户角色是否在允许列表中
- 支持多角色权限
- 返回 403 Forbidden 错误

### 8. Timeout 中间件

**功能**: 请求超时控制

```go
// 使用方式
router.Use(middleware.Timeout(30 * time.Second))
```

**特性**:
- 设置请求处理超时时间
- 自动取消超时的请求
- 防止长时间运行的请求占用资源

## 中间件配置

### 全局中间件配置

```go
func (r *Router) setupMiddleware() {
    // 恢复中间件（最外层）
    r.engine.Use(middleware.Recovery())
    
    // 日志中间件
    r.engine.Use(middleware.Logger())
    
    // CORS 中间件
    r.engine.Use(middleware.CORS())
    
    // 请求 ID 中间件
    r.engine.Use(middleware.RequestID())
    
    // 限流中间件
    r.engine.Use(middleware.RateLimiter())
    
    // 超时中间件
    r.engine.Use(middleware.Timeout(30 * time.Second))
}
```

### 路由组中间件

```go
// API v1 路由组
v1 := api.Group("/v1")
{
    // 公开路由（无需认证）
    auth := v1.Group("/auth")
    {
        auth.POST("/login", userHandler.Login)
        auth.POST("/register", userHandler.Register)
    }
    
    // 需要认证的路由
    user := v1.Group("/user")
    user.Use(middleware.JWTAuth(jwtManager))
    {
        user.GET("/profile", userHandler.GetProfile)
        user.PUT("/profile", userHandler.UpdateProfile)
    }
    
    // 需要管理员权限的路由
    admin := v1.Group("/admin")
    admin.Use(middleware.JWTAuth(jwtManager))
    admin.Use(middleware.RequireRole("admin"))
    {
        admin.GET("/users", userHandler.GetUsers)
        admin.DELETE("/users/:id", userHandler.DeleteUser)
    }
}
```

## 自定义中间件开发

### 中间件模板

```go
// 自定义中间件模板
func CustomMiddleware() gin.HandlerFunc {
    return gin.HandlerFunc(func(c *gin.Context) {
        // 前置处理
        start := time.Now()
        
        // 调用下一个中间件或处理器
        c.Next()
        
        // 后置处理
        duration := time.Since(start)
        klogger.Info("Request processed", zap.Duration("duration", duration))
    })
}
```

### API 密钥认证中间件

```go
func APIKeyAuth(validKeys []string) gin.HandlerFunc {
    return gin.HandlerFunc(func(c *gin.Context) {
        apiKey := c.GetHeader("X-API-Key")
        if apiKey == "" {
            response.Unauthorized(c, "Missing API key")
            c.Abort()
            return
        }
        
        // 验证 API 密钥
        valid := false
        for _, key := range validKeys {
            if apiKey == key {
                valid = true
                break
            }
        }
        
        if !valid {
            response.Unauthorized(c, "Invalid API key")
            c.Abort()
            return
        }
        
        c.Next()
    })
}
```

### 请求体大小限制中间件

```go
func BodySizeLimit(maxSize int64) gin.HandlerFunc {
    return gin.HandlerFunc(func(c *gin.Context) {
        if c.Request.ContentLength > maxSize {
            response.BadRequest(c, "Request body too large")
            c.Abort()
            return
        }
        
        c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSize)
        c.Next()
    })
}
```

### IP 白名单中间件

```go
func IPWhitelist(allowedIPs []string) gin.HandlerFunc {
    return gin.HandlerFunc(func(c *gin.Context) {
        clientIP := c.ClientIP()
        
        allowed := false
        for _, ip := range allowedIPs {
            if clientIP == ip {
                allowed = true
                break
            }
        }
        
        if !allowed {
            response.Forbidden(c, "IP not allowed")
            c.Abort()
            return
        }
        
        c.Next()
    })
}
```

## 性能优化

### 1. 中间件顺序优化

```go
// 推荐顺序：轻量级中间件在前，重量级中间件在后
router.Use(middleware.Recovery())      // 必须在最外层
router.Use(middleware.CORS())          // 轻量级，处理预检请求
router.Use(middleware.RequestID())     // 轻量级，生成 ID
router.Use(middleware.Logger())        // 中等，记录日志
router.Use(middleware.RateLimiter())   // 重量级，可能涉及外部存储
router.Use(middleware.JWTAuth())       // 重量级，解析和验证令牌
```

### 2. 条件中间件

```go
// 根据环境条件应用中间件
if cfg.App.Debug {
    router.Use(gin.Logger())
    router.Use(gin.Recovery())
} else {
    router.Use(middleware.Logger())
    router.Use(middleware.Recovery())
}
```

### 3. 缓存优化

```go
// 缓存中间件处理结果
var jwtCache = make(map[string]*auth.Claims)
var cacheMutex sync.RWMutex

func CachedJWTAuth(jwtManager *auth.JWTManager) gin.HandlerFunc {
    return gin.HandlerFunc(func(c *gin.Context) {
        token := extractToken(c)
        
        // 检查缓存
        cacheMutex.RLock()
        claims, exists := jwtCache[token]
        cacheMutex.RUnlock()
        
        if !exists {
            // 解析令牌并缓存结果
            var err error
            claims, err = jwtManager.ParseToken(token)
            if err != nil {
                response.Unauthorized(c, "Invalid token")
                c.Abort()
                return
            }
            
            cacheMutex.Lock()
            jwtCache[token] = claims
            cacheMutex.Unlock()
        }
        
        // 设置用户信息到上下文
        c.Set("user_id", claims.UserID)
        c.Set("username", claims.Username)
        c.Set("role", claims.Role)
        
        c.Next()
    })
}
```

## 错误处理

### 中间件错误处理模式

```go
func SafeMiddleware() gin.HandlerFunc {
    return gin.HandlerFunc(func(c *gin.Context) {
        defer func() {
            if err := recover(); err != nil {
                klogger.Error("Middleware panic",
                    zap.Any("error", err),
                    zap.String("middleware", "SafeMiddleware"),
                )
                response.ServerError(c, "Middleware error")
                c.Abort()
            }
        }()
        
        // 中间件逻辑
        c.Next()
    })
}
```

### 错误传播

```go
// 在中间件中设置错误信息
func ValidationMiddleware() gin.HandlerFunc {
    return gin.HandlerFunc(func(c *gin.Context) {
        if err := validateRequest(c); err != nil {
            c.Set("validation_error", err)
            response.BadRequest(c, err.Error())
            c.Abort()
            return
        }
        c.Next()
    })
}

// 在处理器中获取错误信息
func handler(c *gin.Context) {
    if err, exists := c.Get("validation_error"); exists {
        // 处理验证错误
        klogger.Warn("Validation failed", zap.Error(err.(error)))
    }
}
```

## 测试

### 中间件单元测试

```go
func TestCORSMiddleware(t *testing.T) {
    router := gin.New()
    router.Use(middleware.CORS())
    router.GET("/test", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "ok"})
    })
    
    req := httptest.NewRequest("OPTIONS", "/test", nil)
    req.Header.Set("Origin", "http://localhost:3000")
    
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusNoContent, w.Code)
    assert.Equal(t, "http://localhost:3000", w.Header().Get("Access-Control-Allow-Origin"))
}
```

### 中间件集成测试

```go
func TestMiddlewareChain(t *testing.T) {
    router := gin.New()
    
    // 设置中间件链
    router.Use(middleware.Recovery())
    router.Use(middleware.Logger())
    router.Use(middleware.CORS())
    router.Use(middleware.RequestID())
    
    router.GET("/test", func(c *gin.Context) {
        requestID := c.GetString("request_id")
        assert.NotEmpty(t, requestID)
        c.JSON(200, gin.H{"request_id": requestID})
    })
    
    req := httptest.NewRequest("GET", "/test", nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusOK, w.Code)
    assert.NotEmpty(t, w.Header().Get("X-Request-ID"))
}
```

## 监控和指标

### 中间件性能监控

```go
func MetricsMiddleware() gin.HandlerFunc {
    return gin.HandlerFunc(func(c *gin.Context) {
        start := time.Now()
        
        c.Next()
        
        duration := time.Since(start)
        path := c.FullPath()
        method := c.Request.Method
        status := c.Writer.Status()
        
        // 记录指标
        klogger.Info("Request metrics",
            zap.String("method", method),
            zap.String("path", path),
            zap.Int("status", status),
            zap.Duration("duration", duration),
        )
        
        // 发送到监控系统
        // metrics.RecordHTTPRequest(method, path, status, duration)
    })
}
```

## 相关文档

- [JWT 认证系统文档](JWT_AUTH.md)
- [统一响应格式文档](RESPONSE.md)
- [路由管理文档](ROUTER.md)
- [日志系统文档](LOGGING.md)
- [Gin 中间件官方文档](https://gin-gonic.com/docs/examples/using-middleware/)

---

**最佳实践**: 合理安排中间件顺序，轻量级中间件在前，重量级中间件在后；使用条件中间件减少不必要的处理；为自定义中间件添加完善的错误处理和日志记录。