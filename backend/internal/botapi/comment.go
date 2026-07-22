package botapi

import (
	"context"
	"tiny-forum/internal/infra/lua/sdk"
)

// ─── Comment ──────────────────────────────────────────────────────────────

func (a *forumAPIImpl) GetComment(ctx context.Context, commentID uint) (*sdk.CommentVO, error) {
	c, err := a.commentRepo.FindByID(uint(commentID))
	if err != nil {
		return nil, err
	}
	return &sdk.CommentVO{
		ID:        c.ID,
		Content:   c.Content,
		AuthorID:  c.AuthorID,
		PostID:    c.CreationsID,
		CreatedAt: c.CreatedAt.Unix(),
	}, nil
}

func (a *forumAPIImpl) DeleteComment(ctx context.Context, commentID uint) error {
	return a.commentRepo.Delete(uint(commentID))
}
