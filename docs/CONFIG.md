# 配置管理系统文档

## 概述

本项目提供了完整的配置管理系统，基于 Viper 构建，支持多种配置格式、环境变量、配置热重载等功能。配置系统采用结构化设计，支持应用配置、服务器配置、数据库配置、Redis配置、日志配置、JWT配置等模块化管理。

## 系统架构

### 配置系统组件

```
Config System
├── Core Config              # 核心配置管理
│   ├── Viper Integration    # Viper 集成
│   ├── File Watching        # 文件监控
│   └── Hot Reload           # 热重载
├── App Config               # 应用配置
├── Server Config            # 服务器配置
├── Database Config          # 数据库配置
│   ├── Single Database      # 单数据库配置
│   └── Multi Database       # 多数据库配置
├── Redis Config             # Redis 配置
├── Logging Config           # 日志配置
│   ├── Basic Logging        # 基础日志
│   ├── GORM Logging         # GORM 日志
│   └── Logging Hooks        # 日志钩子
└── JWT Config               # JWT 配置
```

## 配置结构定义

### 主配置结构

```go
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
```

### 应用配置

```go
// App 应用基础配置
type App struct {
    Name        string `mapstructure:"name"`        // 应用名称
    Version     string `mapstructure:"version"`     // 应用版本
    Environment string `mapstructure:"environment"` // 运行环境 (development, testing, production)
    Debug       bool   `mapstructure:"debug"`       // 调试模式
}
```

**配置示例**:
```yaml
app:
  name: "kunpeng-server"
  version: "1.0.0"
  environment: "development"
  debug: true
```

### 服务器配置

```go
// Server HTTP服务器配置
type Server struct {
    Host         string        `mapstructure:"host"`          // 监听地址
    Port         int           `mapstructure:"port"`          // 监听端口
    ReadTimeout  time.Duration `mapstructure:"read_timeout"`  // 读取超时
    WriteTimeout time.Duration `mapstructure:"write_timeout"` // 写入超时
    IdleTimeout  time.Duration `mapstructure:"idle_timeout"`  // 空闲超时
}

// GetServerAddr 获取服务器地址
func (s *Server) GetServerAddr() string {
    return fmt.Sprintf("%s:%d", s.Host, s.Port)
}
```

**配置示例**:
```yaml
server:
  host: "0.0.0.0"
  port: 8080
  read_timeout: "30s"
  write_timeout: "30s"
  idle_timeout: "60s"
```

### 数据库配置

```go
// Database 数据库配置
type Database struct {
    // 基础连接配置
    Driver          string        `mapstructure:"driver"`           // 数据库驱动
    Host            string        `mapstructure:"host"`             // 主机地址
    Port            int           `mapstructure:"port"`             // 端口
    Username        string        `mapstructure:"username"`         // 用户名
    Password        string        `mapstructure:"password"`         // 密码
    Database        string        `mapstructure:"database"`         // 数据库名
    Schema          string        `mapstructure:"schema"`           // 模式名（PostgreSQL等）
    SSLMode         string        `mapstructure:"ssl_mode"`         // SSL模式
    Timezone        string        `mapstructure:"timezone"`         // 时区
    Charset         string        `mapstructure:"charset"`          // 字符集
    
    // SQLite特有配置
    FilePath        string        `mapstructure:"file_path"`        // SQLite文件路径
    
    // 连接池配置
    MaxOpenConns    int           `mapstructure:"max_open_conns"`   // 最大打开连接数
    MaxIdleConns    int           `mapstructure:"max_idle_conns"`   // 最大空闲连接数
    ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"` // 连接最大生命周期
    ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time"` // 连接最大空闲时间
    
    // 自定义DSN（可选）
    DSN             string        `mapstructure:"dsn"`              // 自定义数据源名称
    
    // gRPC支持配置
    GRPCEnabled     bool          `mapstructure:"grpc_enabled"`     // 是否启用gRPC支持
    GRPCPoolSize    int           `mapstructure:"grpc_pool_size"`   // gRPC连接池大小
}
```

**DSN 生成方法**:
```go
// GetDSN 获取数据源名称
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
    case "sqlserver":
        return d.getSQLServerDSN()
    case "clickhouse":
        return d.getClickHouseDSN()
    default:
        return ""
    }
}

