# Windows 环境下配置 GoCV 开发环境

## 问题说明

错误 `undefined: gocv.IMRead` 是因为：
1. CGO 被禁用（`CGO_ENABLED=0`）
2. 缺少 C/C++ 编译器（GCC）
3. 缺少 OpenCV 库

GoCV 需要通过 CGO 调用底层的 C++ OpenCV 库，因此必须正确配置这些依赖。

---

## 安装步骤

### 步骤 1：启用 CGO

在 PowerShell 中执行（管理员权限）：

```powershell
# 永久设置用户环境变量
[Environment]::SetEnvironmentVariable("CGO_ENABLED", "1", "User")

# 当前会话临时设置
$env:CGO_ENABLED="1"

# 验证
go env CGO_ENABLED  # 应该输出 1
```

### 步骤 2：安装 MinGW-w64（C/C++ 编译器）

#### 方法 A：使用 MSYS2（推荐）

1. 下载 MSYS2 安装包：https://www.msys2.org/
2. 安装到默认路径 `C:\msys64`
3. 打开 MSYS2 MinGW64 terminal
4. 运行以下命令安装 GCC：
   ```bash
   pacman -Sy
   pacman -S mingw-w64-x86_64-gcc
   ```
5. 将 `C:\msys64\mingw64\bin` 添加到系统 PATH

#### 方法 B：使用 WinLibs

1. 下载 WinLibs：https://winlibs.sourceforge.io/
2. 选择 `mingw-w64ucrt-static` 版本
3. 解压到 `C:\mingw64`
4. 将 `C:\mingw64\bin` 添加到系统 PATH

**验证安装：**
```powershell
gcc --version
# 应该显示 gcc 版本信息
```

### 步骤 3：安装 OpenCV

#### 方法 A：使用 GoCV 官方脚本（最简单）

```powershell
# 克隆 GoCV 仓库
git clone https://github.com/hybridgroup/gocv.git
cd gocv

# 运行 Windows 构建脚本（需要管理员权限）
.\win_build_opencv.cmd
```

这个脚本会：
- 自动下载 OpenCV 4.x
- 编译 OpenCV
- 设置必要的环境变量

#### 方法 B：手动安装 OpenCV

1. **下载 OpenCV**
   - 访问：https://opencv.org/releases/
   - 下载 Windows 版本（例如 `opencv-4.8.0-windows.exe`）

2. **解压到指定目录**
   ```powershell
   # 例如解压到 C:\opencv
   ```

3. **设置环境变量**
   ```powershell
   # 设置 OPENCV_DIR
   [Environment]::SetEnvironmentVariable("OPENCV_DIR", "C:\opencv\build", "User")
   
   # 添加 DLL 路径到 PATH
   $currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
   [Environment]::SetEnvironmentVariable("Path", "$currentPath;C:\opencv\build\x64\vc15\bin", "User")
   ```

4. **设置 GoCV 相关环境变量**
   ```powershell
   [Environment]::SetEnvironmentVariable("CGO_CXXFLAGS", "--std=c++11", "User")
   [Environment]::SetEnvironmentVariable("CGO_CPPFLAGS", "-IC:\opencv\build\include", "User")
   [Environment]::SetEnvironmentVariable("CGO_LDFLAGS", "-LC:\opencv\build\x64\vc15\lib -lopencv_world480", "User")
   ```

### 步骤 4：重新加载环境变量并测试

```powershell
# 重新启动 PowerShell 或刷新环境变量
$env:CGO_ENABLED="1"

# 进入项目目录
cd c:\Users\zhang\Github\ScanRectifier

# 清理并重新下载依赖
go clean -modcache
go mod tidy

# 尝试编译
go build ./deskew/
```

---

## 常见问题

### Q1: 编译时提示 "gcc not found"
**解决：** 确保 MinGW-w64 的 bin 目录已添加到 PATH，并且可以运行 `gcc --version`

### Q2: 编译时提示 "opencv2/core.hpp: No such file or directory"
**解决：** 检查 `CGO_CPPFLAGS` 是否正确指向 OpenCV 的 include 目录

### Q3: 运行时提示 "找不到 opencv_world*.dll"
**解决：** 确保 OpenCV 的 bin 目录已添加到 PATH

### Q4: 仍然显示 "undefined: gocv.xxx"
**解决：** 
1. 确认 `go env CGO_ENABLED` 输出为 `1`
2. 确认 `gcc --version` 可以正常执行
3. 运行 `go clean -cache -modcache` 后重新编译

---

## 验证安装

创建测试文件 `test_gocv.go`：

```go
package main

import (
    "fmt"
    "gocv.io/x/gocv"
)

func main() {
    // 创建一个空白图像
    mat := gocv.NewMatWithSize(100, 100, gocv.MatTypeCV8UC3)
    defer mat.Close()
    
    if !mat.Empty() {
        fmt.Println("✓ GoCV 环境配置成功！")
        fmt.Printf("  创建的 Mat: %dx%d\n", mat.Rows(), mat.Cols())
    } else {
        fmt.Println("✗ GoCV 初始化失败")
    }
}
```

运行测试：
```powershell
go run test_gocv.go
```

如果输出 "✓ GoCV 环境配置成功！"，说明环境配置正确。

---

## 参考资源

- GoCV 官方文档：https://gocv.io/
- GoCV GitHub：https://github.com/hybridgroup/gocv
- OpenCV 下载：https://opencv.org/releases/
- MSYS2 下载：https://www.msys2.org/
