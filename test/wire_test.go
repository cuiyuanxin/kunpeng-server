package test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/cuiyuanxin/kunpeng/internal/config"
	"github.com/cuiyuanxin/kunpeng/cmd/wire"
)

// TestWireInitialization 测试Wire依赖注入初始化
func TestWireInitialization(t *testing.T) {
	// 创建测试配置
	cfg := &config.Config{
		App: config.App{
			Name:        "test-app",
			Version:     "1.0.0",
			Environment: "test",
			Debug:       true,
		},
		JWT: config.JWT{
			Secret:     "test-secret",
			ExpireTime: 3600 * time.Second,
			Issuer:     "test-issuer",
		},
	}

	// 创建内存数据库
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// 使用Wire初始化应用程序
	app := wire.InitializeApp(cfg, db)

	// 验证应用程序实例
	assert.NotNil(t, app)
	assert.NotNil(t, app.Config)
	assert.NotNil(t, app.DB)
	assert.NotNil(t, app.JWTManager)
	assert.NotNil(t, app.UserService)
	assert.NotNil(t, app.RoleService)
	assert.NotNil(t, app.PermissionService)
	assert.NotNil(t, app.DepartmentService)
	assert.NotNil(t, app.FileService)
	assert.NotNil(t, app.StatService)

	// 验证配置
	assert.Equal(t, "test-app", app.Config.App.Name)
	assert.Equal(t, "1.0.0", app.Config.App.Version)
	assert.Equal(t, "test", app.Config.App.Environment)

	// 验证JWT管理器功能
	token, err := app.JWTManager.GenerateToken(1, "testuser", "user")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// 验证令牌解析
	claims, err := app.JWTManager.ParseToken(token)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), claims.UserID)
	assert.Equal(t, "testuser", claims.Username)
	assert.Equal(t, "user", claims.Role)
}

// TestWireServiceDependencies 测试Wire服务依赖关系
func TestWireServiceDependencies(t *testing.T) {
	// 创建测试配置
	cfg := &config.Config{
		App: config.App{
			Name:        "test-app",
			Version:     "1.0.0",
			Environment: "test",
			Debug:       true,
		},
		JWT: config.JWT{
			Secret:     "test-secret",
			ExpireTime: 3600 * time.Second,
			Issuer:     "test-issuer",
		},
	}

	// 创建内存数据库
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// 使用Wire初始化应用程序
	app := wire.InitializeApp(cfg, db)

	// 验证所有服务都使用相同的数据库实例
	assert.Same(t, app.DB, db)

	// 验证所有服务都使用相同的配置实例
	assert.Same(t, app.Config, cfg)

	// 验证服务实例的唯一性（每次调用InitializeApp应该创建新的服务实例）
	app2 := wire.InitializeApp(cfg, db)
	assert.NotSame(t, app.UserService, app2.UserService)
	assert.NotSame(t, app.RoleService, app2.RoleService)
}