// getMySQLDSN 生成MySQL DSN
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

// getPostgresDSN 生成PostgreSQL DSN
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

// getSQLiteDSN 生成SQLite DSN
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

// getSQLServerDSN 生成SQL Server DSN
func (d *Database) getSQLServerDSN() string {
    return fmt.Sprintf("server=%s;port=%d;user id=%s;password=%s;database=%s",
        d.Host, d.Port, d.Username, d.Password, d.Database)
}
```

**配置示例**:
```yaml
# 单数据库配置（向后兼容）
database:
  driver: "mysql"
  host: "localhost"
  port: 3306
  username: "root"
  password: "password"
  database: "kunpeng"
  charset: "utf8mb4"
  timezone: "Asia/Shanghai"
  max_open_conns: 100
  max_idle_conns: 10
  conn_max_lifetime: "1h"
  conn_max_idle_time: "10m"

# 多数据库配置
databases:
  main:
    driver: "mysql"
    host: "localhost"
    port: 3306
    username: "root"
    password: "password"
    database: "kunpeng_main"
    grpc_enabled: true
    grpc_pool_size: 5
  
  analytics:
    driver: "clickhouse"
    host: "localhost"
    port: 9000
    username: "default"
    password: ""
    database: "analytics"
  
  cache:
    driver: "sqlite"
    file_path: "./data/cache.db"
```

### Redis 配置

```go
// Redis Redis配置
type Redis struct {
    Host         string `mapstructure:"host"`           // Redis主机
    Port         int    `mapstructure:"port"`           // Redis端口
    Password     string `mapstructure:"password"`       // Redis密码
    Database     int    `mapstructure:"database"`       // Redis数据库编号
    PoolSize     int    `mapstructure:"pool_size"`      // 连接池大小
    MinIdleConns int    `mapstructure:"min_idle_conns"` // 最小空闲连接数
}
```

**配置示例**:
```yaml
redis:
  host: "localhost"
  port: 6379
  password: ""
  database: 0
  pool_size: 10
  min_idle_conns: 5
```

### 日志配置

```go
// Logging 日志配置
type Logging struct {
    Level         string `mapstructure:"level"`          // 日志级别
    Format        string `mapstructure:"format"`         // 日志格式
    Output        string `mapstructure:"output"`         // 输出方式
    FilePath      string `mapstructure:"file_path"`      // 文件路径
    MaxSize       int    `mapstructure:"max_size"`       // 最大文件大小(MB)
    MaxBackups    int    `mapstructure:"max_backups"`    // 最大备份数
    MaxAge        int    `mapstructure:"max_age"`        // 最大保存天数
    Compress      bool   `mapstructure:"compress"`       // 是否压缩
    
    // 分级别日志文件配置
    SeparateFiles bool   `mapstructure:"separate_files"` // 是否启用分级别日志文件
    LogDir        string `mapstructure:"log_dir"`        // 日志目录
    
    // 环境自适应配置
    AutoMode      bool   `mapstructure:"auto_mode"`      // 是否启用环境自适应模式
    ForceConsole  *bool  `mapstructure:"force_console"`  // 强制使用控制台输出
    ForceFile     *bool  `mapstructure:"force_file"`     // 强制使用文件输出
    
    // GORM日志配置
    GormLogging   GormLogging   `mapstructure:"gorm"`    // GORM日志配置
    
    // 扩展钩子配置
    Hooks         LoggingHooks  `mapstructure:"hooks"`   // 日志钩子配置
}

