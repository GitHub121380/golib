package redis

import (
	"errors"
	"github.com/GitHub121380/golib/zlog"
	"math"
	"strconv"

	"github.com/gomodule/redigo/redis"
)

func (objRedis *Redis) HSet(key, field string, val interface{}) (int, error) {
	valStr := parseToString(val)
	return redis.Int(objRedis.Do("HSET", key, field, valStr))
}

func (objRedis *Redis) HGet(key, field string) ([]byte, error) {
	if res, err := redis.Bytes(objRedis.Do("HGET", key, field)); err == redis.ErrNil {
		return nil, nil
	} else {
		return res, err
	}
}

func (objRedis *Redis) HMGet(key string, fields ...string) ([][]byte, error) {
	//1.初始化返回结果
	res := make([][]byte, 0, len(fields))
	var resErr error
	//2.将多个key分批获取（每次32个）
	pageNum := int(math.Ceil(float64(len(fields)) / float64(_CHUNK_SIZE)))
	for i := 0; i < pageNum; i++ {
		//2.1创建分批切片 []string
		var end int
		if i == (pageNum - 1) {
			end = len(fields)
		} else {
			end = (i + 1) * _CHUNK_SIZE
		}
		chunk := fields[i*_CHUNK_SIZE : end]
		//2.2分批切片的类型转换 => [][]byte
		chunkLength := len(chunk)
		fieldList := make([]interface{}, 0, chunkLength)
		for _, v := range chunk {
			fieldList = append(fieldList, v)
		}
		cacheRes, err := redis.ByteSlices(objRedis.Do("HMGET", redis.Args{}.Add(key).AddFlat(fieldList)...))
		if err != nil {
			for i := 0; i < chunkLength; i++ {
				res = append(res, nil)
			}
			zlog.Warn(nil, "cache_mget_error: ", err)
			continue
		} else {
			res = append(res, cacheRes...)
		}
	}
	return res, resErr
}

// HMSet 将一个map存到Redis hash
func (objRedis *Redis) HMSet(key string, fvmap map[string]interface{}) error {
	_, err := objRedis.Do("HMSET", redis.Args{}.Add(key).AddFlat(fvmap)...)
	return err
}

func (objRedis *Redis) HKeys(key string) ([][]byte, error) {
	if res, err := redis.ByteSlices(objRedis.Do("HKEYS", key)); err == redis.ErrNil {
		return nil, nil
	} else {
		return res, err
	}
}

func (objRedis *Redis) HGetAll(key string) ([][]byte, error) {
	if res, err := redis.ByteSlices(objRedis.Do("HGETALL", key)); err == redis.ErrNil {
		return nil, nil
	} else {
		return res, err
	}
}

func (objRedis *Redis) HLen(key string) (int64, error) {
	if res, err := redis.Int64(objRedis.Do("HLEN", key)); err == redis.ErrNil {
		return 0, nil
	} else {
		return res, err
	}
}

func (objRedis *Redis) HVals(key string) ([][]byte, error) {
	if res, err := redis.ByteSlices(objRedis.Do("HVALS", key)); err == redis.ErrNil {
		return nil, nil
	} else {
		return res, err
	}
}

func (objRedis *Redis) HIncrBy(key, field string, value int64) (int64, error) {
	return redis.Int64(objRedis.Do("HINCRBY", key, field, value))
}

func (objRedis *Redis) HExists(key string, field string) (bool, error) {
	if res, err := redis.Bool(objRedis.Do("HEXISTS", key, field)); err == redis.ErrNil {
		return false, nil
	} else {
		return res, err
	}
}

func (objRedis *Redis) HDel(key string, fields ...string) (int64, error) {
	args := packArgs(key, fields)
	if res, err := redis.Int64(objRedis.Do("HDEL", args...)); err == redis.ErrNil {
		return 0, nil
	} else {
		return res, err
	}
}

// 基于游标的迭代器，每次被调用会返回新的游标，在下次迭代时，需要使用这个新游标作为游标参数，以此来延续之前的迭代过程
// param: key
// param: cursor 游标 传""表示开始新迭代
// param: count 每次迭代返回元素的最大值，limit hint，实际数量并不准确=count
// param: pattern 模式参数，符合glob风格  ? (一个字符) * （任意个字符） [] (匹配其中的任意一个字符)  \x (转义字符)
// return: 新的cursor，filed-value map  当返回""，空map时，表示迭代已结束
func (objRedis *Redis) HScan(key string, cursor uint64, pattern string, count int) (uint64, map[string][]byte, error) {
	args := packArgs(key, cursor)
	if pattern != "" {
		args = append(args, "MATCH", pattern)
	}
	if count > 0 {
		args = append(args, "COUNT", count)
	}
	values, err := redis.Values(objRedis.Do("HSCAN", args...))
	if err == redis.ErrNil {
		return 0, nil, nil
	} else if err != nil {
		return 0, nil, err
	}
	return parseScanResults(values)
}

func parseScanResults(results []interface{}) (uint64, map[string][]byte, error) {
	if len(results) != 2 {
		return 0, nil, errors.New("hscan err length")
	}

	cursorIndex, err := strconv.ParseInt(string(results[0].([]byte)), 10, 64)
	if err != nil {
		return 0, nil, err
	}
	result := make(map[string][]byte)
	scanData := results[1].([]interface{})
	for i := 0; i < len(scanData); i = i + 2 {
		key := string(scanData[i].([]byte))
		result[key] = scanData[i+1].([]byte)
	}
	return uint64(cursorIndex), result, nil
}
