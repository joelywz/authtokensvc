package authtokensvc

import (
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

type Service struct {
	dao                   Dao
	accessExpireDuration  time.Duration
	refreshExpireDuration time.Duration
}

type Config struct {
	Dao                   Dao
	AccessExpireDuration  time.Duration
	RefreshExpireDuration time.Duration
}

func New(config *Config) *Service {
	return &Service{
		dao:                   config.Dao,
		accessExpireDuration:  config.AccessExpireDuration,
		refreshExpireDuration: config.RefreshExpireDuration,
	}
}

func (service *Service) Issue(userId string) (*IssueResponse, error) {

	// Generate token
	accessTokenId := gonanoid.Must(32)
	accessExpiry := time.Now().Add(service.accessExpireDuration)
	refreshTokenId := gonanoid.Must(32)
	refershExpiry := time.Now().Add(service.refreshExpireDuration)

	// Save token
	err := service.dao.SaveToken(&Token{
		RefreshID:        refreshTokenId,
		AccessID:         accessTokenId,
		RefreshExpiresAt: refershExpiry,
		AccessExpiresAt:  accessExpiry,
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
	token, err := service.dao.GetRefreshToken(refreshTokenId)

	if err != nil {
		return nil, err
	}

	if token == nil {
		return nil, ErrTokenNotFound
	}

	// Check if refresh token is expired
	if time.Now().After(token.RefreshExpiresAt) {
		return nil, ErrTokenExpired
	}

	// Delete token
	if err := service.dao.Delete(refreshTokenId); err != nil {
		return nil, err
	}

	// Issue new token
	return service.Issue(token.UserID)
}

func (service *Service) Verify(accessTokenId string) (bool, error) {

	// Obtain token by access token id from database
	token, err := service.dao.GetAccessToken(accessTokenId)

	if err != nil {
		return false, err
	}

	// Check if access token is expired
	if time.Now().After(token.AccessExpiresAt) {
		return false, ErrTokenExpired
	}

	return true, nil
}

func (service *Service) Revoke(accessTokenId string) error {

	// Obtain token by access token id from database
	token, err := service.dao.GetAccessToken(accessTokenId)

	if err != nil {
		return err
	}

	if token == nil {
		return ErrTokenNotFound
	}

	// Delete token
	return service.dao.Delete(token.RefreshID)

}

func (service *Service) RevokeAll(accessTokenId string) error {

	// Obtain token by access token id from database
	token, err := service.dao.GetAccessToken(accessTokenId)

	if err != nil {
		return err
	}

	if token == nil {
		return ErrTokenNotFound
	}

	// Delete tokens
	return service.dao.DeleteAll(token.UserID)

}

func (service *Service) AccessExpireDuration() time.Duration {
	return service.accessExpireDuration
}

func (service *Service) RefreshExpireDuration() time.Duration {
	return service.refreshExpireDuration
}
