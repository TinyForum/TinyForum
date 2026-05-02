package user

import (
	"context"
	"tiny-forum/internal/model/request"
	"tiny-forum/internal/model/vo"
)

// service/user_service.go
func (s *userService) GetUserPosts(ctx context.Context, req request.GetUserPostsRequest, userID uint) (*vo.BasicPageData, error) {
	// 1. 排序规则白名单（职责：Service 层决定业务允许的排序方式）
	sortBy := s.resolveSortBy(req.SortBy, req.Order)

	// 2. 调用 Repo 获取帖子数据
	posts, total, err := s.postRepo.ListUserPosts(ctx, req, userID, sortBy)
	if err != nil {
		return nil, err
	}
	if total == 0 {
		return &vo.BasicPageData{
			List:     []vo.UserPosts{},
			Total:    0,
			Page:     req.Page,
			PageSize: req.PageSize,
			HasMore:  false,
		}, nil
	}

	// 3. 批量查询评论数（Service 层聚合其他数据）
	postIDs := make([]uint, len(posts))
	for i, p := range posts {
		postIDs[i] = p.ID
	}
	commentCounts, err := s.commentRepo.BatchCountByPostIDs(ctx, postIDs)
	if err != nil {
		return nil, err
	}

	// 4. 转换为 VO
	voList := make([]vo.UserPosts, 0, len(posts))
	for _, p := range posts {
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

	return &vo.BasicPageData{
		List:     voList,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
		HasMore:  int64(req.Page*req.PageSize) < total,
	}, nil
}

// 排序规则白名单（业务规则）
func (s *userService) resolveSortBy(sortBy, order string) string {
	if order != "asc" && order != "desc" {
		order = "desc"
	}
	switch sortBy {
	case "view_count":
		return "view_count " + order
	case "like_count":
		return "like_count " + order
	case "created_at":
		return "created_at " + order
	default:
		return "created_at DESC" // 默认按创建时间倒序
	}
}
