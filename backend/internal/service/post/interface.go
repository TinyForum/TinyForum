package post

import (
	"tiny-forum/internal/model/dto"
	"tiny-forum/internal/model/po"
	boardRepo "tiny-forum/internal/repository/board"
	postRepo "tiny-forum/internal/repository/post"
	tagRepo "tiny-forum/internal/repository/tag"
	userRepo "tiny-forum/internal/repository/user"

	"tiny-forum/internal/service/check"
	"tiny-forum/internal/service/notification"

	"github.com/gin-gonic/gin"
)

type PostService interface {
	AdminList(page, pageSize int, opts dto.PostListOptions) ([]po.Post, int64, error)
	SetStatus(postID uint, status po.PostStatus) error
	TogglePin(postID uint) error
	AdminSetReviewPost(postID uint, status po.ModerationStatus) error
	// crud
	Create(ctx *gin.Context, authorID uint, input CreatePostInput) (*po.Post, error)
	Update(postID, userID uint, isAdmin bool, input UpdatePostInput) (*po.Post, error)
	Delete(postID, userID uint, isAdmin bool) error
	GetByID(postID, viewerID uint) (*po.Post, bool, error)
	List(page, pageSize int, opts dto.PostListOptions) ([]po.Post, int64, error)
	// like
	Like(userID, postID uint) error
	Unlike(userID, postID uint) error
}

type postService struct {
	postRepo  postRepo.PostRepository
	tagRepo   tagRepo.TagRepository
	boardRepo boardRepo.BoardRepository
	userRepo  userRepo.UserRepository
	notifSvc  notification.NotificationService // 需导入 "tiny-forum/internal/service/notification"
	// riskSvc         *risk.RiskService
	contentcheckSvc check.ContentCheckService
}

func NewPostService(
	postRepo postRepo.PostRepository,
	tagRepo tagRepo.TagRepository,
	userRepo userRepo.UserRepository,
	boardRepo boardRepo.BoardRepository,
	notifSvc notification.NotificationService,
	// riskSvc *risk.RiskService,
	contentcheckSvc check.ContentCheckService,
) PostService {
	return &postService{
		postRepo:  postRepo,
		tagRepo:   tagRepo,
		userRepo:  userRepo,
		boardRepo: boardRepo,
		notifSvc:  notifSvc,
		// riskSvc:         riskSvc,
		contentcheckSvc: contentcheckSvc,
	}
}
