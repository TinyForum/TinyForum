package wire

import (
	"fmt"
	"strings"
	"tiny-forum/config"
)

// 构建 DSN，自动跳过值为空的参数，避免因空密码导致的解析问题
func buildDSN(cfg *config.ConfigPostgres) string {
	var parts []string

	if cfg.Host != "" {
		parts = append(parts, fmt.Sprintf("host=%s", cfg.Host))
	}
	if cfg.Port != 0 {
		parts = append(parts, fmt.Sprintf("port=%d", cfg.Port))
	}
	if cfg.User != "" {
		parts = append(parts, fmt.Sprintf("user=%s", cfg.User))
	}
	// 关键优化：只有密码非空时才添加 password 参数
	if cfg.Password != "" {
		parts = append(parts, fmt.Sprintf("password=%s", cfg.Password))
	}
	if cfg.DBName != "" {
		parts = append(parts, fmt.Sprintf("dbname=%s", cfg.DBName))
	}
	if cfg.SSLMode != "" {
		parts = append(parts, fmt.Sprintf("sslmode=%s", cfg.SSLMode))
	}
	// TimeZone 可选，但注意值中可能含斜杠，建议使用 URL 格式替代
	if cfg.TimeZone != "" {
		parts = append(parts, fmt.Sprintf("TimeZone=%s", cfg.TimeZone))
	}

	return strings.Join(parts, " ")
}
