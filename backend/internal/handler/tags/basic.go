package tag

import (
	tagService "tiny-forum/internal/service/tag"
)

type TagHandler struct {
	tagSvc tagService.TagService
}

func NewTagHandler(tagSvc tagService.TagService) *TagHandler {
	return &TagHandler{tagSvc: tagSvc}
}
