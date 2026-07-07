#!/bin/bash

# 简单直接的纠偏测试
# 用法: ./quick_deskew_test.sh /path/to/your/image.jpg

set -e

echo "=========================================="
echo "快速纠偏测试"
echo "=========================================="
echo ""

# 检查参数
if [ -z "$1" ]; then
    echo "用法: $0 <图片路径>"
    echo ""
    echo "示例:"
    echo "  $0 ./test_images/doc.jpg"
    echo "  $0 /app/test_images/scan.png"
    exit 1
fi

IMAGE_PATH="$1"

# 检查文件是否存在
if [ ! -f "$IMAGE_PATH" ]; then
    echo "✗ 文件不存在: $IMAGE_PATH"
    exit 1
fi

echo "📄 测试图片: $IMAGE_PATH"
echo ""

# 设置环境变量
export CGO_ENABLED=1

# 创建输出目录
OUTPUT_DIR="/tmp/deskew_output"
mkdir -p "$OUTPUT_DIR"

# 获取文件名
FILENAME=$(basename "$IMAGE_PATH")
OUTPUT_PATH="$OUTPUT_DIR/$FILENAME"

echo "🔄 开始纠偏处理..."
echo ""

# 设置测试图片目录环境变量
export TEST_IMAGE_DIR=$(dirname "$IMAGE_PATH")
echo "📁 测试目录: $TEST_IMAGE_DIR"
echo ""

# 运行 Go 测试
cd /app
go test -v ./deskew -run TestDeskewImage \
    2>&1 | tee /tmp/test_output.txt

# 如果 go test 失败，直接用 main 程序测试
if ! grep -q "PASS" /tmp/test_output.txt; then
    echo ""
    echo "⚠️  go test 未通过，尝试直接编译运行..."
    echo ""
    
    # 编译程序
    echo "🔨 编译程序..."
    go build -o /tmp/scanrectifier .
    
    # 运行 deskew 命令
    echo "🚀 执行纠偏..."
    IMAGE_DIR=$(dirname "$IMAGE_PATH")
    /tmp/scanrectifier deskew --dir "$IMAGE_DIR"
    
    # 检查输出
    OUTPUT_FILE="$IMAGE_DIR/deskewed/$FILENAME"
    if [ -f "$OUTPUT_FILE" ]; then
        echo ""
        echo "✅ 纠偏成功！"
        echo "   输入: $IMAGE_PATH"
        echo "   输出: $OUTPUT_FILE"
        echo ""
        echo "📊 查看结果："
        ls -lh "$OUTPUT_FILE"
    else
        echo ""
        echo "✗ 纠偏失败，未找到输出文件"
        exit 1
    fi
else
    echo ""
    echo "✅ 测试通过！"
    echo "   输出目录: $OUTPUT_DIR"
    echo ""
    ls -lh "$OUTPUT_DIR/"
fi

echo ""
echo "=========================================="
echo "完成！"
echo "=========================================="
