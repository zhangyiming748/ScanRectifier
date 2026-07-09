package code

import (
	"ScanRectifier/deskew"
	"ScanRectifier/masking"
)

func Deskew(dir string) {
	deskew.ProcessDirectory(dir)
}

func Masking(dir string) {
	masking.ProcessDirectory(dir)
}

// ProcessAll 依次执行纠偏和边缘漂白
func ProcessAll(dir string) {
	deskew.ProcessDirectory(dir)
	masking.ProcessDirectory(dir)
}
