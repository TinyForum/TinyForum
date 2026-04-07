package service

import (
	"bbs-forum/internal/model"
	"bbs-forum/internal/repository"
)

type TagService struct {
	tagRepo *repository.TagRepository
}

func NewTagService(tagRepo *repository.TagRepository) *TagService {
	return &TagService{tagRepo: tagRepo}
}

type CreateTagInput struct {
	Name        string `json:"name" binding:"required,min=1,max=50"`
	Description string `json:"description"`
	Color       string `json:"color"`
}

func (s *TagService) Create(input CreateTagInput) (*model.Tag, error) {
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

func (s *TagService) List() ([]model.Tag, error) {
	return s.tagRepo.List()
}

func (s *TagService) Update(id uint, input CreateTagInput) (*model.Tag, error) {
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

func (s *TagService) Delete(id uint) error {
	return s.tagRepo.Delete(id)
}
