package authtokensvc

import (
	"time"
)

type Token struct {
	ID        string
	ExpiresAt time.Time
	UserID    string
}
