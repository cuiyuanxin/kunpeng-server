# 代码结构重构说明

本文档说明了项目代码结构的重构结果，包括 `pkg` 和 `internal` 目录的划分，以及 API 版本控制的实现。

## 目录结构划分

### pkg 目录（可复用的公共库）

`pkg` 目录包含可以被其他项目复用的通用工具和库：

```
pkg/
├── utils/           # 通用工具函数
│   └── string.go    # 字符串处理工具
├── response/        # HTTP 响应处理
│   └── response.go  # 统一响应格式
├── auth/           # 认证相关
│   └── jwt.go      # JWT 令牌管理
└── middleware/     # 通用中间件
    └── common.go   # 通用 Gin 中间件
```

#### pkg/utils
- **用途**: 通用的工具函数，如字符串处理、数据转换等
- **特点**: 无业务逻辑依赖，纯函数式工具
- **示例**: `IsEmpty()`, `Capitalize()`, `Reverse()`, `ContainsIgnoreCase()`

#### pkg/response
- **用途**: 统一的 HTTP 响应格式处理
- **特点**: 标准化的 API 响应结构
- **功能**: 
  - 统一响应格式 (`Response`, `PageResponse`)
  - 常用状态码响应函数 (`Success`, `BadRequest`, `Unauthorized` 等)
  - 支持自定义响应和分页响应

#### pkg/auth
- **用途**: JWT 认证相关的通用功能
- **特点**: 可配置的 JWT 管理器
- **功能**:
  - JWT 令牌生成和解析
  - 令牌刷新和验证
  - 支持自定义配置

#### pkg/middleware
- **用途**: 通用的 Gin 中间件
- **特点**: 与具体业务无关的中间件
- **功能**:
  - CORS 跨域处理
  - 请求日志记录
  - 错误恢复
  - 请求 ID 生成
  - 超时控制
  - 限流控制

### internal 目录（项目私有代码）

`internal` 目录包含项目特定的业务逻辑和实现：

```
internal/
├── config/         # 配置管理
├── model/          # 数据模型
├── service/        # 业务逻辑层
├── handler/        # HTTP 处理器
│   ├── v1/         # V1 版本 API
│   └── v2/         # V2 版本 API
├── router/         # 路由配置
├── database/       # 数据库相关
├── logger/         # 日志配置
├── middleware/     # 项目特定中间件
├── auth/           # 认证模块（兼容层）
└── response/       # 响应模块（兼容层）
```

#### 项目特定组件

- **config/**: 项目配置结构和加载逻辑
- **model/**: 数据库模型和业务实体
- **service/**: 业务逻辑实现，包含具体的业务规则
- **database/**: 数据库连接和迁移
- **logger/**: 项目特定的日志配置

#### 中间件分层

- **pkg/middleware**: 通用中间件（CORS、日志、恢复等）
- **internal/middleware**: 项目特定中间件（JWT 认证、权限验证等）

## API 版本控制

### 版本控制策略

项目采用 URL 路径版本控制策略：

```
/api/v1/...  # V1 版本 API
/api/v2/...  # V2 版本 API
```

### 版本化处理器

#### V1 版本 (`internal/handler/v1/`)

- **特点**: 基础功能实现
- **路由**: `/api/v1/admin/*`
- **功能**:
  - 基本的 CRUD 操作
  - 标准的认证和授权
  - 简单的响应格式

**示例路由**:
```
POST   /api/v1/admin/auth/login
GET    /api/v1/admin/users
POST   /api/v1/admin/users
GET    /api/v1/admin/roles
```

#### V2 版本 (`internal/handler/v2/`)

- **特点**: 增强功能和新特性
- **路由**: `/api/v2/admin/*`
- **新增功能**:
  - 批量操作支持
  - 数据分析和统计
  - 增强的响应格式
  - 更详细的错误处理
  - 缓存支持

**示例路由**:
```
POST   /api/v2/admin/auth/login      # 增强的登录（设备信息记录）
GET    /api/v2/admin/users          # 高级筛选和排序
PUT    /api/v2/admin/users/batch    # 批量操作
GET    /api/v2/admin/users/analytics # 数据分析
```

### V2 版本增强特性

#### 1. 增强的响应格式
```json
{
  "code": 200,
  "message": "获取成功",
  "data": [...],
  "version": "v2",
  "meta": {
    "total_pages": 10,
    "has_next": true,
    "has_previous": false,
    "request_time": 1640995200,
    "response_time": 1640995200123
  }
}
```

#### 2. 批量操作
```json
{
  "user_ids": [1, 2, 3],
  "action": "enable",
  "data": {}
}
```

#### 3. 数据分析
```json
{
  "period": "7d",
  "total_users": {"count": 100, "growth": "+5.2%"},
  "active_users": {"count": 85, "growth": "+3.1%"},
  "user_distribution": {...},
  "login_trends": [...]
}
```

## 向后兼容性

为了保持向后兼容性，原有的 `internal/auth` 和 `internal/response` 模块被保留作为兼容层：

```go
// internal/auth/jwt.go
package auth

import "github.com/cuiyuanxin/kunpeng/pkg/auth"

// 重新导出pkg/auth中的类型和函数
type Claims = auth.Claims
type JWTManager = auth.JWTManager

var NewJWTManager = auth.NewJWTManager
```

这样现有代码无需修改即可继续工作。

## 使用建议

### 1. 新功能开发
- 优先使用 `pkg` 目录下的通用组件
- 业务逻辑放在 `internal` 目录下
- 新的 API 建议使用最新版本

### 2. 版本选择
- **V1**: 适用于基础功能和简单需求
- **V2**: 适用于需要高级功能的场景（批量操作、数据分析等）

### 3. 迁移策略
- 逐步将现有代码迁移到新的结构
- 保持 API 版本的向后兼容性
- 新功能优先在新版本中实现

## 代码质量提升

### 1. 模块化设计
- 清晰的职责分离
- 可复用的组件设计
- 降低耦合度

### 2. 可维护性
- 统一的代码结构
- 标准化的响应格式
- 完善的错误处理

### 3. 可扩展性
- 版本化的 API 设计
- 插件化的中间件
- 灵活的配置管理

### 4. 测试友好
- 纯函数式的工具库
- 依赖注入的设计
- 模块化的组件

## 总结

通过这次重构，项目实现了：

1. **清晰的代码组织**: `pkg` 和 `internal` 目录的合理划分
2. **版本控制支持**: V1/V2 API 版本的并行支持
3. **代码复用**: 通用组件的提取和复用
4. **向后兼容**: 现有代码的平滑迁移
5. **可维护性**: 标准化的结构和规范

这为项目的长期发展和维护奠定了良好的基础。