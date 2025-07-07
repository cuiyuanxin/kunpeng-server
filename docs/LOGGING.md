# 日志系统使用指南

本项目采用基于 Zap 的高性能日志系统，支持环境自适应、GORM 集成和扩展钩子功能。

## 功能特性

### 1. 环境自适应配置

日志系统可以根据运行环境自动调整输出方式：

- **开发环境** (`development`, `dev`, `test`)：默认输出到控制台，使用 console 格式
- **生产环境** (`production`, `prod`)：默认输出到文件，使用 JSON 格式

### 2. GORM 日志集成

支持 GORM ORM 的日志记录：

- SQL 查询日志记录
- 慢查询检测和记录
- 数据库错误日志
- 支持分离 SQL 日志和错误日志到不同文件

### 3. 扩展钩子系统

为未来功能预留钩子：

- **Tracing 钩子**：用于链路追踪集成（OpenTelemetry、Jaeger 等）
- **gRPC 钩子**：用于 gRPC 服务日志记录
- **Custom 钩子**：用于自定义日志处理逻辑

## 配置说明

### 基础日志配置

```yaml
logging:
  level: "info"              # 日志级别：debug, info, warn, error
  format: "json"             # 日志格式：json, console
  output: "file"             # 输出方式：file, stdout
  file_path: "logs/app.log"  # 日志文件路径
  max_size: 100              # 单个日志文件最大大小(MB)
  max_backups: 3             # 保留的日志文件数量
  max_age: 28                # 日志文件保留天数
  compress: true             # 是否压缩旧日志文件
```

### 环境自适应配置

```yaml
logging:
  auto_mode: true            # 启用环境自适应模式
  force_console: null        # 强制使用控制台输出 (true/false/null)
  force_file: null           # 强制使用文件输出 (true/false/null)
```

**配置优先级**：
1. `force_console` 或 `force_file` 强制配置（最高优先级）
2. `auto_mode` 环境自适应配置
3. 默认配置

### GORM 日志配置

```yaml
logging:
  gorm:
    enabled: true                    # 是否启用 GORM 日志
    level: "info"                    # GORM 日志级别：silent, error, warn, info
    slow_threshold: "200ms"          # 慢查询阈值
    sql_file: "logs/sql.log"         # SQL 日志文件路径
    error_file: "logs/gorm_error.log" # GORM 错误日志文件路径
    auto_mode: true                  # 启用环境自适应
    force_console: null              # 强制使用控制台输出
    force_file: null                 # 强制使用文件输出
```

### 钩子配置

```yaml
logging:
  hooks:
    tracing:                         # 链路追踪钩子
      enabled: false                 # 是否启用
      level: "info"                  # 日志级别
      format: "json"                 # 日志格式
      output: "file"                 # 输出方式
      file_path: "logs/tracing.log"  # 日志文件路径
      auto_mode: true                # 环境自适应
      force_console: null
      force_file: null
    grpc:                            # gRPC 钩子
      enabled: false
      level: "info"
      format: "json"
      output: "file"
      file_path: "logs/grpc.log"
      auto_mode: true
      force_console: null
      force_file: null
    custom:                          # 自定义钩子
      enabled: false
      level: "info"
      format: "json"
      output: "file"
      file_path: "logs/custom.log"
      auto_mode: true
      force_console: null
      force_file: null
```

## 使用示例

### 基础日志记录

```go
import (
    klogger "github.com/kunpeng-server/internal/logger"
    "go.uber.org/zap"
)

// 基础日志
klogger.Info("用户登录成功", zap.String("user_id", "123"))
klogger.Error("数据库连接失败", zap.Error(err))

// 带上下文的日志
klogger.Info("处理请求",
    klogger.WithRequestID("req-123"),
    klogger.WithUserID("user-456"),
    zap.String("action", "create_order"),
)
```

### 链路追踪日志

```go
// 添加链路追踪信息
klogger.Info("处理业务逻辑",
    klogger.WithTraceID("trace-123"),
    klogger.WithSpanID("span-456"),
    zap.String("operation", "user_service"),
)
```

### gRPC 日志

