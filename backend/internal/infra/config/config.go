package config

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Load 加载配置文件
func Load(configDir string) (*Config, error) {
	// 定义需要加载的配置文件列表
	configFiles := []string{"basic", "private", "risk_control", "postgres", "redis", "ai", "app"}
	vipers := newViperInstances(configDir, configFiles...)

	// 基本配置
	var basicConfig ConfigBasic
	if err := vipers["basic"].UnmarshalExact(&basicConfig); err != nil {
		return nil, fmt.Errorf("加载基础配置失败: %w", err)
	}

	// 私有配置
	var privateConfig ConfigPrivate
	if err := vipers["private"].UnmarshalExact(&privateConfig); err != nil {
		fmt.Printf("⚠️  加载私有配置失败: %v，使用空配置\n", err)
		// privateConfig = ConfigPrivate{}
		return nil, fmt.Errorf("加载私有配置失败: %w", err)
	}

	// 风控配置
	var riskConfig ConfigRiskControl
	if err := vipers["risk_control"].UnmarshalExact(&riskConfig); err != nil {
		fmt.Printf("⚠️  加载风控配置失败: %v，使用空配置\n", err)
		return nil, fmt.Errorf("加载风控配置失败: %w", err)
	}

	// Postgres 配置
	var postgresConfig ConfigPostgres
	if err := vipers["postgres"].UnmarshalExact(&postgresConfig); err != nil {
		fmt.Printf("⚠️  加载 Postgres 配置失败: %v，使用空配置\n", err)
		return nil, fmt.Errorf("加载 Postgres 配置失败: %w", err)
	}

	// Redis 配置
	var redisConfig ConfigRedis
	if err := vipers["redis"].UnmarshalExact(&redisConfig); err != nil {
		fmt.Printf("⚠️  加载 Redis 配置失败: %v，使用空配置\n", err)
		return nil, fmt.Errorf("加载 Redis 配置失败: %w", err)
	}

	// AI 配置
	var aiConfig ConfigAI
	if err := vipers["ai"].UnmarshalExact(&aiConfig); err != nil {
		log.Printf("⚠️  加载 AI 配置失败: %v，使用空配置\n", err)
		return nil, fmt.Errorf("加载 AI 配置失败: %w", err)
	}
	// app 配置
	var appConfig ConfigApp
	if err := vipers["app"].UnmarshalExact(&appConfig); err != nil {
		log.Printf("⚠️  加载 app 配置失败: %v，使用空配置\n", err)
		return nil, fmt.Errorf("加载 app 配置失败: %w", err)
	}

	cfg := &Config{
		Basic:       basicConfig,
		Private:     privateConfig,
		RiskControl: riskConfig,
		Postgres:    postgresConfig,
		Redis:       redisConfig,
		AI:          aiConfig,
		App:         appConfig,
	}

	fmt.Printf("✅ 所有配置文件成功加载\n")
	fmt.Printf("✅ 配置信息: %s\n", cfg)
	cfg.setDefaults()
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("配置验证失败: %w", err)
	}
	return cfg, nil
}

// newViperInstances 创建并配置viper实例
func newViperInstances(configDir string, fileNames ...string) map[string]*viper.Viper {
	instances := make(map[string]*viper.Viper)
	for _, name := range fileNames {
		configPath := filepath.Join(configDir, name+".yml")
		instances[name] = newViper(strings.ToUpper(name), configPath)
	}
	return instances
}

// newViper 创建配置好的viper实例
func newViper(prefix, configPath string) *viper.Viper {
	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigFile(configPath)

	// 环境变量支持
	v.AutomaticEnv()
	v.SetEnvPrefix(prefix)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 读取配置文件（如果不存在也不报错，因为可能通过环境变量提供）
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 配置文件不存在，使用环境变量
			fmt.Printf("ℹ️  配置文件 %s 不存在，将使用环境变量\n", configPath)
		} else {
			fmt.Printf("⚠️  读取配置文件 %s 出错: %v\n", configPath, err)
		}
	}

	return v
}

// // LoadWithOverrides 加载配置并支持覆盖
// func LoadWithOverrides(configDir string, overrides map[string]interface{}) (*Config, error) {
// 	cfg, err := Load(configDir)
// 	if err != nil {
// 		return nil, err
// 	}

// 	applyOverrides(cfg, overrides)
// 	return cfg, nil
// }

// // applyOverrides 应用配置覆盖
// func applyOverrides(cfg *Config, overrides map[string]interface{}) {
// 	for key, value := range overrides {
// 		switch key {
// 		case "server.port":
// 			if v, ok := value.(int); ok {
// 				cfg.Basic.Server.Port = v
// 			}
// 		case "database.host":
// 			if v, ok := value.(string); ok {
// 				cfg.Postgres.Host = v
// 			}
// 			// 可扩展更多覆盖项
// 		}
// 	}
// }

