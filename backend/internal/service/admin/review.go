package admin

import (
	"context"
	"tiny-forum/internal/model/bo"
	"tiny-forum/internal/model/common"
	"tiny-forum/internal/model/do"
)

func (s *adminService) ListReviewRequire(ctx context.Context, ListPostsBO *common.PageQuery[bo.ListPosts]) ([]do.Article, int64, error) {
	return s.postSvc.AdminLists(ctx, ListPostsBO)
}
