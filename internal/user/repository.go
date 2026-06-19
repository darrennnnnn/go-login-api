package user

import "gorm.io/gorm"

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetUser(username string) (*User, error) {
	user := &User{}
	err := r.db.Where("username = ?", username).First(user).Error
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
	return user, err
}

func (r *Repository) DeleteUser(id string) error {
	user := &User{}
	err := r.db.Delete(user, "id = ?", id).Error
	return err
}