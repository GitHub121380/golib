package redis

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var r *Redis
var ServiceName = "zbcourse"
var Hosts = []string{"127.0.0.1:6379"}
var RedisConnConfig = &RedisConfig{
	MaxIdle:     25,
	MaxActive:   25,
	IdleTimeout: 2400,
	Passwd:      "",
	Database:    0,
	Timeout:     5,
}

func init() {
	// log.Setup("testredis", "/home/homework/go-common/conf/redis", "./")
	setup()
}

func setup() {
	RedisInit(ServiceName, Hosts, RedisConnConfig)
	objRedis, err := GetInstance(nil, ServiceName)
	if err != nil {
		fmt.Println("setup fail")
	} else {
		r = objRedis
	}
}

func TestGetInstance(t *testing.T) {
	t.Run("GetInstance", func(t *testing.T) {
		objRedis, err := GetInstance(nil, ServiceName)
		assert.Equal(t, nil, err, fmt.Sprintf("GetInstance error：%v", err))
		fmt.Println("objRedis:", objRedis)

	})
}

func TestDo(t *testing.T) {
	setup()
	t.Run("Do", func(t *testing.T) {
		reply, err := r.Do("ping")
		assert.Equal(t, "PONG", reply)
		assert.Equal(t, nil, err)
	})
}

func TestSend(t *testing.T) {

}
