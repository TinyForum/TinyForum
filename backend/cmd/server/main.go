package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"tiny-forum/config"
	"tiny-forum/internal/wire"
	"tiny-forum/pkg/logger"
)

// @title           Tiny Forum API
// @version         1.0
// @description     一个基于 Gin 的论坛系统 API
// @host            localhost:8080
// @BasePath        /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and the JWT token.

// 主函数，程序的入口点
func main() {
	// 加载配置
	cfg, err := loadConfig()
	if err != nil {
		printConfigError(err)
		os.Exit(1)
	}

	// 初始化日志
	if err := logger.Init(cfg.ToLoggerConfig()); err != nil {
		log.Fatalf("Failed to init logger: %v", err)
	}

	// 打印配置信息（调试用）
	printConfigInfo(cfg)

	// 初始化应用（数据库 + 路由）
	app, err := wire.InitApp(cfg)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Failed to initialize app: %v", err))
	}

	// 启动服务器
	startServer(cfg, app)
}

// loadConfig 加载配置文件
func loadConfig() (*config.Config, error) {
	// 配置文件目录
	configDir := "config"

	// 检查配置文件是否存在
	basicConfigPath := filepath.Join(configDir, "basic.yaml")
	privateConfigPath := filepath.Join(configDir, "private.yaml")

	// 检查配置文件
	printConfigFileStatus(basicConfigPath, "Basic config")
	printConfigFileStatus(privateConfigPath, "Private config")

	// 加载配置
	cfg, err := config.Load(configDir)
	if err != nil {
		return nil, fmt.Errorf("failed to load config from '%s' directory", configDir)
	}

	// 验证并提示配置问题
	if err := validateConfigWithHints(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// printConfigFileStatus 打印配置文件状态
func printConfigFileStatus(filePath, configName string) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Printf("ℹ️  %s file not found: %s", configName, filePath)
		log.Printf("   → Using environment variables or defaults for this configuration")
	} else {
		log.Printf("✓ %s file found: %s", configName, filePath)
	}
}

// validateConfigWithHints 验证配置并提供友好的提示
func validateConfigWithHints(cfg *config.Config) error {
	var errors []string
	var warnings []string

	// 1. 验证服务器配置
	if cfg.Basic.Server.Port <= 0 || cfg.Basic.Server.Port > 65535 {
		errors = append(errors, fmt.Sprintf(
			"❌ Invalid server port: %d (must be between 1 and 65535)\n   → Set 'server.port' in config/basic.yaml or BASIC_SERVER_PORT environment variable",
			cfg.Basic.Server.Port))
	}

	// 2. 验证数据库配置
	if cfg.Private.Database.Host == "" {
		errors = append(errors, "❌ Database host is required\n   → Set 'database.host' in config/basic.yaml or BASIC_DATABASE_HOST environment variable")
	}

	if cfg.Private.Database.DBName == "" {
		errors = append(errors, "❌ Database name is required\n   → Set 'database.dbname' in config/basic.yaml or BASIC_DATABASE_DBNAME environment variable")
	}

	// 3. 验证 JWT 配置（关键配置，必须有值）
	if cfg.Private.JWT.Secret == "" {
		// 检查是否通过环境变量设置
		envSecret := os.Getenv("BASIC_JWT_SECRET")
		if envSecret == "" {
			errors = append(errors, fmt.Sprintf(
				"❌ JWT secret is required for security\n"+
					"   → Set 'jwt.secret' in config/basic.yaml\n"+
					"   → Or set BASIC_JWT_SECRET environment variable\n"+
					"   → Example for development: echo 'BASIC_JWT_SECRET=dev-secret-key' > .env\n"+
					"   → Example for production: Use a strong random string (at least 32 characters)"))
		} else {
			warnings = append(warnings, "⚠️  JWT secret is set via environment variable (this is fine)")
		}
	} else if len(cfg.Private.JWT.Secret) < 32 && cfg.IsProduction() {
		warnings = append(warnings, fmt.Sprintf(
			"⚠️  JWT secret is too short (%d characters) for production\n   → Use at least 32 characters for security",
			len(cfg.Private.JWT.Secret)))
	}

	// 4. 验证邮件配置（可选，仅提示）
	if cfg.Private.Email.Host != "" {
		if cfg.Private.Email.Port <= 0 || cfg.Private.Email.Port > 65535 {
			warnings = append(warnings, fmt.Sprintf(
				"⚠️  Invalid email port: %d, email features may not work",
				cfg.Private.Email.Port))
		}
		if cfg.Private.Email.Password == "" {
			warnings = append(warnings, "⚠️  Email password is not set, email sending will fail\n   → Set email.password in config/private.yaml or PRIVATE_EMAIL_PASSWORD environment variable")
		}
	}

	// 5. 验证 Redis 配置（可选）
	if cfg.Basic.Redis.Host != "" && cfg.Basic.Redis.Port == 0 {
		warnings = append(warnings, "⚠️  Redis host is set but port is 0, Redis features may not work")
	}

	// 输出警告信息
	for _, warning := range warnings {
		log.Printf("%s", warning)
	}

	// 如果有错误，返回友好的错误信息
	if len(errors) > 0 {
		return fmt.Errorf("\n%s\n\n💡 Quick Fix:\n  1. Create config/basic.yaml with required settings\n  2. Or set environment variables (see above)\n  3. Run 'make setup-config' to generate a sample config\n\n📖 Documentation: https://github.com/yourusername/tiny-forum#configuration\n",
			strings.Join(errors, "\n"))
	}

	return nil
}

