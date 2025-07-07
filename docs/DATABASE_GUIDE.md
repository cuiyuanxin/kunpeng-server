# 数据库系统完整指南

## 概述

本项目提供了完整的多数据库支持系统，基于 GORM 构建，支持 MySQL、PostgreSQL、SQLite、SQL Server、ClickHouse 等多种数据库类型。系统支持单数据库和多数据库配置，提供连接池管理、健康检查、事务支持、迁移工具等功能。

## 🚀 支持的数据库类型

### 1. MySQL
- **用途**: 生产环境推荐，主要业务数据
- **特点**: 高性能、稳定可靠、生态完善
- **配置**: 支持完整的连接池和SSL配置

### 2. PostgreSQL
- **用途**: 企业级应用首选，复杂查询和分析
- **特点**: 功能强大、标准兼容、扩展性好
- **配置**: 支持模式(Schema)配置和高级特性

### 3. SQLite
- **用途**: 轻量级应用、开发测试、本地缓存
- **特点**: 无服务器、零配置、文件数据库
- **配置**: 仅需文件路径配置

### 4. SQL Server
- **用途**: 微软生态系统、企业报表
- **特点**: 与.NET集成良好、企业级功能
- **配置**: 支持Windows认证和SQL认证

### 5. ClickHouse
- **用途**: 大数据分析、OLAP场景、日志存储
- **特点**: 列式存储、高性能分析、压缩率高
- **配置**: 支持集群和分布式配置

## 📦 数据库驱动安装

### 快速安装

本项目默认只包含 MySQL 驱动。如需使用其他数据库，请按照以下步骤安装相应的驱动：

#### PostgreSQL
```bash
go get gorm.io/driver/postgres
```

#### SQLite
```bash
go get gorm.io/driver/sqlite
```

#### SQL Server
```bash
go get gorm.io/driver/sqlserver
```

#### ClickHouse
```bash
go get gorm.io/driver/clickhouse
```

### 完整安装

如果你想安装所有支持的数据库驱动：

```bash
# 安装所有数据库驱动
go get gorm.io/driver/postgres
go get gorm.io/driver/sqlite
go get gorm.io/driver/sqlserver
go get gorm.io/driver/clickhouse
```

### 代码配置

安装驱动后，需要在 `internal/database/database.go` 中启用相应的驱动：

```go
// 更新导入部分
import (
    "gorm.io/driver/clickhouse"
    "gorm.io/driver/mysql"
    "gorm.io/driver/postgres"
    "gorm.io/driver/sqlite"
    "gorm.io/driver/sqlserver"
    "gorm.io/gorm"
)

// 更新 getDialector 函数
func getDialector(cfg *config.Database) (gorm.Dialector, error) {
    dsn := cfg.GetDSN()
    switch strings.ToLower(cfg.Driver) {
    case "mysql":
        return mysql.Open(dsn), nil
    case "postgres", "postgresql":
        return postgres.Open(dsn), nil
    case "sqlite", "sqlite3":
        return sqlite.Open(dsn), nil
    case "sqlserver", "mssql":
        return sqlserver.Open(dsn), nil
    case "clickhouse":
        return clickhouse.Open(dsn), nil
    default:
        return nil, fmt.Errorf("unsupported database driver: %s", cfg.Driver)
    }
}
```

## ⚙️ 配置说明

### 单数据库配置（向后兼容）

```yaml
database:
  driver: mysql
  host: localhost
  port: 3306
  username: root
  password: password
  database: kunpeng
  charset: utf8mb4
  timezone: Asia/Shanghai
  max_open_conns: 100
  max_idle_conns: 10
  conn_max_lifetime: 3600s
  conn_max_idle_time: 1800s
```

### 多数据库配置

