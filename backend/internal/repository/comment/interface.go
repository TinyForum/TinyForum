package comment

import (
	"context"
	"time"
	"tiny-forum/internal/model/do"

	"gorm.io/gorm"
)

type CommentRepository interface {
	Create(comment *do.Comment) error
	FindByID(id uint) (*do.Comment, error)
	Update(comment *do.Comment) error
	Delete(id uint) error
	ValidateParentComment(parentID uint, postID uint) error
	// stats
	CountByPost(postID uint) (int64, error)
	Count(ctx context.Context) (int64, error)
	CountByDateRange(ctx context.Context, startDate, endDate time.Time) (int64, error)
	IncrLikeCount(id uint, delta int) error
	UpdateVoteCount(commentID uint, voteCount int) error
	// list
	ListByPost(postID uint, page, pageSize int) ([]do.Comment, int64, error)
	GetAnswersByPostID(postID uint, limit, offset int) ([]do.Comment, int64, error)
	GetAnswersByPostIDOrderByNewest(postID uint, limit, offset int) ([]do.Comment, int64, error)
	GetAnswersByPostIDOrderByOldest(postID uint, limit, offset int) ([]do.Comment, int64, error)
	// answer
	MarkAsAccepted(commentID uint) error
	MarkAsAnswer(commentID uint, isAnswer bool) error
	UnacceptAnswer(commentID uint) error
	GetAcceptedAnswer(postID uint) (*do.Comment, error)
}
type commentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepository{db: db}
}
