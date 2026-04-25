package post

import (
	postService "tiny-forum/internal/service/post"
)

type PostHandler struct {
	postSvc postService.PostService
}

// 注意：原构造函数接收两个参数但只使用了第一个，这里改为只接收 postSvc
func NewPostHandler(postSvc postService.PostService) *PostHandler {
	return &PostHandler{
		postSvc: postSvc,
	}
}
