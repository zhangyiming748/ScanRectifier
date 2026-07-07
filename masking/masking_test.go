package masking

import (
	"os"
	"path/filepath"
	"testing"
)

// TestMaskImage 测试单张图片边缘漂白
// 用法: TEST_IMAGE_PATH=/path/to/image.jpg go test -v -run TestMaskImage
func TestMaskImage(t *testing.T) {
	// 获取测试图片路径
	imagePath := os.Getenv("TEST_IMAGE_PATH")
	if imagePath == "" {
		// 默认测试图片路径
		imagePath = "test_images/Scan_0014.jpg"
	}

	// 检查文件是否存在
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		t.Fatalf("❌ 图片文件不存在: %s", imagePath)
	}

	t.Logf("📄 输入图片: %s", imagePath)

	// 输出路径（在同一目录下生成 xxx_净.jpg）
	outputDir := filepath.Dir(imagePath)
	filename := filepath.Base(imagePath)
	ext := filepath.Ext(filename)
	nameWithoutExt := filename[:len(filename)-len(ext)]
	outputPath := filepath.Join(outputDir, nameWithoutExt+"_净"+ext)

	t.Logf("📤 输出图片: %s", outputPath)
	t.Logf("📊 DPI设置: 300 (边缘漂白范围: 约12像素，1mm)")

	// 执行边缘漂白
	err := MaskImage(imagePath, outputPath, 300)
	if err != nil {
		t.Fatalf("❌ 边缘漂白失败: %v", err)
	}

	// 验证输出文件
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Fatal("❌ 输出文件未生成")
	}

	info, _ := os.Stat(outputPath)
	t.Logf("✅ 边缘漂白成功！")
	t.Logf("   输出文件大小: %d bytes", info.Size())
	t.Logf("   请查看: %s", outputPath)
}

// TestProcessDirectory 测试目录批量处理
// 用法: TEST_IMAGE_DIR=/path/to/images go test -v -run TestProcessDirectory
func TestProcessDirectory(t *testing.T) {
	// 获取测试图片目录
	testDir := os.Getenv("TEST_IMAGE_DIR")
	if testDir == "" {
		// 默认测试目录
		testDir = "test_images"
	}

	// 检查目录是否存在
	if _, err := os.Stat(testDir); os.IsNotExist(err) {
		t.Fatalf("❌ 目录不存在: %s", testDir)
	}

	t.Logf("📁 测试目录: %s", testDir)
	t.Logf("📊 DPI设置: 300")

	// 执行批量处理
	ProcessDirectory(testDir)

	// 检查是否生成了带"净"字的文件
	files, _ := os.ReadDir(testDir)
	cleanedCount := 0
	for _, f := range files {
		if !f.IsDir() && filepath.Ext(f.Name()) != "" {
			name := f.Name()
			if len(name) > 2 && name[len(name)-5:len(name)-4] == "净" {
				cleanedCount++
			}
		}
	}

	if cleanedCount == 0 {
		t.Log("⚠️  未找到生成的漂白文件")
	} else {
		t.Logf("✅ 批量处理成功！")
		t.Logf("   生成漂白文件数: %d", cleanedCount)
	}
}
