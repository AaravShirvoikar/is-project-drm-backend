package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type License struct {
	LicenseID uuid.UUID `json:"license_id"`
	UserID    uuid.UUID `json:"user_id"`
	ContentID uuid.UUID `json:"content_id"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}
