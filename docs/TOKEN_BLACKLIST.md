# Token 黑名单功能说明

## 功能概述

Token黑名单功能是一个安全机制，用于管理已失效的JWT令牌。系统采用双token机制（Access Token + Refresh Token），当用户退出登录时，系统会将其token加入黑名单，确保该token无法再次被使用，从而提高系统安全性。

### 双Token机制

- **Access Token**: 用于API访问认证，有效期较短（默认2小时，记住我7天）
- **Refresh Token**: 用于刷新Access Token，有效期较长（默认14小时，记住我14天）
- **Token类型验证**: JWT中间件只接受Access Token类型的令牌进行API访问

## 功能特性

1. **自动黑名单管理**：用户退出登录时自动将 token 加入黑名单
2. **中间件验证**：JWT 中间件会检查 token 是否在黑名单中
3. **过期清理**：定时任务自动清理过期的黑名单记录
4. **完整日志**：记录退出登录和黑名单操作的详细日志

## 数据库表结构

### kp_token_blacklist 表

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | bigint | 主键ID |
| token | text | JWT token |
| user_id | bigint | 用户ID |
| username | varchar(50) | 用户名 |
| reason | varchar(255) | 加入黑名单原因 |
| expires_at | datetime | token过期时间 |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |
| deleted_at | datetime | 删除时间 |

## API 接口

### 用户登录
- **接口**: `POST /api/v1/login`
- **功能**: 用户登录并获取JWT token对
- **返回**: 包含access_token、refresh_token和过期时间的响应
- **响应字段**:
  - `access_token`: 访问令牌
  - `refresh_token`: 刷新令牌
  - `expires_in`: access_token过期时间（秒）
  - `refresh_expires_in`: refresh_token过期时间（秒）

### 刷新Token
- **接口**: `POST /api/v1/refresh-token`
- **功能**: 使用refresh token刷新获取新的token对
- **请求体**: `{"refresh_token": "your_refresh_token"}`
- **返回**: 新的token对（格式同登录接口）

### 退出登录

**接口地址**：`POST /api/v1/logout`

**请求头**：
```
Authorization: Bearer <your_jwt_token>
```

**响应示例**：
```json
{
  "code": 200,
  "message": "退出登录成功",
  "data": null
}
```

**功能说明**：
- 获取请求头中的 JWT token
- 将 token 加入黑名单
- 记录退出登录日志
- 返回成功响应

## 中间件验证流程

1. 从请求头获取 Authorization 字段
2. 解析 Bearer token
3. 验证 token 格式和有效性
4. **检查 token 是否在黑名单中**
5. 如果在黑名单中，返回 "无效的令牌" 错误
6. 如果不在黑名单中，继续正常流程

## 定时清理任务

系统提供了 `TokenCleanupTask` 定时任务，用于清理过期的黑名单记录：

- **执行时间**：每天凌晨 2 点
- **清理规则**：删除已过期的 token 记录
- **日志记录**：记录清理过程和结果

### 启动清理任务

```go
// 在应用启动时添加以下代码
cleanupTask := task.NewTokenCleanupTask()
cleanupTask.StartCleanupScheduler()
```

## 安全优势

1. **防止 token 重放攻击**：退出登录的 token 无法再次使用
2. **提高安全性**：即使 token 被泄露，退出登录后也无法使用
3. **完整审计**：所有黑名单操作都有详细日志记录
4. **自动清理**：避免数据库存储过多无用数据

## 注意事项

1. **性能考虑**：每次请求都会查询黑名单，建议对黑名单表建立适当索引
2. **存储空间**：token 是 text 类型，会占用较多存储空间
3. **清理策略**：定时任务只清理过期记录，可根据需要调整清理策略
4. **错误处理**：黑名单添加失败不会影响退出登录流程

## 扩展功能

可以基于此功能扩展以下特性：

1. **管理员强制下线**：管理员可以将指定用户的 token 加入黑名单
2. **批量下线**：支持批量将用户 token 加入黑名单
3. **黑名单查询**：提供接口查询用户的黑名单 token
4. **白名单机制**：为特殊场景提供 token 白名单功能