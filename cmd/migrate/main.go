package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/cuiyuanxin/kunpeng/internal/config"
	"github.com/cuiyuanxin/kunpeng/internal/database"
	klogger "github.com/cuiyuanxin/kunpeng/internal/logger"
	"github.com/cuiyuanxin/kunpeng/internal/migrate"
)

func main() {
	// 定义命令行参数
	var (
		configPath = flag.String("config", "configs/config.yaml", "配置文件路径")
		action     = flag.String("action", "migrate", "执行的操作 (migrate, drop, reset)")
		targetDB   = flag.String("database", "all", "目标数据库 (all, primary, 或指定数据库名称)")
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
	if err := database.InitWithConfig(cfg); err != nil {
		log.Fatal("Failed to init database:", err)
	}
	defer database.Close()

	// 显示可用的数据库
	databases := database.ListDatabases()
	fmt.Printf("Available databases: %v\n", databases)

	// 根据操作类型执行相应的功能
	target := *targetDB
	switch *action {
	case "migrate":
		performMigration(target)
	case "drop":
		performDrop(target)
	case "reset":
		performReset(target)
	default:
		fmt.Printf("未知操作: %s\n", *action)
		os.Exit(1)
	}
}

func performMigration(target string) {
	fmt.Printf("执行数据库迁移，目标: %s\n", target)
	
	switch target {
	case "all":
		if err := database.AutoMigrateAll(); err != nil {
			log.Fatal("Failed to migrate all databases:", err)
		}
		fmt.Println("所有数据库迁移完成")
	case "primary":
		if err := migrate.AutoMigrate(); err != nil {
			log.Fatal("Failed to migrate primary database:", err)
		}
		fmt.Println("主数据库迁移完成")
	default:
		if err := database.AutoMigrateOnDatabase(target); err != nil {
			log.Fatal("Failed to migrate database:", err)
		}
		fmt.Printf("数据库 %s 迁移完成\n", target)
	}
}

func performDrop(target string) {
	fmt.Printf("删除数据库表，目标: %s\n", target)
	
	switch target {
	case "all":
		if err := migrate.DropAllTables(); err != nil {
			log.Fatal("Failed to drop tables on all databases:", err)
		}
		fmt.Println("所有数据库的表已删除")
	case "primary":
		if err := migrate.DropAllTables(); err != nil {
			log.Fatal("Failed to drop tables on primary database:", err)
		}
		fmt.Println("主数据库的表已删除")
	default:
		// 对于指定数据库，暂时使用主数据库的方法
		if err := migrate.DropAllTables(); err != nil {
			log.Fatal("Failed to drop tables on database:", err)
		}
		fmt.Printf("数据库 %s 的表已删除\n", target)
	}
}

func performReset(target string) {
	fmt.Printf("重置数据库，目标: %s\n", target)
	
	switch target {
	case "all":
		if err := migrate.ResetDatabase(); err != nil {
			log.Fatal("Failed to reset all databases:", err)
		}
		fmt.Println("所有数据库已重置")
	case "primary":
		if err := migrate.ResetDatabase(); err != nil {
			log.Fatal("Failed to reset primary database:", err)
		}
		fmt.Println("主数据库已重置")
	default:
		// 对于指定数据库，暂时使用主数据库的方法
		if err := migrate.ResetDatabase(); err != nil {
			log.Fatal("Failed to reset database:", err)
		}
		fmt.Printf("数据库 %s 已重置\n", target)
	}
}

func printHelp() {
	fmt.Println("数据库迁移工具（支持多数据库）")
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
	fmt.Println("  -database string")
	fmt.Println("        目标数据库 (默认: all)")
	fmt.Println("        可选值:")
	fmt.Println("          all      - 所有数据库")
	fmt.Println("          primary  - 主数据库")
	fmt.Println("          <name>   - 指定数据库名称")
	fmt.Println("  -help")
	fmt.Println("        显示此帮助信息")
	fmt.Println("")
	fmt.Println("环境变量:")
	fmt.Println("  KUNPENG_CONFIG_PATH  配置文件路径（优先级高于-config参数）")
	fmt.Println("")
	fmt.Println("示例:")
	fmt.Println("  # 在所有数据库上执行迁移")
	fmt.Println("  go run cmd/migrate/main.go")
	fmt.Println("")
	fmt.Println("  # 在主数据库上执行迁移")
	fmt.Println("  go run cmd/migrate/main.go -database=primary")
	fmt.Println("")
	fmt.Println("  # 在指定数据库上执行迁移")
	fmt.Println("  go run cmd/migrate/main.go -database=user_db")
	fmt.Println("")
	fmt.Println("  # 使用指定配置文件")
	fmt.Println("  go run cmd/migrate/main.go -config=configs/config.dev.yaml")
	fmt.Println("")
	fmt.Println("  # 重置所有数据库")
	fmt.Println("  go run cmd/migrate/main.go -action=reset -database=all")
	fmt.Println("")
	fmt.Println("  # 删除指定数据库的所有表")
	fmt.Println("  go run cmd/migrate/main.go -action=drop -database=user_db")
}