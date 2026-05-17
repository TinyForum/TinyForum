package tag

import (
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/request"
)

func (s *tagService) Create(input request.CreateTagRequest) (*do.Tag, error) {
	color := input.Color
	if color == "" {
		color = "#6366f1"
	}
	tag := &do.Tag{
		Name:        input.Name,
		Description: input.Description,
		Color:       color,
	}
	if err := s.tagRepo.Create(tag); err != nil {
		return nil, err
	}
	return tag, nil
}

func (s *tagService) List() ([]do.Tag, error) {
	return s.tagRepo.List()
}
func (s *tagService) Get(id uint) (*do.Tag, error) {
	return s.tagRepo.FindByID(id)
}

func (s *tagService) Update(id uint, input request.CreateTagRequest) (*do.Tag, error) {
	tag, err := s.tagRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if input.Name != "" {
		tag.Name = input.Name
	}
	if input.Description != "" {
		tag.Description = input.Description
	}
	if input.Color != "" {
		tag.Color = input.Color
	}
	if err := s.tagRepo.Update(tag); err != nil {
		return nil, err
	}
	return tag, nil
}

func (s *tagService) Delete(id uint) error {
	return s.tagRepo.Delete(id)
}
