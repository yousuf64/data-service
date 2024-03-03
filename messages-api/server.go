package main

import (
	"fmt"
	"github.com/yousuf64/request-coalescing/gen/messages"
	"github.com/yousuf64/request-coalescing/gen/query"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type MessagesServer struct {
	messages.UnimplementedMessagesServer

	QueryClient query.QueryClient
}

func (ms *MessagesServer) GetMessages(ctx context.Context, request *messages.GetMessagesRequest) (*messages.GetMessagesReply, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.New(map[string]string{})
	}
	md.Set("x-route-key", fmt.Sprintf("%d_%d", request.ChannelId, request.LastMessageId))
	ctx = metadata.NewOutgoingContext(ctx, md)

	reply, err := ms.QueryClient.GetMessages(ctx, &query.GetMessagesRequest{
		ChannelId: request.ChannelId,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	msgs := make([]*messages.Message, len(reply.Messages))
	for i, message := range reply.Messages {
		msgs[i] = &messages.Message{
			MessageId: message.MessageId,
			UserId:    message.UserId,
			Message:   message.Message,
			Timestamp: message.Timestamp,
		}
	}

	return &messages.GetMessagesReply{Messages: msgs}, nil
}

func (ms *MessagesServer) CreateMessage(ctx context.Context, request *messages.CreateMessageRequest) (*messages.CreateMessageReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateMessage not implemented")
}
