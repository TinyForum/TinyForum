package tag

import (
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/request"
	tagRepo "tiny-forum/internal/repository/tag"
)

type TagService interface {
	Create(input request.CreateTagRequest) (*do.Tag, error)
	List() ([]do.Tag, error)
	Get(id uint) (*do.Tag, error)
	Update(id uint, input request.CreateTagRequest) (*do.Tag, error)
	Delete(id uint) error
}
type tagService struct {
	tagRepo tagRepo.TagRepository
}

func NewTagService(tagRepo tagRepo.TagRepository) TagService {
	return &tagService{tagRepo: tagRepo}
}
