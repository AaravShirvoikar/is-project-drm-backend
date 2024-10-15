package repositories

import (
	"database/sql"

	"github.com/AaravShirvoikar/is-project-drm-backend/internal/models"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByEmail(username string) (*models.User, error)
}

type userRepo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(user *models.User) error {
	query := "INSERT INTO users (id, email, name, password) VALUES ($1, $2, $3, $4)"

	_, err := r.db.Exec(query, user.UserID, user.Email, user.UserName, user.Password)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepo) GetByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, name, email, password FROM users WHERE email = $1`

	err := r.db.QueryRow(query, email).Scan(&user.UserID, &user.UserName, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}

	return user, nil
}
