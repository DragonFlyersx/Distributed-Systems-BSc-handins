package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"main/Handin5"
	"os"
	"strings"

	"google.golang.org/grpc"
)

// call for the client to bid on the auction

// Listen for messages from the server

// Request the current highest bid from the server








func main() {
	// Set up a connection to the server
	// needs to connect to 3 server Nodes
	NodeOneAddress := "localhost:50051" // Address to the server
	NodeTwoAddress := "localhost:50052" // Address to the server
	NodeThreeAddress := "localhost:50053" // Address to the server


	var clientID string
	clientID, err := os.Hostname() // 
	if err != nil {
		//log.Fatalf("Error getting hostname: %v", err)
		log.Printf("Error getting hostname: %v", err)
		clientID = "Unknown"
	}

	// Connect to the gRPC server
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect SERVER: %v", err)
	}
	defer conn.Close()

	// Create a client
	client := new AuctionClient)
	var userCommand string

	for {
		fmt.Scan(&userCommand)

		if userCommand == "Result" { // Request the current highest bid from the server

		} else if userCommand == "Bid" { // send Bid to the server
			 
		} 
	}
}