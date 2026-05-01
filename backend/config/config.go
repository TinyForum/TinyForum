package config

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"tiny-forum/internal/model"

	"github.com/spf13/viper"
)

// MARK: Config 一级配置结构
type Config struct {
	Basic       ConfigBasic       //基础配置
	Private     ConfigPrivate     // 私有配置
	RiskControl ConfigRiskControl // 风控配置
}

// =========================================== MARK: 风控配置
type ConfigRiskControl struct {
	RateLimit RateLimitConfig `mapstructure:"rate_limit"` // 风控
}
type RateLimitConfig struct {
	RiskLevels  map[string]map[string]QuotaConfig `yaml:"risk_levels"` // 风控等级
	Enabled     bool                              `yaml:"enabled" json:"enabled"`
	IPWhitelist []string                          `yaml:"ip_whitelist" json:"ip_whitelist"`
}

type QuotaConfig struct {
	Limit  int    `yaml:"limit"`
	Window string `yaml:"window"` // 如 "1h"
}

// =========================================== MARK: 基础配置
type ConfigBasic struct {
	Server       ServerConfig    `mapstructure:"server"`
	Frontend     FrontendConfig  `mapstructure:"frontend"`
	API          APIConfig       `mapstructure:"api"`
	Log          LogConfig       `mapstructure:"log"`
	RateLimit    RateLimitConfig `mapstructure:"rate_limit"`
	Ollama       Ollama          `mapstructure:"ollama"`
	AllowOrigins []string        `mapstructure:"allow_origins"`
	Upload       UploadConfig    `mapstructure:"upload"`
}
type FrontendConfig struct {
	Protocol string `mapstructure:"protocol"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
}

type UploadConfig struct {
	UploadDir  string
	URLPrefix  string
	MaxSize    int64
	AllowedExt []string
}

type Ollama struct {
	BaseURL string `mapstructure:"base_url"`
	Model   string `mapstructure:"model"`
	APIKey  string `mapstructure:"api_key"`
	Timeout uint   `mapstructure:"timeout"`
}

type APIConfig struct {
	Protocol string `mapstructure:"protocol"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Version  string `mapstructure:"version"`
	Prefix   string `mapstructure:"prefix"`
}

// =========================================== MARK: 私有配置
type ConfigPrivate struct {
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Email    EmailConfig    `mapstructure:"email"`
	OAuth    OAuthConfig    `mapstructure:"oauth"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Admin    AdminConfig    `mapstructure:"admin"`
}

// AdminConfig 管理员配置
type AdminConfig struct {
	Username string         `mapstructure:"username"`
	Email    string         `mapstructure:"email"`
	Password string         `mapstructure:"password"`
	Role     model.UserRole `mapstructure:"role"`
	Score    int            `mapstructure:"score"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host           string        `mapstructure:"host"`
	Port           int           `mapstructure:"port"`
	Mode           string        `mapstructure:"mode"`
	ReadTimeout    time.Duration `mapstructure:"read_timeout"`
	WriteTimeout   time.Duration `mapstructure:"write_timeout"`
	MaxHeaderBytes int           `mapstructure:"max_header_bytes"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	DBName          string        `mapstructure:"dbname"`
	SSLMode         string        `mapstructure:"sslmode"`
	TimeZone        string        `mapstructure:"timezone"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret        string        `mapstructure:"secret"`
	Expire        time.Duration `mapstructure:"expire"`
	RefreshExpire time.Duration `mapstructure:"refresh_expire"`
	Issuer        string        `mapstructure:"issuer"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
	Console    bool   `mapstructure:"console"`
	JSONFormat bool   `mapstructure:"json_format"`
	// DB 可选。非零值时自动初始化 SQLite 日志数据库并接入日志管道。
	// 若已单独调用过 InitDB，此字段可留空（不会重复初始化）。
	DB *DBConfig `mapstructure:"db"`
}

// DBConfig SQLite 日志数据库配置
type DBConfig struct {
	// DSN 数据库文件路径，例如 "./logs/app.db"
	DSN string `mapstructure:"dsn"`

	// MaxBuffer 异步写入队列深度，默认 512
	// 队列满时新日志静默丢弃，不阻塞业务
	MaxBuffer int `mapstructure:"max_buffer"`

	// BatchSize 单次事务批量写入条数，默认 50
	BatchSize int `mapstructure:"batch_size"`

	// FlushEvery 定时强制刷盘间隔，默认 2s
	FlushEvery time.Duration `mapstructure:"flush_every"`

	// Retention 日志保留天数，0 表示永久保留
	// 到期表以 DROP TABLE 整表删除，比按行 DELETE 快得多
	Retention int `mapstructure:"retention"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	DB       int    `mapstructure:"db"`
	Password string `mapstructure:"password"` // 合并密码到主配置
	PoolSize int    `mapstructure:"pool_size"`
}

