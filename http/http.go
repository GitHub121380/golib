package http

import (
	"github.com/GitHub121380/golib/env"
	"github.com/GitHub121380/golib/ral"
	"github.com/GitHub121380/golib/utils"
	"github.com/GitHub121380/golib/zlog"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"fmt"
	"time"
)

var mod *ral.Module

func init() {
	m := &ral.Module{
		Type: ral.TYPE_HTTP,
		Config: func(res *ral.Resource) *ral.Resource {
			if res.ReadTimeOut <= 0 {
				res.ReadTimeOut = 3 * time.Second
			} else if res.ReadTimeOut > 20*time.Second {
				res.ReadTimeOut = 20 * time.Second
			}
			return res

		},
		Append: func(self *ral.Resource, res *ral.Resource, ins *ral.Instance) (*ral.Instance, error) {
			i := &ral.Instance{
				IP:   ins.IP,
				Port: ins.Port,
				Client: &Client{
					IP:           ins.IP,
					Port:         ins.Port,
					Host:         fmt.Sprintf("%s:%d", ins.IP, ins.Port),
					ConnTimeOut:  res.ConnTimeOut,
					ReadTimeOut:  res.ReadTimeOut,
					WriteTimeOut: res.WriteTimeOut,
				}}

			return i, nil

		},
		Remove: func(self *ral.Instance, res *ral.Resource, ins *ral.Instance) bool {
			return true
		},
		Method: func(ctx *gin.Context, res *ral.Resource, ins *ral.Instance, method string,
			data map[string]interface{}, head map[string]string) (interface{}, error) {
			start := time.Now()
			fields := []zap.Field{
				zap.String("prot", res.Type),
				zap.String("method", method),
				zap.String("service", res.Name),
				zap.String("conv", res.Encode),
				zap.String("requestStartTime", utils.GetFormatRequestTime(start)),
				zap.String("req_uri", head[HEAD_PATH]),
			}
			client, ok := ins.Client.(*Client)
			if !ok {
				zlog.WarnLogger(ctx, "instance type error", fields...)
				return nil, ral.ERR_NOT_FOUND_CLIENT
			}
			req := NewRequest(method, head, data, res)
			if err := client.Check(ctx, req); err != nil {
				zlog.WarnLogger(ctx, "req_error: "+err.Error(), fields...)
				return nil, err
			}

			buf, err := client.Send(req)

			fields = append(fields,
				zap.String("remoteIp", client.Host),
				zap.String("prot_status", req.Status),
				zap.Int("prot_code", req.StatusCode),
			)

			if err != nil {
				zlog.WarnLogger(ctx, "res_error: "+err.Error(), fields...)
			} else if buf, err = client.After(req, buf); err != nil {
				zlog.WarnLogger(ctx, "decode_error: "+err.Error(), fields...)
			}

			zlog.DebugLogger(ctx, string(buf), fields...)

			// 结束时间
			end := time.Now()

			fields = append(fields,
				zap.Int("res_len", len(buf)),
				zap.String("requestEndTime", utils.GetFormatRequestTime(end)),
				zap.Float64("cost", utils.GetRequestCost(start, end)),
			)

			zlog.InfoLogger(ctx, "http end", fields...)

			return buf, err
		}}

	mod = ral.AddModule(m)
}

func Call(ctx *gin.Context, method string, service string, data map[string]interface{},
	head map[string]string) (buf []byte, err error) {
	ins, err := ral.GetInstance(ctx, ral.TYPE_HTTP, service)
	if err != nil {
		zlog.WarnLogger(ctx, "GetInstance error: "+err.Error(), zap.String("service", service), zap.String("prot", "http"))
		return nil, err
	}
	defer ins.Release()

	ins.Retry(func(res *ral.Resource, ins *ral.Instance) bool {
		if b, e := ins.Request(ctx, method, data, head); e == nil {
			buf, _ = b.([]byte)
			return false
		} else {
			err = e
		}
		return true
	})

	if err != nil {
		field := []zap.Field{
			zap.String("method", method),
			zap.String("service", service),
			zap.String("remoteIp", fmt.Sprintf("%s:%d", ins.IP, ins.Port)),
			zap.String("module", env.GetAppName()),
			zap.String("error", err.Error()),
		}
		zlog.WarnLogger(ctx, "call failed", field...)
	}

	return buf, err
}
func Post(ctx *gin.Context, service string, data map[string]interface{}, head map[string]string) ([]byte, error) {
	if ctx != nil {
		if data != nil {
			data[ral.DATA_CONTEXT] = ctx
		} else {
			data = map[string]interface{}{
				ral.DATA_CONTEXT: ctx,
			}
		}
	}
	if head == nil {
		head = make(map[string]string)
	}

	return Call(ctx, METHOD_POST, service, data, head)
}
func Get(ctx *gin.Context, service string, data map[string]interface{}, head map[string]string) ([]byte, error) {
	if ctx != nil {
		if data != nil {
			data[ral.DATA_CONTEXT] = ctx
		} else {
			data = map[string]interface{}{
				ral.DATA_CONTEXT: ctx,
			}
		}
	}

	if head == nil {
		head = make(map[string]string)
	}
	return Call(ctx, METHOD_GET, service, data, head)
}
