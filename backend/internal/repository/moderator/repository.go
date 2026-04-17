package moderator

import (
	"context"
	"tiny-forum/internal/model"

	"gorm.io/gorm"
)

// ModeratorRepository 版主数据访问接口
type ModeratorRepository interface {
	// 基础 CRUD
	Create(ctx context.Context, moderator *model.Moderator) error
	Update(ctx context.Context, moderator *model.Moderator) error
	Delete(ctx context.Context, id uint) error
	GetByID(ctx context.Context, id uint) (*model.Moderator, error)

	// 查询
	GetByUserAndBoard(ctx context.Context, userID, boardID uint) (*model.Moderator, error)
	GetByBoard(ctx context.Context, boardID uint) ([]model.Moderator, error)
	GetByUser(ctx context.Context, userID uint) ([]model.Moderator, error)

	// 权限相关
	UpdatePermissions(ctx context.Context, moderatorID uint, permissions model.Permission) error
	HasPermission(ctx context.Context, userID, boardID uint, permission string) (bool, error)

	// 批量操作
	DeleteByBoard(ctx context.Context, boardID uint) error
	DeleteByUser(ctx context.Context, userID uint) error

	// 检查
	Exists(ctx context.Context, userID, boardID uint) (bool, error)
	IsModerator(ctx context.Context, userID, boardID uint) (bool, error)

	// 分页列表
	List(ctx context.Context, page, pageSize int, boardID *uint) ([]model.Moderator, int64, error)
}

type moderatorRepository struct {
	db *gorm.DB
}

func NewModeratorRepository(db *gorm.DB) ModeratorRepository {
	return &moderatorRepository{db: db}
}
