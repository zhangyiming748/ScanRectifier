#!/bin/bash

# ScanRectifier Linux 容器环境配置脚本
# 用于在 Debian/Ubuntu 容器中配置 GoCV 开发环境

set -e

echo "=========================================="
echo "ScanRectifier 容器环境配置"
echo "=========================================="
echo ""

# 1. 更新包管理器
echo "📦 更新软件包列表..."
apt-get update -qq

# 2. 安装基础依赖
echo "🔧 安装基础构建工具..."
apt-get install -y -qq \
    build-essential \
    cmake \
    pkg-config \
    git \
    wget > /dev/null 2>&1

# 3. 安装 OpenCV 开发库
echo "📷 安装 OpenCV 4.x..."
apt-get install -y -qq \
    libopencv-dev \
    libopencv-core-dev \
    libopencv-imgproc-dev \
    libopencv-imgcodecs-dev > /dev/null 2>&1

# 验证 OpenCV 安装
if pkg-config --exists opencv4; then
    OPENCV_VERSION=$(pkg-config --modversion opencv4)
    echo "✓ OpenCV $OPENCV_VERSION 已安装"
else
    echo "✗ OpenCV 安装失败"
    exit 1
fi

# 4. 配置 CGO 环境变量
echo "⚙️  配置 CGO 环境变量..."
export CGO_ENABLED=1
export CGO_CXXFLAGS="--std=c++11"
export CGO_CPPFLAGS="-I/usr/include/opencv4"
export CGO_LDFLAGS="$(pkg-config --libs opencv4)"

echo "✓ CGO_ENABLED=$CGO_ENABLED"
echo "✓ OpenCV 库路径: $(pkg-config --libs opencv4)"

# 5. 验证 Go 环境
echo ""
echo "🔍 检查 Go 环境..."
go version
echo "✓ Go 版本正常"

# 6. 下载 Go 模块依赖
echo ""
echo "📥 下载 Go 模块依赖..."
cd /app
go mod download
go mod tidy
echo "✓ 依赖下载完成"

# 7. 测试编译
echo ""
echo "🔨 测试编译..."
if go build -o /tmp/scanrectifier .; then
    echo "✓ 编译成功"
else
    echo "✗ 编译失败"
    exit 1
fi

echo ""
echo "=========================================="
echo "✅ 环境配置完成！"
echo "=========================================="
echo ""
echo "现在可以运行测试:"
echo "  ./run_tests.sh"
echo ""
echo "或者手动运行:"
echo "  export CGO_ENABLED=1"
echo "  TEST_IMAGE_DIR=./test_images go test -v ./deskew"
echo ""
