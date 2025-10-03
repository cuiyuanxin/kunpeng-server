# i18n 与统一错误包集成指南

## 概述

本项目已成功将 i18n（国际化）功能直接集成到统一错误处理包中。所有现有的错误创建和处理方法现在都自动支持多语言，无需修改任何现有代码。

## 核心功能

### 1. 自动国际化支持
- 所有现有的错误创建方法（`New`, `NewWithMessage`, `NewWithData`）现在都自动支持 i18n
- 自动根据当前语言环境返回对应的错误消息
- 如果没有找到翻译，自动回退到默认消息

### 2. 错误消息翻译
- `Error()` 方法现在自动返回国际化的错误字符串
- 支持动态语言切换
- 完全向后兼容，无需修改现有代码

## 在控制器中使用

### 基本用法（无需修改现有代码）
```go
// 现有代码自动支持 i18n
err := errors.New(20001, originalErr)
response.FailWithMessage(err.Error(), c)

// 或者直接使用
response.FailWithError(errors.New(20001, originalErr), c)

// 带自定义消息
err := errors.NewWithMessage(20001, "自定义消息", originalErr)

// 带附加数据
err := errors.NewWithData(20001, originalErr, userData)
```

### 语言切换示例
```go
// 设置为中文
i18n.SetLanguage("zh")
err := errors.New(20001, nil)
fmt.Println(err.Error()) // 输出: 错误码: 20001, 错误信息: 用户不存在

// 切换到英文
i18n.SetLanguage("en")
err2 := errors.New(20001, nil)
fmt.Println(err2.Error()) // 输出: 错误码: 20001, 错误信息: User not found
```

## 支持的错误码

### 系统级错误码 (10000-19999)
- 10000: 成功
- 10001: 未知错误
- 10002: 数据库错误
- 10003: 缓存错误
- 10004: 认证失败
- 10005: 权限不足
- 10006: 文件操作失败

### 业务级错误码 (20000-29999)
- 20001: 用户不存在
- 20002: 用户已存在
- 20003: 密码错误
- 20004: 角色不存在
- 20005: 菜单不存在
- 20006: 部门不存在
- 20007: 岗位不存在
- 20008: 字典不存在

## 语言文件配置

### 中文 (zh.yaml)
```yaml
error:
  10001: "未知错误"
  10002: "数据库错误"
  20001: "用户不存在"
  20002: "用户已存在"
```

### 英文 (en.yaml)
```yaml
error:
  10001: "Unknown error"
  10002: "Database error"
  20001: "User not found"
  20002: "User already exists"
```

## API 使用示例

### 使用现有方法（自动支持 i18n）
```go
// 基本用法
err := errors.New(20001, nil)

// 带原始错误
originalErr := fmt.Errorf("database connection failed")
err := errors.New(10002, originalErr)

// 带自定义消息
err := errors.NewWithMessage(20001, "用户名已存在", nil)

// 带附加数据
err := errors.NewWithData(20001, nil, map[string]string{"field": "username"})

// 在 HTTP 响应中使用
func GetUser(c *gin.Context) {
    // ... 业务逻辑
    if userNotFound {
        err := errors.New(20001, nil)
        response.FailWithMessage(err.Error(), c)
        return
    }
    // ...
}
```

## 最佳实践

1. **无缝迁移**: 现有代码无需修改，自动获得 i18n 支持
2. **保持错误码一致性**: 确保错误码在不同语言文件中都有对应的翻译
3. **合理的回退机制**: 当翻译缺失时，系统会自动使用默认消息
4. **性能考虑**: i18n 查找有轻微的性能开销，但在实际应用中可以忽略
5. **测试覆盖**: 为不同语言环境编写测试用例
6. **自定义消息优先**: 在 `NewWithMessage` 中，自定义消息优先于 i18n 翻译

## 兼容性

- **完全向后兼容**: 现有的 `New`、`NewWithMessage`、`NewWithData` 方法现在自动支持 i18n
- **零代码修改**: 无需修改任何现有代码即可获得 i18n 支持
- **透明集成**: i18n 功能完全透明，不影响现有的错误处理逻辑

## 性能考虑

- i18n 查找操作有轻微的性能开销
- 建议在高并发场景下进行性能测试
- 可以考虑添加缓存机制来优化性能

## 扩展支持

### 添加新语言
1. 在 `pkg/i18n/locales/` 目录下创建新的语言文件
2. 添加对应的错误码翻译
3. 重新编译项目

### 添加新错误码
1. 在 `pkg/errors/code.go` 中定义新的错误码
2. 在所有语言文件中添加对应的翻译
3. 确保错误码的一致性

## 测试验证

项目包含了完整的测试用例来验证 i18n 集成功能：

```bash
# 运行错误包测试
go test ./pkg/errors -v

# 运行 i18n 包测试
go test ./pkg/i18n -v
```

## 总结

通过将 i18n 功能直接集成到错误处理包中，我们实现了：

1. **零侵入性**: 现有代码无需任何修改
2. **自动化**: 错误消息自动根据当前语言环境进行翻译
3. **向后兼容**: 完全保持与现有 API 的兼容性
4. **灵活性**: 支持动态语言切换和自定义消息
5. **可扩展性**: 易于添加新语言和新错误码

这种设计确保了项目的国际化需求得到满足，同时保持了代码的简洁性和可维护性。