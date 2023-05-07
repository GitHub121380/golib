package redis

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRedis_HSet_HGet_HLen(t *testing.T) {
	setup()
	key := "TestHSetHGet"
	fieldMap := map[string]string{
		"TestHSetHGet_key1": "TestHSetHGet_Value1",
		"TestHSetHGet_key2": "TestHSetHGet_Value2",
		"TestHSetHGet_key3": "TestHSetHGet_Value3",
	}
	fieldNum := int64(0)

	r.Del(key)

	time.Sleep(100 * time.Millisecond)
	n, err := r.HLen(key)
	assert.NoError(t, err)
	assert.Equal(t, fieldNum, n)

	for k, v := range fieldMap {
		_, err = r.HSet(key, k, v)
		assert.NoError(t, err)

		fieldNum += 1
		time.Sleep(50 * time.Millisecond)
		data, err := r.HGet(key, k)
		assert.NoError(t, err)
		assert.Equal(t, v, string(data))

		data, err = r.HGet(key, k)
		assert.NoError(t, err)
		assert.Equal(t, v, string(data))

		n, err = r.HLen(key)
		assert.NoError(t, err)
		assert.Equal(t, fieldNum, n)

		n, err = r.HLen(key)
		assert.NoError(t, err)
		assert.Equal(t, fieldNum, n)
	}

}

func TestRedis_HMSetHMGet(t *testing.T) {
	setup()
	key := "TestRedis_HMSetHMGet"
	fieldMap := map[string]interface{}{
		"TestRedis_HMSetHMGet_field1": "TestRedis_HMSetHMGet_Value1",
		"TestRedis_HMSetHMGet_field2": "TestRedis_HMSetHMGet_Value2",
		"TestRedis_HMSetHMGet_field3": "TestRedis_HMSetHMGet_Value3",
	}
	fields := []string{"TestRedis_HMSetHMGet_field1", "TestRedis_HMSetHMGet_field2", "TestRedis_HMSetHMGet_field3"}

	r.Del(key)

	values, err := r.HMGet(key, fields...)
	// values保证不会返回nil
	assert.NotEqual(t, 0, len(values))
	assert.Equal(t, 3, len(values))
	assert.NoError(t, err)
	err = r.HMSet(key, fieldMap)
	assert.NoError(t, err)

	time.Sleep(50 * time.Millisecond)
	values, err = r.HMGet(key, fields...)
	assert.NoError(t, err)
	assert.Equal(t, fieldMap[fields[0]], string(values[0]))
	assert.Equal(t, fieldMap[fields[1]], string(values[1]))
	assert.Equal(t, fieldMap[fields[2]], string(values[2]))

	values, err = r.HMGet(key, fields...)
	assert.NoError(t, err)
	assert.Equal(t, fieldMap[fields[0]], string(values[0]))
	assert.Equal(t, fieldMap[fields[1]], string(values[1]))
	assert.Equal(t, fieldMap[fields[2]], string(values[2]))
}

func TestRedis_HDel(t *testing.T) {
	setup()
	key := "TestRedis_HDel"
	fieldMap := map[string]interface{}{
		"TestRedis_HDel_field1": "TestRedis_HDel_Value1",
		"TestRedis_HDel_field2": "TestRedis_HDel_Value2",
		"TestRedis_HDel_field3": "TestRedis_HDel_Value3",
	}

	r.Del(key)

	// key 不存在时删除key中的field
	n, err := r.HDel(key, "TestRedis_HDel_field1")
	assert.NoError(t, err)

	err = r.HMSet(key, fieldMap)
	assert.NoError(t, err)

	n, err = r.HDel(key, "TestRedis_HDel_field1")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), n)

	// 检查被删除的key是否存在
	time.Sleep(50 * time.Millisecond)
	value, err := r.HGet(key, "TestRedis_HDel_field1")
	assert.NoError(t, err)
	assert.Equal(t, "", string(value))

	value, err = r.HGet(key, "TestRedis_HDel_field1")
	assert.NoError(t, err)
	assert.Equal(t, "", string(value))
}

func TestRedis_HKeys_HVals(t *testing.T) {
	setup()
	key := "TestRedis_HKeys"
	fieldMap := map[string]interface{}{
		"TestRedis_HKeys_field1": "TestRedis_HDel_Value1",
		"TestRedis_HKeys_field2": "TestRedis_HDel_Value2",
		"TestRedis_HKeys_field3": "TestRedis_HDel_Value3",
	}
	_, _ = r.Del(key)
	err := r.HMSet(key, fieldMap)
	assert.NoError(t, err)

	time.Sleep(50 * time.Millisecond)
	keys, err := r.HKeys(key)
	assert.NoError(t, err)
	assert.Equal(t, len(fieldMap), len(keys))
	for i := 0; i != len(fieldMap); i++ {
		_, ok := fieldMap[string(keys[i])]
		assert.True(t, ok)
	}

	values, err := r.HVals(key)
	assert.NoError(t, err)
	assert.Equal(t, len(fieldMap), len(values))
	for i := 0; i != len(fieldMap); i++ {
		assert.Truef(t, strings.HasPrefix(string(values[i]), "TestRedis_HDel_Value"), "value not has prefix:%s", string(values[i]))
	}

}

func TestRedis_HIncrBy(t *testing.T) {
	setup()
	key := "TestRedis_HIncrBy_Key"
	fieldMap := map[string]interface{}{
		"TestRedis_HIncrBy_field1": 0,
	}

	r.Del(key)

	err := r.HMSet(key, fieldMap)
	assert.NoError(t, err)
	value, err := r.HIncrBy(key, "TestRedis_HIncrBy_field1", 1)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), value)

	time.Sleep(100 * time.Millisecond)
	data, err := r.HGet(key, "TestRedis_HIncrBy_field1")
	assert.NoError(t, err)
	assert.Equal(t, "1", string(data))

}

func TestRedis_HScan(t *testing.T) {
	setup()
	key := "TestRedis_HScan_Key"
	fieldMap := map[string]interface{}{
		"TestRedis_HScan_field1": "1",
		"TestRedis_HScan_field2": "2",
		"TestRedis_HScan_field3": "3",
		"TestRedis_HScan_field4": "4",
	}
	r.Del(key)

	_ = r.HMSet(key, fieldMap)
	time.Sleep(50 * time.Millisecond)
	n, fields, err := r.HScan(key, 1, "", 4)
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), n)
	for k, v := range fieldMap {
		assert.Equal(t, v, string(fields[k]))
	}

}
