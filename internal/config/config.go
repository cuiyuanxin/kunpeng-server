package config

import (
	"fmt"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Config 应用配置结构
type Config struct {
	App      App      `mapstructure:"app"`
	Server   Server   `mapstructure:"server"`
	Database Database `mapstructure:"database"`
	Redis    Redis    `mapstructure:"redis"`
	Logging  Logging  `mapstructure:"logging"`
	JWT      JWT      `mapstructure:"jwt"`
}

// 全局viper实例
var globalViper *viper.Viper

// App 应用配置
type App struct {
	Name        string `mapstructure:"name"`
	Version     string `mapstructure:"version"`
	Environment string `mapstructure:"environment"`
	Debug       bool   `mapstructure:"debug"`
}

// Server 服务器配置
type Server struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}

// Database 数据库配置
type Database struct {
	Driver          string        `mapstructure:"driver"`
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	Username        string        `mapstructure:"username"`
	Password        string        `mapstructure:"password"`
	Database        string        `mapstructure:"database"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

// Redis 配置
type Redis struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Password     string `mapstructure:"password"`
	Database     int    `mapstructure:"database"`
	PoolSize     int    `mapstructure:"pool_size"`
	MinIdleConns int    `mapstructure:"min_idle_conns"`
}

// Logging 日志配置
type Logging struct {
	Level         string `mapstructure:"level"`
	Format        string `mapstructure:"format"`
	Output        string `mapstructure:"output"`
	FilePath      string `mapstructure:"file_path"`
	MaxSize       int    `mapstructure:"max_size"`
	MaxBackups    int    `mapstructure:"max_backups"`
	MaxAge        int    `mapstructure:"max_age"`
	Compress      bool   `mapstructure:"compress"`
	// 分级别日志文件配置
	SeparateFiles bool   `mapstructure:"separate_files"` // 是否启用分级别日志文件
	LogDir        string `mapstructure:"log_dir"`        // 日志目录
}

// JWT 配置
type JWT struct {
	Secret     string        `mapstructure:"secret"`
	ExpireTime time.Duration `mapstructure:"expire_time"`
	Issuer     string        `mapstructure:"issuer"`
}

// Init 初始化配置
func Init(configPath string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	// 不自动读取环境变量，避免配置同步到环境变量
	// v.SetEnvPrefix("KUNPENG")
	// v.AutomaticEnv()

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// 保存全局viper实例
	globalViper = v
	return &config, nil
}

// StartWatching 开始监控配置文件变更
func StartWatching(logger *zap.Logger, callback func()) {
	if globalViper == nil {
		if logger != nil {
			logger.Error("Global viper not initialized")
		}
		return
	}

	globalViper.OnConfigChange(func(e fsnotify.Event) {
		if logger != nil {
			logger.Info("Config file changed", zap.String("file", e.Name))
		}
		if callback != nil {
			callback()
		}
	})

	globalViper.WatchConfig()
	if logger != nil {
		logger.Info("Config file watching started")
	}
}

// GetConfig 获取当前配置
func GetConfig() *Config {
	if globalViper == nil {
		return nil
	}

	var config Config
	if err := globalViper.Unmarshal(&config); err != nil {
		return nil
	}
	return &config
}

// Load 加载配置文件（兼容性保留）
func Load(configPath string) (*Config, error) {
	return Init(configPath)
}

// GetDSN 获取数据库连接字符串
func (d *Database) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		d.Username, d.Password, d.Host, d.Port, d.Database)
}

// GetRedisAddr 获取Redis地址
func (r *Redis) GetRedisAddr() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

// GetServerAddr 获取服务器地址
func (s *Server) GetServerAddr() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}