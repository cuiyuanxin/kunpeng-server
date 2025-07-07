# 数据库迁移系统文档

## 概述

本项目提供了完整的数据库迁移系统，基于 GORM 的 AutoMigrate 功能构建，支持自动迁移、表结构管理、数据重置等功能。迁移系统确保数据库结构与代码模型保持同步，支持开发、测试和生产环境的数据库管理。

## 系统架构

### 迁移系统组件

```
Migration System
├── Model Registry          # 模型注册表
├── Auto Migration          # 自动迁移
├── Table Management        # 表管理
├── Schema Validation       # 结构验证
├── Data Seeding           # 数据填充
├── Backup & Restore       # 备份恢复
└── Migration History      # 迁移历史
```

## 核心功能

### 模型注册

```go
// AllModels 存储所有需要迁移的模型
var AllModels = []interface{}{
    &model.User{},
    &model.Role{},
    &model.Permission{},
    &model.Article{},
    &model.Category{},
    &model.Tag{},
    &model.Comment{},
    &model.File{},
    &model.Log{},
    &model.Config{},
}
```

**模型注册说明**:
- 所有需要迁移的模型都必须在 `AllModels` 中注册
- 模型必须实现 GORM 的标准结构
- 支持关联关系的自动处理
- 按依赖顺序排列模型

### 自动迁移

```go
// AutoMigrate 自动迁移所有数据库表
func AutoMigrate() error {
    // 获取所有数据库连接
    databases := database.GetAllDatabases()
    if len(databases) == 0 {
        return errors.New("no database connections available")
    }
    
    var errors []string
    
    // 对每个数据库执行迁移
    for name, db := range databases {
        klogger.Info("Starting migration for database", zap.String("database", name))
        
        err := migrateDatabase(db, name)
        if err != nil {
            errorMsg := fmt.Sprintf("Migration failed for database %s: %v", name, err)
            errors = append(errors, errorMsg)
            klogger.Error(errorMsg)
            continue
        }
        
        klogger.Info("Migration completed for database", zap.String("database", name))
    }
    
    if len(errors) > 0 {
        return fmt.Errorf("migration errors: %s", strings.Join(errors, "; "))
    }
    
    klogger.Info("All database migrations completed successfully")
    return nil
}

// migrateDatabase 迁移单个数据库
func migrateDatabase(db *gorm.DB, dbName string) error {
    // 开始事务
    tx := db.Begin()
    if tx.Error != nil {
        return fmt.Errorf("failed to begin transaction: %w", tx.Error)
    }
    
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
            klogger.Error("Migration panic, rolling back",
                zap.String("database", dbName),
                zap.Any("error", r),
            )
        }
    }()
    
    // 执行迁移
    for _, model := range AllModels {
        modelName := getModelName(model)
        klogger.Debug("Migrating model",
            zap.String("database", dbName),
            zap.String("model", modelName),
        )
        
        err := tx.AutoMigrate(model)
        if err != nil {
            tx.Rollback()
            return fmt.Errorf("failed to migrate model %s: %w", modelName, err)
        }
    }
    
    // 提交事务
    if err := tx.Commit().Error; err != nil {
        return fmt.Errorf("failed to commit migration transaction: %w", err)
    }
    
    return nil
}

// getModelName 获取模型名称
func getModelName(model interface{}) string {
    modelType := reflect.TypeOf(model)
    if modelType.Kind() == reflect.Ptr {
        modelType = modelType.Elem()
    }
    return modelType.Name()
}
```

### 表管理

