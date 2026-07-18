package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	// 命令行参数
	configDir = flag.String("config-dir", "config", "配置文件目录")
	port      = flag.Int("port", 0, "服务端口（覆盖配置文件）")
	env       = flag.String("env", "", "运行环境（覆盖配置文件）")
	verbosity = flag.String("verbosity", "", "日志级别（覆盖配置文件）")
	showHelp  = flag.Bool("help", false, "显示帮助信息")
	showVer   = flag.Bool("version", false, "显示版本信息")
	audit     = flag.Bool("audit", false, "开启审查模式（记录异常操作）")
)

// applyCommandLineOverrides 将命令行参数转换为环境变量
func applyCommandLineOverrides() {
	if *port > 0 {
		os.Setenv("BASIC_SERVER_PORT", fmt.Sprintf("%d", *port))
	}
	if *env != "" {
		os.Setenv("ENV", *env)
	}
	if *verbosity != "" {
		os.Setenv("BASIC_VERBOSITY", *verbosity)
	}
	if *audit {
		os.Setenv("BASIC_AUDIT", "true")
	}
}
