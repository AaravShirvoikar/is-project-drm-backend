package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"

	"github.com/AaravShirvoikar/is-project-drm-backend/internal/models"
	"github.com/AaravShirvoikar/is-project-drm-backend/internal/repositories"
	"github.com/AaravShirvoikar/is-project-drm-backend/pkg/storage"
	"github.com/gofrs/uuid"
)

type ContentService interface {
	Create(content *models.Content, file io.Reader, fileExt string, fileSize int64) (string, bool, float64, error)
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

func (s *contentService) Create(content *models.Content, file io.Reader, fileExt string, fileSize int64) (string, bool, float64, error) {
	if content.Title == "" {
		return "", false, 0, errors.New("content title cannot be empty")
	}
	if content.Price < 0 {
		return "", false, 0, errors.New("content price cannot be negative")
	}

	contentId, err := uuid.NewV4()
	if err != nil {
		return "", false, 0, err
	}
	content.ContentID = contentId

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", false, 0, err
	}

	url := "http://localhost:8000/compare-video-bytes/"
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	part1, err := writer.CreateFormFile("file", filepath.Base("file.mp4"))
	if err != nil {
		return "", false, 0, err
	}

	if _, err = part1.Write(fileBytes); err != nil {
		return "", false, 0, err
	}

	if err = writer.WriteField("file_id", contentId.String()); err != nil {
		return "", false, 0, err
	}

	if err = writer.Close(); err != nil {
		return "", false, 0, err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return "", false, 0, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		return "", false, 0, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", false, 0, err
	}
	var resp struct {
		VideoID       string  `json:"video_id"`
		MaxSimilarity float64 `json:"max_similarity"`
		Similar       bool    `json:"similar"`
	}

	if err = json.Unmarshal(body, &resp); err != nil {
		return "", false, 0, err
	}

	if resp.Similar {
		return resp.VideoID, false, resp.MaxSimilarity, nil
	}

	fileReader := bytes.NewReader(fileBytes)
	fileId, err := s.storage.UploadFile(context.Background(), fileReader, fileExt, fileSize)
	if err != nil {
		return "", false, 0, err
	}

	content.FileID = fileId
	content.FileSize = fileSize

	err = s.contentRepo.Create(content)
	if err != nil {
		return "", false, 0, err
	}

	return resp.VideoID, true, resp.MaxSimilarity, nil
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
