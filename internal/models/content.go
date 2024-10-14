package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type Content struct {
	ContentID   uuid.UUID `json:"content_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatorID   uuid.UUID `json:"creator_id"`
	Price       float64   `json:"price"`
	FileID      string    `json:"file_id"`
	FileSize    int64     `json:"file_size"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type File struct {
	FileID string `json:"file_id"`
	Size   int64  `json:"size"`
}
