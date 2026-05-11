package botapi

import (
	"context"
	"fmt"
	"time"

	"tiny-forum/internal/infra/lua/sdk"
	"tiny-forum/internal/model/do"
	commentrepo "tiny-forum/internal/repository/comment"
	notificationrepo "tiny-forum/internal/repository/notification"
	postrepo "tiny-forum/internal/repository/post"
	userrepo "tiny-forum/internal/repository/user"
)

// forumAPIImpl 实现 luaengine.ForumAPI，对接现有 repo 层。
type forumAPIImpl struct {
	botActorID  uint // 机器人在论坛中以哪个用户身份操作（SystemBotID）
	postRepo    postrepo.PostRepository
	commentRepo commentrepo.CommentRepository
	userRepo    userrepo.UserRepository
	notifRepo   notificationrepo.NotificationRepository
}

// newForumAPI 构造 ForumAPI 实现，由 service 层在 executeBot 前调用。
func NewForumAPI(
	botActorID uint,
	postRepo postrepo.PostRepository,
	commentRepo commentrepo.CommentRepository,
	userRepo userrepo.UserRepository,
	notifRepo notificationrepo.NotificationRepository,
) sdk.ForumAPI {
	return &forumAPIImpl{
		botActorID:  botActorID,
		postRepo:    postRepo,
		commentRepo: commentRepo,
		userRepo:    userRepo,
		notifRepo:   notifRepo,
	}
}

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
		PostID:    c.PostID,
		CreatedAt: c.CreatedAt.Unix(),
	}, nil
}

func (a *forumAPIImpl) DeleteComment(ctx context.Context, commentID uint) error {
	return a.commentRepo.Delete(uint(commentID))
}

// ─── User ─────────────────────────────────────────────────────────────────

func (a *forumAPIImpl) GetUser(ctx context.Context, userID uint) (*sdk.UserVO, error) {
	u, err := a.userRepo.FindByID(uint(userID))
	if err != nil {
		return nil, err
	}
	return &sdk.UserVO{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
		Role:     string(u.Role),
	}, nil
}

// BanUser 通过 UserRepository.UpdateBlocked 封禁用户，并发送通知。
// durationSec 目前实现为永久封禁（UpdateBlocked），如需定时解封可后续扩展。
func (a *forumAPIImpl) BanUser(ctx context.Context, userID uint, reason string, durationSec int) error {
	uid := uint(userID)
	if err := a.userRepo.UpdateBlocked(ctx, uid, true); err != nil {
		return err
	}
	// 发送封禁通知
	until := time.Now().Add(time.Duration(durationSec) * time.Second)
	notif := &do.Notification{
		UserID:     uid,
		Type:       do.NotifyBan,
		Content:    fmt.Sprintf("您的账号因 [%s] 已被封禁至 %s", reason, until.Format("2006-01-02 15:04")),
		TargetType: "user",
	}
	return a.notifRepo.Create(notif)
}

// ─── Message ──────────────────────────────────────────────────────────────

// SendMessage 通过 NotificationRepository 创建私信通知。
func (a *forumAPIImpl) SendMessage(ctx context.Context, toUserID uint, content string) error {
	senderID := a.botActorID
	notif := &do.Notification{
		UserID:     uint(toUserID),
		SenderID:   &senderID,
		Type:       do.NotifyPrivateMessage,
		Content:    content,
		TargetType: "bot_message",
	}
	return a.notifRepo.Create(notif)
}

// ─── Stats ────────────────────────────────────────────────────────────────

func (a *forumAPIImpl) GetForumStats(ctx context.Context) (*sdk.StatsVO, error) {
	postCount, err := a.postRepo.Count(ctx)
	if err != nil {
		return nil, err
	}
	userCount, err := a.userRepo.Count(ctx)
	if err != nil {
		return nil, err
	}
	commentCount, err := a.commentRepo.Count(ctx)
	if err != nil {
		return nil, err
	}
	// 活跃用户（近 24h 内有注册的用户数作为近似值）
	yesterday := time.Now().Add(-24 * time.Hour)
	activeToday, _ := a.userRepo.CountActiveByDateRange(ctx, yesterday, time.Now())

	return &sdk.StatsVO{
		PostCount:    postCount,
		UserCount:    userCount,
		CommentCount: commentCount,
		ActiveToday:  activeToday,
	}, nil
}
