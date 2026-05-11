package user

import (
	"context"
	"tiny-forum/internal/model/common"
	"tiny-forum/internal/model/request"
	"tiny-forum/internal/model/vo"
)

func (s *userService) ListUserViolation(ctx context.Context, req request.ListUserViolationRequest, userID uint) (*common.PageResult[vo.UserPosts], error) {
	// 1. 排序规则白名单（职责：Service 层决定业务允许的排序方式）

	// 2. 调用 Repo 获取帖子数据
	violations, err := s.violationSvc.ListUserViolation(ctx, req, userID)
	if err != nil {
		return nil, err
	}
	if len(violations) == 0 {

		return &common.PageResult[vo.UserPosts]{
			List:     []vo.UserPosts{},
			Total:    0,
			Page:     req.Page,
			PageSize: req.PageSize,
			HasMore:  false,
		}, nil
	}

	// 3. 批量查询评论数（Service 层聚合其他数据）
	postIDs := make([]uint, len(violations))
	for i, p := range violations {
		postIDs[i] = p.ID
	}
	commentCounts, err := s.commentRepo.BatchCountByPostIDs(ctx, postIDs)
	if err != nil {
		return nil, err
	}

	// 4. 转换为 VO
	voList := make([]vo.UserPosts, 0, len(violations))
	for _, p := range violations {
		tagNames := make([]string, len(p.Tags))
		for i, t := range p.Tags {
			tagNames[i] = t.Name
		}
		boardName := ""
		if p.Board.ID != 0 {
			boardName = p.Board.Name
		}
		voList = append(voList, vo.UserPosts{
			ID:               p.ID,
			CreatedAt:        p.CreatedAt,
			UpdatedAt:        p.UpdatedAt,
			Title:            p.Title,
			Summary:          p.Summary,
			Cover:            p.Cover,
			Type:             p.Type,
			PostStatus:       p.PostStatus,
			ModerationStatus: p.ModerationStatus,
			ViewCount:        p.ViewCount,
			LikeCount:        p.LikeCount,
			CommentCount:     commentCounts[p.ID],
			PinTop:           p.PinTop,
			Tags:             tagNames,
			BoardName:        boardName,
			PinInBoard:       p.PinInBoard,
		})
	}

	return &common.PageResult[vo.UserPosts]{
		List:     voList,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
		HasMore:  int64(req.Page*req.PageSize) < total,
	}, nil
}
