package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	flag.Parse()

	if *showHelp {
		printHelp()
		os.Exit(0)
	}

	version := os.Getenv("TINYFORUM_VERSION")
	if version == "" {
		version = "1.0.0"
	}

	if *showVer {
		fmt.Printf("Tiny Forum v%s\n", version)
		os.Exit(0)
	}

	// 应用命令行覆盖到环境变量
	applyCommandLineOverrides()

	// 启动应用
	if err := runApp(*configDir, version); err != nil {
		log.Fatalf("Failed to start: %v", err)
	}
}
