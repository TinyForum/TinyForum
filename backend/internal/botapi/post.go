package botapi

import (
	"context"
	"fmt"
	"tiny-forum/internal/infra/lua/sdk"
	"tiny-forum/internal/model/do"
)

// ─── Post ─────────────────────────────────────────────────────────────────

func (a *forumAPIImpl) GetPost(ctx context.Context, postID uint) (*sdk.PostVO, error) {
	p, err := a.postRepo.FindByID(uint(postID))
	if err != nil {
		return nil, err
	}
	return &sdk.PostVO{
		ID:        p.ID,
		Title:     p.Title,
		Content:   p.Content,
		AuthorID:  p.AuthorID,
		BoardID:   p.BoardID,
		CreatedAt: p.CreatedAt.Unix(),
	}, nil
}

func (a *forumAPIImpl) CreatePost(ctx context.Context, req sdk.CreatePostReq) (*sdk.PostVO, error) {
	p := &do.Post{
		Title:      req.Title,
		Content:    req.Content,
		AuthorID:   a.botActorID,
		BoardID:    req.BoardID,
		PostStatus: do.PostStatusPublished,
	}
	if err := a.postRepo.Create(p); err != nil {
		return nil, err
	}
	return &sdk.PostVO{
		ID:        p.ID,
		Title:     p.Title,
		Content:   p.Content,
		AuthorID:  p.AuthorID,
		BoardID:   p.BoardID,
		CreatedAt: p.CreatedAt.Unix(),
	}, nil
}

func (a *forumAPIImpl) ReplyPost(ctx context.Context, postID uint, content string) (*sdk.CommentVO, error) {
	c := &do.Comment{
		PostID:   uint(postID),
		AuthorID: a.botActorID,
		Content:  content,
		Status:   do.CommentStatusVisible,
	}
	if err := a.commentRepo.Create(c); err != nil {
		return nil, err
	}
	return &sdk.CommentVO{
		ID:        c.ID,
		Content:   c.Content,
		AuthorID:  c.AuthorID,
		PostID:    c.PostID,
		CreatedAt: c.CreatedAt.Unix(),
	}, nil
}

func (a *forumAPIImpl) DeletePost(ctx context.Context, postID uint) error {
	return a.postRepo.Delete(uint(postID))
}

// ModeratePost 支持 action: hide | pin | lock | delete
func (a *forumAPIImpl) ModeratePost(ctx context.Context, postID uint, action, reason string) error {
	p, err := a.postRepo.FindByID(uint(postID))
	if err != nil {
		return err
	}
	switch action {
	case "hide":
		p.PostStatus = do.PostStatusHidden
		return a.postRepo.Update(p)
	case "pin":
		return a.postRepo.TogglePinInBoard(uint(postID), true)
	case "lock":
		// do.Post 没有 locked 字段，用 Hidden 作降级处理
		// 如有需要可扩展 Post 模型
		p.PostStatus = do.PostStatusHidden
		return a.postRepo.Update(p)
	case "delete":
		return a.postRepo.Delete(uint(postID))
	default:
		return fmt.Errorf("unknown moderate action: %s", action)
	}
}
