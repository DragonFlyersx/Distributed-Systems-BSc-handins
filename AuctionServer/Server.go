package main

import (
	"fmt"
	"log"
	"net"
	"time"

	AuctionHouse "main/Handin5"

	"google.golang.org/grpc"
)

var CurrentHighestBidder string = ""
var CurrentHighestBid int32 = 0
var AuctionStatus string = "Closed"
var AuctionServer server // Server instance

type server struct {
	AuctionHouse.AuctionServiceServer
	clients map[AuctionHouse.AuctionService_SendResultServer]bool
}

type auctionServiceServer struct {
	AuctionHouse.UnimplementedAuctionServiceServer
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

	// Register the client if not already registered
	if _, exists := s.clients[stream]; !exists {
		s.clients[stream] = true
	}

	log.Printf("Bid received from %s : %v", userBidRequest.BidAmount, userBidRequest.BidderName)
	if userBidRequest.BidAmount > CurrentHighestBid { // If the bid is higher than the current highest bid
		CurrentHighestBid = userBidRequest.BidAmount // Update the current highest bid
		CurrentHighestBidder = userBidRequest.BidderName
		log.Printf("Current highest bid: %v by %s", CurrentHighestBid, userBidRequest.BidderName)

		return stream.Send(&AuctionHouse.BidResponse{Ack: true})
	}

	// If the bid is lower than the current highest bid then send a response to the client
	return stream.Send(&AuctionHouse.BidResponse{Ack: false})
}

func (s *server) SendResult(stream AuctionHouse.AuctionService_SendResultServer) error { // sends result to client
	result := CurrentHighestBid
	Winning_Person := CurrentHighestBidder
	if AuctionStatus == "Open" { // If the auction is still open
		log.Printf("Result sent to clients: %v", result)
		for client := range s.clients {
			client.Send(&AuctionHouse.ResultResponse{Status: "Open", Result: result, WinnerName: Winning_Person})
		}
	} else { // If the auction is closed
		log.Printf("Result sent to clients: %v", result)
		for client := range s.clients {
			client.Send(&AuctionHouse.ResultResponse{Status: "Closed", Result: result, WinnerName: Winning_Person})
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
	service := &auctionServiceServer{}
	AuctionHouse.RegisterAuctionServiceServer(grpcServer, service)

	//listen and server
	log.Printf("Ready to receive and listening on port %s", port)

	err = grpcServer.Serve(listen)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func (s *server) ClearRegisteredUsers() {
	s.clients = make(map[AuctionHouse.AuctionService_SendResultServer]bool)
}

func main() {
	ip := "Local:50051" // Ip of the server
	port := "50051"
	go startServer(port, ip)

	var userCommand string
	for {
		fmt.Scan(&userCommand)

		if userCommand == "Start" && AuctionStatus == "Closed" { // Request the current highest bid from the server
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
