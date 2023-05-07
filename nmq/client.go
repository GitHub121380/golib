package nmq

import (
	"github.com/GitHub121380/golib/gomcpack/mcpacknpc"
	"github.com/GitHub121380/golib/ral"
	"github.com/GitHub121380/golib/utils"
	"github.com/gin-gonic/gin"

	"time"
)

const (
	HEAD_PRODUCT = "_product"
	HEAD_TOPIC   = "_topic"
	HEAD_CMD     = "_cmd"

	// 编码方式使用的配置文件中的 Encode 字段，为防止歧义，不要使用该字段指定了
	HEAD_ENCODE = "_encode"
	HEAD_UAGENT = "User-Agent"
)

const (
	METHOD_LONG  = "LONG"
	METHOD_SHORT = "SHORT"
)

const (
	ENCODE_MCPACK1 = "mcpack1"
	ENCODE_MCPACK2 = "mcpack2"
)

type Request struct {
	Method string
	Header map[string]string
	Data   map[string]interface{}

	EncodeType string
	EncodeData interface{}

	Meta map[string]interface{}

	BeginTime time.Time
	ReplyTime time.Time
}

func NewRequest(method string, header map[string]string, data map[string]interface{}, res *ral.Resource) *Request {
	return &Request{
		Method:     method,
		Header:     header,
		Data:       data,
		Meta:       map[string]interface{}{},
		EncodeType: res.Encode,
	}
}

type Client struct {
	IP   string
	Port int
	Host string

	ConnTimeOut  time.Duration
	ReadTimeOut  time.Duration
	WriteTimeOut time.Duration
	RetryCount   int

	agent *mcpacknpc.Client
}

func (client *Client) Check(ctx *gin.Context, req *Request) error {
	req.Header[HEAD_UAGENT] = ral.AppInfo.Agent

	if req.EncodeType == "" {
		return ral.ERR_NOT_FOUND_ENCODE
	}

	data := make(map[string]interface{})
	for k, v := range req.Data {
		switch k {
		case ral.DATA_CONTEXT:
		default:
			data[k] = v
		}
	}

	// 全链路压测标记透传
	data[utils.HttpUrlPressureCallerKey], data[utils.HttpUrlPressureMarkKey] = utils.GetPressureFlag(ctx)

	if res, err := Encode(req.EncodeType, data); err != nil {
		return err
	} else {
		req.EncodeData = res
	}

	req.BeginTime = time.Now()
	return nil
}

func (client *Client) After(req *Request, buf []byte, err error, res interface{}) error {
	if err != nil {
		return err
	}
	req.ReplyTime = time.Now()
	req.Meta["cost"] = req.ReplyTime.Sub(req.BeginTime)
	return Decode(req.EncodeType, buf, res)
}

func NewClient(ip string, port int, retryCount int, connTimeOut time.Duration, readTimeOut time.Duration, writeTimeOut time.Duration) *Client {
	return &Client{
		IP:           ip,
		Port:         port,
		RetryCount:   retryCount,
		ConnTimeOut:  connTimeOut,
		ReadTimeOut:  readTimeOut,
		WriteTimeOut: writeTimeOut,
	}
}
