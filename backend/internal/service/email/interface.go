package email

import (
	"fmt"
	"sync"
	"time"

	"tiny-forum/internal/infra/config"
	"tiny-forum/pkg/logger"

	"gopkg.in/gomail.v2"
)

// EmailService 邮件服务接口
type EmailService interface {
	SendResetPasswordEmail(to, token string, tokenExpiresIn time.Duration, username, appURL, apiVersion, ip, userAgent, locale string) error // 重置密码
	SendWelcomeEmail(to, username, locale, appURL string) error                                                                              // 欢迎邮件
	SendVerificationEmail(to, token, username, locale, appURL string) error                                                                  // 验证邮件
	TestConnection(to string) error                                                                                                          // 测试连接
	Close() error                                                                                                                            // 关闭服务
	IsEnabled() bool                                                                                                                         // 是否启用
}

// emailService 邮件服务实现（私有）
type emailService struct {
	dialer      *gomail.Dialer
	from        string
	fromName    string
	pool        chan *gomail.Message
	poolSize    int
	wg          sync.WaitGroup
	stopChan    chan struct{}
	templateDir string
	enabled     bool
}

// NewEmailService 创建邮件服务实例
func NewEmailService(cfg *config.EmailConfig) EmailService {
	if cfg == nil || cfg.Host == "" {
		logger.Warn("Email service disabled: missing configuration")
		return &emailService{enabled: false}
	}

	// 验证必要配置
	if cfg.Username == "" || cfg.Password == "" {
		logger.Warn("Email service disabled: missing username or password")
		return &emailService{enabled: false}
	}

	// 设置默认端口
	port := cfg.Port
	if port == 0 {
		port = 587
	}

	// 创建 dialer
	dialer := gomail.NewDialer(cfg.Host, port, cfg.Username, cfg.Password)

	// 配置 SSL
	if cfg.SSL {
		dialer.SSL = true
	}

	// 设置连接池大小
	poolSize := cfg.PoolSize
	if poolSize <= 0 {
		poolSize = 5
	}

	service := &emailService{
		dialer:      dialer,
		from:        cfg.From,
		fromName:    cfg.FromName,
		pool:        make(chan *gomail.Message, poolSize),
		poolSize:    poolSize,
		stopChan:    make(chan struct{}),
		templateDir: "templates/emails",
		enabled:     true,
	}

	// 启动邮件发送工作池
	service.startWorkerPool()

	logger.Info(fmt.Sprintf("Email service initialized: %s:%d", cfg.Host, port))

	return service
}
