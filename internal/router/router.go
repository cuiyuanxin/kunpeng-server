package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/cuiyuanxin/kunpeng/internal/auth"
	"github.com/cuiyuanxin/kunpeng/internal/config"
	"github.com/cuiyuanxin/kunpeng/internal/handler"
	"github.com/cuiyuanxin/kunpeng/internal/middleware"
	"github.com/cuiyuanxin/kunpeng/internal/response"
	"github.com/cuiyuanxin/kunpeng/internal/service"
)

// Router 路由管理器
type Router struct {
	engine      *gin.Engine
	config      *config.Config
	userHandler *handler.UserHandler
	jwtManager  *auth.JWTManager
}

// NewRouter 创建路由管理器
func NewRouter(cfg *config.Config) *Router {
	// 设置Gin模式
	if !cfg.App.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建Gin引擎
	engine := gin.New()

	// 创建JWT管理器
	jwtManager := auth.NewJWTManager(&cfg.JWT)

	// 创建服务
	userService := service.NewUserService()

	// 创建处理器
	userHandler := handler.NewUserHandler(userService, jwtManager)

	return &Router{
		engine:      engine,
		config:      cfg,
		userHandler: userHandler,
		jwtManager:  jwtManager,
	}
}

// Setup 设置路由
func (r *Router) Setup() {
	// 设置全局中间件
	r.setupMiddleware()

	// 设置API路由
	r.setupAPIRoutes()

	// 设置健康检查路由
	r.setupHealthRoutes()

	// 设置Swagger文档路由
	r.setupSwaggerRoutes()

	// 设置静态文件路由
	r.setupStaticRoutes()

	// 设置404处理
	r.setupNotFoundHandler()
}

// setupMiddleware 设置中间件
func (r *Router) setupMiddleware() {
	// 恢复中间件
	r.engine.Use(middleware.Recovery())

	// 日志中间件
	r.engine.Use(middleware.Logger())

	// CORS中间件
	r.engine.Use(middleware.CORS())

	// 请求ID中间件
	r.engine.Use(middleware.RequestID())

	// 限流中间件
	r.engine.Use(middleware.RateLimiter())
}

// setupAPIRoutes 设置API路由
func (r *Router) setupAPIRoutes() {
	api := r.engine.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			// 认证相关路由（无需JWT验证）
			auth := v1.Group("/auth")
			{
				auth.POST("/register", r.userHandler.Register)
				auth.POST("/login", r.userHandler.Login)
			}

			// 用户相关路由（需要JWT验证）
			user := v1.Group("/user")
			user.Use(middleware.JWTAuth(r.jwtManager))
			{
				user.GET("/profile", r.userHandler.GetProfile)
				user.PUT("/profile", r.userHandler.UpdateProfile)
				user.PUT("/password", r.userHandler.ChangePassword)
			}

			// 管理员相关路由（需要JWT验证和管理员权限）
			admin := v1.Group("/admin")
			admin.Use(middleware.JWTAuth(r.jwtManager))
			admin.Use(middleware.RequireRole("admin"))
			{
				// 用户管理
				users := admin.Group("/users")
				{
					users.GET("", r.userHandler.GetUsers)
					users.GET("/:id", r.userHandler.GetUser)
					users.PUT("/:id", r.userHandler.UpdateUser)
					users.DELETE("/:id", r.userHandler.DeleteUser)
				}
			}
		}
	}
}

// setupHealthRoutes 设置健康检查路由
func (r *Router) setupHealthRoutes() {
	health := r.engine.Group("/health")
	{
		health.GET("/ping", func(c *gin.Context) {
			response.Success(c, gin.H{
				"status":  "ok",
				"message": "pong",
			})
		})

		health.GET("/check", func(c *gin.Context) {
			// 这里可以添加数据库、Redis等健康检查
			response.Success(c, gin.H{
				"status":   "healthy",
				"database": "ok",
				"redis":    "ok",
				"version":  r.config.App.Version,
			})
		})
	}
}

// setupSwaggerRoutes 设置Swagger文档路由
func (r *Router) setupSwaggerRoutes() {
	if r.config.App.Debug {
		r.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
}

// setupStaticRoutes 设置静态文件路由
func (r *Router) setupStaticRoutes() {
	// 静态文件服务
	r.engine.Static("/static", "./web/static")
	r.engine.StaticFile("/favicon.ico", "./web/static/favicon.ico")

	// 上传文件服务
	r.engine.Static("/uploads", "./uploads")
}

// setupNotFoundHandler 设置404处理
func (r *Router) setupNotFoundHandler() {
	r.engine.NoRoute(func(c *gin.Context) {
		response.NotFound(c, "Route not found")
	})

	r.engine.NoMethod(func(c *gin.Context) {
		response.Custom(c, http.StatusMethodNotAllowed, 405, "Method not allowed", nil)
	})
}

// GetEngine 获取Gin引擎
func (r *Router) GetEngine() *gin.Engine {
	return r.engine
}

// Run 启动服务器
func (r *Router) Run() error {
	return r.engine.Run(r.config.Server.GetServerAddr())
}
