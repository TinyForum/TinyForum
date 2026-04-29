package email

import (
	"fmt"
	"sync"
	"time"

	"tiny-forum/config"
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

// EmailData 邮件数据
type EmailData struct {
	Username     string
	ResetURL     string
	ExpiresIn    string
	Year         int
	AppName      string
	SupportEmail string
	SiteURL      string
	RequestTime  string // 请求时间
	RequestIP    string // 请求 IP
	UserAgent    string // 用户代理
	Location     string // IP 地理位置（可选）

}

// FormatDuration 将一个时间持续时间格式化为易读的字符串表示
// 参数:
//
//	d - time.Duration类型的时间持续时间
//
// 返回值:
//
//	string - 格式化后的时间字符串，可能包含天、小时和分钟
func FormatDuration(d time.Duration) string {
	// 计算天数，将总小时数除以24并取整
	days := int(d.Hours() / 24)
	// 计算剩余的小时数，取总小时数除以24的余数
	hours := int(d.Hours()) % 24
	// 计算剩余的分钟数，取总分钟数除以60的余数
	minutes := int(d.Minutes()) % 60

	// 如果天数大于0
	if days > 0 {
		// 如果分钟数大于0，返回包含天、小时和分钟的格式化字符串
		if minutes > 0 {
			return fmt.Sprintf("%d天%d小时%d分钟", days, hours, minutes)
		}
		// 如果分钟数为0，返回只包含天和小时的格式化字符串
		return fmt.Sprintf("%d天%d小时", days, hours)
	}

	// 如果小时数大于0（天数等于0）
	if hours > 0 {
		// 如果分钟数大于0，返回包含小时和分钟的格式化字符串
		if minutes > 0 {
			return fmt.Sprintf("%d小时%d分钟", hours, minutes)
		}
		// 如果分钟数为0，返回只包含小时的格式化字符串
		return fmt.Sprintf("%d小时", hours)
	}

	// 如果天和小时都为0，只返回分钟数的格式化字符串
	return fmt.Sprintf("%d分钟", minutes)
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
