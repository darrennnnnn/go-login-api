package users

import (
	"github.com/darrennnnnn/go-login-api/domains/users/entities"
	"github.com/darrennnnnn/go-login-api/domains/users/models/requests"
)

type UserUseCase interface {
	Login(request *requests.LoginRequest) (string, error)
}

type UserRepository interface {
	FindUser(username string) (*entities.User, error)
	CreateUser(user *entities.User) (*entities.User, error)
	UpdateUser(user *entities.User) (*entities.User, error)
	DeleteUser(id string)
}