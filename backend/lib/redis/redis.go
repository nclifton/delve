package redis

import (
	"errors"
	"log"

	"github.com/go-redis/cache/v7"
	"github.com/go-redis/redis/v7"
)

// aliases to remove redis dependency requirement elsewhere
type Connection struct {
	Client *redis.Client
	cache  *cache.Codec
}

const RedisNil = redis.Nil

// Connect helper
func Connect(url string) (*Connection, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     url,
		PoolSize: 5,
	})

	if _, err := client.Ping().Result(); err != nil {
		log.Println("fail to initialize redis client: ", err)
		return nil, errors.New("redis client initial ping failed")
	}

	return &Connection{Client: client}, nil
}