```go
// CreateTables 创建所有表
func CreateTables() error {
    databases := database.GetAllDatabases()
    if len(databases) == 0 {
        return errors.New("no database connections available")
    }
    
    for name, db := range databases {
        klogger.Info("Creating tables for database", zap.String("database", name))
        
        for _, model := range AllModels {
            modelName := getModelName(model)
            
            // 检查表是否已存在
            if db.Migrator().HasTable(model) {
                klogger.Debug("Table already exists",
                    zap.String("database", name),
                    zap.String("table", modelName),
                )
                continue
            }
            
            // 创建表
            err := db.Migrator().CreateTable(model)
            if err != nil {
                return fmt.Errorf("failed to create table %s in database %s: %w", modelName, name, err)
            }
            
            klogger.Info("Table created",
                zap.String("database", name),
                zap.String("table", modelName),
            )
        }
    }
    
    return nil
}

// DropAllTables 删除所有表（危险操作）
func DropAllTables() error {
    databases := database.GetAllDatabases()
    if len(databases) == 0 {
        return errors.New("no database connections available")
    }
    
    // 确认操作
    klogger.Warn("DANGER: Dropping all tables from all databases")
    
    for name, db := range databases {
        klogger.Warn("Dropping tables from database", zap.String("database", name))
        
        // 反向删除表（避免外键约束问题）
        for i := len(AllModels) - 1; i >= 0; i-- {
            model := AllModels[i]
            modelName := getModelName(model)
            
            if !db.Migrator().HasTable(model) {
                klogger.Debug("Table does not exist",
                    zap.String("database", name),
                    zap.String("table", modelName),
                )
                continue
            }
            
            err := db.Migrator().DropTable(model)
            if err != nil {
                klogger.Error("Failed to drop table",
                    zap.String("database", name),
                    zap.String("table", modelName),
                    zap.Error(err),
                )
                continue
            }
            
            klogger.Info("Table dropped",
                zap.String("database", name),
                zap.String("table", modelName),
            )
        }
    }
    
    return nil
}

// ResetDatabase 重置数据库（删除后重新创建）
func ResetDatabase() error {
    klogger.Warn("DANGER: Resetting all databases")
    
    // 删除所有表
    err := DropAllTables()
    if err != nil {
        return fmt.Errorf("failed to drop tables: %w", err)
    }
    
    // 重新创建表
    err = CreateTables()
    if err != nil {
        return fmt.Errorf("failed to create tables: %w", err)
    }
    
    klogger.Info("Database reset completed")
    return nil
}
```

## 高级功能

### 迁移验证

```go
// ValidateMigration 验证迁移结果
func ValidateMigration() error {
    databases := database.GetAllDatabases()
    
    for name, db := range databases {
        klogger.Info("Validating migration for database", zap.String("database", name))
        
        for _, model := range AllModels {
            modelName := getModelName(model)
            
            // 检查表是否存在
            if !db.Migrator().HasTable(model) {
                return fmt.Errorf("table %s does not exist in database %s", modelName, name)
            }
            
            // 验证表结构
            err := validateTableStructure(db, model, modelName)
            if err != nil {
                return fmt.Errorf("table structure validation failed for %s in database %s: %w", modelName, name, err)
            }
        }
        
        klogger.Info("Migration validation passed", zap.String("database", name))
    }
    
    return nil
}

// validateTableStructure 验证表结构
func validateTableStructure(db *gorm.DB, model interface{}, modelName string) error {
    // 获取模型的字段信息
    stmt := &gorm.Statement{DB: db}
    err := stmt.Parse(model)
    if err != nil {
        return fmt.Errorf("failed to parse model: %w", err)
    }
    
    // 检查每个字段是否存在
    for _, field := range stmt.Schema.Fields {
        if !db.Migrator().HasColumn(model, field.DBName) {
            return fmt.Errorf("column %s does not exist", field.DBName)
        }
    }
    
    // 检查索引
    for _, index := range stmt.Schema.ParseIndexes() {
        if !db.Migrator().HasIndex(model, index.Name) {
            klogger.Warn("Index does not exist",
                zap.String("model", modelName),
                zap.String("index", index.Name),
            )
        }
    }
    
    return nil
}
```

### 数据填充

```go
// Seeder 数据填充器接口
type Seeder interface {
    Seed(db *gorm.DB) error
    GetName() string
}

// UserSeeder 用户数据填充器
type UserSeeder struct{}

func (s *UserSeeder) GetName() string {
    return "UserSeeder"
}

func (s *UserSeeder) Seed(db *gorm.DB) error {
    // 检查是否已有数据
    var count int64
    db.Model(&model.User{}).Count(&count)
    if count > 0 {
        klogger.Info("Users already exist, skipping seed")
        return nil
    }
    
    // 创建默认用户
    users := []model.User{
        {
            Username: "admin",
            Email:    "admin@example.com",
            Password: "$2a$10$...", // 加密后的密码
            Role:     "admin",
            Status:   1,
        },
        {
            Username: "user",
            Email:    "user@example.com",
            Password: "$2a$10$...", // 加密后的密码
            Role:     "user",
            Status:   1,
        },
    }
    
    for _, user := range users {
        err := db.Create(&user).Error
        if err != nil {
            return fmt.Errorf("failed to create user %s: %w", user.Username, err)
        }
        klogger.Info("User created", zap.String("username", user.Username))
    }
    
    return nil
}

// RunSeeders 运行所有数据填充器
func RunSeeders() error {
    seeders := []Seeder{
        &UserSeeder{},
        &RoleSeeder{},
        &ConfigSeeder{},
    }
    
    databases := database.GetAllDatabases()
    
    for name, db := range databases {
        klogger.Info("Running seeders for database", zap.String("database", name))
        
        for _, seeder := range seeders {
            klogger.Info("Running seeder",
                zap.String("database", name),
                zap.String("seeder", seeder.GetName()),
            )
            
            err := seeder.Seed(db)
            if err != nil {
                return fmt.Errorf("seeder %s failed for database %s: %w", seeder.GetName(), name, err)
            }
        }
    }
    
    return nil
}
```

