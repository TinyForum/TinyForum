package article

import (
	articleService "tiny-forum/internal/service/article"
)

type ArticleHandler struct {
	articleSvc articleService.ArticleService
}

// 注意：原构造函数接收两个参数但只使用了第一个，这里改为只接收 postSvc
func NewPostHandler(postSvc articleService.ArticleService) *ArticleHandler {
	return &ArticleHandler{
		articleSvc: postSvc,
	}
}
