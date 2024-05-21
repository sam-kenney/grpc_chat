// Client to test gRPC server. Will update with something proper at some point...
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/sam-kenney/grpc_chat/chat"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func printMessage(msg *pb.Message) {
	log.Printf("%s: '%s' in %s at %d", msg.Author, msg.Content, msg.Channel.Name, *msg.Timestamp)
}

func main() {
	ctx := context.Background()

	// TODO: Real credentials
	conn, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewChatClient(conn)

	channel := &pb.Channel{Name: "My Channel"}
	ch, err := client.CreateChannel(ctx, channel)
	if err != nil {
		log.Printf("Channel already exists: %s", channel.Name)
	} else {
		log.Println(fmt.Sprintf("Created channel %s", ch.Name))
	}

	channels, err := client.ListChannels(ctx, &pb.ListChannelsRequest{})
	if err != nil {
		log.Fatalf("Error: %s", err.Error())
	}
	log.Println("Channels")
	for _, c := range channels.Channels {
		log.Printf("%s", c)
	}

	ts := time.Now().Unix()

	_, err = client.SendMessage(ctx, &pb.Message{
		Content:   "Hello, world!",
		Author:    "Sam",
		Channel:   channel,
		Timestamp: &ts,
	})
	if err != nil {
		log.Fatalf("Error: %s", err.Error())
	}

	messages, err := client.ListMessages(ctx, &pb.ListMessagesRequest{Channel: channel})
	if err != nil {
		log.Fatalf("Error: %s", err.Error())
	}
	for _, m := range messages.Messages {
		printMessage(m)
	}
}