// GormLogging GORM日志配置
type GormLogging struct {
    Enabled       bool   `mapstructure:"enabled"`        // 是否启用GORM日志
    Level         string `mapstructure:"level"`          // GORM日志级别
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

// LoggingHook 日志钩子配置
type LoggingHook struct {
    Enabled   bool   `mapstructure:"enabled"`    // 是否启用
    Level     string `mapstructure:"level"`     // 日志级别
    Format    string `mapstructure:"format"`    // 日志格式
    Output    string `mapstructure:"output"`    // 输出方式
    FilePath  string `mapstructure:"file_path"` // 文件路径
    AutoMode  bool   `mapstructure:"auto_mode"` // 环境自适应
}
```

**配置示例**:
```yaml
logging:
  level: "info"
  format: "json"
  output: "both"  # console, file, both
  file_path: "./logs/app.log"
  max_size: 100
  max_backups: 10
  max_age: 30
  compress: true
  separate_files: true
  log_dir: "./logs"
  auto_mode: true
  
  # GORM日志配置
  gorm:
    enabled: true
    level: "info"
    slow_threshold: "200ms"
    sql_file: "./logs/sql.log"
    error_file: "./logs/gorm_error.log"
    auto_mode: true
  
  # 日志钩子配置
  hooks:
    tracing:
      enabled: false
      level: "info"
      format: "json"
      output: "file"
      file_path: "./logs/tracing.log"
    grpc:
      enabled: false
      level: "info"
      format: "json"
      output: "file"
      file_path: "./logs/grpc.log"
    custom:
      enabled: false
      level: "info"
      format: "json"
      output: "console"
```

### JWT 配置

```go
// JWT JWT认证配置
type JWT struct {
    Secret     string        `mapstructure:"secret"`      // JWT密钥
    ExpireTime time.Duration `mapstructure:"expire_time"` // 过期时间
    Issuer     string        `mapstructure:"issuer"`      // 签发者
}
```

**配置示例**:
```yaml
jwt:
  secret: "your-secret-key-here"
  expire_time: "24h"
  issuer: "kunpeng-server"
```

## 配置管理功能

### 配置初始化

```go
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

    // 保存全局viper实例用于热重载
    globalViper = v
    
    return &config, nil
}

// Load 加载配置（别名）
func Load(configPath string) (*Config, error) {
    return Init(configPath)
}
```

### 配置热重载

```go
var globalViper *viper.Viper

// StartWatching 启动配置文件监控
func StartWatching(logger *zap.Logger, callback func()) {
    if globalViper == nil {
        if logger != nil {
            logger.Error("Global viper not initialized")
        }
        return
    }

    globalViper.WatchConfig()
    globalViper.OnConfigChange(func(e fsnotify.Event) {
        if logger != nil {
            logger.Info("Config file changed", zap.String("file", e.Name))
        }
        
        // 执行回调函数
        if callback != nil {
            callback()
        }
    })
    
    if logger != nil {
        logger.Info("Started watching config file")
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
```

## 使用示例

### 基础配置加载

```go
package main

import (
    "log"
    "github.com/your-project/internal/config"
)

func main() {
    // 加载配置
    cfg, err := config.Load("./configs/config.yaml")
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }
    
    // 使用配置
    fmt.Printf("App Name: %s\n", cfg.App.Name)
    fmt.Printf("Server Address: %s\n", cfg.Server.GetServerAddr())
    fmt.Printf("Database DSN: %s\n", cfg.Database.GetDSN())
}
```

### 配置热重载示例

```go
package main

import (
    "log"
    "github.com/your-project/internal/config"
    "github.com/your-project/pkg/klogger"
)

func main() {
    // 初始化日志
    logger, _ := zap.NewDevelopment()
    defer logger.Sync()
    
    // 加载配置
    cfg, err := config.Load("./configs/config.yaml")
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }
    
    // 启动配置监控
    config.StartWatching(logger, func() {
        // 配置文件变更时的回调
        newCfg := config.GetConfig()
        if newCfg != nil {
            logger.Info("Config reloaded", 
                zap.String("app_name", newCfg.App.Name),
                zap.String("environment", newCfg.App.Environment),
            )
            
            // 这里可以重新初始化相关组件
            // 例如：重新初始化日志系统、数据库连接等
            reinitializeComponents(newCfg)
        }
    })
    
    // 应用启动逻辑
    startApplication(cfg)
}

