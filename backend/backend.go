package backend

import (
	"errors"
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	"newspopper/loader"
	"time"
)

type Backend interface {
	Get(key string) (int, error)
	Set(key string) error
}

var defaultTTl time.Duration = time.Hour * 24 * 7

func NewBackend(loader loader.BackendConf) (Backend, error) {
	if loader.RedisConf.Uri != "" {
		persistDuration, err := time.ParseDuration(loader.PersistDuration)
		if err != nil {
			log.Infoln("error_while getting redis persist duration, fallback to default (1 week)")
			persistDuration = defaultTTl
		}
		return Redis{Client: redis.NewClient(&redis.Options{Addr: loader.RedisConf.Uri}), TTL: persistDuration}, nil
	}
	return nil, errors.New("no backend detected")
}
