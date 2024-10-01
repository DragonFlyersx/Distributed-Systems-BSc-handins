package main

import (
	"context"
	"log"
	"net"
	"time"

	"bob/Exercise05" // Import generated from the .proto file

	"google.golang.org/grpc"
)

// Define the server struct
type server struct {
    Exercise05.UnimplementedTimeServiceServer
	//Timepb.UnimplementedTimeServiceServer
}

// Implement the GetTime method
func (s *server) GetTime(ctx context.Context, req *Exercise05.TimeRequest) (*Exercise05.TimeResponse, error) {
	// Get the current time
	currentTime := time.Now().Format(time.RFC3339)

	// Return the current time in the response
	return &Exercise05.TimeResponse{CurrentTime: currentTime}, nil
}

func main() {
	// Define the port for the server to listen on
	port := ":50051"
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	// Create a new gRPC server
	grpcServer := grpc.NewServer()

	// Register the TimeService with the gRPC server
	Exercise05.RegisterTimeServiceServer(grpcServer, &server{})
    

	log.Printf("TimeService server is running on port %s", port)

	// Start serving requests
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
