package main

import (
	"fmt"
	"log"
	"net"
	"time"

	AuctionHouse "main/Handin5"

	"google.golang.org/grpc"
)

var CurrentHighestBid int32 = 0
var AuctionStatus string = "Closed"

type server struct {
	AuctionHouse.AuctionServiceServer
	// Handin3.ChatMessage latestMessage := nil
	clients     map[AuctionHouse.AuctionService_SendResultServer]bool
	lamportTime int64
}

// Constructor for the server
func newServer() *server {
	return &server{
		clients: make(map[AuctionHouse.AuctionService_SendResultServer]bool), // Instantiate the map
	}
}

func SendBid(stream AuctionHouse.AuctionService_SendBidServer) error { // receive bid from client
	userBidRequest, err := stream.Recv()
	if err != nil {
		log.Fatalf("Failed to receive a bid: %v", err)
	}
	log.Printf("Bid received from client: %v", userBidRequest)
	CurrentHighestBid = userBidRequest.BidAmount // Update the current highest bid
	log.Printf("Current highest bid: %v", CurrentHighestBid)

	return stream.Send(&AuctionHouse.BidResponse{Ack: true})
}

func SendResult(stream AuctionHouse.AuctionService_SendResultServer) error { // sends result to client
	result := CurrentHighestBid
	if AuctionStatus == "Open" { // If the auction is still open
		log.Printf("Result sent to client: %v", result)
		return stream.Send(&AuctionHouse.ResultResponse{Status: "Open", Result: result})
	} else { // If the auction is closed
		log.Printf("Result sent to client: %v", result)
		return stream.Send(&AuctionHouse.ResultResponse{Status: "Closed", Result: result})
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
	AuctionHouse.RegisterAuctionServiceServer(grpcServer, &ResultResponse{})

	//listen and server
	log.Printf("Ready to receive and listening on port %s", port)

	err = grpcServer.Serve(listen)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func main() {
	ip := "Local:50051" // Ip of the server
	port := "50051"
	go startServer(port, ip)

	var userCommand string
	for {
		fmt.Scan(&userCommand)

		if userCommand == "Start Auction" && AuctionStatus == "Closed" { // Request the current highest bid from the server
			AuctionStatus = "Open"
			CurrentHighestBid = 0
			log.Printf("Auction started")
			time.Sleep(100 * time.Second) // The auction runs for 100 seconds
			// send result of auction to all clients

			AuctionStatus = "Closed"
			log.Printf("Auction closed")
		}
	}
}