func reinitializeComponents(cfg *config.Config) {
    // 重新初始化日志系统
    klogger.Init(&cfg.Logging)
    
    // 重新初始化数据库连接
    // database.Reinit(&cfg.Database)
    
    // 重新初始化Redis连接
    // redis.Reinit(&cfg.Redis)
}
```

### 多数据库配置示例

```go
func initializeDatabases(cfg *config.Config) error {
    // 初始化主数据库（向后兼容）
    if cfg.Database.Driver != "" {
        mainDB, err := initDatabase(&cfg.Database)
        if err != nil {
            return fmt.Errorf("failed to init main database: %w", err)
        }
        database.RegisterDatabase("main", mainDB)
    }
    
    // 初始化多数据库
    for name, dbConfig := range cfg.Databases {
        db, err := initDatabase(&dbConfig)
        if err != nil {
            return fmt.Errorf("failed to init database %s: %w", name, err)
        }
        database.RegisterDatabase(name, db)
    }
    
    return nil
}

func initDatabase(dbConfig *config.Database) (*gorm.DB, error) {
    dsn := dbConfig.GetDSN()
    
    var db *gorm.DB
    var err error
    
    switch strings.ToLower(dbConfig.Driver) {
    case "mysql":
        db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
    case "postgres", "postgresql":
        db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    case "sqlite", "sqlite3":
        db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
    default:
        return nil, fmt.Errorf("unsupported database driver: %s", dbConfig.Driver)
    }
    
    if err != nil {
        return nil, err
    }
    
    // 配置连接池
    sqlDB, err := db.DB()
    if err != nil {
        return nil, err
    }
    
    sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConns)
    sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConns)
    sqlDB.SetConnMaxLifetime(dbConfig.ConnMaxLifetime)
    sqlDB.SetConnMaxIdleTime(dbConfig.ConnMaxIdleTime)
    
    return db, nil
}
```

### 环境特定配置

```go
// 根据环境加载不同配置文件
func LoadConfigByEnvironment() (*config.Config, error) {
    env := os.Getenv("GO_ENV")
    if env == "" {
        env = "development"
    }
    
    var configFile string
    switch env {
    case "production":
        configFile = "./configs/config.prod.yaml"
    case "testing":
        configFile = "./configs/config.test.yaml"
    default:
        configFile = "./configs/config.dev.yaml"
    }
    
    return config.Load(configFile)
}
```

## 配置文件示例

### 开发环境配置

```yaml
# config.dev.yaml
app:
  name: "kunpeng-server"
  version: "1.0.0"
  environment: "development"
  debug: true

server:
  host: "0.0.0.0"
  port: 8080
  read_timeout: "30s"
  write_timeout: "30s"
  idle_timeout: "60s"

database:
  driver: "mysql"
  host: "localhost"
  port: 3306
  username: "root"
  password: "password"
  database: "kunpeng_dev"
  charset: "utf8mb4"
  timezone: "Asia/Shanghai"
  max_open_conns: 50
  max_idle_conns: 10
  conn_max_lifetime: "1h"
  conn_max_idle_time: "10m"

redis:
  host: "localhost"
  port: 6379
  password: ""
  database: 0
  pool_size: 10
  min_idle_conns: 5

logging:
  level: "debug"
  format: "console"
  output: "console"
  auto_mode: true
  separate_files: false
  
  gorm:
    enabled: true
    level: "info"
    slow_threshold: "100ms"
    auto_mode: true

