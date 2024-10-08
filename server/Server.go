package main

import (
	"io"
	"log"
	"net"

	"main/Handin3" // Import generated from the .proto file

	"google.golang.org/grpc"
)

// Define the server struct Keep the server struct in the server package
type server struct {
	Handin3.ChittyChatServer
}

// Implement the BroadcastMessage method of the ChittyChatServer interface
func (s *server) BroadcastMessage(*Handin3.Empty, grpc.ServerStreamingServer[Handin3.ChatMessage]) error {
	var latestMessage *Handin3.ChatMessage

	for {
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

	return nil
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
