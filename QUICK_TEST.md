# 在 Linux 容器中快速运行测试的命令

## 方法 1: 使用修复后的脚本（推荐）

```bash
cd /app
chmod +x run_tests.sh
./run_tests.sh
```

## 方法 2: 直接运行 go test

```bash
cd /app
export TEST_IMAGE_DIR=./test_images
go test -v ./deskew
```

## 方法 3: 一行命令

```bash
cd /app && TEST_IMAGE_DIR=./test_images go test -v ./deskew
```

## 如果看到 "go.mod file not found" 错误

确认你在正确的项目根目录：

```bash
# 检查是否在包含 go.mod 的目录
ls -la /app/go.mod

# 如果不在，找到正确的目录
find / -name "go.mod" -path "*/ScanRectifier/*" 2>/dev/null

# 然后 cd 到那个目录
cd /path/to/ScanRectifier
```

## 完整的测试流程

```bash
# 1. 进入项目目录
cd /app

# 2. 确认 go.mod 存在
ls go.mod

# 3. 准备测试图片
mkdir -p test_images
# 复制你的测试图片到 test_images 目录

# 4. 运行测试
TEST_IMAGE_DIR=./test_images go test -v ./deskew

# 5. 查看覆盖率
go test -coverprofile=coverage.out ./deskew
go tool cover -func=coverage.out
```

## 常见问题

### Q: 找不到 go.mod
**A**: 确保你在项目根目录（包含 go.mod 文件的目录）

### Q: 找不到 gocv 包
**A**: 运行 `go mod download` 下载依赖

### Q: OpenCV 相关错误
**A**: 确保已安装 OpenCV: `apt-get install -y libopencv-dev`

### Q: CGO 未启用
**A**: 设置环境变量: `export CGO_ENABLED=1`
