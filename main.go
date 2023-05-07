package main

import (
	"context"
	"data-service/query"
	"data-service/requestctx"
	"data-service/workrunner"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/logutils"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

func main() {
	port := os.Args[1]

	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "WARN", "ERROR", "INFO"},
		MinLevel: "DEBUG",
		Writer:   os.Stdout,
	}
	log.SetOutput(filter)

	querySvr := &queryServer{
		Messages: workrunner.New[int32, []*query.Message](),
	}

	addr := fmt.Sprintf("%s:%s", "127.0.0.1", port)
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
