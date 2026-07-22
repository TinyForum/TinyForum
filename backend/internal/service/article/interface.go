package article

import (
	"context"
	"tiny-forum/internal/model/bo"
	"tiny-forum/internal/model/common"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/request"
	postRepo "tiny-forum/internal/repository/article"
	boardRepo "tiny-forum/internal/repository/board"
	tagRepo "tiny-forum/internal/repository/tag"
	userRepo "tiny-forum/internal/repository/user"

	"tiny-forum/internal/service/check"
	"tiny-forum/internal/service/notification"

	"github.com/gin-gonic/gin"
)

type ArticleService interface {
	// admin
	AdminLists(ctx context.Context, listPostsBO *common.PageQuery[bo.ListPosts]) ([]do.Article, int64, error)
	SetStatus(postID uint, status do.PostStatus) error
	TogglePin(postID uint) error
	AdminSetReviewPost(postID uint, status do.ModerationStatus) error
	// crud
	Create(ctx *gin.Context, authorID uint, input request.CreatePostRequest) (*do.Article, error)
	Update(postID, userID uint, isAdmin bool, input request.UpdatePostRequest) (*do.Article, error)
	Delete(postID, userID uint, isAdmin bool) error
	GetByID(postID, viewerID uint) (*do.Article, bool, error)
	// List(ctx context.Context, page, pageSize int, opts bo.ListPosts) ([]do.Post, int64, error)
	List(ctx context.Context, ListPostsBO *common.PageQuery[bo.ListPosts]) ([]do.Article, int64, error)
	// like
	Like(userID, postID uint) error
	Unlike(userID, postID uint) error
}

type articleService struct {
	postRepo  postRepo.ArticleRepository
	tagRepo   tagRepo.TagRepository
	boardRepo boardRepo.BoardRepository
	userRepo  userRepo.UserRepository
	notifSvc  notification.NotificationService
	// riskSvc         *risk.RiskService
	contentcheckSvc check.ContentCheckService
}

func NewPostService(
	postRepo postRepo.ArticleRepository,
	tagRepo tagRepo.TagRepository,
	userRepo userRepo.UserRepository,
	boardRepo boardRepo.BoardRepository,
	notifSvc notification.NotificationService,
	// riskSvc *risk.RiskService,
	contentcheckSvc check.ContentCheckService,
) ArticleService {
	return &articleService{
		postRepo:  postRepo,
		tagRepo:   tagRepo,
		userRepo:  userRepo,
		boardRepo: boardRepo,
		notifSvc:  notifSvc,
		// riskSvc:         riskSvc,
		contentcheckSvc: contentcheckSvc,
	}
}
