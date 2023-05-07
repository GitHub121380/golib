package redis

import (
	"errors"
	"github.com/GitHub121380/golib/zlog"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gomodule/redigo/redis"
)

var redisRefreshRWLock = sync.RWMutex{}

type redisPools struct {
	hostport string
	pool     *redis.Pool
}

type RedisConfig struct {
	HostPort    string
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
	Passwd      string
	Database    int
	Timeout     time.Duration
}

var serviceMap = &sync.Map{}

func hasServerName(name string) bool {
	_, exist := serviceMap.Load(name)
	return exist
}

// 调用Redis前必须初始化该方法
func RedisInit(name string, hosts []string, redisConf *RedisConfig) error {
	if len(name) <= 0 || len(hosts) <= 0 {
		return ParamsErr
	}

	if hasServerName(name) {
		return HasServerNameErr
	}

	zlog.Infof(nil, "RedisInit name: %s, host: %v, config: %v ", name, hosts, redisConf)
	//4.写入配置及对应的连接池，添加到slice中,声明大小为当前2倍，避免bns刷新，切片扩容
	rPools := make([]*redisPools, len(hosts), 2*len(hosts))
	hIndex := make(map[string]int, len(hosts))
	for index, host := range hosts {
		conf := RedisConfig{
			HostPort:    host,
			MaxIdle:     redisConf.MaxIdle,
			MaxActive:   redisConf.MaxActive,
			IdleTimeout: redisConf.IdleTimeout,
			Passwd:      redisConf.Passwd,
			Database:    redisConf.Database,
			Timeout:     redisConf.Timeout,
		}
		rPools[index] = &redisPools{
			hostport: host,
			pool:     getPool(conf),
		}
		hIndex[host] = index
	}
	//5.写入全局变量并获取连接池
	serviceMap.Store(name, &RedisClient{
		Service: name,
		config:  redisConf,
		pools:   rPools,
		hIndex:  hIndex,
		mu:      &sync.RWMutex{},
	})
	return nil
}

func getPool(conf RedisConfig) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     conf.MaxIdle,
		MaxActive:   conf.MaxActive,
		IdleTimeout: conf.IdleTimeout * time.Second,
		Wait:        true,
		Dial: func() (conn redis.Conn, e error) {
			con, err := redis.Dial("tcp", conf.HostPort,
				redis.DialPassword(conf.Passwd),
				redis.DialDatabase(conf.Database),
				redis.DialConnectTimeout(conf.Timeout*time.Second),
				redis.DialReadTimeout(conf.Timeout*time.Second),
				redis.DialWriteTimeout(conf.Timeout*time.Second),
			)

			if err != nil {
				zlog.Warn(nil, "get_redis_conn_fail: ", err)
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
}

// 获取连接
func (objRedis *RedisClient) getConn() (redis.Conn, error) {
	defer func() {
		if err := recover(); err != nil {
			zlog.Errorf(nil, "Redis getConn err : %s", err)
		}
	}()

	objRedis.mu.RLock()
	defer objRedis.mu.RUnlock()
	redises := objRedis.pools
	redisLen := uint64(len(redises))

	if len(redises) <= 0 {
		return nil, errors.New("Invalid redis host list")
	}

	tryCount := 0
	for tryCount < 3 {
		tryCount++
		indexIncr := atomic.AddUint64(&objRedis.roundRobin, 1)
		index := indexIncr % redisLen
		redisConn := redises[index]
		conn := redisConn.getConn()
		if _, err := conn.Do("ping"); err != nil {
			zlog.Warn(nil, "ping redis host conn fail, ip port : ", redisConn.hostport)
			conn.Close()
			continue
		}
		return conn, nil
	}
	return nil, GetRedisConnErr
}

// 获取连接
func (rp *redisPools) getConn() redis.Conn {
	conn := rp.pool.Get()
	return conn
}

// 对外暴露：新增，删除redis连接池函数
func RedisRefresh(serviceName string, slAdds []string, slDels []string) error {
	defer func() {
		if err := recover(); err != nil {
			zlog.Errorf(nil, "RedisBNSRefresh err: %s", err)
		}
	}()
	if !hasServerName(serviceName) {
		return NotExistServerNameErr
	}
	//无需操作任何东西
	if len(slAdds) <= 0 && len(slDels) <= 0 {
		return nil
	}
	redisRefreshRWLock.Lock()
	defer redisRefreshRWLock.Unlock()
	if val, exist := serviceMap.Load(serviceName); exist {
		sRedis := val.(*RedisClient)
		//删除节点比现存加新增的还多，认为bns刷新异常
		if len(sRedis.pools)+len(slAdds) < len(slDels) {
			zlog.Warnf(nil, "server name : %s, remove node is more than now", serviceName)
			return RemoveNodeErr
		}
		newAddRedis(sRedis, slAdds)
		removeRedis(sRedis, slDels)
	}
	return NotExistServerNameErr
}

func newAddRedis(serviceRedis *RedisClient, slAdds []string) {
	//遍历新加的节点,并初始化redis链接池
	for _, host := range slAdds {
		if _, exist := serviceRedis.hIndex[host]; !exist {
			conf := RedisConfig{
				HostPort:    host,
				MaxIdle:     serviceRedis.config.MaxIdle,
				MaxActive:   serviceRedis.config.MaxActive,
				IdleTimeout: serviceRedis.config.IdleTimeout,
				Passwd:      serviceRedis.config.Passwd,
				Database:    serviceRedis.config.Database,
				Timeout:     serviceRedis.config.Timeout,
			}
			newRedisPool := &redisPools{
				hostport: host,
				pool:     getPool(conf),
			}
			serviceRedis.pools = append(serviceRedis.pools, newRedisPool)
			serviceRedis.hIndex[host] = len(serviceRedis.pools) - 1
			zlog.Info(nil, "RedisBNSRefresh New Redis Host:", host, " current Length:", len(serviceRedis.pools))
		}
	}
	//重新构建hostToRedisIndex
	rebuildHostToIndex(serviceRedis)
}

func removeRedis(serviceRedis *RedisClient, slDels []string) {
	//遍历下线的节点,并从缓存池中下线
	for _, host := range slDels {
		if removeIndex, exist := serviceRedis.hIndex[host]; exist {
			if len(serviceRedis.pools) > removeIndex {
				//下线节点IP
				closePool := serviceRedis.pools[removeIndex].pool
				serviceRedis.pools = append(serviceRedis.pools[:removeIndex], serviceRedis.pools[removeIndex+1:]...)
				//重新构建hostToRedisIndex
				rebuildHostToIndex(serviceRedis)

				err := closePool.Close()
				if err != nil {
					zlog.Error(nil, "RedisBNSRefresh closePool Close host :", host, "err :", err)
				} else {
					zlog.Info(nil, "RedisBNSRefresh Remove Redis Host:", host)
				}
			}
		}
	}
}

func rebuildHostToIndex(serviceRedis *RedisClient) {
	newHostToRedisIndex := make(map[string]int)
	for index, interRedis := range serviceRedis.pools {
		newHostToRedisIndex[interRedis.hostport] = index
	}
	serviceRedis.hIndex = newHostToRedisIndex
	zlog.Info(nil, "RedisRefresh rebuildHostToIndex current hIndex: ", serviceRedis.hIndex)
}
