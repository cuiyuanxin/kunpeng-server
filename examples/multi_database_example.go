package examples

import (
	"fmt"
	"log"

	"github.com/cuiyuanxin/kunpeng/internal/config"
	"github.com/cuiyuanxin/kunpeng/internal/database"
	klogger "github.com/cuiyuanxin/kunpeng/internal/logger"
	"gorm.io/gorm"
)

// ExampleUser 用户模型示例
type ExampleUser struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `gorm:"size:100;not null"`
	Email string `gorm:"size:100;uniqueIndex"`
}

// ExampleOrder 订单模型示例
type ExampleOrder struct {
	ID     uint   `gorm:"primaryKey"`
	UserID uint   `gorm:"not null"`
	Amount float64 `gorm:"not null"`
	Status string `gorm:"size:50;default:'pending'"`
}

// ExampleAnalytics 分析数据模型示例
type ExampleAnalytics struct {
	ID        uint   `gorm:"primaryKey"`
	EventName string `gorm:"size:100;not null"`
	EventData string `gorm:"type:text"`
	Timestamp int64  `gorm:"not null"`
}

func main() {
	// 初始化配置
	cfg, err := config.Init("configs/config.dev.yaml")
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// 初始化日志
	if err := klogger.InitWithEnvironment(&cfg.Logging, cfg.App.Environment); err != nil {
		log.Fatal("Failed to init logger:", err)
	}
	defer klogger.Sync()

	// 初始化多数据库
	if err := database.InitWithConfig(cfg); err != nil {
		log.Fatal("Failed to init databases:", err)
	}
	defer database.Close()

	// 显示可用的数据库
	databases := database.ListDatabases()
	fmt.Printf("Available databases: %v\n", databases)
	fmt.Printf("Total databases: %d\n", database.GetDatabaseCount())

	// 健康检查所有数据库
	healthResults := database.HealthCheckAll()
	for name, err := range healthResults {
		if err != nil {
			fmt.Printf("Database %s health check failed: %v\n", name, err)
		} else {
			fmt.Printf("Database %s is healthy\n", name)
		}
	}

	// 示例1: 在不同数据库上执行迁移
	fmt.Println("\n=== 数据库迁移示例 ===")
	
	// 在用户数据库上迁移用户和订单表
	if err := database.AutoMigrateOnDatabase("user_db", &ExampleUser{}, &ExampleOrder{}); err != nil {
		log.Printf("Failed to migrate user_db: %v", err)
	} else {
		fmt.Println("User database migration completed")
	}

	// 在分析数据库上迁移分析表
	if err := database.AutoMigrateOnDatabase("analytics_db", &ExampleAnalytics{}); err != nil {
		log.Printf("Failed to migrate analytics_db: %v", err)
	} else {
		fmt.Println("Analytics database migration completed")
	}

	// 示例2: 在不同数据库上执行操作
	fmt.Println("\n=== 数据库操作示例 ===")
	
	// 在用户数据库上创建用户
	userDB := database.GetDatabase("user_db")
	if userDB != nil {
		user := &ExampleUser{
			Name:  "张三",
			Email: "zhangsan@example.com",
		}
		if err := userDB.Create(user).Error; err != nil {
			log.Printf("Failed to create user: %v", err)
		} else {
			fmt.Printf("User created with ID: %d\n", user.ID)
			
			// 创建订单
			order := &ExampleOrder{
				UserID: user.ID,
				Amount: 99.99,
				Status: "completed",
			}
			if err := userDB.Create(order).Error; err != nil {
				log.Printf("Failed to create order: %v", err)
			} else {
				fmt.Printf("Order created with ID: %d\n", order.ID)
			}
		}
	}

	// 在分析数据库上记录事件
	analyticsDB := database.GetDatabase("analytics_db")
	if analyticsDB != nil {
		analytics := &ExampleAnalytics{
			EventName: "user_registration",
			EventData: `{"source": "web", "campaign": "summer2024"}`,
			Timestamp: 1234567890,
		}
		if err := analyticsDB.Create(analytics).Error; err != nil {
			log.Printf("Failed to create analytics: %v", err)
		} else {
			fmt.Printf("Analytics event created with ID: %d\n", analytics.ID)
		}
	}

	// 示例3: 事务操作
	fmt.Println("\n=== 事务操作示例 ===")
	
	// 在用户数据库上执行事务
	err = database.TransactionOnDatabase("user_db", func(tx *gorm.DB) error {
		// 创建用户
		user := &ExampleUser{
			Name:  "李四",
			Email: "lisi@example.com",
		}
		if err := tx.Create(user).Error; err != nil {
			return err
		}

		// 创建订单
		order := &ExampleOrder{
			UserID: user.ID,
			Amount: 199.99,
			Status: "pending",
		}
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		fmt.Printf("Transaction completed: User ID %d, Order ID %d\n", user.ID, order.ID)
		return nil
	})

	if err != nil {
		log.Printf("Transaction failed: %v", err)
	} else {
		fmt.Println("Transaction completed successfully")
	}

	// 示例4: gRPC 数据库连接
	fmt.Println("\n=== gRPC 数据库连接示例 ===")
	
	// 获取支持 gRPC 的数据库连接
	grpcDB, err := database.GetGRPCDatabase("user_db")
	if err != nil {
		log.Printf("Failed to get gRPC database: %v", err)
	} else {
		fmt.Println("gRPC database connection obtained successfully")
		
		// 使用 gRPC 数据库连接执行查询
		var userCount int64
		if err := grpcDB.Model(&ExampleUser{}).Count(&userCount).Error; err != nil {
			log.Printf("Failed to count users: %v", err)
		} else {
			fmt.Printf("Total users in user_db: %d\n", userCount)
		}
	}

	// 获取 gRPC 数据库连接池
	grpcPool, err := database.GetGRPCDatabasePool("user_db")
	if err != nil {
		log.Printf("Failed to get gRPC database pool: %v", err)
	} else {
		fmt.Printf("gRPC database pool size: %d\n", len(grpcPool))
	}

	// 示例5: 数据库信息
	fmt.Println("\n=== 数据库信息示例 ===")
	
	dbInfos := database.GetDatabaseInfo()
	for _, info := range dbInfos {
		fmt.Printf("Database: %s, Driver: %s, Connected: %v, gRPC: %v\n",
			info.Name, info.Driver, info.Connected, info.GRPCEnabled)
	}

	fmt.Println("\n=== 多数据库示例完成 ===")
}