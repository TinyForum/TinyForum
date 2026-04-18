package email

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"
	"sync"
	"time"

	"tiny-forum/config"
	"tiny-forum/pkg/logger"

	"gopkg.in/gomail.v2"
)

// EmailService 邮件服务接口
type EmailService interface {
	SendResetPasswordEmail(to, token, username, locale, appURL string) error
	SendWelcomeEmail(to, username, locale, appURL string) error
	SendVerificationEmail(to, token, username, locale, appURL string) error
	TestConnection(to string) error
	Close() error
	IsEnabled() bool
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

// startWorkerPool 启动工作池
func (s *emailService) startWorkerPool() {
	if !s.enabled {
		return
	}

	for i := 0; i < s.poolSize; i++ {
		s.wg.Add(1)
		go s.worker()
	}
}

// worker 邮件发送工作协程
func (s *emailService) worker() {
	defer s.wg.Done()

	for {
		select {
		case msg, ok := <-s.pool:
			if !ok {
				return
			}
			if err := s.dialer.DialAndSend(msg); err != nil {
				logger.Error(fmt.Sprintf("Failed to send email: %v", err))
			}
		case <-s.stopChan:
			return
		}
	}
}

// SendEmail 发送邮件（异步）
func (s *emailService) sendEmail(to, subject, body string) error {
	if !s.enabled {
		logger.Debug("Email service disabled, skipping send")
		return nil
	}

	if s.dialer == nil {
		return fmt.Errorf("email service not initialized")
	}

	msg := gomail.NewMessage()

	// 设置发件人
	from := s.from
	if s.fromName != "" {
		from = fmt.Sprintf("%s <%s>", s.fromName, s.from)
	}
	msg.SetHeader("From", from)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)

	// 异步发送
	select {
	case s.pool <- msg:
		return nil
	default:
		// 如果队列满了，同步发送
		logger.Warn("Email queue full, sending synchronously")
		return s.dialer.DialAndSend(msg)
	}
}

// SendEmailSync 同步发送邮件
func (s *emailService) sendEmailSync(to, subject, body string) error {
	if !s.enabled {
		return nil
	}

	if s.dialer == nil {
		return fmt.Errorf("email service not initialized")
	}

	msg := gomail.NewMessage()
	from := s.from
	if s.fromName != "" {
		from = fmt.Sprintf("%s <%s>", s.fromName, s.from)
	}
	msg.SetHeader("From", from)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)

	return s.dialer.DialAndSend(msg)
}

// SendResetPasswordEmail 发送重置密码邮件
func (s *emailService) SendResetPasswordEmail(to, token, username, locale, appURL string) error {
	if !s.enabled {
		logger.Warn("Email service disabled, skipping password reset email")
		return nil
	}

	// 构建重置链接
	resetURL := fmt.Sprintf("%s/%s/auth/reset-password?token=%s", appURL, locale, token)

	// 准备邮件数据
	data := EmailData{
		Username:     username,
		ResetURL:     resetURL,
		ExpiresIn:    "60 minutes",
		Year:         time.Now().Year(),
		AppName:      "TinyForum",
		SupportEmail: "support@tinyforum.com",
		SiteURL:      appURL,
	}

	// 尝试加载模板
	body, err := s.renderTemplate("reset_password.html", data)
	if err != nil {
		logger.Warn(fmt.Sprintf("Failed to load email template, using simple email: %v", err))
		body = s.buildSimpleResetEmail(data)
	}

	subject := "Reset Your Password - TinyForum"

	return s.sendEmail(to, subject, body)
}

// SendWelcomeEmail 发送欢迎邮件
func (s *emailService) SendWelcomeEmail(to, username, locale, appURL string) error {
	if !s.enabled {
		return nil
	}

	data := EmailData{
		Username:     username,
		Year:         time.Now().Year(),
		AppName:      "TinyForum",
		SupportEmail: "support@tinyforum.com",
		SiteURL:      appURL,
	}

	body, err := s.renderTemplate("welcome.html", data)
	if err != nil {
		body = s.buildSimpleWelcomeEmail(username)
	}

	subject := fmt.Sprintf("Welcome to %s!", data.AppName)

	return s.sendEmail(to, subject, body)
}

