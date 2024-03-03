package main

import (
	"fmt"
	"github.com/hashicorp/logutils"
	"github.com/yousuf64/request-coalescing/gen/messages"
	"github.com/yousuf64/request-coalescing/gen/query"
	"github.com/yousuf64/request-coalescing/pkg/configloader"
	"github.com/yousuf64/request-coalescing/pkg/grpcinterceptors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"os"
)

type Config struct {
	App struct {
		Hostname string `yaml:"hostname"`
		Port     int    `yaml:"port"`
	} `yaml:"app"`
	External struct {
		MessageQueryServiceAddr string `yaml:"message_query_service_addr"`
	} `yaml:"external"`
}

func main() {
	conf := configloader.Load[Config]("env", "config/env.yaml")

	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "WARN", "ERROR", "INFO"},
		MinLevel: "DEBUG",
		Writer:   os.Stdout,
	}
	log.SetOutput(filter)

	var opts = []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.Dial(conf.External.MessageQueryServiceAddr, opts...)
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	queryClient := query.NewQueryClient(conn)

	messagesSvr := &MessagesServer{
		QueryClient: queryClient,
	}

	addr := fmt.Sprintf("%s:%d", conf.App.Hostname, conf.App.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("listening on %s\n", addr)

	grpcSvr := grpc.NewServer(grpc.UnaryInterceptor(grpcinterceptors.RequestId))
	messages.RegisterMessagesServer(grpcSvr, messagesSvr)

	err = grpcSvr.Serve(lis)
	if err != nil {
		log.Fatal(err)
	}
}
