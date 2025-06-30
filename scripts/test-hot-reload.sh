#!/bin/bash

# 配置热重载测试脚本
# 此脚本用于演示配置文件热重载功能

set -e

echo "=== Kunpeng 配置热重载测试 ==="
echo

# 检查是否已构建应用
if [ ! -f "build/kunpeng" ]; then
    echo "应用未构建，正在构建..."
    make build
    echo
fi

# 备份原始测试配置
cp configs/config.test.yaml configs/config.test.yaml.backup

echo "1. 启动应用（使用测试配置）..."
echo "   应用将在后台运行，PID将保存到 /tmp/kunpeng-test.pid"
echo

# 启动应用
./build/kunpeng -config=configs/config.test.yaml > /tmp/kunpeng-test.log 2>&1 &
APP_PID=$!
echo $APP_PID > /tmp/kunpeng-test.pid

echo "应用已启动，PID: $APP_PID"
echo "日志文件: /tmp/kunpeng-test.log"
echo

# 等待应用启动
echo "等待应用启动..."
sleep 3

# 检查应用是否正在运行
if ! kill -0 $APP_PID 2>/dev/null; then
    echo "错误: 应用启动失败"
    cat /tmp/kunpeng-test.log
    exit 1
fi

echo "2. 应用已成功启动，现在测试配置热重载..."
echo

# 修改配置文件
echo "修改应用版本号从 1.0.1 到 1.0.2..."
sed -i '' 's/version: "1.0.1"/version: "1.0.2"/' configs/config.test.yaml

echo "等待配置重载..."
sleep 2

echo "修改日志级别从 debug 到 info..."
sed -i '' 's/level: "debug"/level: "info"/' configs/config.test.yaml

echo "等待配置重载..."
sleep 2

echo "修改JWT密钥..."
sed -i '' 's/secret: "test-secret-key-for-hot-reload"/secret: "new-secret-key-after-reload"/' configs/config.test.yaml

echo "等待配置重载..."
sleep 2

echo
echo "3. 配置修改完成，查看应用日志中的配置重载信息:"
echo "========================================"
tail -20 /tmp/kunpeng-test.log | grep -E "(Config|config|Configuration|reload)"
echo "========================================"
echo

echo "4. 清理测试环境..."

# 停止应用
if kill -0 $APP_PID 2>/dev/null; then
    echo "停止应用 (PID: $APP_PID)..."
    kill $APP_PID
    sleep 2
    
    # 如果应用仍在运行，强制停止
    if kill -0 $APP_PID 2>/dev/null; then
        echo "强制停止应用..."
        kill -9 $APP_PID
    fi
fi

# 恢复原始配置
echo "恢复原始配置文件..."
mv configs/config.test.yaml.backup configs/config.test.yaml

# 清理临时文件
rm -f /tmp/kunpeng-test.pid

echo
echo "=== 测试完成 ==="
echo "完整的应用日志保存在: /tmp/kunpeng-test.log"
echo "你可以查看该文件了解配置热重载的详细过程"
echo