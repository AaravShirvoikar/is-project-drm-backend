package services

import (
	"crypto/rand"
	"database/sql"
	"errors"
	"time"

	"github.com/AaravShirvoikar/is-project-drm-backend/internal/models"
	"github.com/AaravShirvoikar/is-project-drm-backend/internal/repositories"
	"github.com/gofrs/uuid"
)

type SessionKeyService interface {
	GetOrCreate(userId, contentId string) ([]byte, error)
}

type sessionKeyService struct {
	sessionKeyRepo repositories.SessionKeyRepository
}

func NewSessionKeyService(sessionKeyRepo repositories.SessionKeyRepository) SessionKeyService {
	return &sessionKeyService{sessionKeyRepo: sessionKeyRepo}
}

func (s *sessionKeyService) GetOrCreate(userId, contentId string) ([]byte, error) {
	sessionKey, err := s.sessionKeyRepo.Get(userId, contentId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			key, err := generateSessionKey()
			if err != nil {
				return nil, err
			}

			keyId, err := uuid.NewV4()
			if err != nil {
				return nil, err
			}

			newSessionKey := &models.SessionKey{
				KeyID:      keyId,
				UserID:     uuid.FromStringOrNil(userId),
				ContentID:  uuid.FromStringOrNil(contentId),
				SessionKey: key,
				CreatedAt:  time.Now(),
				ExpiresAt:  time.Now().Add(24 * time.Hour),
			}

			err = s.sessionKeyRepo.Create(newSessionKey)
			if err != nil {
				return nil, err
			}
			return newSessionKey.SessionKey, nil
		}
		return nil, err
	}

	if sessionKey.ExpiresAt.Before(time.Now()) {
		s.sessionKeyRepo.Delete(sessionKey.KeyID.String())

		key, err := generateSessionKey()
		if err != nil {
			return nil, err
		}
		newSessionKey := &models.SessionKey{
			KeyID:      uuid.Must(uuid.NewV4()),
			UserID:     uuid.Must(uuid.FromString(userId)),
			ContentID:  uuid.Must(uuid.FromString(contentId)),
			SessionKey: key,
			CreatedAt:  time.Now(),
			ExpiresAt:  time.Now().Add(24 * time.Hour),
		}

		err = s.sessionKeyRepo.Create(newSessionKey)
		if err != nil {
			return nil, err
		}
		return newSessionKey.SessionKey, nil
	}

	return sessionKey.SessionKey, nil
}

func generateSessionKey() ([]byte, error) {
	key := make([]byte, 32)

	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}

	return key, nil
}