// SendVerificationEmail 发送邮箱验证邮件
func (s *emailService) SendVerificationEmail(to, token, username, locale, appURL string) error {
	if !s.enabled {
		return nil
	}

	verifyURL := fmt.Sprintf("%s/%s/auth/verify-email?token=%s", appURL, locale, token)

	data := EmailData{
		Username:     username,
		ResetURL:     verifyURL,
		Year:         time.Now().Year(),
		AppName:      "TinyForum",
		SupportEmail: "support@tinyforum.com",
		SiteURL:      appURL,
	}

	body, err := s.renderTemplate("verify_email.html", data)
	if err != nil {
		body = s.buildSimpleVerificationEmail(username, verifyURL)
	}

	subject := "Verify Your Email Address - TinyForum"

	return s.sendEmail(to, subject, body)
}

// renderTemplate 渲染邮件模板
func (s *emailService) renderTemplate(templateName string, data EmailData) (string, error) {
	// 尝试多个模板路径
	templatePaths := []string{
		filepath.Join(s.templateDir, templateName),
		filepath.Join("templates", "emails", templateName),
		filepath.Join("..", "templates", "emails", templateName),
		filepath.Join("../..", "templates", "emails", templateName),
	}

	var tmpl *template.Template
	var err error

	for _, path := range templatePaths {
		tmpl, err = template.ParseFiles(path)
		if err == nil {
			logger.Debug(fmt.Sprintf("Loaded email template from: %s", path))
			break
		}
	}

	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// buildSimpleResetEmail 构建简单的重置密码邮件
func (s *emailService) buildSimpleResetEmail(data EmailData) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Reset Your Password</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
            line-height: 1.6;
            color: #333;
            margin: 0;
            padding: 0;
            background-color: #f4f4f5;
        }
        .container {
            max-width: 560px;
            margin: 40px auto;
            padding: 20px;
            background-color: #ffffff;
            border-radius: 12px;
            box-shadow: 0 2px 8px rgba(0,0,0,0.1);
        }
        .header {
            text-align: center;
            padding-bottom: 20px;
            border-bottom: 2px solid #e4e4e7;
        }
        .logo {
            font-size: 28px;
            font-weight: bold;
            color: #3b82f6;
        }
        .content {
            padding: 30px 20px;
        }
        .button {
            display: inline-block;
            padding: 12px 28px;
            background-color: #3b82f6;
            color: #ffffff;
            text-decoration: none;
            border-radius: 8px;
            margin: 20px 0;
            font-weight: 500;
        }
        .button:hover {
            background-color: #2563eb;
        }
        .link {
            word-break: break-all;
            color: #3b82f6;
            font-size: 14px;
        }
        .footer {
            margin-top: 30px;
            padding-top: 20px;
            border-top: 1px solid #e4e4e7;
            font-size: 12px;
            color: #71717a;
            text-align: center;
        }
        .warning {
            background-color: #fef3c7;
            padding: 12px;
            border-radius: 8px;
            font-size: 13px;
            margin: 16px 0;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <div class="logo">%s</div>
        </div>
        <div class="content">
            <h2>Reset Your Password</h2>
            <p>Hello <strong>%s</strong>,</p>
            <p>We received a request to reset your password. Click the button below to create a new password:</p>
            <div style="text-align: center;">
                <a href="%s" class="button">Reset Password</a>
            </div>
            <p>Or copy and paste this link into your browser:</p>
            <p class="link">%s</p>
            <div class="warning">
                <strong>⚠️ This link will expire in %s</strong>
            </div>
            <p>If you didn't request this, please ignore this email. Your password will remain unchanged.</p>
        </div>
        <div class="footer">
            <p>&copy; %d %s. All rights reserved.</p>
            <p>Need help? Contact us at <a href="mailto:%s">%s</a></p>
        </div>
    </div>
</body>
</html>
`, data.AppName, data.Username, data.ResetURL, data.ResetURL, data.ExpiresIn, data.Year, data.AppName, data.SupportEmail, data.SupportEmail)
}

// buildSimpleWelcomeEmail 构建简单的欢迎邮件
func (s *emailService) buildSimpleWelcomeEmail(username string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Welcome to TinyForum</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 560px; margin: 40px auto; padding: 20px; }
        .header { text-align: center; padding-bottom: 20px; border-bottom: 2px solid #e4e4e7; }
        .logo { font-size: 28px; font-weight: bold; color: #3b82f6; }
        .content { padding: 30px 20px; }
        .footer { margin-top: 30px; padding-top: 20px; border-top: 1px solid #e4e4e7; font-size: 12px; color: #71717a; text-align: center; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <div class="logo">TinyForum</div>
        </div>
        <div class="content">
            <h2>Welcome to TinyForum!</h2>
            <p>Hello <strong>%s</strong>,</p>
            <p>Thank you for joining TinyForum! We're excited to have you on board.</p>
            <p>Get started by:</p>
            <ul>
                <li>Completing your profile</li>
                <li>Exploring discussions</li>
                <li>Creating your first post</li>
            </ul>
            <p>If you have any questions, feel free to reach out to our support team.</p>
            <br>
            <p>Best regards,<br>TinyForum Team</p>
        </div>
        <div class="footer">
            <p>&copy; 2024 TinyForum. All rights reserved.</p>
        </div>
    </div>
</body>
</html>
`, username)
}

