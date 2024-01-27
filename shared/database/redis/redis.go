package redis

import (
	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
)

var conn *Redis

func GetConnection() *Redis {
	if conn != nil {
		return conn
	}

	panic("redis cache is not initialized")
}

func Connect(host, port string, db int, password string) *Redis {
	redisConn := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password,
		DB:       db,
	})

	pool := goredis.NewPool(redisConn)
	rs := redsync.New(pool)
	mutex := rs.NewMutex("redismutex")

	conn = &Redis{
		mutex:  mutex,
		client: redisConn,
	}

	return conn
}
