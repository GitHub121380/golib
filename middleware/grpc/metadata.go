package grpc

import (
	"context"

	"github.com/GitHub121380/golib/utils/metadata"
	"google.golang.org/grpc"
)

func Metadata() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md := metadata.MD(map[string]interface{}{metadata.Notice: make(map[string]interface{})})
		ctx = metadata.NewContext(ctx, md)
		return handler(ctx, req)
	}
}
