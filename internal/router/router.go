package router

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"

	"github.com/cuiyuanxin/kunpeng/internal/config"
	"github.com/cuiyuanxin/kunpeng/internal/controller"
	"github.com/cuiyuanxin/kunpeng/internal/logger"
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
	// 根据环境选择不同的恢复和日志中间件
	if r.config.App.Environment == "production" {
		// 生产环境使用自定义中间件
		r.engine.Use(middleware.CustomRecovery(logger.Logger, func(c *gin.Context, msg string) {
			response.ServerError(c, msg)
		}))
		r.engine.Use(middleware.CustomLogger(logger.Logger))
	} else {
		// 开发环境使用gin框架自带的中间件
		r.engine.Use(gin.Recovery())
		r.engine.Use(gin.Logger())
	}

	// CORS中间件
	r.engine.Use(pkgMiddleware.CORS())

	// 限流中间件
	r.engine.Use(pkgMiddleware.RateLimiter())

	// 请求ID中间件
	r.engine.Use(pkgMiddleware.RequestID())
}

// GetEngine 获取Gin引擎
func (r *Router) GetEngine() *gin.Engine {
	return r.engine
}

// Run 启动服务器
func (r *Router) Run() error {
	return r.engine.Run(r.config.Server.GetServerAddr())
}