// buildSimpleVerificationEmail 构建简单的验证邮件
func (s *emailService) buildSimpleVerificationEmail(username, verifyURL string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Verify Your Email</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 560px; margin: 40px auto; padding: 20px; }
        .button { display: inline-block; padding: 10px 20px; background-color: #3b82f6; color: white; text-decoration: none; border-radius: 5px; }
    </style>
</head>
<body>
    <div class="container">
        <h2>Verify Your Email Address</h2>
        <p>Hello <strong>%s</strong>,</p>
        <p>Please click the button below to verify your email address:</p>
        <p style="text-align: center;">
            <a href="%s" class="button">Verify Email</a>
        </p>
        <p>Or copy and paste this link into your browser:</p>
        <p>%s</p>
        <p>This link will expire in 24 hours.</p>
        <p>If you didn't create an account, please ignore this email.</p>
        <br>
        <p>Best regards,<br>TinyForum Team</p>
    </div>
</body>
</html>
`, username, verifyURL, verifyURL)
}

// TestConnection 测试邮件服务连接
func (s *emailService) TestConnection(to string) error {
	if !s.enabled {
		return fmt.Errorf("email service not enabled")
	}

	if s.dialer == nil {
		return fmt.Errorf("email service not initialized")
	}

	testBody := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <title>Test Email</title>
</head>
<body>
    <h2>Test Email</h2>
    <p>This is a test email from TinyForum. If you're receiving this, your email configuration is working correctly!</p>
    <p>Time: %s</p>
    <p>Best regards,<br>TinyForum Team</p>
</body>
</html>
`, time.Now().Format("2006-01-02 15:04:05"))

	return s.sendEmailSync(to, "Test Email from TinyForum", testBody)
}

// Close 关闭邮件服务
func (s *emailService) Close() error {
	if !s.enabled {
		return nil
	}

	if s.stopChan != nil {
		close(s.stopChan)
	}
	s.wg.Wait()

	if s.pool != nil {
		close(s.pool)
	}

	return nil
}

// IsEnabled 检查邮件服务是否启用
func (s *emailService) IsEnabled() bool {
	return s.enabled
}
