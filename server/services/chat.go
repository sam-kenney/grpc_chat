package services

import (
	"context"
	"fmt"
	"sync"

	pb "github.com/sam-kenney/grpc_chat/chat"
	"github.com/sam-kenney/logging"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

type ChatServer struct {
	pb.UnimplementedChatServer
	channels []*pb.Channel
	messages []*pb.Message
	mu       sync.Mutex // lock when mutating channels or messages
	log      logging.Logger
}

func NewChatServer(log logging.Logger) *ChatServer {
	return &ChatServer{log: log}
}

func (s *ChatServer) CreateChannel(ctx context.Context, channel *pb.Channel) (*pb.Channel, error) {
	s.log.Info(fmt.Sprintf("Chat.CreateChannel: Creating channel '%s'", channel.Name))

	for _, ch := range s.channels {
		if proto.Equal(ch, channel) {
			s.log.Warning(fmt.Sprintf("Chat.CreateChannel: Channel '%s' already exists", channel.Name))
			return nil, status.Errorf(codes.AlreadyExists, channel.Name)
		}
	}

	s.mu.Lock()
	s.channels = append(s.channels, channel)
	s.mu.Unlock()

	s.log.Debug(fmt.Sprintf("Chat.CreateChannel: Created channel '%s'", channel.Name))

	return channel, nil
}

func (s *ChatServer) ListChannels(ctx context.Context, req *pb.ListChannelsRequest) (*pb.Channels, error) {
	s.log.Info("Chat.ListChannels: Listing channels")
	return &pb.Channels{
		Channels: s.channels,
	}, nil
}

func (s *ChatServer) SendMessage(ctx context.Context, message *pb.Message) (*pb.Message, error) {
	s.log.Info(
		fmt.Sprintf(
			"Chat.SendMessage: Recieved message of length %db from '%s' for channel '%s'",
			len(message.Content),
			message.Author,
			message.Channel.Name,
		),
	)

	s.mu.Lock()
	s.messages = append(s.messages, message)
	s.mu.Unlock()

	s.log.Debug("Chat.SendMessage: Accepted message")
	return message, nil
}

func (s *ChatServer) ListMessages(ctx context.Context, req *pb.ListMessagesRequest) (*pb.Messages, error) {
	s.log.Info(fmt.Sprintf("Chat.ListMessages: Listing messages for channel '%s'", req.Channel.Name))
	var messages []*pb.Message
	ch := req.GetChannel()

	for _, m := range s.messages {
		if m.Channel.Name == ch.Name {
			messages = append(messages, m)
		}
	}

	return &pb.Messages{
		Messages: messages,
	}, nil
}
