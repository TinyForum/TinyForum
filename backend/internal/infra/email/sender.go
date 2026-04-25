// pkg/email/email.go
package email

import (
	"crypto/tls"
	"fmt"
	"gopkg.in/gomail.v2"
)

// Config 邮件配置
type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	FromName string
	SSL      bool
	TLS      bool
}

// Client 邮件客户端
type Client struct {
	dialer *gomail.Dialer
	from   string
	name   string
}

// NewClient 创建邮件客户端
func NewClient(cfg *Config) *Client {
	dialer := gomail.NewDialer(cfg.Host, cfg.Port, cfg.Username, cfg.Password)

	if cfg.TLS {
		dialer.TLSConfig = &tls.Config{
			InsecureSkipVerify: false,
			ServerName:         cfg.Host,
		}
	}

	if cfg.SSL {
		dialer.SSL = true
	}

	return &Client{
		dialer: dialer,
		from:   cfg.From,
		name:   cfg.FromName,
	}
}

// Send 发送 html 邮件
func (c *Client) Send(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("%s <%s>", c.name, c.from))
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	return c.dialer.DialAndSend(m)
}

// SendSimple 发送简单文本邮件
func (c *Client) SendSimple(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("%s <%s>", c.name, c.from))
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	return c.dialer.DialAndSend(m)
}
