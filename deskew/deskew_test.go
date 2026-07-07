package deskew

import (
	"os"
	"path/filepath"
	"testing"
)

// TestDeskewMyImage 最简单的纠偏测试
// 用法: TEST_IMAGE_PATH=/app/test_images/Scan_0014.jpg go test -v -run TestDeskewMyImage
func TestDeskewMyImage(t *testing.T) {
	// 你的图片路径 - 修改这里！
	imagePath := os.Getenv("TEST_IMAGE_PATH")
	if imagePath == "" {
		// 默认测试图片路径，你可以改成自己的
		imagePath = "test_images/Scan_0014.jpg"
	}

	// 检查文件是否存在
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		t.Fatalf("❌ 图片文件不存在: %s", imagePath)
	}

	t.Logf("📄 输入图片: %s", imagePath)

	// 输出路径（在同一目录下生成 deskewed_xxx.jpg）
	outputDir := filepath.Dir(imagePath)
	filename := filepath.Base(imagePath)
	outputPath := filepath.Join(outputDir, "deskewed_"+filename)

	t.Logf("📤 输出图片: %s", outputPath)

	// 执行纠偏
	angle, err := DeskewImage(imagePath, outputPath)
	if err != nil {
		t.Fatalf("❌ 纠偏失败: %v", err)
	}

	t.Logf("✅ 检测到倾斜角度: %.2f°", angle)

	// 验证输出文件
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Fatal("❌ 输出文件未生成")
	}

	info, _ := os.Stat(outputPath)
	t.Logf("✅ 纠偏成功！")
	t.Logf("   输出文件大小: %d bytes", info.Size())
	t.Logf("   请查看: %s", outputPath)
}

// TestProcessDirectory 测试目录批量处理
// 用法: TEST_IMAGE_DIR=/app/test_images go test -v -run TestProcessDirectory
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

	// 执行批量处理
	ProcessDirectory(testDir)

	// 检查是否生成了 deskewed 子目录
	deskewedDir := filepath.Join(testDir, "deskewed")
	if _, err := os.Stat(deskewedDir); os.IsNotExist(err) {
		t.Fatal("❌ deskewed 目录未生成")
	}

	// 列出生成的文件
	files, _ := os.ReadDir(deskewedDir)
	if len(files) == 0 {
		t.Fatal("❌ deskewed 目录为空")
	}

	t.Logf("✅ 批量处理成功！")
	t.Logf("   输出目录: %s", deskewedDir)
	t.Logf("   生成文件数: %d", len(files))
	t.Logf("   文件列表:")
	for _, f := range files {
		info, _ := f.Info()
		t.Logf("     - %s (%d bytes)", f.Name(), info.Size())
	}
}
