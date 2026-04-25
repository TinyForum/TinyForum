package tag

import (
	tagRepo "tiny-forum/internal/repository/tag"
)

type TagService struct {
	tagRepo tagRepo.TagRepository
}

func NewTagService(tagRepo tagRepo.TagRepository) *TagService {
	return &TagService{tagRepo: tagRepo}
}
