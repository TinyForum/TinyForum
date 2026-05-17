package reports

import (
	"context"
	"tiny-forum/internal/model/bo"
	"tiny-forum/internal/model/common"
	"tiny-forum/internal/model/do"
)

func (r *reportsRepository) List(ctx context.Context, listReportBO *common.PageQuery[bo.ListReportBO]) ([]do.Report, int64, error) {
	var reports []do.Report
	var total int64

	// 1. 基础查询
	query := r.db.Model(&do.Report{})

	// 2. 动态添加过滤条件（基于 listReportBO.Data）
	if listReportBO != nil {
		data := listReportBO.Data

		if data.ID != 0 {
			query = query.Where("id = ?", data.ID)
		}
		if data.ReporterID != 0 {
			query = query.Where("reporter_id = ?", data.ReporterID)
		}
		if data.TargetID != 0 {
			query = query.Where("target_id = ?", data.TargetID)
		}
		if data.TargetType != "" {
			query = query.Where("target_type = ?", data.TargetType)
		}
		if data.Type != "" {
			query = query.Where("type = ?", data.Type)
		}
		if data.Reason != "" {
			query = query.Where("reason LIKE ?", "%"+data.Reason+"%")
		}
		if data.Status != "" {
			query = query.Where("status = ?", data.Status)
		}
		if data.HandlerID != nil {
			// 注意：如果传了具体的 uint 指针，筛选等于该值；若需要查 NULL 可设计额外字段
			query = query.Where("handler_id = ?", *data.HandlerID)
		}
		if data.HandleNote != "" {
			query = query.Where("handle_note LIKE ?", "%"+data.HandleNote+"%")
		}
		if data.HandleAt != nil {
			query = query.Where("handle_at = ?", data.HandleAt)
		}
		if data.ContentSnapshot != "" {
			query = query.Where("content_snapshot LIKE ?", "%"+data.ContentSnapshot+"%")
		}
		if data.ReporterIP != "" {
			query = query.Where("reporter_ip = ?", data.ReporterIP)
		}
		if data.IsAnonymous {
			query = query.Where("is_anonymous = ?", true)
		}
		if data.Priority != 0 {
			query = query.Where("priority = ?", data.Priority)
		}
	}

	// 3. 关键字搜索（如果提供，可搜索 Reasons、ContentSnapshot 等）
	if listReportBO != nil && listReportBO.Keyword != "" {
		keyword := "%" + listReportBO.Keyword + "%"
		query = query.Where("reason LIKE ? OR content_snapshot LIKE ? OR reporter_ip LIKE ?", keyword, keyword, keyword)
	}

	// 4. 获取总记录数（在分页和 Preload 之前）
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 5. 预加载关联
	query = query.Preload("Reporter").Preload("Handler")

	// 6. 排序（使用白名单防止 SQL 注入）
	sortBy := "created_at" // 默认排序字段
	order := "DESC"        // 默认排序方向

	if listReportBO != nil && listReportBO.SortBy != "" {
		// 允许排序的字段白名单
		allowedSortFields := map[string]bool{
			"id": true, "created_at": true, "updated_at": true, "status": true,
			"priority": true, "handle_at": true, "reporter_id": true,
		}
		if allowedSortFields[listReportBO.SortBy] {
			sortBy = listReportBO.SortBy
		}
	}
	if listReportBO != nil && listReportBO.Order != "" {
		if listReportBO.Order == "ASC" || listReportBO.Order == "asc" {
			order = "ASC"
		} else if listReportBO.Order == "DESC" || listReportBO.Order == "desc" {
			order = "DESC"
		}
	}
	query = query.Order(sortBy + " " + order)

	// 7. 分页参数
	page := 1
	pageSize := 10 // 默认每页10条
	if listReportBO != nil {
		if listReportBO.Page > 0 {
			page = listReportBO.Page
		}
		if listReportBO.PageSize > 0 {
			pageSize = listReportBO.PageSize
		}
	}
	offset := (page - 1) * pageSize

	// 8. 执行查询
	err := query.Limit(pageSize).Offset(offset).Find(&reports).Error
	if err != nil {
		return nil, 0, err
	}

	return reports, total, nil
}
