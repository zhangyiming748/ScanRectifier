# 在 Linux 容器中运行测试 - 完整指南

## 🚨 问题诊断

你遇到的错误是因为：
1. **CGO 未启用** - GoCV 需要 CGO 支持
2. **OpenCV 未安装** - 缺少 OpenCV 4.x 库

## ✅ 解决方案（3 步）

### 第 1 步：配置容器环境

在容器中运行：

```bash
cd /app
chmod +x setup_container.sh
./setup_container.sh
```

这个脚本会自动：
- ✓ 安装 build-essential, cmake, pkg-config
- ✓ 安装 OpenCV 4.x 开发库
- ✓ 配置 CGO 环境变量
- ✓ 下载 Go 模块依赖
- ✓ 验证编译成功

### 第 2 步：准备测试图片

```bash
mkdir -p test_images
# 将你的扫描图片复制到 test_images 目录
cp /path/to/your/images/*.jpg test_images/
```

### 第 3 步：运行测试

```bash
./run_tests.sh
```

或者手动运行：

```bash
export CGO_ENABLED=1
TEST_IMAGE_DIR=./test_images go test -v ./deskew
```

## 🔧 手动安装（如果脚本失败）

如果自动化脚本失败，可以手动执行：

```bash
# 1. 安装系统依赖
apt-get update
apt-get install -y \
    build-essential \
    cmake \
    pkg-config \
    libopencv-dev \
    libopencv-core-dev \
    libopencv-imgproc-dev \
    libopencv-imgcodecs-dev

# 2. 验证 OpenCV 安装
pkg-config --modversion opencv4
# 应该输出类似: 4.5.4

# 3. 设置环境变量
export CGO_ENABLED=1
export CGO_CXXFLAGS="--std=c++11"
export CGO_CPPFLAGS="-I/usr/include/opencv4"
export CGO_LDFLAGS="$(pkg-config --libs opencv4)"

# 4. 下载依赖
cd /app
go mod download
go mod tidy

# 5. 测试编译
go build .
```

## 📋 验证清单

运行以下命令确认环境正确：

```bash
# 检查 CGO
echo $CGO_ENABLED  # 应该输出: 1

# 检查 OpenCV
pkg-config --modversion opencv4  # 应该输出版本号

# 检查 GoCV
go list -m gocv.io/x/gocv  # 应该输出: gocv.io/x/gocv v0.43.0

# 尝试编译
go build .  # 应该成功，无错误
```

## 🐛 常见问题

### Q1: `undefined: gocv.Point`
**原因**: GoCV v0.43.0 使用 `image.Point` 而不是 `gocv.Point`  
**解决**: 代码已修复，确保使用最新版本

### Q2: `package gocv.io/x/gocv: cannot find package`
**原因**: 依赖未下载  
**解决**: 
```bash
go mod download
go mod tidy
```

### Q3: `opencv4.pc not found`
**原因**: OpenCV 未安装或 pkg-config 找不到  
**解决**:
```bash
apt-get install -y libopencv-dev pkg-config
export PKG_CONFIG_PATH=/usr/lib/x86_64-linux-gnu/pkgconfig:$PKG_CONFIG_PATH
```

### Q4: `CGO_ENABLED=0`
**原因**: CGO 被禁用  
**解决**:
```bash
export CGO_ENABLED=1
```

### Q5: 编译时出现 C++ 链接错误
**原因**: 缺少 C++ 标准库或编译器  
**解决**:
```bash
apt-get install -y build-essential g++
```

## 📊 预期输出

成功配置后，运行测试应该看到：

```
==========================================
ScanRectifier Deskew Module Unit Tests
==========================================
📂 项目根目录: /app
✓ OpenCV 4.5.4 已安装
📁 测试图片目录: ./test_images
📊 找到图片数量: 5

开始运行测试...
------------------------------------------

1. 运行单元测试 (go test -v)
------------------------------------------
=== RUN   TestDeskewImage
=== RUN   TestDeskewImage/doc1.jpg
    deskew_test.go:57: Image: doc1.jpg, Detected angle: 3.45°
--- PASS: TestDeskewImage (1.23s)

PASS
ok      ScanRectifier/deskew    1.234s

✅ 测试完成！
```

## 💡 提示

- 首次运行 `setup_container.sh` 可能需要几分钟下载和安装依赖
- 确保容器有足够的磁盘空间（OpenCV 约 200MB）
- 如果使用 Docker，建议在 Dockerfile 中预安装这些依赖
- 测试图片应该是真实的扫描文档，包含文字或线条效果最佳

## 🎯 下一步

环境配置完成后：
1. 放置测试图片到 `test_images/` 目录
2. 运行 `./run_tests.sh`
3. 查看测试结果和覆盖率报告
4. 根据需要调整纠偏算法参数