// printConfigError 打印友好的配置错误信息
func printConfigError(err error) {
	log.Printf("%s", "\n"+strings.Repeat("=", 60))
	log.Printf("🚫 Configuration Error")
	log.Printf("%s", strings.Repeat("=", 60))
	log.Printf("\n%v\n", err)

	// 提供解决建议
	log.Printf("\n💡 Troubleshooting Tips:")
	log.Printf("   1. Make sure you have a config/basic.yaml file")
	log.Printf("   2. Check if config/private.yaml exists (optional)")
	log.Printf("   3. Verify all required fields are filled")
	log.Printf("\n📝 Quick Setup:")
	log.Printf("   # Create sample config file")
	log.Printf("   mkdir -p config")
	log.Printf("   cat > config/basic.yaml << EOF")
	log.Printf("   server:")
	log.Printf("     port: 8080")
	log.Printf("     mode: debug")
	log.Printf("   database:")
	log.Printf("     host: localhost")
	log.Printf("     port: 5432")
	log.Printf("     user: postgres")
	log.Printf("     password: yourpassword")
	log.Printf("     dbname: tiny_forum")
	log.Printf("   jwt:")
	log.Printf("     secret: your-super-secret-key-change-this")
	log.Printf("   EOF")
	log.Printf("\n   # Or use environment variables")
	log.Printf("   export BASIC_JWT_SECRET=\"your-secret-key\"")
	log.Printf("   export BASIC_DATABASE_HOST=\"localhost\"")
	log.Printf("   export BASIC_DATABASE_DBNAME=\"tiny_forum\"")
	log.Printf("%s", "\n"+strings.Repeat("=", 60))
}

// startServer 启动 HTTP 服务器
func startServer(cfg *config.Config, app *wire.App) {
	// 构建服务器地址
	addr := fmt.Sprintf(":%d", cfg.Basic.Server.Port)

	// 打印启动信息
	printStartupInfo(cfg, addr)

	// 启动服务器
	if err := app.Engine.Run(addr); err != nil {
		logger.Fatal(fmt.Sprintf("Server failed to start: %v", err))
	}
}

// printStartupInfo 打印启动信息
func printStartupInfo(cfg *config.Config, addr string) {
	logger.Info("========================================")
	logger.Info("🚀 Tiny Forum Server Starting...")
	logger.Info("========================================")
	logger.Info(fmt.Sprintf("🌍 Environment: %s", getEnvironment(cfg)))
	logger.Info(fmt.Sprintf("🔧 Server Address: http://localhost%s", addr))
	logger.Info("📚 API Base Path: /api/v1")
	logger.Info(fmt.Sprintf("🗄️  Database: %s@%s:%d/%s",
		cfg.Private.Database.User,
		cfg.Private.Database.Host,
		cfg.Private.Database.Port,
		cfg.Private.Database.DBName))

	if cfg.Private.Email.Host != "" {
		logger.Info(fmt.Sprintf("📧 Email Service: %s:%d", cfg.Private.Email.Host, cfg.Private.Email.Port))
	} else {
		logger.Info("📧 Email Service: Disabled")
	}

	if cfg.Basic.Redis.Host != "" {
		logger.Info(fmt.Sprintf("📡 Redis: %s:%d", cfg.Basic.Redis.Host, cfg.Basic.Redis.Port))
	}

	logger.Info("========================================")
}

// printConfigInfo 打印配置信息（调试用）
func printConfigInfo(cfg *config.Config) {
	logger.Debug("Configuration loaded:")
	logger.Debug(fmt.Sprintf("  Server Mode: %s", cfg.Basic.Server.Mode))
	logger.Debug(fmt.Sprintf("  Database: %s:%d/%s",
		cfg.Private.Database.Host,
		cfg.Private.Database.Port,
		cfg.Private.Database.DBName))
	logger.Debug(fmt.Sprintf("  Log Level: %s", cfg.Basic.Log.Level))

	if cfg.Private.JWT.Secret != "" {
		maskedSecret := maskString(cfg.Private.JWT.Secret, 4)
		logger.Debug(fmt.Sprintf("  JWT Secret: %s", maskedSecret))
		logger.Debug(fmt.Sprintf("  JWT Expire: %v", cfg.Private.JWT.Expire))
	}

	if cfg.Private.Email.Host != "" {
		logger.Debug(fmt.Sprintf("  Email: %s:%d", cfg.Private.Email.Host, cfg.Private.Email.Port))
	}
}

// maskString 隐藏敏感字符串
func maskString(s string, showChars int) string {
	if len(s) <= showChars {
		return "***"
	}
	return s[:showChars] + strings.Repeat("*", len(s)-showChars)
}

// getEnvironment 获取当前环境
func getEnvironment(cfg *config.Config) string {
	if cfg.IsProduction() {
		return "🚀 Production"
	}
	if cfg.IsDevelopment() {
		return "💻 Development"
	}
	return "❓ Unknown"
}

// fileExists 检查文件是否存在
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