jwt:
  secret: "dev-secret-key"
  expire_time: "24h"
  issuer: "kunpeng-server"
```

### 生产环境配置

```yaml
# config.prod.yaml
app:
  name: "kunpeng-server"
  version: "1.0.0"
  environment: "production"
  debug: false

server:
  host: "0.0.0.0"
  port: 8080
  read_timeout: "30s"
  write_timeout: "30s"
  idle_timeout: "120s"

# 多数据库配置
databases:
  main:
    driver: "mysql"
    host: "${DB_HOST}"
    port: 3306
    username: "${DB_USER}"
    password: "${DB_PASSWORD}"
    database: "kunpeng_prod"
    charset: "utf8mb4"
    timezone: "UTC"
    max_open_conns: 100
    max_idle_conns: 20
    conn_max_lifetime: "1h"
    conn_max_idle_time: "10m"
    grpc_enabled: true
    grpc_pool_size: 10
  
  analytics:
    driver: "clickhouse"
    host: "${CLICKHOUSE_HOST}"
    port: 9000
    username: "${CLICKHOUSE_USER}"
    password: "${CLICKHOUSE_PASSWORD}"
    database: "analytics"
    max_open_conns: 50
    max_idle_conns: 10
  
  cache:
    driver: "redis"
    host: "${REDIS_HOST}"
    port: 6379
    password: "${REDIS_PASSWORD}"
    database: 1

redis:
  host: "${REDIS_HOST}"
  port: 6379
  password: "${REDIS_PASSWORD}"
  database: 0
  pool_size: 20
  min_idle_conns: 10

logging:
  level: "info"
  format: "json"
  output: "file"
  file_path: "/var/log/kunpeng/app.log"
  max_size: 100
  max_backups: 30
  max_age: 90
  compress: true
  separate_files: true
  log_dir: "/var/log/kunpeng"
  
  gorm:
    enabled: true
    level: "warn"
    slow_threshold: "500ms"
    sql_file: "/var/log/kunpeng/sql.log"
    error_file: "/var/log/kunpeng/gorm_error.log"
    force_file: true

jwt:
  secret: "${JWT_SECRET}"
  expire_time: "2h"
  issuer: "kunpeng-server"
```

### 测试环境配置

```yaml
# config.test.yaml
app:
  name: "kunpeng-server"
  version: "1.0.0"
  environment: "testing"
  debug: true

server:
  host: "127.0.0.1"
  port: 8081
  read_timeout: "10s"
  write_timeout: "10s"
  idle_timeout: "30s"

database:
  driver: "sqlite"
  file_path: ":memory:"
  max_open_conns: 10
  max_idle_conns: 5

redis:
  host: "localhost"
  port: 6379
  password: ""
  database: 15  # 使用测试专用数据库
  pool_size: 5
  min_idle_conns: 2

logging:
  level: "debug"
  format: "console"
  output: "console"
  
  gorm:
    enabled: false  # 测试环境关闭GORM日志

jwt:
  secret: "test-secret-key"
  expire_time: "1h"
  issuer: "kunpeng-server-test"
