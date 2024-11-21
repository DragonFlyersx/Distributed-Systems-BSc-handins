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

	userBidRequest, err := stream.Recv()
	if err != nil {
		log.Fatalf("Failed to receive a bid: %v", err)
	}

	log.Printf("Bid received from %s : %v", userBidRequest.BidAmount, userBidRequest.BidderName)
	var response AuctionHouse.GeneralResponse

	// Check if bid is higher than current highest bid
	if userBidRequest.BidAmount > CurrentHighestBid {
		CurrentHighestBid = userBidRequest.BidAmount
		CurrentHighestBidder = userBidRequest.BidderName
		log.Printf("Current highest bid updated: %v by %s", CurrentHighestBid, CurrentHighestBidder)

		response.Response = &AuctionHouse.GeneralResponse_BidResponse{
			BidResponse: &AuctionHouse.BidResponse{Ack: true},
		}
	} else {
		response.Response = &AuctionHouse.GeneralResponse_BidResponse{
			BidResponse: &AuctionHouse.BidResponse{Ack: false},
		}
	}

	// Register the client
	s.clients[stream] = true

	if response.Response.(*AuctionHouse.GeneralResponse_BidResponse).BidResponse.Ack == true {
		// If the bid is lower than the current highest bid then send a response to the client along with new result
		stream.Send(&response)
		receivedBid = true
		return s.SendResult(&AuctionHouse.Empty{}, stream)
	}

	// If ack is false then send a response to the client with no new result
	return stream.Send(&response)
}

func (s *server) SendResult(in *AuctionHouse.Empty, stream AuctionHouse.AuctionService_SendResultServer) error { // sends result to client

	for {
		if receivedBid == true {
			result := CurrentHighestBid
			Winning_Person := CurrentHighestBidder
			var response AuctionHouse.GeneralResponse

			if AuctionStatus == "Open" { // If the auction is still open
				log.Printf("Result sent to clients: %v", result)
				response.Response = &AuctionHouse.GeneralResponse_ResultResponse{
					ResultResponse: &AuctionHouse.ResultResponse{
						Result:     result,
						Status:     "Open",
						WinnerName: Winning_Person,
					},
				}
			} else { // If the auction is closed
				log.Printf("Result sent to clients: %v", result)
				response.Response = &AuctionHouse.GeneralResponse_ResultResponse{
					ResultResponse: &AuctionHouse.ResultResponse{
						Result:     result,
						Status:     "Closed",
						WinnerName: Winning_Person,
					},
				}
			}

			if err := stream.Send(&response); err != nil {
				log.Printf("Failed to send result to client: %v", err)
				return err
			}
			receivedBid = false
		}
		time.Sleep(5 * time.Second) // Adjust the sleep duration as needed
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
