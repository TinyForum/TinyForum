package botapi

import (
	"tiny-forum/internal/infra/lua/sdk"
	postrepo "tiny-forum/internal/repository/article"
	commentrepo "tiny-forum/internal/repository/comment"
	notificationrepo "tiny-forum/internal/repository/notification"
	userrepo "tiny-forum/internal/repository/user"
)

// forumAPIImpl 实现 luaengine.ForumAPI，对接现有 repo 层。
type forumAPIImpl struct {
	botActorID  uint // 机器人在论坛中以哪个用户身份操作（SystemBotID）
	postRepo    postrepo.ArticleRepository
	commentRepo commentrepo.CommentRepository
	userRepo    userrepo.UserRepository
	notifRepo   notificationrepo.NotificationRepository
}

// newForumAPI 构造 ForumAPI 实现，由 service 层在 executeBot 前调用。
func NewForumAPI(
	botActorID uint,
	postRepo postrepo.ArticleRepository,
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
