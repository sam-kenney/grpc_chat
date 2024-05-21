package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	pb "github.com/sam-kenney/grpc_chat/chat"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

type ChatServer struct {
	pb.UnimplementedChatServer
	channels []*pb.Channel
	messages []*pb.Message
	mu       sync.Mutex // lock when mutating channels or messages
}

func NewChatServer() *ChatServer {
	return &ChatServer{}
}

func (s *ChatServer) CreateChannel(ctx context.Context, channel *pb.Channel) (*pb.Channel, error) {
	for _, ch := range s.channels {
		if proto.Equal(ch, channel) {
			return nil, status.Errorf(codes.AlreadyExists, channel.Name)
		}
	}

	s.mu.Lock()
	s.channels = append(s.channels, channel)
	s.mu.Unlock()

	return channel, nil
}

func (s *ChatServer) ListChannels(ctx context.Context, req *pb.ListChannelsRequest) (*pb.Channels, error) {
	return &pb.Channels{
		Channels: s.channels,
	}, nil
}

func (s *ChatServer) SendMessage(ctx context.Context, message *pb.Message) (*pb.Message, error) {
	s.mu.Lock()
	s.messages = append(s.messages, message)
	s.mu.Unlock()

	return message, nil
}

func (s *ChatServer) ListMessages(ctx context.Context, req *pb.ListMessagesRequest) (*pb.Messages, error) {
	var messages []*pb.Message
	ch := req.GetChannel()

	for _, m := range s.messages {
		if m.Channel.Name == ch.GetName() {
			messages = append(messages, m)
		}
	}

	return &pb.Messages{
		Messages: messages,
	}, nil
}

func main() {
	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 8080))
	if err != nil {
		log.Fatalf("Error: %s", err.Error())
	}
	defer listener.Close()

	server := grpc.NewServer()
	defer server.Stop()
	pb.RegisterChatServer(server, NewChatServer())

	log.Printf("Listening at http://%s", listener.Addr())
	server.Serve(listener)
}
