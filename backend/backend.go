package backend

import (
	"errors"
	"github.com/go-redis/redis"
	"newspopper/loader"
)

type Backend interface {
	Get(key string) (int, error)
	Set(key string) error
}

func NewBackend(loader loader.BackendConf) (Backend, error) {
	if loader.RedisConf.Uri != "" {
		return Redis{Client: redis.NewClient(&redis.Options{Addr: loader.RedisConf.Uri})}, nil
	}
	return nil, errors.New("no backend detected")
}