```yaml
# 主数据库配置（向后兼容）
database:
  driver: "mysql"
  host: "localhost"
  port: 3306
  username: "root"
  password: "password"
  database: "kunpeng_main"
  charset: "utf8mb4"
  timezone: "Local"
  max_open_conns: 100
  max_idle_conns: 10
  conn_max_lifetime: 3600s
  conn_max_idle_time: 1800s
  grpc_enabled: true
  grpc_pool_size: 5

# 多数据库配置
databases:
  # MySQL 用户数据库
  user_db:
    driver: "mysql"
    host: "localhost"
    port: 3306
    username: "root"
    password: "password"
    database: "kunpeng_users"
    charset: "utf8mb4"
    timezone: "Local"
    max_open_conns: 50
    max_idle_conns: 5
    conn_max_lifetime: 3600s
    grpc_enabled: true
    grpc_pool_size: 3

  # PostgreSQL 分析数据库
  analytics_db:
    driver: "postgres"
    host: "localhost"
    port: 5432
    username: "postgres"
    password: "password"
    database: "kunpeng_analytics"
    schema: "public"
    ssl_mode: "disable"
    timezone: "UTC"
    max_open_conns: 30
    max_idle_conns: 3
    conn_max_lifetime: 3600s
    grpc_enabled: true
    grpc_pool_size: 2

  # SQLite 缓存数据库
  cache_db:
    driver: "sqlite"
    file_path: "./data/cache.db"
    max_open_conns: 10
    max_idle_conns: 2
    conn_max_lifetime: 1800s
    grpc_enabled: false

  # SQL Server 报表数据库
  report_db:
    driver: "sqlserver"
    host: "localhost"
    port: 1433
    username: "sa"
    password: "YourPassword123"
    database: "kunpeng_reports"
    max_open_conns: 20
    max_idle_conns: 2
    conn_max_lifetime: 3600s
    grpc_enabled: true
    grpc_pool_size: 2

  # ClickHouse 日志数据库
  log_db:
    driver: "clickhouse"
    host: "localhost"
    port: 9000
    username: "default"
    password: ""
    database: "kunpeng_logs"
    max_open_conns: 15
    max_idle_conns: 2
    conn_max_lifetime: 3600s
    grpc_enabled: false
```

### 数据库特定配置

#### MySQL 配置示例
```yaml
mysql_db:
  driver: "mysql"
  host: "localhost"
  port: 3306
  username: "root"
  password: "password"
  database: "kunpeng"
  charset: "utf8mb4"              # 字符集
  timezone: "Asia/Shanghai"       # 时区
  # 或使用自定义DSN
  dsn: "root:password@tcp(localhost:3306)/kunpeng?charset=utf8mb4&parseTime=True&loc=Local"
```

#### PostgreSQL 配置示例
```yaml
postgres_db:
  driver: "postgres"
  host: "localhost"
  port: 5432
  username: "postgres"
  password: "password"
  database: "kunpeng"
  schema: "public"                # PostgreSQL 模式
  ssl_mode: "disable"             # SSL 模式
  timezone: "UTC"                 # 时区
```

#### SQLite 配置示例
```yaml
sqlite_db:
  driver: "sqlite"
  file_path: "./data/app.db"       # 数据库文件路径
  # 注意：SQLite 不需要 host, port, username, password
```

#### SQL Server 配置示例
```yaml
sqlserver_db:
  driver: "sqlserver"
  host: "localhost"
  port: 1433
  username: "sa"
  password: "YourPassword123"      # 需符合复杂性要求
  database: "kunpeng"
```

#### ClickHouse 配置示例
```yaml
clickhouse_db:
  driver: "clickhouse"
  host: "localhost"
  port: 9000
  username: "default"
  password: ""                     # ClickHouse 默认无密码
  database: "kunpeng"
```

## 🔧 使用方法

### 1. 初始化数据库

```go
// 使用配置初始化多数据库
if err := database.InitWithConfig(cfg); err != nil {
    log.Fatal("Failed to init databases:", err)
}
defer database.Close()
```

### 2. 获取数据库连接

