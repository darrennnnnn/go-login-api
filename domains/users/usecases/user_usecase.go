package usecases

import (
	"errors"
	"strconv"
	"time"

	users "github.com/darrennnnnn/go-login-api/domains"
	"github.com/darrennnnnn/go-login-api/domains/users/models/requests"
	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v5"
)

type userUseCase struct {
	repo users.UserRepository
}

func NewUserUseCase(repo users.UserRepository) users.UserUseCase {
	return userUseCase{repo: repo}
}

func (u userUseCase) Login(request *requests.LoginRequest) (string, error) {
	user, err := u.repo.FindUser(request.Username)

	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		return "", errors.New("Unauthorized")
	}

	// simple secret (same as middleware)
	secret := []byte("change-me")

	idInt, _ := strconv.Atoi(user.Id)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       idInt,
		"username": user.UserName,
		"name":     user.Name,
		"exp":      time.Now().Add(1 * time.Hour).Unix(),
	})

	accessToken, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}