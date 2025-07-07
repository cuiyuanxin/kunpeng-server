package migrate

import (
	"fmt"

	"github.com/cuiyuanxin/kunpeng/internal/database"
	"github.com/cuiyuanxin/kunpeng/internal/model"
	klogger "github.com/cuiyuanxin/kunpeng/internal/logger"
	"go.uber.org/zap"
)

// AllModels 所有需要迁移的模型
var AllModels = []interface{}{
	// 前台用户模型
	&model.User{},
	
	// 后台管理员模型
	&model.AdminUser{},
	&model.AdminRole{},
	&model.AdminPermission{},
	&model.AdminDepartment{},
	&model.AdminLoginLog{},
	&model.AdminOperationLog{},
	&model.AdminConfig{},
	
	// 在这里添加新的模型
}

// AutoMigrate 自动迁移所有数据库表
func AutoMigrate() error {
	if database.DB == nil {
		return fmt.Errorf("database not initialized")
	}

	klogger.Info("Starting database migration...")

	for _, model := range AllModels {
		klogger.Info("Migrating model", zap.String("model", fmt.Sprintf("%T", model)))
		if err := database.DB.AutoMigrate(model); err != nil {
			return fmt.Errorf("failed to migrate model %T: %w", model, err)
		}
	}

	klogger.Info("Database migration completed successfully")
	return nil
}

// DropAllTables 删除所有表（危险操作，仅用于开发环境）
func DropAllTables() error {
	if database.DB == nil {
		return fmt.Errorf("database not initialized")
	}

	klogger.Warn("Dropping all database tables...")

	// 反向删除表，避免外键约束问题
	for i := len(AllModels) - 1; i >= 0; i-- {
		model := AllModels[i]
		klogger.Info("Dropping table for model", zap.String("model", fmt.Sprintf("%T", model)))
		if err := database.DB.Migrator().DropTable(model); err != nil {
			return fmt.Errorf("failed to drop table for model %T: %w", model, err)
		}
	}

	klogger.Info("All tables dropped successfully")
	return nil
}

// CreateTables 创建所有表
func CreateTables() error {
	return AutoMigrate()
}

// ResetDatabase 重置数据库（删除所有表后重新创建）
func ResetDatabase() error {
	klogger.Warn("Resetting database...")

	if err := DropAllTables(); err != nil {
		return fmt.Errorf("failed to drop tables: %w", err)
	}

	if err := CreateTables(); err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	klogger.Info("Database reset completed successfully")
	return nil
}