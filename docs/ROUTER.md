# 路由管理文档

## 概述

本项目基于 Gin 框架构建了完整的路由管理系统，采用模块化设计，支持路由分组、中间件集成、版本控制、静态文件服务和 API 文档等功能。路由系统提供了清晰的结构和灵活的配置方式。

## 路由架构

### 路由结构图

```
/
├── /api/v1/                    # API v1 版本
│   ├── /auth/                  # 认证相关路由
│   │   ├── POST /login         # 用户登录
│   │   ├── POST /register      # 用户注册
│   │   └── POST /refresh       # 刷新令牌
│   ├── /user/                  # 用户相关路由（需认证）
│   │   ├── GET /profile        # 获取用户信息
│   │   ├── PUT /profile        # 更新用户信息
│   │   └── DELETE /profile     # 删除用户账户
│   └── /admin/                 # 管理员路由（需管理员权限）
│       ├── GET /users          # 获取用户列表
│       ├── GET /users/:id      # 获取指定用户
│       ├── PUT /users/:id      # 更新用户信息
│       └── DELETE /users/:id   # 删除用户
├── /health                     # 健康检查
├── /docs/                      # API 文档
│   └── /*any                   # Swagger UI
└── /static/                    # 静态文件服务
    └── /*filepath              # 静态资源
```

## 核心组件

### Router 结构体

```go
type Router struct {
    engine     *gin.Engine
    jwtManager *auth.JWTManager
}
```

**字段说明**:
- `engine`: Gin 引擎实例
- `jwtManager`: JWT 管理器，用于认证中间件

### 构造函数

```go
func NewRouter(jwtManager *auth.JWTManager) *Router {
    return &Router{
        engine:     gin.New(),
        jwtManager: jwtManager,
    }
}
```

## 路由设置

### 主设置函数

```go
func (r *Router) Setup() {
    // 设置全局中间件
    r.setupMiddleware()
    
    // 设置 API 路由
    r.setupAPIRoutes()
    
    // 设置健康检查路由
    r.setupHealthRoutes()
    
    // 设置 Swagger 文档路由
    r.setupSwaggerRoutes()
    
    // 设置静态文件路由
    r.setupStaticRoutes()
    
    // 设置 404 处理器
    r.setupNotFoundHandler()
}
```

### 中间件设置

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

### API 路由设置

```go
func (r *Router) setupAPIRoutes() {
    // 创建 API 路由组
    api := r.engine.Group("/api")
    
    // API v1 版本
    v1 := api.Group("/v1")
    {
        // 认证路由（无需认证）
        auth := v1.Group("/auth")
        {
            auth.POST("/login", userHandler.Login)
            auth.POST("/register", userHandler.Register)
            auth.POST("/refresh", userHandler.RefreshToken)
        }
        
        // 用户路由（需要认证）
        user := v1.Group("/user")
        user.Use(middleware.JWTAuth(r.jwtManager))
        {
            user.GET("/profile", userHandler.GetProfile)
            user.PUT("/profile", userHandler.UpdateProfile)
            user.DELETE("/profile", userHandler.DeleteProfile)
        }
        
        // 管理员路由（需要管理员权限）
        admin := v1.Group("/admin")
        admin.Use(middleware.JWTAuth(r.jwtManager))
        admin.Use(middleware.RequireRole("admin"))
        {
            admin.GET("/users", userHandler.GetUsers)
            admin.GET("/users/:id", userHandler.GetUser)
            admin.PUT("/users/:id", userHandler.UpdateUser)
            admin.DELETE("/users/:id", userHandler.DeleteUser)
        }
    }
}
```

### 健康检查路由

```go
func (r *Router) setupHealthRoutes() {
    r.engine.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "status":    "ok",
            "timestamp": time.Now().Unix(),
            "version":   "1.0.0",
        })
    })
}
```

### Swagger 文档路由

```go
func (r *Router) setupSwaggerRoutes() {
    // 开发环境才启用 Swagger
    if gin.Mode() != gin.ReleaseMode {
        r.engine.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
    }
}
```

### 静态文件路由

```go
func (r *Router) setupStaticRoutes() {
    // 静态文件服务
    r.engine.Static("/static", "./web/static")
    
    // 上传文件服务
    r.engine.Static("/uploads", "./uploads")
    
    // 前端应用（SPA）
    r.engine.StaticFile("/", "./web/dist/index.html")
    r.engine.StaticFS("/assets", http.Dir("./web/dist/assets"))
}
```

### 404 处理器

```go
func (r *Router) setupNotFoundHandler() {
    r.engine.NoRoute(func(c *gin.Context) {
        // API 路由返回 JSON 错误
        if strings.HasPrefix(c.Request.URL.Path, "/api/") {
            response.NotFound(c, "API 接口不存在")
            return
        }
        
        // 其他路由返回前端应用
        c.File("./web/dist/index.html")
    })
}
```

