package main

import (
	"context"
	"fmt"
	"log"
	"main/Handin4"
)

// NodeServerImpl is the server implementation for handling token passing
type NodeServerImpl struct {
	Handin4.UnimplementedNodeServer
}

// ReceiveToken is the server method to handle incoming tokens
func (s *NodeServerImpl) ReceiveToken(*Handin4.Empty, grpc.ServerStreamingServer[Handin4.NodeMessage]) {
	// Simulate receiving the token (in a real system, you would process it)
	stream, err := node.ReceiveToken(context.Background(), &Handin4.Empty{})
	if err != nil {
		log.Fatalf("Error receiving token: %v", err)
		fmt.Println("Receiving token...")
		// For now, we just return a dummy token message
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
	}
}
