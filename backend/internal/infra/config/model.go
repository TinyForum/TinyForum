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
	// RateLimit   RateLimitConfig   // 限流配置文件
}

// MARK: 基础配置集合
type ConfigBasic struct {
	Server   ServerConfig   `mapstructure:"server" validate:"required"` // 服务配置
	Frontend FrontendConfig `mapstructure:"frontend" validate:"required"`
	API      APIConfig      `mapstructure:"api" validate:"required"`
	Log      LogConfig      `mapstructure:"log" validate:"required"`
	// RateLimit    RateLimitConfig `mapstructure:"rate_limit" validate:"required"`
	Ollama       Ollama       `mapstructure:"ollama" validate:"required"`
	AllowOrigins []string     `mapstructure:"allow_origins" validate:"omitempty,dive,url"` // 允许跨域请求的域名
	Upload       UploadConfig `mapstructure:"upload" validate:"required"`                  // 上传配置
	Version      string       `mapstructure:"version" validate:"required,semver"`
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
	RateLimit RateLimitConfig `mapstructure:"rate_limit" validate:"required"`
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
