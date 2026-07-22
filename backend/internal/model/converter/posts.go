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
		AuthorID:         bo.AuthorID,
		PostStatus:       bo.PostStatus,
		ModerationStatus: bo.ModerationStatus,
		Type:             do.PostType(bo.Type), // Type 为 string，需转换
		// 不包含 Keyword, TagNames, SortBy
	}
}
