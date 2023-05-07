package redis

import (
	"github.com/gomodule/redigo/redis"
)

func (objRedis *Redis) Expire(key string, time int64) (bool, error) {
	return redis.Bool(objRedis.Do("EXPIRE", key, time))
}

func (objRedis *Redis) Exists(key string) (bool, error) {
	return redis.Bool(objRedis.Do("EXISTS", key))
}

func (objRedis *Redis) Del(keys ...interface{}) (int64, error) {
	return redis.Int64(objRedis.Do("DEL", keys...))
}

func (objRedis *Redis) Ttl(key string) (int64, error) {
	return redis.Int64(objRedis.Do("TTL", key))
}

func (objRedis *Redis) Pttl(key string) (int64, error) {
	return redis.Int64(objRedis.Do("PTTL", key))
}
