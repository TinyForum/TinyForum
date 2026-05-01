package tag

import (
	"tiny-forum/internal/model/po"
	tagRepo "tiny-forum/internal/repository/tag"
)

type TagService interface {
	Create(input CreateTagInput) (*po.Tag, error)
	List() ([]po.Tag, error)
	Get(id uint) (*po.Tag, error)
	Update(id uint, input CreateTagInput) (*po.Tag, error)
	Delete(id uint) error
}
type tagService struct {
	tagRepo tagRepo.TagRepository
}

func NewTagService(tagRepo tagRepo.TagRepository) TagService {
	return &tagService{tagRepo: tagRepo}
}
