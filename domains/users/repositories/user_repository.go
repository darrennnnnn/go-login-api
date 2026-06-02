package repositories

import (
	"errors"

	users "github.com/darrennnnnn/go-login-api/domains"
	"github.com/darrennnnnn/go-login-api/domains/users/entities"
	"golang.org/x/crypto/bcrypt"
)

type userRepository struct {
}

var userList []entities.User

func NewUserRepository() users.UserRepository {
	hashed1, _ := bcrypt.GenerateFromPassword([]byte("password1"), 10)
	userList = append(userList, entities.User{
		Id: "1",
		UserName: "user1",
		Name: "User 1",
		Password: string(hashed1),
	})

	hashed2, _ := bcrypt.GenerateFromPassword([]byte("password2"), 10)
	userList = append(userList, entities.User{
		Id: "2",
		UserName: "user2",
		Name: "User 2",
		Password: string(hashed2),
	})

	return userRepository{}
}


func (repo userRepository) FindUser(username string) (*entities.User, error) {
	for i := 0; i < len(userList); i++ {
		if userList[i].UserName == username {
			return &userList[i], nil
		}
	}

	return nil, errors.New("User not found")
}

func (repo userRepository) CreateUser(user *entities.User) (*entities.User, error) {
	return nil, errors.New("not implemented")
}

func (repo userRepository) UpdateUser(user *entities.User) (*entities.User, error) {
	return nil, errors.New("not implemented")
}

func (repo userRepository) DeleteUser(id string) {
	// TODO: implement delete
}