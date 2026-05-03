package config

import (
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/go-playground/validator/v10"
)

// =========================================== MARK: 辅助方法

// GetServerPort returns the server port.
func (c *Config) GetServerPort() int {
	return c.Basic.Server.Port
}

// GetServerMode returns the server mode (debug/release/test).
func (c *Config) GetServerMode() string {
	return c.Basic.Server.Mode
}

// IsProduction returns true if mode is "release" or "production".
func (c *Config) IsProduction() bool {
	mode := c.Basic.Server.Mode
	return mode == "release" || mode == "production"
}

// GetDatabaseInfo returns host, port, user, database name from Postgres config.
func (c *Config) GetDatabaseInfo() (host string, port int, user, dbName string) {
	return c.Postgres.Host, c.Postgres.Port, c.Postgres.User, c.Postgres.DBName
}

// GetEmailInfo returns host, port, and whether email is enabled.
func (c *Config) GetEmailInfo() (host string, port int, enabled bool) {
	return c.Private.Email.Host, c.Private.Email.Port, c.Private.Email.Host != ""
}

// GetRedisInfo returns host, port, and whether redis is enabled.
func (c *Config) GetRedisInfo() (host string, port int, enabled bool) {
	return c.Redis.Host, c.Redis.Port, c.Redis.Host != ""
}

// GetOllamaInfo returns base URL and model name.
func (c *Config) GetOllamaInfo() (baseURL, model string) {
	return c.Basic.Ollama.BaseURL, c.Basic.Ollama.Model
}

// GetJWTSecretMasked returns a masked version of JWT secret.
func (c *Config) GetJWTSecretMasked() string {
	secret := c.Private.JWT.Secret
	if len(secret) <= 4 {
		return "***"
	}
	return secret[:4] + "****"
}

// GetLogLevel returns the configured log level.
func (c *Config) GetLogLevel() string {
	return c.Basic.Log.Level
}

// GetDSN 获取 PostgreSQL 数据库连接字符串
func (c *Config) GetDSN() string {
	return c.Postgres.GetDSN()
}

// GetRedisAddr 获取 Redis 地址（host:port）
func (c *Config) GetRedisAddr() string {
	return c.Redis.GetAddr()
}

// IsDevelopment 是否为开发环境（debug 模式）
func (c *Config) IsDevelopment() bool {
	return c.Basic.Server.Mode == "debug"
}

// =========================================== MARK: 结构体自身方法

// GetDSN 获取 PostgreSQL 的完整连接字符串
func (p *ConfigPostgres) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		p.Host, p.Port, p.User, p.Password, p.DBName, p.SSLMode, p.TimeZone)
}

// GetAddr 获取 PostgreSQL 的主机地址（host:port），用于某些需要直接连接的场景
func (p *ConfigPostgres) GetAddr() string {
	return p.Host + ":" + strconv.Itoa(p.Port)
}

// GetAddr 获取 Redis 的主机地址（host:port）
func (r *ConfigRedis) GetAddr() string {
	redisAddr := r.Host + ":" + strconv.Itoa(r.Port)
	log.Printf("config: %s, %s, %s, %s", r.Host, r.Port, r.Password, r.DB)
	log.Printf("connect Redis address: %s", redisAddr)
	return redisAddr
}

// =========================================== MARK: 配置验证

// ConfigError 自定义错误类型（可选）
type ConfigError struct {
	Field   string
	Message string
}

func (e *ConfigError) Error() string {
	return e.Message
}

var validate *validator.Validate

func init() {
	validate = validator.New()
	// 注册自定义 regexp 验证器
	validate.RegisterValidation("regexp", func(fl validator.FieldLevel) bool {
		pattern := fl.Param()
		if pattern == "" {
			return false
		}
		matched, err := regexp.MatchString(pattern, fl.Field().String())
		if err != nil {
			return false
		}
		return matched
	})
}

// validate 执行配置的递归验证
func (c *Config) validate() error {
	return validate.Struct(c)
}
