package authtokensvc

import "time"

type TokenType = string

var (
	TypeRefresh = "refresh"
	TypeAccess  = "access"
)

type Token struct {
	RefreshID        string    `bson:"refresh_token_id"`
	AccessID         string    `bson:"access_token_id"`
	RefreshExpiresAt time.Time `bson:"refresh_expires_at"`
	AccessExpiresAt  time.Time `bson:"access_expires_at"`
	UserID           string    `bson:"user_id"`
}
