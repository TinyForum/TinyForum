package config

import (
	"time"
	"tiny-forum/pkg/logger"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Log      LogConfig      `mapstructure:"log"`
}

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
	TimeZone string `mapstructure:"timezone"`
}

type JWTConfig struct {
	Secret string        `mapstructure:"secret"`
	Expire time.Duration `mapstructure:"expire"`
}

type LogConfig struct {
	Level      string `mapstructure:"level" json:"level"`
	Filename   string `mapstructure:"filename" json:"filename"`
	MaxSize    int    `mapstructure:"max_size" json:"max_size"`
	MaxBackups int    `mapstructure:"max_backups" json:"max_backups"`
	MaxAge     int    `mapstructure:"max_age" json:"max_age"`
	Compress   bool   `mapstructure:"compress" json:"compress"`
	Console    bool   `mapstructure:"console" json:"console"`
	JSONFormat bool   `mapstructure:"json_format" json:"json_format"`
}

func Load(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c *Config) ToLoggerConfig() logger.Config {
	return logger.Config{
		Level:      c.Log.Level,
		Filename:   c.Log.Filename,
		MaxSize:    c.Log.MaxSize,
		MaxBackups: c.Log.MaxBackups,
		MaxAge:     c.Log.MaxAge,
		Compress:   c.Log.Compress,
		Console:    c.Log.Console,
		JSONFormat: c.Log.JSONFormat,
	}
}
