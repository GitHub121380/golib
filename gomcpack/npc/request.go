package npc

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math/rand"
	"strconv"
	"strings"

	"github.com/GitHub121380/golib/zlog"
	"github.com/gin-gonic/gin"
)

type Request struct {
	Header     Header
	Body       io.Reader
	RemoteAddr string
}

func (r *Request) Write(w io.Writer) (n int, err error) {
	n, err = r.Header.Write(w)
	if err != nil {
		return 0, err
	}
	written, err := io.Copy(w, r.Body)
	return int(written), err
}

func ReadRequest(r io.Reader) (req *Request, err error) {
	req = new(Request)
	_, err = req.Header.Read(r)
	if err != nil {
		return nil, err
	}
	if req.Header.MagicNum != HEADER_MAGICNUM {
		return nil, fmt.Errorf("invalid magic number %x", req.Header.MagicNum)
	}
	req.Body = io.LimitReader(r, int64(req.Header.BodyLen))
	return req, nil
}

func NewRequest(ctx context.Context, body io.Reader) *Request {
	req := new(Request)
	// 优先从ctx 中获取logID，否则按照原来逻辑随机生成
	var logID uint32
	if c, ok := ctx.(*gin.Context); ok && c != nil {
		if ids, ok := ctx.Value(zlog.ContextKeyLogID).(string); ok {
			if id, err := strconv.Atoi(ids); err == nil {
				logID = uint32(id)
			}
		}
	}

	if logID == 0 {
		logID = rand.Uint32()
	}

	req.Header.LogId = logID
	req.Header.MagicNum = HEADER_MAGICNUM
	if body != nil {
		switch v := body.(type) {
		case *bytes.Buffer:
			req.Header.BodyLen = uint32(v.Len())
			req.Body = io.LimitReader(body, int64(req.Header.BodyLen))
		case *bytes.Reader:
			req.Header.BodyLen = uint32(v.Len())
			req.Body = io.LimitReader(body, int64(req.Header.BodyLen))
		case *strings.Reader:
			req.Header.BodyLen = uint32(v.Len())
			req.Body = io.LimitReader(body, int64(req.Header.BodyLen))
		default:
			panic("unsupported io.Reader")
		}
	}
	return req
}
