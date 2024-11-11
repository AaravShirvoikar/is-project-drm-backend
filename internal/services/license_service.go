package services

import (
	"database/sql"
	"errors"
	"time"

	"github.com/AaravShirvoikar/is-project-drm-backend/internal/models"
	"github.com/AaravShirvoikar/is-project-drm-backend/internal/repositories"
	"github.com/gofrs/uuid"
)

type LicenseService interface {
	Generate(userId, contentId string, expiresAt time.Time) error
	Verify(userId, contentId string) bool
	Revoke(licenseId string) error
}

type licenseService struct {
	licenseRepo repositories.LicenseRepository
}

func NewLicenseService(licenseRepo repositories.LicenseRepository) LicenseService {
	return &licenseService{licenseRepo: licenseRepo}
}

func (s *licenseService) Generate(userId, contentId string, expiresAt time.Time) error {
	licenseId, err := uuid.NewV4()
	if err != nil {
		return err
	}
	
	license := &models.License{
		LicenseID: licenseId,
		UserID:    uuid.FromStringOrNil(userId),
		ContentID: uuid.FromStringOrNil(contentId),
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}

	err = s.licenseRepo.Create(license)
	if err != nil {
		return err
	}

	return nil
}

func (s *licenseService) Verify(userId, contentId string) bool {
	license, err := s.licenseRepo.Get(userId, contentId)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return false
		}
	}

	if license == nil {
		return false
	}

	if license.ExpiresAt.Before(time.Now()) {
		return false
	}

	return true
}

func (s *licenseService) Revoke(licenseId string) error {
	err := s.licenseRepo.Delete(licenseId)
	if err != nil {
		return err
	}

	return nil
}