// setDefaults 设置默认值
func (c *Config) setDefaults() {
	c.setServerDefaults()
	c.setDatabaseDefaults()
	c.setJWTDefaults()
	c.setLogDefaults()
	c.setEmailDefaults()
}

func (c *Config) setServerDefaults() {
	fmt.Println("设置服务默认值")
	if c.Basic.Server.Port == 0 {
		c.Basic.Server.Port = 8080
	}
	if c.Basic.Server.Mode == "" {
		c.Basic.Server.Mode = "debug"
	}
	if c.Basic.Server.ReadTimeout == 0 {
		c.Basic.Server.ReadTimeout = 30 * time.Second
	}
	if c.Basic.Server.WriteTimeout == 0 {
		c.Basic.Server.WriteTimeout = 30 * time.Second
	}
}

func (c *Config) setDatabaseDefaults() {
	fmt.Println("设置数据库默认值")
	if c.Postgres.SSLMode == "" {
		c.Postgres.SSLMode = "disable"
	}
	if c.Postgres.TimeZone == "" {
		c.Postgres.TimeZone = "Asia/Shanghai"
	}
	if c.Postgres.MaxIdleConns == 0 {
		c.Postgres.MaxIdleConns = 10
	}
	if c.Postgres.MaxOpenConns == 0 {
		c.Postgres.MaxOpenConns = 100
	}
}

func (c *Config) setJWTDefaults() {
	fmt.Println("设置 JWT 默认值")
	if c.Private.JWT.Expire == 0 {
		c.Private.JWT.Expire = 24 * time.Hour
	}
	if c.Private.JWT.RefreshExpire == 0 {
		c.Private.JWT.RefreshExpire = 7 * 24 * time.Hour
	}
}

func (c *Config) setLogDefaults() {
	fmt.Println("设置日志默认值")
	if c.Basic.Log.Level == "" {
		c.Basic.Log.Level = "info"
	}
	if c.Basic.Log.MaxSize == 0 {
		c.Basic.Log.MaxSize = 100
	}
	if c.Basic.Log.MaxBackups == 0 {
		c.Basic.Log.MaxBackups = 10
	}
	if c.Basic.Log.MaxAge == 0 {
		c.Basic.Log.MaxAge = 30
	}
}

func (c *Config) setEmailDefaults() {
	fmt.Println("设置邮箱默认值")
	if c.Private.Email.PoolSize == 0 {
		c.Private.Email.PoolSize = 5
	}
}

// // validate 验证配置
// // func (c *Config) validate() error {
// // 	validators := []func() error{
// // 		c.validateJWT,
// // 		c.validateDatabase,
// // 	}

// // 	for _, validator := range validators {
// // 		if err := validator(); err != nil {
// // 			return err
// // 		}
// // 	}
// // 	return nil
// // }

// // validateJWT 验证JWT配置
// // func (c *Config) validateJWT() error {
// // 	// JWT secret 必须存在
// // 	if c.Private.JWT.Secret == "" {
// // 		// 检查是否通过环境变量设置（这里只是提示，实际环境变量已经在viper中读取）
// // 		return &ConfigError{
// // 			Field:   "jwt.secret",
// // 			Message: "JWT secret is required for security. Please set it in config/basic.yaml or via BASIC_JWT_SECRET environment variable",
// // 		}
// // 	}
// // 	return nil
// // }

// func (c *Config) validateDatabase() error {
// 	if c.Postgres.Host == "" {
// 		return &ConfigError{Field: "database.host", Message: "database host is required"}
// 	}
// 	if c.Postgres.User == "" {
// 		return &ConfigError{Field: "database.user", Message: "database user is required"}
// 	}
// 	if c.Postgres.DBName == "" {
// 		return &ConfigError{Field: "database.dbname", Message: "database name is required"}
// 	}
// 	return nil
// }

// ToLoggerConfig 转换为日志配置
func (c *Config) ToLoggerConfig() LogConfig {
	cfg := LogConfig{
		Level:      c.Basic.Log.Level,
		Filename:   c.Basic.Log.Filename,
		MaxSize:    c.Basic.Log.MaxSize,
		MaxBackups: c.Basic.Log.MaxBackups,
		MaxAge:     c.Basic.Log.MaxAge,
		Compress:   c.Basic.Log.Compress,
		Console:    c.Basic.Log.Console,
		JSONFormat: c.Basic.Log.JSONFormat,
	}
	if c.Basic.Log.DB != nil && c.Basic.Log.DB.DSN != "" {

		cfg.DB = &DBConfig{
			DSN:        c.Basic.Log.DB.DSN,
			MaxBuffer:  c.Basic.Log.DB.MaxBuffer,
			BatchSize:  c.Basic.Log.DB.BatchSize,
			FlushEvery: c.Basic.Log.DB.FlushEvery,
			Retention:  c.Basic.Log.DB.Retention,
		}
	}
	return cfg
}
