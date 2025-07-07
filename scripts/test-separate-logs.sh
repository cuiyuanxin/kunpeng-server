#!/bin/bash

# 分级别日志功能测试脚本

set -e

echo "=== 分级别日志功能测试 ==="

# 清理之前的日志文件
echo "清理之前的日志文件..."
rm -rf logs/
mkdir -p logs

# 构建应用
echo "构建应用..."
make build

# 启动应用（后台运行）
echo "启动应用（使用分级别日志配置）..."
./build/kunpeng -config "configs/config.separate-logs.yaml" &
APP_PID=$!

echo "应用已启动，PID: $APP_PID"
echo "等待应用初始化..."
sleep 3

# 检查日志文件是否创建
echo "\n=== 检查日志文件创建情况 ==="
ls -la logs/ || echo "日志目录不存在"

# 发送一些测试请求来生成不同级别的日志
echo "\n=== 发送测试请求生成日志 ==="

# 测试正常请求（会产生info日志）
echo "发送正常请求..."
curl -s http://localhost:8080/health > /dev/null || echo "健康检查请求失败（这是正常的，因为路由可能未定义）"

# 测试错误请求（会产生error日志）
echo "发送错误请求..."
curl -s http://localhost:8080/nonexistent > /dev/null || echo "错误请求已发送"

# 等待日志写入
echo "等待日志写入..."
sleep 2

# 检查生成的日志文件
echo "\n=== 检查生成的日志文件 ==="
echo "日志目录内容:"
ls -la logs/ 2>/dev/null || echo "没有日志文件生成"

echo "\n=== 查看各级别日志内容 ==="

# 查看debug日志
if [ -f "logs/debug.log" ]; then
    echo "\n--- Debug日志内容 ---"
    tail -5 logs/debug.log
else
    echo "\n--- Debug日志文件不存在 ---"
fi

# 查看info日志
if [ -f "logs/info.log" ]; then
    echo "\n--- Info日志内容 ---"
    tail -5 logs/info.log
else
    echo "\n--- Info日志文件不存在 ---"
fi

# 查看warn日志
if [ -f "logs/warn.log" ]; then
    echo "\n--- Warn日志内容 ---"
    tail -5 logs/warn.log
else
    echo "\n--- Warn日志文件不存在 ---"
fi

# 查看error日志
if [ -f "logs/error.log" ]; then
    echo "\n--- Error日志内容 ---"
    tail -5 logs/error.log
else
    echo "\n--- Error日志文件不存在 ---"
fi

# 停止应用
echo "\n=== 停止应用 ==="
kill $APP_PID 2>/dev/null || echo "应用已停止"
wait $APP_PID 2>/dev/null || true

echo "\n=== 测试完成 ==="
echo "\n说明:"
echo "1. 只有产生了对应级别的日志时，才会创建相应的日志文件"
echo "2. debug.log - 包含debug级别的日志"
echo "3. info.log - 包含info级别的日志"
echo "4. warn.log - 包含warn级别的日志"
echo "5. error.log - 包含error、panic、fatal级别的日志"
echo "6. 所有日志文件都支持自动轮转（根据大小、时间、备份数量）"
echo "\n可以通过修改 configs/config.separate-logs.yaml 中的 separate_files 为 false 来禁用分级别日志"