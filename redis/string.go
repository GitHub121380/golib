package redis

import (
	"math"
	"strings"

	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
)

const (
	_CHUNK_SIZE = 32
)

func (objRedis *Redis) Get(key string) ([]byte, error) {
	if res, err := redis.Bytes(objRedis.Do("GET", key)); err == redis.ErrNil {
		return nil, nil
	} else {
		return res, err
	}
}

func (objRedis *Redis) MGet(keys ...string) [][]byte {
	//1.初始化返回结果
	res := make([][]byte, 0, len(keys))

	//2.将多个key分批获取（每次32个）
	pageNum := int(math.Ceil(float64(len(keys)) / float64(_CHUNK_SIZE)))
	for n := 0; n < pageNum; n++ {
		//2.1创建分批切片 []string
		var end int
		if n != (pageNum - 1) {
			end = (n + 1) * _CHUNK_SIZE
		} else {
			end = len(keys)
		}
		chunk := keys[n*_CHUNK_SIZE : end]
		//2.2分批切片的类型转换 => []interface{}
		chunkLength := len(chunk)
		keyList := make([]interface{}, 0, chunkLength)
		for _, v := range chunk {
			keyList = append(keyList, v)
		}
		cacheRes, err := redis.ByteSlices(objRedis.Do("MGET", keyList...))
		if err != nil {
			for i := 0; i < len(keyList); i++ {
				res = append(res, nil)
			}
		} else {
			res = append(res, cacheRes...)
		}
	}
	return res
}

func (objRedis *Redis) MSet(values ...interface{}) error {
	_, err := objRedis.Do("MSET", values...)
	return err
}

func (objRedis *Redis) Set(key string, value interface{}, expire ...int64) error {
	var res string
	var err error
	if expire == nil {
		res, err = redis.String(objRedis.Do("SET", key, value))
	} else {
		res, err = redis.String(objRedis.Do("SET", key, value, "EX", expire[0]))
	}
	if err != nil {
		return err
	} else if strings.ToLower(res) != "ok" {
		return errors.New("set result not OK")
	}
	return nil
}

func (objRedis *Redis) SetEx(key string, value interface{}, expire int64) error {
	return objRedis.Set(key, value, expire)
}

func (objRedis *Redis) Append(key string, value interface{}) (int, error) {
	return redis.Int(objRedis.Do("APPEND", key, value))
}

func (objRedis *Redis) Incr(key string) (int64, error) {
	return redis.Int64(objRedis.Do("INCR", key))
}

func (objRedis *Redis) IncrBy(key string, value int64) (int64, error) {
	return redis.Int64(objRedis.Do("INCRBY", key, value))
}

func (objRedis *Redis) IncrByFloat(key string, value float64) (float64, error) {
	return redis.Float64(objRedis.Do("INCRBYFLOAT", key, value))
}

func (objRedis *Redis) Decr(key string) (int64, error) {
	return redis.Int64(objRedis.Do("DECR", key))
}

func (objRedis *Redis) DecrBy(key string, value int64) (int64, error) {
	return redis.Int64(objRedis.Do("DECRBY", key, value))
}
