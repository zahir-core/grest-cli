package app

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"grest.dev/grest"
)

func Cache() CacheInterface {
	if cache == nil {
		cache = &cacheImpl{}
		cache.configure()
	}
	return cache
}

type CacheInterface interface {
	grest.CacheInterface
	Invalidate(endPoint string, keys ...string)
}

var cache *cacheImpl

// cacheImpl implement CacheInterface embed from grest.Cache for simplicity
type cacheImpl struct {
	grest.Cache
}

func (c *cacheImpl) configure() {
	c.Exp = 24 * time.Hour
	c.RedisClient = redis.NewClient(&redis.Options{
		Addr:     REDIS_HOST + ":" + REDIS_PORT,
		Username: REDIS_USERNAME,
		Password: REDIS_PASSWORD,
		DB:       REDIS_CACHE_DB,
	})
	c.Ctx = context.Background()
	err := c.RedisClient.Ping(c.Ctx).Err()
	if err != nil {
		Logger().Error().
			Err(err).
			Str("REDIS_HOST", REDIS_HOST).
			Str("REDIS_PORT", REDIS_PORT).
			Str("REDIS_USERNAME", REDIS_USERNAME).
			Str("REDIS_PASSWORD", REDIS_PASSWORD).
			Int("REDIS_CACHE_DB", REDIS_CACHE_DB).
			Msg("Failed to connect to redis. The cache will be use in-memory local storage.")
	} else {
		c.IsUseRedis = true
		Logger().Info().Msg("Cache configured with redis.")
	}
}

func (c *cacheImpl) Invalidate(endPoint string, keys ...string) {
	for _, k := range keys {
		c.Delete(endPoint + "." + k)
	}
	go c.DeleteWithPrefix(endPoint + "?")
}
