package hbase

import (
	"errors"
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/GitHub121380/golib/pool"
	"github.com/GitHub121380/golib/ral"
	"github.com/GitHub121380/golib/utils"
	"github.com/GitHub121380/golib/zlog"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net"
	"strconv"
	"time"
)

type HbaseClientModule struct {
	Username string
	Password string
	Service  string

	pool pool.Pool

	ins *ral.Instance
	res *ral.Resource
}

type Pool struct {
	pool pool.Pool
}

type HBasePoolClient struct {
	Client *HbaseClient
	Trans  thrift.TTransport
	pool   *Pool
}

func (h *HbaseClientModule) Exec(ctx *gin.Context, efunc func(c *HbaseClient) error) (err error) {
	zlog.CreateSpan(ctx)
	start := time.Now()
	remoteIp, remotePort := h.ins.IP, h.ins.Port
	retry := h.ins.Retry(func(res *ral.Resource, ins *ral.Instance) bool {
		if r, ok := ins.Client.(*HbaseClientModule); ok {
			var conn interface{}
			conn, err = r.pool.Get(ctx)
			if err != nil {
				return true
			}

			defer r.pool.Put(conn)
			c, ok := conn.(*HBasePoolClient)
			if !ok || c == nil || c.Client == nil {
				return true
			}

			if err = efunc(c.Client); err == nil {
				return false
			} else {
				zlog.WarnLogger(ctx, err.Error(), zap.String("service", h.Service))
			}
		}
		return true
	})
	end := time.Now()

	fields := []zap.Field{
		zap.String("service", h.Service),
		zap.String("requestStartTime", utils.GetFormatRequestTime(start)),
		zap.String("requestEndTime", utils.GetFormatRequestTime(end)),
		zap.Float64("cost", utils.GetRequestCost(start, end)), // 执行时间 单位:毫秒
		zap.String("remoteAddr", fmt.Sprintf("%s:%d", remoteIp, remotePort)),
		zap.Int("retry", retry),
	}

	msg := "hbase exec success"
	if err != nil {
		msg = fmt.Sprintf("hbase exec error: %s", err.Error())
	}

	zlog.InfoLogger(ctx, msg, fields...)
	return
}

var mod *ral.Module

func init() {
	m := &ral.Module{
		Type: ral.TYPE_HBASE,
		Config: func(res *ral.Resource) *ral.Resource {
			if res.HBase.IdleTimeout == 0 {
				res.HBase.IdleTimeout = 15 * time.Second
			}
			if res.HBase.MaxActive == 0 {
				res.HBase.MaxActive = 10
			}
			if res.HBase.MaxIdle == 0 {
				res.HBase.MaxIdle = 5
			}
			return res
		},
		Append: func(self *ral.Resource, res *ral.Resource, ins *ral.Instance) (*ral.Instance, error) {
			// 为当前实例打开一个连接池
			addr := net.JoinHostPort(ins.IP, strconv.Itoa(ins.Port))
			max := func(a, b time.Duration) time.Duration {
				if a > b {
					return a
				}
				return b
			}
			timeout := max(max(res.ReadTimeOut, res.WriteTimeOut), res.ConnTimeOut)

			createConn := func() (interface{}, error) {
				trans, err := thrift.NewTSocket(addr)
				if err != nil {
					return nil, err
				}

				_ = trans.SetTimeout(timeout)

				if err = trans.Open(); err != nil {
					return nil, errors.New("HBase open error: " + err.Error())
				}

				f := thrift.NewTBinaryProtocolFactoryDefault()
				std := thrift.NewTStandardClient(f.GetProtocol(trans), f.GetProtocol(trans))
				c := NewHbaseClient(std)

				h := &HBasePoolClient{
					Client: c,
					Trans:  trans,
				}
				zlog.DebugLogger(nil, "open a new connection at: "+time.Now().String())
				return h, nil
			}

			closeConn := func(client interface{}) error {
				if client != nil {
					zlog.DebugLogger(nil, "close connection at: "+time.Now().String())
					c, ok := client.(*HBasePoolClient)
					if ok && c != nil && c.Trans != nil {
						return c.Trans.Close()
					}
				}
				return nil
			}

			poolConfig := &pool.Config{
				Factory:    createConn,
				Close:      closeConn,
				Ping:       nil,
				InitialCap: res.HBase.MaxIdle,   // 资源池初始连接数
				MaxIdle:    res.HBase.MaxIdle,   // 最大空闲连接数
				MaxCap:     res.HBase.MaxActive, // 最大并发连接数
				// 连接最大空闲时间，超过该时间的连接 将会关闭，可避免空闲时连接EOF，自动失效的问题
				IdleTimeout: res.HBase.IdleTimeout,
				WaitTimeOut: res.HBase.MaxWaitTimeout,
			}

			p, err := pool.NewChannelPool(poolConfig)
			if err != nil {
				zlog.WarnLogger(nil, "init pool error: "+err.Error())
				return nil, err
			}

			client := &HbaseClientModule{
				pool:    p,
				Service: res.Name,
			}
			sub := &ral.Instance{
				IP:     ins.IP,
				Port:   ins.Port,
				Client: client,
			}
			client.ins = sub
			client.res = res
			return sub, nil
		},
		Remove: func(self *ral.Instance, res *ral.Resource, ins *ral.Instance) bool {
			// 关闭当前实例的连接池
			if r, ok := ins.Client.(*HbaseClientModule); ok && r.pool != nil {
				r.pool.Release()
			}
			zlog.DebugLogger(nil, "release connections at: "+time.Now().String())
			return true
		},
		Method: func(ctx *gin.Context, res *ral.Resource, ins *ral.Instance, method string, data map[string]interface{}, head map[string]string) (interface{}, error) {
			return nil, nil
		}}
	mod = ral.AddModule(m)
}

func GetInstance(ctx *gin.Context, service string) (*HbaseClientModule, error) {
	ins, err := ral.GetInstance(ctx, ral.TYPE_HBASE, service)
	if err != nil {
		zlog.ErrorLogger(ctx, "module hbase not found clients", zap.String("prot", "hbase"))
		return nil, ral.ERR_NOT_FOUND_INSTANCE
	}
	p, ok := ins.Client.(*HbaseClientModule)
	if !ok {
		zlog.ErrorLogger(ctx, "module hbase get clients error", zap.String("prot", "hbase"))
		return nil, ral.ERR_NOT_FOUND_CLIENT
	}
	return p, nil
}
