package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"main/Handin3" // Import generated from the .proto file

	"google.golang.org/grpc"
)

// Define the server struct Keep the server struct in the server package
type server struct {
	Handin3.ChittyChatServer
	// Handin3.ChatMessage latestMessage := nil
	latestMessage *Handin3.ChatMessage
}

// Implement the BroadcastMessage method of the ChittyChatServer interface
func (s *server) BroadcastMessage(empty *Handin3.Empty, stream Handin3.ChittyChat_BroadcastMessageServer) error {
	//var latestMessage *Handin3.ChatMessage

	/*for {
		// Receive message from the stream
		message, err := stream.Recv()
		if err == io.EOF {
			// If the stream is finished, break the loop
			break
		}
		if err != nil {
			// Handle error if needed
			return err
		}

		// Process the received message (save it as the latest message)
		latestMessage = message
	}

	// After the stream ends, latestMessage will hold the last message received
	if latestMessage != nil {
		log.Printf("Latest message: %v", latestMessage.GetMessage())
	} else {
		log.Println("No messages received.")
	}
	*/
	return nil
}

func (s *server) PublishMessage(ctx context.Context, msg *Handin3.ChatMessage) (*Handin3.Empty, error) {
	// increment The logical timestamp
	s.latestMessage = msg
	log.Println("Message received:", msg.GetMessage())

	return &Handin3.Empty{}, nil
}

// This function implements the server-side logic for streaming messages to the client
func (s *server) BroadcastMessages(empty *Handin3.Empty, stream Handin3.ChittyChat_BroadcastMessageServer) error {
	messages := []string{
		"Hello from server",
		"This is the second message",
		"And here comes the third message",
	}

	for _, msg := range messages {
		// Send each message to the client over the stream
		chatMessage := &Handin3.ChatMessage{
			Message: msg,
		}

		log.Printf("Sending message: %v", chatMessage.Message)

		// Send the message
		if err := stream.Send(chatMessage); err != nil {
			log.Printf("Failed to send message: %v", err)
			return fmt.Errorf("failed to send message: %v", err)
		}
	}

	log.Println("BroadcastMessages completed successfully")

	return nil // Close the stream when done
}

// Keep the main function in the server package
func main() {
	// Define the port for the server to listen on
	port := ":50051"
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	// Create a new gRPC server
	grpcServer := grpc.NewServer()

	// Register the TimeService with the gRPC server  // changed
	Handin3.RegisterChittyChatServer(grpcServer, &server{})

	log.Printf("ChittyChat server is running on port %s", port)

	// Start serving requests
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
