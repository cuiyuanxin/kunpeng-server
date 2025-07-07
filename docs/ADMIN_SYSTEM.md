# 后台管理系统文档

## 概述

本项目实现了完整的后台管理系统，支持管理员用户管理、角色权限控制、部门管理、操作日志记录等功能。系统采用 RBAC（基于角色的访问控制）模型，提供灵活的权限管理机制。

## 系统架构

### 核心组件

```
后台管理系统
├── 数据模型层 (Model)
│   ├── AdminUser - 管理员用户
│   ├── AdminRole - 角色
│   ├── AdminPermission - 权限
│   ├── AdminDepartment - 部门
│   ├── AdminLoginLog - 登录日志
│   ├── AdminOperationLog - 操作日志
│   └── AdminConfig - 系统配置
├── 服务层 (Service)
│   ├── AdminUserService - 用户服务
│   ├── AdminRoleService - 角色服务
│   ├── AdminPermissionService - 权限服务
│   └── AdminDepartmentService - 部门服务
├── 处理器层 (Handler)
│   └── AdminHandler - 管理员处理器
└── 路由层 (Router)
    └── AdminRoutes - 管理员路由
```

## 数据库设计

### 核心表结构

#### 1. 管理员用户表 (admin_users)
```sql
CREATE TABLE admin_users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(50) UNIQUE NOT NULL COMMENT '用户名',
    phone VARCHAR(20) UNIQUE COMMENT '手机号',
    email VARCHAR(100) UNIQUE COMMENT '邮箱',
    password VARCHAR(255) NOT NULL COMMENT '密码',
    real_name VARCHAR(50) COMMENT '真实姓名',
    avatar VARCHAR(255) COMMENT '头像',
    status TINYINT DEFAULT 1 COMMENT '状态：1-正常，2-禁用',
    department_id BIGINT COMMENT '部门ID',
    last_login_at TIMESTAMP COMMENT '最后登录时间',
    last_login_ip VARCHAR(45) COMMENT '最后登录IP',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

#### 2. 角色表 (admin_roles)
```sql
CREATE TABLE admin_roles (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(50) UNIQUE NOT NULL COMMENT '角色名称',
    code VARCHAR(50) UNIQUE NOT NULL COMMENT '角色代码',
    description TEXT COMMENT '角色描述',
    status TINYINT DEFAULT 1 COMMENT '状态：1-正常，2-禁用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

#### 3. 权限表 (admin_permissions)
```sql
CREATE TABLE admin_permissions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL COMMENT '权限名称',
    code VARCHAR(100) UNIQUE NOT NULL COMMENT '权限代码',
    type TINYINT NOT NULL COMMENT '类型：1-菜单，2-按钮，3-接口',
    parent_id BIGINT DEFAULT 0 COMMENT '父级ID',
    path VARCHAR(255) COMMENT '路由路径',
    component VARCHAR(255) COMMENT '组件路径',
    icon VARCHAR(100) COMMENT '图标',
    sort_order INT DEFAULT 0 COMMENT '排序',
    status TINYINT DEFAULT 1 COMMENT '状态：1-正常，2-禁用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

### RBAC 权限模型

```
用户 (AdminUser) ←→ 用户角色关联 (admin_user_roles) ←→ 角色 (AdminRole)
                                                           ↓
                                                    角色权限关联 (admin_role_permissions)
                                                           ↓
                                                    权限 (AdminPermission)
```

## API 接口

### 认证接口

#### 登录
```http
POST /api/v1/admin/auth/login
Content-Type: application/json

{
    "account": "admin",     // 用户名、手机号或邮箱
    "password": "123456",
    "captcha": "1234",     // 验证码（可选）
    "captcha_id": "xxx"    // 验证码ID（可选）
}
```

#### 登出
```http
POST /api/v1/admin/auth/logout
Authorization: Bearer <token>
```

#### 获取个人信息
```http
GET /api/v1/admin/auth/profile
Authorization: Bearer <token>
```

### 用户管理接口

#### 创建用户
```http
POST /api/v1/admin/users
Authorization: Bearer <token>
Content-Type: application/json

{
    "username": "testuser",
    "phone": "13800138000",
    "email": "test@example.com",
    "password": "123456",
    "real_name": "测试用户",
    "department_id": 1,
    "role_ids": [1, 2]
}
```

#### 获取用户列表
```http
GET /api/v1/admin/users?page=1&page_size=10&keyword=test&status=1&department_id=1
Authorization: Bearer <token>
```

#### 更新用户
```http
PUT /api/v1/admin/users/:id
Authorization: Bearer <token>
Content-Type: application/json

{
    "real_name": "新的真实姓名",
    "department_id": 2,
    "status": 1
}
```

### 角色管理接口

#### 创建角色
```http
POST /api/v1/admin/roles
Authorization: Bearer <token>
Content-Type: application/json

{
    "name": "编辑员",
    "code": "editor",
    "description": "内容编辑员角色",
    "permission_ids": [1, 2, 3]
}
```

#### 获取角色列表
```http
GET /api/v1/admin/roles?page=1&page_size=10
Authorization: Bearer <token>
```

### 权限管理接口

#### 获取权限列表
```http
GET /api/v1/admin/permissions
Authorization: Bearer <token>
```

#### 获取权限树
```http
GET /api/v1/admin/permissions/tree
Authorization: Bearer <token>
```

### 部门管理接口

#### 获取部门列表
```http
GET /api/v1/admin/departments
Authorization: Bearer <token>
```

#### 获取部门树
```http
GET /api/v1/admin/departments/tree
Authorization: Bearer <token>
```

## 使用示例

### 1. 初始化管理员账户

```sql
-- 插入超级管理员
INSERT INTO admin_users (username, phone, email, password, real_name, status) 
VALUES ('admin', '13800138000', 'admin@example.com', '$2a$10$...', '超级管理员', 1);

-- 插入超级管理员角色
INSERT INTO admin_roles (name, code, description) 
VALUES ('超级管理员', 'super_admin', '拥有所有权限的超级管理员');

-- 关联用户和角色
INSERT INTO admin_user_roles (user_id, role_id) VALUES (1, 1);
```

### 2. 登录获取 Token

```bash
curl -X POST http://localhost:8080/api/v1/admin/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "account": "admin",
    "password": "123456"
  }'
