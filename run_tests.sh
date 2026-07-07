#!/bin/bash

# ScanRectifier 单元测试运行脚本
# 用于在 Linux 容器中运行 deskew 模块的测试

set -e

# 获取脚本所在目录的绝对路径
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$SCRIPT_DIR"

echo "=========================================="
echo "ScanRectifier Deskew Module Unit Tests"
echo "=========================================="
echo "📂 项目根目录: $PROJECT_ROOT"
echo ""

# 检查 CGO 是否启用
if [ -z "$CGO_ENABLED" ] || [ "$CGO_ENABLED" != "1" ]; then
    echo "⚠️  CGO_ENABLED 未设置或不为 1"
    echo "   GoCV 需要 CGO 支持，正在设置..."
    export CGO_ENABLED=1
    echo "✓ CGO_ENABLED=1"
fi

# 检查 OpenCV 是否安装
if ! pkg-config --exists opencv4 2>/dev/null; then
    echo "⚠️  OpenCV 4.x 未安装或未配置"
    echo "   请先运行: ./setup_container.sh"
    echo ""
    echo "   或者手动安装:"
    echo "   apt-get update && apt-get install -y libopencv-dev pkg-config"
    echo ""
    exit 1
else
    OPENCV_VERSION=$(pkg-config --modversion opencv4)
    echo "✓ OpenCV $OPENCV_VERSION 已安装"
fi

# 检查是否设置了测试图片目录
if [ -z "$TEST_IMAGE_DIR" ]; then
    echo "⚠️  TEST_IMAGE_DIR 环境变量未设置"
    echo "   用法: TEST_IMAGE_DIR=/path/to/images ./run_tests.sh"
    echo ""
    echo "   或者将测试图片放在以下目录之一："
    echo "   - ./test_images/"
    echo "   - /tmp/test_images/"
    echo ""
    
    # 尝试使用默认目录
    if [ -d "./test_images" ]; then
        export TEST_IMAGE_DIR="./test_images"
        echo "✓ 使用默认测试目录: ./test_images"
    elif [ -d "/tmp/test_images" ]; then
        export TEST_IMAGE_DIR="/tmp/test_images"
        echo "✓ 使用默认测试目录: /tmp/test_images"
    else
        echo "✗ 未找到测试图片目录，将只运行不依赖图片的测试"
        echo ""
    fi
fi

# 显示测试配置
if [ -n "$TEST_IMAGE_DIR" ]; then
    echo "📁 测试图片目录: $TEST_IMAGE_DIR"
    
    # 统计图片数量
    IMAGE_COUNT=$(find "$TEST_IMAGE_DIR" -type f \( -name "*.jpg" -o -name "*.jpeg" -o -name "*.png" -o -name "*.bmp" \) | wc -l)
    echo "📊 找到图片数量: $IMAGE_COUNT"
    
    if [ "$IMAGE_COUNT" -eq 0 ]; then
        echo "⚠️  警告: 测试目录中没有找到图片文件"
        echo ""
    fi
fi

echo ""
echo "开始运行测试..."
echo "------------------------------------------"

# 运行单元测试
echo ""
echo "1. 运行单元测试 (go test -v)"
echo "------------------------------------------"
cd "$PROJECT_ROOT"
go test -v ./deskew -run "^Test"

# 运行基准测试（可选）
echo ""
echo "2. 运行性能测试 (go test -bench)"
echo "------------------------------------------"
go test -bench=. -benchmem ./deskew || echo "⚠️  性能测试跳过（可能需要 TEST_IMAGE_DIR）"

# 生成测试覆盖率报告
echo ""
echo "3. 生成测试覆盖率报告"
echo "------------------------------------------"
go test -coverprofile=coverage.out ./deskew
go tool cover -func=coverage.out | tail -1

echo ""
echo "=========================================="
echo "✅ 测试完成！"
echo "=========================================="
echo ""
echo "查看详细覆盖率报告:"
echo "  go tool cover -html=coverage.out -o coverage.html"
echo ""
