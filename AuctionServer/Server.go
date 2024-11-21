package main

import (
	"context"
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
	AuctionHouse.UnimplementedAuctionServiceServer
	clients map[AuctionHouse.AuctionService_SendBidServer]bool
}

// Constructor for the server
func newServer() *server {
	return &server{
		clients: make(map[AuctionHouse.AuctionService_SendBidServer]bool), // Instantiate the map
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

		} else {
			log.Printf("Bid too low")
			// Send NACK to the client (NACK means Not Acknowledged)
			response := &AuctionHouse.GeneralResponse{
				Response: &AuctionHouse.GeneralResponse_BidResponse{
					BidResponse: &AuctionHouse.BidResponse{Ack: false},
				},
			}
			if err := stream.Send(response); err != nil {
				log.Printf("Error sending NACK: %v", err)
				return err
			}
		}
	}
}

func (s *server) SendResult(ctx context.Context, in *AuctionHouse.Empty) (*AuctionHouse.GeneralResponse, error) {
	log.Printf("SendResult function called")

	// Send immediate result when queried
	result := &AuctionHouse.GeneralResponse{
		Response: &AuctionHouse.GeneralResponse_ResultResponse{
			ResultResponse: &AuctionHouse.ResultResponse{
				Result:     CurrentHighestBid,
				Status:     AuctionStatus,
				WinnerName: CurrentHighestBidder,
			},
		},
	}

	return result, nil
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
	s.clients = make(map[AuctionHouse.AuctionService_SendBidServer]bool)
}

func (s *server) BroadcastWinner() {

	log.Printf("Auction winner announced")
	for client := range s.clients {
		log.Printf(client.Context().Value("BidderName").(string))
		response := &AuctionHouse.GeneralResponse{
			Response: &AuctionHouse.GeneralResponse_ResultResponse{
				ResultResponse: &AuctionHouse.ResultResponse{
					Result:     CurrentHighestBid,
					Status:     AuctionStatus,
					WinnerName: CurrentHighestBidder,
				},
			},
		}
		if err := client.Send(response); err != nil {
			log.Printf("Error sending result: %v", err)
		}
	}
	s.ClearRegisteredUsers()
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

			timeleft := 5
			for timeleft > 0 {
				timeleft--
				time.Sleep(1 * time.Second)
				log.Printf("Time left: %d", timeleft)
			}

			AuctionStatus = "Closed"
			log.Printf("Auction closed")
			AuctionServer.BroadcastWinner()
		}
	}
}
