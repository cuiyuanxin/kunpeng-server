package config

import "time"

// Config 应用配置结构
type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Log      LogConfig      `mapstructure:"log"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Casbin   CasbinConfig   `mapstructure:"casbin"`
}

// AppConfig 应用基础配置
type AppConfig struct {
	Name        string `mapstructure:"name"`
	Mode        string `mapstructure:"mode"`
	Version     string `mapstructure:"version"`
	Language    string `mapstructure:"language"`
	TraceEnable bool   `mapstructure:"trace_enable"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver          string `mapstructure:"driver"`
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	Database        string `mapstructure:"database"`
	Username        string `mapstructure:"username"`
	Password        string `mapstructure:"password"`
	Charset         string `mapstructure:"charset"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
	ShowSql         bool   `mapstructure:"show_sql"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
	// SQL日志配置
	SqlFilename   string `mapstructure:"sql_filename"`
	SqlMaxSize    int    `mapstructure:"sql_max_size"`
	SqlMaxBackups int    `mapstructure:"sql_max_backups"`
	SqlMaxAge     int    `mapstructure:"sql_max_age"`
	SqlCompress   bool   `mapstructure:"sql_compress"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret               string        `mapstructure:"secret"`
	Issuer               string        `mapstructure:"issuer"`
	ExpireTime           time.Duration `mapstructure:"expire_time"`             // 普通登录token有效期
	RememberMeExpireTime time.Duration `mapstructure:"remember_me_expire_time"` // 记住我token有效期
}

// CasbinConfig Casbin配置
type CasbinConfig struct {
	ModelPath string `mapstructure:"model_path"`
	Enable    bool   `mapstructure:"enable"`
}
