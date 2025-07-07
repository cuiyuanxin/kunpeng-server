# 处理器与服务层文档

## 概述

本项目采用分层架构设计，将业务逻辑分为处理器层(Handler)和服务层(Service)。处理器层负责HTTP请求的处理、参数验证和响应格式化，服务层负责具体的业务逻辑实现和数据操作。这种分层设计提高了代码的可维护性、可测试性和可扩展性。

## 系统架构

### 分层架构图

```
请求 → 路由 → 中间件 → 处理器(Handler) → 服务层(Service) → 数据库
                                ↓
                            响应格式化
```

### 核心组件

- **处理器层(Handler)**: HTTP请求处理、参数验证、响应格式化
- **服务层(Service)**: 业务逻辑实现、数据操作、事务管理
- **模型层(Model)**: 数据结构定义、数据传输对象
- **数据库层(Database)**: 数据持久化、查询优化

## 处理器层(Handler)

### 设计原则

1. **单一职责**: 每个处理器只负责特定资源的HTTP操作
2. **参数验证**: 统一的请求参数验证机制
3. **错误处理**: 标准化的错误响应格式
4. **依赖注入**: 通过构造函数注入服务层依赖

### 处理器结构

```go
// UserHandler 用户处理器
type UserHandler struct {
    userService *service.UserService
    jwtManager  *auth.JWTManager
    validator   *validator.Validate
}

// NewUserHandler 创建用户处理器
func NewUserHandler(userService *service.UserService, jwtManager *auth.JWTManager) *UserHandler {
    return &UserHandler{
        userService: userService,
        jwtManager:  jwtManager,
        validator:   validator.New(),
    }
}
```

### 处理器方法模式

```go
// 标准处理器方法模式
func (h *UserHandler) MethodName(c *gin.Context) {
    // 1. 参数绑定
    var req model.RequestStruct
    if err := c.ShouldBindJSON(&req); err != nil {
        response.BadRequest(c, "Invalid request format")
        return
    }
    
    // 2. 参数验证
    if err := h.validator.Struct(&req); err != nil {
        response.BadRequestWithData(c, "Validation failed", err.Error())
        return
    }
    
    // 3. 权限检查（如需要）
    userID, exists := c.Get("user_id")
    if !exists {
        response.Unauthorized(c, "User not authenticated")
        return
    }
    
    // 4. 调用服务层
    result, err := h.userService.SomeMethod(&req)
    if err != nil {
        response.BadRequest(c, err.Error())
        return
    }
    
    // 5. 返回响应
    response.SuccessWithMessage(c, "Operation successful", result)
}
```

### 用户处理器示例

#### 用户注册

```go
// Register 用户注册
func (h *UserHandler) Register(c *gin.Context) {
    var req model.UserCreateRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.BadRequest(c, "Invalid request format")
        return
    }

    // 验证请求参数
    if err := h.validator.Struct(&req); err != nil {
        response.BadRequestWithData(c, "Validation failed", err.Error())
        return
    }

    // 创建用户
    user, err := h.userService.Create(&req)
    if err != nil {
        response.BadRequest(c, err.Error())
        return
    }

    response.SuccessWithMessage(c, "User registered successfully", user.ToResponse())
}
```

#### 用户登录

```go
// Login 用户登录
func (h *UserHandler) Login(c *gin.Context) {
    var req model.UserLoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.BadRequest(c, "Invalid request format")
        return
    }

    // 获取用户并验证密码
    user, err := h.userService.GetByUsername(req.Username)
    if err != nil {
        response.Unauthorized(c, "Invalid username or password")
        return
    }

    if !h.userService.VerifyPassword(user, req.Password) {
        response.Unauthorized(c, "Invalid username or password")
        return
    }

    // 生成JWT令牌
    token, err := h.jwtManager.GenerateToken(user.ID, user.Username, user.Role)
    if err != nil {
        response.ServerError(c, "Failed to generate token")
        return
    }

    // 返回登录响应
    loginResp := model.LoginResponse{
        Token: token,
        User:  user.ToResponse(),
    }

    response.SuccessWithMessage(c, "Login successful", loginResp)
}
```

## 服务层(Service)

### 设计原则

1. **业务逻辑封装**: 将复杂的业务逻辑封装在服务层
2. **数据操作抽象**: 提供高级的数据操作接口
3. **事务管理**: 处理复杂的数据库事务
4. **错误处理**: 提供详细的错误信息

### 服务层结构

```go
// UserService 用户服务
type UserService struct {
    db *gorm.DB
}

// NewUserService 创建用户服务
func NewUserService() *UserService {
    return &UserService{
        db: database.GetDB(),
    }
}
```

### 服务层方法模式

```go
// 标准服务层方法模式
func (s *UserService) MethodName(params interface{}) (*Model, error) {
    // 1. 参数验证
    if params == nil {
        return nil, errors.New("invalid parameters")
    }
    
    // 2. 业务逻辑处理
    // ...
    
    // 3. 数据库操作
    var result Model
    if err := s.db.Create(&result).Error; err != nil {
        return nil, fmt.Errorf("failed to create: %w", err)
    }
    
    return &result, nil
}
```