```go
// 记录 gRPC 调用
klogger.Info("gRPC 调用完成",
    klogger.WithGRPCMethod("/user.UserService/GetUser"),
    klogger.WithGRPCCode(0),
    zap.Duration("duration", time.Since(start)),
)
```

### GORM 日志

GORM 日志会自动记录，无需手动调用：

```go
// GORM 操作会自动记录日志
var user User
db.Where("id = ?", userID).First(&user)
// 自动记录：SQL 语句、执行时间、影响行数等
```

## 环境配置示例

### 开发环境 (config.dev.yaml)

```yaml
app:
  environment: "development"

logging:
  level: "debug"
  format: "console"    # 便于阅读
  auto_mode: true      # 自动输出到控制台
  gorm:
    enabled: true
    level: "info"       # 显示所有 SQL
    slow_threshold: "100ms"  # 更严格的慢查询检测
  hooks:
    tracing:
      enabled: true     # 开发环境启用追踪
```

### 生产环境 (config.prod.yaml)

```yaml
app:
  environment: "production"

logging:
  level: "info"
  format: "json"       # 便于日志收集
  auto_mode: true      # 自动输出到文件
  max_backups: 10      # 保留更多备份
  max_age: 30          # 保留更长时间
  gorm:
    enabled: true
    level: "warn"       # 只记录警告和错误
    slow_threshold: "500ms"  # 更宽松的慢查询阈值
  hooks:
    tracing:
      enabled: true     # 生产环境启用追踪
    grpc:
      enabled: true     # 生产环境启用 gRPC 日志
```

## 钩子扩展开发

### 注册自定义钩子

```go
import (
    klogger "github.com/kunpeng-server/internal/logger"
    "go.uber.org/zap/zapcore"
)

// 定义钩子函数
func myCustomHook(entry zapcore.Entry, fields []zapcore.Field) error {
    // 自定义日志处理逻辑
    // 例如：发送到外部监控系统、格式化特殊字段等
    return nil
}

// 注册钩子
klogger.RegisterHook("my_custom_hook", myCustomHook)

// 注销钩子
klogger.UnregisterHook("my_custom_hook")
```

### 链路追踪集成示例

```go
// 集成 OpenTelemetry
func tracingHook(entry zapcore.Entry, fields []zapcore.Field) error {
    span := trace.SpanFromContext(ctx)
    if span.IsRecording() {
        // 将日志信息添加到 span
        span.AddEvent(entry.Message, trace.WithAttributes(
            attribute.String("level", entry.Level.String()),
            attribute.String("logger", entry.LoggerName),
        ))
    }
    return nil
}
```

## 性能优化建议

1. **生产环境**：使用 `info` 或更高级别，避免 `debug` 级别
2. **文件轮转**：合理设置 `max_size`、`max_backups` 和 `max_age`
3. **异步日志**：Zap 默认支持异步写入，无需额外配置
4. **结构化日志**：使用 `zap.Field` 而不是字符串拼接
5. **避免频繁日志**：在高频操作中适当降低日志级别

## 故障排查

### 常见问题

1. **日志文件无法创建**：检查目录权限和磁盘空间
2. **GORM 日志不显示**：确认 `gorm.enabled` 为 `true`
3. **环境自适应不生效**：检查 `app.environment` 配置
4. **钩子不工作**：确认钩子已正确注册且 `enabled` 为 `true`

### 调试方法

```go
// 检查当前日志配置
klogger.Info("当前环境", zap.String("env", cfg.App.Environment))

// 测试 GORM 日志
gormLogger := klogger.GetGormLogger()
if gormLogger != nil {
    klogger.Info("GORM 日志器已初始化")
} else {
    klogger.Warn("GORM 日志器未初始化")
}
```

## 最佳实践

1. **统一日志格式**：使用结构化字段，便于日志分析
2. **合理分级**：根据重要性选择合适的日志级别
3. **上下文信息**：记录请求 ID、用户 ID 等关键信息
4. **敏感信息**：避免记录密码、令牌等敏感数据
5. **性能监控**：利用慢查询日志优化数据库性能
6. **日志聚合**：在生产环境中使用 ELK、Prometheus 等工具收集分析日志