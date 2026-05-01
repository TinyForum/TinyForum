package email

import (
	"fmt"
	"time"
)

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
