package board

import (
	"context"
	"time"

	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/request"
	"tiny-forum/internal/model/vo"
	apperrors "tiny-forum/pkg/errors"
)

// snapshotContent 截断长文本为摘要，避免版主面板响应过大
func snapshotContent(s string) string {
	runes := []rune(s)
	if len(runes) > 500 {
		return string(runes[:500]) + "…"
	}
	return s
}

// ListModeratorPosts 版主面板：分页获取板块帖子（不做已发布过滤）
func (s *boardService) ListModeratorPosts(ctx context.Context, boardID uint, req *request.ListModeratorPostsRequest) ([]vo.ModeratorPostVO, int64, error) {
	status := do.CreationStatus(req.Status)
	if req.Status != "" && !status.IsValid() {
		return nil, 0, apperrors.ErrValidation
	}
	limit := req.PageSize
	offset := (req.Page - 1) * req.PageSize

	posts, total, err := s.postRepo.ListPostsByBoardID(ctx, boardID, req.Keyword, status, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	postIDs := make([]uint, 0, len(posts))
	for _, p := range posts {
		postIDs = append(postIDs, p.CreationID)
	}
	commentCounts, err := s.postRepo.CountCommentsByPostID(ctx, postIDs)
	if err != nil {
		return nil, 0, err
	}

	vos := make([]vo.ModeratorPostVO, 0, len(posts))
	for _, p := range posts {
		c := p.Creation
		v := vo.ModeratorPostVO{
			ID:           c.ID,
			Title:        c.Title,
			Content:      snapshotContent(c.Content),
			BoardID:      c.BoardID,
			IsPinned:     c.PinInBoard,
			ViewCount:    c.ViewCount,
			LikeCount:    c.LikeCount,
			CommentCount: commentCounts[c.ID],
			Status:       string(c.CreationStatus),
			CreatedAt:    c.CreatedAt,
		}
		if c.AuthorID != 0 {
			v.AuthorID = c.AuthorID
			v.AuthorName = c.Author.Username
			v.AuthorAvatar = c.Author.AvatarUrl
		}
		vos = append(vos, v)
	}
	return vos, total, nil
}

// ListModeratorReports 版主面板：分页获取板块举报（帖子或评论属于该板块）
func (s *boardService) ListModeratorReports(ctx context.Context, boardID uint, req *request.ListModeratorReportsRequest) ([]vo.ModeratorReportVO, int64, error) {
	status := do.ReportStatus(req.Status)
	if req.Status != "" && status != do.ReportPending && status != do.ReportResolved && status != do.ReportRejected {
		return nil, 0, apperrors.ErrValidation
	}
	limit := req.PageSize
	offset := (req.Page - 1) * req.PageSize

	reports, total, err := s.reportRepo.ListByBoardID(ctx, boardID, status, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	vos := make([]vo.ModeratorReportVO, 0, len(reports))
	for _, r := range reports {
		v := vo.ModeratorReportVO{
			ID:         r.ID,
			TargetType: string(r.TargetType),
			TargetID:   r.TargetID,
			Reason:     r.Reason,
			Status:     string(r.Status),
			CreatedAt:  r.CreatedAt,
		}
		if !r.IsAnonymous && r.Reporter != nil {
			v.ReporterID = r.Reporter.ID
			v.ReporterName = r.Reporter.Username
		}
		switch r.TargetType {
		case do.ReportTargetTypePost:
			post, err := s.postRepo.FindByArticleID(r.TargetID)
			if err == nil {
				v.TargetTitle = post.Creation.Title
				v.TargetContent = snapshotContent(post.Creation.Content)
			}
		case do.ReportTargetTypeComment:
			comment, err := s.commentRepo.FindByCommentID(r.TargetID)
			if err == nil {
				if comment.Reply != nil {
					v.TargetContent = snapshotContent(comment.Reply.Content)
				}
				post, perr := s.postRepo.FindByArticleID(comment.WorksID)
				if perr == nil {
					v.TargetTitle = post.Creation.Title
				}
			}
		}
		vos = append(vos, v)
	}
	return vos, total, nil
}

// ListModeratorBans 版主面板：分页获取板块封禁记录
func (s *boardService) ListModeratorBans(ctx context.Context, boardID uint, req *request.ListModeratorBansRequest) ([]vo.BannedUserVO, int64, error) {
	limit := req.PageSize
	offset := (req.Page - 1) * req.PageSize

	bans, total, err := s.boardRepo.ListBans(ctx, boardID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	now := time.Now()
	vos := make([]vo.BannedUserVO, 0, len(bans))
	for _, b := range bans {
		isActive := b.ExpiresAt == nil || b.ExpiresAt.After(now)
		v := vo.BannedUserVO{
			ID:        b.ID,
			UserID:    b.UserID,
			BoardID:   b.BoardID,
			Reason:    b.Reason,
			ExpiresAt: b.ExpiresAt,
			IsActive:  isActive,
			CreatedAt: b.CreatedAt,
		}
		if b.User != nil {
			v.Username = b.User.Username
			v.AvatarURL = b.User.AvatarUrl
		}
		vos = append(vos, v)
	}
	return vos, total, nil
}