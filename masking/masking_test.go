package masking

import (
	"os"
	"testing"
)

// TestProcessDirectory 测试目录批量处理
func TestProcessDirectory(t *testing.T) {
	// 测试目录
	testDir := "/app"

	// 检查目录是否存在
	if _, err := os.Stat(testDir); os.IsNotExist(err) {
		t.Fatalf("目录不存在: %s", testDir)
	}

	// 执行批量处理
	ProcessDirectory(testDir)

	t.Log("处理完成")
}
