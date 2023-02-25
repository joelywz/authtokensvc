package authtokensvc

import (
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

type Service struct {
	accessDao             AccessDao
	refreshDao            RefreshDao
	accessExpireDuration  time.Duration
	refreshExpireDuration time.Duration
}

type Config struct {
	AccessDao             AccessDao
	RefreshDao            RefreshDao
	AccessExpireDuration  time.Duration
	RefreshExpireDuration time.Duration
}

func New(config *Config) *Service {
	return &Service{
		accessDao:             config.AccessDao,
		refreshDao:            config.RefreshDao,
		accessExpireDuration:  config.AccessExpireDuration,
		refreshExpireDuration: config.RefreshExpireDuration,
	}
}

func (service *Service) Issue(userId string) (*IssueResponse, error) {

	// Generate token
	accessTokenId := gonanoid.Must(32)
	accessExpiry := time.Now().Add(service.accessExpireDuration)
	refreshTokenId := gonanoid.Must(32)
	refreshExpiry := time.Now().Add(service.refreshExpireDuration)

	// Create refresh token
	err := service.refreshDao.Create(&Token{
		ID:        refreshTokenId,
		ExpiresAt: refreshExpiry,
		UserID:    userId,
	})

	if err != nil {
		return nil, err
	}

	// Create access token2
	err = service.accessDao.Create(&Token{
		ID:        accessTokenId,
		ExpiresAt: accessExpiry,
		UserID:    userId,
	})

	if err != nil {
		return nil, err
	}

	// Create response
	res := &IssueResponse{
		RefreshToken: refreshTokenId,
		AccessToken:  accessTokenId,
	}

	return res, nil
}

func (service *Service) Refresh(refreshTokenId string) (*RefreshResponse, error) {

	// Obtain token by refresh token id from database
	token, err := service.refreshDao.Get(refreshTokenId)

	if err != nil {
		return nil, err
	}

	if token == nil {
		return nil, ErrTokenNotFound
	}

	// Check if refresh token is expired
	if time.Now().After(token.ExpiresAt) {
		return nil, ErrTokenExpired
	}

	// Delete token
	if err := service.refreshDao.DeleteByID(refreshTokenId); err != nil {
		return nil, err
	}

	// Issue new token
	return service.Issue(token.UserID)
}

func (service *Service) Verify(accessTokenId string) (*Token, error) {

	// Obtain token by access token id from database
	token, err := service.accessDao.Get(accessTokenId)

	if err != nil {
		return nil, err
	}

	if token == nil {
		return nil, ErrTokenNotFound
	}

	// Check if access token is expired
	if time.Now().After(token.ExpiresAt) {
		return nil, ErrTokenExpired
	}

	return token, nil
}

func (service *Service) Revoke(refreshTokenId string) error {

	// Obtain token by access token id from database
	token, err := service.refreshDao.Get(refreshTokenId)

	if err != nil {
		return err
	}

	if token == nil {
		return ErrTokenNotFound
	}

	// Delete token
	return service.refreshDao.DeleteByID(refreshTokenId)

}

func (service *Service) RevokeAll(refreshTokenId string) error {

	// Obtain token by access token id from database
	token, err := service.refreshDao.Get(refreshTokenId)

	if err != nil {
		return err
	}

	if token == nil {
		return ErrTokenNotFound
	}

	// Delete refresh tokens
	if err := service.refreshDao.DeleteByUserID(token.UserID); err != nil {
		return err
	}

	// Delete access tokens
	if err := service.accessDao.DeleteByUserID(token.UserID); err != nil {
		return err
	}

	return nil
}

func (service *Service) AccessExpireDuration() time.Duration {
	return service.accessExpireDuration
}

func (service *Service) RefreshExpireDuration() time.Duration {
	return service.refreshExpireDuration
}
