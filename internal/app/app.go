package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cuiyuanxin/kunpeng/internal/middleware"
	"github.com/cuiyuanxin/kunpeng/internal/repository"
	"github.com/cuiyuanxin/kunpeng/pkg/casbin"
	"github.com/cuiyuanxin/kunpeng/pkg/config"
	"github.com/cuiyuanxin/kunpeng/pkg/database"
	"github.com/cuiyuanxin/kunpeng/pkg/logger"
	"github.com/cuiyuanxin/kunpeng/pkg/tracer"
	"github.com/cuiyuanxin/kunpeng/pkg/validator"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// App 应用程序
type App struct {
	engine *gin.Engine
	server *http.Server
}

// New 创建应用程序
func New(configPath string) *App {
	// 初始化配置
	if err := config.Init(configPath); err != nil {
		panic(fmt.Sprintf("初始化配置失败: %v", err))
	}

	// 初始化日志
	_ = logger.Init()

	// 初始化数据库
	if err := database.Init(); err != nil {
		panic(fmt.Sprintf("初始化数据库失败: %v", err))
	}

	// 初始化仓储层
	repository.Init()

	// 初始化验证器（使用i18n支持）
	validator.InitI18n()

	// 初始化Casbin
	if err := casbin.Init(); err != nil {
		panic(fmt.Sprintf("初始化Casbin失败: %v", err))
	}

	// 初始化链路追踪（如果启用）
	if config.GetAppConfig().TraceEnable {
		tracer.Init()
	}

	// 设置运行模式
	appMode := config.Get().App.Mode
	var ginMode string
	switch appMode {
	case "production":
		ginMode = gin.ReleaseMode
	case "test":
		ginMode = gin.TestMode
	default:
		ginMode = gin.DebugMode
	}
	gin.SetMode(ginMode)

	// 创建Gin引擎
	engine := gin.New()

	// 注册中间件
	registerMiddlewares(engine)

	return &App{
		engine: engine,
	}
}

// Run 运行应用程序
func (a *App) Run() {
	// 注册路由
	a.RegisterRoutes()

	// 获取服务器配置
	serverConfig := config.Get().Server

	// 创建HTTP服务器
	a.server = &http.Server{
		Addr:           fmt.Sprintf("%s:%d", serverConfig.Host, serverConfig.Port),
		Handler:        a.engine,
		ReadTimeout:    serverConfig.ReadTimeout,
		WriteTimeout:   serverConfig.WriteTimeout,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	// 优雅启动HTTP服务器
	logger.GetLogger().Info("正在启动服务器...",
		zap.String("地址", a.server.Addr),
		zap.String("模式", gin.Mode()),
	)

	go func() {
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.GetLogger().Fatal("启动服务器失败", zap.Error(err))
		}
	}()

	logger.GetLogger().Info("服务器已启动", zap.String("地址", a.server.Addr))

	// 优雅关闭
	a.gracefulShutdown()
}

// gracefulShutdown 优雅关闭
func (a *App) gracefulShutdown() {
	// 等待中断信号
	quit := make(chan os.Signal, 1)
	// 监听 SIGINT, SIGTERM, SIGQUIT 信号
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	sig := <-quit

	logger.GetLogger().Info("接收到关闭信号", zap.String("信号", sig.String()))

	// 设置关闭超时时间
	timeout := 10 * time.Second
	logger.GetLogger().Info("正在关闭服务器...", zap.Duration("超时时间", timeout))

	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// 关闭HTTP服务器
	if err := a.server.Shutdown(ctx); err != nil {
		logger.GetLogger().Error("服务器关闭异常", zap.Error(err))
	}

	// 等待未完成的请求处理完毕
	select {
	case <-ctx.Done():
		logger.GetLogger().Info("关闭超时")
	default:
		logger.GetLogger().Info("服务器已正常关闭")
	}

	// 关闭数据库连接
	logger.GetLogger().Info("正在关闭数据库连接...")
	database.Close()

	logger.GetLogger().Info("应用程序已完全关闭")
	os.Exit(0)
}

// registerMiddlewares 注册中间件
func registerMiddlewares(engine *gin.Engine) {
	// 注册全局中间件
	if config.IsProduction() {
		engine.Use(middleware.Recovery())
		engine.Use(middleware.Logger())
	} else {
		engine.Use(gin.Recovery())
		engine.Use(gin.Logger())
	}
	// 根据配置决定是否启用链路追踪中间件
	if config.GetAppConfig().TraceEnable {
		engine.Use(middleware.Tracer())
	}
	// 暂时禁用超时中间件，因为它会导致所有请求超时
	engine.Use(middleware.TimeoutWithDefault())
	engine.Use(middleware.RateLimiter())
	// 操作日志中间件
	engine.Use(middleware.OperationLog())
}
