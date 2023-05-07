package redis

import (
	"fmt"
	"github.com/GitHub121380/golib/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"sync"
	"time"

	"github.com/GitHub121380/golib/ral"
	"github.com/GitHub121380/golib/zlog"
	"github.com/gomodule/redigo/redis"
)

// 日志打印Do args部分支持的最大长度
const logForRedisValue = 50

type RedisClient struct {
	ctx  *gin.Context
	ins  *ral.Instance
	pool *redis.Pool

	Service    string         // 服务单元名称
	roundRobin uint64         // 轮询计数
	config     *RedisConfig   // 配置记录
	mu         *sync.RWMutex  // 操作锁
	pools      []*redisPools  // redis连接池
	hIndex     map[string]int // host对应的pools的下标位置，用于变更记录
}

type Redis struct {
	r   *RedisClient
	ctx *gin.Context
}

func (objRedis *Redis) Send(commandName string, args ...interface{}) (err error) {
	objRedis.r.ins.Retry(func(res *ral.Resource, ins *ral.Instance) bool {
		if r, ok := ins.Client.(*RedisClient); ok {
			conn := r.pool.Get()
			defer conn.Close()

			if err = conn.Send(commandName, args...); err == nil {
				conn.Flush()
				_, err = conn.Receive()
				return false
			}
		}
		return true
	})
	return
}

func (objRedis *Redis) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	start := time.Now()
	remoteIp, remotePort := objRedis.r.ins.IP, objRedis.r.ins.Port
	retry := objRedis.r.ins.Retry(func(res *ral.Resource, ins *ral.Instance) bool {
		if r, ok := ins.Client.(*RedisClient); ok {
			conn := r.pool.Get()
			defer conn.Close()

			if reply, err = conn.Do(commandName, args...); err == nil {
				remoteIp, remotePort = ins.IP, ins.Port
				return false
			} else {
				zlog.WarnLogger(objRedis.ctx, err.Error(), zap.String("service", objRedis.r.Service))
			}
		}
		return true
	})

	end := time.Now()

	fields := []zap.Field{
		zap.String("service", objRedis.r.Service),
		zap.String("requestStartTime", utils.GetFormatRequestTime(start)),
		zap.String("requestEndTime", utils.GetFormatRequestTime(end)),
		zap.Float64("cost", utils.GetRequestCost(start, end)),
		zap.String("command", commandName),
		zap.String("commandVal", utils.JoinArgs(logForRedisValue, args)),
		zap.String("remoteAddr", fmt.Sprintf("%s:%d", remoteIp, remotePort)),
		zap.Int("retry", retry),
	}

	// 执行时间 单位:毫秒
	msg := "redis do success"
	if err != nil {
		msg = fmt.Sprintf("redis do error: %s", err.Error())
	}
	zlog.InfoLogger(objRedis.ctx, msg, fields...)
	return reply, nil
}

func (objRedis *Redis) Release() {
	if objRedis.r.ins != nil {
		objRedis.r.ins.Release()
	}
	objRedis.ctx = nil
}

var mod *ral.Module

func init() {
	appendFun := func(self *ral.Resource, res *ral.Resource, ins *ral.Instance) (*ral.Instance, error) {
		p := &redis.Pool{
			MaxIdle:     res.Redis.MaxIdle,
			MaxActive:   res.Redis.MaxActive,
			IdleTimeout: res.Redis.IdleTimeout,
			Wait:        res.Redis.Wait,
			Dial: func() (conn redis.Conn, e error) {
				con, err := redis.Dial(
					"tcp",
					fmt.Sprintf("%s:%d", ins.IP, ins.Port),
					redis.DialPassword(res.Redis.Password),
					redis.DialDatabase(res.Redis.Database),
					redis.DialConnectTimeout(res.ConnTimeOut),
					redis.DialReadTimeout(res.ReadTimeOut),
					redis.DialWriteTimeout(res.WriteTimeOut),
				)
				if err != nil {
					zlog.WarnLogger(nil, "get_redis_conn_fail: "+err.Error(), zap.String("prot", "redis"))
					return nil, err
				}
				return con, nil
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				if time.Since(t) < time.Minute {
					return nil
				}
				_, err := c.Do("PING")
				return err
			},
		}
		client := &RedisClient{
			pool:    p,
			Service: res.Name,
		}
		sub := &ral.Instance{
			IP:     ins.IP,
			Port:   ins.Port,
			Client: client,
		}
		client.ins = sub
		return sub, nil
	}

	removeFun := func(self *ral.Instance, res *ral.Resource, ins *ral.Instance) bool {
		if r, ok := ins.Client.(*RedisClient); ok && r.pool != nil {
			if err := r.pool.Close(); err != nil {
				zlog.WarnLogger(nil, "redis pool close error: "+err.Error(), zap.String("prot", "redis"))
			}
		}
		return true
	}

	m := &ral.Module{
		Type: ral.TYPE_REDIS,
		Config: func(res *ral.Resource) *ral.Resource {
			return res
		},
		Append: appendFun,
		Remove: removeFun,
		Method: func(ctx *gin.Context, res *ral.Resource, ins *ral.Instance, method string, data map[string]interface{}, head map[string]string) (interface{}, error) {
			return nil, nil
		}}

	mod = ral.AddModule(m)
}

func GetInstance(ctx *gin.Context, service string) (*Redis, error) {
	ins, err := ral.GetInstance(ctx, ral.TYPE_REDIS, service)
	if err != nil {
		zlog.WarnLogger(ctx, "redis GetInstance error: not found client", zap.String("prot", "redis"))
		return nil, ral.ERR_NOT_FOUND_CLIENT
	}

	r, ok := ins.Client.(*RedisClient)
	if !ok {
		zlog.WarnLogger(ctx, "redis GetInstance error: ins.client type invalid", zap.String("prot", "redis"))
		return nil, ral.ERR_NOT_FOUND_CLIENT
	}

	c := &Redis{
		r:   r,
		ctx: ctx,
	}
	return c, nil
}
