package plugin

import (
	"context"
	"mime/multipart"

	"tiny-forum/internal/infra/config"
	"tiny-forum/internal/model/bo"
	"tiny-forum/internal/model/common"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/vo"
	pluginRepo "tiny-forum/internal/repository/plugin"
	"tiny-forum/internal/storage"
	apperrors "tiny-forum/pkg/errors"
	"tiny-forum/pkg/logger"
)

type pluginService struct {
	repo    pluginRepo.PluginRepository
	storage storage.StorageDriver
	cfg     *config.ConfigPlugins
}
type PluginService interface {
	ListPlugins(ctx context.Context, queryBO *bo.PageQuery[bo.PluginQueryBO]) (*common.PageResult[vo.PluginMetaVO], error)
	Create(ctx context.Context, fileHeader *multipart.FileHeader, userID uint) (*do.PluginMeta, error)
	ListUserPlugins(ctx context.Context, userID uint) ([]*do.PluginMeta, error)
	DeletePlugin(ctx context.Context, pluginID uint,userID uint) error
}

func NewPluginService(repo pluginRepo.PluginRepository, storage storage.StorageDriver, cfg *config.ConfigPlugins) PluginService {
	return &pluginService{repo: repo, storage: storage,cfg :cfg}
}


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



// ListUserPlugins 获取用户安装的插件
func (s *pluginService) ListUserPlugins(ctx context.Context, userID uint) ([]*do.PluginMeta, error) {
	return s.repo.ListByAuthorID(ctx, userID)
}
