package services

import (
	"errors"

	"github.com/AaravShirvoikar/is-project-drm-backend/internal/models"
	"github.com/AaravShirvoikar/is-project-drm-backend/internal/repositories"
)

type ContentService interface {
	Create(content *models.Content) error
}

type contentService struct {
	contentRepo repositories.ContentRepository
}

func NewContentService(contentRepo repositories.ContentRepository) ContentService {
	return &contentService{contentRepo: contentRepo}
}

func (s *contentService) Create(content *models.Content) error {
	if content.Title == "" {
		return errors.New("content title cannot be empty")
	}
	if content.Price < 0 {
		return errors.New("content price cannot be negative")
	}
	
	err := s.contentRepo.Create(content)
	if err != nil {
		return err
	}
	
	return nil
}
