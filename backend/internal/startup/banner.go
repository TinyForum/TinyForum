package startup

import (
	_ "embed"
	"fmt"

	"tiny-forum/pkg/logger"
)

// 定义 ANSI 颜色码（不依赖外部库）
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorCyan   = "\033[36m"
	colorBold   = "\033[1m"
)

//go:embed banner
var bannerRaw string

func PrintBanner(version string) {
	banner := colorCyan + bannerRaw + colorReset + "\n" +
		colorYellow + "✨ Tiny Forum v" + version + " ✨" + colorReset + "\n"
	logger.Infof(banner)
}

// PrintStartupInfo 打印统一的启动信息
func PrintStartupInfo(cfg interface {
	GetServerPort() int
	GetServerMode() string
	GetDatabaseInfo() (host string, port int, user, dbName string)
	GetEmailInfo() (host string, port int, enabled bool)
	GetRedisInfo() (host string, port int, enabled bool)
	GetOllamaInfo() (baseURL string, model string)
	IsProduction() bool
}) {
	// 获取各项配置
	port := cfg.GetServerPort()
	mode := cfg.GetServerMode()
	dbHost, dbPort, dbUser, dbName := cfg.GetDatabaseInfo()
	emailHost, emailPort, emailEnabled := cfg.GetEmailInfo()
	redisHost, redisPort, redisEnabled := cfg.GetRedisInfo()
	ollamaURL, ollamaModel := cfg.GetOllamaInfo()

	// 使用 logger 或 fmt 输出
	logger.Info("=================================================")
	logger.Info(colorGreen + "🚀 Tiny Forum Server Starting..." + colorReset)
	logger.Info("=================================================")
	logger.Info(fmt.Sprintf("🌍 Environment:     %s", getEnvText(mode, cfg.IsProduction())))
	logger.Info(fmt.Sprintf("🔧 Server Address:  http://localhost:%d", port))
	logger.Info(fmt.Sprintf("📚 API Base Path:   /api/v1"))
	logger.Info(fmt.Sprintf("🗄️  Database:        %s@%s:%d/%s",
		dbUser, dbHost, dbPort, dbName))
	if emailEnabled {
		logger.Info(fmt.Sprintf("📧 Email Service:   %s:%d", emailHost, emailPort))
	} else {
		logger.Info("📧 Email Service:   Disabled")
	}
	if redisEnabled {
		logger.Info(fmt.Sprintf("📡 Redis:           %s:%d", redisHost, redisPort))
	} else {
		logger.Info("📡 Redis:           Disabled")
	}
	if ollamaURL != "" {
		logger.Info(fmt.Sprintf("🧠 Ollama:          %s (model: %s)", ollamaURL, ollamaModel))
	}
	logger.Info("=================================================")
}

// PrintConfigSummary 打印配置摘要（调试用，可隐藏敏感信息）
func PrintConfigSummary(cfg interface {
	GetJWTSecretMasked() string
	GetLogLevel() string
}) {
	logger.Debug("Configuration Summary:")
	logger.Debug(fmt.Sprintf("  • Log Level:      %s", cfg.GetLogLevel()))
	logger.Debug(fmt.Sprintf("  • JWT Secret:     %s", cfg.GetJWTSecretMasked()))
	// 可继续添加其他非敏感字段
}

// getEnvText 返回带颜色的环境文本
func getEnvText(mode string, isProd bool) string {
	if isProd {
		return colorRed + "🚀 Production" + colorReset
	}
	return colorGreen + "💻 Development" + colorReset
}
