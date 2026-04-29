package email

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"
	"time"

	"tiny-forum/pkg/logger"

	"gopkg.in/gomail.v2"
)

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
	logger.Debugf("attempting to send reset email",
		"to", to,
		"enabled", s.enabled,
		// "has_token", token != ""
	)

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
func (s *emailService) SendResetPasswordEmail(to, token, username, appURL, apiVersion, ip, userAgent, locale string) error {
	if !s.enabled {
		logger.Warn("Email service disabled, skipping password reset email")
		return nil
	}

	// 构建重置链接
	// resetURL := fmt.Sprintf("%s/%s/auth/password/validate-token?token=%s", appURL, apiVersion, token)
// http://localhost:3000/zh-CN/auth/reset-password?token=ec2d50e6a812530e3baca78a8b1e0ba405fa4b8932185fd1235a206569cb44a4
resetURL := fmt.Sprintf("%s/auth/reset-password?token=%s", appURL, token)

	// 准备邮件数据
	data := EmailData{
		Username:     username,
		ResetURL:     resetURL,
		ExpiresIn:    "60 minutes",
		Year:         time.Now().Year(),
		AppName:      "TinyForum",
		SupportEmail: "support@tinyforum.com",
		SiteURL:      appURL,
		RequestTime:  time.Now().Format("2006-01-02 15:04:05"),
		RequestIP:    ip,
		UserAgent:    userAgent,
		Location:     locale,
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

	path := filepath.Join(s.templateDir, templateName)
	var tmpl *template.Template
	tmpl, err := template.ParseFiles(path)

	if err == nil {
		logger.Debug(fmt.Sprintf("Loaded email template from: %s", path))
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
	// 确保 SupportEmail 有默认值
	supportEmail := data.SupportEmail
	if supportEmail == "" {
		supportEmail = "support@example.com"
	}

	// 格式化请求时间
	requestTime := data.RequestTime
	if requestTime == "" {
		requestTime = time.Now().Format("2006-01-02 15:04:05")
	}

	// 格式化 IP 和 UserAgent
	requestIP := data.RequestIP
	if requestIP == "" {
		requestIP = "未知"
	}

	userAgent := data.UserAgent
	if userAgent == "" {
		userAgent = "未知"
	}

	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>重置您的密码</title>
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
            background-color: #ff8c8c;
        }
        .info-box {
            background-color: #f0f9ff;
            border-left: 4px solid #ff7954;
            padding: 12px;
            margin: 16px 0;
            font-size: 13px;
            border-radius: 4px;
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
            <div class="logo">🔐 %s</div>
        </div>
        <div class="content">
            <h2>重置您的密码</h2>
            
            <p>您好，<strong>%s</strong>：</p>
            
            <p>我们收到了重置您 <strong>%s</strong> 账户密码的请求。</p>
            
            <div class="info-box">
                <strong>📋 请求详情：</strong><br>
                • 请求时间：%s<br>
                • 来源 IP：<code>%s</code><br>
                • 设备信息：%s
            </div>
            
            <p>点击下面的按钮创建新密码：</p>
            
            <div style="text-align: center;">
                <a href="%s" class="button">重置密码</a>
            </div>
            
            <p>或者复制以下链接到浏览器：</p>
            <p class="link">%s</p>
            
            <div class="warning">
                <strong>⚠️ 此链接将在 %s 后失效</strong><br>
                如果您没有请求重置密码，请忽略此邮件，您的密码将保持不变。
            </div>
            
            <p>为了账户安全，请不要将密码告诉任何人。</p>
        </div>
        <div class="footer">
            <p>&copy; %d %s. 保留所有权利。</p>
            <p>需要帮助？请联系 <a href="mailto:%s" style="color: #3b82f6;">%s</a></p>
        </div>
    </div>
</body>
</html>
`,
		data.AppName,   // logo
		data.Username,  // 用户名
		data.AppName,   // 应用名称
		requestTime,    // 请求时间
		requestIP,      // 来源 IP
		userAgent,      // 设备信息
		data.ResetURL,  // 按钮链接
		data.ResetURL,  // 文本链接
		data.ExpiresIn, // 过期时间
		data.Year,      // 年份
		data.AppName,   // 版权名称
		supportEmail,   // 支持邮箱
		supportEmail,   // 支持邮箱（显示）
	)
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
