#!/bin/bash

# Kunpeng 构建脚本

set -e

APP_NAME="kunpeng"
BUILD_DIR="bin"
MAIN_FILE="cmd/main.go"

echo "开始构建 $APP_NAME..."

# 创建构建目录
mkdir -p $BUILD_DIR

# 获取版本信息
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(git rev-parse HEAD 2>/dev/null || echo "unknown")

echo "版本: $VERSION"
echo "构建时间: $BUILD_TIME"
echo "Git提交: $GIT_COMMIT"

# 设置构建标志
LDFLAGS="-X main.Version=$VERSION -X main.BuildTime=$BUILD_TIME -X main.GitCommit=$GIT_COMMIT"

# 构建应用程序
echo "正在编译..."
go build -ldflags "$LDFLAGS" -o $BUILD_DIR/$APP_NAME $MAIN_FILE

echo "构建完成: $BUILD_DIR/$APP_NAME"

# 显示文件信息
ls -la $BUILD_DIR/$APP_NAME

echo "构建成功!"