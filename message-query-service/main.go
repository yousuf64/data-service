package main

import (
	"embed"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/hashicorp/logutils"
	"github.com/yousuf64/request-coalescing/gen/query"
	"github.com/yousuf64/request-coalescing/message-query-service/workrunner"
	"github.com/yousuf64/request-coalescing/pkg/configloader"
	"github.com/yousuf64/request-coalescing/pkg/grpcinterceptors"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

//go:embed config/*
var configFiles embed.FS

type Config struct {
	App struct {
		Hostname string `yaml:"hostname"`
		Port     int    `yaml:"port"`
	} `yaml:"app"`
}

func main() {
	conf := configloader.Load[Config]("env", "env.yaml")
	cc := gocql.NewCluster("127.0.0.1:19042")
	fallback := gocql.RoundRobinHostPolicy()
	cc.Polpolicy := gocql.TokenAwareHostPolicy(fallback)
	consistency := gocql.LocalQuorum
	s, _ := gocql.NewSession(*cc)
	s.Query("")
	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "WARN", "ERROR", "INFO"},
		MinLevel: "DEBUG",
		Writer:   os.Stdout,
	}
	log.SetOutput(filter)

	querySvr := &queryServer{
		Messages: workrunner.New[string, []*query.Message](),
	}

	addr := fmt.Sprintf("%s:%d", conf.App.Hostname, conf.App.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("listening on %s\n", addr)

	grpcSvr := grpc.NewServer(grpc.UnaryInterceptor(grpcinterceptors.RequestId))
	query.RegisterQueryServer(grpcSvr, querySvr)

	err = grpcSvr.Serve(lis)
	if err != nil {
		log.Fatal(err)
	}
}