## 路由分组策略

### 按功能分组

```go
// 用户管理相关路由
func setupUserRoutes(router *gin.RouterGroup, jwtManager *auth.JWTManager) {
    userHandler := handler.NewUserHandler()
    
    // 公开路由
    public := router.Group("/users")
    {
        public.POST("/register", userHandler.Register)
        public.POST("/login", userHandler.Login)
        public.POST("/forgot-password", userHandler.ForgotPassword)
    }
    
    // 需要认证的路由
    protected := router.Group("/users")
    protected.Use(middleware.JWTAuth(jwtManager))
    {
        protected.GET("/profile", userHandler.GetProfile)
        protected.PUT("/profile", userHandler.UpdateProfile)
        protected.POST("/change-password", userHandler.ChangePassword)
    }
}

// 文章管理相关路由
func setupArticleRoutes(router *gin.RouterGroup, jwtManager *auth.JWTManager) {
    articleHandler := handler.NewArticleHandler()
    
    articles := router.Group("/articles")
    {
        // 公开路由
        articles.GET("", articleHandler.GetList)
        articles.GET("/:id", articleHandler.GetByID)
        
        // 需要认证的路由
        protected := articles.Group("")
        protected.Use(middleware.JWTAuth(jwtManager))
        {
            protected.POST("", articleHandler.Create)
            protected.PUT("/:id", articleHandler.Update)
            protected.DELETE("/:id", articleHandler.Delete)
        }
    }
}
```

### 按权限分组

```go
// 按权限级别分组
func setupPermissionBasedRoutes(api *gin.RouterGroup, jwtManager *auth.JWTManager) {
    // 公开 API（无需认证）
    public := api.Group("/public")
    {
        public.GET("/info", handler.GetSystemInfo)
        public.GET("/announcements", handler.GetAnnouncements)
    }
    
    // 用户 API（需要登录）
    user := api.Group("/user")
    user.Use(middleware.JWTAuth(jwtManager))
    {
        user.GET("/dashboard", handler.GetUserDashboard)
        user.GET("/notifications", handler.GetNotifications)
    }
    
    // 管理员 API（需要管理员权限）
    admin := api.Group("/admin")
    admin.Use(middleware.JWTAuth(jwtManager))
    admin.Use(middleware.RequireRole("admin"))
    {
        admin.GET("/stats", handler.GetSystemStats)
        admin.GET("/logs", handler.GetSystemLogs)
    }
    
    // 超级管理员 API（需要超级管理员权限）
    superAdmin := api.Group("/super-admin")
    superAdmin.Use(middleware.JWTAuth(jwtManager))
    superAdmin.Use(middleware.RequireRole("superadmin"))
    {
        superAdmin.POST("/backup", handler.CreateBackup)
        superAdmin.POST("/restore", handler.RestoreBackup)
    }
}
```

## 版本控制

### API 版本管理

```go
// API 版本控制
func setupVersionedAPI(engine *gin.Engine, jwtManager *auth.JWTManager) {
    api := engine.Group("/api")
    
    // API v1
    v1 := api.Group("/v1")
    {
        setupV1Routes(v1, jwtManager)
    }
    
    // API v2（新版本）
    v2 := api.Group("/v2")
    {
        setupV2Routes(v2, jwtManager)
    }
    
    // 默认版本（指向最新版本）
    api.Any("/*path", func(c *gin.Context) {
        // 重定向到 v2
        newPath := "/api/v2" + c.Param("path")
        c.Redirect(http.StatusMovedPermanently, newPath)
    })
}

func setupV1Routes(v1 *gin.RouterGroup, jwtManager *auth.JWTManager) {
    // V1 版本的路由实现
    users := v1.Group("/users")
    {
        users.GET("", handler.V1GetUsers)      // 旧版本实现
        users.POST("", handler.V1CreateUser)   // 旧版本实现
    }
}

func setupV2Routes(v2 *gin.RouterGroup, jwtManager *auth.JWTManager) {
    // V2 版本的路由实现
    users := v2.Group("/users")
    {
        users.GET("", handler.V2GetUsers)      // 新版本实现
        users.POST("", handler.V2CreateUser)   // 新版本实现
    }
}
```

### 版本兼容性处理

