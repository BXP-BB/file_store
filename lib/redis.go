package lib

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

var RedisPool *redis.Pool

func init() {
	config := LoadServerConfig()
	// config := Config

	RedisPool = &redis.Pool{
		MaxIdle:     5,
		MaxActive:   0,
		Wait:        true,
		IdleTimeout: 200 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", config.RedisHost)
			if err != nil {
				return nil, err
			}
			if _, err := c.Do("SELECT", config.RedisIndex); err != nil {
				_ = c.Close()
				return nil, err
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, lastUsed time.Time) error {
			if time.Since(lastUsed) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
	fmt.Println("Redis init on port ", config.RedisHost)
}

func GetKey(key string) (string, error) {
	rds := RedisPool.Get()
	defer rds.Close()
	return redis.String(rds.Do("GET", key))
}

func SetKey(Key, Value interface{}, expires int) error {
	rds := RedisPool.Get()
	defer rds.Close()

	if expires == 0 {
		_, err := rds.Do("SET", Key, Value)
		return err
	} else {
		_, err := rds.Do("SETEX", Key, expires, Value)
		return err
	}
}

func DelKey(key string) error {
	rds := RedisPool.Get()
	defer rds.Close()
	_, err := rds.Do("DEL", key)
	return err
}

func LRange(key string, start, stop int64) ([]string, error) {
	rds := RedisPool.Get()
	defer rds.Close()
	return redis.Strings(rds.Do("LRANGE", key, start, stop))
}

func LPop(key string) (string, error) {
	rds := RedisPool.Get()
	defer rds.Close()
	return redis.String(rds.Do("LPOP", key))
}

func LPushAndTrimKey(key, value interface{}, size int64) error {
	rds := RedisPool.Get()
	defer rds.Close()
	_ = rds.Send("MULTI")
	_ = rds.Send("LPUSH", key, value)
	_ = rds.Send("LTRIM", key, size-2*size, -1)
	_, err := rds.Do("EXEC")
	return err
}

func RPushAndTrimKey(key, value interface{}, size int64) error {
	rds := RedisPool.Get()
	defer rds.Close()
	_ = rds.Send("MULTI")
	_ = rds.Send("RPUSH", key, value)
	_ = rds.Send("LTRIM", key, size-2*size, -1)
	_, err := rds.Do("EXEC")
	return err
}

func ExistsKey(key string) (bool, error) {
	rds := RedisPool.Get()
	defer rds.Close()
	return redis.Bool(rds.Do("EXISTS", key))
}

func TTLKey(key string) (int64, error) {
	rds := RedisPool.Get()
	defer rds.Close()
	return redis.Int64(rds.Do("TTL", key))
}

func Decr(key string) (int64, error) {
	rds := RedisPool.Get()
	defer rds.Close()
	return redis.Int64(rds.Do("DECR", key))
}

func MsetKey(key_value ...interface{}) error {
	rds := RedisPool.Get()
	defer rds.Close()
	_, err := rds.Do("MSET", key_value...)
	return err
}

func MgetKey(keys ...interface{}) map[interface{}]string {
	rds := RedisPool.Get()
	defer rds.Close()
	values, _ := redis.Strings(rds.Do("MGET", keys...))
	resultMap := map[interface{}]string{}
	keyLen := len(keys)
	for i := 0; i < keyLen; i++ {
		resultMap[keys[i]] = values[i]
	}
	return resultMap
}
