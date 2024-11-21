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
var receivedBid = false

type server struct {
	AuctionHouse.UnimplementedAuctionServiceServer
	clients map[AuctionHouse.AuctionService_SendResultServer]bool
}

// Constructor for the server
func newServer() *server {
	return &server{
		clients: make(map[AuctionHouse.AuctionService_SendResultServer]bool), // Instantiate the map
	}
}

func (s *server) SendBid(stream AuctionHouse.AuctionService_SendBidServer) error { // receive bid from client
	log.Printf("Server SendBid function called")

	// Add the client stream to the clients map
	s.clients[stream] = true

	for {
		userBidRequest, err := stream.Recv()
		if err != nil {
			log.Printf("Failed to receive a bid: %v", err)
			return err
		}

		log.Printf("Bid received from %s : %v", userBidRequest.BidderName, userBidRequest.BidAmount)

		// Check if bid is higher than current highest bid
		if userBidRequest.BidAmount > CurrentHighestBid {
			CurrentHighestBid = userBidRequest.BidAmount
			CurrentHighestBidder = userBidRequest.BidderName
			log.Printf("New highest bid: %v by %s", CurrentHighestBid, CurrentHighestBidder)

			// Trigger sending results to all clients
			receivedBid = true
		}

		// Send ACK to the client
		response := &AuctionHouse.GeneralResponse{
			Response: &AuctionHouse.GeneralResponse_BidResponse{
				BidResponse: &AuctionHouse.BidResponse{Ack: true},
			},
		}
		if err := stream.Send(response); err != nil {
			log.Printf("Error sending ACK: %v", err)
			return err
		}
	}
}

func (s *server) SendResult(in *AuctionHouse.Empty, stream AuctionHouse.AuctionService_SendResultServer) error { // send results to the client
	log.Printf("SendResult function called")

	for {
		if receivedBid {
			// Send the latest auction result to the client
			result := &AuctionHouse.GeneralResponse{
				Response: &AuctionHouse.GeneralResponse_ResultResponse{
					ResultResponse: &AuctionHouse.ResultResponse{
						Result:     CurrentHighestBid,
						Status:     AuctionStatus,
						WinnerName: CurrentHighestBidder,
					},
				},
			}

			log.Printf("Sending updated auction result to client: %v", result)
			if err := stream.Send(result); err != nil {
				log.Printf("Error sending result: %v", err)
				return err
			}

			// Reset the flag
			receivedBid = false
		}

		time.Sleep(1 * time.Second) // Wait before checking for the next update
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
	service := newServer()
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
			receivedBid = true
			time.Sleep(100 * time.Second) // The auction runs for 100 seconds
			// send result of auction to all clients

			AuctionStatus = "Closed"
			log.Printf("Auction closed")
		}
	}
}
