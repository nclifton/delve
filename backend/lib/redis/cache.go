package redis

import (
	"errors"
	"reflect"
	"time"

	"github.com/go-redis/cache/v7"
	"github.com/vmihailenco/msgpack/v4"
)

func (db *Connection) EnableCache() error {

	if db.cache != nil {
		return errors.New("cache already enabled")
	}

	db.cache = &cache.Codec{
		Redis: db.Client,
		Marshal: func(v interface{}) ([]byte, error) {
			return msgpack.Marshal(v)
		},
		Unmarshal: func(b []byte, v interface{}) error {
			return msgpack.Unmarshal(b, v)
		},
	}

	return nil
}

type QueryFunc func() (interface{}, error)

func (db *Connection) Cached(key string, obj interface{}, expiry time.Duration, fn QueryFunc) (err error) {
	// if we don't have cache enabled we still need the wrap to work
	// but the reflecting below feels fragile \_o_/
	if db.cache == nil {
		res, err := fn()
		if err != nil {
			return err
		}

		reflect.ValueOf(obj).Elem().Set(reflect.ValueOf(res))
		return err
	}

	err = db.cache.Once(&cache.Item{
		Key:        key,
		Object:     obj,
		Expiration: expiry,
		Func:       fn,
	})

	return
}

func (db *Connection) PurgeCache(key string) (err error) {
	// if we don't have cache enabled we will make this
	// still work
	if db.cache == nil {
		return nil
	}

	err = db.cache.Delete(key)
	// we dont care if key didnt exist
	if err != nil && err != cache.ErrCacheMiss {
		return err
	}

	return nil
}
