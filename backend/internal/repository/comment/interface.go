package comment

import (
	"context"
	"time"
	"tiny-forum/internal/model/do"

	"gorm.io/gorm"
)

type CommentRepository interface {
	CreateComment(comment *do.Comment) error
	CreateAnswer(comment *do.Answer) error
	FindByCommentID(id uint) (*do.Comment, error)
	FindByAnswerID(id uint) (*do.Answer, error)
	Update(comment *do.Comment) error
	Delete(id uint) error
	ValidateParentComment(parentID uint, postID uint) error
	// stats
	CountByPost(postID uint) (int64, error)
	Count(ctx context.Context) (int64, error)
	CountByDateRange(ctx context.Context, startDate, endDate time.Time) (int64, error)
	UpdateVoteCount(commentID uint, voteCount int) error
	// list
	ListByPost(postID uint, page, pageSize int) ([]do.Comment, int64, error)
	// GetAnswersByPostID(postID uint, limit, offset int) ([]do.Answer, int64, error)
	GetAnswersByQuestionID(questionID uint, limit, offset int) ([]do.Answer, int64, error)
	GetAnswersByPostIDOrderByNewest(postID uint, limit, offset int) ([]do.Answer, int64, error)
	GetAnswersByPostIDOrderByOldest(postID uint, limit, offset int) ([]do.Answer, int64, error)
	// answer
	MarkAsAccepted(commentID uint) error
	MarkAsAnswer(commentID uint, isAnswer bool) error
	UnacceptAnswer(commentID uint) error
	GetAcceptedAnswer(postID uint) (*do.Answer, error)
	// 查询评论数
	BatchCountByPostIDs(ctx context.Context, postIDs []uint) (map[uint]int64, error)

	GetCommentTree(postID uint) ([]do.Comment, error)

	// count
	IncrViewCount(id uint) error
	IncrLikeCount(id uint, delta int) error

	// like
	AddLike(userID, articleID uint) error
	RemoveLike(userID, postID uint) error
	IsLiked(userID, postID uint) bool
}

type commentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepository{db: db}
}
