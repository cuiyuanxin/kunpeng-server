# 鲲鹏后台管理系统

鲲鹏是一个通用型的前后端分离模式的admin后台管理系统，基于Golang实现，采用企业级目录结构设计。

## 技术栈

- 框架：Gin
- 配置管理：Viper (支持热更新)
- 数据库ORM：GORM + MySQL 8
- 日志：Zap + Lumberjack
- 参数验证：go-playground/validator
- 权限管理：Casbin (RBAC + RESTful)
- 认证：JWT
- API文档：Swagger
- 包管理：Go Modules

## 项目结构

```
kunpeng/
├── cmd/                    # 应用程序入口
│   └── server/             # 服务器入口
├── configs/                # 配置文件
├── internal/               # 私有应用和库代码
│   ├── app/                # 应用服务
│   ├── controller/         # 控制器层
│   ├── interfaces/         # 接口定义层
│   │   ├── repository/     # 数据访问接口
│   │   └── service/        # 业务服务接口
│   ├── middleware/         # 中间件
│   ├── model/              # 数据模型
│   │   └── dto/            # 数据传输对象
│   ├── repository/         # 数据访问层
│   │   └── impl/           # 仓储实现
│   └── service/            # 业务逻辑层
│       └── impl/           # 服务实现
├── pkg/                    # 公共库代码
│   ├── casbin/             # 权限管理
│   ├── config/             # 配置处理
│   ├── constants/          # 常量定义
│   ├── database/           # 数据库连接
│   ├── errors/             # 错误码定义
│   ├── i18n/               # 国际化支持
│   ├── jwt/                # JWT处理
│   ├── logger/             # 日志处理
│   ├── response/           # 响应处理
│   ├── tracer/             # 链路追踪
│   ├── utils/              # 工具函数
│   └── validator/          # 参数验证
├── scripts/                # 脚本文件
│   └── mysql/              # 数据库脚本
├── docs/                   # 文档
├── api/                    # API定义
├── test/                   # 测试文件
├── build/                  # 构建相关
│   └── docker/             # Docker配置
├── deploy/                 # 部署配置
├── logs/                   # 日志文件
├── tmp/                    # 临时文件
├── Makefile                # 构建脚本
├── go.mod                  # 依赖管理
└── go.sum                  # 依赖校验
```

## 功能特性

### 架构设计
1. **分层架构**：采用DDD分层架构，接口与实现分离
2. **依赖注入**：通过接口定义实现依赖倒置，提高可测试性
3. **模块化设计**：清晰的目录结构，便于维护和扩展
4. **避免循环依赖**：统一的接口管理，解决模块间依赖问题

### 核心功能
1. **配置管理**：使用Viper实现配置热更新
2. **日志系统**：使用Zap+Lumberjack实现日志分级和分割
3. **数据库处理**：使用GORM连接MySQL，支持扩展接口
4. **错误码标准化**：实现了系统级和业务级错误码
5. **统一响应处理**：标准化API响应格式
6. **参数验证**：使用validator实现参数验证，支持国际化
7. **权限管理**：使用Casbin实现RBAC+RESTful权限控制
8. **JWT认证**：支持签发、验证、续签，实现二次加密
9. **IP获取优化**：支持代理环境下的真实IP获取
10. **登录安全**：实现登录失败拉黑机制，防止暴力破解

### 中间件系统
- **JWT验证中间件**：统一的身份认证
- **访问日志中间件**：记录请求响应日志
- **异常捕获中间件**：全局异常处理
- **限流中间件**：API访问频率控制
- **超时控制中间件**：请求超时保护
- **链路追踪中间件**：分布式链路追踪
- **操作日志中间件**：业务操作审计

### 业务模块
- **用户管理**：用户CRUD、状态管理、密码重置
- **角色管理**：角色权限分配、菜单授权、API授权
- **菜单管理**：动态菜单树、权限控制
- **API管理**：接口资源管理、权限绑定
- **部门管理**：组织架构树形管理
- **岗位管理**：职位信息维护
- **字典管理**：系统字典数据维护
- **日志管理**：登录日志、操作日志查询

## 快速开始

### 环境要求

- Go 1.18+
- MySQL 8.0+

### 本地开发

1. 克隆项目

```bash
git clone https://github.com/cuiyuanxin/kunpeng.git
cd kunpeng
```

2. 安装依赖

```bash
go mod tidy
```

3. 修改配置文件

```bash
cp configs/config.yaml.example configs/config.yaml
# 编辑 configs/config.yaml 文件，修改数据库连接信息等
```

4. 运行项目

```bash
# 直接运行
go run cmd/server/main.go

# 或者使用Makefile
make run

# 开发模式（热重载）
make dev
```

### Docker部署

1. 构建Docker镜像

```bash
make docker-build
```

2. 运行Docker容器

```bash
make docker-run
```

或者使用Docker Compose

```bash
docker-compose up -d
```

## 使用Makefile

项目提供了Makefile来简化常用操作：

```bash
# 安装依赖
make deps

# 格式化代码
make fmt

# 生成Swagger文档
make swagger

# 构建应用
make build

# 运行应用
make run

# 开发模式运行（热重载）
make dev

# 清理构建产物
make clean

# 测试
make test

# 构建Docker镜像
make docker-build

# 运行Docker容器
make docker-run

# 停止Docker容器
make docker-stop

# 显示帮助信息
make help
```

## 项目配置

配置文件位于`configs/config.yaml`，支持以下配置项：

- 应用配置：端口、环境、超时时间等
- 数据库配置：连接信息、连接池设置等
- 日志配置：日志级别、输出路径、分割设置等
- JWT配置：密钥、过期时间、签发者等
- Casbin配置：模型路径等

## 功能备忘录
- [ ] 实现单点登录功能，A浏览器登录后，B浏览器再登录，A浏览器的账户会自动被踢出登录状态
- [ ✅ ] jwt token，在单点登录切换登录后需要把之前的token失效
- [ ] 增加手机号+验证码登录
- [ ] 增加短信验证平台+模版管理
- [ ] 增加图形验证码验证插件
- [ ] 增加第三方登录授权


## 贡献指南

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建Pull Request

## 许可证

[MIT](LICENSE)