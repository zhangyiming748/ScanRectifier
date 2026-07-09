package main

import (
	"fmt"
	"os"
	"runtime"

	"ScanRectifier/code"

	"github.com/spf13/cobra"
)

// 版本信息（通过 -ldflags 在编译时注入）
var (
	Version   = "dev"
	GitCommit = "unknown"
	BuildTime = "unknown"
)

var rootCmd = &cobra.Command{
	Use:     "scanfix",
	Short:   "扫描图片修复工具",
	Long:    `ScanRectifier 是一个用于修复扫描图片的命令行工具，支持矫正倾斜和去除边缘黑线。`,
	Version: Version,
}

var deskewCmd = &cobra.Command{
	Use:   "deskew",
	Short: "矫正倾斜的扫描图片",
	Long:  `将扫描后稍微倾斜的图片修复为横平竖直的正向图片。`,
	Run: func(cmd *cobra.Command, args []string) {
		dir, _ := cmd.Flags().GetString("dir")
		code.Deskew(dir)
		fmt.Printf("Deskew: 处理目录 [%s] 下的图片...\n", dir)
	},
}

var maskingCmd = &cobra.Command{
	Use:   "masking",
	Short: "去除图片边缘黑线",
	Long:  `去除扫描图片边缘的黑线，使图片边缘更加干净。`,
	Run: func(cmd *cobra.Command, args []string) {
		dir, _ := cmd.Flags().GetString("dir")
		code.Masking(dir)
		fmt.Printf("Masking: 处理目录 [%s] 下的图片...\n", dir)
	},
}

var processCmd = &cobra.Command{
	Use:   "process",
	Short: "完整处理流程（纠偏 + 边缘漂白）",
	Long:  `依次执行图片纠偏和边缘漂白处理，处理完成后直接覆盖原文件。`,
	Run: func(cmd *cobra.Command, args []string) {
		dir, _ := cmd.Flags().GetString("dir")
		fmt.Printf("Process: 开始处理目录 [%s] 下的图片...\n\n", dir)
		code.ProcessAll(dir)
	},
}

func init() {
	deskewCmd.Flags().StringP("dir", "d", "./", "图片所在的根目录（必填）")

	maskingCmd.Flags().StringP("dir", "d", "./", "图片所在的根目录（必填）")

	processCmd.Flags().StringP("dir", "d", "./", "图片所在的根目录（必填）")

	rootCmd.AddCommand(deskewCmd, maskingCmd, processCmd, versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "显示版本信息",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("ScanRectifier %s\n", Version)
		fmt.Printf("Git Commit: %s\n", GitCommit)
		fmt.Printf("Build Time: %s\n", BuildTime)
		fmt.Printf("Go Version: %s %s/%s\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
