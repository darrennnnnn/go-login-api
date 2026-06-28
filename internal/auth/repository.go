package auth

import "gorm.io/gorm"

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateAccessToken(accessToken *AccessToken) error {
	return r.db.Create(accessToken).Error
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
