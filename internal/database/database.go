package database

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/cuiyuanxin/kunpeng/internal/config"
	klogger "github.com/cuiyuanxin/kunpeng/internal/logger"
	"github.com/cuiyuanxin/kunpeng/pkg/auth"

	"go.uber.org/zap"
	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DatabaseManager 数据库管理器
type DatabaseManager struct {
	mu        sync.RWMutex
	databases map[string]*gorm.DB
	primary   *gorm.DB // 主数据库
	casbin    *auth.CasbinEnforcer // Casbin权限管理器
}

// 全局数据库管理器实例
var (
	manager *DatabaseManager
	once    sync.Once
	DB      *gorm.DB // 向后兼容的全局变量
)

// Init 初始化数据库连接（向后兼容）
func Init(cfg *config.Database) error {
	return InitWithConfig(&config.Config{
		Database: *cfg,
	})
}

// InitWithConfig 使用完整配置初始化数据库
func InitWithConfig(cfg *config.Config) error {
	once.Do(func() {
		manager = &DatabaseManager{
			databases: make(map[string]*gorm.DB),
		}
	})

	// 初始化主数据库（向后兼容）
	if err := manager.initDatabase("primary", &cfg.Database); err != nil {
		return fmt.Errorf("failed to initialize primary database: %w", err)
	}

	// 初始化多数据库
	for name, dbCfg := range cfg.Databases {
		if err := manager.initDatabase(name, &dbCfg); err != nil {
			klogger.Error("Failed to initialize database", zap.String("name", name), zap.Error(err))
			continue
		}
	}

	// 设置主数据库和向后兼容的全局变量
	manager.mu.Lock()
	manager.primary = manager.databases["primary"]
	DB = manager.primary
	manager.mu.Unlock()

	// 初始化Casbin权限管理器
	if err := manager.initCasbin(); err != nil {
		klogger.Error("Failed to initialize Casbin", zap.Error(err))
		return fmt.Errorf("failed to initialize Casbin: %w", err)
	}

	klogger.Info("Database initialization completed", zap.Int("databases", len(manager.databases)))
	return nil
}

// initDatabase 初始化单个数据库连接
func (dm *DatabaseManager) initDatabase(name string, cfg *config.Database) error {
	// 获取数据库驱动
	dialector, err := getDialector(cfg)
	if err != nil {
		return fmt.Errorf("failed to get dialector for %s: %w", name, err)
	}

	// 使用自定义的GORM日志器
	gormLogger := klogger.GetGormLogger()
	if gormLogger == nil {
		// 如果GORM日志器未初始化，使用默认日志器
		gormLogger = logger.Default.LogMode(logger.Info)
	}

	// 创建GORM配置
	gormConfig := &gorm.Config{
		Logger: gormLogger,
	}

	// 连接数据库
	db, err := gorm.Open(dialector, gormConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to database %s: %w", name, err)
	}

	// 获取底层的sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB for %s: %w", name, err)
	}

	// 设置连接池参数
	if cfg.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	}
	if cfg.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	if cfg.ConnMaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	}
	if cfg.ConnMaxIdleTime > 0 {
		sqlDB.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)
	}

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database %s: %w", name, err)
	}

	// 存储数据库连接
	dm.mu.Lock()
	dm.databases[name] = db
	dm.mu.Unlock()

	klogger.Info("Database connected successfully", zap.String("name", name), zap.String("driver", cfg.Driver))
	return nil
}

// initCasbin 初始化Casbin权限管理器
func (dm *DatabaseManager) initCasbin() error {
	if dm.primary == nil {
		return fmt.Errorf("primary database not initialized")
	}

	// 创建Casbin权限管理器
	casbinEnforcer, err := auth.NewCasbinEnforcer(dm.primary)
	if err != nil {
		return fmt.Errorf("failed to create Casbin enforcer: %w", err)
	}

	dm.casbin = casbinEnforcer
	klogger.Info("Casbin enforcer initialized successfully")
	return nil
}

