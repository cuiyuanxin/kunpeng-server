package app

import (
	"fmt"
	"log"
)

// App 应用程序结构体
type App struct {
	Name    string
	Version string
}

// New 创建新的应用程序实例
func New(name, version string) *App {
	return &App{
		Name:    name,
		Version: version,
	}
}

// Start 启动应用程序
func (a *App) Start() error {
	log.Printf("Starting %s v%s", a.Name, a.Version)
	fmt.Printf("Welcome to %s!\n", a.Name)
	return nil
}

// Stop 停止应用程序
func (a *App) Stop() error {
	log.Printf("Stopping %s", a.Name)
	return nil
}