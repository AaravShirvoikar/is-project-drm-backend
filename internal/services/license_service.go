package services

import (
	"errors"
	"time"

	"github.com/AaravShirvoikar/is-project-drm-backend/internal/models"
	"github.com/AaravShirvoikar/is-project-drm-backend/internal/repositories"
	"github.com/gofrs/uuid"
)

type LicenseService interface {
	Generate(userId, contentId string, expiresAt time.Time) (*models.License, error)
	Verify(userId, contentId string) (bool, error)
	Revoke(licenseId string) error
}

type licenseService struct {
	licenseRepo repositories.LicenseRepository
}

func NewLicenseService(licenseRepo repositories.LicenseRepository) LicenseService {
	return &licenseService{licenseRepo: licenseRepo}
}

func (s *licenseService) Generate(userId, contentId string, expiresAt time.Time) (*models.License, error) {
	license := &models.License{
		UserID:    uuid.FromStringOrNil(userId),
		ContentID: uuid.FromStringOrNil(contentId),
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}

	err := s.licenseRepo.Create(license)
	if err != nil {
		return nil, err
	}

	return license, nil
}

func (s *licenseService) Verify(userId, contentId string) (bool, error) {
	license, err := s.licenseRepo.Get(userId, contentId)
	if err != nil {
		return false, err
	}

	if license == nil {
		return false, errors.New("license not found")
	}

	if license.ExpiresAt.Before(time.Now()) {
		return false, errors.New("license has expired")
	}

	return true, nil
}

func (s *licenseService) Revoke(licenseId string) error {
	err := s.licenseRepo.Delete(licenseId)
	if err != nil {
		return err
	}

	return nil
}
