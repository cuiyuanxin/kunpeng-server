# RBAC + RESTful 权限控制系统使用指南

## 概述

本系统基于 Casbin 实现了 RBAC（基于角色的访问控制）与 RESTful API 权限的融合控制模型，支持细粒度的权限管理和灵活的权限策略配置。

## 功能特性

### 1. 权限模型特性
- **RBAC 模型**: 支持用户-角色-权限的三层权限模型
- **RESTful 支持**: 支持基于 HTTP 方法和 URL 路径的权限控制
- **通配符匹配**: 支持路径通配符匹配（如 `/api/users/*`）
- **动态权限**: 支持运行时动态添加、删除权限策略
- **权限继承**: 支持角色继承和权限传递

### 2. 核心组件
- **权限模型配置**: 内嵌在 `internal/auth/casbin.go` 中的 `getRBACModelText()` 函数
- **Casbin 权限管理器**: `internal/auth/casbin.go`
- **权限服务**: `internal/service/permission.go`
- **数据库表**: `casbin_rule` 表存储权限策略

## 权限模型配置

### 内嵌模型配置

权限模型配置已经内嵌到代码中，位于 `internal/auth/casbin.go` 文件的 `getRBACModelText()` 函数：

```go
func getRBACModelText() string {
	return `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act, eft

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && keyMatch2(r.obj, p.obj) && regexMatch(r.act, p.act)
`
}
```

**优势**:
- 无需外部配置文件，简化部署
- 避免配置文件丢失或路径错误
- 模型配置与代码版本同步管理

### 配置说明
- `sub`: 主体（用户或角色）
- `obj`: 对象（资源路径）
- `act`: 动作（HTTP 方法）
- `keyMatch2`: 支持通配符路径匹配
- `regexMatch`: 支持正则表达式方法匹配

## 数据库表结构

### casbin_rule 表

```sql
CREATE TABLE `casbin_rule` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `ptype` varchar(100) NOT NULL COMMENT '策略类型: p-权限策略, g-角色继承',
  `v0` varchar(100) NOT NULL COMMENT '主体(用户/角色)',
  `v1` varchar(100) NOT NULL COMMENT '对象(资源/权限)',
  `v2` varchar(100) NOT NULL COMMENT '动作(操作/HTTP方法)',
  `v3` varchar(100) DEFAULT NULL COMMENT '效果(allow/deny)',
  `v4` varchar(100) DEFAULT NULL COMMENT '扩展字段4',
  `v5` varchar(100) DEFAULT NULL COMMENT '扩展字段5',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_casbin_rule` (`ptype`, `v0`, `v1`, `v2`, `v3`)
);
```

## API 接口说明

### 1. 权限检查

```go
// 检查用户是否有权限访问指定资源
func (s *PermissionService) CheckPermission(userID, resource, action string) (bool, error)

// 检查 RESTful API 权限
func (s *PermissionService) CheckRESTfulPermission(userID, path, method string) (bool, error)
```

### 2. 角色权限管理

```go
// 为角色添加权限
func (s *PermissionService) AddRolePermission(role, resource, action string) error

// 移除角色权限
func (s *PermissionService) RemoveRolePermission(role, resource, action string) error

// 为角色添加 RESTful API 权限
func (s *PermissionService) AddRESTfulPermission(role, path, method string) error

// 移除角色的 RESTful API 权限
func (s *PermissionService) RemoveRESTfulPermission(role, path, method string) error
```

### 3. 权限同步和查询

```go
// 同步角色的 RESTful 权限
func (s *PermissionService) SyncRoleRESTfulPermissions(role string, permissions []RESTfulPermission) error

// 获取角色的 RESTful 权限列表
func (s *PermissionService) GetRoleRESTfulPermissions(role string) ([]RESTfulPermission, error)
```

## 使用示例

### 1. 基本权限检查

```go
// 检查用户是否有访问用户列表的权限
hasPermission, err := permissionService.CheckRESTfulPermission("user1", "/api/users", "GET")
if err != nil {
    log.Printf("权限检查失败: %v", err)
    return
}

if hasPermission {
    fmt.Println("用户有权限访问用户列表")
} else {
    fmt.Println("用户没有权限访问用户列表")
}
```

### 2. 动态添加权限

```go
// 为管理员角色添加用户管理权限
err := permissionService.AddRESTfulPermission("admin", "/api/users/*", "GET|POST|PUT|DELETE")
if err != nil {
    log.Printf("添加权限失败: %v", err)
    return
}

