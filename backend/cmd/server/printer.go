package main

import (
	"flag"
	"fmt"
)

// printHelp 打印帮助信息
func printHelp() {
	fmt.Println(`Tiny Forum - 轻量级论坛系统

Usage:
  tiny-forum [flags]

Flags:`)
	flag.PrintDefaults()
	fmt.Print(`
Environment Variables:
  BASIC_SERVER_PORT       服务端口
  BASIC_LOG_LEVEL         日志级别 (debug/info/warn/error)
  ENV                     运行环境 (development/production)
  BASIC_DATABASE_HOST     数据库主机
  BASIC_DATABASE_PORT     数据库端口
  BASIC_DATABASE_USER     数据库用户
  BASIC_DATABASE_PASSWORD 数据库密码
  BASIC_DATABASE_DBNAME   数据库名称
  BASIC_JWT_SECRET        JWT密钥

Examples:
  tiny-forum                                    # 使用默认配置
  tiny-forum --config-dir /etc/myapp            # 使用自定义配置目录
  tiny-forum --port 9090 --env production       # 覆盖端口和环境
  tiny-forum --log-level debug                  # 设置日志级别
  tiny-forum --version                          # 显示版本信息
  BASIC_SERVER_PORT=8080 ./tiny-forum           # 使用环境变量

Note:
  命令行参数优先级高于配置文件，环境变量优先级最高
`)
}
