# 数据层级标准

> 原则：各层级优先横向调用，例如 adminService 需要操作数据时，先看 xxxSvc 是否有相关的操作，如果没有，就去对应的 Svc 中创建，然后调用，避免直接调用非自己业务的 Repo。

## Handler

1. 在 handler 中应该进行 request 的原始绑定，然后直接交由 service 处理。
2. 优先横向调用

```go
func (h *AnnouncementHandler) Create(c *gin.Context) {
	var req request.CreateAnnouncement
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	userID := c.GetUint("user_id")
	announcement, err := h.service.Create(c.Request.Context(), &req, userID)
	if err != nil {
		response.InternalError(c, err.Error()) 
		return
	}
	response.Success(c, announcement)
}
```



## Service

1. Service 中应该进行请求合规验证和数据清洗，避免数据泄漏。
2. Service 只处理请求和响应的数据，不能直接读写数据库。
3. Service 横向优先

```go

func (s *announcementService) Create(ctx context.Context, req *request.CreateAnnouncement, userID uint) (*do.Announcement, error) {
	if err := s.validateTime(req.PublishedAt, req.ExpiredAt); err != nil {
		return nil, err
	}
	now := time.Now()
	announcement := &do.Announcement{
		Title:       req.Title,
		Content:     req.Content,
		Summary:     req.Summary,
		Cover:       req.Cover,
		Type:        req.Type,
		IsPinned:    req.IsPinned,
		IsGlobal:    req.IsGlobal,
		BoardID:     req.BoardID,
		PublishedAt: req.PublishedAt,
		ExpiredAt:   req.ExpiredAt,
		Status:      req.Status,
		ViewCount:   0,
		CreatedBy:   userID,
		UpdatedBy:   userID,
	}
	if announcement.PublishedAt == nil {
		announcement.PublishedAt = &now
	}
	if err := s.repo.Create(ctx, announcement); err != nil {
		return nil, err
	}
	return announcement, nil
}
```

## Repository

1. Response 中绝对信任 Service，不再进行数据检查。

```go

func (r *pluginRepository) List(ctx context.Context, query *dto.PluginQuery, pageParam common.PageParam) (*common.PageResult[do.PluginMeta], error){
db := r.db.WithContext(ctx).Model(&do.PluginMeta{})

    // 动态条件
    if query.AuthorID != 0 {
        db = db.Where("author_id = ?", query.AuthorID)
    }
    if query.TagID != 0 {
        db = db.Where("tag_id = ?", query.TagID)
    }
    if query.PostType != "" {
        db = db.Where("post_type = ?", query.PostType)
    }
    if query.Keyword != "" {
        db = db.Where("name LIKE ?", "%"+query.Keyword+"%")
    }
    if query.Status != 0 {
        db = db.Where("status = ?", query.Status)
    }

    // 总数
    var total int64
    if err := db.Count(&total).Error; err != nil {
        return nil, err
    }

    // 列表
    var list []do.PluginMeta
    offset := (pageParam.Page - 1) * pageParam.PageSize
    order := "created_at DESC"
    if query.SortBy != "" {
        order = query.SortBy
    }
    err := db.Offset(offset).Limit(pageParam.PageSize).Order(order).Find(&list).Error
    if err != nil {
        return nil, err
    }

    hasMore := int64(pageParam.Page*pageParam.PageSize) < total

    return &common.PageResult[do.PluginMeta]{
        Total:    total,
        Page:     pageParam.Page,
        PageSize: pageParam.PageSize,
        List:     list,
        HasMore:  hasMore,
    }, nil
}

```

