package models

import "github.com/gofrs/uuid"

type User struct {
	UserID   uuid.UUID `json:"user_id"`
	Email    string    `json:"email"`
	UserName string    `json:"user_name"`
	Password string    `json:"password"`
}
