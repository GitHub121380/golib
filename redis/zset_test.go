package redis

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRedis_ZAdd_ZCard_ZRange_ZRevRange(t *testing.T) {
	setup()
	key := "TestRedis_ZAdd_ZCard_ZRange_ZRevRange"
	pairs := map[string]float64{
		"one":   1.0,
		"three": 3.0,
		"two":   1.0,
	}
	keys := []string{"one", "two", "three"}
	value := []string{"1", "2", "3"}
	keys1 := []string{"three", "two", "one"}
	value1 := []string{"3", "2", "1"}

	_, _ = r.Del(key)
	n, err := r.ZAdd(key, pairs)

	assert.NoError(t, err)
	assert.Equal(t, int64(len(pairs)), n)

	n, err = r.ZAdd(key, map[string]float64{"two": 2.0})
	assert.NoError(t, err)
	assert.Equal(t, 0, int(n))

	num, err := r.ZCard(key)
	assert.NoError(t, err)
	assert.Equal(t, int64(len(pairs)), num)

	time.Sleep(50 * time.Millisecond)
	values, err := r.ZRange(key, 0, -1, true)
	assert.NoError(t, err)
	for i := range values {
		if i%2 == 0 {
			assert.Equal(t, keys[i/2], string(values[i]))
		} else {
			assert.Equal(t, value[(i-1)/2], string(values[i]))
		}
	}

	time.Sleep(50 * time.Millisecond)
	values, err = r.ZRange(key, -2, -1, false)
	assert.NoError(t, err)
	for i := range values {
		assert.Equal(t, keys[1:][i], string(values[i]))
	}

	time.Sleep(50 * time.Millisecond)
	values, err = r.ZRevRange(key, 0, -1, true)
	assert.NoError(t, err)
	for i := range values {
		if i%2 == 0 {
			assert.Equal(t, keys1[i/2], string(values[i]))
		} else {
			assert.Equal(t, value1[(i-1)/2], string(values[i]))
		}
	}

}

func TestRedis_ZRangeByScore_ZRevRangeByScore_ZCount(t *testing.T) {
	setup()
	key := "TestRedis_ZRangeByScore_ZRevRangeByScore_ZCount"
	pairs := map[string]float64{
		"one":   1.0,
		"two":   2.0,
		"three": 3.0,
		"four":  4.0,
	}

	_, _ = r.Del(key)
	n, err := r.ZAdd(key, pairs)

	assert.NoError(t, err)
	assert.Equal(t, int64(len(pairs)), n)

	time.Sleep(50 * time.Millisecond)
	value, err := r.ZRangeByScore(key, "-inf", "(5", true, true, 2, 2)
	assert.NoError(t, err)
	keys1 := []string{"three", "four"}
	value1 := []string{"3", "4"}
	for i := range value {
		if i%2 == 0 {
			assert.Equal(t, keys1[i/2], string(value[i]))
		} else {
			assert.Equal(t, value1[(i-1)/2], string(value[i]))
		}
	}

	time.Sleep(50 * time.Millisecond)
	value, err = r.ZRevRangeByScore(key, "+inf", "3", true, true, 0, 2)
	assert.NoError(t, err)
	keys2 := []string{"four", "three"}
	value2 := []string{"4", "3"}
	for i := range value {
		if i%2 == 0 {
			assert.Equal(t, keys2[i/2], string(value[i]))
		} else {
			assert.Equal(t, value2[(i-1)/2], string(value[i]))
		}
	}

	num, err := r.ZCount(key, "1.0", "2.0")
	assert.NoError(t, err)
	assert.Equal(t, 2, int(num))

}

func TestRedis_ZRank_ZRevRank(t *testing.T) {
	setup()
	key := "TestRedis_ZRank_ZRevRank"
	pairs := map[string]float64{
		"one":   1.0,
		"two":   2.0,
		"three": 3.0,
	}

	_, _ = r.Del(key)
	n, err := r.ZAdd(key, pairs)

	assert.NoError(t, err)
	assert.Equal(t, int64(len(pairs)), n)
	time.Sleep(100 * time.Millisecond)
	num, err := r.ZRank(key, "two")
	assert.NoError(t, err)
	assert.Equal(t, 1, int(num))

	num, err = r.ZRevRank(key, "one")
	assert.NoError(t, err)
	assert.Equal(t, 2, int(num))

}