// getDialector 根据数据库驱动获取 GORM 驱动
func getDialector(cfg *config.Database) (gorm.Dialector, error) {
	dsn := cfg.GetDSN()
	switch strings.ToLower(cfg.Driver) {
	case "mysql":
		return mysql.Open(dsn), nil
	// 注意：以下数据库驱动需要先安装相应的包
	case "postgres", "postgresql":
		return postgres.Open(dsn), nil
	case "sqlite", "sqlite3":
		return sqlite.Open(dsn), nil
	case "sqlserver", "mssql":
		return sqlserver.Open(dsn), nil
	case "clickhouse":
		return clickhouse.Open(dsn), nil
	default:
		return nil, fmt.Errorf("unsupported database driver: %s (currently only mysql is available)", cfg.Driver)
	}
}

// Close 关闭所有数据库连接
func Close() error {
	if manager == nil {
		return nil
	}

	manager.mu.Lock()
	defer manager.mu.Unlock()

	var lastErr error
	for name, db := range manager.databases {
		if db != nil {
			sqlDB, err := db.DB()
			if err != nil {
				lastErr = err
				continue
			}
			if err := sqlDB.Close(); err != nil {
				lastErr = err
				klogger.Error("Failed to close database", zap.String("name", name), zap.Error(err))
			} else {
				klogger.Info("Database closed successfully", zap.String("name", name))
			}
		}
	}

	// 清空数据库映射
	manager.databases = make(map[string]*gorm.DB)
	manager.primary = nil
	DB = nil

	return lastErr
}

// GetDB 获取主数据库实例（向后兼容）
func GetDB() *gorm.DB {
	return DB
}

// GetDatabase 根据名称获取数据库实例
func GetDatabase(name string) *gorm.DB {
	if manager == nil {
		return nil
	}

	manager.mu.RLock()
	defer manager.mu.RUnlock()

	return manager.databases[name]
}

// GetPrimaryDatabase 获取主数据库实例
func GetPrimaryDatabase() *gorm.DB {
	if manager == nil {
		return nil
	}

	manager.mu.RLock()
	defer manager.mu.RUnlock()

	return manager.primary
}

// GetCasbin 获取Casbin权限管理器实例
func GetCasbin() *auth.CasbinEnforcer {
	if manager == nil {
		return nil
	}

	manager.mu.RLock()
	defer manager.mu.RUnlock()

	return manager.casbin
}

// ListDatabases 列出所有已初始化的数据库名称
func ListDatabases() []string {
	if manager == nil {
		return nil
	}

	manager.mu.RLock()
	defer manager.mu.RUnlock()

	names := make([]string, 0, len(manager.databases))
	for name := range manager.databases {
		names = append(names, name)
	}
	return names
}

// GetDatabaseCount 获取数据库连接数量
func GetDatabaseCount() int {
	if manager == nil {
		return 0
	}

	manager.mu.RLock()
	defer manager.mu.RUnlock()

	return len(manager.databases)
}

// AutoMigrate 在主数据库上自动迁移数据库表（向后兼容）
func AutoMigrate(models ...interface{}) error {
	return AutoMigrateOnDatabase("primary", models...)
}

// AutoMigrateOnDatabase 在指定数据库上自动迁移数据库表
func AutoMigrateOnDatabase(dbName string, models ...interface{}) error {
	db := GetDatabase(dbName)
	if db == nil {
		return fmt.Errorf("database %s not initialized", dbName)
	}

	for _, model := range models {
		if err := db.AutoMigrate(model); err != nil {
			return fmt.Errorf("failed to migrate model %T on database %s: %w", model, dbName, err)
		}
	}

	klogger.Info("Database migration completed successfully", zap.String("database", dbName))
	return nil
}

