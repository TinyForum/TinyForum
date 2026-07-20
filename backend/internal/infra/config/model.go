package config

import (
	"time"
	"tiny-forum/internal/model/do"
)

// MARK: Config 一级配置结构
type Config struct {
	Basic       ConfigBasic       // 基础配置文件
	Private     ConfigPrivate     // 私有配置文件
	RiskControl ConfigRiskControl // 风控配置文件
	Postgres    ConfigPostgres    // 数据库配置文件
	Redis       ConfigRedis       // Redis 配置文件
	AI          ConfigAI          // AI 配置文件
	App         ConfigApp         // 应用配置文件
}

// MARK: 应用配置
type ConfigApp struct {
	// 是否启用调试模式
	Debug bool `mapstructure:"debug" validate:"required"`
}

// MARK: 基础配置集合
type ConfigBasic struct {
	Server       ServerConfig   `mapstructure:"server" validate:"required"` // 服务配置
	Frontend     FrontendConfig `mapstructure:"frontend" validate:"required"`
	API          APIConfig      `mapstructure:"api" validate:"required"`
	Log          LogConfig      `mapstructure:"log" validate:"required"`
	AllowOrigins []string       `mapstructure:"allow_origins" validate:"omitempty,dive,url"` // 允许跨域请求的域名
	Attachment   UploadConfig   `mapstructure:"attachment" validate:"required"`              // 上传配置
	Version      string         `mapstructure:"version" validate:"required,semver"`
	Plugins      ConfigPlugins  `mapstructure:"plugins" validate:"required"`
}

type ConfigPlugins struct {
	StorageDir string `mapstructure:"storage_dir" validate:"required"`
}

type ConfigPrivate struct {
	JWT        JWTConfig        `mapstructure:"jwt" validate:"required"`
	Email      EmailConfig      `mapstructure:"email" validate:"required"`
	OAuth      OAuthConfig      `mapstructure:"oauth" validate:"required"`
	AdminUser  AdminUserConfig  `mapstructure:"admin" validate:"required"`
	SystemUser SystemUserConfig `mapstructure:"system" validate:"required"`
}
type RateLimitConfig struct {
	RiskControlLevels map[string]map[string]QuotaConfig `yaml:"risk_control_levels" json:"risk_control_levels" mapstructure:"risk_control_levels" validate:"required,min=1"`
	Enabled           bool                              `yaml:"enabled" json:"enabled" mapstructure:"enabled"`
	IPWhitelist       []string                          `yaml:"ip_whitelist" json:"ip_whitelist" mapstructure:"ip_whitelist" validate:"omitempty,dive,cidr"`
}
type ContentCheckConfig struct {
	Enabled bool `yaml:"enabled" json:"enabled" mapstructure:"enabled"`
}

// =========================================== MARK: AI
// ConfigAI 是 AI 配置的根结构
type ConfigAI struct {
	// 是否启用 AI 模块
	Enable bool `yaml:"enable" mapstructure:"enable" json:"enable" validate:"required"` // 是否启用 AI 模块 （默认启用）
	// 当前激活的服务提供商
	Provider string `mapstructure:"provider" validate:"omitempty"` // 当前激活的服务提供商（默认为 openai）

	// 当前提供商的详细配置（所有提供商结构一致，兼容 OpenAI 格式）
	Config ProviderConfig `mapstructure:"config" validate:"omitempty"` // 当前提供商的详细配置（所有提供商结构一致，兼容 OpenAI 格式）

	// 全局通用配置
	Defaults Defaults    `mapstructure:"defaults" validate:"omitempty"`
	Retry    RetryConfig `mapstructure:"retry" validate:"omitempty"`
	Logging  Logging     `mapstructure:"logging" validate:"omitempty"`
	Cache    CacheConfig `mapstructure:"cache" validate:"omitempty"`
}

// ProviderConfig 适用于所有 OpenAI 兼容服务
type ProviderConfig struct {
	APIKey           string        `mapstructure:"api-key" validate:"omitempty"`
	BaseURL          string        `mapstructure:"base-url" validate:"omitempty"`
	Model            string        `mapstructure:"model" validate:"omitempty"`
	Timeout          time.Duration `mapstructure:"timeout" validate:"omitempty"`
	MaxRetries       int           `mapstructure:"max-retries" validate:"omitempty"`
	Temperature      float64       `mapstructure:"temperature" validate:"omitempty"`
	MaxTokens        int           `mapstructure:"max-tokens" validate:"omitempty"`
	TopP             float64       `mapstructure:"top-p" validate:"omitempty"`
	FrequencyPenalty float64       `mapstructure:"frequency-penalty" validate:"omitempty"`
	PresencePenalty  float64       `mapstructure:"presence-penalty" validate:"omitempty"`
	EmbeddingModel   string        `mapstructure:"embedding-model" validate:"omitempty"`
	ModerationModel  string        `mapstructure:"moderation-model" validate:"omitempty"`
}

