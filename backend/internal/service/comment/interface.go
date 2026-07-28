package comment

import (
	"tiny-forum/internal/model/bo"
	"tiny-forum/internal/model/do"
	postRepo "tiny-forum/internal/repository/article"
	commentRepo "tiny-forum/internal/repository/comment"
	userRepo "tiny-forum/internal/repository/user"
	voteRepo "tiny-forum/internal/repository/vote"
	"tiny-forum/internal/service/notification"
)

type CommentService interface {
	// answer
	MarkAsAnswer(commentID, userID uint, isAdmin bool, isAnswer bool) error
	UnacceptAnswer(answerID, userID uint, isAdmin bool) error
	// create
	CreateComment(authorID uint, input bo.CreateCommentInput) (*do.Comment, error)
	CreateAnswer(authorID uint, input bo.CreateAnswerInput) (*do.Answer, error)
	// delete
	Delete(commentID, userID uint, isAdmin bool) error
	DeleteAnswer(commentID, userID uint, isAdmin bool) error
	// query
	List(postID uint, page, pageSize int) ([]do.Comment, int64, error)
	GetCommentCount(postID uint) (int64, error)
	GetAnswerByID(commentID uint) (*do.Answer, error)
	GetAnswersByPostID(postID uint, page, pageSize int, sortBy string) ([]do.Answer, int64, error)
	GetAnswerVoteCount(commentID uint) (int, error)
	GetVoteStatistics(answerID uint) (upCount, downCount int, err error)
	// vote
	VoteAnswer(answerID uint, userID uint, voteType do.AnswerVoteType) (*do.Answer, error)
	RemoveVote(answerID uint, userID uint) (*do.Answer, error)
	GetUserVoteStatus(answerID uint, userID uint) (*do.AnswerVoteType, error)

	// like
	Like(userID, postID uint) error
	Unlike(userID, postID uint) error
}

type commentService struct {
	commentRepo commentRepo.CommentRepository
	postRepo    postRepo.ArticleRepository
	userRepo    userRepo.UserRepository
	notifSvc    notification.NotificationService
	voteRepo    voteRepo.VoteRepository
}

func NewCommentService(
	commentRepo commentRepo.CommentRepository,
	postRepo postRepo.ArticleRepository,
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
