package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Repository struct {
	db    *gorm.DB
	cache *redis.Client
}

func NewRepository(db *gorm.DB, cache *redis.Client) *Repository {
	return &Repository{db: db, cache: cache}
}

func (r *Repository) CreateAccessToken(accessToken *AccessToken) error {
	if err := r.db.Create(accessToken).Error; err != nil {
		return err
	}

	if err := r.cacheSetAccessToken(accessToken); err != nil {
		return nil
	}

	return nil
}

func (r *Repository) GetAccessTokenByID(tokenID string) (*AccessToken, error) {
	if accessToken, err := r.getAccessTokenFromCache(tokenID); err == nil && accessToken != nil {
		return accessToken, nil
	}

	accessToken := &AccessToken{}
	err := r.db.Where("id = ?", tokenID).First(accessToken).Error
	if err != nil {
		return nil, err
	}

	if err := r.cacheSetAccessToken(accessToken); err != nil {
		return accessToken, nil
	}

	return accessToken, nil
}

func (r *Repository) RevokeAccessToken(tokenID string) error {
	if err := r.db.Model(&AccessToken{}).Where("id = ?", tokenID).Update("revoked", true).Error; err != nil {
		return err
	}

	_ = r.cache.Del(context.Background(), accessTokenCacheKey(tokenID)).Err()

	return nil
}

func (r *Repository) getAccessTokenFromCache(tokenID string) (*AccessToken, error) {
	data, err := r.cache.Get(context.Background(), accessTokenCacheKey(tokenID)).Bytes()
	if err != nil {
		return nil, err
	}

	accessToken := &AccessToken{}
	if err := json.Unmarshal(data, accessToken); err != nil {
		return nil, err
	}

	return accessToken, nil
}

func (r *Repository) cacheSetAccessToken(accessToken *AccessToken) error {
	data, err := json.Marshal(accessToken)
	if err != nil {
		return err
	}

	ttl := time.Until(accessToken.ExpiresAt)
	if ttl <= 0 {
		ttl = time.Second
	}

	return r.cache.Set(context.Background(), accessTokenCacheKey(accessToken.ID), data, ttl).Err()
}

func accessTokenCacheKey(tokenID string) string {
	return fmt.Sprintf("auth:access_token:%s", tokenID)
}
