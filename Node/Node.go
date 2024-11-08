package main

import (
	"context"
	"fmt"
	"log"
	"main/Handin4"
	"math/rand"
	"net"
	"time"

	"google.golang.org/grpc"
)

var NodeNumber int32
var ReceivingMode bool
var RequestedAccess bool
var isLeader bool

// SendToken sends a token to the next node in the chain
func SendToken(node Handin4.NodeClient, message Handin4.NodeMessage) {
	_, err := node.SendToken(context.Background(), &message)
	if err != nil {
		log.Fatalf("Error sending token: %v", err)
	}
}

// ReceiveToken listens for messages from other nodes
/*func ReceiveToken(node Handin4.NodeClient) {
	stream, err := node.ReceiveToken(context.Background(), &Handin4.Empty{})
	if err != nil {
		log.Fatalf("Error receiving token: %v", err)
	}

	for {
		nodeMessage, err := stream.Recv()
		if err != nil {
			log.Fatalf("Error receiving token: %v", err)
		}
		if nodeMessage.Value != 0 {
			ReceivingMode = true
		}

		// Implement token passing logic
		if nodeMessage.Value != NodeNumber {
			if !RequestedAccess || (RequestedAccess && NodeNumber < nodeMessage.Value) {
				SendToken(node, *nodeMessage)
			} else if RequestedAccess && NodeNumber > nodeMessage.Value {
				nodeMessage.Value = NodeNumber
				SendToken(node, *nodeMessage)
			}
		} else {
			if !isLeader {
				fmt.Println("You are elected as the leader")
				isLeader = true
			}
		}
	}
}*/

// RegisterServer sets up the gRPC server to listen for incoming tokens
func RegisterServer() {
	port := ":50051" // Change this for each node to run on a different port
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	grpcServer := grpc.NewServer()
	Handin4.RegisterNodeServer(grpcServer, &NodeServerImpl{})

	log.Printf("NodeServer is running on port %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func main() {
	go RegisterServer()

	// Each node's connection to its neighbor (next node)
	nextNodeAddress := "25.8.113.191:50051" // Replace with actual next node's address
	var conn *grpc.ClientConn
	var err error

	// Retry mechanism for connecting to the next node
	for {
		conn, err = grpc.Dial(nextNodeAddress, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(5*time.Second))
		if err == nil {
			log.Printf("Successfully connected to %s", nextNodeAddress)
			break
		}
		log.Printf("Failed to connect to %s: %v. Retrying in 5 seconds...", nextNodeAddress, err)
		time.Sleep(5 * time.Second)
	}
	defer conn.Close()

	node := Handin4.NewNodeClient(conn)

	// Start receiving tokens from other nodes

	for {
		var userCommand string
		fmt.Print("To enter the critical section (Access/Done): ")
		fmt.Scan(&userCommand)

		if userCommand == "Access" {
			queue := rand.Int31() + 1
			RequestedAccess = true
			if RequestedAccess && !ReceivingMode {
				fmt.Println("Starting token passing system with initial token.")
				ReceivingMode = true
				NodeNumber = queue

				// creating the token
				NewMessage := &Handin4.NodeMessage{
					Value: NodeNumber,
				}
				SendToken(node, *NewMessage)
			} else if RequestedAccess && ReceivingMode {
				NodeNumber = queue
			}
		}

		if userCommand == "Done" {
			RequestedAccess = false
			NodeNumber = 0
			fmt.Println("Done with the Access to the critical section.")
			isLeader = false
		}

		select {}
	}
}
