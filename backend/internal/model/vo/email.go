package vo

type ResetPasswordEmailDataHtml struct {
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
