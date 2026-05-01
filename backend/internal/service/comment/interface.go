package comment

import (
	"tiny-forum/internal/model/do"
	commentRepo "tiny-forum/internal/repository/comment"
	postRepo "tiny-forum/internal/repository/post"
	userRepo "tiny-forum/internal/repository/user"
	voteRepo "tiny-forum/internal/repository/vote"
	"tiny-forum/internal/service/notification"
)

type CommentService interface {
	// answer
	MarkAsAnswer(commentID, userID uint, isAdmin bool, isAnswer bool) error
	UnacceptAnswer(answerID, userID uint, isAdmin bool) error
	// create
	Create(authorID uint, input CreateCommentInput) (*do.Comment, error)
	CreateAnswer(authorID uint, input CreateCommentInput) (*do.Comment, error)
	// delete
	Delete(commentID, userID uint, isAdmin bool) error
	DeleteAnswer(commentID, userID uint, isAdmin bool) error
	// query
	List(postID uint, page, pageSize int) ([]do.Comment, int64, error)
	GetCommentCount(postID uint) (int64, error)
	GetAnswerByID(commentID uint) (*do.Comment, error)
	GetAnswersByPostID(postID uint, page, pageSize int, sortBy string) ([]do.Comment, int64, error)
	GetAnswerVoteCount(commentID uint) (int, error)
	GetVoteStatistics(answerID uint) (upCount, downCount int, err error)
	// vote
	VoteAnswer(answerID uint, userID uint, voteType int) (*do.Comment, error)
	RemoveVote(answerID uint, userID uint) (*do.Comment, error)
	GetUserVoteStatus(answerID uint, userID uint) (int, error)
}

type commentService struct {
	commentRepo commentRepo.CommentRepository
	postRepo    postRepo.PostRepository
	userRepo    userRepo.UserRepository
	notifSvc    notification.NotificationService
	voteRepo    voteRepo.VoteRepository
}

func NewCommentService(
	commentRepo commentRepo.CommentRepository,
	postRepo postRepo.PostRepository,
	userRepo userRepo.UserRepository,
	notifSvc notification.NotificationService,
	voteRepo voteRepo.VoteRepository,
) CommentService {
	return &commentService{
		commentRepo: commentRepo,
		postRepo:    postRepo,
		userRepo:    userRepo,
		notifSvc:    notifSvc,
		voteRepo:    voteRepo,
	}
}
