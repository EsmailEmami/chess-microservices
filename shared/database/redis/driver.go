package redis

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
)

type Redis struct {
	client *redis.Client
	mutex  *redsync.Mutex
}

func (driver *Redis) Set(key string, value interface{}, expiration time.Duration) error {
	switch value.(type) {
	case string:
		{
			st := driver.client.Set(context.Background(), key, value, expiration)
			if st.Err() != nil {
				return st.Err()
			}
			return nil
		}
	default:
		{
			bts, err := json.Marshal(&value)
			if err != nil {
				return err
			}
			st := driver.client.Set(context.Background(), key, string(bts), expiration)
			return st.Err()

		}
	}
}

func (driver *Redis) Get(key string) (str string, err error) {
	value := driver.client.Get(context.Background(), key)
	if value == nil {
		err := errors.New("ErrCacheRecordNotFound")
		return "", err
	}
	if value.Err() != nil {
		return "", value.Err()
	}
	return value.Val(), nil
}

func (driver *Redis) UnmarshalToObject(key string, object interface{}) error {
	value, err := driver.Get(key)
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(value), object)
}

func (driver *Redis) Delete(key string) error {
	res := driver.client.Del(context.Background(), key)
	if res.Err() != nil {
		return res.Err()
	}
	d, err := res.Result()
	if err != nil {
		return err
	}
	if d == 0 {
		return err
	}
	return nil
}

func (driver *Redis) DeleteByPattern(key string) (deletedCount int64, err error) {
	scan := driver.client.Scan(context.Background(), 0, key+"*", 0)
	res, _, err := scan.Result()
	if err != nil {
		return 0, err
	}
	driver.client.Del(context.Background(), res...)

	return int64(len(res)), nil
}

func (driver *Redis) ResetDB() error {
	status := driver.client.FlushDB(context.Background())
	return status.Err()
}

func (driver *Redis) Lock() error {
	return driver.mutex.Lock()
}
func (driver *Redis) Unlock() error {
	_, err := driver.mutex.Unlock()
	return err
}
