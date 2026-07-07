package main

import (
	"fmt"
	"os"

	"ScanRectifier/code"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "scanfix",
	Short: "扫描图片修复工具",
	Long:  `ScanRectifier 是一个用于修复扫描图片的命令行工具，支持矫正倾斜和去除边缘黑线。`,
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

func init() {
	deskewCmd.Flags().StringP("dir", "d", "./", "图片所在的根目录（必填）")

	maskingCmd.Flags().StringP("dir", "d", "./", "图片所在的根目录（必填）")

	rootCmd.AddCommand(deskewCmd, maskingCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
