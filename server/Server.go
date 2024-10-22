package main

import (
	"context"
	"log"
	"net"

	"main/Handin3" // Import generated from the .proto file

	"google.golang.org/grpc"
)

// Define the server struct Keep the server struct in the server package
type server struct {
	Handin3.ChittyChatServer
	// Handin3.ChatMessage latestMessage := nil
	clients map[Handin3.ChittyChat_BroadcastMessageServer]bool
}

// Constructor for the server
func newServer() *server {
	return &server{
		clients: make(map[Handin3.ChittyChat_BroadcastMessageServer]bool), // Initialize the map
	}
}

// Implement the BroadcastMessage method of the ChittyChatServer interface
// Implement the BroadcastMessage method of the ChittyChatServer interface
func (s *server) BroadcastMessage(empty *Handin3.Empty, stream Handin3.ChittyChat_BroadcastMessageServer) error {
	// // After the stream ends, latestMessage will hold the last message received
	// if s.latestMessage != nil {
	// 	log.Printf("Latest message: %v", s.latestMessage.GetMessage())
	// } else {
	// 	log.Println("No messages received.")
	// }
	// // put in stream

	// return nil
	// Register the client
	s.clients[stream] = true
	defer func() { delete(s.clients, stream) }() // Clean up on disconnect

	// Keep the stream open to send messages to this client
	for {
		// Optionally wait for a disconnection
		if err := stream.Context().Err(); err != nil {
			break // Exit the loop if the context is done
		}
	}
	return nil
}

func (s *server) PublishMessage(ctx context.Context, msg *Handin3.ChatMessage) (*Handin3.Empty, error) {
	// increment The logical timestamp

	// Log the message received
	log.Println("Message received:", msg.GetMessage())

	// Broadcast the message to all connected clients
	for client := range s.clients {
		if err := client.Send(msg); err != nil {
			log.Printf("Failed to send message to a client: %v", err)
			// Optionally handle client disconnection
		}
	}

	return &Handin3.Empty{}, nil
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
	Handin3.RegisterChittyChatServer(grpcServer, newServer())

	log.Printf("ChittyChat server is running on port %s", port)
	log.Println("Server is ready to accept connections on port", port)
	// Start serving requests
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
