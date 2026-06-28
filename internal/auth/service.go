package auth

import (
	"errors"
	"strings"
	"time"

	"github.com/darrennnnnn/go-login-api/config"
	"github.com/darrennnnnn/go-login-api/internal/user"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service struct {
	repo     *Repository
	userRepo *user.Repository
	cfg      *config.Config
}

func NewService(repo *Repository, userRepo *user.Repository, cfg *config.Config) *Service {
	return &Service{
		repo:     repo,
		userRepo: userRepo,
		cfg:      cfg,
	}
}

func (s *Service) Login(requestBody LoginRequest) (string, error) {
	email := strings.ToLower(strings.TrimSpace(requestBody.Email))

	userRecord, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("User not found")
		}
		return "", err
	}

	if userRecord == nil {
		return "", errors.New("User not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userRecord.Password), []byte(requestBody.Password)); err != nil {
		return "", errors.New("Unauthorized")
	}

	expirationDate := time.Now().Add(30 * time.Minute)
	accessTokenID := uuid.NewString()

	accessTokenRecord := &AccessToken{
		ID:        accessTokenID,
		UserID:    userRecord.Id,
		Revoked:   false,
		ExpiresAt: expirationDate,
	}

	if err := s.repo.CreateAccessToken(accessTokenRecord); err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       accessTokenID,
		"username": userRecord.Username,
		"email":    userRecord.Email,
		"exp":      expirationDate.Unix(),
	})

	accessToken, err := token.SignedString(s.cfg.JWT.Secret)
	if err != nil {
		return "", err
	}

	return accessToken, nil
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
