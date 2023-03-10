package authtokensvc

type AccessDao interface {
	DeleteByUserID(id string) error
	DeleteByID(id string) error
	DeleteAll() error
	Get(id string) (*Token, error)
	Create(token *Token) error
	Has(id string) (bool, error)
}
