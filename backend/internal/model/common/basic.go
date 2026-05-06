package common

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"` // 插入时自动填充
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index" swaggerignore:"true"`
}

type BasicResponse struct {
	Code      int    `json:"code"`                 // 业务状态码
	Message   string `json:"message"`              // 业务信息
	Data      any    `json:"data,omitempty"`       // 业务数据
	Timestamp int64  `json:"timestamp"`            // 时间戳
	RequestID string `json:"request_id,omitempty"` // 请求ID
	TraceID   string `json:"trace_id,omitempty"`   // 追踪ID
}
