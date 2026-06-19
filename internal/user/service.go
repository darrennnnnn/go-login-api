package user

import (
	"errors"
	"strconv"
	"time"

	"github.com/darrennnnnn/go-login-api/config"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
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
	user, err := s.repo.GetUser(requestBody.Username)

	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestBody.Password))
	if err != nil {
		return "", errors.New("Unauthorized")
	}

	idInt, _ := strconv.Atoi(user.Id)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       idInt,
		"username": user.UserName,
		"name":     user.Name,
		"exp":      time.Now().Add(1 * time.Hour).Unix(),
	})

	accessToken, err := token.SignedString(s.cfg.JWT.Secret)

	if err != nil {
		return "", err
	}

	return accessToken, nil

}