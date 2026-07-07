#!/bin/bash
# ScanRectifier - Debian 环境配置脚本
# 用于在 Debian/Ubuntu Docker 容器中安装 GoCV 所需依赖

set -e  # 遇到错误立即退出

echo "========================================="
echo "  ScanRectifier 环境配置 (Debian)"
echo "========================================="
echo ""

# 检查是否在 Debian/Ubuntu 系统
if [ ! -f /etc/debian_version ] && [ ! -f /etc/os-release ]; then
    echo "⚠️  警告: 这可能不是 Debian/Ubuntu 系统"
fi

# 步骤 1: 更新包列表
echo "📦 步骤 1/4: 更新包列表..."
apt-get update -y
echo "✓ 包列表更新完成"
echo ""

# 步骤 2: 安装编译工具
echo "🔨 步骤 2/4: 安装编译工具 (GCC, G++, Make, CMake)..."
apt-get install -y gcc g++ make cmake pkg-config
echo "✓ 编译工具安装完成"
echo ""

# 步骤 3: 安装 OpenCV 开发库
echo "📚 步骤 3/4: 安装 OpenCV 开发库..."
apt-get install -y libopencv-dev
echo "✓ OpenCV 安装完成"
echo ""

# 步骤 4: 验证安装
echo "✅ 步骤 4/4: 验证安装..."
echo ""

# 检查 GCC
if command -v gcc &> /dev/null; then
    GCC_VERSION=$(gcc --version | head -n1)
    echo "✓ GCC: $GCC_VERSION"
else
    echo "✗ GCC 未安装"
    exit 1
fi

# 检查 OpenCV
if pkg-config --exists opencv4; then
    OPENCV_VERSION=$(pkg-config --modversion opencv4)
    echo "✓ OpenCV: $OPENCV_VERSION"
elif pkg-config --exists opencv; then
    OPENCV_VERSION=$(pkg-config --modversion opencv)
    echo "✓ OpenCV: $OPENCV_VERSION (旧版本)"
else
    echo "✗ OpenCV 未正确安装"
    exit 1
fi

# 检查 CGO
CGO_ENABLED=$(go env CGO_ENABLED)
echo "✓ CGO_ENABLED: $CGO_ENABLED"

if [ "$CGO_ENABLED" != "1" ]; then
    echo "⚠️  警告: CGO 未启用，正在设置..."
    export CGO_ENABLED=1
    go env -w CGO_ENABLED=1
    echo "✓ CGO 已启用"
fi

echo ""
echo "========================================="
echo "  ✅ 环境配置完成！"
echo "========================================="
echo ""
echo "现在可以运行以下命令测试："
echo "  cd /path/to/ScanRectifier"
echo "  go build ./deskew/"
echo ""
