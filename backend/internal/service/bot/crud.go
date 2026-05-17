package bot

import (
	"context"
	"errors"
	"fmt"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/request"
	"tiny-forum/internal/model/vo"
)

// ─── CRUD ─────────────────────────────────────────────────────────────────

func (s *service) Create(ctx context.Context, creatorID uint, req *request.CreateBotRequest) (*do.Bot, error) {
	// 零代码模式：预先校验 Flow JSON
	if req.ScriptCode == "" && req.ConfigValues != nil {
		if flowRaw, ok := req.ConfigValues["flow"]; ok {
			reqFlow := parseFlowRequestRaw(flowRaw)
			if errs := s.ValidateFlowRequest((*request.ValidateFlowRequest)(toFlow((*request.ValidateFlowRequest)(reqFlow)))); len(errs) > 0 {
				return nil, fmt.Errorf("invalid nocode flow: %v", errs[0])
			}
		}
	}

	// 查询创建者用户名
	creatorName := "user"
	if u, err := s.userRepo.FindByID(creatorID); err == nil {
		creatorName = u.Username
	}

	bot := &do.Bot{
		Name:          req.Name,
		Version:       req.Version,
		Description:   req.Description,
		Summary:       req.Summary,
		AvatarURL:     req.AvatarURL,
		Screenshots:   orStrSlice(req.Screenshots),
		HomepageURL:   req.HomepageURL,
		Type:          req.Type,
		Tags:          orStrSlice(req.Tags),
		CreatorID:     creatorID,
		CreatorName:   creatorName,
		ScriptCode:    req.ScriptCode,
		ScriptURL:     req.ScriptURL,
		TriggerType:   req.TriggerType,
		CronExpr:      req.CronExpr,
		EventFilter:   req.EventFilter,
		TimeoutSec:    req.TimeoutSec,
		RetryTimes:    req.RetryTimes,
		EnvVars:       orStrMap(req.EnvVars),
		ResourceLimit: req.ResourceLimit,
		Pricing:       req.Pricing,
		Permissions:   req.Permissions,
		Enabled:       false,
		Status:        do.BotStatusInactive,
		ConfigSchema:  req.ConfigSchema,
		ConfigValues:  orAnyMap(req.ConfigValues),
	}
	if bot.TimeoutSec == 0 {
		bot.TimeoutSec = 10
	}
	if err := s.repo.Create(ctx, bot); err != nil {
		return nil, err
	}
	return bot, nil
}

func (s *service) Update(ctx context.Context, userID uint, botID uint, req *request.UpdateBotRequest) error {
	bot, err := s.repo.GetByID(ctx, botID)
	if err != nil {
		return err
	}
	if bot.CreatorID != userID {
		return errors.New("permission denied")
	}

	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Summary != nil {
		updates["summary"] = *req.Summary
	}
	if req.AvatarURL != nil {
		updates["avatar_url"] = *req.AvatarURL
	}
	if req.ScriptCode != nil {
		updates["script_code"] = *req.ScriptCode
	}
	if req.ScriptURL != nil {
		updates["script_url"] = *req.ScriptURL
	}
	if req.TriggerType != nil {
		updates["trigger_type"] = *req.TriggerType
	}
	if req.CronExpr != nil {
		updates["cron_expr"] = *req.CronExpr
	}
	if req.EventFilter != nil {
		updates["event_filter"] = *req.EventFilter
	}
	if req.TimeoutSec != nil {
		updates["timeout_sec"] = *req.TimeoutSec
	}
	if req.RetryTimes != nil {
		updates["retry_times"] = *req.RetryTimes
	}
	if req.EnvVars != nil {
		updates["env_vars"] = req.EnvVars
	}
	if req.ResourceLimit != nil {
		updates["resource_limit"] = req.ResourceLimit
	}
	if req.Pricing != nil {
		updates["pricing"] = req.Pricing
	}
	if req.Permissions != nil {
		updates["permissions"] = req.Permissions
	}
	if req.ConfigSchema != nil {
		updates["config_schema"] = req.ConfigSchema
	}
	if req.ConfigValues != nil {
		updates["config_values"] = req.ConfigValues
	}
	if req.Enabled != nil {
		updates["enabled"] = *req.Enabled
		if *req.Enabled {
			updates["status"] = do.BotStatusActive
		} else {
			updates["status"] = do.BotStatusInactive
		}
	}
	return s.repo.Update(ctx, botID, updates)
}

func (s *service) Delete(ctx context.Context, userID uint, botID uint) error {
	bot, err := s.repo.GetByID(ctx, botID)
	if err != nil {
		return err
	}
	if bot.CreatorID != userID {
		return errors.New("permission denied")
	}
	return s.repo.Delete(ctx, botID)
}

func (s *service) Get(ctx context.Context, id uint) (*vo.BotResponse, error) {
	bot, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toResponse(bot), nil
}

func (s *service) ListByUser(ctx context.Context, userID uint, page, pageSize int) ([]*vo.BotResponse, int64, error) {
	offset := (page - 1) * pageSize
	bots, total, err := s.repo.ListByUser(ctx, userID, offset, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return mapToResponse(bots), total, nil
}

func (s *service) List(ctx context.Context, page, pageSize int) ([]*vo.BotResponse, int64, error) {
	offset := (page - 1) * pageSize
	bots, total, err := s.repo.List(ctx, offset, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return mapToResponse(bots), total, nil
}
