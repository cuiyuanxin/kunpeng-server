//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"gorm.io/gorm"

	"github.com/cuiyuanxin/kunpeng/internal/config"
	"github.com/cuiyuanxin/kunpeng/internal/controller"
	"github.com/cuiyuanxin/kunpeng/internal/middleware"
	"github.com/cuiyuanxin/kunpeng/internal/router"
	"github.com/cuiyuanxin/kunpeng/internal/service"
	"github.com/cuiyuanxin/kunpeng/pkg/auth"
)

// App 应用程序结构体
type App struct {
	Config              *config.Config
	DB                  *gorm.DB
	Engine              *gin.Engine
	JWTManager          *auth.JWTManager
	UserService         *service.UserService
	RoleService         *service.RoleService
	PermissionService   *service.PermissionService
	DepartmentService   *service.DepartmentService
	FileService         *service.FileService
	StatService         *service.StatService
	UserController      *controller.UserController
	AdminController     *controller.AdminController
	AuthMiddleware      *middleware.AuthMiddleware
	CasbinMiddleware    *middleware.CasbinMiddleware
	RouteRegistry       *router.RouteRegistry
}

// ProvideJWTManager 提供JWT管理器
func ProvideJWTManager(cfg *config.Config) *auth.JWTManager {
	jwtConfig := &auth.JWTConfig{
		Secret:     cfg.JWT.Secret,
		ExpireTime: cfg.JWT.ExpireTime,
		Issuer:     cfg.JWT.Issuer,
	}
	return auth.NewJWTManager(jwtConfig)
}

// ProvideUserService 提供用户服务
func ProvideUserService(db *gorm.DB) *service.UserService {
	return service.NewUserServiceWithDB(db)
}

// ProvideGinEngine 提供Gin引擎
func ProvideGinEngine(cfg *config.Config) *gin.Engine {
	if !cfg.App.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	return gin.New()
}

// ProvideRoleService 提供角色服务
func ProvideRoleService(db *gorm.DB) *service.RoleService {
	return service.NewRoleService(db)
}

// ProvidePermissionService 提供权限服务
func ProvidePermissionService(db *gorm.DB) *service.PermissionService {
	return service.NewPermissionService(db)
}

// ProvideDepartmentService 提供部门服务
func ProvideDepartmentService(db *gorm.DB) *service.DepartmentService {
	return service.NewDepartmentService(db)
}

// ProvideFileService 提供文件服务
func ProvideFileService(cfg *config.Config, db *gorm.DB) *service.FileService {
	return service.NewFileService(cfg, db)
}

// ProvideStatService 提供统计服务
func ProvideStatService(db *gorm.DB) *service.StatService {
	return service.NewStatService(db)
}

// ProvideUserController 提供用户控制器
func ProvideUserController(userService *service.UserService, jwtManager *auth.JWTManager) *controller.UserController {
	return controller.NewUserController(userService, jwtManager)
}

// ProvideAdminController 提供管理员控制器
func ProvideAdminController(
	userService *service.UserService,
	roleService *service.RoleService,
	permissionService *service.PermissionService,
	departmentService *service.DepartmentService,
	fileService *service.FileService,
	statService *service.StatService,
) *controller.AdminController {
	return controller.NewAdminController(
		userService,
		roleService,
		permissionService,
		departmentService,
		fileService,
		statService,
	)
}

// ProvideAuthMiddleware 提供认证中间件
func ProvideAuthMiddleware(jwtManager *auth.JWTManager, db *gorm.DB) *middleware.AuthMiddleware {
	return middleware.NewAuthMiddleware(jwtManager, db)
}

// ProvideCasbinMiddleware 提供Casbin中间件
func ProvideCasbinMiddleware(permissionService *service.PermissionService) *middleware.CasbinMiddleware {
	return middleware.NewCasbinMiddleware(permissionService)
}

// ProvideRouteRegistry 提供路由注册器
func ProvideRouteRegistry(
	engine *gin.Engine,
	userController *controller.UserController,
	adminController *controller.AdminController,
	authMiddleware *middleware.AuthMiddleware,
	casbinMiddleware *middleware.CasbinMiddleware,
) *router.RouteRegistry {
	return router.NewRouteRegistry(
		engine,
		userController,
		adminController,
		authMiddleware,
		casbinMiddleware,
	)
}

// ProvideApp 提供应用程序实例
func ProvideApp(
	cfg *config.Config,
	db *gorm.DB,
	engine *gin.Engine,
	jwtManager *auth.JWTManager,
	userService *service.UserService,
	roleService *service.RoleService,
	permissionService *service.PermissionService,
	departmentService *service.DepartmentService,
	fileService *service.FileService,
	statService *service.StatService,
	userController *controller.UserController,
	adminController *controller.AdminController,
	authMiddleware *middleware.AuthMiddleware,
	casbinMiddleware *middleware.CasbinMiddleware,
	routeRegistry *router.RouteRegistry,
) *App {
	return &App{
		Config:            cfg,
		DB:                db,
		Engine:            engine,
		JWTManager:        jwtManager,
		UserService:       userService,
		RoleService:       roleService,
		PermissionService: permissionService,
		DepartmentService: departmentService,
		FileService:       fileService,
		StatService:       statService,
		UserController:    userController,
		AdminController:   adminController,
		AuthMiddleware:    authMiddleware,
		CasbinMiddleware:  casbinMiddleware,
		RouteRegistry:     routeRegistry,
	}
}

// Wire集合定义
var AppSet = wire.NewSet(
	ProvideGinEngine,
	ProvideJWTManager,
	ProvideUserService,
	ProvideRoleService,
	ProvidePermissionService,
	ProvideDepartmentService,
	ProvideFileService,
	ProvideStatService,
	ProvideUserController,
	ProvideAdminController,
	ProvideAuthMiddleware,
	ProvideCasbinMiddleware,
	ProvideRouteRegistry,
	ProvideApp,
)

// InitializeApp 初始化应用程序（Wire将生成此函数的实现）
func InitializeApp(cfg *config.Config, db *gorm.DB) *App {
	wire.Build(AppSet)
	return &App{}
}