package user

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateUser(requestBody CreateUserRequest) (*User, error) {
	email := strings.ToLower(strings.TrimSpace(requestBody.Email))
	username := strings.TrimSpace(requestBody.Username)

	if existing, err := s.repo.GetUserByEmail(email); err == nil && existing != nil {
		return nil, ErrEmailInUse
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if existing, err := s.repo.GetUserByUsername(username); err == nil && existing != nil {
		return nil, ErrUsernameInUse
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(requestBody.Password), 10)
	if err != nil {
		return nil, err
	}

	user := User{
		ID:       uuid.NewString(),
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
	}

	result, err := s.repo.CreateUser(&user)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Service) DeleteUser(id string) error {
	rowsAffected, err := s.repo.DeleteUser(id)

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (s *Service) GetUsers() ([]User, error) {
	return s.repo.GetUsers()
}

func (s *Service) GetUserByID(id string) (*User, error) {
	user, err := s.repo.GetUserByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}
