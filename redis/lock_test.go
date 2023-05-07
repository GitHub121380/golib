package redis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRedis_SetNxByEX(t *testing.T) {
	setup()

	type args struct {
		key    string
		value  interface{}
		expire uint64
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "setLock",
			args: args{
				key:    "setex",
				value:  "1",
				expire: 20,
			},
		},
		{
			name: "setLocked",
			args: args{
				key:    "setex2",
				value:  "2",
				expire: 20,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := r.SetNxByEX(tt.args.key, tt.args.value, tt.args.expire)
			assert.NoError(t, err)
			assert.True(t, res)
		})
	}
}

func TestRedis_SetNxByPX(t *testing.T) {
	setup()
	objRedis, err := GetInstance("")
	if err != nil {
		log.Error("[FixTaskMedal]Get redis error.ERROR:", err)
		return
	}
	type args struct {
		key    string
		value  interface{}
		expire uint64
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "setLock",
			args: args{
				key:    "setex",
				value:  "1",
				expire: 20000,
			},
		},
		{
			name: "setLocked",
			args: args{
				key:    "setex",
				value:  "2",
				expire: 20000,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := objRedis.SetNxByPX(tt.args.key, tt.args.value, tt.args.expire)
			if err != nil {
				log.Errorf("Redis.SetNxByPX() error = %v, res %v", err, res)
				return
			}
			if res == true {
				log.Errorf("Redis.SetNxByPX() error = %v, res %v", err, res)
			} else {
				log.Errorf("Redis.SetNxByPX() error = %v, res %v", err, res)
			}
		})
	}
}