fmt.Println("成功为管理员角色添加用户管理权限")
```

### 3. 批量权限同步

```go
// 定义角色权限
permissions := []service.RESTfulPermission{
    {Path: "/api/users", Method: "GET"},
    {Path: "/api/users/*", Method: "GET"},
    {Path: "/api/users", Method: "POST"},
    {Path: "/api/roles", Method: "GET"},
}

// 同步角色权限
err := permissionService.SyncRoleRESTfulPermissions("admin", permissions)
if err != nil {
    log.Printf("同步权限失败: %v", err)
    return
}

fmt.Println("成功同步管理员角色权限")
```

### 4. 查询角色权限

```go
// 获取角色的所有 RESTful 权限
permissions, err := permissionService.GetRoleRESTfulPermissions("admin")
if err != nil {
    log.Printf("获取权限失败: %v", err)
    return
}

fmt.Printf("管理员角色权限列表:\n")
for _, perm := range permissions {
    fmt.Printf("- %s %s\n", perm.Method, perm.Path)
}
```

## 权限策略示例

### 1. 角色权限策略 (p 策略)

```sql
-- 超级管理员权限
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'super_admin', '/api/*', '*');

-- 管理员权限
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES 
('p', 'admin', '/api/users', 'GET'),
('p', 'admin', '/api/users', 'POST'),
('p', 'admin', '/api/users/*', 'GET|PUT|DELETE');

-- 普通用户权限
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES 
('p', 'user', '/api/users/profile', 'GET|PUT'),
('p', 'user', '/api/departments', 'GET');
```

### 2. 用户角色分配 (g 策略)

```sql
-- 用户角色分配
INSERT INTO casbin_rule (ptype, v0, v1) VALUES 
('g', 'admin', 'super_admin'),
('g', 'user1', 'admin'),
('g', 'user2', 'user');
```

## 中间件集成

### Gin 中间件示例

```go
func AuthMiddleware(permissionService *service.PermissionService) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.GetHeader("User-ID")
        if userID == "" {
            c.JSON(401, gin.H{"error": "未授权访问"})
            c.Abort()
            return
        }

        path := c.Request.URL.Path
        method := c.Request.Method

        hasPermission, err := permissionService.CheckRESTfulPermission(userID, path, method)
        if err != nil {
            c.JSON(500, gin.H{"error": "权限检查失败"})
            c.Abort()
            return
        }

        if !hasPermission {
            c.JSON(403, gin.H{"error": "权限不足"})
            c.Abort()
            return
        }

        c.Next()
    }
}
```

## 最佳实践

### 1. 权限设计原则
- **最小权限原则**: 用户只获得完成工作所需的最小权限
- **职责分离**: 不同角色承担不同的职责和权限
- **权限继承**: 合理使用角色继承减少权限配置复杂度

### 2. 性能优化
- **权限缓存**: 对频繁检查的权限进行缓存
- **批量操作**: 使用批量权限同步减少数据库操作
- **索引优化**: 确保 casbin_rule 表有适当的索引

### 3. 安全建议
- **权限审计**: 定期审计用户权限，及时回收不必要的权限
- **权限日志**: 记录权限变更和访问日志
- **权限测试**: 在权限变更后进行充分测试

## 故障排查

### 1. 常见问题

**问题**: 权限检查总是返回 false
**解决**: 检查权限策略配置和用户角色分配

**问题**: 通配符匹配不生效
**解决**: 确认模型文件中使用了 keyMatch2 函数

**问题**: 权限变更不生效
**解决**: 检查是否调用了 LoadPolicy() 重新加载策略

### 2. 调试方法

```go
// 启用 Casbin 日志
casbin.NewEnforcer("path/to/model.conf", "path/to/policy.csv").EnableLog(true)

// 打印所有策略
policies := enforcer.GetPolicy()
for _, policy := range policies {
    fmt.Printf("Policy: %v\n", policy)
}
```

## 扩展功能

### 1. 自定义匹配函数
可以在模型文件中添加自定义匹配函数来支持更复杂的权限控制逻辑。

### 2. 多租户支持
通过在权限策略中添加租户标识来支持多租户权限隔离。

### 3. 动态权限
结合业务逻辑实现基于数据的动态权限控制。

---

更多详细信息请参考 [Casbin 官方文档](https://casbin.org/docs/zh-CN/overview)。