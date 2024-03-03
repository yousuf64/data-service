package grpcinterceptors

import (
	"github.com/google/uuid"
	"github.com/yousuf64/request-coalescing/pkg/requestctx"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func RequestId(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	var id uuid.UUID
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if idstr, ok := md["x-request-id"]; ok {
			id, _ = uuid.Parse(idstr[0])
		}
	}

	if id == uuid.Nil {
		id = uuid.New()
	}

	ctx = requestctx.New(ctx, id)
	return handler(ctx, req)
}
