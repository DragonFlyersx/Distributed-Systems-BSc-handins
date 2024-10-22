package main

import (
	"context"
	"fmt"
	"log"

	"main/Handin3"

	"google.golang.org/grpc"
)

func ListenForMessage(client Handin3.ChittyChatClient) {

	for {
		var message string
		fmt.Scanln(&message)

		if message == "server.exit" {
			// Client should leave from server
			break
		} else if message != "" {
			PublishMessage(client, message)
		}

	}
}

func PublishMessage(client Handin3.ChittyChatClient, message string) {
	// set local lamport from message
	// Create a ChatMessage object
	chatMessage := &Handin3.ChatMessage{
		Message: message,
	}

	// Call the publish message method
	client.PublishMessage(context.Background(), chatMessage)

}

func ReceiveMessage(client Handin3.ChittyChatClient) {
	// set local lamport from message

	// Call the server
	stream, err := client.BroadcastMessage(context.Background(), &Handin3.Empty{})
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	// Read message from stream
	for {
		chatMessage, err := stream.Recv()
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		fmt.Printf("Received message: %s\n", chatMessage.Message)
	}
}

func main() {
	address := "192.168.0.72:50051" // address to server
	//var lamportTime *int

	// Connect to the gRPC server
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Create a client
	client := Handin3.NewChittyChatClient(conn)

	go ReceiveMessage(client) // Listen for any messages from server

	go ListenForMessage(client) // Listen for any messages sent by client

	select {} // This will block the main goroutine indefinitely
}
