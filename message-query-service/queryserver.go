package main

import (
	"context"
	"github.com/yousuf64/request-coalescing/gen/query"
	"github.com/yousuf64/request-coalescing/message-query-service/workrunner"
	"github.com/yousuf64/request-coalescing/pkg/requestctx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"time"
)

type queryServer struct {
	query.UnimplementedQueryServer

	Messages *workrunner.WorkRunner[string, []*query.Message]
}

func (svr *queryServer) GetMessages(ctx context.Context, request *query.GetMessagesRequest) (*query.GetMessagesReply, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "Missing x-route-key header")
	}

	routeKeyValues, ok := md["x-route-key"]
	if !ok || len(routeKeyValues) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Missing x-route-key header")
	}

	response := <-svr.Messages.Observe(routeKeyValues[0], func() []*query.Message {
		log.Printf("[INFO] <ID: %s, ChannelId: %d> Executing GetMessages query", requestctx.FromContext(ctx).String(), request.ChannelId)
		time.Sleep(5 * time.Second)
		return []*query.Message{
			{
				MessageId: 222,
				UserId:    912,
				Message:   "hello",
				Timestamp: 0,
			},
			{
				MessageId: 225,
				UserId:    911,
				Message:   "sup?",
				Timestamp: 0,
			},
		}
	})

	return &query.GetMessagesReply{Messages: response}, nil
}
