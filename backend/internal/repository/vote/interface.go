package vote

import (
	"tiny-forum/internal/model/do"

	"gorm.io/gorm"
)

type VoteRepository interface {
	// crud
	CreateOrUpdateVote(commentID, userID uint, value do.AnswerVoteType) error
	RemoveVote(commentID, userID uint) error
	// query
	GetUserVote(commentID, userID uint) (*do.AnswerVoteType, error)
	GetVoteCount(commentID uint) (int, error)
	GetVoteUsers(commentID uint, voteType do.AnswerVoteType) ([]do.User, error)
}

type voteRepository struct {
	db *gorm.DB
}

func NewVoteRepository(db *gorm.DB) VoteRepository {
	return &voteRepository{db: db}
}
