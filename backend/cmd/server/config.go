package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"tiny-forum/internal/infra/config"
)

// loadConfig 加载配置文件
func loadConfig(configDir string) (*config.Config, error) {
	// 获取当前工作目录
	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "获取当前目录失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("当前路径:", dir)

	// 构建配置文件路径
	basicConfigPath := filepath.Join(configDir, "basic.yml")
	privateConfigPath := filepath.Join(configDir, "private.yml")
	riskConfigPath := filepath.Join(configDir, "risk_control.yml")
	postgresPath := filepath.Join(configDir, "postgres.yml")
	redisPath := filepath.Join(configDir, "redis.yml")

	printConfigFileStatus(basicConfigPath, "Basic config")
	printConfigFileStatus(privateConfigPath, "Private config")
	printConfigFileStatus(riskConfigPath, "Risk control config")
	printConfigFileStatus(postgresPath, "Postgres config")
	printConfigFileStatus(redisPath, "Redis config")

	cfg, err := config.Load(configDir)
	if err != nil {
		return nil, fmt.Errorf("❌ failed to load config from '%s' directory. \nERROR: %s", configDir, err)
	}

	if err := validateConfigWithHints(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// printConfigFileStatus 打印配置文件状态
func printConfigFileStatus(filePath, configName string) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Printf("ℹ️  %s file not found: %s\n", configName, filePath)
		log.Printf("   → Using environment variables or defaults for this configuration\n")
	} else {
		log.Printf("✓ %s file found: %s\n", configName, filePath)
	}
}

// validateConfigWithHints 验证配置并提供修复提示
func validateConfigWithHints(cfg *config.Config) error {
	var errors []string
	var warnings []string

	// 服务器端口
	if cfg.Basic.Server.Port <= 0 || cfg.Basic.Server.Port > 65535 {
		errors = append(errors, fmt.Sprintf(
			"❌ Invalid server port: %d (must be between 1 and 65535)\n   → Set 'server.port' in config/basic.yaml or BASIC_SERVER_PORT environment variable",
			cfg.Basic.Server.Port))
	}

	// 数据库配置
	if cfg.Postgres.Host == "" {
		errors = append(errors, "❌ Database host is required\n   → Set 'database.host' in config/basic.yaml or BASIC_DATABASE_HOST environment variable")
	}
	if cfg.Postgres.DBName == "" {
		errors = append(errors, "❌ Database name is required\n   → Set 'database.dbname' in config/basic.yaml or BASIC_DATABASE_DBNAME environment variable")
	}

	// JWT 校验
	if cfg.Private.JWT.Secret == "" {
		envSecret := os.Getenv("BASIC_JWT_SECRET")
		if envSecret == "" {
			errors = append(errors, "❌ JWT secret is required for security\n   → Set 'jwt.secret' in config/basic.yaml\n   → Or set BASIC_JWT_SECRET environment variable")
		} else {
			warnings = append(warnings, "⚠️  JWT secret is set via environment variable (this is fine)")
		}
	} else if len(cfg.Private.JWT.Secret) < 32 && cfg.IsProduction() {
		warnings = append(warnings, fmt.Sprintf("⚠️  JWT secret is too short (%d characters) for production", len(cfg.Private.JWT.Secret)))
	}

	// 邮件配置（可选）
	if cfg.Private.Email.Host != "" && cfg.Private.Email.Port <= 0 {
		warnings = append(warnings, "⚠️  Invalid email port, email features may not work")
	}

	// 输出警告
	for _, warning := range warnings {
		log.Printf("%s", warning)
	}

	if len(errors) > 0 {
		return fmt.Errorf("\n%s\n\n💡 Quick Fix:\n  1. Create config/basic.yaml with required settings\n  2. Or set environment variables (see above)\n  3. Run 'make init-dev' to generate sample config",
			strings.Join(errors, "\n"))
	}
	return nil
}

// printConfigError 打印配置错误
func printConfigError(err error) {
	log.Printf("%s", strings.Repeat("=", 60))
	log.Printf("🚫 Configuration Error")
	log.Printf("%s", strings.Repeat("=", 60))
	log.Printf("\n%v\n", err)
	log.Printf("\nPlease run 'make init-dev' to initialize configuration files.\n")
}
