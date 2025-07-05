.PHONY: build clean test help run dev deps swagger docker-build docker-run lint fmt vet mod-tidy

# 应用名称
APP_NAME := kunpeng

# 构建目录
BUILD_DIR := build

# 版本信息
VERSION := $(shell git describe --tags --always --dirty)
BUILD_TIME := $(shell date +%Y-%m-%d_%H:%M:%S)
GIT_COMMIT := $(shell git rev-parse HEAD)

# Go 构建标志
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)"

# 默认目标
all: build

# 安装依赖
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

# 构建应用
build: deps
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME) cmd/main.go
	@echo "Build completed: $(BUILD_DIR)/$(APP_NAME)"

# 运行应用
run: build
	@echo "Running $(APP_NAME)..."
	@./$(BUILD_DIR)/$(APP_NAME)

# 开发模式运行
dev:
	@echo "Running in development mode..."
	@go run cmd/main.go

# 清理构建文件
clean:
	@echo "Cleaning build files..."
	@rm -rf $(BUILD_DIR)
	@go clean
	@echo "Clean completed"

# 运行测试
test:
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

# 代码格式化
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# 代码检查
vet:
	@echo "Running go vet..."
	@go vet ./...

# 代码静态分析
lint:
	@echo "Running golangci-lint..."
	@golangci-lint run

# 整理模块
mod-tidy:
	@echo "Tidying modules..."
	@go mod tidy

# 生成Swagger文档
swagger:
	@echo "Generating Swagger docs..."
	@swag init -g cmd/main.go -o docs

# Docker构建
docker-build:
	@echo "Building Docker image..."
	@docker build -t $(APP_NAME):$(VERSION) -f deployments/Dockerfile .

# Docker运行
docker-run:
	@echo "Running Docker container..."
	@docker run -p 8081:8080 --name $(APP_NAME) $(APP_NAME):$(VERSION)

# Docker Compose启动
docker-up:
	@echo "Starting services with Docker Compose..."
	@docker-compose -f deployments/docker-compose.yml up -d

# Docker Compose停止
docker-down:
	@echo "Stopping services with Docker Compose..."
	@docker-compose -f deployments/docker-compose.yml down

# 安装开发工具
install-tools:
	@echo "Installing development tools..."
	@go install github.com/swaggo/swag/cmd/swag@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/air-verse/air@latest

# 热重载开发
air:
	@echo "Starting with hot reload..."
	@air

# 使用自定义配置运行
run-config:
	@echo "Running with custom config..."
	@./$(BUILD_DIR)/$(APP_NAME) -config=$(CONFIG_PATH)

# 使用环境变量配置运行
run-env:
	@echo "Running with environment config..."
	@KUNPENG_CONFIG_PATH=$(CONFIG_PATH) ./$(BUILD_DIR)/$(APP_NAME)

# 测试配置热重载功能
test-reload:
	@echo "Testing config hot reload..."
	@echo "Starting application with test config..."
	@echo "You can modify configs/config.test.yaml in another terminal to test hot reload"
	./build/kunpeng -config=configs/config.test.yaml

# 测试分级别日志功能
test-separate-logs:
	@echo "Testing separate log files by level..."
	@./scripts/test-separate-logs.sh

# 数据库迁移
migrate:
	@echo "Running database migration..."
	@./scripts/migrate.sh

# 数据库迁移 - 使用指定配置
migrate-config:
	@echo "Running database migration with custom config..."
	@./scripts/migrate.sh -c $(CONFIG)

# 重置数据库
migrate-reset:
	@echo "Resetting database..."
	@./scripts/migrate.sh -a reset

# 删除所有数据库表
migrate-drop:
	@echo "Dropping all database tables..."
	@./scripts/migrate.sh -a drop

# 测试Docker配置
test-docker:
	@echo "Testing Docker configuration..."
	@./scripts/test-docker.sh

# Docker部署（带构建）
docker-deploy:
	@echo "Deploying with Docker Compose..."
	@docker-compose -f deployments/docker-compose.yml up --build -d

# 查看Docker服务状态
docker-status:
	@echo "Docker services status:"
	@docker-compose -f deployments/docker-compose.yml ps

# 查看Docker服务日志
docker-logs:
	@echo "Docker services logs:"
	@docker-compose -f deployments/docker-compose.yml logs -f

# 显示帮助信息
help:
	@echo "Available targets:"
	@echo "  deps              - Install dependencies"
	@echo "  build             - Build the application"
	@echo "  run               - Build and run the application"
	@echo "  dev               - Run in development mode"
	@echo "  clean             - Clean build files"
	@echo "  test              - Run tests with coverage"
	@echo "  fmt               - Format code"
	@echo "  vet               - Run go vet"
	@echo "  lint              - Run golangci-lint"
	@echo "  mod-tidy          - Tidy modules"
	@echo "  swagger           - Generate Swagger docs"
	@echo "  docker-build      - Build Docker image"
	@echo "  docker-run        - Run Docker container"
	@echo "  docker-up         - Start with Docker Compose"
	@echo "  docker-down       - Stop Docker Compose services"
	@echo "  install-tools     - Install development tools"
	@echo "  air               - Start with hot reload"
	@echo "  run-config        - Run with custom config (CONFIG_PATH=path/to/config.yaml)"
	@echo "  run-env           - Run with environment config (CONFIG_PATH=path/to/config.yaml)"
	@echo "  test-reload       - Test config hot reload functionality"
	@echo "  test-separate-logs- Test separate log files by level"
	@echo "  migrate            - 执行数据库迁移"
	@echo "  migrate-config     - 使用指定配置执行迁移 (CONFIG=path)"
	@echo "  migrate-reset      - 重置数据库"
	@echo "  migrate-drop       - 删除所有数据库表"
	@echo "  test-docker        - 测试Docker配置"
	@echo "  docker-deploy      - Docker部署（带构建）"
	@echo "  docker-status      - 查看Docker服务状态"
	@echo "  docker-logs        - 查看Docker服务日志"
	@echo "  help              - Show this help message"