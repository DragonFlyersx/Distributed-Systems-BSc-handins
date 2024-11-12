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
var AuctionServer server // Server instance

type server struct {
	AuctionHouse.AuctionServiceServer
	clients map[AuctionHouse.AuctionService_SendResultServer]bool
}

// Constructor for the server
func newServer() *server {
	return &server{
		clients: make(map[AuctionHouse.AuctionService_SendResultServer]bool), // Instantiate the map
	}
}

func (s *server) SendBid(stream AuctionHouse.AuctionService_SendBidServer) error { // receive bid from client
	userBidRequest, err := stream.Recv()
	if err != nil {
		log.Fatalf("Failed to receive a bid: %v", err)
	}
	log.Printf("Bid received from %s : %v", userBidRequest.BidAmount, userBidRequest.BidderName)
	CurrentHighestBid = userBidRequest.BidAmount // Update the current highest bid
	log.Printf("Current highest bid: %v by %s", CurrentHighestBid, userBidRequest.BidderName)

	return stream.Send(&AuctionHouse.BidResponse{Ack: true})
}

func SendResult(stream AuctionHouse.AuctionService_SendResultServer) error { // sends result to client
	result := CurrentHighestBid
	if AuctionStatus == "Open" { // If the auction is still open
		log.Printf("Result sent to clients: %v", result)
		for client := range AuctionServer.clients {
			client.Send(&AuctionHouse.ResultResponse{Status: "Open", Result: result})
		}
	} else { // If the auction is closed
		log.Printf("Result sent to clients: %v", result)
		for client := range AuctionServer.clients {
			client.Send(&AuctionHouse.ResultResponse{Status: "Closed", Result: result})
		}
	}
	return nil
}

func startServer(port string, ip string) {
	//init listener
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	//gRPC server instance
	grpcServer := grpc.NewServer()
	AuctionHouse.RegisterAuctionServiceServer(grpcServer)

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
