package request

import "github.com/gin-gonic/gin"

// PageRequest 纯分页参数
type PageRequest struct {
	Page     int    `form:"page" binding:"omitempty,min=1" json:"page"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100" json:"page_size"`
	SortBy   string `form:"sort_by" json:"sort_by"`
	Order    string `form:"order" binding:"omitempty,oneof=asc desc" json:"order"`
}

// Bind 从 Gin Context 绑定并规范化分页参数
func (p *PageRequest) Bind(c *gin.Context) error {
	// 1. 绑定查询参数（自动校验 binding 标签）
	if err := c.ShouldBindQuery(p); err != nil {
		return err
	}
	// 2. 补全默认值
	p.normalize()
	return nil
}

// normalize 内部方法，设置默认值
func (p *PageRequest) normalize() {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 20
	}
	if p.PageSize > 100 {
		p.PageSize = 100
	}
	if p.Order == "" {
		p.Order = "desc"
	}
}

// Offset / Limit 保持不变
func (p *PageRequest) Offset() int { return (p.Page - 1) * p.PageSize }
func (p *PageRequest) Limit() int  { return p.PageSize }