### 备份和恢复

```go
// BackupManager 备份管理器
type BackupManager struct {
    backupDir string
}

// NewBackupManager 创建备份管理器
func NewBackupManager(backupDir string) *BackupManager {
    return &BackupManager{
        backupDir: backupDir,
    }
}

// CreateBackup 创建数据库备份
func (bm *BackupManager) CreateBackup(dbName string) error {
    db := database.GetDatabase(dbName)
    if db == nil {
        return fmt.Errorf("database %s not found", dbName)
    }
    
    // 确保备份目录存在
    err := os.MkdirAll(bm.backupDir, 0755)
    if err != nil {
        return fmt.Errorf("failed to create backup directory: %w", err)
    }
    
    // 生成备份文件名
    timestamp := time.Now().Format("20060102_150405")
    backupFile := filepath.Join(bm.backupDir, fmt.Sprintf("%s_%s.sql", dbName, timestamp))
    
    // 根据数据库类型执行备份
    switch database.GetDatabaseType(dbName) {
    case "mysql":
        err = bm.backupMySQL(dbName, backupFile)
    case "postgres":
        err = bm.backupPostgreSQL(dbName, backupFile)
    case "sqlite":
        err = bm.backupSQLite(dbName, backupFile)
    default:
        return fmt.Errorf("unsupported database type for backup")
    }
    
    if err != nil {
        return fmt.Errorf("backup failed: %w", err)
    }
    
    klogger.Info("Backup created",
        zap.String("database", dbName),
        zap.String("file", backupFile),
    )
    
    return nil
}

// backupMySQL 备份 MySQL 数据库
func (bm *BackupManager) backupMySQL(dbName, backupFile string) error {
    config := database.GetDatabaseConfig(dbName)
    
    cmd := exec.Command("mysqldump",
        "-h", config.Host,
        "-P", strconv.Itoa(config.Port),
        "-u", config.User,
        "-p"+config.Password,
        "--single-transaction",
        "--routines",
        "--triggers",
        config.DBName,
    )
    
    output, err := os.Create(backupFile)
    if err != nil {
        return err
    }
    defer output.Close()
    
    cmd.Stdout = output
    return cmd.Run()
}

// RestoreBackup 恢复数据库备份
func (bm *BackupManager) RestoreBackup(dbName, backupFile string) error {
    if !filepath.IsAbs(backupFile) {
        backupFile = filepath.Join(bm.backupDir, backupFile)
    }
    
    // 检查备份文件是否存在
    if _, err := os.Stat(backupFile); os.IsNotExist(err) {
        return fmt.Errorf("backup file does not exist: %s", backupFile)
    }
    
    klogger.Warn("DANGER: Restoring database backup",
        zap.String("database", dbName),
        zap.String("backup_file", backupFile),
    )
    
    // 根据数据库类型执行恢复
    switch database.GetDatabaseType(dbName) {
    case "mysql":
        return bm.restoreMySQL(dbName, backupFile)
    case "postgres":
        return bm.restorePostgreSQL(dbName, backupFile)
    case "sqlite":
        return bm.restoreSQLite(dbName, backupFile)
    default:
        return fmt.Errorf("unsupported database type for restore")
    }
}
```

### 迁移历史