// AutoMigrateAll 在所有数据库上自动迁移数据库表
func AutoMigrateAll(models ...interface{}) error {
	if manager == nil {
		return fmt.Errorf("database manager not initialized")
	}

	manager.mu.RLock()
	databases := make(map[string]*gorm.DB)
	for name, db := range manager.databases {
		databases[name] = db
	}
	manager.mu.RUnlock()

	for name := range databases {
		if err := AutoMigrateOnDatabase(name, models...); err != nil {
			klogger.Error("Failed to migrate on database", zap.String("name", name), zap.Error(err))
			continue
		}
	}

	return nil
}

// Transaction 在主数据库上执行事务（向后兼容）
func Transaction(fn func(*gorm.DB) error) error {
	return TransactionOnDatabase("primary", fn)
}

// TransactionOnDatabase 在指定数据库上执行事务
func TransactionOnDatabase(dbName string, fn func(*gorm.DB) error) error {
	db := GetDatabase(dbName)
	if db == nil {
		return fmt.Errorf("database %s not initialized", dbName)
	}
	return db.Transaction(fn)
}

// Paginate 分页查询
func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

// HealthCheck 主数据库健康检查（向后兼容）
func HealthCheck() error {
	return HealthCheckDatabase("primary")
}

// HealthCheckDatabase 指定数据库健康检查
func HealthCheckDatabase(dbName string) error {
	db := GetDatabase(dbName)
	if db == nil {
		return fmt.Errorf("database %s not initialized", dbName)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return sqlDB.PingContext(ctx)
}

// HealthCheckAll 所有数据库健康检查
func HealthCheckAll() map[string]error {
	results := make(map[string]error)

	if manager == nil {
		results["manager"] = fmt.Errorf("database manager not initialized")
		return results
	}

	manager.mu.RLock()
	databases := make(map[string]*gorm.DB)
	for name, db := range manager.databases {
		databases[name] = db
	}
	manager.mu.RUnlock()

	for name := range databases {
		results[name] = HealthCheckDatabase(name)
	}

	return results
}

// gRPC支持相关功能

// GetGRPCDatabase 获取支持gRPC的数据库连接
func GetGRPCDatabase(dbName string) (*gorm.DB, error) {
	db := GetDatabase(dbName)
	if db == nil {
		return nil, fmt.Errorf("database %s not initialized", dbName)
	}

	// 这里可以添加gRPC特定的配置或连接池管理
	// 目前直接返回数据库连接，后续可以根据需要扩展
	return db, nil
}

// GetGRPCDatabasePool 获取gRPC数据库连接池
func GetGRPCDatabasePool(dbName string) ([]*gorm.DB, error) {
	db := GetDatabase(dbName)
	if db == nil {
		return nil, fmt.Errorf("database %s not initialized", dbName)
	}

	// 简单实现：返回单个连接的切片
	// 在实际应用中，可以根据配置创建多个连接
	return []*gorm.DB{db}, nil
}

// DatabaseInfo 数据库信息结构
type DatabaseInfo struct {
	Name        string `json:"name"`
	Driver      string `json:"driver"`
	Host        string `json:"host"`
	Port        int    `json:"port"`
	Database    string `json:"database"`
	Connected   bool   `json:"connected"`
	GRPCEnabled bool   `json:"grpc_enabled"`
}

// GetDatabaseInfo 获取数据库信息
func GetDatabaseInfo() []DatabaseInfo {
	var infos []DatabaseInfo

	if manager == nil {
		return infos
	}

	manager.mu.RLock()
	defer manager.mu.RUnlock()

	for name, db := range manager.databases {
		info := DatabaseInfo{
			Name:      name,
			Connected: db != nil,
		}

		if db != nil {
			// 尝试获取连接信息（这里需要根据实际配置来填充）
			// 由于我们没有直接访问配置的方式，这里只设置基本信息
			info.Driver = "unknown" // 实际应用中应该从配置中获取
		}

		infos = append(infos, info)
	}

	return infos
}
