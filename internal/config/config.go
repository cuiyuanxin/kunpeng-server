package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Config 应用配置结构
type Config struct {
	App       App                 `mapstructure:"app"`
	Server    Server              `mapstructure:"server"`
	Database  Database            `mapstructure:"database"`  // 主数据库配置（向后兼容）
	Databases map[string]Database `mapstructure:"databases"` // 多数据库配置
	Redis     Redis               `mapstructure:"redis"`
	Logging   Logging             `mapstructure:"logging"`
	JWT       JWT                 `mapstructure:"jwt"`
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
	// 基础连接配置
	Driver          string        `mapstructure:"driver"`           // 数据库驱动: mysql, postgres, sqlite, sqlserver, clickhouse
	Host            string        `mapstructure:"host"`             // 主机地址
	Port            int           `mapstructure:"port"`             // 端口
	Username        string        `mapstructure:"username"`         // 用户名
	Password        string        `mapstructure:"password"`         // 密码
	Database        string        `mapstructure:"database"`         // 数据库名
	Schema          string        `mapstructure:"schema"`           // 模式名（PostgreSQL等）
	SSLMode         string        `mapstructure:"ssl_mode"`         // SSL模式（PostgreSQL等）
	Timezone        string        `mapstructure:"timezone"`         // 时区
	Charset         string        `mapstructure:"charset"`          // 字符集
	
	// SQLite特有配置
	FilePath        string        `mapstructure:"file_path"`        // SQLite文件路径
	
	// 连接池配置
	MaxOpenConns    int           `mapstructure:"max_open_conns"`   // 最大打开连接数
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`   // 最大空闲连接数
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"` // 连接最大生命周期
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time"` // 连接最大空闲时间
	
	// 高级配置
	DSN             string        `mapstructure:"dsn"`              // 自定义DSN（优先级最高）
	Options         map[string]interface{} `mapstructure:"options"`   // 额外选项
	
	// gRPC支持配置
	GRPCEnabled     bool          `mapstructure:"grpc_enabled"`     // 是否启用gRPC支持
	GRPCPoolSize    int           `mapstructure:"grpc_pool_size"`   // gRPC连接池大小
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
	// 环境自适应配置
	AutoMode      bool   `mapstructure:"auto_mode"`      // 是否启用环境自适应模式
	ForceConsole  *bool  `mapstructure:"force_console"`  // 强制使用控制台输出（可选配置）
	ForceFile     *bool  `mapstructure:"force_file"`     // 强制使用文件输出（可选配置）
	// GORM日志配置
	Gorm          GormLogging `mapstructure:"gorm"`
	// 扩展钩子配置
	Hooks         LoggingHooks `mapstructure:"hooks"`
}

// GormLogging GORM日志配置
type GormLogging struct {
	Enabled       bool   `mapstructure:"enabled"`        // 是否启用GORM日志
	Level         string `mapstructure:"level"`          // GORM日志级别 (silent, error, warn, info)
	SlowThreshold string `mapstructure:"slow_threshold"` // 慢查询阈值
	SQLFile       string `mapstructure:"sql_file"`       // SQL日志文件路径
	ErrorFile     string `mapstructure:"error_file"`     // GORM错误日志文件路径
	AutoMode      bool   `mapstructure:"auto_mode"`      // 是否启用环境自适应模式
	ForceConsole  *bool  `mapstructure:"force_console"`  // 强制使用控制台输出
	ForceFile     *bool  `mapstructure:"force_file"`     // 强制使用文件输出
}

// LoggingHooks 日志钩子配置
type LoggingHooks struct {
	Tracing LoggingHook `mapstructure:"tracing"` // 链路追踪钩子
	GRPC    LoggingHook `mapstructure:"grpc"`    // gRPC日志钩子
	Custom  LoggingHook `mapstructure:"custom"`  // 自定义钩子
}

// LoggingHook 单个钩子配置
type LoggingHook struct {
	Enabled   bool   `mapstructure:"enabled"`    // 是否启用
	Level     string `mapstructure:"level"`     // 日志级别
	Format    string `mapstructure:"format"`    // 日志格式
	Output    string `mapstructure:"output"`    // 输出方式
	FilePath  string `mapstructure:"file_path"` // 文件路径
	AutoMode  bool   `mapstructure:"auto_mode"` // 环境自适应
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
	// 如果设置了自定义DSN，直接返回
	if d.DSN != "" {
		return d.DSN
	}
	
	switch strings.ToLower(d.Driver) {
	case "mysql":
		return d.getMySQLDSN()
	case "postgres", "postgresql":
		return d.getPostgresDSN()
	case "sqlite", "sqlite3":
		return d.getSQLiteDSN()
	case "sqlserver", "mssql":
		return d.getSQLServerDSN()
	case "clickhouse":
		return d.getClickHouseDSN()
	default:
		// 默认使用MySQL格式
		return d.getMySQLDSN()
	}
}

// getMySQLDSN 获取MySQL DSN
func (d *Database) getMySQLDSN() string {
	charset := d.Charset
	if charset == "" {
		charset = "utf8mb4"
	}
	timezone := d.Timezone
	if timezone == "" {
		timezone = "Local"
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=%s",
		d.Username, d.Password, d.Host, d.Port, d.Database, charset, timezone)
}

// getPostgresDSN 获取PostgreSQL DSN
func (d *Database) getPostgresDSN() string {
	sslMode := d.SSLMode
	if sslMode == "" {
		sslMode = "disable"
	}
	timezone := d.Timezone
	if timezone == "" {
		timezone = "UTC"
	}
	
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		d.Host, d.Port, d.Username, d.Password, d.Database, sslMode, timezone)
	
	if d.Schema != "" {
		dsn += fmt.Sprintf(" search_path=%s", d.Schema)
	}
	
	return dsn
}

// getSQLiteDSN 获取SQLite DSN
func (d *Database) getSQLiteDSN() string {
	filePath := d.FilePath
	if filePath == "" {
		filePath = d.Database
	}
	if filePath == "" {
		filePath = "./data.db"
	}
	return filePath
}

// getSQLServerDSN 获取SQL Server DSN
func (d *Database) getSQLServerDSN() string {
	return fmt.Sprintf("server=%s;port=%d;user id=%s;password=%s;database=%s",
		d.Host, d.Port, d.Username, d.Password, d.Database)
}

// getClickHouseDSN 获取ClickHouse DSN
func (d *Database) getClickHouseDSN() string {
	return fmt.Sprintf("tcp://%s:%d?username=%s&password=%s&database=%s",
		d.Host, d.Port, d.Username, d.Password, d.Database)
}

// GetRedisAddr 获取Redis地址
func (r *Redis) GetRedisAddr() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

// GetServerAddr 获取服务器地址
func (s *Server) GetServerAddr() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}