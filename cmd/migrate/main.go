package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/cuiyuanxin/kunpeng/internal/config"
	"github.com/cuiyuanxin/kunpeng/internal/database"
	klogger "github.com/cuiyuanxin/kunpeng/internal/logger"
	"github.com/cuiyuanxin/kunpeng/internal/migrate"
	"go.uber.org/zap"
)

func main() {
	// 定义命令行参数
	var (
		configPath = flag.String("config", "configs/config.yaml", "配置文件路径")
		action     = flag.String("action", "migrate", "执行的操作: migrate(迁移), drop(删除表), reset(重置数据库)")
		help       = flag.Bool("help", false, "显示帮助信息")
	)
	flag.Parse()

	// 显示帮助信息
	if *help {
		printHelp()
		return
	}

	// 检查环境变量中的配置路径
	if envConfigPath := os.Getenv("KUNPENG_CONFIG_PATH"); envConfigPath != "" {
		*configPath = envConfigPath
	}

	// 初始化配置
	cfg, err := config.Init(*configPath)
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// 初始化日志
	if err := klogger.Init(&cfg.Logging); err != nil {
		fmt.Printf("Failed to init logger: %v\n", err)
		os.Exit(1)
	}
	defer klogger.Sync()

	// 初始化数据库
	if err := database.Init(&cfg.Database); err != nil {
		klogger.Fatal("Failed to init database", zap.Error(err))
	}
	defer database.Close()

	// 执行相应的操作
	switch *action {
	case "migrate":
		if err := migrate.AutoMigrate(); err != nil {
			klogger.Fatal("Migration failed", zap.Error(err))
		}
		klogger.Info("Migration completed successfully")

	case "drop":
		if err := migrate.DropAllTables(); err != nil {
			klogger.Fatal("Drop tables failed", zap.Error(err))
		}
		klogger.Info("All tables dropped successfully")

	case "reset":
		if err := migrate.ResetDatabase(); err != nil {
			klogger.Fatal("Reset database failed", zap.Error(err))
		}
		klogger.Info("Database reset completed successfully")

	default:
		fmt.Printf("Unknown action: %s\n", *action)
		printHelp()
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println("数据库迁移工具")
	fmt.Println("")
	fmt.Println("用法:")
	fmt.Println("  go run cmd/migrate/main.go [选项]")
	fmt.Println("")
	fmt.Println("选项:")
	fmt.Println("  -config string")
	fmt.Println("        配置文件路径 (默认: configs/config.yaml)")
	fmt.Println("  -action string")
	fmt.Println("        执行的操作 (默认: migrate)")
	fmt.Println("        可选值:")
	fmt.Println("          migrate  - 自动迁移数据库表")
	fmt.Println("          drop     - 删除所有数据库表")
	fmt.Println("          reset    - 重置数据库（删除后重新创建）")
	fmt.Println("  -help")
	fmt.Println("        显示此帮助信息")
	fmt.Println("")
	fmt.Println("环境变量:")
	fmt.Println("  KUNPENG_CONFIG_PATH  配置文件路径（优先级高于-config参数）")
	fmt.Println("")
	fmt.Println("示例:")
	fmt.Println("  # 执行数据库迁移")
	fmt.Println("  go run cmd/migrate/main.go")
	fmt.Println("")
	fmt.Println("  # 使用指定配置文件")
	fmt.Println("  go run cmd/migrate/main.go -config=configs/config.prod.yaml")
	fmt.Println("")
	fmt.Println("  # 重置数据库")
	fmt.Println("  go run cmd/migrate/main.go -action=reset")
	fmt.Println("")
	fmt.Println("  # 删除所有表")
	fmt.Println("  go run cmd/migrate/main.go -action=drop")
}