package redis

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRedis_SAdd_SIsMember_SMembers(t *testing.T) {
	setup()
	key, values := "TestRedis_SAdd_SIsMember_SMembers", []string{"1", "2", "3", "3"}

	r.Del(key)

	n, err := r.SAdd(key, values...)
	assert.NoError(t, err)
	assert.Equal(t, int64(len(values))-1, n)

	time.Sleep(50 * time.Millisecond)
	ok, err := r.SIsMember(key, values[0])
	assert.NoError(t, err)
	assert.True(t, ok)

	members, err := r.SMembers(key)
	assert.NoError(t, err)
	fmt.Println("members len :", len(members))
	assert.Equal(t, len(values)-1, len(members))

}

func TestRedis_SRem_SCard(t *testing.T) {
	setup()
	key, values := "TestRedis_SAdd_SIsMember_SMembers", []string{"1", "2", "3", "3"}

	r.Del(key)

	n, err := r.SAdd(key, values...)
	assert.NoError(t, err)
	assert.Equal(t, int64(len(values))-1, n)

	time.Sleep(50 * time.Millisecond)
	ok, err := r.SIsMember(key, values[0])
	assert.NoError(t, err)
	assert.True(t, ok)

	members, err := r.SMembers(key)
	assert.NoError(t, err)
	assert.Equal(t, len(values)-1, len(members))

}

func TestRedis_SPop_SRandMember(t *testing.T) {
	setup()
	key, values := "TestRedis_SAdd_SIsMember_SMembers", []string{"1", "2", "3"}

	r.Del(key)

	_, err := r.SAdd(key, values...)

	assert.NoError(t, err)

	value, err := r.SRandMember(key)

	assert.NoError(t, err)
	assert.True(t, string(value) >= "1" && string(value) <= "3")
	members, err := r.SMembers(key)
	assert.NoError(t, err)
	assert.Equal(t, len(values), len(members))

	value, err = r.SPop(key)
	assert.NoError(t, err)
	assert.True(t, string(value) >= "1" && string(value) <= "3")
	time.Sleep(100 * time.Millisecond)
	members, err = r.SMembers(key)
	assert.NoError(t, err)
	assert.Equal(t, len(values)-1, len(members))

}
func TestSCard(t *testing.T) {
	setup()
	key := "TestSCard"
	r.Del(key)
	r.SAdd(key, "value")
	//主从同步延迟
	time.Sleep(100 * time.Millisecond)
	if n, err := r.SCard(key); err != nil {
		t.Error(err)
	} else if n != 1 {
		t.Fail()
	}
}

func TestSPop(t *testing.T) {
	setup()
	key := "TestSPop"
	r.Del(key)
	r.SAdd(key, "value")
	//主从同步延迟
	time.Sleep(100 * time.Millisecond)
	if item, err := r.SPop(key); err != nil {
		t.Error(err)
	} else if item == nil {
		t.Fail()
	} else if string(item) != "value" {
		t.Fail()
	}
	if item, _ := r.SPop(key); item != nil {
		t.Fail()
	}
}

func TestSRem(t *testing.T) {
	setup()
	key := "TestSRem"
	r.Del(key)
	r.SAdd(key, "one", "two", "three")
	//主从同步延迟
	time.Sleep(100 * time.Millisecond)
	if n, err := r.SRem(key, "one", "four"); err != nil {
		t.Error(err)
	} else if n != 1 {
		t.Fail()
	}
}

func TestSScan(t *testing.T) {
	setup()
	key := "TestSScan"
	r.Del(key)
	r.SAdd(key, "one", "two", "three")
	//主从同步延迟
	time.Sleep(100 * time.Millisecond)
	if _, list, err := r.SScan(key, 0, "", 0); err != nil {
		t.Error(err)
	} else if len(list) == 0 {
		t.Fail()
	}
}
