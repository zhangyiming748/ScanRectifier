# ScanRectifier (扫描件自动纠偏与边缘漂白工具)

项目简介
ScanRectifier 是一个基于 Go 语言与 OpenCV (gocv) 开发的图像自动化处理工具。专为解决扫描仪扫描 A4 纸张时出现的倾斜问题和边缘固定黑线问题而设计。

本项目目前处于单图测试阶段，旨在验证核心算法的可行性，为后续的全自动化批量处理打下基础。

核心功能
1. 自动纠偏 (Auto Deskew)：通过霍夫直线变换（HoughLinesP）智能检测文档倾斜角度，并进行仿射变换，使文档恢复横平竖直。
2. 边缘掩膜漂白 (Edge Masking)：针对扫描仪在纸张边缘（如距边缘 1~5 毫米处）产生的固定黑线，通过计算 DPI 动态生成边缘掩膜（Mask），强制将边缘区域的黑色像素覆盖为纯白色，同时完美保留文档中心的灰度细节（如印章、照片）。

️ 技术栈
- 编程语言: Go (Golang)
- 核心依赖: gocv (Go bindings for OpenCV)
- 图像处理算法: Canny 边缘检测, HoughLinesP, 仿射变换 (WarpAffine), 掩膜操作 (SetTo with Mask)

核心代码在core文件夹中实现
main函数位置使用cobra实现命令行工具

快速开始

1. 环境准备
由于依赖 OpenCV，请确保您的系统已安装 OpenCV 4.x 版本。
# macOS
brew install opencv

# Ubuntu / Debian
sudo apt-get install libopencv-dev

2. 安装依赖
go mod tidy

3. 运行单图测试
将待测试的扫描件放入 test_images/ 目录，然后运行：
go run main.go -input test_images/scan_01.jpg -dpi 300

️ 核心算法流程 (伪代码逻辑)

Step 1: 自动纠偏 (Deskew)
1. 将原图转为灰度图，并进行高斯模糊（去噪）。
2. 使用 Canny 提取边缘，通过 HoughLinesP 找出直线。
3. 计算所有直线角度的中位数，得出精准倾斜角。
4. 若倾斜角 > 0.5°，生成旋转矩阵并进行仿射变换。

Step 2: 边缘 5mm 强制漂白 (Edge Masking)
1. 根据传入的 DPI 参数，计算 5 毫米对应的像素值：
   marginPx = int(5 * DPI / 25.4) （300 DPI 下约为 59 像素）
2. 创建与原图同尺寸的全白掩膜（Mask）。
3. 将掩膜中间的安全区域（上下左右各缩进 marginPx）填充为黑色。
4. 使用 SetTo(纯白, mask)：将掩膜为白色的边缘区域，在原图上强制设为纯白（255）。

️ 注意事项
- DPI 参数至关重要：边缘漂白的像素范围强依赖于扫描时的 DPI。请确保传入正确的 DPI（如 300），否则可能导致漂白范围过大（吃掉正文）或过小（黑线残留）。
- 内存管理：在 Go 中使用 gocv 时，请务必注意 defer mat.Close()，防止内存泄漏。
- 色彩空间：处理前请确保图像已统一转换为 BGR 或灰度图，避免 Alpha 通道干扰。

后续规划 (Roadmap)
- [x] 单图测试验证（当前阶段）
- [ ] 批量文件夹处理与并发优化 (Goroutine + Channel)
- [ ] 增加对倾斜角度过大导致边缘出现黑边的二次裁切功能
- [ ] 提供 CLI 命令行参数支持

贡献与反馈
如果您在测试中发现边缘漂白效果不佳，或者纠偏算法在某些复杂排版下失效，请随时记录问题并反馈！

这份 README 涵盖了项目背景、技术栈、目录结构、运行指南以及核心算法的详细说明，非常适合作为单图测试阶段的指导文档。

您觉得这份文档的结构和细节符合您的预期吗？
1. 需要我为您补充一些具体的 Go 依赖安装命令（比如 gocv 的环境配置）吗？
2. 或者需要我直接基于这份 README 中的伪代码，为您生成完整的 main.go 和 deskew.go 代码，方便您立刻跑起来测试？
