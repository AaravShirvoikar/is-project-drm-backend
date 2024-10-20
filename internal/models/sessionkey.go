package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type SessionKey struct {
	KeyID      uuid.UUID `json:"key_id"`
	UserID     uuid.UUID `json:"user_id"`
	ContentID  uuid.UUID `json:"content_id"`
	SessionKey []byte    `json:"session_key"`
	ExpiresAt  time.Time `json:"expires_at"`
	CreatedAt  time.Time `json:"created_at"`
}