```go
// 版本兼容性中间件
func VersionCompatibility() gin.HandlerFunc {
    return gin.HandlerFunc(func(c *gin.Context) {
        version := c.GetHeader("API-Version")
        if version == "" {
            // 从 URL 路径中提取版本
            path := c.Request.URL.Path
            if strings.HasPrefix(path, "/api/v1") {
                version = "v1"
            } else if strings.HasPrefix(path, "/api/v2") {
                version = "v2"
            } else {
                version = "v2" // 默认最新版本
            }
        }
        
        c.Set("api_version", version)
        c.Next()
    })
}
```

## 路由参数处理

### 路径参数

```go
// 路径参数示例
func setupParameterRoutes(router *gin.RouterGroup) {
    // 单个参数
    router.GET("/users/:id", func(c *gin.Context) {
        userID := c.Param("id")
        // 参数验证
        id, err := strconv.ParseUint(userID, 10, 32)
        if err != nil {
            response.BadRequest(c, "无效的用户ID")
            return
        }
        // 处理逻辑...
    })
    
    // 多个参数
    router.GET("/users/:id/articles/:articleId", func(c *gin.Context) {
        userID := c.Param("id")
        articleID := c.Param("articleId")
        // 处理逻辑...
    })
    
    // 通配符参数
    router.GET("/files/*filepath", func(c *gin.Context) {
        filepath := c.Param("filepath")
        // 文件服务逻辑...
    })
}
```

### 查询参数

```go
// 查询参数处理
func handleQueryParams(c *gin.Context) {
    // 基础查询参数
    page := c.DefaultQuery("page", "1")
    size := c.DefaultQuery("size", "10")
    keyword := c.Query("keyword")
    
    // 参数验证和转换
    pageNum, err := strconv.Atoi(page)
    if err != nil || pageNum < 1 {
        response.BadRequest(c, "无效的页码")
        return
    }
    
    sizeNum, err := strconv.Atoi(size)
    if err != nil || sizeNum < 1 || sizeNum > 100 {
        response.BadRequest(c, "每页大小必须在1-100之间")
        return
    }
    
    // 数组参数
    tags := c.QueryArray("tags")
    
    // 处理逻辑...
}
```

## 中间件集成

### 路由级中间件

```go
// 为特定路由组添加中间件
func setupMiddlewareIntegration(api *gin.RouterGroup, jwtManager *auth.JWTManager) {
    // 文件上传路由（添加文件大小限制）
    upload := api.Group("/upload")
    upload.Use(middleware.JWTAuth(jwtManager))
    upload.Use(middleware.BodySizeLimit(10 << 20)) // 10MB 限制
    {
        upload.POST("/image", handler.UploadImage)
        upload.POST("/document", handler.UploadDocument)
    }
    
    // 支付相关路由（添加 IP 白名单）
    payment := api.Group("/payment")
    payment.Use(middleware.IPWhitelist([]string{"192.168.1.100", "10.0.0.1"}))
    {
        payment.POST("/webhook", handler.PaymentWebhook)
        payment.POST("/notify", handler.PaymentNotify)
    }
    
    // API 密钥认证路由
    apiKey := api.Group("/external")
    apiKey.Use(middleware.APIKeyAuth([]string{"key1", "key2"}))
    {
        apiKey.GET("/data", handler.GetExternalData)
        apiKey.POST("/sync", handler.SyncData)
    }
}
```

### 条件中间件

```go
// 条件中间件应用
func setupConditionalMiddleware(router *gin.RouterGroup) {
    // 开发环境才启用的中间件
    if gin.Mode() == gin.DebugMode {
        router.Use(gin.Logger())
        router.Use(middleware.RequestDump()) // 请求转储中间件
    }
    
    // 生产环境才启用的中间件
    if gin.Mode() == gin.ReleaseMode {
        router.Use(middleware.SecurityHeaders()) // 安全头中间件
        router.Use(middleware.RateLimiter())     // 限流中间件
    }
}
```

## 错误处理

### 路由级错误处理

```go
// 路由错误处理中间件
func RouteErrorHandler() gin.HandlerFunc {
    return gin.HandlerFunc(func(c *gin.Context) {
        defer func() {
            if err := recover(); err != nil {
                klogger.Error("Route panic",
                    zap.Any("error", err),
                    zap.String("path", c.Request.URL.Path),
                    zap.String("method", c.Request.Method),
                )
                response.ServerError(c, "路由处理错误")
            }
        }()
        
        c.Next()
        
        // 检查是否有错误
        if len(c.Errors) > 0 {
            err := c.Errors.Last()
            klogger.Error("Route error",
                zap.Error(err),
                zap.String("path", c.Request.URL.Path),
            )
            
            if !c.Writer.Written() {
                response.ServerError(c, "请求处理失败")
            }
        }
    })
}
```

## 性能优化

### 路由性能优化

