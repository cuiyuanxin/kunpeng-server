# 环境自适应中间件配置指南

本文档说明如何在不同环境下使用不同的恢复和日志中间件，参考《Go程序编程之旅》的实现方式。

## 概述

项目现在支持根据环境配置自动选择合适的中间件：
- **开发环境**：使用 Gin 框架自带的 `gin.Recovery()` 和 `gin.Logger()` 中间件
- **生产环境**：使用自定义的 `CustomRecovery` 和 `CustomLogger` 中间件

## 实现原理

### 1. 环境判断

在 `internal/router/router.go` 的 `setupMiddleware` 方法中，通过读取配置文件中的 `app.environment` 字段来判断当前运行环境：

```go
func (r *Router) setupMiddleware() {
    // 根据环境选择不同的恢复和日志中间件
    if r.config.App.Environment == "production" {
        // 生产环境使用自定义中间件
        r.engine.Use(middleware.CustomRecovery(logger.Logger, func(c *gin.Context, msg string) {
            response.InternalServerError(c, msg)
        }))
        r.engine.Use(middleware.CustomLogger(logger.Logger))
    } else {
        // 开发环境使用gin框架自带的中间件
        r.engine.Use(gin.Recovery())
        r.engine.Use(gin.Logger())
    }
    
    // 其他中间件...
}
```

### 2. 自定义中间件特性

#### CustomLogger 中间件

生产环境的自定义日志中间件提供以下增强功能：

- **请求ID追踪**：自动生成或使用现有的请求ID
- **详细的请求信息记录**：包括请求开始和结束时间
- **智能日志级别**：根据HTTP状态码自动选择日志级别
  - 5xx 错误：Error 级别
  - 4xx 错误：Warn 级别
  - 其他：Info 级别
- **丰富的上下文信息**：记录响应体大小、用户代理、引用页等

#### CustomRecovery 中间件

生产环境的自定义恢复中间件提供以下增强功能：

- **详细的panic信息记录**：包括堆栈跟踪信息
- **请求上下文保存**：记录发生panic时的完整请求信息
- **结构化错误日志**：使用zap结构化日志格式
- **时间戳记录**：精确记录panic发生时间

## 配置文件设置

### 开发环境配置 (config.dev.yaml)

```yaml
app:
  name: "kunpeng-server"
  version: "1.0.0"
  environment: "development"  # 开发环境标识
  debug: true
```

### 生产环境配置 (config.prod.yaml)

```yaml
app:
  name: "kunpeng-server"
  version: "1.0.0"
  environment: "production"   # 生产环境标识
  debug: false
```

## 使用方法

### 1. 启动开发环境

```bash
# 使用开发环境配置
go run cmd/main.go -config configs/config.dev.yaml
```

开发环境下的日志输出示例：
```
[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.
[GIN-debug] GET    /health/ping              --> main.main.func1 (5 handlers)
[GIN] 2024/01/15 - 14:30:25 | 200 |      1.2345ms |       127.0.0.1 | GET      "/health/ping"
```

### 2. 启动生产环境

```bash
# 使用生产环境配置
go run cmd/main.go -config configs/config.prod.yaml
```

生产环境下的日志输出示例：
```json
{
  "level": "info",
  "timestamp": "2024-01-15T14:30:25.123Z",
  "caller": "middleware/common.go:85",
  "msg": "Request started",
  "request_id": "20240115143025-random",
  "method": "GET",
  "path": "/health/ping",
  "client_ip": "127.0.0.1",
  "user_agent": "curl/7.68.0"
}
```

## 优势对比

### 开发环境优势

1. **简单直观**：使用Gin内置中间件，输出格式简洁易读
2. **快速调试**：彩色输出，便于开发时快速定位问题
3. **零配置**：无需额外配置，开箱即用

### 生产环境优势

1. **结构化日志**：JSON格式便于日志收集和分析
2. **详细追踪**：请求ID支持分布式链路追踪
3. **智能分级**：根据响应状态自动调整日志级别
4. **完整上下文**：记录更多请求和响应信息
5. **错误恢复**：详细的panic信息和堆栈跟踪

## 扩展配置

### 自定义环境判断

如果需要支持更多环境（如测试环境），可以修改判断逻辑：

```go
func (r *Router) setupMiddleware() {
    switch r.config.App.Environment {
    case "production":
        // 生产环境中间件
        r.engine.Use(middleware.CustomRecovery(logger.Logger, func(c *gin.Context, msg string) {
            response.InternalServerError(c, msg)
        }))
        r.engine.Use(middleware.CustomLogger(logger.Logger))
    case "testing":
        // 测试环境中间件
        r.engine.Use(middleware.CustomRecovery(logger.Logger, func(c *gin.Context, msg string) {
            response.InternalServerError(c, msg)
        }))
        r.engine.Use(gin.Logger()) // 测试环境使用简单日志
    default:
        // 开发环境中间件
        r.engine.Use(gin.Recovery())
        r.engine.Use(gin.Logger())
    }
}
```

### 中间件参数配置

可以通过配置文件控制中间件的行为：

```yaml
middleware:
  custom_logger:
    enable_request_body: true
    enable_response_body: false
    max_body_size: 1024
  custom_recovery:
    enable_stack_trace: true
    stack_size: 4096
```

## 注意事项

1. **性能考虑**：生产环境的自定义中间件会记录更多信息，可能对性能有轻微影响
2. **日志存储**：生产环境建议配置日志轮转和归档策略
3. **敏感信息**：注意不要在日志中记录敏感信息（如密码、token等）
4. **监控集成**：生产环境建议集成APM工具进行性能监控

## 参考资料

- 《Go程序编程之旅》- 中间件设计模式
- [Gin框架官方文档](https://gin-gonic.com/docs/)
- [Zap日志库文档](https://pkg.go.dev/go.uber.org/zap)