// RedisPrivateConfig 已废弃，密码已合并到 RedisConfig
// type RedisPrivateConfig struct {
// 	Password string `mapstructure:"password"`
// }

// EmailConfig 邮件配置
type EmailConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	From     string `mapstructure:"from"`
	FromName string `mapstructure:"from_name"`
	SSL      bool   `mapstructure:"ssl"`
	TLS      bool   `mapstructure:"tls"`
	PoolSize int    `mapstructure:"pool_size"`
}

// OAuthConfig OAuth配置
type OAuthConfig struct {
	Github GithubOAuthConfig `mapstructure:"github"`
	Google GoogleOAuthConfig `mapstructure:"google"`
	Wechat WechatOAuthConfig `mapstructure:"wechat"`
}

type GithubOAuthConfig struct {
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	RedirectURL  string `mapstructure:"redirect_url"`
}

type GoogleOAuthConfig struct {
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	RedirectURL  string `mapstructure:"redirect_url"`
}

type WechatOAuthConfig struct {
	AppID       string `mapstructure:"app_id"`
	AppSecret   string `mapstructure:"app_secret"`
	RedirectURL string `mapstructure:"redirect_url"`
}

// UploadConfig 文件上传配置
// type UploadConfig struct {
// 	MaxSize        int64    `mapstructure:"max_size"`
// 	AllowedTypes   []string `mapstructure:"allowed_types"`
// 	StoragePath    string   `mapstructure:"storage_path"`
// 	CDNDomain      string   `mapstructure:"cdn_domain"`
// 	EnableCompress bool     `mapstructure:"enable_compress"`
// }

// RateLimitConfig 限流配置
// type RateLimitConfig struct {
// 	Enabled  bool `mapstructure:"enabled"`
// 	Requests int  `mapstructure:"requests"`
// 	Duration int  `mapstructure:"duration"`
// 	Burst    int  `mapstructure:"burst"`
// }

