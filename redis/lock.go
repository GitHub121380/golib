package redis

import (
	"github.com/GitHub121380/golib/zlog"
	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
)

const (
	EXSECONDS       = "EX"
	PXMILLISSECONDS = "PX"
	NOTEXISTS       = "NX"
)

// 设置过期时间为秒级的redis分布式锁
func (objRedis Redis) SetNxByEX(key string, value interface{}, expire uint64) (bool, error) {
	return objRedis.tryLock(key, value, expire, EXSECONDS)
}

// 设置过期时间为毫秒的redis分布式锁
func (objRedis Redis) SetNxByPX(key string, value interface{}, expire uint64) (bool, error) {
	return objRedis.tryLock(key, value, expire, PXMILLISSECONDS)
}

func (objRedis Redis) tryLock(key string, value interface{}, expire uint64, exType string) (bool, error) {
	str := parseToString(value)
	if str == "" {
		zlog.Warn(nil, "tryLock: [parase value is empty] detail: ", key, value)
		return false, errors.New("value is empty")
	}

	_, err := redis.String(objRedis.Do("SET", key, str, exType, expire, NOTEXISTS))

	if err == redis.ErrNil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
