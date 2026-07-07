#!/bin/bash

# 最直接的纠偏测试 - 不使用 go test，直接编译运行程序

set -e

echo "=========================================="
echo "直接纠偏测试（不使用单元测试）"
echo "=========================================="
echo ""

# 检查参数
if [ -z "$1" ]; then
    echo "用法: $0 <图片路径>"
    exit 1
fi

IMAGE_PATH="$1"

if [ ! -f "$IMAGE_PATH" ]; then
    echo "✗ 文件不存在: $IMAGE_PATH"
    exit 1
fi

echo "📄 输入图片: $IMAGE_PATH"
echo ""

# 设置环境变量
export CGO_ENABLED=1

# 获取图片所在目录和文件名
IMAGE_DIR=$(dirname "$IMAGE_PATH")
FILENAME=$(basename "$IMAGE_PATH")

echo "🔨 编译程序..."
cd /app
go build -o /tmp/scanrectifier .

echo ""
echo "🚀 执行纠偏..."
echo ""

# 运行 deskew 命令
/tmp/scanrectifier deskew --dir "$IMAGE_DIR"

echo ""
echo "📊 检查结果..."
echo ""

# 检查输出文件
OUTPUT_FILE="$IMAGE_DIR/deskewed/$FILENAME"
if [ -f "$OUTPUT_FILE" ]; then
    echo "✅ 纠偏成功！"
    echo ""
    echo "   输入文件: $IMAGE_PATH"
    echo "   输出文件: $OUTPUT_FILE"
    echo ""
    echo "   文件大小对比:"
    ls -lh "$IMAGE_PATH" | awk '{print "   输入: " $5 " " $9}'
    ls -lh "$OUTPUT_FILE" | awk '{print "   输出: " $5 " " $9}'
    echo ""
    echo "=========================================="
    echo "完成！请查看输出文件"
    echo "=========================================="
else
    echo "✗ 未找到输出文件: $OUTPUT_FILE"
    echo ""
    echo "检查目录内容:"
    ls -la "$IMAGE_DIR/"
    if [ -d "$IMAGE_DIR/deskewed" ]; then
        echo ""
        echo "deskewed 目录内容:"
        ls -la "$IMAGE_DIR/deskewed/"
    fi
    exit 1
fi
