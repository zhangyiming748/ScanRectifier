package test

import (
	"os"
	"testing"

	"ScanRectifier/deskew"
	"ScanRectifier/masking"
)

// go test -v ./test -run TestProcessDirectory 测试目录批量处理
func TestProcessDirectory(t *testing.T) {
	testDir := "/app"

	// 检查目录是否存在
	if _, err := os.Stat(testDir); os.IsNotExist(err) {
		t.Fatalf("❌ 目录不存在: %s", testDir)
	}

	t.Logf("📁 测试目录: %s", testDir)

	// 第一步：纠偏
	t.Log("=== 第一步：纠偏处理 ===")
	deskew.ProcessDirectory(testDir)

	

	t.Logf("✅ 纠偏完成")


	// 第二步：边缘漂白
	t.Log("=== 第二步：边缘漂白处理 ===")
	masking.ProcessDirectory(testDir)

	t.Log("✅ 全部处理完成")
}
