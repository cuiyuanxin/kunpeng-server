package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

var config *Config

// Init 初始化配置
func Init(configPath string) error {
	workDir, err := os.Getwd()
	if err != nil {
		return err
	}

	// 如果指定了配置文件路径，则使用指定的路径
	if configPath != "" {
		// 检查配置文件是否存在
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			return fmt.Errorf("指定的配置文件不存在: %s", configPath)
		}

		// 设置配置文件
		viper.SetConfigFile(configPath)
	} else {
		// 使用默认配置文件路径
		viper.AddConfigPath(filepath.Join(workDir, "configs"))
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	// 读取环境变量
	viper.AutomaticEnv()

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 解析配置到结构体
	if err := viper.Unmarshal(&config); err != nil {
		return fmt.Errorf("解析配置失败: %w", err)
	}

	// 监听配置文件变化
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("配置文件已更改: %s\n", e.Name)
		if err := viper.Unmarshal(&config); err != nil {
			fmt.Printf("重新加载配置失败: %v\n", err)
		}
	})

	return nil
}

// Get 获取配置
func Get() *Config {
	return config
}

// GetAppConfig 获取应用配置
func GetAppConfig() AppConfig {
	return config.App
}

// GetServerConfig 获取服务器配置
func GetServerConfig() ServerConfig {
	return config.Server
}

// GetDatabaseConfig 获取数据库配置
func GetDatabaseConfig() DatabaseConfig {
	return config.Database
}

// GetLogConfig 获取日志配置
func GetLogConfig() LogConfig {
	return config.Log
}

// GetJWTConfig 获取JWT配置
func GetJWTConfig() JWTConfig {
	return config.JWT
}

// GetCasbinConfig 获取Casbin配置
func GetCasbinConfig() CasbinConfig {
	return config.Casbin
}

// IsProduction 是否为生产环境
func IsProduction() bool {
	return config.App.Mode == "production"
}

// IsDevelopment 是否为开发环境
func IsDevelopment() bool {
	return config.App.Mode == "development"
}

// IsTest 是否为测试环境
func IsTest() bool {
	return config.App.Mode == "test"
}
