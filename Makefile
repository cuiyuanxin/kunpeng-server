# 鲲鹏后台管理系统 Makefile

# 变量定义
APP_NAME=kunpeng
GO=go
GOPATH=$(shell go env GOPATH)
GOBIN=$(shell go env GOBIN)
ifeq ($(GOBIN),)
	GOBIN=$(GOPATH)/bin
endif
SWAG=$(GOBIN)/swag
AIR=air
GOFMT=gofmt -w
MAIN_FILE=cmd/server/main.go
BUILD_DIR=build
DOCKER_IMAGE=$(APP_NAME)
DOCKER_TAG=latest

# 版本信息
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME = $(shell date '+%Y-%m-%d %H:%M:%S')
GIT_COMMIT = $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# 构建标志
LDFLAGS = -ldflags "-X 'main.Version=$(VERSION)' -X 'main.BuildTime=$(BUILD_TIME)' -X 'main.GitCommit=$(GIT_COMMIT)'"

# 默认目标
.PHONY: all
all: clean build

# 安装依赖
.PHONY: deps
deps:
	@echo "安装依赖..."
	$(GO) mod tidy
	@if ! command -v $(SWAG) > /dev/null; then \
		echo "安装 swag..."; \
		$(GO) install github.com/swaggo/swag/cmd/swag@latest; \
	fi
	@echo "安装 air..."; \
	$(GO) install github.com/air-verse/air@latest; \

# 格式化代码
.PHONY: fmt
fmt:
	@echo "格式化代码..."
	$(GOFMT) .

# 生成Swagger文档
.PHONY: swagger
swagger:
	@echo "生成Swagger文档..."
	$(SWAG) init -g $(MAIN_FILE) -o ./docs

# 构建应用
.PHONY: build
build: fmt swagger
	@echo "构建应用..."
	$(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_FILE)

# 运行应用
.PHONY: run
run: build
	@echo "运行应用..."
	./$(BUILD_DIR)/$(APP_NAME)

# 开发模式运行（热重载）
.PHONY: dev
dev:
	@echo "开发模式运行..."
	$(AIR)

# 清理构建产物
.PHONY: clean
clean:
	@echo "清理构建产物..."
	rm -rf $(BUILD_DIR)
	rm -rf ./docs

# 测试
.PHONY: test
test:
	@echo "运行测试..."
	$(GO) test -v ./...

# 构建Docker镜像
.PHONY: docker-build
docker-build:
	@echo "构建Docker镜像..."
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) -f build/docker/Dockerfile .

# 运行Docker容器
.PHONY: docker-run
docker-run:
	@echo "运行Docker容器..."
	docker run -d -p 8080:8080 --name $(APP_NAME) $(DOCKER_IMAGE):$(DOCKER_TAG)

# 停止Docker容器
.PHONY: docker-stop
docker-stop:
	@echo "停止Docker容器..."
	docker stop $(APP_NAME)
	docker rm $(APP_NAME)

# 帮助信息
.PHONY: help
help:
	@echo "鲲鹏后台管理系统 Makefile 帮助"
	@echo ""
	@echo "可用命令:"
	@echo "  make deps         - 安装依赖"
	@echo "  make fmt          - 格式化代码"
	@echo "  make swagger      - 生成Swagger文档"
	@echo "  make build        - 构建应用"
	@echo "  make run          - 运行应用"
	@echo "  make dev          - 开发模式运行（热重载）"
	@echo "  make clean        - 清理构建产物"
	@echo "  make test         - 运行测试"
	@echo "  make docker-build - 构建Docker镜像"
	@echo "  make docker-run   - 运行Docker容器"
	@echo "  make docker-stop  - 停止Docker容器"
	@echo "  make help         - 显示帮助信息"