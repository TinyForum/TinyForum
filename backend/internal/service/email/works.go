package email

import (
	"fmt"
	"time"

	"tiny-forum/internal/model/vo"
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
func (s *emailService) SendResetPasswordEmail(
	to string,
	token string,
	tokenExpiresIn time.Duration,
	username, appURL, apiVersion, ip, userAgent, locale string,
) error {
	if !s.enabled {
		logger.Warn("Email service disabled, skipping password reset email")
		return nil
	}

	logger.Infof("user email:", to,
		"token:", token,
		"tokenExpiresIn:", tokenExpiresIn,
		"username:", username,
		"appURL:", appURL,
		"apiVersion: %s\n", apiVersion,
		"ip:", ip,
		"userAgent:", userAgent,
		"locale:", locale)
	// 构建重置链接
	// resetURL := fmt.Sprintf("%s/%s/auth/password/validate-token?token=%s", appURL, apiVersion, token)
	// http://localhost:3000/zh-CN/auth/reset-password?token=ec2d50e6a812530e3baca78a8b1e0ba405fa4b8932185fd1235a206569cb44a4
	resetURL := fmt.Sprintf("%s/auth/reset-password?token=%s", appURL, token)

	logger.Infof(resetURL)

	// 准备邮件数据
	data := vo.ResetPasswordEmailDataHtml{
		Username:     username,
		ResetURL:     resetURL,
		ExpiresIn:    FormatDuration(tokenExpiresIn),
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
	body, err := s.renderResetPasswordTemplate("reset_password.html", data)
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
