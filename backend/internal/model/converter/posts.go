package converter

import (
	"tiny-forum/internal/model/bo"
	"tiny-forum/internal/model/do"
)

func ListPostsBOToPostDO(bo *bo.ListPosts) *do.Article {
	if bo == nil {
		return &do.Article{}
	}
	return &do.Article{

		Creation: do.Creation{
			AuthorID:         bo.AuthorID,
			CreationStatus:   bo.PostStatus,
			ModerationStatus: bo.ModerationStatus,
			Type:             do.CreationType(bo.Type),
		},
	}
}
