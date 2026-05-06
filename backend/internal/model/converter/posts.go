package converter

import (
	"tiny-forum/internal/model/bo"
	"tiny-forum/internal/model/do"
)

func ListPostsBOToPostDO(b *bo.ListPosts) *do.Post {
	if b == nil {
		return nil
	}
	return &do.Post{
		// BaseModel: common.BaseModel{ // 通过嵌入类型名初始化

		// 	// UpdatedAt 和 DeletedAt 通常由数据库自动维护，不需要从 BO 传递
		// },
		AuthorID:         b.AuthorID,
		Type:             b.Type,
		PostStatus:       b.PostStatus,
		ModerationStatus: b.ModerationStatus,
	}
}
