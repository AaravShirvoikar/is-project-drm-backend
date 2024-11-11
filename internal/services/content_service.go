package services

import (
	"context"
	"errors"
	"io"

	"github.com/AaravShirvoikar/is-project-drm-backend/internal/models"
	"github.com/AaravShirvoikar/is-project-drm-backend/internal/repositories"
	"github.com/AaravShirvoikar/is-project-drm-backend/pkg/storage"
	"github.com/gofrs/uuid"
)

type ContentService interface {
	Create(content *models.Content, file io.Reader, fileExt string, fileSize int64) error
	Get(id string) (*models.Content, []byte, error)
	List() ([]*models.Content, error)
}

type contentService struct {
	contentRepo repositories.ContentRepository
	storage     *storage.FileStorage
}

func NewContentService(contentRepo repositories.ContentRepository, storage *storage.FileStorage) ContentService {
	return &contentService{contentRepo: contentRepo, storage: storage}
}

func (s *contentService) Create(content *models.Content, file io.Reader, fileExt string, fileSize int64) error {
	if content.Title == "" {
		return errors.New("content title cannot be empty")
	}
	if content.Price < 0 {
		return errors.New("content price cannot be negative")
	}

	fileId, err := s.storage.UploadFile(context.Background(), file, fileExt, fileSize)
	if err != nil {
		return err
	}

	contentId, err := uuid.NewV4()
	if err != nil {
		return err
	}
	content.ContentID = contentId

	content.FileID = fileId
	content.FileSize = fileSize

	err = s.contentRepo.Create(content)
	if err != nil {
		return err
	}

	return nil
}

func (s *contentService) List() ([]*models.Content, error) {
	return s.contentRepo.GetAll()
}

func (s *contentService) Get(id string) (*models.Content, []byte, error) {
    content, err := s.contentRepo.GetById(id)
    if err != nil {
        return nil, nil, err
    }

    file, err := s.storage.DownloadFile(context.Background(), content.FileID)
    if err != nil {
        return nil, nil, err
    }
    defer file.Close()

    fileContent, err := io.ReadAll(file)
    if err != nil {
        return nil, nil, err
    }

    return content, fileContent, nil
}