package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"main/TokenRing"
	"math/rand"

	"google.golang.org/grpc"
)

var ReceivingMode bool
var nodeID int32

type Token struct {
	TokenRing.UnimplementedNodeServer
	tokenID       int32
	electionState bool
	isLeader      bool
}

func (token *Token) SendToken(stream TokenRing.Node_SendTokenServer) error {
	for {
		token, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// Process the received token
		fmt.Printf("Received token from node in SendToken() %d\n", token.TokenID)
		electionState := token.ElectionState

		if electionState {
			// Handle election message
			if token.TokenID == nodeID {
				log.Printf("Node %d has been elected as leader", nodeID)

				// Send elected notification
				response := &TokenRing.Token{
					TokenID:       nodeID,
					ElectionState: false,
					IsLeader:      true,
				}
				if err := stream.Send(response); err != nil {
					return err
				}
			} else if token.TokenID > nodeID {
				// Forward the election message
				if err := stream.Send(token); err != nil {
					return err
				}
			} else {
				// Regular token passing
				if token.IsLeader {
					log.Printf("Leader (Node %d) processing token", nodeID)
					time.Sleep(1 * time.Second) // Simulate processing
				}

				// Forward the token
				if err := stream.Send(token); err != nil {
					return err
				}
				log.Printf("Node %d forwarded token", token.TokenID)
			}
		}
	}

	// Forward token to the next node (handled in client logic below)

	return nil
}

// sendTokenToNextNode handles sending a token to the next node in the ring.
func sendTokenToNextNode(nextNodeAddress string, nodeID int32, election bool) {
	conn, err := grpc.Dial(nextNodeAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to next node: %v", err)
	}
	defer conn.Close()

	client := TokenRing.NewNodeClient(conn)
	stream, err := client.SendToken(context.Background())
	var token *TokenRing.Token
	if err != nil {
		log.Fatalf("Failed to open stream: %v", err)
	}
	if !ReceivingMode {
		token = &TokenRing.Token{
			TokenID:       nodeID,
			ElectionState: election,
		}
	} else {
		// needs to update token with values

	}
	// Sending a token to the next node

	if err := stream.Send(token); err != nil {
		log.Fatalf("Failed to send token: %v", err)
	}
	fmt.Printf("Sent token from node %d to %s\n", nodeID, nextNodeAddress)

	// Close the stream after sending the token
	if err := stream.CloseSend(); err != nil {
		log.Fatalf("Failed to close stream: %v", err)
	}
}

func startServer(port string, ip string) {
	//init listener
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	//gRPC server instance
	grpcServer := grpc.NewServer()
	TokenRing.RegisterNodeServer(grpcServer, &Token{})

	//listen and server
	err = grpcServer.Serve(listen)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	log.Printf("Listening on port %s", port)
}

func main() {
	ip := "25.8.121.243:50051"
	port := "50051"
	nodeID = rand.Int31() + 1
	nextNodeAddress := "25.11.126.45:50051"
	//isLeader := true // Change this to false for non-leader nodes as there can only be on to start the passing of token

	go startServer(port, ip)

	var userCommand string
	for {
		fmt.Scan(&userCommand)

		if userCommand == "AskForElection" { // Ie Enter the Crititcal Section. Needs to update Token and send First one
			sendTokenToNextNode(nextNodeAddress, nodeID, true)
		}
	}

}
