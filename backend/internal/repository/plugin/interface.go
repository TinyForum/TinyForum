// package plugin

// import (
// 	"context"
// 	"tiny-forum/internal/model/common"
// 	"tiny-forum/internal/model/do"
// 	"tiny-forum/internal/model/dto"

// 	"gorm.io/gorm"
// )

// // PostRepository 帖子数据访问接口
// type PluginRepository interface {
// 	// 基础 CRUD
// 	// List(ctx context.Context, page, pageSize int, opts dto.PluginListOptionsQuery) ([]dto.PluginList, int64, error)
// 	List(ctx context.Context, query *dto.PluginQueryDTO, pageParam common.PageParam) (*common.PageResult[do.PluginMeta], error)
// 	GetByName(ctx context.Context, name string) (*do.PluginMeta, error)
//     Create(ctx context.Context, plugin *do.PluginMeta) error
// 	 ListByAuthorID(ctx context.Context, authorID int64) ([]*do.PluginMeta, error)
// }

// type pluginRepository struct {
// 	db *gorm.DB
// }

//	func NewPostRepository(db *gorm.DB) PluginRepository {
//		return &pluginRepository{
//			db: db,
//		}
//	}
package plugin