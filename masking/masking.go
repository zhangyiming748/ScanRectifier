package masking

import (
	"fmt"
	"image"
	"image/color"
	"path/filepath"

	"github.com/zhangyiming748/finder"
	"gocv.io/x/gocv"
)

// MaskImage 对单张图片进行边缘掩膜漂白处理
// 输入: srcPath - 源图片路径, dstPath - 输出图片路径, dpi - 扫描DPI
// 返回: 错误信息
func MaskImage(srcPath, dstPath string, dpi int) error {
	// 读取图片
	src := gocv.IMRead(srcPath, gocv.IMReadColor)
	if src.Empty() {
		return fmt.Errorf("无法读取图片: %s", srcPath)
	}
	defer src.Close()

	// 计算边缘漂白范围（1mm对应的像素值）
	// 公式: marginPx = 1 * DPI / 25.4
	marginPx := int(float64(dpi) * 1.0 / 25.4)

	fmt.Printf("DPI: %d, 边缘漂白范围: %d 像素 (约1mm)\n", dpi, marginPx)

	height := src.Rows()
	width := src.Cols()

	// 检查图片尺寸是否足够
	if width <= 2*marginPx || height <= 2*marginPx {
		return fmt.Errorf("图片尺寸太小，无法应用边缘漂白")
	}

	// 创建与原图同尺寸的全白图像作为结果
	result := gocv.NewMatWithSize(height, width, src.Type())
	defer result.Close()

	// 先复制原图到结果
	src.CopyTo(&result)

	// 在结果图上绘制白色矩形覆盖边缘区域
	whiteColor := color.RGBA{R: 255, G: 255, B: 255, A: 255}

	// 上边缘
	gocv.Rectangle(&result,
		image.Rect(0, 0, width, marginPx),
		whiteColor,
		-1) // -1 表示填充

	// 下边缘
	gocv.Rectangle(&result,
		image.Rect(0, height-marginPx, width, height),
		whiteColor,
		-1)

	// 左边缘
	gocv.Rectangle(&result,
		image.Rect(0, 0, marginPx, height),
		whiteColor,
		-1)

	// 右边缘
	gocv.Rectangle(&result,
		image.Rect(width-marginPx, 0, width, height),
		whiteColor,
		-1)

	// 保存结果
	if !gocv.IMWrite(dstPath, result) {
		return fmt.Errorf("无法保存处理后的图片: %s", dstPath)
	}

	fmt.Printf("✅ 边缘漂白完成\n")
	return nil
}

// ProcessDirectory 处理目录下的所有图片
// dir: 图片所在目录
// dpi: 扫描DPI（用于计算边缘漂白范围）
// 在原目录下生成文件名后缀带有"净"字的漂白后图片
func ProcessDirectory(dir string) {

	dpi := 300 // 默认300 DPI

	images := finder.FindAllImages(dir)
	if len(images) == 0 {
		fmt.Printf("⚠️  在目录 %s 中未找到图片文件\n", dir)
		return
	}

	fmt.Printf("📁 找到 %d 个图片文件\n", len(images))
	fmt.Printf("📊 DPI设置: %d (边缘漂白范围: %d 像素，约1mm)\n\n", dpi, int(float64(dpi)*1.0/25.4))

	successCount := 0
	failCount := 0

	for _, img := range images {
		// 获取文件名和扩展名
		ext := filepath.Ext(img)
		nameWithoutExt := img[:len(img)-len(ext)]

		// 生成输出文件名：原文件名_净.扩展名
		dst := nameWithoutExt + "_净" + ext

		fmt.Printf("处理: %s -> %s\n", filepath.Base(img), filepath.Base(dst))

		err := MaskImage(img, dst, dpi)
		if err != nil {
			fmt.Printf("  ✗ 失败: %v\n", err)
			failCount++
		} else {
			fmt.Printf("  ✓ 成功\n")
			successCount++
		}
	}

	fmt.Printf("\n处理完成! 成功: %d, 失败: %d\n", successCount, failCount)
}
