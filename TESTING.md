# Deskew 模块单元测试指南

## 📋 测试概述

本测试套件包含以下测试内容：

1. **TestDeskewImage** - 单张图片纠偏功能测试
2. **TestProcessDirectory** - 目录批量处理测试
3. **TestMedian** - 中位数计算算法测试（纯逻辑，无需图片）
4. **BenchmarkDeskewImage** - 性能基准测试

## 🐳 在 Linux 容器中运行测试

### 方法 1: 使用自动化脚本（推荐）

```bash
# 赋予执行权限
chmod +x run_tests.sh

# 方式 A: 指定测试图片目录
TEST_IMAGE_DIR=/path/to/your/images ./run_tests.sh

# 方式 B: 将图片放在默认目录
mkdir -p test_images
cp your_images/*.jpg test_images/
./run_tests.sh
```

### 方法 2: 手动运行测试

```bash
# 进入项目根目录
cd /path/to/ScanRectifier

# 设置测试图片目录环境变量
export TEST_IMAGE_DIR=/path/to/test/images

# 运行所有测试（详细输出）
go test -v ./deskew

# 只运行特定测试
go test -v ./deskew -run TestMedian
go test -v ./deskew -run TestDeskewImage

# 运行性能测试
go test -bench=. -benchmem ./deskew

# 生成覆盖率报告
go test -coverprofile=coverage.out ./deskew
go tool cover -html=coverage.out -o coverage.html
```

### 方法 3: 在 Docker 容器中运行

```bash
# 构建开发容器
docker build -f Dockerfile.dev -t scanrectifier-dev .

# 运行测试（挂载测试图片目录）
docker run -it \
  -v $(pwd):/app \
  -v /path/to/your/images:/test_images \
  -e TEST_IMAGE_DIR=/test_images \
  scanrectifier-dev \
  bash -c "cd /app && ./run_tests.sh"
```

## 📁 准备测试图片

### 推荐的测试图片要求：

- **格式**: JPG, JPEG, PNG, BMP
- **数量**: 至少 3-5 张不同倾斜角度的图片
- **内容**: 包含文字或线条的文档扫描件效果最佳
- **倾斜角度**: 建议包含 -10° 到 +10° 范围内的图片

### 示例目录结构：

```
ScanRectifier/
├── deskew/
│   ├── deskew.go
│   └── deskew_test.go
├── test_images/          # 测试图片目录
│   ├── doc1.jpg
│   ├── doc2.png
│   └── doc3.jpg
├── run_tests.sh
└── TESTING.md
```

## 🔍 测试输出示例

```
==========================================
ScanRectifier Deskew Module Unit Tests
==========================================

📁 测试图片目录: /test_images
📊 找到图片数量: 5

开始运行测试...
------------------------------------------

1. 运行单元测试 (go test -v)
------------------------------------------
=== RUN   TestDeskewImage
=== RUN   TestDeskewImage/doc1.jpg
    deskew_test.go:57: Image: doc1.jpg, Detected angle: 3.45°
=== RUN   TestDeskewImage/doc2.png
    deskew_test.go:57: Image: doc2.png, Detected angle: -2.18°
--- PASS: TestDeskewImage (1.23s)

=== RUN   TestMedian
=== RUN   TestMedian/odd_count
=== RUN   TestMedian/even_count
--- PASS: TestMedian (0.00s)

PASS
ok      ScanRectifier/deskew    1.234s

✅ 测试完成！
```

## ⚙️ 环境要求

### Linux 容器依赖

确保容器中已安装：

```bash
# Go 语言环境
go version

# OpenCV 库（用于 gocv）
pkg-config --modversion opencv4

# CGO 支持
echo $CGO_ENABLED  # 应该为 1

# GCC 编译器
gcc --version
```

### 如果缺少依赖，在 Debian/Ubuntu 中安装：

```bash
apt-get update
apt-get install -y \
    build-essential \
    cmake \
    pkg-config \
    libopencv-dev \
    golang-go
```

## 📊 查看测试覆盖率

```bash
# 生成 HTML 覆盖率报告
go test -coverprofile=coverage.out ./deskew
go tool cover -html=coverage.out -o coverage.html

# 在浏览器中打开
xdg-open coverage.html  # Linux
open coverage.html      # macOS
start coverage.html     # Windows
```

## 🐛 故障排除

### 问题 1: `TEST_IMAGE_DIR not set`

**解决方案**: 
```bash
export TEST_IMAGE_DIR=/path/to/images
# 或者直接在命令中设置
TEST_IMAGE_DIR=/path/to/images go test -v ./deskew
```

### 问题 2: `no required module provides package gocv.io/x/gocv`

**解决方案**:
```bash
go mod download
go mod tidy
```

### 问题 3: OpenCV 相关错误

**解决方案**:
```bash
# 检查 OpenCV 是否正确安装
pkg-config --cflags opencv4
pkg-config --libs opencv4

# 如果没有输出，重新安装 OpenCV
apt-get install -y libopencv-dev
```

### 问题 4: CGO 未启用

**解决方案**:
```bash
export CGO_ENABLED=1
go test -v ./deskew
```

## 📝 添加新的测试用例

编辑 `deskew/deskew_test.go`，添加新的测试函数：

```go
func TestYourNewFeature(t *testing.T) {
    // 测试代码
}
```

然后运行：
```bash
go test -v ./deskew -run TestYourNewFeature
```

## 🎯 最佳实践

1. **隔离测试**: 每个测试使用独立的临时目录 (`t.TempDir()`)
2. **跳过机制**: 使用 `t.Skip()` 当环境不满足时优雅跳过
3. **详细日志**: 使用 `t.Logf()` 输出关键信息便于调试
4. **性能测试**: 对核心算法添加 Benchmark 测试
5. **覆盖率目标**: 保持测试覆盖率 > 80%
