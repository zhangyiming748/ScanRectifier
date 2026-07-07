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
