package services

import (
	"sync"

	"github.com/attributeerror/steam-api-service/configuration"
	"github.com/attributeerror/steam-api-service/handlers/response_models"
	"github.com/patrickmn/go-cache"
)

type (
	allCache struct {
		steamHexs *cache.Cache
	}
)

var (
	cacheOnce     sync.Once
	steamHexCache *allCache
)

func GetCache() *allCache {
	cacheOnce.Do(func() {
		configuration := configuration.GetConfiguration()
		cacheStore := cache.New(configuration.CachingExpirationDuration, configuration.CachingPurgeDuration)

		steamHexCache = &allCache{
			steamHexs: cacheStore,
		}
	})

	return steamHexCache
}

func (c *allCache) Get(query string) (item *response_models.GetUserHexResponse, ok bool) {
	steamHexResponse, ok := c.steamHexs.Get(query)
	if ok {
		res, ok := steamHexResponse.(*response_models.GetUserHexResponse)
		if ok {
			return res, true
		}
	}
	return nil, false
}

func (c *allCache) Update(query string, item *response_models.GetUserHexResponse) {
	c.steamHexs.Set(query, item, cache.DefaultExpiration)
}