```

## 最佳实践

### 1. 环境变量支持

```go
// 支持环境变量替换
func expandEnvVars(cfg *Config) {
    // 数据库密码
    if strings.HasPrefix(cfg.Database.Password, "${")
        && strings.HasSuffix(cfg.Database.Password, "}") {
        envVar := cfg.Database.Password[2 : len(cfg.Database.Password)-1]
        cfg.Database.Password = os.Getenv(envVar)
    }
    
    // JWT密钥
    if strings.HasPrefix(cfg.JWT.Secret, "${")
        && strings.HasSuffix(cfg.JWT.Secret, "}") {
        envVar := cfg.JWT.Secret[2 : len(cfg.JWT.Secret)-1]
        cfg.JWT.Secret = os.Getenv(envVar)
    }
}
```

### 2. 配置验证

```go
// ValidateConfig 验证配置
func ValidateConfig(cfg *Config) error {
    // 验证应用配置
    if cfg.App.Name == "" {
        return fmt.Errorf("app name is required")
    }
    
    // 验证服务器配置
    if cfg.Server.Port <= 0 || cfg.Server.Port > 65535 {
        return fmt.Errorf("invalid server port: %d", cfg.Server.Port)
    }
    
    // 验证数据库配置
    if cfg.Database.Driver != "" {
        if err := validateDatabaseConfig(&cfg.Database); err != nil {
            return fmt.Errorf("invalid database config: %w", err)
        }
    }
    
    // 验证多数据库配置
    for name, dbConfig := range cfg.Databases {
        if err := validateDatabaseConfig(&dbConfig); err != nil {
            return fmt.Errorf("invalid database config for %s: %w", name, err)
        }
    }
    
    // 验证JWT配置
    if cfg.JWT.Secret == "" {
        return fmt.Errorf("JWT secret is required")
    }
    
    return nil
}

func validateDatabaseConfig(db *Database) error {
    supportedDrivers := []string{"mysql", "postgres", "postgresql", "sqlite", "sqlite3", "sqlserver", "clickhouse"}
    
    found := false
    for _, driver := range supportedDrivers {
        if strings.ToLower(db.Driver) == driver {
            found = true
            break
        }
    }
    
    if !found {
        return fmt.Errorf("unsupported database driver: %s", db.Driver)
    }
    
    if db.Driver != "sqlite" && db.Driver != "sqlite3" {
        if db.Host == "" {
            return fmt.Errorf("database host is required")
        }
        if db.Username == "" {
            return fmt.Errorf("database username is required")
        }
    }
    
    return nil
}
```

### 3. 配置加密

```go
// 敏感配置加密存储
func EncryptSensitiveConfig(cfg *Config, key string) error {
    // 加密数据库密码
    if cfg.Database.Password != "" {
        encrypted, err := utils.AESEncryptString(cfg.Database.Password, key)
        if err != nil {
            return err
        }
        cfg.Database.Password = encrypted
    }
    
    // 加密JWT密钥
    if cfg.JWT.Secret != "" {
        encrypted, err := utils.AESEncryptString(cfg.JWT.Secret, key)
        if err != nil {
            return err
        }
        cfg.JWT.Secret = encrypted
    }
    
    return nil
}

// 解密敏感配置
func DecryptSensitiveConfig(cfg *Config, key string) error {
    // 解密数据库密码
    if cfg.Database.Password != "" {
        decrypted, err := utils.AESDecryptString(cfg.Database.Password, key)
        if err != nil {
            return err
        }
        cfg.Database.Password = decrypted
    }
    
    // 解密JWT密钥
    if cfg.JWT.Secret != "" {
        decrypted, err := utils.AESDecryptString(cfg.JWT.Secret, key)
        if err != nil {
            return err
        }
        cfg.JWT.Secret = decrypted
    }
    
    return nil
}
```

### 4. 配置缓存

```go
var (
    configCache     *Config
    configCacheMux  sync.RWMutex
    configCacheTime time.Time
    cacheTTL        = 5 * time.Minute
)

// GetCachedConfig 获取缓存的配置
func GetCachedConfig() *Config {
    configCacheMux.RLock()
    defer configCacheMux.RUnlock()
    
    if configCache != nil && time.Since(configCacheTime) < cacheTTL {
        return configCache
    }
    
    return nil
}

// SetCachedConfig 设置缓存的配置
func SetCachedConfig(cfg *Config) {
    configCacheMux.Lock()
    defer configCacheMux.Unlock()
    
    configCache = cfg
    configCacheTime = time.Now()
}
```

## 测试

### 配置测试示例

```go
package config

