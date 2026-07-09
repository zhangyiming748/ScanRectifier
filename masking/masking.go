package masking

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"path/filepath"

	"github.com/zhangyiming748/finder"
	"gocv.io/x/gocv"
)

// MaskImage 对单张图片进行边缘掩膜漂白处理
// 输入: srcPath - 源图片路径, dpi - 扫描DPI
// 返回: 错误信息
// 处理成功后会用漂白结果覆盖原文件
func MaskImage(srcPath string, dpi int) error {
	// 生成临时文件名（在扩展名前添加"_临时"标记）
	ext := filepath.Ext(srcPath)
	tmpPath := srcPath[:len(srcPath)-len(ext)] + "_临时" + ext
	// 读取图片
	src := gocv.IMRead(srcPath, gocv.IMReadColor)
	if src.Empty() {
		return fmt.Errorf("无法读取图片: %s", srcPath)
	}
	defer src.Close()

	// 计算边缘漂白范围（3mm对应的像素值）
	// 公式: marginPx = 3 * DPI / 25.4
	marginPx := int(float64(dpi) * 3.0 / 25.4)

	fmt.Printf("DPI: %d, 边缘漂白范围: %d 像素 (约3mm)\n", dpi, marginPx)

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

	// 保存结果到临时文件
	if !gocv.IMWrite(tmpPath, result) {
		return fmt.Errorf("无法保存处理后的图片: %s", tmpPath)
	}

	// 处理成功：删除原文件，将临时文件重命名为原文件名
	if err := os.Remove(srcPath); err != nil {
		return fmt.Errorf("无法删除原文件: %s, %v", srcPath, err)
	}
	if err := os.Rename(tmpPath, srcPath); err != nil {
		return fmt.Errorf("无法重命名临时文件: %s -> %s, %v", tmpPath, srcPath, err)
	}

	fmt.Printf("✅ 边缘漂白完成\n")
	return nil
}

// ProcessDirectory 处理目录下的所有图片
// dir: 图片所在目录
// dpi: 扫描DPI（用于计算边缘漂白范围）
// 漂白成功后直接覆盖原文件
func ProcessDirectory(dir string) {

	dpi := 300 // 默认300 DPI

	images := finder.FindAllImages(dir)
	if len(images) == 0 {
		fmt.Printf("⚠️  在目录 %s 中未找到图片文件\n", dir)
		return
	}

	fmt.Printf("📁 找到 %d 个图片文件\n", len(images))
	fmt.Printf("📊 DPI设置: %d (边缘漂白范围: %d 像素，约3mm)\n\n", dpi, int(float64(dpi)*3.0/25.4))

	successCount := 0
	failCount := 0

	for _, img := range images {
		fmt.Printf("处理: %s\n", filepath.Base(img))

		err := MaskImage(img, dpi)
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