### 用户服务示例

#### 创建用户

```go
// Create 创建用户
func (s *UserService) Create(req *model.UserCreateRequest) (*model.User, error) {
    // 检查用户名是否已存在
    var existUser model.User
    if err := s.db.Where("username = ? OR email = ?", req.Username, req.Email).First(&existUser).Error; err == nil {
        if existUser.Username == req.Username {
            return nil, errors.New("username already exists")
        }
        if existUser.Email == req.Email {
            return nil, errors.New("email already exists")
        }
    }

    // 加密密码
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, fmt.Errorf("failed to hash password: %w", err)
    }

    // 创建用户
    user := &model.User{
        Username: req.Username,
        Email:    req.Email,
        Password: string(hashedPassword),
        Nickname: req.Nickname,
        Phone:    req.Phone,
        Role:     req.Role,
        Status:   1, // 默认启用
    }

    if user.Role == "" {
        user.Role = "user" // 默认角色
    }

    if err := s.db.Create(user).Error; err != nil {
        return nil, fmt.Errorf("failed to create user: %w", err)
    }

    return user, nil
}
```

#### 用户列表查询

```go
// List 获取用户列表
func (s *UserService) List(req *model.UserListRequest) ([]model.User, int64, error) {
    var users []model.User
    var total int64

    // 构建查询
    query := s.db.Model(&model.User{})

    // 关键词搜索
    if req.Keyword != "" {
        query = query.Where("username LIKE ? OR email LIKE ? OR nickname LIKE ?",
            "%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
    }

    // 角色筛选
    if req.Role != "" {
        query = query.Where("role = ?", req.Role)
    }

    // 状态筛选
    if req.Status != nil {
        query = query.Where("status = ?", *req.Status)
    }

    // 获取总数
    if err := query.Count(&total).Error; err != nil {
        return nil, 0, fmt.Errorf("failed to count users: %w", err)
    }

    // 分页查询
    if err := query.Scopes(database.Paginate(req.Page, req.PageSize)).Find(&users).Error; err != nil {
        return nil, 0, fmt.Errorf("failed to get users: %w", err)
    }

    return users, total, nil
}
```

## 应用程序启动

### 主程序结构

```go
func main() {
    // 1. 解析命令行参数
    configPath := flag.String("config", "configs/config.yaml", "配置文件路径")
    flag.Parse()

    // 2. 初始化配置
    cfg, err := config.Init(*configPath)
    if err != nil {
        panic(fmt.Sprintf("Failed to init config: %v", err))
    }

    // 3. 初始化日志
    if err := klogger.InitWithEnvironment(&cfg.Logging, cfg.App.Environment); err != nil {
        panic(fmt.Sprintf("Failed to init logger: %v", err))
    }
    defer klogger.Sync()

    // 4. 初始化数据库
    if err := database.InitWithConfig(cfg); err != nil {
        klogger.Fatal("Failed to init database", zap.Error(err))
    }
    defer database.Close()

    // 5. 初始化Redis
    if err := redis.Init(&cfg.Redis); err != nil {
        klogger.Fatal("Failed to init redis", zap.Error(err))
    }
    defer redis.Close()

    // 6. 使用 Wire 初始化应用程序
    app := wire.InitializeApp(cfg, database.GetDB())

    // 7. 使用 Wire 应用程序创建路由
    r := router.NewRouterWithWire(app)
    r.Setup()

    // 7. 启动HTTP服务器
    srv := &http.Server{
        Addr:         cfg.Server.GetServerAddr(),
        Handler:      r.GetEngine(),
        ReadTimeout:  cfg.Server.ReadTimeout,
        WriteTimeout: cfg.Server.WriteTimeout,
        IdleTimeout:  cfg.Server.IdleTimeout,
    }

    // 8. 优雅关闭
    // ...
}
```

### 应用程序生命周期

1. **初始化阶段**: 配置加载、日志初始化、数据库连接、Redis连接
2. **启动阶段**: 路由设置、中间件注册、服务器启动
3. **运行阶段**: 请求处理、配置热更新、健康检查
4. **关闭阶段**: 优雅关闭、资源清理、连接释放

## 依赖注入

### 服务注册

```go
// 在路由初始化时注册服务
func (r *Router) setupUserRoutes() {
    // 创建服务实例
    userService := service.NewUserService()
    jwtManager := auth.NewJWTManager(r.config.JWT.Secret)
    
    // 创建处理器实例
    userHandler := handler.NewUserHandler(userService, jwtManager)
    
    // 注册路由
    userGroup := r.engine.Group("/api/v1/user")
    {
        userGroup.POST("/register", userHandler.Register)
        userGroup.POST("/login", userHandler.Login)
        userGroup.GET("/profile", middleware.JWTAuth(), userHandler.GetProfile)
        userGroup.PUT("/profile", middleware.JWTAuth(), userHandler.UpdateProfile)
    }
}
```