func TestRedis_ZRem_ZRemRangeByRank_ZRemRangeByScore(t *testing.T) {
	setup()
	key := "TestRedis_ZRem_ZRemRangeByRank_ZRemRangeByScore"
	pairs := map[string]float64{
		"one":   1.0,
		"two":   2.0,
		"three": 3.0,
	}
	_, _ = r.Del(key)
	n, err := r.ZAdd(key, pairs)
	assert.NoError(t, err)
	assert.Equal(t, int64(len(pairs)), n)

	num, err := r.ZRem(key, "one")
	assert.NoError(t, err)
	assert.Equal(t, 1, int(num))

	_, err = r.ZAdd(key, map[string]float64{"one": 1.0})
	assert.NoError(t, err)
	value, err := r.ZRemRangeByRank(key, 0, 1)
	assert.NoError(t, err)
	assert.Equal(t, 2, int(value))

	_, err = r.ZAdd(key, map[string]float64{"one": 1.0, "two": 2.0})
	assert.NoError(t, err)
	values, err := r.ZRemRangeByScore(key, "-inf", "2.0")
	assert.NoError(t, err)
	assert.Equal(t, 2, int(values))

}

func TestZIncrBy(t *testing.T) {
	setup()
	key := "TestRedis_ZIncrBy"
	pairs := map[string]float64{
		"one":   1.0,
		"two":   2.0,
		"three": 3.0,
	}

	_, _ = r.Del(key)
	n, err := r.ZAdd(key, pairs)

	assert.NoError(t, err)
	assert.Equal(t, int64(len(pairs)), n)

	num, err := r.ZIncrBy(key, 3.0, "one")
	assert.NoError(t, err)
	assert.Equal(t, float64(4.0), num)

}

func TestZLexCount(t *testing.T) {
	setup()
	key := "TestRedis_ZLexCount"
	pairs := map[string]float64{
		"one":   1.0,
		"two":   1.0,
		"three": 1.0,
	}

	_, _ = r.Del(key)
	n, err := r.ZAdd(key, pairs)

	assert.NoError(t, err)
	assert.Equal(t, int64(len(pairs)), n)
	time.Sleep(100 * time.Millisecond)
	num, err := r.ZLexCount(key, "-", "+")
	assert.NoError(t, err)
	assert.Equal(t, int64(3), num)

}

func TestRedis_ZRemRangeByLex(t *testing.T) {
	setup()
	key := "TestRedis_ZRemRangeByLex"
	pairs := map[string]float64{
		"v1":  1.0,
		"v11": 1.0,
		"v12": 1.0,
	}
	_, _ = r.Del(key)
	n, err := r.ZAdd(key, pairs)

	assert.NoError(t, err)
	assert.Equal(t, int64(len(pairs)), n)

	num, err := r.ZRemRangeByLex(key, "-", "[v11")
	assert.NoError(t, err)
	assert.Equal(t, 2, int(num))

}

func TestRedis_ZScore(t *testing.T) {
	setup()
	key := "TestRedis_ZScore"
	pairs := map[string]float64{
		"one":   1.0,
		"two":   2.0,
		"three": 3.0,
	}

	_, _ = r.Del(key)
	n, err := r.ZAdd(key, pairs)

	assert.NoError(t, err)
	assert.Equal(t, int64(len(pairs)), n)
	time.Sleep(100 * time.Millisecond)
	value, err := r.ZScore(key, "one")
	value1, err := strconv.ParseFloat(value, 64)
	assert.NoError(t, err)
	assert.Equal(t, 1.0, float64(value1))

}

func TestRedis_ZScan(t *testing.T) {
	setup()
	key := "TestRedis_ZScan"
	pairs := map[string]float64{
		"one":   1.0,
		"two":   2.0,
		"three": 3.0,
	}

	_, _ = r.Del(key)
	n, err := r.ZAdd(key, pairs)

	assert.NoError(t, err)
	assert.Equal(t, int64(len(pairs)), n)

	num, values, err := r.ZScan(key, 0, "*o*", 3)
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), num)
	list := []string{"one", "1", "two", "2"}
	for i := range values {
		assert.Equal(t, list[i], values[i])
	}

}
