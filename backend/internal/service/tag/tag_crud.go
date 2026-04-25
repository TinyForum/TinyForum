package tag

import (
	"tiny-forum/internal/model"
)

type CreateTagInput struct {
	Name        string `json:"name" binding:"required,min=1,max=50"`
	Description string `json:"description"`
	Color       string `json:"color"`
}

func (s *tagService) Create(input CreateTagInput) (*model.Tag, error) {
	color := input.Color
	if color == "" {
		color = "#6366f1"
	}
	tag := &model.Tag{
		Name:        input.Name,
		Description: input.Description,
		Color:       color,
	}
	if err := s.tagRepo.Create(tag); err != nil {
		return nil, err
	}
	return tag, nil
}

func (s *tagService) List() ([]model.Tag, error) {
	return s.tagRepo.List()
}

func (s *tagService) Update(id uint, input CreateTagInput) (*model.Tag, error) {
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
