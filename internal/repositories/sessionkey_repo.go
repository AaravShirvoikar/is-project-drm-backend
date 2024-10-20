package repositories

import (
	"database/sql"
	"errors"

	"github.com/AaravShirvoikar/is-project-drm-backend/internal/models"
)

type SessionKeyRepository interface {
	Create(sessionKey *models.SessionKey) error
	Get(userId, contentId string) (*models.SessionKey, error)
	Delete(keyId string) error
}

type sessionKeyRepo struct {
	db *sql.DB
}

func NewSessionKeyRepo(db *sql.DB) SessionKeyRepository {
	return &sessionKeyRepo{db: db}
}

var ErrSessionKeyNotFound = errors.New("session key not found")

func (r *sessionKeyRepo) Create(sessionKey *models.SessionKey) error {
	query := `INSERT INTO session_keys (id, user_id, content_id, key, expires_at, created_at)
			VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.db.Exec(query, sessionKey.KeyID, sessionKey.UserID, sessionKey.ContentID, sessionKey.SessionKey, sessionKey.ExpiresAt, sessionKey.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *sessionKeyRepo) Get(userId, contentId string) (*models.SessionKey, error) {
	query := `SELECT id, user_id, content_id, key, expires_at, created_at FROM session_keys 
			WHERE user_id = $1 AND content_id = $2`

	row := r.db.QueryRow(query, userId, contentId)

	var sessionKey models.SessionKey
	err := row.Scan(&sessionKey.KeyID, &sessionKey.UserID, &sessionKey.ContentID, &sessionKey.SessionKey, &sessionKey.ExpiresAt, &sessionKey.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrSessionKeyNotFound
		}
		return nil, err
	}

	return &sessionKey, nil
}

func (r *sessionKeyRepo) Delete(keyId string) error {
	query := "DELETE FROM session_keys WHERE id = $1"

	_, err := r.db.Exec(query, keyId)
	if err != nil {
		return err
	}

	return nil
}