```go
// 获取主数据库
mainDB := database.GetPrimaryDatabase()

// 获取指定数据库
userDB := database.GetDatabase("user_db")
analyticsDB := database.GetDatabase("analytics_db")

// 获取 gRPC 数据库连接
grpcDB, err := database.GetGRPCDatabase("user_db")
if err != nil {
    log.Fatal("Failed to get gRPC database:", err)
}
```

### 3. 数据库操作

```go
// 在指定数据库上执行迁移
database.AutoMigrateOnDatabase("user_db", &User{}, &Order{})

// 在指定数据库上执行事务
database.TransactionOnDatabase("user_db", func(tx *gorm.DB) error {
    // 事务操作
    return tx.Create(&user).Error
})

// 健康检查
healthResults := database.HealthCheckAll()
for dbName, healthy := range healthResults {
    fmt.Printf("Database %s: %v\n", dbName, healthy)
}
```

### 4. 数据库信息查看

```go
// 获取所有数据库信息
dbInfos := database.GetDatabaseInfo()
for _, info := range dbInfos {
    fmt.Printf("Database: %s, Driver: %s, Connected: %v\n",
        info.Name, info.Driver, info.Connected)
}

// 获取数据库列表
databases := database.ListDatabases()
fmt.Printf("Available databases: %v\n", databases)

// 获取数据库数量
count := database.GetDatabaseCount()
fmt.Printf("Total databases: %d\n", count)
```

## 🔄 数据库迁移

### 命令行迁移

```bash
# 在所有数据库上执行迁移
go run cmd/migrate/main.go

# 在主数据库上执行迁移
go run cmd/migrate/main.go -database=primary

# 在指定数据库上执行迁移
go run cmd/migrate/main.go -database=user_db

# 重置所有数据库
go run cmd/migrate/main.go -action=reset -database=all

# 删除所有表
go run cmd/migrate/main.go -action=drop -database=all
```

### 程序化迁移

```go
// 自动迁移所有数据库
err := migrate.AutoMigrate()
if err != nil {
    log.Fatal("Migration failed:", err)
}

// 在指定数据库上迁移
err = migrate.AutoMigrateOnDatabase("user_db", &User{}, &Order{})
if err != nil {
    log.Fatal("Migration failed:", err)
}
```

## 📊 性能优化

### 连接池配置

```yaml
database:
  max_idle_conns: 10        # 最大空闲连接数
  max_open_conns: 100       # 最大打开连接数
  conn_max_lifetime: 3600s  # 连接最大生存时间
  conn_max_idle_time: 1800s # 连接最大空闲时间
```

### gRPC 优化

```yaml
database:
  grpc_enabled: true        # 启用 gRPC 支持
  grpc_pool_size: 10        # gRPC 连接池大小
```

### 读写分离配置

```yaml
databases:
  # 主数据库（写库）
  master_db:
    driver: "mysql"
    host: "mysql-master"
    port: 3306
    username: "root"
    password: "password"
    database: "kunpeng_master"
    max_open_conns: 50
    max_idle_conns: 5
    grpc_enabled: true
    grpc_pool_size: 3

  # 从数据库（读库）
  slave_db:
    driver: "mysql"
    host: "mysql-slave"
    port: 3306
    username: "readonly"
    password: "password"
    database: "kunpeng_slave"
    max_open_conns: 80        # 读库可以更多连接
    max_idle_conns: 8
    grpc_enabled: true
    grpc_pool_size: 5
```

## 🔍 监控和调试

### 健康检查

```go
// 检查所有数据库
healthResults := database.HealthCheckAll()
for dbName, healthy := range healthResults {
    if !healthy {
        log.Printf("Database %s is unhealthy", dbName)
    }
}

// 检查指定数据库
err := database.HealthCheckDatabase("user_db")
if err != nil {
    log.Printf("Database health check failed: %v", err)
}
```

### 连接状态监控

