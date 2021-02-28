package cacheredis

import (
	"time"

	"github.com/go-redis/redis"
)

type service struct {
	cache *redis.Client
}

func New(cache *redis.Client) *service {
	return &service{cache}
}

func (s *service) Get(key string) ([]byte, error) {
	return s.cache.Get(key).Bytes()
}

func (s *service) Del(key string) error {
	return s.cache.Del(key).Err()
}

func (s *service) Set(key string, data []byte, expires time.Duration) error {
	return s.cache.Set(key, data, expires).Err()
}
