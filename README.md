# kunpeng

一个基于Go语言的现代化Web应用程序脚手架项目，使用Gin+GORM+Viper+Zap+Lumberjack.v2技术栈。

## 技术栈

- **Web框架**: [Gin](https://github.com/gin-gonic/gin) - 高性能的Go Web框架
- **ORM**: [GORM](https://gorm.io/) - Go语言ORM库
- **配置管理**: [Viper](https://github.com/spf13/viper) - Go应用程序配置解决方案
- **日志**: [Zap](https://github.com/uber-go/zap) + [Lumberjack](https://github.com/natefinch/lumberjack) - 高性能日志库
- **数据库**: MySQL + Redis
- **认证**: JWT (JSON Web Token)
- **API文档**: Swagger
- **容器化**: Docker + Docker Compose

## 项目结构

本项目采用Go官方推荐的项目布局结构：

```
kunpeng/
├── api/                    # API定义文件（OpenAPI/Swagger规范）
├── assets/                 # 其他资源文件
├── build/                  # 打包和持续集成相关文件
├── cmd/                    # 主应用程序入口
│   └── main.go            # 应用程序主入口
├── configs/                # 配置文件模板或默认配置
│   └── config.yaml        # 应用配置文件
├── deployments/            # 系统和容器编排部署配置
│   ├── Dockerfile         # Docker镜像构建文件
│   └── docker-compose.yml # Docker Compose配置
├── docs/                   # 详细文档（配置、数据库、认证、权限等）
├── examples/               # 应用程序或公共库的示例
├── githooks/               # Git钩子
├── init/                   # 系统初始化配置
├── internal/               # 私有应用程序和库代码
│   ├── auth/              # 认证相关
│   │   └── jwt.go         # JWT认证
│   ├── config/            # 配置管理
│   │   └── config.go      # 配置结构和加载
│   ├── database/          # 数据库连接管理
│   │   └── database.go    # 数据库初始化和操作
│   ├── handler/           # HTTP处理器
│   │   └── user.go        # 用户相关API处理
│   ├── logger/            # 日志管理
│   │   └── logger.go      # 日志初始化和配置
│   ├── middleware/        # 中间件
│   │   └── middleware.go  # 各种中间件实现
│   ├── model/             # 数据模型
│   │   └── user.go        # 用户模型
│   ├── redis/             # Redis连接管理
│   │   └── redis.go       # Redis操作封装
│   ├── response/          # 响应处理
│   │   └── response.go    # 统一响应格式
│   ├── router/            # 路由管理
│   │   └── router.go      # 路由配置
│   └── service/           # 业务逻辑层
│       └── user.go        # 用户业务逻辑
├── pkg/                    # 外部应用程序可以使用的库代码
│   └── utils/             # 工具函数
│       └── string.go      # 字符串工具函数
├── scripts/                # 执行各种构建、安装、分析等操作的脚本
│   └── build.sh           # 构建脚本
├── test/                   # 额外的外部测试应用程序和测试数据
│   └── integration_test.go # 集成测试
├── third_party/            # 外部辅助工具、分叉代码和其他第三方工具
├── tools/                  # 项目的支持工具
├── web/                    # Web应用程序特定的组件
├── .air.toml              # Air热重载配置
├── .gitignore             # Git忽略文件配置
├── .golangci.yml          # GolangCI-Lint配置
├── go.mod                 # Go模块文件
├── Makefile               # 构建和开发任务
└── README.md              # 项目说明文档
```

## 主要特性

- **标准目录结构**: 遵循Go官方推荐的项目布局
- **完整的Web框架**: 基于Gin的高性能Web服务
- **数据库支持**: GORM ORM + MySQL + Redis
- **数据库迁移**: 独立的数据库迁移工具，支持自动迁移、重置和删除操作
- **配置管理**: Viper配置管理，支持多种配置格式
- **配置热重载**: 实时监控配置文件变更，自动重新加载配置
- **日志系统**: Zap高性能日志 + Lumberjack日志轮转
- **分级别日志**: 支持按日志级别分文件记录，自动轮转
- **认证授权**: JWT认证 + 角色权限控制
- **中间件**: CORS、日志、恢复、限流、认证等中间件
- **API文档**: Swagger自动生成API文档
- **开发工具**: 热重载、代码检查、格式化等
- **容器化**: Docker + Docker Compose支持
- **测试覆盖**: 单元测试 + 集成测试
- **完整文档**: 提供详细的功能文档和开发指南

## 快速开始

### 环境要求

- Go 1.21+
- MySQL 8.0+
- Redis 6.0+
- Docker (可选)

### 安装依赖

```bash
make deps
```

### 配置文件

复制并修改配置文件：

```bash
cp configs/config.example.yaml configs/config.yaml
# 根据实际环境修改配置
```

#### 配置文件路径覆盖

支持通过以下方式指定配置文件路径：

1. **命令行参数**（优先级最低）：
   ```bash
   ./build/kunpeng -config /path/to/your/config.yaml
   ```

2. **环境变量**（优先级最高）：
   ```bash
   export KUNPENG_CONFIG_PATH=/path/to/your/config.yaml
   ./build/kunpeng
   ```

优先级顺序：环境变量 > 命令行参数 > 默认路径(`configs/config.yaml`)

#### 配置热重载

应用程序支持基于Viper的简单配置文件热重载功能：

- **实时监控**: 基于Viper内置的文件监控机制
- **自动重载**: 配置文件修改后自动重新加载
- **简单回调**: 支持配置变更回调，动态更新日志配置
- **安全性**: 配置不会同步到环境变量，避免敏感信息泄露
- **轻量级**: 移除复杂的ConfigManager，直接使用Viper原生功能

**支持热重载的配置项**:
- 日志配置（级别、格式、输出路径等）
- 应用基本信息（版本、环境等）
- JWT配置（密钥、过期时间等）

**需要重启的配置项**:
- 服务器端口和地址
- 数据库连接参数
- Redis连接参数

**测试热重载功能**:
```bash
# 启动应用
./build/kunpeng -config "configs/config.test.yaml"

# 在另一个终端修改配置文件
sed -i '' 's/level: "debug"/level: "info"/' configs/config.test.yaml

# 观察应用日志，会看到配置重载的信息
```

### 分级别日志

应用程序支持按日志级别分文件记录功能：

- **按级别分文件**: 不同级别的日志记录到不同文件
- **按需创建**: 只有产生对应级别日志时才创建文件
- **自动轮转**: 每个日志文件都支持大小、时间、备份数量控制
- **灵活配置**: 可通过配置开关启用/禁用分级别功能

**日志文件说明**:
- `debug.log` - 仅包含debug级别日志
- `info.log` - 仅包含info级别日志  
- `warn.log` - 仅包含warn级别日志
- `error.log` - 包含error、panic、fatal级别日志

**配置示例**:
```yaml
logging:
  level: "debug"
  format: "json"
  output: "both"          # console/file/both
  separate_files: true    # 启用分级别日志
  log_dir: "logs"         # 日志目录
  max_size: 100           # 每个文件最大100MB
  max_backups: 7          # 保留7个备份文件
  max_age: 30             # 保留30天
  compress: true          # 压缩旧文件
```

**测试分级别日志**:
```bash
# 使用分级别日志配置启动
./build/kunpeng -config "configs/config.separate-logs.yaml"

# 或使用测试脚本
make test-separate-logs
```

## 数据库迁移

项目提供了独立的数据库迁移工具，将数据库迁移操作从应用主程序中分离出来，提供更好的部署和维护体验。

### 迁移工具特性

- **独立运行**: 不依赖主应用程序，可单独执行
- **多种操作**: 支持迁移、重置、删除等操作
- **安全确认**: 危险操作需要用户确认
- **配置灵活**: 支持指定配置文件和环境变量
- **详细日志**: 提供详细的操作日志和错误信息

### 使用方法

#### 1. 使用Makefile命令（推荐）

```bash
# 执行数据库迁移
make migrate

# 使用指定配置文件
make migrate-config CONFIG=configs/config.prod.yaml

# 重置数据库（删除后重新创建）
make migrate-reset

# 删除所有数据库表
make migrate-drop
```

#### 2. 直接使用脚本

```bash
# 基本迁移
./scripts/migrate.sh

# 指定配置文件
./scripts/migrate.sh -c configs/config.prod.yaml

# 重置数据库
./scripts/migrate.sh -a reset

# 删除所有表
./scripts/migrate.sh -a drop

# 查看帮助
./scripts/migrate.sh -h
```

#### 3. 直接运行Go程序

```bash
# 基本迁移
go run cmd/migrate/main.go

# 指定配置和操作
go run cmd/migrate/main.go -config configs/config.yaml -action migrate

# 查看帮助
go run cmd/migrate/main.go -help
```

### 操作说明

- **migrate**: 自动迁移数据库表结构，这是最常用的操作
- **reset**: 删除所有表后重新创建，用于开发环境重置
- **drop**: 仅删除所有表，用于清理数据库

### 环境变量

可以通过环境变量指定配置文件路径：

```bash
export KUNPENG_CONFIG_PATH=configs/config.prod.yaml
./scripts/migrate.sh
```

### 部署建议

1. **生产环境**: 在应用启动前单独执行迁移
   ```bash
   # 部署脚本中
   make migrate-config CONFIG=configs/config.prod.yaml
   # 然后启动应用
   ./build/kunpeng -config configs/config.prod.yaml
   ```

2. **开发环境**: 可以使用重置功能快速重建数据库
   ```bash
   make migrate-reset
   ```

3. **CI/CD**: 在部署流水线中集成迁移步骤
   ```yaml
   # 示例CI配置
   - name: Database Migration
     run: make migrate-config CONFIG=configs/config.prod.yaml
   ```

### 数据库初始化

确保MySQL和Redis服务正在运行，应用启动时会自动创建数据库表。

### 运行项目

#### 开发模式（推荐）

```bash
make dev
```

#### 热重载模式

```bash
make air
```

#### 构建并运行

```bash
make build
make run
```

### API文档

启动服务后，访问 http://localhost:8080/swagger/index.html 查看API文档。

### 健康检查

```bash
curl http://localhost:8080/health/ping
curl http://localhost:8080/health/check
```

## 开发指南

### 项目结构说明

- `/cmd`: 应用程序入口点
- `/internal`: 私有代码，不对外暴露
  - `/config`: 配置管理
  - `/handler`: HTTP请求处理器
  - `/service`: 业务逻辑层
  - `/model`: 数据模型
  - `/middleware`: 中间件
  - `/auth`: 认证相关
  - `/database`: 数据库操作
  - `/redis`: Redis操作
  - `/logger`: 日志管理
  - `/response`: 响应处理
  - `/router`: 路由管理
- `/pkg`: 可对外使用的库代码
- `/configs`: 配置文件
- `/deployments`: 部署相关文件

### 开发工具

```bash
# 安装开发工具
make install-tools

# 代码格式化
make fmt

# 代码检查
make vet
make lint

# 运行测试
make test

# 生成API文档
make swagger
```

### Docker部署

#### 使用Docker Compose（推荐）

```bash
# 启动所有服务
make docker-up

# 停止所有服务
make docker-down
```

#### 单独构建Docker镜像

```bash
make docker-build
make docker-run
```

### API使用示例

#### 用户注册

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123",
    "nickname": "Test User"
  }'
```

#### 用户登录

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
```

#### 获取用户信息（需要JWT Token）

```bash
curl -X GET http://localhost:8080/api/v1/user/profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## 详细文档

项目提供了完整的文档体系，涵盖各个功能模块的详细说明：

### 核心功能文档
- [配置管理](docs/CONFIG.md) - 详细的配置系统说明和示例
- [数据库指南](docs/DATABASE_GUIDE.md) - 多数据库支持、驱动安装和配置
- [数据库迁移](docs/MIGRATION.md) - 数据库迁移系统和最佳实践
- [JWT认证](docs/JWT_AUTH.md) - JWT认证机制和使用方法
- [RBAC权限控制](docs/RBAC_RESTful_Permission_Guide.md) - 基于角色的访问控制指南
- [日志系统](docs/LOGGING.md) - 日志配置和管理
- [Redis缓存](docs/REDIS.md) - Redis配置和使用

### 开发文档
- [处理器和服务](docs/HANDLER_SERVICE.md) - HTTP处理器和业务服务层
- [数据模型](docs/MODELS.md) - 数据模型定义和关系
- [中间件](docs/MIDDLEWARE.md) - 中间件系统和自定义中间件
- [路由管理](docs/ROUTER.md) - 路由配置和管理
- [响应处理](docs/RESPONSE.md) - 统一响应格式和错误处理
- [工具函数](docs/UTILS.md) - 通用工具函数库

### 系统文档
- [管理系统](docs/ADMIN_SYSTEM.md) - 后台管理系统功能

## 配置说明

主要配置项说明：

```yaml
app:
  name: "kunpeng"           # 应用名称
  version: "1.0.0"          # 版本号
  environment: "development" # 环境：development/production
  debug: true               # 调试模式

server:
  host: "0.0.0.0"          # 服务器地址
  port: 8080               # 端口号
  read_timeout: 30s        # 读取超时
  write_timeout: 30s       # 写入超时
  idle_timeout: 60s        # 空闲超时

database:
  driver: "mysql"          # 数据库驱动
  host: "localhost"        # 数据库地址
  port: 3306              # 数据库端口
  username: "root"         # 用户名
  password: "password"     # 密码
  database: "kunpeng"      # 数据库名

redis:
  host: "localhost"        # Redis地址
  port: 6379              # Redis端口
  password: ""            # Redis密码
  database: 0             # Redis数据库

jwt:
  secret: "your-secret"    # JWT密钥（生产环境请修改）
  expire_time: 24h        # Token过期时间
  issuer: "kunpeng"        # 签发者

logging:
  level: "info"           # 日志级别
  format: "json"          # 日志格式
  output: "file"          # 输出方式
```

## 贡献

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

## 许可证

本项目采用 MIT 许可证。详情请参阅 [LICENSE](LICENSE) 文件。

## 致谢

本项目参考了以下优秀的Go项目脚手架：

- [eagle](https://github.com/go-eagle/eagle) - 一个基于Go的微服务框架
- [nunu](https://github.com/go-nunu/nunu) - 一个Go应用脚手架

感谢这些项目提供的灵感和最佳实践。

## 目录说明

### `/cmd`
本项目的主干。每个应用程序的目录名应该与你想要的可执行文件的名称相匹配。

### `/internal`
私有应用程序和库代码。这是你不希望其他人在其应用程序或库中导入代码。

### `/pkg`
外部应用程序可以使用的库代码。其他项目会导入这些库，所以在这里放东西之前要三思。

### `/api`
OpenAPI/Swagger规范，JSON模式文件，协议定义文件。

### `/web`
Web应用程序特定的组件：静态Web资产，服务器端模板和SPAs。

### `/configs`
配置文件模板或默认配置。

### `/scripts`
执行各种构建，安装，分析等操作的脚本。

### `/build`
打包和持续集成。将你的云（AMI），容器（Docker），操作系统（deb，rpm，pkg）包配置和脚本放在/build/package目录中。

### `/deployments`
IaaS，PaaS，系统和容器编排部署配置和模板（docker-compose，kubernetes/helm，mesos，terraform，bosh）。

### `/test`
额外的外部测试应用程序和测试数据。你可以随时根据需求构造/test目录。

### `/docs`
详细的功能文档和开发指南，包括配置管理、数据库、认证、权限控制等各个模块的完整说明。

### `/examples`
你的应用程序和/或公共库的示例。

### `/third_party`
外部辅助工具，分叉代码和其他第三方工具（例如Swagger UI）。

### `/githooks`
Git钩子。

### `/assets`
与存储库一起使用的其他资产（图像，徽标等）。

### `/tools`
这个项目的支持工具。注意，这些工具可以从/pkg和/internal目录导入代码。

### `/init`
系统初始化（systemd，upstart，sysv）和进程管理器（runit，supervisord）配置。

## 开始使用

1. 在`cmd/`目录下创建你的应用程序入口点
2. 在`internal/`目录下放置私有代码
3. 在`pkg/`目录下放置可重用的公共库代码
4. 根据需要使用其他目录

## 构建和运行

```bash
# 构建应用程序
go build -o bin/app cmd/main.go

# 运行应用程序
./bin/app
```