```go
// 路由性能优化配置
func optimizeRouter(engine *gin.Engine) {
    // 设置路由缓存
    engine.RedirectTrailingSlash = true
    engine.RedirectFixedPath = true
    
    // 设置最大多部分内存
    engine.MaxMultipartMemory = 8 << 20 // 8MB
    
    // 禁用不必要的功能
    if gin.Mode() == gin.ReleaseMode {
        gin.DisableConsoleColor()
    }
}

// 路由预编译
func precompileRoutes(engine *gin.Engine) {
    // 预编译正则表达式路由
    engine.GET("/users/:id([0-9]+)", handler.GetUser)
    engine.GET("/articles/:slug([a-z0-9-]+)", handler.GetArticle)
}
```

### 路由缓存

```go
// 路由响应缓存中间件
func RouteCache(duration time.Duration) gin.HandlerFunc {
    cache := make(map[string]CacheItem)
    var mutex sync.RWMutex
    
    return gin.HandlerFunc(func(c *gin.Context) {
        // 只缓存 GET 请求
        if c.Request.Method != "GET" {
            c.Next()
            return
        }
        
        key := c.Request.URL.String()
        
        // 检查缓存
        mutex.RLock()
        item, exists := cache[key]
        mutex.RUnlock()
        
        if exists && time.Now().Before(item.ExpireAt) {
            // 返回缓存的响应
            c.Data(item.StatusCode, item.ContentType, item.Body)
            c.Abort()
            return
        }
        
        // 捕获响应
        writer := &responseWriter{ResponseWriter: c.Writer}
        c.Writer = writer
        
        c.Next()
        
        // 缓存响应
        if c.Writer.Status() == http.StatusOK {
            mutex.Lock()
            cache[key] = CacheItem{
                Body:        writer.body.Bytes(),
                StatusCode:  writer.status,
                ContentType: writer.Header().Get("Content-Type"),
                ExpireAt:    time.Now().Add(duration),
            }
            mutex.Unlock()
        }
    })
}
```

## 服务器管理

### 启动服务器

```go
func (r *Router) Run(addr string) error {
    klogger.Info("Starting server", zap.String("address", addr))
    return r.engine.Run(addr)
}

// 优雅关闭
func (r *Router) RunWithGracefulShutdown(addr string) error {
    srv := &http.Server{
        Addr:    addr,
        Handler: r.engine,
    }
    
    // 启动服务器
    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            klogger.Fatal("Server startup failed", zap.Error(err))
        }
    }()
    
    // 等待中断信号
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    
    klogger.Info("Shutting down server...")
    
    // 优雅关闭
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    if err := srv.Shutdown(ctx); err != nil {
        klogger.Error("Server forced to shutdown", zap.Error(err))
        return err
    }
    
    klogger.Info("Server exited")
    return nil
}
```

### 获取引擎实例

```go
func (r *Router) GetEngine() *gin.Engine {
    return r.engine
}
```

## 测试

### 路由测试

```go
func TestRouterSetup(t *testing.T) {
    jwtManager := auth.NewJWTManager("test-secret", time.Hour)
    router := NewRouter(jwtManager)
    router.Setup()
    
    // 测试健康检查路由
    req := httptest.NewRequest("GET", "/health", nil)
    w := httptest.NewRecorder()
    router.GetEngine().ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusOK, w.Code)
    
    var resp map[string]interface{}
    err := json.Unmarshal(w.Body.Bytes(), &resp)
    assert.NoError(t, err)
    assert.Equal(t, "ok", resp["status"])
}

func TestAPIRoutes(t *testing.T) {
    router := setupTestRouter()
    
    tests := []struct {
        method   string
        path     string
        expected int
    }{
        {"GET", "/api/v1/users", http.StatusUnauthorized}, // 需要认证
        {"POST", "/api/v1/auth/login", http.StatusBadRequest}, // 缺少参数
        {"GET", "/health", http.StatusOK},
        {"GET", "/nonexistent", http.StatusNotFound},
    }
    
    for _, test := range tests {
        req := httptest.NewRequest(test.method, test.path, nil)
        w := httptest.NewRecorder()
        router.ServeHTTP(w, req)
        
        assert.Equal(t, test.expected, w.Code,
            "Expected %d for %s %s, got %d",
            test.expected, test.method, test.path, w.Code)
    }
}
```

## 相关文档

- [中间件系统文档](MIDDLEWARE.md)
- [JWT 认证系统文档](JWT_AUTH.md)
- [统一响应格式文档](RESPONSE.md)
- [Gin 框架官方文档](https://gin-gonic.com/docs/)

---

**最佳实践**: 合理组织路由结构，按功能和权限分组；使用中间件实现横切关注点；为不同环境配置不同的路由策略；实现优雅的错误处理和服务器关闭；为路由添加完整的测试覆盖。