type Defaults struct {
	Timeout          int     `mapstructure:"timeout" validate:"omitempty"`
	MaxRetries       int     `mapstructure:"max-retries" validate:"omitempty"`
	Temperature      float64 `mapstructure:"temperature" validate:"omitempty"`
	MaxTokens        int     `mapstructure:"max-tokens" validate:"omitempty"`
	TopP             float64 `mapstructure:"top-p" validate:"omitempty"`
	FrequencyPenalty float64 `mapstructure:"frequency-penalty" validate:"omitempty"`
	PresencePenalty  float64 `mapstructure:"presence-penalty" validate:"omitempty"`
}

type RetryConfig struct {
	Enabled     bool          `mapstructure:"enabled" validate:"omitempty"`
	MaxAttempts int           `mapstructure:"max-attempts" validate:"omitempty"`
	Backoff     BackoffConfig `mapstructure:"backoff" validate:"omitempty"`
}

type BackoffConfig struct {
	InitialInterval int     `mapstructure:"initial-interval" validate:"omitempty"`
	Multiplier      float64 `mapstructure:"multiplier" validate:"omitempty"`
	MaxInterval     int     `mapstructure:"max-interval" validate:"omitempty"`
}

type Logging struct {
	Enabled    bool   `mapstructure:"enabled" validate:"omitempty"`
	Level      string `mapstructure:"level" validate:"omitempty"`
	LogPayload bool   `mapstructure:"log-payload" validate:"omitempty"`
}

type CacheConfig struct {
	Enabled bool `mapstructure:"enabled" validate:"omitempty"`
	TTL     int  `mapstructure:"ttl" validate:"omitempty"`
	MaxSize int  `mapstructure:"max-size" validate:"omitempty"`
}

// =========================================== MARK: 数据库 / 存储

type ConfigPostgres struct {
	Host            string        `mapstructure:"host" validate:"required"`
	Port            int           `mapstructure:"port" validate:"required,min=1,max=65535"`
	User            string        `mapstructure:"user" validate:"required"`
	Password        string        `mapstructure:"password" validate:"required"`
	DBName          string        `mapstructure:"dbname" validate:"required"`
	SSLMode         string        `mapstructure:"sslmode" validate:"omitempty,oneof=disable allow prefer require verify-ca verify-full"`
	TimeZone        string        `mapstructure:"timezone" validate:"omitempty"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns" validate:"omitempty,min=0"`
	MaxOpenConns    int           `mapstructure:"max_open_conns" validate:"omitempty,min=0"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime" validate:"omitempty,min=0"`
	Logger          *PostLogger   `mapstructure:"logger"` // 日志配置
}

type PostLogger struct {
	SlowThreshold             string `mapstructure:"slow_threshold" validate:"omitempty"`                        // 如 "200ms"
	LogLevel                  string `mapstructure:"log_level" validate:"required,oneof=silent error warn info"` // silent, error, warn, info
	IgnoreRecordNotFoundError bool   `mapstructure:"ignore_record_not_found_error"`
	Colorful                  bool   `mapstructure:"colorful"`
}

type ConfigRedis struct {
	Host     string `mapstructure:"host" validate:"required,hostname|ip"`
	Port     int    `mapstructure:"port" validate:"required,min=1,max=65535"`
	DB       int    `mapstructure:"db" validate:"min=0"`
	Password string `mapstructure:"password"`
	PoolSize int    `mapstructure:"pool_size" validate:"omitempty,min=1"`
}

// =========================================== MARK: 风控配置

type ConfigRiskControl struct {
	RateLimit    RateLimitConfig    `mapstructure:"rate_limit" validate:"required"`
	ContentCheck ContentCheckConfig `mapstructure:"content_check" validate:"required"`
}

type QuotaConfig struct {
	Limit  int    `yaml:"limit" validate:"required,min=1"`
	Window string `yaml:"window" validate:"required,regexp=^[0-9]+[smhdw]?$"` // 简单正则匹配如 "1h", "30m", "7d"
}

// =========================================== MARK: 基础子配置

type FrontendConfig struct {
	Protocol string `mapstructure:"protocol" validate:"required,oneof=http https"`
	Host     string `mapstructure:"host" validate:"required,hostname|ip"`
	Port     int    `mapstructure:"port" validate:"required,min=1,max=65535"`
}

type UploadConfig struct {
	UploadDir  string   `mapstructure:"upload_dir" validate:"required"`
	URLPrefix  string   `mapstructure:"url_prefix" validate:"required"`
	MaxSize    int64    `mapstructure:"max_size" validate:"required,min=1"`
	AllowedExt []string `mapstructure:"allowed_ext" validate:"required,min=1,dive,required"`
}

type Ollama struct {
	BaseURL string `mapstructure:"base_url" validate:"required,url"`
	Model   string `mapstructure:"model" validate:"required"`
	APIKey  string `mapstructure:"api_key"`
	Timeout uint   `mapstructure:"timeout" validate:"omitempty,min=1"`
}

