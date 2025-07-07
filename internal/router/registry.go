package router

import (
	"github.com/gin-gonic/gin"
	"github.com/cuiyuanxin/kunpeng/internal/controller"
	"github.com/cuiyuanxin/kunpeng/internal/middleware"
)

// RouteRegistry 路由注册器
type RouteRegistry struct {
	engine         *gin.Engine
	userController *controller.UserController
	adminController *controller.AdminController
	authMiddleware *middleware.AuthMiddleware
	casbinMiddleware *middleware.CasbinMiddleware
}

// NewRouteRegistry 创建路由注册器
func NewRouteRegistry(
	engine *gin.Engine,
	userController *controller.UserController,
	adminController *controller.AdminController,
	authMiddleware *middleware.AuthMiddleware,
	casbinMiddleware *middleware.CasbinMiddleware,
) *RouteRegistry {
	return &RouteRegistry{
		engine:           engine,
		userController:   userController,
		adminController:  adminController,
		authMiddleware:   authMiddleware,
		casbinMiddleware: casbinMiddleware,
	}
}

// RegisterRoutes 注册所有路由
func (rr *RouteRegistry) RegisterRoutes() {
	// 注册API路由
	rr.registerAPIRoutes()
	
	// 注册健康检查路由
	rr.registerHealthRoutes()
	
	// 注册Swagger路由
	rr.registerSwaggerRoutes()
	
	// 注册静态文件路由
	rr.registerStaticRoutes()
	
	// 注册404处理
	rr.registerNotFoundHandler()
}

// registerAPIRoutes 注册API路由
func (rr *RouteRegistry) registerAPIRoutes() {
	api := rr.engine.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			// 认证相关路由（无需认证）
			auth := v1.Group("/auth")
			{
				auth.POST("/register", rr.userController.Register)
				auth.POST("/login", rr.userController.Login)
			}

			// 用户相关路由（需要认证）
			user := v1.Group("/user")
			user.Use(rr.authMiddleware.RequireAuth())
			{
				user.GET("/profile", rr.userController.GetProfile)
				user.PUT("/profile", rr.userController.UpdateProfile)
				user.POST("/change-password", rr.userController.ChangePassword)
			}

			// 管理员相关路由（需要认证和管理员权限）
			admin := v1.Group("/admin")
			admin.Use(rr.authMiddleware.RequireAuth())
			admin.Use(rr.casbinMiddleware.RequirePermission("admin"))
			{
				// 用户管理
				users := admin.Group("/users")
				{
					users.GET("", rr.adminController.GetUsers)
					users.POST("", rr.adminController.CreateUser)
					users.PUT("/:id", rr.adminController.UpdateUser)
					users.DELETE("/:id", rr.adminController.DeleteUser)
				}

				// 角色管理
				roles := admin.Group("/roles")
				{
					roles.GET("", rr.adminController.GetRoles)
					roles.POST("", rr.adminController.CreateRole)
				}

				// 权限管理
				permissions := admin.Group("/permissions")
				{
					permissions.GET("", rr.adminController.GetPermissions)
				}

				// 部门管理
				departments := admin.Group("/departments")
				{
					departments.GET("", rr.adminController.GetDepartments)
					departments.POST("", rr.adminController.CreateDepartment)
				}

				// 文件管理
				files := admin.Group("/files")
				{
					files.GET("", rr.adminController.GetFiles)
					files.POST("/upload", rr.adminController.UploadFile)
				}

				// 统计分析
				stats := admin.Group("/stats")
				{
					stats.GET("/dashboard", rr.adminController.GetDashboardStats)
					stats.GET("/users", rr.adminController.GetUserStats)
				}
			}
		}
	}
}

// registerHealthRoutes 注册健康检查路由
func (rr *RouteRegistry) registerHealthRoutes() {
	health := rr.engine.Group("/health")
	{
		health.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "pong"})
		})
		health.GET("/status", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": "ok",
				"service": "kunpeng-server",
				"version": "1.0.0",
			})
		})
	}
}

// registerSwaggerRoutes 注册Swagger路由
func (rr *RouteRegistry) registerSwaggerRoutes() {
	// Swagger文档路由
	rr.engine.GET("/swagger/*any", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Swagger documentation"})
	})
}

// registerStaticRoutes 注册静态文件路由
func (rr *RouteRegistry) registerStaticRoutes() {
	// 静态文件服务
	rr.engine.Static("/static", "./static")
	rr.engine.Static("/uploads", "./uploads")
}

// registerNotFoundHandler 注册404处理
func (rr *RouteRegistry) registerNotFoundHandler() {
	rr.engine.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"code":    404,
			"message": "Route not found",
			"path":    c.Request.URL.Path,
		})
	})
}