// Load 加载配置文件
func Load(configDir string) (*Config, error) {
	basicViper, privateViper := newViperInstances(configDir)
	fmt.Println("开始加载配置文件\n")

	// 加载基础配置
	var basicConfig ConfigBasic
	if err := basicViper.Unmarshal(&basicConfig); err != nil {
		fmt.Printf("加载基础配置文件失败: %v\n", err)
		return nil, err
	}
	fmt.Printf("加载基础配置文件成功\n")

	// 加载私有配置
	var privateConfig ConfigPrivate
	if err := privateViper.Unmarshal(&privateConfig); err != nil {
		fmt.Printf("加载私有配置文件失败: %v\n", err)
		privateConfig = ConfigPrivate{}
	}
	fmt.Printf("加载私有配置文件成功\n")

	// 加载风控配置（独立文件）
	riskViper := viper.New()
	riskViper.SetConfigType("yaml")
	riskViper.SetConfigFile(filepath.Join(configDir, "risk_control.yaml"))
	_ = riskViper.ReadInConfig() // 允许文件不存在

	fmt.Printf("加载风控配置文件成功\n")
	var riskControl ConfigRiskControl
	if err := riskViper.Unmarshal(&riskControl); err != nil {
		// 解析失败则使用空配置（不阻塞启动）
		riskControl = ConfigRiskControl{}
	}

	cfg := &Config{
		Basic:       basicConfig,
		Private:     privateConfig,
		RiskControl: riskControl,
	}
	fmt.Println("加载配置文件成功\n")
	cfg.setDefaults()
	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

// newViperInstances 创建并配置viper实例
func newViperInstances(configDir string) (*viper.Viper, *viper.Viper) {
	basicViper := newViper("BASIC", filepath.Join(configDir, "basic.yaml"))
	privateViper := newViper("PRIVATE", filepath.Join(configDir, "private.yaml"))
	return basicViper, privateViper
}

// newViper 创建配置好的viper实例
func newViper(prefix, configPath string) *viper.Viper {
	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigFile(configPath)
	v.AutomaticEnv()
	v.SetEnvPrefix(prefix)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	_ = v.ReadInConfig()
	return v
}

// LoadWithOverrides 加载配置并支持覆盖
func LoadWithOverrides(configDir string, overrides map[string]interface{}) (*Config, error) {
	cfg, err := Load(configDir)
	if err != nil {
		return nil, err
	}

	applyOverrides(cfg, overrides)
	return cfg, nil
}

// applyOverrides 应用配置覆盖
func applyOverrides(cfg *Config, overrides map[string]interface{}) {
	for key, value := range overrides {
		switch key {
		case "server.port":
			if v, ok := value.(int); ok {
				cfg.Basic.Server.Port = v
			}
		case "database.host":
			if v, ok := value.(string); ok {
				cfg.Private.Database.Host = v
			}
			// 可扩展更多覆盖项
		}
	}
}

// setDefaults 设置默认值
func (c *Config) setDefaults() {
	c.setServerDefaults()
	c.setDatabaseDefaults()
	c.setJWTDefaults()
	c.setLogDefaults()
	c.setEmailDefaults()
}

func (c *Config) setServerDefaults() {
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
	if c.Private.Database.SSLMode == "" {
		c.Private.Database.SSLMode = "disable"
	}
	if c.Private.Database.TimeZone == "" {
		c.Private.Database.TimeZone = "Asia/Shanghai"
	}
	if c.Private.Database.MaxIdleConns == 0 {
		c.Private.Database.MaxIdleConns = 10
	}
	if c.Private.Database.MaxOpenConns == 0 {
		c.Private.Database.MaxOpenConns = 100
	}
}

func (c *Config) setJWTDefaults() {
	if c.Private.JWT.Expire == 0 {
		c.Private.JWT.Expire = 24 * time.Hour
	}
	if c.Private.JWT.RefreshExpire == 0 {
		c.Private.JWT.RefreshExpire = 7 * 24 * time.Hour
	}
}

func (c *Config) setLogDefaults() {
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
	if c.Private.Email.PoolSize == 0 {
		c.Private.Email.PoolSize = 5
	}
}

// validate 验证配置
func (c *Config) validate() error {
	validators := []func() error{
		c.validateJWT,
		c.validateDatabase,
	}

	for _, validator := range validators {
		if err := validator(); err != nil {
			return err
		}
	}
	return nil
}

// validateJWT 验证JWT配置
func (c *Config) validateJWT() error {
	// JWT secret 必须存在
	if c.Private.JWT.Secret == "" {
		// 检查是否通过环境变量设置（这里只是提示，实际环境变量已经在viper中读取）
		return &ConfigError{
			Field:   "jwt.secret",
			Message: "JWT secret is required for security. Please set it in config/basic.yaml or via BASIC_JWT_SECRET environment variable",
		}
	}
	return nil
}

func (c *Config) validateDatabase() error {
	if c.Private.Database.Host == "" {
		return &ConfigError{Field: "database.host", Message: "database host is required"}
	}
	if c.Private.Database.User == "" {
		return &ConfigError{Field: "database.user", Message: "database user is required"}
	}
	if c.Private.Database.DBName == "" {
		return &ConfigError{Field: "database.dbname", Message: "database name is required"}
	}
	return nil
}

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

// GetDSN 获取数据库连接字符串
func (c *Config) GetDSN() string {
	return c.Private.Database.GetDSN()
}

// GetRedisAddr 获取Redis地址
func (c *Config) GetRedisAddr() string {
	return c.Private.Redis.GetAddr()
}

// IsProduction 是否为生产环境
func (c *Config) IsProduction() bool {
	return c.Basic.Server.Mode == "release"
}

// IsDevelopment 是否为开发环境
func (c *Config) IsDevelopment() bool {
	return c.Basic.Server.Mode == "debug"
}

// GetDSN 获取数据库连接字符串
func (d *DatabaseConfig) GetDSN() string {
	// TODO: 实现完整的DSN构建
	return d.Host
}

// GetAddr 获取Redis地址
func (r *RedisConfig) GetAddr() string {
	addr := r.Host + ":" + strconv.Itoa(r.Port)
	return addr
}

// ConfigError 配置错误
type ConfigError struct {
	Field   string
	Message string
}

func (e *ConfigError) Error() string {
	return e.Message
}
