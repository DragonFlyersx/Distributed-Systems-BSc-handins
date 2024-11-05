package main

import (
	"context"
	"fmt"
	"log"
	"main/Handin4"
	"os/user"

	"google.golang.org/grpc"
)

// implemt logic for critical server access
// implement logic for token passing

func SendToken(node Handin4.NodeClient, message Handin4.NodeMessage) {
	// Send token to next node

}

// This function listens for any messages from other nodes
func ReceiveToken(node Handin4.NodeMessage) {
	stream, err := node.SendToken(context.Background())
	if err != nil {
		log.Fatalf("Error sending token: %v", err)
	}

	// Read message from stream
	for {
		NodeMessage, err := stream.Recv()
		if err != nil {
			log.Fatalf("Recieved Token %v", err)
		}

		// Implement logic for receiving token
		// Send token to next node
		SendToken(node, NodeMessage)
	}
}

func main() {
	nextNodeAddress := "localhost:50051" // Address of the next node in the chain

	// Connect to the next node in the chain
	conn, err := grpc.Dial(nextNodeAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect to N: %v", err)
	}
	defer conn.Close()

	// Create a Node
	node := Handin4.NodeClient(conn)

	go ReceiveToken(node) // Listen for any messages from Nodes

	var userCommand string
	fmt.Scan(&userCommand)

	if userCommand == "Access" {
		// send Token with value for access global Int
		// Send token to next node
	}

	select {} // This will block the main goroutine indefinitely
}
