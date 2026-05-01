package moderator

import (
	"context"
	"tiny-forum/internal/model/po"

	"gorm.io/gorm"
)

// ModeratorRepository 版主数据访问接口
type ModeratorRepository interface {
	// 基础 CRUD
	Create(ctx context.Context, moderator *po.Moderator) error
	Update(ctx context.Context, moderator *po.Moderator) error
	Delete(ctx context.Context, id uint) error
	GetByID(ctx context.Context, id uint) (*po.Moderator, error)

	// 查询
	GetByUserAndBoard(ctx context.Context, userID, boardID uint) (*po.Moderator, error)
	GetByBoard(ctx context.Context, boardID uint) ([]po.Moderator, error)
	GetByUser(ctx context.Context, userID uint) ([]po.Moderator, error)

	// 权限相关
	UpdatePermissions(ctx context.Context, moderatorID uint, permissions po.Permission) error
	HasPermission(ctx context.Context, userID, boardID uint, permission string) (bool, error)

	// 批量操作
	DeleteByBoard(ctx context.Context, boardID uint) error
	DeleteByUser(ctx context.Context, userID uint) error

	// 检查
	Exists(ctx context.Context, userID, boardID uint) (bool, error)
	IsModerator(ctx context.Context, userID, boardID uint) (bool, error)

	// 分页列表
	List(ctx context.Context, page, pageSize int, boardID *uint) ([]po.Moderator, int64, error)
}

type moderatorRepository struct {
	db *gorm.DB
}

func NewModeratorRepository(db *gorm.DB) ModeratorRepository {
	return &moderatorRepository{db: db}
}
