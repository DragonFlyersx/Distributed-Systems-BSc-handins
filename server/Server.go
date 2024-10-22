package main

import (
	"context"
	"log"
	"net"
	"strconv"

	"main/Handin3" // Import generated from the .proto file

	"google.golang.org/grpc"
)

// Define the server struct Keep the server struct in the server package
type server struct {
	Handin3.ChittyChatServer
	// Handin3.ChatMessage latestMessage := nil
	clients     map[Handin3.ChittyChat_BroadcastMessageServer]bool
	lamportTime int64
}

// Constructor for the server
func newServer() *server {
	return &server{
		clients: make(map[Handin3.ChittyChat_BroadcastMessageServer]bool), // Instantiate the map
	}
}

// Implement the BroadcastMessage method of the ChittyChatServer interface
// Implement the BroadcastMessage method of the ChittyChatServer interface
func (s *server) BroadcastMessage(empty *Handin3.Empty, stream Handin3.ChittyChat_BroadcastMessageServer) error {
	// Register the client
	s.clients[stream] = true
	s.lamportTime += 1
	JoinMessage := &Handin3.ChatMessage{
		Message:   "Participant Joined Chitty-Chat at Lamport time " + strconv.FormatInt(int64(s.lamportTime), 10),
		Timestamp: s.lamportTime,
	}
	s.PublishMessage(stream.Context(), JoinMessage)

	defer func() {
		delete(s.clients, stream)
		s.lamportTime += 1
		DisconnectMessage := &Handin3.ChatMessage{
			Message:   "Participant Left Chitty-Chat at Lamport time " + strconv.FormatInt(int64(s.lamportTime), 10),
			Timestamp: s.lamportTime,
		}
		s.PublishMessage(stream.Context(), DisconnectMessage)
	}() // Clean up on disconnect

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
	log.Printf("Server Received message: \"%s\" - %d \n", msg.Message, s.lamportTime)

	// Find maximum lamport time and increment
	if s.lamportTime < msg.Timestamp {
		s.lamportTime = msg.Timestamp + 1
		log.Print("incrementing lamport time")

	} else { 
		s.lamportTime += 1
		log.Print("incrementing lamport time")
	}
	
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
