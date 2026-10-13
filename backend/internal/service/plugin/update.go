package plugin

import (
	"context"

	"tiny-forum/internal/model/request"
	apperrors "tiny-forum/pkg/errors"
)

func (s *pluginService) TogglePluginStatus(ctx context.Context, pluginSlug string) error {
	return s.repo.TogglePluginStatus(ctx, pluginSlug)
}

// UpdatePlugin 更新插件元数据（管理员）
func (s *pluginService) UpdatePlugin(ctx context.Context, pluginSlug string, req *request.UpdatePluginRequest) error {
	if pluginSlug == "" || req == nil {
		return apperrors.ErrValidation
	}
	p, err := s.repo.GetBySlug(ctx, pluginSlug)
	if err != nil {
		return err
	}
	if req.Name != "" {
		p.Name = req.Name
	}
	if req.Version != "" {
		p.Version = req.Version
	}
	if req.Description != "" {
		p.Description = req.Description
	}
	if req.Summary != "" {
		p.Summary = req.Summary
	}
	if req.Enabled != nil {
		p.Enabled = *req.Enabled
	}
	if req.Slots != nil {
		p.Slots = req.Slots
	}
	if req.Config != nil {
		p.Config = req.Config
	}
	return s.repo.Update(ctx, p)
}
