package nmq

import (
	"fmt"
	"strconv"
	"time"

	"github.com/GitHub121380/golib/gomcpack/mcpacknpc"
	"github.com/GitHub121380/golib/ral"
	"github.com/GitHub121380/golib/zlog"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type NmqResponse struct {
	ErrNo  int    `mcpack:"_error_no" binding:"required"`
	ErrStr string `mcpack:"_error_msg" binding:"required"`
}

var mod *ral.Module

func init() {
	mod = ral.AddModule(&ral.Module{Type: ral.TYPE_NMQ, Config: func(res *ral.Resource) *ral.Resource {
		// 检查配置
		if res.ReadTimeOut <= 0 {
			res.ReadTimeOut = 3 * time.Second
		} else if res.ReadTimeOut > 20*time.Second {
			res.ReadTimeOut = 20 * time.Second
		}
		return res

	}, Append: func(self *ral.Resource, res *ral.Resource, ins *ral.Instance) (*ral.Instance, error) {
		// 创建实例
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

	}, Remove: func(self *ral.Instance, res *ral.Resource, ins *ral.Instance) bool {
		// 删除实例
		return true

	}, Method: func(ctx *gin.Context, res *ral.Resource, ins *ral.Instance, method string, data map[string]interface{}, head map[string]string) (buf interface{}, err error) {
		// 访问资源
		client, ok := ins.Client.(*Client)

		fields := []zap.Field{
			zap.String("prot", "nmq"),
			zap.String("method", method),
			zap.String("service", res.Name),
			zap.String("remoteIp", client.Host),
		}

		if !ok {
			zlog.WarnLogger(ctx, "instance type error", fields...)
			return nil, ral.ERR_NOT_FOUND_CLIENT
		}

		// 创建连接
		agent := client.agent
		if agent == nil || !res.LongConnect {
			agent = mcpacknpc.NewClient([]string{client.Host})
			agent.Timeout = client.ReadTimeOut
		}

		// 保存连接
		if !res.LongConnect {
			defer agent.Close()
		} else if client.agent == nil {
			client.agent = agent
		}

		fields = append(fields,
			zap.String("conv", res.Encode),
			zap.Reflect("topic", data[HEAD_TOPIC]),
			zap.Reflect("cmdno", data[HEAD_CMD]),
		)

		// 构造请求
		req := NewRequest(method, head, data, res)
		if err := client.Check(ctx, req); err != nil {
			zlog.WarnLogger(ctx, "check error"+err.Error(), fields...)
			return nil, err
		}

		// 发送请求
		enData, _ := req.EncodeData.([]byte)
		result, err := agent.Send(ctx, enData)

		// 处理响应
		var resp NmqResponse
		if err = client.After(req, result, err, &resp); err == nil {
			fields = append(fields,
				zap.Reflect("req_data", data),
				zap.Int("err_no", resp.ErrNo),
				zap.String("err_info", resp.ErrStr),
			)

			zlog.InfoLogger(ctx, "call succ", fields...)
		} else {
			fields = append(fields, zap.Reflect("req_data", data))
			zlog.WarnLogger(ctx, "call failed"+err.Error(), fields...)
		}
		return &resp, err
	}})
}

func Call(ctx *gin.Context, method string, service string, data map[string]interface{}, head map[string]string) (resp *NmqResponse, err error) {
	ins, err := ral.GetInstance(ctx, ral.TYPE_NMQ, service)
	if err != nil {
		zlog.WarnLogger(ctx, err.Error(), zap.String("service", service), zap.String("prot", "nmq"))
		return nil, err
	}
	defer ins.Release()

	ins.Retry(func(res *ral.Resource, ins *ral.Instance) bool {
		if r, e := ins.Request(ctx, method, data, head); e == nil {
			resp, _ = r.(*NmqResponse)
			return false
		} else {
			err = e
		}
		return true
	})

	return resp, err
}
func SendCmd(ctx *gin.Context, service string, cmd int64, topic string, product string, data map[string]interface{}, head map[string]string) (*NmqResponse, error) {
	data[HEAD_CMD] = strconv.FormatInt(cmd, 10)
	data[HEAD_TOPIC] = topic
	data[HEAD_PRODUCT] = product

	return Call(ctx, METHOD_SHORT, service, data, head)
}
func Send(ctx *gin.Context, service string, data map[string]interface{}, head map[string]string) (*NmqResponse, error) {
	return Call(ctx, METHOD_SHORT, service, data, head)
}
func SendByLongConnect(ctx *gin.Context, service string, data map[string]interface{}, head map[string]string) (*NmqResponse, error) {
	return Call(ctx, METHOD_LONG, service, data, head)
}