type APIConfig struct {
	Protocol string `mapstructure:"protocol" validate:"required,oneof=http https"`     // 协议：	http 或 https
	Host     string `mapstructure:"host" validate:"required,hostname|ip"`              // 主机名或 IP 地址
	Port     int    `mapstructure:"port" validate:"required,min=1,max=65535"`          // API 端口
	Version  string `mapstructure:"version" validate:"required,regexp=^v[1-9][0-9]*$"` // 简单正则匹配如 "v1", "v2", "v3"
	Prefix   string `mapstructure:"prefix" validate:"required,startswith=/"`           // 必须以 / 开头
}

// =========================================== MARK: 私有子配置

type AdminUserConfig struct {
	Username string      `mapstructure:"username" validate:"required,min=3,max=32"`
	Email    string      `mapstructure:"email" validate:"required,email"`
	Password string      `mapstructure:"password" validate:"required,min=8"`
	Role     do.UserRole `mapstructure:"role" validate:"required"`
	Score    int         `mapstructure:"score" validate:"min=0"`
}

type SystemUserConfig struct {
	Username string      `mapstructure:"username" validate:"required,min=3,max=32"`
	Email    string      `mapstructure:"email" validate:"required,email"`
	Password string      `mapstructure:"password" validate:"required,min=8"`
	Role     do.UserRole `mapstructure:"role" validate:"required"`
	Score    int         `mapstructure:"score" validate:"min=0"`
}

type ServerConfig struct {
	Host           string        `mapstructure:"host" validate:"required,hostname|ip"`
	Port           int           `mapstructure:"port" validate:"required,min=1,max=65535"`
	Mode           string        `mapstructure:"mode" validate:"required,oneof=debug release test"`
	ReadTimeout    time.Duration `mapstructure:"read_timeout" validate:"required,min=1s"`
	WriteTimeout   time.Duration `mapstructure:"write_timeout" validate:"required,min=1s"`
	MaxHeaderBytes int           `mapstructure:"max_header_bytes" validate:"omitempty,min=1024"`
}

type JWTConfig struct {
	Secret        string        `mapstructure:"secret" validate:"required,min=32"`
	Expire        time.Duration `mapstructure:"expire" validate:"required,min=1m"`
	RefreshExpire time.Duration `mapstructure:"refresh_expire" validate:"required,min=1m,gtfield=Expire"`
	Issuer        string        `mapstructure:"issuer" validate:"required"`
}

type LogConfig struct {
	Level      string    `mapstructure:"level" validate:"required,oneof=debug info warn error panic fatal"`
	Filename   string    `mapstructure:"filename" validate:"required"`
	MaxSize    int       `mapstructure:"max_size" validate:"required,min=1"`
	MaxBackups int       `mapstructure:"max_backups" validate:"omitempty,min=0"`
	MaxAge     int       `mapstructure:"max_age" validate:"omitempty,min=0"`
	Compress   bool      `mapstructure:"compress"`
	Console    bool      `mapstructure:"console"`
	JSONFormat bool      `mapstructure:"json_format"`
	DB         *DBConfig `mapstructure:"db" validate:"omitempty"` // nil 时跳过验证
}

type DBConfig struct {
	DSN        string        `mapstructure:"dsn" validate:"required"`
	MaxBuffer  int           `mapstructure:"max_buffer" validate:"omitempty,min=1"`
	BatchSize  int           `mapstructure:"batch_size" validate:"omitempty,min=1"`
	FlushEvery time.Duration `mapstructure:"flush_every" validate:"omitempty,min=100ms"`
	Retention  int           `mapstructure:"retention" validate:"omitempty,min=0"`
}

type EmailConfig struct {
	Host     string `mapstructure:"host" validate:"required,hostname|ip"`
	Port     int    `mapstructure:"port" validate:"required,min=1,max=65535"`
	Username string `mapstructure:"username" validate:"required"`
	Password string `mapstructure:"password" validate:"required"`
	From     string `mapstructure:"from" validate:"required,email"`
	FromName string `mapstructure:"from_name" validate:"required"`
	SSL      bool   `mapstructure:"ssl"`
	TLS      bool   `mapstructure:"tls"`
	PoolSize int    `mapstructure:"pool_size" validate:"omitempty,min=1"`
}

type OAuthConfig struct {
	// Github GithubOAuthConfig `mapstructure:"github" validate:"required"`
	// Google GoogleOAuthConfig `mapstructure:"google" validate:"required"`
	// Wechat WechatOAuthConfig `mapstructure:"wechat" validate:"required"`
}

// type GithubOAuthConfig struct {
// 	ClientID     string `mapstructure:"client_id" validate:"required"`
// 	ClientSecret string `mapstructure:"client_secret" validate:"required"`
// 	RedirectURL  string `mapstructure:"redirect_url" validate:"required,url"`
// }

// type GoogleOAuthConfig struct {
// 	ClientID     string `mapstructure:"client_id" validate:"required"`
// 	ClientSecret string `mapstructure:"client_secret" validate:"required"`
// 	RedirectURL  string `mapstructure:"redirect_url" validate:"required,url"`
// }

// type WechatOAuthConfig struct {
// 	AppID       string `mapstructure:"app_id" validate:"required"`
// 	AppSecret   string `mapstructure:"app_secret" validate:"required"`
// 	RedirectURL string `mapstructure:"redirect_url" validate:"required,url"`
// }
