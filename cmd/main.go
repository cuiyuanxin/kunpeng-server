package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/cuiyuanxin/kunpeng/internal/config"
	"github.com/cuiyuanxin/kunpeng/internal/database"
	klogger "github.com/cuiyuanxin/kunpeng/internal/logger"
	"github.com/cuiyuanxin/kunpeng/internal/redis"
	"github.com/cuiyuanxin/kunpeng/internal/router"
)

// @title Kunpeng API
// @version 1.0
// @description 基于Gin+GORM+Viper+Zap的Go Web框架脚手架
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// 解析命令行参数
	configPath := flag.String("config", "configs/config.yaml", "配置文件路径")
	flag.Parse()

	// 检查环境变量覆盖
	if envConfigPath := os.Getenv("KUNPENG_CONFIG_PATH"); envConfigPath != "" {
		*configPath = envConfigPath
	}

	// 初始化配置
	cfg, err := config.Init(*configPath)
	if err != nil {
		panic(fmt.Sprintf("Failed to init config from %s: %v", *configPath, err))
	}

	// 初始化日志
	if err := klogger.Init(&cfg.Logging); err != nil {
		panic(fmt.Sprintf("Failed to init logger: %v", err))
	}
	defer klogger.Sync()

	klogger.Info("Starting Kunpeng application...",
		zap.String("version", cfg.App.Version),
		zap.String("environment", cfg.App.Environment),
		zap.String("config_path", *configPath),
	)

	// 初始化数据库
	if err := database.Init(&cfg.Database); err != nil {
		klogger.Fatal("Failed to init database", zap.Error(err))
	}
	defer database.Close()

	// 检查数据库连接
	if err := database.HealthCheck(); err != nil {
		klogger.Fatal("Database health check failed", zap.Error(err))
	}
	klogger.Info("Database connection verified successfully")

	// 初始化Redis
	if err := redis.Init(&cfg.Redis); err != nil {
		klogger.Fatal("Failed to init redis", zap.Error(err))
	}
	defer redis.Close()

	// 创建路由
	r := router.NewRouter(cfg)
	r.Setup()

	// 启动配置文件监控
	config.StartWatching(klogger.Logger, func() {
		// 配置文件变更时的回调
		newCfg := config.GetConfig()
		if newCfg != nil {
			klogger.Info("Configuration updated",
				zap.String("version", newCfg.App.Version),
				zap.String("environment", newCfg.App.Environment),
			)

			// 更新日志配置
			if err := klogger.Init(&newCfg.Logging); err != nil {
				klogger.Error("Failed to reinit logger with new config", zap.Error(err))
			} else {
				klogger.Info("Logger configuration updated successfully")
			}
		}
	})

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:         cfg.Server.GetServerAddr(),
		Handler:      r.GetEngine(),
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}



	// 启动服务器
	go func() {
		klogger.Info("Server starting", zap.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			klogger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// 等待中断信号以优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	klogger.Info("Shutting down server...")

	// 设置5秒的超时时间来关闭服务器
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		klogger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	klogger.Info("Server exited")
}
