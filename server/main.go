package main

import (
	"context"
	"fmt"
	"net"

	pb "github.com/sam-kenney/grpc_chat/chat"
	"github.com/sam-kenney/logging"
	"github.com/sam-kenney/server/services"
	"google.golang.org/grpc"
)

var log logging.Logger

func main() {
	var err error
	log, err = logging.NewLoggerFromEnv(context.Background(), logging.Default, "chat")
	if err != nil {
		panic(fmt.Sprintf("Failed to initialise logger: %s", err.Error()))
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 8080))
	if err != nil {
		log.Err(err)
	}
	defer listener.Close()

	server := grpc.NewServer()
	defer server.Stop()
	pb.RegisterChatServer(server, services.NewChatServer(log))

	log.Info(fmt.Sprintf("Listening at http://%s", listener.Addr()))
	server.Serve(listener)
}
