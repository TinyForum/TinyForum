# 调用方法



> 原则：各层级优先横向调用，例如 adminService 需要操作数据时，先看 xxxSvc 是否有相关的操作，如果没有，就去对应的 Svc 中创建，然后调用
>
> 避免 AxRepo 调用 BxRepos。
>
> 避免 AxHandler 调用 BxHandler

## 基本原则

1. **调用方向必须单向**：严格遵循 `Handler → Service → Repository` 的依赖方向，禁止反向调用（如 Repository 调用 Service）。
2. **同一层内可以相互调用**：Service 可以调用其他 Service，Repository 可以调用其他 Repository，但需谨慎。
3. **禁止跨层调用**：Handler 不能直接调用 Repository；Service 不能直接操作数据库（应通过 Repository）。



## Handler

1. 在 handler 中应该进行 request 的原始绑定，然后直接交由 service 处理。
2. 标准数据类型：
   1. 接收：`Request`
   2. 返回：`VO`

3. **禁止横向调用**：Handler 应保持轻薄，只负责参数解析、验证和响应封装。一个 Handler 不应直接调用另一个 Handler。
4. 如需复用逻辑，应下沉到 Service 层。

```go
func (h *Handler) ListPlugins(c *gin.Context) {
	var req request.ListPluginsRequest
	if err := req.Bind(c); err != nil {
		response.HandleError(c, err)
		return
	}

	// 获取当前登录用户ID（中间件注入，类型为 uint）
	UserID := c.GetUint("user_id")

	// 构建分页查询对象 PageQuery[PluginQueryBO]
	pageQuery := &bo.PageQuery[bo.PluginQueryBO]{
		Page:     req.Page,
		PageSize: req.PageSize,
		SortBy:   req.SortBy,
		Order:    req.Order,
		Options: bo.PluginQueryBO{
			Name:     req.Keyword,
			AuthorID: UserID,     
			Category: req.Category,       
			Tags:     req.Tags,     
			Type:     req.Type,       
			Keyword:  req.Keyword,
			Status:   do.PluginStatus(req.Status), 
			Version:  req.Version,                       
		},
	}

	// 调用 Service（参数类型匹配）
	pageResult, err := h.svc.ListPlugins(c.Request.Context(), pageQuery)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, pageResult)
}

```



## Service

1. Service 中应该进行请求合规验证和数据清洗，避免数据泄漏。
2. Service 只处理请求和响应的数据，不能直接读写数据库。
3. 标准数据类型：
   1. 接收：`BO`
   2. 返回：`VO`
   3. 该层涉及复杂的 `DTO` 模型，
4. **允许横向调用**：一个 Service 可以调用另一个 Service 的方法。
5. **规范要求**：
   - 通过接口依赖，而非直接依赖实现类（便于单元测试和替换）。
   - 避免循环依赖（A Service 调用 B Service，B 又调用 A）。
   - 横向调用时应传递 context（如请求 ID、用户信息），保持链路追踪。
   - 被调用的 Service 方法应职责单一、无副作用的业务逻辑。
   - 尽量避免过深的调用链（不超过 3 层）。

接口定义

```go
ListPlugins(ctx context.Context, queryBO *bo.PageQuery[bo.PluginQueryBO]) (*common.PageResult[vo.PluginMetaVO], error)
```

实现

