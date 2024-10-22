package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"main/Handin3"

	"google.golang.org/grpc"
)

func ListenForMessage(client Handin3.ChittyChatClient) {

	for {
		reader := bufio.NewReader(os.Stdin)

		message, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Error reading input: %v", err)
		}

		message = message[:len(message)-1] // Remove newline character

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
		log.Fatalf("did not connect SERVER CALL RECIEVE: %v", err)
	}

	// Read message from stream
	for {
		chatMessage, err := stream.Recv()
		if err != nil {
			log.Fatalf("No messages to recieve: %v", err)
		}
		fmt.Printf("Received message: %s\n", chatMessage.Message)
		//lamportTime = chatMessage.Timestamp
	}
}

func main() {
	address := "localhost:50051" // Address to the server

	//var lamportTime int

	// Connect to the gRPC server
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect SERVER: %v", err)
	}
	defer conn.Close()

	// Create a client
	client := Handin3.NewChittyChatClient(conn)

	//PublishMessage(client, "Client Joined!")

	go ReceiveMessage(client) // Listen for any messages from server

	go ListenForMessage(client) // Listen for any messages sent by client

	select {} // This will block the main goroutine indefinitely
}
