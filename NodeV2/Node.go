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
var nextNodeAddress string
var port string
var WantsToBeLeader bool

type Token struct {
	TokenRing.UnimplementedNodeServer
	tokenID int32
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
		fmt.Printf("Received token from node in SendToken() With the ID:  %d\n", token.TokenID)
		if !ReceivingMode {
			ReceivingMode = true // able to only recieve and not generate new Tokens need to change the values of the token instead
			log.Printf("Recieving mode is on")
		}

		if WantsToBeLeader {
			// Handle election message
			if token.TokenID == nodeID {
				log.Printf("Node %d has been elected as leader", nodeID)
				log.Printf("Node %d Is accesing the Critical Section", nodeID)
				time.Sleep(1 * time.Second) // Simulate processing
				log.Printf("Node %d Is Done with the Critical Section", nodeID)
				WantsToBeLeader = false
				// is elected as leader needs to Time for critical section then pass the node along with zero as cliend Id
				//calls sendTokenNExtNode updated ID to 0
				token.TokenID = 0
			} else if token.TokenID < nodeID {
				// Forward the Token since we cant become leader, needs to be only if we want to enter the criticak
				// update The internal value of the Token
				token.TokenID = nodeID
			}
		}

		sendTokenToNextNode(token)
	}
	return nil
}

// sendTokenToNextNode handles sending a token to the next node in the ring.
func sendTokenToNextNode(receivedToken *TokenRing.Token) {
	conn, err := grpc.Dial(nextNodeAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to next node: %v", err)
	}
	defer conn.Close()

	client := TokenRing.NewNodeClient(conn)
	stream, err := client.SendToken(context.Background())
	if err != nil {
		log.Fatalf("Failed to open stream: %v", err)
	}

	// needs to update token with values
	if err := stream.Send(receivedToken); err != nil {
		log.Fatalf("Failed to send token: %v", err)
	}

	// Sending a token to the next node
	fmt.Printf("Sent received token from node %d to %s\n", nodeID, nextNodeAddress)

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
	port = "50051"
	nodeID = rand.Int31() + 1
	nextNodeAddress = "25.11.126.45:50051"

	log.Printf("Node was assigned id %v", nodeID)

	go startServer(port, ip)

	var userCommand string
	for {
		fmt.Scan(&userCommand)

		if userCommand == "AskForElection" { // Ie Enter the Critical Section. Needs to update Token and send First one
			WantsToBeLeader = true
			if !ReceivingMode {

				// Genereate first Token
				var token *TokenRing.Token

				token = &TokenRing.Token{
					TokenID: nodeID,
				}

				sendTokenToNextNode(token)
			}

		}
	}

}
