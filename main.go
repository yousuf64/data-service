package main

import (
	"context"
	"data-service/query"
	"data-service/requestctx"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	querySvr := queryServer{
		getMessagesOrch:    NewOrchestrator[[]*query.Message](),
		getMessagesWorkers: nil,
	}

	addr := fmt.Sprintf("%s:%d", "localhost", 3000)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("listening on %s\n", addr)

	interceptor := grpc.UnaryInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		ctx = requestctx.New(ctx, uuid.New())
		return handler(ctx, req)
	})

	grpcSvr := grpc.NewServer(interceptor)
	query.RegisterQueryServer(grpcSvr, querySvr)

	err = grpcSvr.Serve(lis)
	if err != nil {
		log.Fatal(err)
	}
}
