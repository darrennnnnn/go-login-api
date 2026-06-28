package user

import (
	"errors"
	"strings"
	"time"

	"github.com/darrennnnnn/go-login-api/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service struct {
	repo *Repository
	cfg *config.Config
}

func NewService(repo *Repository, cfg *config.Config) *Service {
	return &Service{
		repo: repo,
		cfg: cfg,
	}
}

func (s *Service) Login(requestBody LoginRequest) (string, error) {
	// normalize email
	email := strings.ToLower(strings.TrimSpace(requestBody.Email))

	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("User not found")
		}
		return "", err
	}

	if user == nil {
		return "", errors.New("User not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestBody.Password)); err != nil {
		return "", errors.New("Unauthorized")
	}

	expirationDate := time.Now().Add(30 * time.Minute)
	accessTokenID := uuid.NewString()

	accessTokenVar := &AccessToken{
		ID: accessTokenID,
		UserID: user.Id,
		Revoked: false,
		ExpiresAt: expirationDate,
	}

	err = s.repo.CreateAccessToken(accessTokenVar)

	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       accessTokenID,
		"username": user.Username,
		"email":    user.Email,
		"exp":      expirationDate.Unix(),
	})

	accessToken, err := token.SignedString(s.cfg.JWT.Secret)

	if err != nil {
		return "", err
	}

	return accessToken, nil

}

func (s *Service) CreateUser(requestBody CreateUserRequest) (*User, error) {
	// normalize input
	email := strings.ToLower(strings.TrimSpace(requestBody.Email))
	username := strings.TrimSpace(requestBody.Username)

	// check for existing email
	if existing, err := s.repo.GetUserByEmail(email); err == nil && existing != nil {
		return nil, errors.New("Email already in use")
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// check for existing username
	if existing, err := s.repo.GetUserByUsername(username); err == nil && existing != nil {
		return nil, errors.New("Username already in use")
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(requestBody.Password), 10)
	if err != nil {
		return nil, err
	}

	user := User{
		Id:       uuid.NewString(),
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
		return errors.New("User not found")
	}

	return err
}

func (s *Service) GetUsers() ([]User, error) {
	users, err := s.repo.GetUsers()

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *Service) GetUserByID(id string) (*User, error) {
	user, err := s.repo.GetUserByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("User not found")
		}
		return nil, err
	}
	return user, nil
}

func (s *Service) ValidateAccessToken(tokenID string) error {
	accessToken, err := s.repo.GetAccessTokenByID(tokenID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("Token not found, please try to login again")
		}
		return err
	}

	if accessToken.Revoked {
		return errors.New("Token is revoked")
	}

	if accessToken.ExpiresAt.Before(time.Now()) {
		return errors.New("Token is expired")
	}

	return nil
}

func (s *Service) Logout(tokenID string) error {
	return s.repo.RevokeAccessToken(tokenID)
}