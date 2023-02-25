package authtokensvc

import (
	"github.com/muesli/cache2go"
	"time"
)

var _ AccessDao = (*accessCacheDao)(nil)

type accessCacheDao struct {
	cache              *cache2go.CacheTable
	userAccessTokenMap map[string][]string
}

func (a *accessCacheDao) Get(id string) (*Token, error) {
	res, err := a.cache.Value(id)

	if err != nil {
		return nil, nil
	}

	token := res.Data().(*Token)

	return token, nil
}

func (a *accessCacheDao) DeleteByUserID(id string) error {
	accessTokenIDs, ok := a.userAccessTokenMap[id]

	if !ok {
		return nil
	}

	for aId := range accessTokenIDs {
		a.cache.Delete(aId)
	}

	delete(a.userAccessTokenMap, id)

	return nil

}

func (a *accessCacheDao) DeleteByID(id string) error {
	a.cache.Delete(id)
	return nil
}

func (a *accessCacheDao) DeleteAll() error {
	a.cache.Flush()
	return nil
}

func (a *accessCacheDao) Create(token *Token) error {

	lifespan := token.ExpiresAt.Sub(time.Now())

	a.cache.Add(token.ID, lifespan, token)

	return nil
}

func (a *accessCacheDao) Has(id string) (bool, error) {

	_, err := a.cache.Value(id)

	if err != nil {
		return false, nil
	}

	return true, nil
}

func DefaultAccessDao() AccessDao {
	table := cache2go.Cache("accessTokens")

	return &accessCacheDao{cache: table, userAccessTokenMap: map[string][]string{}}
}
