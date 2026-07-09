package deskew

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"os"
	"path/filepath"

	"github.com/zhangyiming748/finder"
	"gocv.io/x/gocv"
)

// DeskewImage 对单张图片进行纠偏处理
// 输入: srcPath - 源图片路径
// 返回: 倾斜角度(度), 错误信息
// 处理成功后会用纠偏结果覆盖原文件
func DeskewImage(srcPath string) (float64, error) {
	// 生成临时文件名
	tmpPath := srcPath + ".tmp"

	// 读取图片
	src := gocv.IMRead(srcPath, gocv.IMReadColor)
	if src.Empty() {
		return 0, fmt.Errorf("无法读取图片: %s", srcPath)
	}
	defer src.Close()

	// 转换为灰度图
	gray := gocv.NewMat()
	defer gray.Close()
	if err := gocv.CvtColor(src, &gray, gocv.ColorBGRToGray); err != nil {
		return 0, fmt.Errorf("转换灰度图失败: %v", err)
	}

	// 高斯模糊去噪
	blurred := gocv.NewMat()
	defer blurred.Close()
	if err := gocv.GaussianBlur(gray, &blurred, image.Point{X: 5, Y: 5}, 0, 0, gocv.BorderDefault); err != nil {
		return 0, fmt.Errorf("高斯模糊失败: %v", err)
	}

	// Canny 边缘检测
	edges := gocv.NewMat()
	defer edges.Close()
	if err := gocv.Canny(blurred, &edges, 50, 150); err != nil {
		return 0, fmt.Errorf("边缘检测失败: %v", err)
	}

	// HoughLinesP 检测直线
	lines := gocv.NewMat()
	defer lines.Close()
	if err := gocv.HoughLinesPWithParams(edges, &lines, 1, math.Pi/180, 50, 50, 10); err != nil {
		return 0, fmt.Errorf("直线检测失败: %v", err)
	}

	// 计算倾斜角度
	angle := calculateAngle(lines)

	fmt.Printf("检测到倾斜角度: %.2f°\n", angle)

	// 如果倾斜角度很小，直接保存原图到临时文件
	if math.Abs(angle) < 0.5 {
		if !gocv.IMWrite(tmpPath, src) {
			return angle, fmt.Errorf("无法保存图片: %s", tmpPath)
		}
	} else {
		// 旋转图片
		rotated := rotateImage(src, angle)
		defer rotated.Close()

		// 保存结果到临时文件
		if !gocv.IMWrite(tmpPath, rotated) {
			return angle, fmt.Errorf("无法保存纠偏后的图片: %s", tmpPath)
		}
	}

	// 处理成功：删除原文件，将临时文件重命名为原文件名
	if err := os.Remove(srcPath); err != nil {
		return angle, fmt.Errorf("无法删除原文件: %s, %v", srcPath, err)
	}
	if err := os.Rename(tmpPath, srcPath); err != nil {
		return angle, fmt.Errorf("无法重命名临时文件: %s -> %s, %v", tmpPath, srcPath, err)
	}

	return angle, nil
}

// calculateAngle 从检测到的直线中计算倾斜角度
func calculateAngle(lines gocv.Mat) float64 {
	if lines.Empty() || lines.Rows() == 0 {
		return 0
	}

	var angles []float64

	// 遍历所有检测到的直线
	for i := 0; i < lines.Rows(); i++ {
		// HoughLinesP 返回的格式: [x1, y1, x2, y2]
		x1 := float64(lines.GetIntAt(i, 0))
		y1 := float64(lines.GetIntAt(i, 1))
		x2 := float64(lines.GetIntAt(i, 2))
		y2 := float64(lines.GetIntAt(i, 3))

		// 计算直线角度（弧度）
		angle := math.Atan2(y2-y1, x2-x1) * 180.0 / math.Pi

		// 过滤掉接近垂直的线（可能是文档边框）
		if math.Abs(angle) < 45 || math.Abs(angle) > 135 {
			// 将角度归一化到 -90 ~ 90 范围
			if angle > 45 {
				angle -= 90
			} else if angle < -45 {
				angle += 90
			}
			angles = append(angles, angle)
		}
	}

	if len(angles) == 0 {
		return 0
	}

	// 使用中位数作为最终角度（更鲁棒）
	return median(angles)
}

// median 计算浮点数切片的中位数
func median(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	// 排序
	sorted := make([]float64, len(values))
	copy(sorted, values)
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[i] > sorted[j] {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	// 取中位数
	n := len(sorted)
	if n%2 == 0 {
		return (sorted[n/2-1] + sorted[n/2]) / 2.0
	}
	return sorted[n/2]
}

// rotateImage 旋转图片（使用白色背景填充）
func rotateImage(src gocv.Mat, angle float64) gocv.Mat {
	height := src.Rows()
	width := src.Cols()

	// 计算旋转中心
	center := image.Point{X: width / 2, Y: height / 2}

	// 获取旋转矩阵
	rotationMatrix := gocv.GetRotationMatrix2D(center, angle, 1.0)
	defer rotationMatrix.Close()

	// 计算旋转后的图像尺寸
	radians := angle * math.Pi / 180.0
	sinVal := math.Abs(math.Sin(radians))
	cosVal := math.Abs(math.Cos(radians))

	newWidth := int(float64(height)*sinVal + float64(width)*cosVal)
	newHeight := int(float64(width)*sinVal + float64(height)*cosVal)

	// 调整旋转矩阵的平移分量
	rotationMatrix.SetFloatAt(0, 2, rotationMatrix.GetFloatAt(0, 2)+float32(newWidth-width)/2)
	rotationMatrix.SetFloatAt(1, 2, rotationMatrix.GetFloatAt(1, 2)+float32(newHeight-height)/2)

	// 执行仿射变换，使用白色背景填充
	dst := gocv.NewMat()
	whiteBorder := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	if err := gocv.WarpAffineWithParams(
		src, &dst, rotationMatrix,
		image.Point{X: newWidth, Y: newHeight},
		gocv.InterpolationLinear,
		gocv.BorderConstant,
		whiteBorder,
	); err != nil {
		// 如果出错，返回空 Mat
		return gocv.NewMat()
	}

	return dst
}

// ProcessDirectory 处理目录下的所有图片
// dir: 图片所在目录
// 纠偏成功后直接覆盖原文件
func ProcessDirectory(dir string) {
	images := finder.FindAllImages(dir)
	for _, img := range images {
		fmt.Printf("处理: %s\n", filepath.Base(img))
		DeskewImage(img)
	}
}
