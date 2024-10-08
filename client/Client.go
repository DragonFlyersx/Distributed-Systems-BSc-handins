package main

import (
	"context"
	"fmt"
	"go/printer"
	"log"
	"os"

	"main/Handin3"

	"google.golang.org/grpc"
)

func PublishMessage(address string, message string) {
	// Connect to the gRPC server
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		print(err)
	}
	defer conn.Close()

	// Create a client
	client := Handin3.NewChittyChatClient(conn)

	// Create a ChatMessage object
	chatMessage := &Handin3.ChatMessage{
		Message: message,
	}

	// Call the publish message method
	client.PublishMessage(context.Background(), chatMessage)
	
}

func ReceiveMessage(address string) {
	// Connect to the gRPC server
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Create a client
	client := Handin3.NewChittyChatClient(conn)

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
	address := "localhost:50051" // address to server

	ReceiveMessage(address) // Maybe start go routine for receiving messages else we just call the method here

	// some sort of loop or condition to keep the method alive indefinitely
}