import (
    "os"
    "testing"
    "time"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
    // 创建临时配置文件
    configContent := `
app:
  name: "test-app"
  version: "1.0.0"
  environment: "testing"
  debug: true

server:
  host: "localhost"
  port: 8080
  read_timeout: "30s"
  write_timeout: "30s"
  idle_timeout: "60s"

database:
  driver: "sqlite"
  file_path: ":memory:"

redis:
  host: "localhost"
  port: 6379
  database: 0

logging:
  level: "debug"
  format: "console"
  output: "console"

jwt:
  secret: "test-secret"
  expire_time: "24h"
  issuer: "test-issuer"
`
    
    tmpFile, err := os.CreateTemp("", "config-*.yaml")
    require.NoError(t, err)
    defer os.Remove(tmpFile.Name())
    
    _, err = tmpFile.WriteString(configContent)
    require.NoError(t, err)
    tmpFile.Close()
    
    // 加载配置
    cfg, err := Load(tmpFile.Name())
    require.NoError(t, err)
    
    // 验证配置
    assert.Equal(t, "test-app", cfg.App.Name)
    assert.Equal(t, "1.0.0", cfg.App.Version)
    assert.Equal(t, "testing", cfg.App.Environment)
    assert.True(t, cfg.App.Debug)
    
    assert.Equal(t, "localhost", cfg.Server.Host)
    assert.Equal(t, 8080, cfg.Server.Port)
    assert.Equal(t, 30*time.Second, cfg.Server.ReadTimeout)
    
    assert.Equal(t, "sqlite", cfg.Database.Driver)
    assert.Equal(t, ":memory:", cfg.Database.FilePath)
    
    assert.Equal(t, "test-secret", cfg.JWT.Secret)
    assert.Equal(t, 24*time.Hour, cfg.JWT.ExpireTime)
}

func TestDatabaseDSN(t *testing.T) {
    tests := []struct {
        name     string
        database Database
        expected string
    }{
        {
            name: "MySQL DSN",
            database: Database{
                Driver:   "mysql",
                Host:     "localhost",
                Port:     3306,
                Username: "root",
                Password: "password",
                Database: "testdb",
                Charset:  "utf8mb4",
                Timezone: "Local",
            },
            expected: "root:password@tcp(localhost:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local",
        },
        {
            name: "SQLite DSN",
            database: Database{
                Driver:   "sqlite",
                FilePath: "./test.db",
            },
            expected: "./test.db",
        },
        {
            name: "Custom DSN",
            database: Database{
                Driver: "mysql",
                DSN:    "custom-dsn-string",
            },
            expected: "custom-dsn-string",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            dsn := tt.database.GetDSN()
            assert.Equal(t, tt.expected, dsn)
        })
    }
}

func TestConfigValidation(t *testing.T) {
    tests := []struct {
        name      string
        config    Config
        expectErr bool
    }{
        {
            name: "Valid config",
            config: Config{
                App: App{
                    Name: "test-app",
                },
                Server: Server{
                    Port: 8080,
                },
                Database: Database{
                    Driver: "mysql",
                    Host:   "localhost",
                    Username: "root",
                },
                JWT: JWT{
                    Secret: "test-secret",
                },
            },
            expectErr: false,
        },
        {
            name: "Missing app name",
            config: Config{
                App: App{},
            },
            expectErr: true,
        },
        {
            name: "Invalid port",
            config: Config{
                App: App{
                    Name: "test-app",
                },
                Server: Server{
                    Port: 70000,
                },
            },
            expectErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateConfig(&tt.config)
            if tt.expectErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

## 相关文档

- [数据库系统完整指南](DATABASE_GUIDE.md)
- [日志系统文档](LOGGING.md)
- [JWT 认证系统文档](JWT_AUTH.md)
- [Redis 缓存系统文档](REDIS.md)
- [Viper 官方文档](https://github.com/spf13/viper)

---

**最佳实践**: 使用环境变量管理敏感配置；为不同环境创建独立的配置文件；启用配置热重载提高开发效率；对配置进行验证确保系统稳定性；使用配置缓存提高性能；为配置系统编写完整的测试用例。