#!/bin/bash

# Docker 部署测试脚本

set -e

# 脚本目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

echo "=== Docker 部署测试 ==="
echo "项目目录: $PROJECT_DIR"
echo ""

# 切换到项目目录
cd "$PROJECT_DIR"

# 检查 Docker 是否运行
echo "检查 Docker 状态..."
if ! docker info >/dev/null 2>&1; then
    echo "❌ Docker 未运行，请启动 Docker Desktop"
    exit 1
fi
echo "✅ Docker 运行正常"

# 检查配置文件
echo "检查配置文件..."
if [[ ! -f "deployments/Dockerfile" ]]; then
    echo "❌ Dockerfile 不存在"
    exit 1
fi

if [[ ! -f "deployments/docker-compose.yml" ]]; then
    echo "❌ docker-compose.yml 不存在"
    exit 1
fi
echo "✅ 配置文件存在"

# 检查必要目录
echo "检查必要目录..."
if [[ ! -d "configs" ]]; then
    echo "❌ configs 目录不存在"
    exit 1
fi

if [[ ! -d "web/static" ]]; then
    echo "❌ web/static 目录不存在"
    exit 1
fi
echo "✅ 必要目录存在"

# 验证 docker-compose 配置
echo "验证 docker-compose 配置..."
if docker-compose -f deployments/docker-compose.yml config >/dev/null 2>&1; then
    echo "✅ docker-compose 配置有效"
else
    echo "❌ docker-compose 配置无效"
    docker-compose -f deployments/docker-compose.yml config
    exit 1
fi

# 检查网络连接
echo "检查网络连接..."
if curl -s --connect-timeout 5 https://registry-1.docker.io >/dev/null 2>&1; then
    echo "✅ Docker Hub 连接正常"
    NETWORK_OK=true
else
    echo "⚠️  Docker Hub 连接失败，可能需要配置镜像源"
    NETWORK_OK=false
fi

# 显示建议
echo ""
echo "=== 部署建议 ==="
if [[ "$NETWORK_OK" == "true" ]]; then
    echo "1. 网络连接正常，可以直接运行:"
    echo "   docker-compose -f deployments/docker-compose.yml up --build -d"
else
    echo "1. 网络连接问题，建议配置 Docker 镜像源:"
    echo "   - 打开 Docker Desktop 设置"
    echo "   - 在 Docker Engine 中添加镜像源配置"
    echo "   - 或使用本地已有镜像"
fi

echo ""
echo "2. 如果遇到权限问题，确保:"
echo "   - Docker Desktop 有足够权限"
echo "   - 当前用户在 docker 组中"

echo ""
echo "3. 服务端口映射:"
echo "   - 应用服务: http://localhost:8081"
echo "   - Nginx: http://localhost:80"
echo "   - MySQL: localhost:3306"
echo "   - Redis: localhost:6379"

echo ""
echo "4. 查看服务状态:"
echo "   docker-compose -f deployments/docker-compose.yml ps"

echo ""
echo "5. 查看服务日志:"
echo "   docker-compose -f deployments/docker-compose.yml logs -f [service_name]"

echo ""
echo "✅ Docker 配置检查完成!"