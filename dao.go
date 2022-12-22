package authtokensvc

type Dao interface {
	DeleteAll(userId string) error
	Delete(refreshTokenId string) error
	SaveToken(token *Token) error
	GetRefreshToken(refreshTokenId string) (*Token, error)
	GetAccessToken(accessTokenId string) (*Token, error)
	HasToken(id string) (bool, error)
}