```go
// 获取数据库统计信息
stats := database.GetDatabaseStats("user_db")
fmt.Printf("Open connections: %d\n", stats.OpenConnections)
fmt.Printf("In use: %d\n", stats.InUse)
fmt.Printf("Idle: %d\n", stats.Idle)
```

## 🔧 验证安装

安装完成后，可以通过以下方式验证：

### 1. 编译检查
```bash
go build ./cmd/main.go
```

### 2. 运行测试
```bash
go run examples/multi_database_example.go
```

### 3. 检查配置
使用配置文件测试不同数据库连接：

```bash
# 使用开发环境配置
go run cmd/main.go -config=configs/config.dev.yaml

# 使用生产环境配置
go run cmd/main.go -config=configs/config.prod.yaml
```

## ⚠️ 注意事项

### 1. 依赖管理
安装新驱动后，记得运行：
```bash
go mod tidy
```

### 2. 版本兼容性
确保所有 GORM 相关包版本兼容：
```bash
go list -m gorm.io/gorm
go list -m gorm.io/driver/mysql
# 其他驱动...
```

### 3. 数据库服务
确保目标数据库服务已启动并可访问。

### 4. 配置文件
根据实际数据库配置更新 DSN 连接字符串。

### 5. 安全配置
- **生产环境**: 启用SSL，使用专用用户
- **开发环境**: 可以禁用SSL，使用简单配置

### 6. 性能考虑
- 多数据库会增加内存使用，请根据实际需求配置
- 确保正确关闭数据库连接以避免资源泄漏
- 跨数据库事务需要特别注意数据一致性

## 🐛 故障排查

### 编译错误
如果遇到编译错误，请检查：
- Go 版本是否满足要求（推荐 Go 1.19+）
- 网络连接是否正常
- 代理设置是否正确

### 连接错误
如果遇到数据库连接错误，请检查：
- 数据库服务是否启动
- 连接参数是否正确
- 防火墙设置
- 数据库用户权限

### 性能问题
如果遇到性能问题，请调整：
- 连接池大小
- 超时设置
- 日志级别

## 🎯 使用场景

### 微服务架构
- 用户服务使用 MySQL
- 分析服务使用 ClickHouse
- 缓存服务使用 SQLite

### 读写分离
- 主数据库用于写操作
- 只读副本用于查询操作

### 多租户系统
- 每个租户使用独立的数据库
- 共享配置和元数据库

### 数据迁移
- 从旧系统（SQL Server）迁移到新系统（PostgreSQL）
- 保持双写确保数据一致性

## 🔄 迁移指南

### 从单数据库迁移
1. 保持现有配置不变（向后兼容）
2. 根据需要添加新的数据库配置
3. 逐步迁移业务逻辑到新的数据库
4. 更新代码使用新的 API

### 代码更新
```go
// 旧方式
db := database.GetDB()

// 新方式
db := database.GetPrimaryDatabase()  // 获取主数据库
userDB := database.GetDatabase("user_db")  // 获取指定数据库
```

## 📚 相关文档

- [GORM 官方文档](https://gorm.io/docs/)
- [MySQL 驱动文档](https://gorm.io/docs/connecting_to_the_database.html#MySQL)
- [PostgreSQL 驱动文档](https://gorm.io/docs/connecting_to_the_database.html#PostgreSQL)
- [SQLite 驱动文档](https://gorm.io/docs/connecting_to_the_database.html#SQLite)
- [SQL Server 驱动文档](https://gorm.io/docs/connecting_to_the_database.html#SQL-Server)
- [ClickHouse 驱动文档](https://gorm.io/docs/connecting_to_the_database.html#ClickHouse)

---

**项目状态**: ✅ 完成  
**版本**: v2.0.0  
**兼容性**: 完全向后兼容  
**测试状态**: 已通过测试

**提示**: 建议在开发环境中先测试单个数据库驱动，确认无误后再在生产环境中部署。