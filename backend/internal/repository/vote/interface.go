package vote

import (
	"tiny-forum/internal/model"

	"gorm.io/gorm"
)

type VoteRepository interface {
	// crud
	CreateOrUpdateVote(commentID, userID uint, value int) error
	RemoveVote(commentID, userID uint) error
	// query
	GetUserVote(commentID, userID uint) (int, error)
	GetVoteCount(commentID uint) (int, error)
	GetVoteUsers(commentID uint, voteType int) ([]model.User, error)
}

type voteRepository struct {
	db *gorm.DB
}

func NewVoteRepository(db *gorm.DB) VoteRepository {
	return &voteRepository{db: db}
}
