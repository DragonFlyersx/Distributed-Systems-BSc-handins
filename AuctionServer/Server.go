package main

import (
	"context"
	"log"
	"net"
	"time"

	AuctionHouse "main/Handin5"

	"google.golang.org/grpc"
)

var CurrentHighestBidder string = ""
var CurrentHighestBid int32 = 0
var AuctionStatus string = "Closed"
var restartAvailable bool = true

// var AuctionServer server // Server instance

type server struct {
	AuctionHouse.UnimplementedAuctionServiceServer
	clients map[AuctionHouse.AuctionService_SendBidServer]bool // chose map instead of slice as it is easier to check if a client is already in the map
	// clients []AuctionHouse.AuctionService_SendBidServer
}

// Constructor for the server
func newServer() *server {
	return &server{
		// clients: make([]AuctionHouse.AuctionService_SendBidServer, 0),
		clients: make(map[AuctionHouse.AuctionService_SendBidServer]bool), // Instantiate the map
	}
}

func (s *server) SendBid(stream AuctionHouse.AuctionService_SendBidServer) error { // receive bid from client
	log.Printf("Server SendBid function called")

	// s.clients = append(s.clients, stream)
	// Add the client stream to the clients map
	// if statement to check if the client is already in the map
	if _, exists := s.clients[stream]; !exists {
		s.clients[stream] = true
		log.Printf("Client added. Total clients: %d", len(s.clients))
	}

	for {
		userBidRequest, err := stream.Recv()
		if err != nil {
			log.Printf("Failed to receive a bid: %v", err)
			return err
		}

		log.Printf("Bid received from %s : %v", userBidRequest.BidderName, userBidRequest.BidAmount)

		if AuctionStatus == "Closed" && restartAvailable == true {
			restartAvailable = false
			CurrentHighestBid = 0
			CurrentHighestBidder = ""
			go s.openAuction() // Start the auction
		}

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

func startServer(port string, ip string, s *server) {
	//init listener
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	//gRPC server instance
	grpcServer := grpc.NewServer()
	AuctionHouse.RegisterAuctionServiceServer(grpcServer, s)

	//listen and server
	log.Printf("Ready to receive and listening on port %s", port)

	err = grpcServer.Serve(listen)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func (s *server) ClearRegisteredUsers() {
	s.clients = make(map[AuctionHouse.AuctionService_SendBidServer]bool)
	log.Printf("Total clients: %d", len(s.clients))
	// s.clients = make([]AuctionHouse.AuctionService_SendBidServer, 0)
}

func (s *server) BroadcastWinner() {
	log.Printf("BroadcastWinner function called")
	log.Printf("Number of clients: %d", len(s.clients))
	for client := range s.clients {
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
			log.Printf("Error sending result to client: %v", err)
		} else {
			log.Printf("Result sent to client")
		}
	}
	s.ClearRegisteredUsers()
	log.Printf("Clearing registered users")
}

// function for opening of the auction house
func (s *server) openAuction() {
	AuctionStatus = "Open"
	log.Printf("Auction started")
	timeleft := 15
	for timeleft > 0 {
		timeleft--
		time.Sleep(1 * time.Second)
		log.Printf("Time left: %d", timeleft)
	}

	AuctionStatus = "Closed"
	log.Printf("Auction closed")
	s.BroadcastWinner()

	// Wait for 5 seconds before allowing new bids to start a new auction
	time.Sleep(5 * time.Second)
	restartAvailable = true
}

func main() {
	ip := "Local:50053" // Ip of the server
	port := "50053"
	s := newServer()
	go startServer(port, ip, s)

	// Ensure the server runs indefinitely
	select {}
}
