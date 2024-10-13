package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type Content struct {
    ID          uuid.UUID `json:"id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    CreatorID   uuid.UUID `json:"creator_id"`
    Price       float64   `json:"price"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