```

### 3. 使用 Token 访问受保护的接口

```bash
curl -X GET http://localhost:8080/api/v1/admin/auth/profile \
  -H "Authorization: Bearer <your-token>"
```

## 权限控制

### 中间件

系统使用以下中间件进行权限控制：

1. **JWT 认证中间件**: 验证用户身份
2. **权限检查中间件**: 检查用户是否有访问特定资源的权限
3. **角色检查中间件**: 检查用户是否具有特定角色

### 权限检查流程

```
请求 → JWT验证 → 获取用户信息 → 获取用户角色 → 获取角色权限 → 权限匹配 → 允许/拒绝访问
```

## 安全特性

### 1. 密码安全
- 使用 bcrypt 加密存储密码
- 支持密码强度验证
- 支持密码过期策略

### 2. 登录安全
- 支持验证码验证
- 登录失败次数限制
- 异常登录检测
- 登录日志记录

### 3. 操作审计
- 记录所有管理员操作
- 包含操作时间、IP、用户、操作内容等信息
- 支持操作日志查询和导出

### 4. 会话管理
- JWT Token 过期时间控制
- 支持 Token 刷新
- 支持强制下线

## 配置说明

### JWT 配置

```yaml
jwt:
  secret: "your-secret-key"
  expire_time: 7200  # 2小时
  refresh_time: 86400  # 24小时
```

### 管理员配置

```yaml
admin:
  default_password: "123456"  # 默认密码
  password_expire_days: 90    # 密码过期天数
  max_login_attempts: 5       # 最大登录尝试次数
  lockout_duration: 1800      # 锁定时长（秒）
```

## 部署说明

### 1. 数据库初始化

```bash
# 执行 SQL 文件创建表结构
mysql -u root -p kunpeng < sql/admin_tables.sql
```

### 2. 运行数据库迁移

```bash
# 自动创建表结构
go run cmd/main.go migrate
```

### 3. 创建初始管理员

```bash
# 使用命令行工具创建
go run cmd/main.go admin:create --username=admin --password=123456
```

## 最佳实践

### 1. 权限设计
- 遵循最小权限原则
- 合理设计权限粒度
- 定期审查权限分配

### 2. 安全建议
- 定期更换 JWT 密钥
- 启用 HTTPS
- 配置防火墙规则
- 定期备份数据

### 3. 监控告警
- 监控异常登录
- 监控权限变更
- 监控系统性能

## 相关文档

- [JWT 认证系统文档](JWT_AUTH.md)
- [中间件系统文档](MIDDLEWARE.md)
- [数据模型文档](MODELS.md)
- [路由管理文档](ROUTER.md)
- [统一响应格式文档](RESPONSE.md)

---

**注意**: 在生产环境中，请确保修改默认密码、配置强密码策略、启用 HTTPS 等安全措施。