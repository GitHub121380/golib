package grpc

import (
	"context"
	"go.uber.org/zap"
	"time"

	"github.com/GitHub121380/golib/zlog"
	"google.golang.org/grpc"
)

// access日志打印
func AccessLog() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// 开始时间
		start := time.Now()

		resp, err := handler(ctx, req)
		// 结束时间
		end := time.Now()
		// 执行时间 单位:微秒
		cost := end.Sub(start).Nanoseconds() / 1e3

		// 用户自定义notice
		var customerFields []zap.Field
		commonFields := []zap.Field{
			zap.String("logId", "logID"),
			zap.String("method", info.FullMethod),
			zap.Reflect("request", req),
			zap.Int64("cost", cost),
		}

		zlog.GetAccessLogger().With(commonFields...).Info("notice", customerFields...)

		return resp, err
	}
}