```go
// MigrationHistory 迁移历史记录
type MigrationHistory struct {
    ID          uint      `gorm:"primaryKey"`
    Version     string    `gorm:"uniqueIndex;not null"`
    Description string    `gorm:"not null"`
    ExecutedAt  time.Time `gorm:"not null"`
    Success     bool      `gorm:"not null"`
    Error       string    `gorm:"type:text"`
}

// RecordMigration 记录迁移历史
func RecordMigration(version, description string, success bool, err error) {
    databases := database.GetAllDatabases()
    
    for name, db := range databases {
        // 确保迁移历史表存在
        db.AutoMigrate(&MigrationHistory{})
        
        history := MigrationHistory{
            Version:     version,
            Description: description,
            ExecutedAt:  time.Now(),
            Success:     success,
        }
        
        if err != nil {
            history.Error = err.Error()
        }
        
        if err := db.Create(&history).Error; err != nil {
            klogger.Error("Failed to record migration history",
                zap.String("database", name),
                zap.Error(err),
            )
        }
    }
}

// GetMigrationHistory 获取迁移历史
func GetMigrationHistory(dbName string) ([]MigrationHistory, error) {
    db := database.GetDatabase(dbName)
    if db == nil {
        return nil, fmt.Errorf("database %s not found", dbName)
    }
    
    var history []MigrationHistory
    err := db.Order("executed_at DESC").Find(&history).Error
    return history, err
}
```

## 命令行工具

### Makefile 集成

```makefile
# 数据库迁移相关命令
.PHONY: migrate migrate-up migrate-down migrate-reset migrate-status

# 执行数据库迁移
migrate:
	@echo "Running database migration..."
	@go run cmd/migrate/main.go

# 创建数据库表
migrate-up:
	@echo "Creating database tables..."
	@go run cmd/migrate/main.go -action=up

# 删除数据库表
migrate-down:
	@echo "Dropping database tables..."
	@go run cmd/migrate/main.go -action=down

# 重置数据库
migrate-reset:
	@echo "Resetting database..."
	@go run cmd/migrate/main.go -action=reset

# 检查迁移状态
migrate-status:
	@echo "Checking migration status..."
	@go run cmd/migrate/main.go -action=status

# 运行数据填充
seed:
	@echo "Running database seeders..."
	@go run cmd/migrate/main.go -action=seed

# 创建备份
backup:
	@echo "Creating database backup..."
	@go run cmd/migrate/main.go -action=backup

# 恢复备份
restore:
	@echo "Restoring database backup..."
	@go run cmd/migrate/main.go -action=restore -file=$(FILE)
```

### 迁移命令行工具

```go
// cmd/migrate/main.go
package main

import (
    "flag"
    "fmt"
    "log"
    "os"
    
    "github.com/your-project/internal/config"
    "github.com/your-project/internal/database"
    "github.com/your-project/internal/migrate"
    "github.com/your-project/pkg/klogger"
)

func main() {
    var (
        action = flag.String("action", "migrate", "Migration action: migrate, up, down, reset, status, seed, backup, restore")
        file   = flag.String("file", "", "Backup file for restore action")
        dbName = flag.String("db", "", "Database name (empty for all databases)")
    )
    flag.Parse()
    
    // 初始化配置
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }
    
    // 初始化日志
    klogger.Init(&cfg.Log)
    defer klogger.Sync()
    
    // 初始化数据库
    err = database.Init(&cfg.Database)
    if err != nil {
        log.Fatalf("Failed to initialize database: %v", err)
    }
    defer database.Close()
    
    // 执行迁移操作
    switch *action {
    case "migrate":
        err = migrate.AutoMigrate()
    case "up":
        err = migrate.CreateTables()
    case "down":
        err = migrate.DropAllTables()
    case "reset":
        err = migrate.ResetDatabase()
    case "status":
        err = migrate.ValidateMigration()
    case "seed":
        err = migrate.RunSeeders()
    case "backup":
        backupManager := migrate.NewBackupManager("./backups")
        if *dbName != "" {
            err = backupManager.CreateBackup(*dbName)
        } else {
            err = fmt.Errorf("database name required for backup")
        }
    case "restore":
        if *file == "" {
            err = fmt.Errorf("backup file required for restore")
        } else if *dbName == "" {
            err = fmt.Errorf("database name required for restore")
        } else {
            backupManager := migrate.NewBackupManager("./backups")
            err = backupManager.RestoreBackup(*dbName, *file)
        }
    default:
        err = fmt.Errorf("unknown action: %s", *action)
    }
    
    if err != nil {
        log.Fatalf("Migration failed: %v", err)
    }
    
    fmt.Printf("Migration action '%s' completed successfully\n", *action)
}
```

## 环境配置

### 开发环境

```yaml
# config/development.yaml
database:
  default:
    type: "mysql"
    host: "localhost"
    port: 3306
    user: "dev_user"
    password: "dev_password"
    dbname: "kunpeng_dev"
    
migration:
  auto_migrate: true
  run_seeders: true
  backup_enabled: false
```