```go
func (s *pluginService) ListPlugins(ctx context.Context, queryBO *bo.PageQuery[bo.PluginQueryBO]) (*common.PageResult[vo.PluginMetaVO], error) {
	// 1. 防御性检查
	if queryBO == nil {
		return nil, apperrors.ErrValidation
	}

	// 2. 业务校验：仅当 Status 非空且无效时报错
	status := queryBO.Options.Status
	if status != "" && !status.IsValid() {
		logger.Warnf("无效的插件状态: %s", status)
		return nil, apperrors.ErrValidation
	}

	// 3. 规范化分页参数
	page, pageSize := queryBO.Page, queryBO.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	// 4. 构建 Repository 查询参数（内联转换，简洁清晰）
	repoQuery := do.PluginMeta{
		Name:     queryBO.Options.Name,
		AuthorID: queryBO.Options.AuthorID,
		Category: queryBO.Options.Category,
		Tags:     queryBO.Options.Tags,
		Type:     do.PluginType(queryBO.Options.Type),
		Status:   queryBO.Options.Status,
		Version:  queryBO.Options.Version,
	}
	queryDO := &common.PageQuery[do.PluginMeta]{
		Page:     page,
		PageSize: pageSize,
		SortBy:   "created_at", 
		Order:    "desc",
		Data:     repoQuery,
		Keyword:  queryBO.Keywords,
	}

	// 5. 调用 Repository 层
	plugins, total, err := s.repo.List(ctx, queryDO)
	if err != nil {
		logger.Errorf("查询插件列表失败: %v, query: %+v", err, repoQuery)
		return nil, apperrors.ErrInternalError
	}

	// 6. 批量转换 DO -> VO
	vos := make([]vo.PluginMetaVO, 0, len(plugins))
	for _, p := range plugins {
		vos = append(vos, vo.PluginMetaVO{
			ID:            p.ID,
			CreatedAt:     p.CreatedAt,
			UpdatedAt:     p.UpdatedAt,
			Name:          p.Name,
			Version:       p.Version,
			Description:   p.Description,
			Summary:       p.Summary,
			IconURL:       p.IconURL,
			Screenshots:   p.Screenshots,
			HomepageURL:   p.HomepageURL,
			Type:          string(p.Type),
			Category:      string(p.Category),
			Tags:          p.Tags, // []string
			AuthorID:      p.AuthorID,
			AuthorURL:     p.AuthorURL,
			ScriptURL:     p.ScriptURL,
			ServerEntry:   p.ServerEntry,
			Slots:         p.Slots,
			Routes:        p.Routes,
			Pricing:       p.Pricing,
			Compatibility: p.Compatibility,
			Permissions:   p.Permissions,
			Enabled:       p.Enabled,
			Status:        string(p.Status),
			InstallCount:  p.InstallCount,
			Rating:        p.Rating,
			ConfigSchema:  p.ConfigSchema,
		})
	}

	return &common.PageResult[vo.PluginMetaVO]{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		List:     vos,
	}, nil
}

```

## Repository

1. Response 中绝对信任 Service，不再进行数据检查。
2. 标准数据类型：
   1. 接收：`DO`
   2. 返回：`DO`
3. **允许横向调用**：一个 Repository 可以调用另一个 Repository（例如跨表操作需要聚合数据）。
4. **规范要求**：
   - 仅限于数据访问层内部的组合查询，不应包含业务逻辑。
   - 避免循环依赖和过长调用链。
   - 可使用事务管理器协调多个 Repository 调用，确保原子性。

接口定义

```go
	List(ctx context.Context, queryDO *common.PageQuery[do.PluginMeta]) ([]*do.PluginMeta, int64, error)
```

实现

```go
func (r *pluginRepo) List(ctx context.Context, queryBO *common.PageQuery[do.PluginMeta]) ([]*do.PluginMeta, int64, error) {
	db := r.db.WithContext(ctx).Model(&do.PluginMeta{})

	// 软删除过滤
	db = db.Where("deleted_at IS NULL")

	// 用户过滤
	if queryBO.Data.AuthorID > 0 {
		db = db.Where("user_id = ?", queryBO.Data.AuthorID)
	}

	// 状态过滤
	if queryBO.Data.Status != "" {
		db = db.Where("status = ?", queryBO.Data.Status)
	}

	// 标签过滤
	if len(queryBO.Data.Tags) != 0 {
		db = db.Where("tag = ?", queryBO.Data.Tags)
	}

	// 关键字模糊搜索
	if queryBO.Keyword != "" {
		keyword := "%" + queryBO.Keyword + "%"
		db = db.Where("name LIKE ? OR description LIKE ?", keyword, keyword)
	}

	// 排序
	sortField := queryBO.SortBy
	if sortField == "" {
		sortField = "created_at"
	}
	// 默认降序
	order := queryBO.Order
	if order == "" {
		order = "desc"
	}
	// 支持自定义排序字段
	db = db.Order(fmt.Sprintf("%s %s", sortField, order))

	// 分页查询
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err

	}

	var plugins []*do.PluginMeta
	offset := (queryBO.Page - 1) * queryBO.PageSize
	err := db.Offset(offset).Limit(queryBO.PageSize).Find(&plugins).Error
	return plugins, total, err
}

```

