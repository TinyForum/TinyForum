package botapi

import (
	"context"
	"fmt"
	"tiny-forum/internal/infra/lua/sdk"
	"tiny-forum/internal/model/do"
	"tiny-forum/pkg/utils"
)

// ─── Post ─────────────────────────────────────────────────────────────────

func (a *forumAPIImpl) GetPost(ctx context.Context, postID uint) (*sdk.PostVO, error) {
	p, err := a.postRepo.FindByArticleID(uint(postID))
	if err != nil {
		return nil, err
	}
	return &sdk.PostVO{

		ID:        p.ID,
		Title:     p.Creation.Title,
		Content:   p.Creation.Content,
		AuthorID:  p.Creation.AuthorID,
		BoardID:   p.Creation.BoardID,
		CreatedAt: p.CreatedAt.Unix(),
	}, nil
}

func (a *forumAPIImpl) CreatePost(ctx context.Context, req sdk.CreatePostReq) (*sdk.PostVO, error) {
	p := &do.Article{
		Creation: do.Creation{
			Title:          req.Title,
			Content:        req.Content,
			Slug:           utils.GenerateSlug(),
			AuthorID:       a.botActorID,
			BoardID:        req.BoardID,
			CreationStatus: do.CreationStatusPublished,
		},
	}
	if err := a.postRepo.Create(p); err != nil {
		return nil, err
	}
	return &sdk.PostVO{
		ID:        p.ID,
		Title:     p.Creation.Title,
		Content:   p.Creation.Content,
		AuthorID:  p.Creation.AuthorID,
		BoardID:   p.Creation.BoardID,
		CreatedAt: p.CreatedAt.Unix(),
	}, nil
}

func (a *forumAPIImpl) ReplyPost(ctx context.Context, postID uint, content string) (*sdk.CommentVO, error) {
	c := &do.Comment{
		WorksID: uint(postID),
		Reply: &do.Reply{
			AuthorID: a.botActorID,
			Content:  content,
			Status:   do.ReplyStatusVisible,
			TargetID: uint(postID),
		},
	}
	if err := a.commentRepo.CreateComment(c); err != nil {
		return nil, err
	}
	return &sdk.CommentVO{
		ID:        c.ID,
		Content:   c.Reply.Content,
		AuthorID:  c.Reply.AuthorID,
		TargetID:  c.Reply.TargetID,
		CreatedAt: c.CreatedAt.Unix(),
	}, nil
}

func (a *forumAPIImpl) DeletePost(ctx context.Context, postID uint) error {
	attachments, _ := a.attachmentRepo.FindByPostID(ctx, int64(postID))
	for _, att := range attachments {
		_ = a.attachmentRepo.SoftDelete(ctx, att.FileID)
	}
	return a.postRepo.Delete(uint(postID))
}

// ModeratePost 支持 action: hide | pin | delete（lock 需数据库 schema 支持，暂不支持）
func (a *forumAPIImpl) ModeratePost(ctx context.Context, postID uint, action, reason string) error {
	p, err := a.postRepo.FindByArticleID(uint(postID))
	if err != nil {
		return err
	}
	switch action {
	case "hide":
		p.Creation.CreationStatus = do.CreationStatusHidden
		return a.postRepo.Update(p)
	case "pin":
		return a.postRepo.TogglePinInBoard(uint(postID), true)
	case "lock":
		// 帖子模型当前没有 locked 字段，无法真正实现"锁定（禁止评论）"。
		// 明确返回错误，禁止降级为隐藏（那会错误地把帖子内容隐藏掉）。
		return fmt.Errorf("moderate action 'lock' is not supported: post model has no lock field")
	case "delete":
		return a.DeletePost(ctx, postID)
	default:
		return fmt.Errorf("unknown moderate action: %s", action)
	}
}
