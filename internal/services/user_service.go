package services

import (
	"errors"

	"github.com/AaravShirvoikar/is-project-drm-backend/internal/models"
	"github.com/AaravShirvoikar/is-project-drm-backend/internal/repositories"
	"github.com/AaravShirvoikar/is-project-drm-backend/pkg/auth"
	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(user *models.User) error
	Authenticate(email, password string) (string, error)
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) Register(user *models.User) error {
	userId, err := uuid.NewV4()
	if err != nil {
		return err
	}
	user.UserID = userId

	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	return s.userRepo.Create(user)
}

func (s *userService) Authenticate(email, password string) (string, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return "", err
	}

	if !verifyPassword(password, user.Password) {
		return "", errors.New("invalid credentials")
	}

	return auth.GenerateJWT(user.UserID.String())
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func verifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
