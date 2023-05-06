package main

import (
	"context"
	"data-service/query"
	"data-service/requestctx"
	"data-service/workrunner"
	"log"
	"time"
)

type queryServer struct {
	query.UnimplementedQueryServer

	Messages *workrunner.WorkRunner[int32, []*query.Message]
}

func (svr *queryServer) GetMessages(ctx context.Context, request *query.GetMessagesRequest) (*query.GetMessagesReply, error) {
	response := <-svr.Messages.Observe(request.ChannelId, func() []*query.Message {
		log.Printf("[INFO] <ID: %s, ChannelId: %d> Executing GetMessages query", requestctx.FromContext(ctx).String(), request.ChannelId)
		time.Sleep(5 * time.Second)
		return []*query.Message{
			{
				FromId:  911,
				Message: "hello",
			},
			{
				FromId:  912,
				Message: "sup?",
			},
		}
	})

	return &query.GetMessagesReply{Messages: response}, nil
}
