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
func (s *server) BroadcastMessage(clientID *Handin3.ClientID, stream Handin3.ChittyChat_BroadcastMessageServer) error {
	// Register the client
	s.clients[stream] = true
	s.lamportTime += 1
	log.Printf("LP: %s: Client: %s connected", strconv.FormatInt(int64(s.lamportTime), 10), clientID)
	// Log when client joins a server
	cMessage := string(clientID.ClientID + " Joined Chitty-Chat")
	JoinMessage := &Handin3.ChatMessage{
		Message:   cMessage,
		Timestamp: s.lamportTime,
		ClientID:  "Server",
	}
	s.PublishMessage(stream.Context(), JoinMessage)

	defer func() {
		delete(s.clients, stream)
		s.lamportTime += 1
		log.Printf("LP: %s: Client: %s Disconnected", strconv.FormatInt(int64(s.lamportTime), 10), clientID)
		dMessage := string(clientID.ClientID + " Left Chitty-Chat")
		DisconnectMessage := &Handin3.ChatMessage{
			Message:   dMessage,
			Timestamp: s.lamportTime,
			ClientID:  "Server",
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
	var receivedMessage = msg
	// Find maximum lamport time and increment
	if s.lamportTime < receivedMessage.Timestamp {
		s.lamportTime = receivedMessage.Timestamp + 1

	} else {
		s.lamportTime += 1
	}

	// increase the lamport time after recieving the message
	log.Printf("LP: %d: Server Received message: \"%s\"\n", s.lamportTime, receivedMessage.Message)

	// Increase the lamport time when resending to the clients
	s.lamportTime += 1
	log.Printf("LP: %d: Server sent message: \"%s\"\n", s.lamportTime, receivedMessage.Message)

	// Update timestamp for received message
	receivedMessage.Timestamp = s.lamportTime

	// Broadcast the message to all connected clients
	for client := range s.clients {
		if err := client.Send(receivedMessage); err != nil {
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
