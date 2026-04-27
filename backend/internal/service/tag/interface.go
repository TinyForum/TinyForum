package tag

import (
	"tiny-forum/internal/model"
	tagRepo "tiny-forum/internal/repository/tag"
)

type TagService interface {
	Create(input CreateTagInput) (*model.Tag, error)
	List() ([]model.Tag, error)
	Get(id uint) (*model.Tag, error)
	Update(id uint, input CreateTagInput) (*model.Tag, error)
	Delete(id uint) error
}
type tagService struct {
	tagRepo tagRepo.TagRepository
}

func NewTagService(tagRepo tagRepo.TagRepository) TagService {
	return &tagService{tagRepo: tagRepo}
}
