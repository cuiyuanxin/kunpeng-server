#!/bin/bash

# 数据库迁移脚本

set -e

# 脚本目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

# 默认配置
CONFIG_PATH="configs/config.yaml"
ACTION="migrate"

# 显示帮助信息
show_help() {
    echo "数据库迁移脚本"
    echo ""
    echo "用法:"
    echo "  $0 [选项]"
    echo ""
    echo "选项:"
    echo "  -c, --config PATH     配置文件路径 (默认: configs/config.yaml)"
    echo "  -a, --action ACTION   执行的操作 (默认: migrate)"
    echo "                        可选值: migrate, drop, reset"
    echo "  -h, --help           显示此帮助信息"
    echo ""
    echo "操作说明:"
    echo "  migrate  - 自动迁移数据库表"
    echo "  drop     - 删除所有数据库表 (危险操作)"
    echo "  reset    - 重置数据库 (删除后重新创建)"
    echo ""
    echo "环境变量:"
    echo "  KUNPENG_CONFIG_PATH  配置文件路径 (优先级最高)"
    echo ""
    echo "示例:"
    echo "  # 执行数据库迁移"
    echo "  $0"
    echo ""
    echo "  # 使用指定配置文件"
    echo "  $0 -c configs/config.prod.yaml"
    echo ""
    echo "  # 重置数据库"
    echo "  $0 -a reset"
    echo ""
    echo "  # 删除所有表"
    echo "  $0 -a drop"
}

# 解析命令行参数
while [[ $# -gt 0 ]]; do
    case $1 in
        -c|--config)
            CONFIG_PATH="$2"
            shift 2
            ;;
        -a|--action)
            ACTION="$2"
            shift 2
            ;;
        -h|--help)
            show_help
            exit 0
            ;;
        *)
            echo "未知选项: $1"
            show_help
            exit 1
            ;;
    esac
done

# 验证操作类型
case $ACTION in
    migrate|drop|reset)
        ;;
    *)
        echo "错误: 无效的操作类型 '$ACTION'"
        echo "有效的操作类型: migrate, drop, reset"
        exit 1
        ;;
esac

# 切换到项目目录
cd "$PROJECT_DIR"

# 检查配置文件是否存在
if [[ -n "$KUNPENG_CONFIG_PATH" ]]; then
    CONFIG_PATH="$KUNPENG_CONFIG_PATH"
fi

if [[ ! -f "$CONFIG_PATH" ]]; then
    echo "错误: 配置文件不存在: $CONFIG_PATH"
    exit 1
fi

echo "=== 数据库迁移工具 ==="
echo "配置文件: $CONFIG_PATH"
echo "执行操作: $ACTION"
echo ""

# 确认危险操作
if [[ "$ACTION" == "drop" || "$ACTION" == "reset" ]]; then
    echo "警告: 此操作将删除数据库中的所有表和数据!"
    read -p "确定要继续吗? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "操作已取消"
        exit 0
    fi
fi

# 执行迁移
echo "开始执行数据库迁移..."
go run cmd/migrate/main.go -config "$CONFIG_PATH" -action "$ACTION"

if [[ $? -eq 0 ]]; then
    echo ""
    echo "✅ 数据库迁移完成!"
else
    echo ""
    echo "❌ 数据库迁移失败!"
    exit 1
fi