package main

import (
	"bufio"
	"context"
	"log"
	"main/Handin3"
	"os"
	"sync"

	"google.golang.org/grpc"
)

var lamportTime int64 = 0
var lamportMutex sync.Mutex

func incrementLamportTime() {
	lamportMutex.Lock()
	lamportTime++
	lamportMutex.Unlock()
}

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
	if len(message) < 128 {
		// Create a ChatMessage object
		lamportTime += 1
		chatMessage := &Handin3.ChatMessage{
			Message:   message,
			Timestamp: lamportTime,
		}

		// Call the publish message method
		client.PublishMessage(context.Background(), chatMessage)
	} else {
		print("Message is over 128 characters.")
	}
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

		// Find maximum lamport time and increment
		if lamportTime < chatMessage.Timestamp {
			lamportTime = chatMessage.Timestamp + 1
		} else {
			lamportTime += 1
		}

		log.Printf("Client Received message: '%s' - %d \n", chatMessage.Message, lamportTime)
	}
}

func main() {
	address := "192.168.127.41:50051" // Address to the server

	//var lamportTime int

	// Connect to the gRPC server
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect SERVER: %v", err)
	}
	defer conn.Close()

	// Create a client
	client := Handin3.NewChittyChatClient(conn)

	go ReceiveMessage(client) // Listen for any messages from server

	go ListenForMessage(client) // Listen for any messages sent by client

	select {} // This will block the main goroutine indefinitely
}
