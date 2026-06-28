package user

import (
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetUserByEmail(email string) (*User, error) {
	user := &User{}
	err := r.db.Where("email = ?", email).First(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *Repository) GetUserByUsername(username string) (*User, error) {
	user := &User{}
	err := r.db.Where("username = ?", username).First(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *Repository) GetUsers() ([]User, error) {
	var users []User
	err := r.db.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *Repository) GetUserByID(id string) (*User, error) {
	user := &User{}
	err := r.db.Where("id = ?", id).First(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *Repository) CreateUser(user *User) (*User, error) {
	err := r.db.Create(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *Repository) DeleteUser(id string) (int, error) {
	res := r.db.Where("id = ?", id).Delete(&User{})
	if res.Error != nil {
		return 0, res.Error
	}
	return int(res.RowsAffected), nil
}

func (r *Repository) CreateAccessToken(accessToken *AccessToken) error {
	err := r.db.Create(accessToken).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetAccessTokenByID(tokenID string) (*AccessToken, error) {
	accessToken := &AccessToken{}
	err := r.db.Where("id = ?", tokenID).First(accessToken).Error
	if err != nil {
		return nil, err
	}
	return accessToken, nil
}

func (r *Repository) RevokeAccessToken(tokenID string) error {
	return r.db.Model(&AccessToken{}).Where("id = ?", tokenID).Update("revoked", true).Error
}