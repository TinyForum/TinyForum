package wire

import (
	"tiny-forum/internal/repository/announcement"
	"tiny-forum/internal/repository/auth"
	"tiny-forum/internal/repository/board"
	"tiny-forum/internal/repository/comment"
	"tiny-forum/internal/repository/notification"
	"tiny-forum/internal/repository/post"
	"tiny-forum/internal/repository/question"
	"tiny-forum/internal/repository/risk"
	"tiny-forum/internal/repository/stats"
	"tiny-forum/internal/repository/tag"
	"tiny-forum/internal/repository/timeline"
	"tiny-forum/internal/repository/token"
	"tiny-forum/internal/repository/topic"
	"tiny-forum/internal/repository/transaction"
	"tiny-forum/internal/repository/user"
	"tiny-forum/internal/repository/vote"

	"gorm.io/gorm"
)

// Repositories 聚合所有 Repository
type Repositories struct {
	Token        token.TokenRepository
	User         user.UserRepository
	Post         post.PostRepository
	Comment      comment.CommentRepository
	Tag          tag.TagRepository
	Notification notification.NotificationRepository
	Board        board.BoardRepository
	Timeline     *timeline.TimelineRepository
	Topic        *topic.TopicRepository
	Question     question.QuestionRepository
	Vote         *vote.VoteRepository
	Announcement announcement.AnnouncementRepository
	Stats        stats.StatsRepository
	Auth         auth.AuthRepository
	Risk         risk.RiskRepository
	Transaction  transaction.TransactionManager
}

// NewRepositories 创建所有 Repository 实例
func NewRepositories(db *gorm.DB) *Repositories {
	tokenRepo := token.NewTokenRepository(db)

	return &Repositories{
		Token:        tokenRepo,
		User:         user.NewUserRepository(db, tokenRepo),
		Post:         post.NewPostRepository(db),
		Comment:      comment.NewCommentRepository(db),
		Tag:          tag.NewTagRepository(db),
		Notification: notification.NewNotificationRepository(db),
		Board:        board.NewBoardRepository(db),
		Timeline:     timeline.NewTimelineRepository(db),
		Topic:        topic.NewTopicRepository(db),
		Question:     question.NewQuestionRepository(db),
		Vote:         vote.NewVoteRepository(db),
		Announcement: announcement.NewAnnouncementRepository(db),
		Stats:        stats.NewStatsRepository(db),
		Auth:         auth.NewAuthRepository(db),
		Risk:         risk.NewRiskRepository(db),
		Transaction:  transaction.NewTransactionManager(db),
	}
}
