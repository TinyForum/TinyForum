package tag

import (
	"tiny-forum/internal/model/do"
	tagRepo "tiny-forum/internal/repository/tag"
)

type TagService interface {
	Create(input CreateTagInput) (*do.Tag, error)
	List() ([]do.Tag, error)
	Get(id uint) (*do.Tag, error)
	Update(id uint, input CreateTagInput) (*do.Tag, error)
	Delete(id uint) error
}
type tagService struct {
	tagRepo tagRepo.TagRepository
}

func NewTagService(tagRepo tagRepo.TagRepository) TagService {
	return &tagService{tagRepo: tagRepo}
}
