package cos

import (
	"context"
	"fmt"
	"github.com/GitHub121380/golib/utils"
	"github.com/GitHub121380/golib/zlog"
	"go.uber.org/zap"
	"net/http"
	"reflect"
	"strconv"
	"time"
)

type cosRequestTransport struct {
	Transport http.RoundTripper
}

// 注意此函数只提取string类型的key
func FromContext(ctx context.Context, key string) string {
	rv := reflect.ValueOf(ctx)
	if !rv.IsNil() {
		u, ok := ctx.Value(key).(string)
		if ok {
			return u
		}
	}

	return ""
}

// TODO: add spanID to log
func CreateSpan(ctx context.Context) string {
	parentSpanID := FromContext(ctx, "parentSpanID")

	childSpanID := 0
	if s := FromContext(ctx, "childSpanID"); s != "" {
		var err error
		childSpanID, err = strconv.Atoi(s)
		if err != nil {
			return ""
		}
	}

	childSpanID++
	ctx = context.WithValue(ctx, "childSpanID", childSpanID)

	newSpan := fmt.Sprintf("%s.%d", parentSpanID, childSpanID)
	return newSpan
}

func (t *cosRequestTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req = cloneRequest(req)
	start := time.Now()
	resp, err := t.transport().RoundTrip(req)
	end := time.Now()

	httpCode := 0
	if resp != nil {
		httpCode = resp.StatusCode
	}
	ralCode := 0
	msg := "cos request success"
	if err != nil {
		ralCode = -1
		msg = err.Error()
		return resp, err
	}

	fields := []zap.Field{
		zap.String("prot", "cos"),
		zap.String("logId", FromContext(req.Context(), zlog.ContextKeyLogID)),
		zap.String("spanId", FromContext(req.Context(), zlog.ContextKeySpanID)),
		zap.String("requestId", FromContext(req.Context(), zlog.ContextKeyRequestID)),
		zap.String("method", req.Method),
		zap.String("domain", req.URL.Host),
		zap.String("requestUri", req.URL.Path),
		zap.String("requestStartTime", utils.GetFormatRequestTime(start)),
		zap.String("requestEndTime", utils.GetFormatRequestTime(end)),
		zap.Float64("cost", utils.GetRequestCost(start, end)),
		zap.Int("httpCode", httpCode),
		zap.Int("ralCode", ralCode),
	}

	zlog.InfoLogger(nil, msg, fields...)
	return resp, err
}

func (t *cosRequestTransport) transport() http.RoundTripper {
	if t.Transport != nil {
		return t.Transport
	}
	return http.DefaultTransport
}

// cloneRequest returns a clone of the provided *http.Request. The clone is a
// shallow copy of the struct and its Header map.
func cloneRequest(r *http.Request) *http.Request {
	// shallow copy of the struct
	r2 := new(http.Request)
	*r2 = *r
	// deep copy of the Header
	r2.Header = make(http.Header, len(r.Header))
	for k, s := range r.Header {
		r2.Header[k] = append([]string(nil), s...)
	}
	return r2
}
