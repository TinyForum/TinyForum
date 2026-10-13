package user

import (
	"context"
	"errors"

	"tiny-forum/internal/model/common"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/request"
	"tiny-forum/internal/model/vo"
	apperrors "tiny-forum/pkg/errors"

	"gorm.io/gorm"
)

func (s *userService) GetUserPosts(ctx context.Context, req request.GetUserPostsRequest, userID uint) (*common.PageResult[vo.UserPosts], error) {
	// 1. 排序规则白名单（职责：Service 层决定业务允许的排序方式）
	sortBy := s.resolveSortBy(req.SortBy, req.Order)

	// 2. 调用 Repo 获取帖子数据
	posts, total, err := s.postRepo.ListUserPosts(ctx, req, userID, sortBy)
	if err != nil {
		return nil, err
	}
	if total == 0 {

		return &common.PageResult[vo.UserPosts]{
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
		voList = append(voList, s.toUserPostsVO(p, commentCounts[p.ID]))
	}

	return &common.PageResult[vo.UserPosts]{
		List:     voList,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
		HasMore:  int64(req.Page*req.PageSize) < total,
	}, nil
}

// toUserPostsVO 将帖子聚合视图转换为前端视图对象
func (s *userService) toUserPostsVO(p do.Article, commentCount int64) vo.UserPosts {
	tagNames := make([]string, 0, len(p.Creation.Tags))
	for _, t := range p.Creation.Tags {
		tagNames = append(tagNames, t.Name)
	}
	boardName := ""
	if p.Creation.Board.ID != 0 {
		boardName = p.Creation.Board.Name
	}
	return vo.UserPosts{
		ID:               p.ID,
		CreatedAt:        p.CreatedAt,
		UpdatedAt:        p.UpdatedAt,
		Title:            p.Creation.Title,
		Summary:          p.Creation.Summary,
		Cover:            p.Creation.CoverUrl,
		Type:             p.Creation.Type,
		CreationStatus:   p.Creation.CreationStatus,
		ModerationStatus: p.Creation.ModerationStatus,
		ViewCount:        p.Creation.ViewCount,
		LikeCount:        p.Creation.LikeCount,
		CommentCount:     commentCount,
		PinTop:           p.Creation.PinTop,
		Tags:             tagNames,
		BoardName:        boardName,
		PinInBoard:       p.Creation.PinInBoard,
	}
}

// GetUserPost 获取当前用户的单个帖子
func (s *userService) GetUserPost(ctx context.Context, postID, userID uint) (*vo.UserPosts, error) {
	post, err := s.postRepo.FindByArticleID(postID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrPostNotFound
		}
		return nil, err
	}
	if post.Creation.AuthorID != userID {
		return nil, apperrors.ErrForbidden.WithDetail("无权查看该帖子")
	}
	commentCounts, err := s.commentRepo.BatchCountByPostIDs(ctx, []uint{postID})
	if err != nil {
		return nil, err
	}
	v := s.toUserPostsVO(*post, commentCounts[postID])
	return &v, nil
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
