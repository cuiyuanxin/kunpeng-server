package examples

import (
	"fmt"
	"log"

	"github.com/cuiyuanxin/kunpeng/internal/service"
	"gorm.io/gorm"
)

// RunRBACRESTfulExample RBAC + RESTful 权限模型示例
func RunRBACRESTfulExample(db *gorm.DB) {
	// 检查数据库连接
	if db == nil {
		log.Fatal("Database connection is required")
		return
	}

	// 创建权限服务
	permissionService := service.NewPermissionService(db)

	// 示例：创建角色和用户
	exampleRoleID := uint(1)
	exampleUserID := uint(1)

	// 1. 传统RBAC权限管理
	fmt.Println("=== 传统RBAC权限管理 ===")
	
	// 为角色添加传统权限
	err := permissionService.AddRolePermission(exampleRoleID, 1) // 假设权限ID为1
	if err != nil {
		fmt.Printf("添加角色权限失败: %v\n", err)
	} else {
		fmt.Println("成功为角色添加传统权限")
	}

	// 为用户分配角色
	err = permissionService.AddUserRole(exampleUserID, exampleRoleID)
	if err != nil {
		fmt.Printf("分配用户角色失败: %v\n", err)
	} else {
		fmt.Println("成功为用户分配角色")
	}

	// 2. RESTful API权限管理
	fmt.Println("\n=== RESTful API权限管理 ===")

	// 为角色添加RESTful API权限
	restfulPermissions := []map[string]string{
		{"path": "/api/v1/users/*", "method": "GET"},
		{"path": "/api/v1/users", "method": "POST"},
		{"path": "/api/v1/users/*", "method": "PUT"},
		{"path": "/api/v1/users/*", "method": "DELETE"},
		{"path": "/api/v1/roles/*", "method": "GET"},
		{"path": "/api/v1/permissions/*", "method": "GET"},
	}

	err = permissionService.SyncRoleRESTfulPermissions(exampleRoleID, restfulPermissions)
	if err != nil {
		fmt.Printf("同步RESTful权限失败: %v\n", err)
	} else {
		fmt.Println("成功同步角色RESTful权限")
	}

	// 3. 权限检查示例
	fmt.Println("\n=== 权限检查示例 ===")

	// 检查RESTful API权限
	testCases := []struct {
		path   string
		method string
		desc   string
	}{
		{"/api/v1/users/123", "GET", "获取用户信息"},
		{"/api/v1/users", "POST", "创建用户"},
		{"/api/v1/users/123", "PUT", "更新用户"},
		{"/api/v1/users/123", "DELETE", "删除用户"},
		{"/api/v1/roles/1", "GET", "获取角色信息"},
		{"/api/v1/admin/settings", "GET", "获取管理员设置（无权限）"},
	}

	for _, tc := range testCases {
		hasPermission, err := permissionService.CheckUserPermissionByPath(exampleUserID, tc.path, tc.method)
		if err != nil {
			fmt.Printf("检查权限失败 [%s %s]: %v\n", tc.method, tc.path, err)
			continue
		}

		status := "拒绝"
		if hasPermission {
			status = "允许"
		}
		fmt.Printf("权限检查 [%s %s] %s: %s\n", tc.method, tc.path, tc.desc, status)
	}

	// 4. 获取角色权限列表
	fmt.Println("\n=== 角色权限列表 ===")

	// 获取传统权限
	traditionalPerms, err := permissionService.GetRolePermissions(exampleRoleID)
	if err != nil {
		fmt.Printf("获取传统权限失败: %v\n", err)
	} else {
		fmt.Printf("传统权限: %v\n", traditionalPerms)
	}

	// 获取RESTful权限
	restfulPerms, err := permissionService.GetRoleRESTfulPermissions(exampleRoleID)
	if err != nil {
		fmt.Printf("获取RESTful权限失败: %v\n", err)
	} else {
		fmt.Println("RESTful权限:")
		for _, perm := range restfulPerms {
			fmt.Printf("  %s %s\n", perm["method"], perm["path"])
		}
	}

	fmt.Println("\n=== 示例完成 ===")
}

// 权限模型说明:
// 1. RBAC模型: 用户 -> 角色 -> 权限
// 2. RESTful模型: 直接匹配HTTP方法和路径
// 3. 融合模型: 优先检查RESTful权限，回退到传统RBAC权限
//
// 配置文件 rbac_model.conf 说明:
// - keyMatch2: 支持路径通配符匹配，如 /api/v1/users/* 可以匹配 /api/v1/users/123
// - regexMatch: 支持正则表达式匹配HTTP方法
// - eft字段: 支持allow/deny效果
//
// 使用场景:
// 1. 细粒度API权限控制
// 2. 动态路径权限管理
// 3. 多租户权限隔离
// 4. 微服务权限统一管理
