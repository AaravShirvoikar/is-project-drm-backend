package repositories

import (
	"database/sql"

	"github.com/AaravShirvoikar/is-project-drm-backend/internal/models"
	"github.com/gofrs/uuid"
)

type LicenseRepository interface {
	Create(license *models.License) error
	Get(userId, contentId string) (*models.License, error)
	Delete(licenseId string) error
}

type licenseRepo struct {
	db *sql.DB
}

func NewLicenseRepository(db *sql.DB) LicenseRepository {
	return &licenseRepo{db: db}
}

func (r *licenseRepo) Create(license *models.License) error {
	query := `INSERT INTO licenses (user_id, content_id, expires_at, created_at)
			VALUES ($1, $2, $3, $4)	RETURNING id`

	var id uuid.UUID
	err := r.db.QueryRow(query, license.UserID, license.ContentID, license.ExpiresAt, license.CreatedAt).Scan(&id)
	if err != nil {
		return err
	}

	license.LicenseID = id

	return nil
}

func (r *licenseRepo) Get(userId, contentId string) (*models.License, error) {
	query := `SELECT id, user_id, content_id, expires_at, created_at FROM licenses 
			WHERE user_id = $1 AND content_id = $2`

	row := r.db.QueryRow(query, userId, contentId)

	var license models.License
	err := row.Scan(&license.LicenseID, &license.UserID, &license.ContentID, &license.ExpiresAt, &license.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &license, nil
}

func (r *licenseRepo) Delete(licenseId string) error {
	query := "DELETE FROM licenses WHERE id = $1"

	_, err := r.db.Exec(query, licenseId)
	if err != nil {
		return err
	}

	return nil
}
