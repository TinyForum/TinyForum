package do

import "time"

// SystemStatus 系统状态
type SystemStatus struct {
	ID           int64     `json:"id"`
	ServiceName  string    `json:"service_name"`  // 服务名称
	Status       string    `json:"status"`        // healthy/unhealthy/warning
	Message      string    `json:"message"`       // 状态描述
	LastCheckAt  time.Time `json:"last_check_at"` // 最后检查时间
	ResponseTime int64     `json:"response_time"` // 响应时间(ms)
}

// SystemHealth 系统健康检查聚合信息
type SystemHealth struct {
	OverallStatus string         `json:"overall_status"` // healthy/unhealthy/warning
	Services      []SystemStatus `json:"services"`
	Metrics       SystemMetrics  `json:"metrics"`
	Issues        []SystemIssue  `json:"issues"`
	LastUpdate    time.Time      `json:"last_update"`
}

// SystemMetrics 系统指标
type SystemMetrics struct {
	CPUUsage    float64 `json:"cpu_usage"`    // CPU使用率
	MemoryUsage float64 `json:"memory_usage"` // 内存使用率
	DiskUsage   float64 `json:"disk_usage"`   // 磁盘使用率
	LoadAvg1    float64 `json:"load_avg_1"`   // 1分钟负载
	LoadAvg5    float64 `json:"load_avg_5"`   // 5分钟负载
	LoadAvg15   float64 `json:"load_avg_15"`  // 15分钟负载
}

// SystemIssue 系统问题
type SystemIssue struct {
	Severity    string    `json:"severity"`    // error/warning/info
	Service     string    `json:"service"`     // 相关服务
	Title       string    `json:"title"`       // 问题标题
	Description string    `json:"description"` // 问题描述
	Suggestion  string    `json:"suggestion"`  // 修复建议
	OccurredAt  time.Time `json:"occurred_at"` // 发生时间
}

// SystemConfig 系统配置（供维护者查看）
type SystemConfig struct {
	ID          int64     `json:"id"`
	Key         string    `json:"key"`         // 配置键
	Value       string    `json:"value"`       // 配置值
	Description string    `json:"description"` // 配置说明
	Category    string    `json:"category"`    // 分类
	UpdatedBy   string    `json:"updated_by"`  // 更新人
	UpdatedAt   time.Time `json:"updated_at"`  // 更新时间
}

// SystemOperationLog 系统操作日志
type SystemOperationLog struct {
	ID        int64     `json:"id"`
	Operator  string    `json:"operator"`   // 操作人
	Operation string    `json:"operation"`  // 操作类型
	Target    string    `json:"target"`     // 操作目标
	Details   string    `json:"details"`    // 详情
	IP        string    `json:"ip"`         // IP地址
	UserAgent string    `json:"user_agent"` // User Agent
	CreatedAt time.Time `json:"created_at"` // 操作时间
}
