package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/files"

	"github.com/cuiyuanxin/kunpeng/internal/config"
	"github.com/cuiyuanxin/kunpeng/internal/controller"
	"github.com/cuiyuanxin/kunpeng/internal/middleware"
	"github.com/cuiyuanxin/kunpeng/pkg/auth"
	pkgMiddleware "github.com/cuiyuanxin/kunpeng/pkg/middleware"
	"github.com/cuiyuanxin/kunpeng/pkg/response"
)

// Router 路由管理器
type Router struct {
	engine        *gin.Engine
	config        *config.Config
	routeRegistry *RouteRegistry
	userHandler   *controller.UserController
	jwtManager    *auth.JWTManager
	db            *gorm.DB
	validator     *validator.Validate
}

// NewRouter 创建路由管理器（使用路由注册器）
func NewRouter(cfg *config.Config, engine *gin.Engine, routeRegistry *RouteRegistry, userHandler *controller.UserController, jwtManager *auth.JWTManager, db *gorm.DB, validator *validator.Validate) *Router {
	return &Router{
		engine:        engine,
		config:        cfg,
		routeRegistry: routeRegistry,
		userHandler:   userHandler,
		jwtManager:    jwtManager,
		db:            db,
		validator:     validator,
	}
}





// Setup 设置路由
func (r *Router) Setup() {
	// 设置中间件
	r.setupMiddleware()

	// 注册所有路由
	r.routeRegistry.RegisterRoutes()
}

// setupMiddleware 设置中间件
func (r *Router) setupMiddleware() {
	// 恢复中间件
	r.engine.Use(gin.Recovery())

	// 日志中间件
	r.engine.Use(gin.Logger())

	// CORS中间件
	r.engine.Use(pkgMiddleware.CORS())

	// 限流中间件
	r.engine.Use(pkgMiddleware.RateLimiter())

	// 请求ID中间件
	r.engine.Use(pkgMiddleware.RequestID())
}

// setupAPIRoutes 设置API路由
func (r *Router) setupAPIRoutes() {
	api := r.engine.Group("/api")
	{
		// V1 API
		v1 := api.Group("/v1")
		{
			// 认证相关路由
			auth := v1.Group("/auth")
			{
				auth.POST("/login", r.userHandler.Login)
				auth.POST("/register", r.userHandler.Register)
			}

			// 用户相关路由（需要认证）
			user := v1.Group("/user")
			user.Use(middleware.JWTAuth(r.jwtManager, r.db))
			{
				user.GET("/profile", r.userHandler.GetProfile)
				user.PUT("/profile", r.userHandler.UpdateProfile)
				user.POST("/change-password", r.userHandler.ChangePassword)
			}


		}
	}
}

// setupAdminRoutes 设置管理员路由
func (r *Router) setupAdminRoutes() {
	// 管理员路由设置
	admin := r.engine.Group("/api/v1/admin")
	admin.Use(middleware.JWTAuth(r.jwtManager, r.db))
	{
		// 这里可以添加管理员相关的路由
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
