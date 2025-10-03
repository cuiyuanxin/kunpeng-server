package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	_ "github.com/cuiyuanxin/kunpeng/docs" // 导入swagger文档
	"github.com/cuiyuanxin/kunpeng/internal/app"
)

// 版本信息变量，在构建时通过 -ldflags 注入
var (
	Version   = "dev"             // 版本号
	BuildTime = "unknown"         // 构建时间
	GitCommit = "unknown"         // Git提交哈希
	GoVersion = runtime.Version() // Go版本
)

// @title 鲲鹏后台管理系统API
// @version 1.0
// @description 鲲鹏后台管理系统API文档
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @schemes http https

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	// 定义命令行参数
	var (
		configPath  = flag.String("config", "configs/config.yaml", "配置文件路径 (例如: -config=/path/to/config.yaml)")
		showHelp    = flag.Bool("help", false, "显示帮助信息")
		showVersion = flag.Bool("version", false, "显示版本信息")
	)

	// 解析命令行参数
	flag.Parse()

	// 显示帮助信息
	if *showHelp {
		fmt.Println("鲲鹏后台管理系统")
		fmt.Println()
		fmt.Println("使用方法:")
		fmt.Printf("  %s [选项]\n", os.Args[0])
		fmt.Println()
		fmt.Println("选项:")
		flag.PrintDefaults()
		fmt.Println()
		fmt.Println("示例:")
		fmt.Printf("  %s                              # 使用默认配置文件 (configs/config.yaml)\n", os.Args[0])
		fmt.Printf("  %s -config=/path/to/config.yaml # 使用指定配置文件\n", os.Args[0])
		fmt.Printf("  %s -help                        # 显示此帮助信息\n", os.Args[0])
		fmt.Printf("  %s -version                     # 显示版本信息\n", os.Args[0])
		return
	}

	// 显示版本信息
	if *showVersion {
		fmt.Printf("鲲鹏后台管理系统 %s\n", Version)
		fmt.Printf("构建时间: %s\n", BuildTime)
		fmt.Printf("Git提交: %s\n", GitCommit)
		fmt.Printf("Go版本: %s\n", GoVersion)
		return
	}

	// 创建应用程序
	var application *app.App

	// 如果没有指定配置文件，使用默认配置文件
	if *configPath == "" {
		*configPath = "configs/config.yaml"
	}

	// 将相对路径转换为绝对路径用于显示
	absPath := *configPath
	if absPath != "" {
		if workDir, err := os.Getwd(); err == nil {
			if !filepath.IsAbs(absPath) {
				absPath = filepath.Join(workDir, absPath)
			}
		}
	}
	fmt.Printf("使用配置文件: %s\n", absPath)
	application = app.New(*configPath)

	// 运行应用程序
	application.Run()
}
