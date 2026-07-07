#!/bin/bash

# 验证 GoCV API 是否正确可用

set -e

echo "=========================================="
echo "验证 GoCV API 可用性"
echo "=========================================="
echo ""

export CGO_ENABLED=1

echo "1. 检查 Go 版本..."
go version

echo ""
echo "2. 检查 GoCV 版本..."
go list -m gocv.io/x/gocv

echo ""
echo "3. 检查 OpenCV..."
if pkg-config --exists opencv4; then
    echo "✓ OpenCV $(pkg-config --modversion opencv4) 已安装"
else
    echo "✗ OpenCV 未找到"
    exit 1
fi

echo ""
echo "4. 尝试编译 deskew 包..."
cd /app
if go build -v ./deskew 2>&1 | tee /tmp/build_output.txt; then
    echo "✓ 编译成功！"
    echo ""
    echo "5. 清理编译产物..."
    rm -f deskew
    echo ""
    echo "=========================================="
    echo "✅ 所有检查通过！代码可以正常编译"
    echo "=========================================="
else
    echo "✗ 编译失败"
    echo ""
    echo "完整错误信息："
    cat /tmp/build_output.txt
    exit 1
fi