## 错误处理

### 统一错误处理

```go
// 服务层错误处理
func (s *UserService) GetByID(id uint) (*model.User, error) {
    var user model.User
    if err := s.db.First(&user, id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("user not found")
        }
        return nil, fmt.Errorf("failed to get user: %w", err)
    }
    return &user, nil
}

// 处理器层错误处理
func (h *UserHandler) GetProfile(c *gin.Context) {
    userID, exists := c.Get("user_id")
    if !exists {
        response.Unauthorized(c, "User not authenticated")
        return
    }

    user, err := h.userService.GetByID(userID.(uint))
    if err != nil {
        if err.Error() == "user not found" {
            response.NotFound(c, "User not found")
        } else {
            response.ServerError(c, "Internal server error")
        }
        return
    }

    response.Success(c, user.ToResponse())
}
```

## 事务处理

### 服务层事务

```go
// 复杂业务逻辑的事务处理
func (s *UserService) CreateUserWithProfile(req *CreateUserWithProfileRequest) error {
    return database.Transaction(func(tx *gorm.DB) error {
        // 创建用户
        user := &model.User{
            Username: req.Username,
            Email:    req.Email,
            // ...
        }
        if err := tx.Create(user).Error; err != nil {
            return err
        }

        // 创建用户资料
        profile := &model.UserProfile{
            UserID: user.ID,
            Bio:    req.Bio,
            // ...
        }
        if err := tx.Create(profile).Error; err != nil {
            return err
        }

        return nil
    })
}
```

## 测试

### 服务层测试

```go
func TestUserService_Create(t *testing.T) {
    // 设置测试数据库
    db := setupTestDB()
    defer cleanupTestDB(db)
    
    service := &UserService{db: db}
    
    req := &model.UserCreateRequest{
        Username: "testuser",
        Email:    "test@example.com",
        Password: "password123",
    }
    
    user, err := service.Create(req)
    assert.NoError(t, err)
    assert.Equal(t, "testuser", user.Username)
    assert.Equal(t, "test@example.com", user.Email)
}
```

### 处理器测试

```go
func TestUserHandler_Register(t *testing.T) {
    // 设置测试环境
    gin.SetMode(gin.TestMode)
    
    // 创建模拟服务
    mockService := &MockUserService{}
    handler := NewUserHandler(mockService, nil)
    
    // 创建测试请求
    reqBody := `{"username":"testuser","email":"test@example.com","password":"password123"}`
    req := httptest.NewRequest("POST", "/register", strings.NewReader(reqBody))
    req.Header.Set("Content-Type", "application/json")
    
    // 执行测试
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Request = req
    
    handler.Register(c)
    
    // 验证结果
    assert.Equal(t, http.StatusOK, w.Code)
}
```

## 性能优化

### 数据库查询优化

```go
// 使用预加载减少N+1查询
func (s *UserService) GetUsersWithProfiles() ([]model.User, error) {
    var users []model.User
    if err := s.db.Preload("Profile").Find(&users).Error; err != nil {
        return nil, err
    }
    return users, nil
}

// 使用索引优化查询
func (s *UserService) GetByUsernameOptimized(username string) (*model.User, error) {
    var user model.User
    // 确保username字段有索引
    if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
        return nil, err
    }
    return &user, nil
}
```

### 缓存策略

```go
// 在服务层添加缓存
func (s *UserService) GetByIDWithCache(id uint) (*model.User, error) {
    // 先从缓存获取
    cacheKey := fmt.Sprintf("user:%d", id)
    if cached := redis.Get(cacheKey); cached != "" {
        var user model.User
        if err := json.Unmarshal([]byte(cached), &user); err == nil {
            return &user, nil
        }
    }
    
    // 从数据库获取
    user, err := s.GetByID(id)
    if err != nil {
        return nil, err
    }
    
    // 存入缓存
    if data, err := json.Marshal(user); err == nil {
        redis.Set(cacheKey, string(data), time.Hour)
    }
    
    return user, nil
}
```

## 最佳实践

### 1. 代码组织

- 按功能模块组织处理器和服务
- 保持处理器轻量，业务逻辑放在服务层
- 使用接口定义服务层契约

### 2. 错误处理

- 服务层返回详细错误信息
- 处理器层转换为用户友好的错误响应
- 使用错误包装提供上下文信息

### 3. 安全考虑

- 输入验证和清理
- 权限检查
- 敏感数据加密
- SQL注入防护

### 4. 性能优化

- 数据库查询优化
- 适当使用缓存
- 连接池配置
- 分页查询

### 5. 测试策略

- 单元测试覆盖服务层
- 集成测试覆盖处理器层
- 使用模拟对象隔离依赖

## 相关文档

- [路由管理文档](ROUTER.md)
- [中间件系统文档](MIDDLEWARE.md)
- [数据模型文档](MODELS.md)
- [统一响应格式文档](RESPONSE.md)
- [JWT认证文档](JWT_AUTH.md)
- [数据库系统完整指南](DATABASE_GUIDE.md)