package main

import (
	"context"
	"fmt"
	"log"
	"main/Handin4"
	"math/rand/v2"

	"google.golang.org/grpc"
)

var NodeNumber int32
var ReceivingMode bool
var RequestedAccess bool
var isLeader bool

// var LeaderID int32

// implemt logic for critical server access
// implement logic for token passing

func SendToken(node Handin4.NodeClient, message Handin4.NodeMessage) {
	// Call the publish message method
	_, err := node.SendToken(context.Background(), &message)
	if err != nil {
		log.Fatalf("Error sending token: %v", err)
	}
}

// This function listens for any messages from other nodes
func ReceiveToken(node Handin4.NodeClient) {
	stream, err := node.ReceiveToken(context.Background(), &Handin4.Empty{})
	if err != nil {
		log.Fatalf("Error recieving token: %v", err)
	}

	// Read message from stream
	for {
		NodeMessage, err := stream.Recv()
		if err != nil {
			log.Fatalf("Recieving Token %v", err)
		}
		if NodeMessage.Value != 0 {
			ReceivingMode = true
		}

		// here we hold the token while in the critical section

		// Implement logic for receiving token
		// Send token to next node
		if NodeMessage.Value != NodeNumber {
			if !RequestedAccess || RequestedAccess && NodeNumber < NodeMessage.Value {
				SendToken(node, *NodeMessage)
			} else if RequestedAccess && NodeNumber > NodeMessage.Value {
				NodeMessage.Value = NodeNumber
				SendToken(node, *NodeMessage)
			}
		} else {
			if !isLeader {
				fmt.Println("You are elected as the leader")
				isLeader = true // Set to true to avoid print multiple times
			}
			// You are elected as the leader

		}
	}
}

func main() {
	nextNodeAddress := "localhost:50051" // Address of the next node in the chain

	// Connect to the next node in the chain
	conn, err := grpc.Dial(nextNodeAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect to N: %v", err)
	}
	defer conn.Close() // Ensure the connection is closed when main exits

	// Create a Node
	node := Handin4.NewNodeClient(conn)

	go ReceiveToken(node) // Listen for any messages from Nodes

	for {
		var userCommand string
		fmt.Print("To enter the critical section (Access/Done)")
		fmt.Scan(&userCommand)

		if userCommand == "Access" {
			// send Token with value for access global Int
			// Send token to next node
			queue := rand.Int32() + 1
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
				// send first Token / call Election

			} else if RequestedAccess && ReceivingMode {
				NodeNumber = queue
				// This wants to acces the critical / call Election
			}

		}
		if userCommand == "Done" {
			RequestedAccess = false
			NodeNumber = 0
			fmt.Println("Done with the Access to the critical section.")
			isLeader = false // set to false since the node is done being the leader
			// send the token along
			// now with the new 0 value in the token / Rekinquish Elected Right

		}

		select {} // This will block the main goroutine indefinitely
	}
}
