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
func (h *PluginHandler) List(c *gin.Context) {
	var req request.PluginListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Infof("绑定错误: ", err)
		response.BadRequest(c, apperrors.ErrInvalidRequest.Error())
		return
	}

	common.ApplyDefaults(&req)

	// Request -> BO
	queryBO := converter.PluginListRequestToBO(&req)

	// 调用 Service
	pageBO, err := h.service.ListPlugins(c.Request.Context(), queryBO)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	// PageResult[BO] -> PageResult[VO]
	pageVO := converter.PageBOToPageVO(pageBO, converter.PluginBOToVO)

	response.SuccessPage(c, pageVO.List, pageVO.Total, pageVO.Page, pageVO.PageSize)

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

```go
func (s *pluginService) ListPlugins(ctx context.Context, queryBO *bo.PluginQueryBO) (*common.PageResult[bo.PluginMeta], error) {
	// BO -> Query DO
	queryDO := converter.PluginQueryBOToQueryDO(queryBO)

	// 调用 Repo
	pageDO, err := s.repo.List(ctx, queryDO, common.PageParam{
		Page:     queryBO.Page,
		PageSize: queryBO.PageSize,
	})
	if err != nil {
		return nil, err
	}

	// DO Page -> BO Page
	pageBO := converter.PageDOToPageBO(pageDO, converter.PluginDOToBO)
	return pageBO, nil
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

```go
func (r *pluginRepository) List(ctx context.Context, query *dto.PluginQueryDTO, pageParam common.PageParam) (*common.PageResult[do.PluginMeta], error) {
	db := r.db.WithContext(ctx).Model(&do.PluginMeta{})

	// 动态条件
	if query.AuthorID != 0 {
		db = db.Where("author_id = ?", query.AuthorID)
	}
	if query.Tags != nil {
		db = db.Where("tag_id = ?", query.Tags)
	}
	if query.Type != "" {
		db = db.Where("post_type = ?", query.Type)
	}
	if query.Keyword != "" {
		db = db.Where("name LIKE ?", "%"+query.Keyword+"%")
	}
	if query.Status != "" {
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

