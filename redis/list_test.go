package redis

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRedis_LPush_LPushX(t *testing.T) {
	setup()
	key := "TestRedis_LPush"
	list := []interface{}{"1", "2", "3", "4", "5"}

	r.Del(key)

	num, err := r.LPushX(key, list[0])
	assert.NoError(t, err)
	assert.Equal(t, 0, num)
	n, err := r.LPush(key, list...)
	assert.NoError(t, err)
	assert.Equal(t, len(list), n)
	time.Sleep(100 * time.Millisecond)
	values, err := r.LRange(key, 0, 100)
	assert.NoError(t, err)
	assert.Equal(t, len(list), len(values))
	for i := range list {
		assert.Equal(t, list[len(list)-i-1], string(values[i]))
	}

}

func TestRedis_RPush_RPushX(t *testing.T) {
	setup()
	key := "TestRedis_RPush"
	list := []interface{}{"1", "2", "3", "4", "5"}

	r.Del(key)

	num, err := r.RPushX(key, list[0])
	assert.NoError(t, err)
	assert.Equal(t, 0, num)
	n, err := r.RPush(key, list...)
	assert.NoError(t, err)
	assert.Equal(t, len(list), n)
	time.Sleep(100 * time.Millisecond)
	values, err := r.LRange(key, 0, 100)
	assert.NoError(t, err)
	assert.Equal(t, len(list), len(values))
	for i := range list {
		assert.Equal(t, list[i], string(values[i]))
	}

}

func TestRedis_LPop_RPop_LLen(t *testing.T) {
	setup()
	key := "TestRedis_LPop_RPop"
	list := []interface{}{"1", "2", "3", "4", "5"}

	r.Del(key)

	_, err := r.RPush(key, list...)
	assert.NoError(t, err)
	value, err := r.LPop(key)
	assert.NoError(t, err)
	assert.Equal(t, list[0], string(value))
	value, err = r.RPop(key)
	assert.NoError(t, err)
	assert.Equal(t, list[4], string(value))
	time.Sleep(50 * time.Millisecond)
	values, err := r.LLen(key)
	assert.NoError(t, err)
	assert.Equal(t, len(list)-2, values)

}

func TestRedis_LIndex_LSet(t *testing.T) {
	setup()
	key := "TestRedis_LIndex_LSet"
	list := []interface{}{"1", "2", "3", "4", "5"}

	r.Del(key)

	_, err := r.RPush(key, list...)
	assert.NoError(t, err)
	value, err := r.LIndex(key, 0)
	assert.NoError(t, err)
	assert.Equal(t, list[0], string(value))
	ok, err := r.LSet(key, 0, "SetValue")
	assert.NoError(t, err)
	assert.True(t, ok)
	time.Sleep(50 * time.Millisecond)
	value, err = r.LIndex(key, 0)
	assert.NoError(t, err)
	assert.Equal(t, "SetValue", string(value))

}
func TestRedis_LRem(t *testing.T) {
	setup()
	key := "TestRedis_LRem"
	list := []interface{}{"1", "1", "1", "4", "5"}

	r.Del(key)

	value, err := r.LRem(key, 2, "1")
	assert.NoError(t, err)
	assert.Equal(t, 0, value)

	_, err1 := r.RPush(key, list...)
	assert.NoError(t, err1)
	value, err = r.LRem(key, 2, "1")
	assert.NoError(t, err)
	assert.Equal(t, 2, value)

	time.Sleep(50 * time.Millisecond)
	values, err := r.LRange(key, 0, 100)
	assert.NoError(t, err)
	assert.Equal(t, len(list)-2, len(values))
	for i := range values {
		assert.Equal(t, list[i+2], string(values[i]))
	}

	_, err1 = r.RPush(key, "1")
	assert.NoError(t, err1)
	value, err = r.LRem(key, -1, "1")
	assert.NoError(t, err)
	assert.Equal(t, 1, value)

	time.Sleep(50 * time.Millisecond)
	values, err = r.LRange(key, 0, 100)
	assert.NoError(t, err)
	assert.Equal(t, len(list)-2, len(values))
	for i := range values {
		assert.Equal(t, string(values[i]), list[i+2])
	}

	_, err = r.RPush(key, "1")
	assert.NoError(t, err)
	value, err = r.LRem(key, 0, "1")
	assert.NoError(t, err)
	assert.Equal(t, 2, value)

	time.Sleep(50 * time.Millisecond)
	values, err = r.LRange(key, 0, 100)
	assert.NoError(t, err)
	assert.Equal(t, len(list)-3, len(values))
	for i := range values {
		assert.Equal(t, string(values[i]), list[i+3])
	}

}

func TestRedis_LInsert(t *testing.T) {
	setup()
	key := "TestRedis_LInsert"
	list := []interface{}{"1", "2", "3", "4", "5"}

	r.Del(key)

	_, err := r.RPush(key, list...)
	assert.NoError(t, err)

	value, err := r.LInsert(key, true, "1", "6")
	assert.NoError(t, err)
	assert.Equal(t, len(list)+1, value)

	time.Sleep(50 * time.Millisecond)
	values, err := r.LRange(key, 0, 100)
	assert.NoError(t, err)
	assert.Equal(t, "6", string(values[0]))

	value, err = r.LInsert(key, false, "7", "6")
	assert.NoError(t, err)
	assert.Equal(t, -1, value)

	r.Del(key)

	value, err = r.LInsert(key, false, "1", "6")
	assert.NoError(t, err)
	assert.Equal(t, 0, value)

}
func TestRedis_LTrim(t *testing.T) {
	setup()
	key := "TestRedis_LTrim"
	list := []interface{}{"1", "2", "3", "4", "5"}

	r.Del(key)

	_, err := r.RPush(key, list...)
	assert.NoError(t, err)
	time.Sleep(100 * time.Millisecond)
	ok, err := r.LTrim(key, 1, -1)
	assert.NoError(t, err)
	assert.True(t, ok)

	values, err := r.LRange(key, 0, 100)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(values))
	list1 := list[1:5]
	for i := range values {
		assert.Equal(t, list1[i], string(values[i]))
	}

}