### 测试环境

```yaml
# config/testing.yaml
database:
  default:
    type: "sqlite"
    dsn: ":memory:"
    
migration:
  auto_migrate: true
  run_seeders: true
  reset_before_test: true
```

### 生产环境

```yaml
# config/production.yaml
database:
  default:
    type: "mysql"
    host: "prod-db.example.com"
    port: 3306
    user: "prod_user"
    password: "${DB_PASSWORD}"
    dbname: "kunpeng_prod"
    
migration:
  auto_migrate: false  # 生产环境手动迁移
  run_seeders: false
  backup_enabled: true
  backup_schedule: "0 2 * * *"  # 每天凌晨2点备份
```

## 最佳实践

### 1. 迁移安全

```go
// 生产环境迁移检查
func SafeMigration() error {
    if gin.Mode() == gin.ReleaseMode {
        // 生产环境需要确认
        fmt.Print("Are you sure you want to run migration in production? (yes/no): ")
        var confirm string
        fmt.Scanln(&confirm)
        
        if confirm != "yes" {
            return fmt.Errorf("migration cancelled")
        }
        
        // 创建备份
        backupManager := NewBackupManager("./backups")
        err := backupManager.CreateBackup("default")
        if err != nil {
            return fmt.Errorf("failed to create backup before migration: %w", err)
        }
    }
    
    return AutoMigrate()
}
```

### 2. 版本控制

```go
// 版本化迁移
func VersionedMigration(version string) error {
    // 检查版本是否已执行
    history, err := GetMigrationHistory("default")
    if err == nil {
        for _, h := range history {
            if h.Version == version && h.Success {
                klogger.Info("Migration already executed", zap.String("version", version))
                return nil
            }
        }
    }
    
    // 执行迁移
    err = AutoMigrate()
    
    // 记录结果
    RecordMigration(version, "Auto migration", err == nil, err)
    
    return err
}
```

### 3. 数据完整性检查

```go
// CheckDataIntegrity 检查数据完整性
func CheckDataIntegrity() error {
    databases := database.GetAllDatabases()
    
    for name, db := range databases {
        klogger.Info("Checking data integrity", zap.String("database", name))
        
        // 检查外键约束
        err := checkForeignKeys(db)
        if err != nil {
            return fmt.Errorf("foreign key check failed for database %s: %w", name, err)
        }
        
        // 检查数据一致性
        err = checkDataConsistency(db)
        if err != nil {
            return fmt.Errorf("data consistency check failed for database %s: %w", name, err)
        }
    }
    
    return nil
}
```

## 故障排查

### 常见问题

1. **迁移失败**
   ```bash
   # 检查数据库连接
   go run cmd/migrate/main.go -action=status
   
   # 查看详细日志
   LOG_LEVEL=debug go run cmd/migrate/main.go
   ```

2. **外键约束错误**
   ```go
   // 临时禁用外键检查（MySQL）
   db.Exec("SET FOREIGN_KEY_CHECKS = 0")
   // 执行迁移
   db.Exec("SET FOREIGN_KEY_CHECKS = 1")
   ```

3. **表已存在错误**
   ```go
   // 检查表是否存在
   if db.Migrator().HasTable(&model.User{}) {
       klogger.Info("Table already exists, skipping")
       return nil
   }
   ```

## 测试

### 迁移测试

```go
func TestMigration(t *testing.T) {
    // 设置测试数据库
    db := setupTestDB(t)
    defer db.Close()
    
    // 测试自动迁移
    err := AutoMigrate()
    assert.NoError(t, err)
    
    // 验证表是否创建
    for _, model := range AllModels {
        assert.True(t, db.Migrator().HasTable(model))
    }
    
    // 测试数据填充
    err = RunSeeders()
    assert.NoError(t, err)
    
    // 验证数据
    var userCount int64
    db.Model(&model.User{}).Count(&userCount)
    assert.Greater(t, userCount, int64(0))
}
```

## 相关文档

- [数据库系统完整指南](DATABASE_GUIDE.md)
- [配置管理文档](../README.md#配置管理)
- [GORM 官方文档](https://gorm.io/docs/)

---

**最佳实践**: 在生产环境执行迁移前创建备份；使用版本控制管理迁移历史；为不同环境配置不同的迁移策略；定期检查数据完整性；使用事务确保迁移的原子性；为迁移操作添加详细的日志记录。