package bns

import (
	"github.com/golang/protobuf/proto"
	"time"
)

//go:generate protoc --go_out=. naming.proto naminglib.proto service.proto

type MsgType int

const (
	ReqService         MsgType = 1
	ReqAuthService     MsgType = 2
	ResService         MsgType = 3
	ResAuthService     MsgType = 4
	ReqServiceList     MsgType = 9
	ReqAuthServiceList MsgType = 10
	ResServiceList     MsgType = 11
	ResAuthServiceList MsgType = 12
	ReqServiceConf     MsgType = 15
	ResServiceConf     MsgType = 16
	ReqRealService     MsgType = 30
)

const (
	defaultRetryTimes = 2
)

type Client struct {
	// socket read/write timeout.
	// if zero, DefaultTimeout is used
	Timeout time.Duration

	local  *client
	remote *client
	conn   *clientConn
}

func New(addr string, timeout time.Duration) *Client {
	return &Client{
		Timeout: timeout,
		local:   &client{Timeout: timeout, Addr: addr},
		remote:  &client{Timeout: timeout, Addr: addr},
	}
}

func (c *Client) Connect() error {
	n, e := c.local.getConn(c.local.Addr)
	c.conn = n
	return e
}
func (c *Client) Close() {
	c.conn.close()
}

func (c *Client) Call(args proto.Message, reply proto.Message) error {
	var reqType MsgType
	switch args.(type) {
	case *LocalNamingRequest:
		reqType = ReqService
	case *LocalNamingAuthRequest:
		reqType = ReqAuthService
	case *LocalServiceConfRequest:
		reqType = ReqServiceConf
	case *LocalNamingListRequest:
		reqType = ReqServiceList
	case *LocalNamingAuthListRequest:
		reqType = ReqAuthServiceList
	}

	content, err := proto.Marshal(args)
	if err != nil {
		return err
	}

	req := newRequest(reqType, content)

	var resp *response
	resp, err = c.local.do(req)
	if err != nil {
		resp, err = c.remote.doWithRetry(req, defaultRetryTimes)
		if err != nil {
			return err
		}
	}

	return proto.Unmarshal(resp.Body, reply)
}
