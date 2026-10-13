package user

import (
	"context"

	"tiny-forum/internal/model/common"
	"tiny-forum/internal/model/request"
	"tiny-forum/internal/model/vo"
)

// GetUserComments 获取当前用户发表的评论（分页）
func (s *userService) GetUserComments(ctx context.Context, req request.GetUserCommentsRequest, userID uint) (*common.PageResult[vo.UserComment], error) {
	comments, total, err := s.commentRepo.ListByAuthor(ctx, userID, req.Offset(), req.Limit())
	if err != nil {
		return nil, err
	}

	list := make([]vo.UserComment, 0, len(comments))
	for _, c := range comments {
		item := vo.UserComment{
			ID:        c.ID,
			CreatedAt: c.CreatedAt,
			WorksID:   c.WorksID,
		}
		if c.Reply != nil {
			item.Content = c.Reply.Content
			item.LikeCount = c.Reply.LikeCount
			item.Author = vo.UserPublicVO{
				ID:        c.Reply.Author.ID,
				Name:      c.Reply.Author.Username,
				AvatarUrl: c.Reply.Author.AvatarUrl,
			}
		}
		list = append(list, item)
	}

	return &common.PageResult[vo.UserComment]{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
		HasMore:  int64(req.Offset()+len(list)) < total,
	}, nil
}
