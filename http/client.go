package http

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/GitHub121380/golib/ral"
	"github.com/GitHub121380/golib/utils"
	"github.com/gin-gonic/gin"
	"github.com/parnurzeal/gorequest"
)

const (
	HEAD_PATH = "_path"
	// 编码方式使用的配置文件中的 Encode 字段，为防止歧义，不要使用该字段指定了
	HEAD_ENCODE = "_encode"
	HEAD_UAGENT = "User-Agent"
	// 控制 GoRequest 的 BounceToRawString 字段，不为空则表表开启
	HEAD_PARAM_BOUNCE = "_bounce"
	// 返回的数据不进行解包
	HEAD_NO_UNPACK = "_no_unpack"
)

const (
	METHOD_POST = "POST"
	METHOD_GET  = "GET"
)

const (
	ENCODE_TEXT    = "text"
	ENCODE_FORM    = "form"
	ENCODE_JSON    = "json"
	ENCODE_MCPACK1 = "mcpack1"
	ENCODE_MCPACK2 = "mcpack2"
)

type Request struct {
	Method string
	Header map[string]string
	Data   map[string]interface{}

	URI          string
	EncodeType   string
	EncodeData   interface{}
	BounceEncode bool
	NoDecode     bool

	Meta map[string]interface{}

	StatusCode int
	Status     string

	BeginTime time.Time
	ReplyTime time.Time
	Cost      time.Duration
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
	Port int
	IP   string
	Host string

	ConnTimeOut  time.Duration
	ReadTimeOut  time.Duration
	WriteTimeOut time.Duration
	RetryCount   int
}

func (client *Client) Check(ctx *gin.Context, req *Request) error {
	callerURI, _ := utils.GetPressureFlag(ctx)
	if callerURI != "" {
		req.Header[utils.HttpXBDCallerURI] = callerURI
	}
	req.Header[HEAD_UAGENT] = ral.AppInfo.Agent

	if p, ok := req.Header[HEAD_PATH]; !ok || p == "" {
		return ral.ERR_NOT_FOUND_PATH
	} else {
		if !strings.HasPrefix(p, "http") {
			req.URI = fmt.Sprintf("http://%s:%d%s", client.IP, client.Port, p)
		} else {
			req.URI = p
		}
	}

	if p, ok := req.Header[HEAD_PARAM_BOUNCE]; ok && p != "" {
		req.BounceEncode = true
		delete(req.Header, HEAD_PARAM_BOUNCE)
	}

	if p, ok := req.Header[HEAD_NO_UNPACK]; ok && p != "" {
		req.NoDecode = true
		delete(req.Header, HEAD_NO_UNPACK)
	}

	if req.EncodeType == "" {
		return ral.ERR_NOT_FOUND_ENCODE
	}

	data := map[string]interface{}{}
	for k, v := range req.Data {
		switch k {
		case ral.DATA_CONTEXT:
		default:
			data[k] = v
		}
	}

	// get 参数只需要拼装,不能做其他类型的格式转换
	encodeType := req.EncodeType
	if req.Method == METHOD_GET {
		encodeType = ENCODE_FORM
	}

	if res, err := Encode(encodeType, data); err != nil {
		return err
	} else {
		req.EncodeData = res
	}

	req.BeginTime = time.Now()
	return nil
}

func (client *Client) Send(req *Request) ([]byte, error) {
	goreq := gorequest.New()
	goreq = goreq.Timeout(client.ReadTimeOut)
	goreq.ClearSuperAgent()

	uri, query := req.URI, ""
	if list := strings.SplitN(req.URI, "?", 2); len(list) > 1 {
		uri, query = list[0], list[1]
	}

	if req.EncodeType == ENCODE_MCPACK1 || req.EncodeType == ENCODE_MCPACK2 {
		goreq.BounceToRawString = true
		body, ok := req.EncodeData.(string)
		if !ok {
			return nil, errors.New("mcpak date error!")
		}
		goreq = goreq.Post(uri).Query(query).Type(gorequest.TypeText).SendString(body)
	} else {
		switch req.Method {
		case METHOD_GET:
			goreq = goreq.Get(uri).Query(query).Query(req.EncodeData)
		case METHOD_POST:
			goreq = goreq.Post(uri).Query(query).Type(req.EncodeType).Send(req.EncodeData)
		}

		if req.BounceEncode {
			goreq.BounceToRawString = true
		}
	}

	for k, v := range req.Header {
		goreq.Set(k, v)
	}

	resp, data, errs := goreq.EndBytes()
	if len(errs) > 0 {
		return data, errs[0]
	}

	req.StatusCode = resp.StatusCode
	req.Status = resp.Status
	return data, nil
}
func (client *Client) After(req *Request, buf []byte) ([]byte, error) {
	req.ReplyTime = time.Now()
	req.Cost = req.ReplyTime.Sub(req.BeginTime)

	return Decode(req, buf)
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
