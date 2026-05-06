package admin

import (
	"context"
	"tiny-forum/internal/model/bo"
	"tiny-forum/internal/model/common"
	"tiny-forum/internal/model/do"
	"tiny-forum/pkg/logger"
)

func (s *adminService) ListPosts(ctx context.Context, listPostsBO *common.PageQuery[bo.ListPosts]) ([]do.Post, int64, error) {
	logger.Infof("查询参数：", listPostsBO)
	return s.postSvc.AdminList(ctx, listPostsBO)
}
