package pool

import (
	"errors"
	"fmt"
	"github.com/GitHub121380/golib/zlog"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"sync"
	"sync/atomic"
	"time"
)

var (
	ErrClosed               = errors.New("pool is closed")
	ErrMaxActiveConnReached = errors.New("MaxActiveConnReached")
	ErrExceedMaxWaitTimeout = errors.New("exceed default maxWait timeout")
)

// Config 连接池相关配置
type Config struct {
	// 连接池中拥有的最小连接数
	InitialCap int
	// 最大并发存活连接数
	MaxCap int
	// 最大空闲连接
	MaxIdle int
	// 生成连接的方法
	Factory func() (interface{}, error)
	// 关闭连接的方法
	Close func(interface{}) error
	// 检查连接是否有效的方法
	Ping func(interface{}) error
	// 连接最大空闲时间，超过该事件则将失效
	IdleTimeout time.Duration
	// 获取连接最大等待时间，<=0 不等待
	WaitTimeOut time.Duration
}

type connReq struct {
	idleConn *idleConn
}

// channelPool 存放连接信息
type channelPool struct {
	countMu      sync.Mutex
	openingConns int

	mu    sync.Mutex
	conns chan *idleConn

	connReqs    sync.Map
	nextRequest uint64 // Next key to use in connRequests.

	factory func() (interface{}, error)
	close   func(interface{}) error
	ping    func(interface{}) error

	idleTimeout time.Duration
	waitTimeOut time.Duration
	wait        bool
	maxActive   int
}

type idleConn struct {
	conn interface{}
	t    time.Time
}

// NewChannelPool 初始化连接
func NewChannelPool(poolConfig *Config) (Pool, error) {
	if !(poolConfig.InitialCap <= poolConfig.MaxIdle &&
		poolConfig.MaxCap >= poolConfig.MaxIdle &&
		poolConfig.InitialCap >= 0) {
		return nil, errors.New("invalid capacity settings")
	}
	if poolConfig.Factory == nil {
		return nil, errors.New("invalid factory func settings")
	}
	if poolConfig.Close == nil {
		return nil, errors.New("invalid close func settings")
	}

	c := &channelPool{
		conns:        make(chan *idleConn, poolConfig.MaxIdle),
		factory:      poolConfig.Factory,
		close:        poolConfig.Close,
		idleTimeout:  poolConfig.IdleTimeout,
		maxActive:    poolConfig.MaxCap,
		openingConns: poolConfig.InitialCap,
		waitTimeOut:  poolConfig.WaitTimeOut,
	}

	if c.waitTimeOut > 0 {
		c.wait = true
	}

	if poolConfig.Ping != nil {
		c.ping = poolConfig.Ping
	}

	initTime := time.Now()
	for i := 0; i < poolConfig.InitialCap; i++ {
		conn, err := c.factory()
		if err != nil {
			c.Release()
			return nil, fmt.Errorf("init pool error, factory resp: %s", err)
		}
		c.conns <- &idleConn{conn: conn, t: initTime}
	}

	return c, nil
}

func (c *channelPool) getConns() chan *idleConn {
	c.mu.Lock()
	conns := c.conns
	c.mu.Unlock()
	return conns
}

// Get 从pool中取一个连接
func (c *channelPool) Get(ctx *gin.Context) (interface{}, error) {
	conns := c.getConns()
	if conns == nil {
		return nil, ErrClosed
	}

	for {
		select {
		case wrapConn := <-conns:
			if wrapConn == nil {
				return nil, ErrClosed
			}
			if timeout := c.idleTimeout; timeout > 0 {
				if wrapConn.t.Add(timeout).Before(time.Now()) {
					_ = c.Close(wrapConn.conn)
					continue
				}
			}
			if c.ping != nil {
				if err := c.Ping(wrapConn.conn); err != nil {
					_ = c.Close(wrapConn.conn)
					continue
				}
			}
			return wrapConn.conn, nil
		default:
			c.countMu.Lock()
			if c.openingConns < c.maxActive {
				c.openingConns++ // optimistically
				c.countMu.Unlock()

				conn, err := c.factory()
				if err != nil {
					c.countMu.Lock()
					c.openingConns-- // correct for earlier optimism
					c.countMu.Unlock()
					return nil, err
				}
				return conn, nil
			}
			c.countMu.Unlock()

			if !c.wait {
				return nil, ErrMaxActiveConnReached
			}

			// wait for a conn to reuse
			req := make(chan connReq, 1)
			reqKey := c.nextRequestKeyLocked()
			c.connReqs.Store(reqKey, req)

			select {
			case ret, ok := <-req:
				if !ok {
					zlog.WarnLogger(ctx, "MaxActiveConnReached!", zap.String("prot", "pool"))
					return nil, ErrMaxActiveConnReached
				}

				if timeout := c.idleTimeout; timeout > 0 {
					if ret.idleConn.t.Add(timeout).Before(time.Now()) {
						_ = c.Close(ret.idleConn.conn)
						continue
					}
				}
				return ret.idleConn.conn, nil

			case <-time.After(c.waitTimeOut):
				c.connReqs.Delete(reqKey)
				select {
				case ret, ok := <-req:
					if ok && ret.idleConn != nil {
						_ = c.Put(ret.idleConn.conn)
					}
				}
				return nil, ErrExceedMaxWaitTimeout
			}
		}
	}
}

func (c *channelPool) nextRequestKeyLocked() uint64 {
	return atomic.AddUint64(&c.nextRequest, 1)
}

// Put 将连接放回pool中
func (c *channelPool) Put(conn interface{}) error {
	if conn == nil {
		return errors.New("connection is nil. rejecting")
	}

	if c.conns == nil {
		return c.Close(conn)
	}

	if c.wait {
		start := time.Now()
		reuse := false
		c.connReqs.Range(func(k, v interface{}) bool {
			reqKey, ok1 := k.(uint64)
			req, ok2 := v.(chan connReq)
			if ok1 && ok2 {
				c.connReqs.Delete(reqKey) // Remove from pending requests.
				req <- connReq{
					idleConn: &idleConn{conn: conn, t: start},
				}
				reuse = true
				//zlog.DebugLogger(nil, "put connection to wait queues cost: %f", cost)
				return false
			}

			return true
		})

		if reuse {
			return nil
		}
	}

	start := time.Now()
	select {
	case c.conns <- &idleConn{conn: conn, t: start}:
		//zlog.DebugLogger(nil, "put the connection back to the pool cost: %v , conn que len: %v", cost, len(c.conns))
		return nil
	default:
		// 连接池已满，直接关闭该连接
		//zlog.DebugLogger(nil, "pool is full, close the connection")
		return c.Close(conn)
	}
}

func (c *channelPool) Close(conn interface{}) error {
	if conn == nil {
		return errors.New("connection is nil. rejecting")
	}
	if c.close == nil {
		return nil
	}
	c.countMu.Lock()
	c.openingConns--
	c.countMu.Unlock()
	return c.close(conn)
}

func (c *channelPool) Ping(conn interface{}) error {
	if conn == nil {
		return errors.New("connection is nil")
	}
	return c.ping(conn)
}

func (c *channelPool) Release() {
	c.mu.Lock()
	conns := c.conns
	c.conns = nil
	c.factory = nil
	c.ping = nil
	closeFun := c.close
	c.close = nil
	c.mu.Unlock()

	if conns == nil {
		return
	}

	close(conns)
	for wrapConn := range conns {
		closeFun(wrapConn.conn)
	}
}

func (c *channelPool) Len() int {
	return len(c.getConns())
}
