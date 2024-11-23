package frontend

import (
	"context"
	"fmt"
	"log"
	AuctionHouse "main/Handin5"

	"google.golang.org/grpc"
)

var clientId string = ""
var bidfromclient int32 = 0

type Frontend struct {
	clients []AuctionHouse.AuctionServiceClient
}

func NewFrontend() *Frontend {
	NodeOneAddress := "localhost:50051" // Address to the server
	// NodeTwoAddress := "localhost:50052"   // Address to the server
	// NodeThreeAddress := "localhost:50053" // Address to the server
	nodeAddresses := []string{NodeOneAddress}

	var clients []AuctionHouse.AuctionServiceClient
	for _, address := range nodeAddresses {
		conn, err := grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Failed to connect to server: %v", err)
		}
		client := AuctionHouse.NewAuctionServiceClient(conn)
		clients = append(clients, client)
	}
	return &Frontend{clients: clients}
}

func (f *Frontend) SendBid(bidRequest *AuctionHouse.UserBidRequest) string {
	log.Printf("SendBid called")
	var response string
	for _, client := range f.clients {
		// Send the bid to the server
		stream, err := client.SendBid(context.Background())
		if err != nil {
			log.Fatalf("Error sending bid: %v", err)
		}

		err = stream.Send(&AuctionHouse.UserBidRequest{
			BidAmount:  bidRequest.BidAmount,  // The amount being bid
			BidderName: bidRequest.BidderName, // The name of the bidder
		})
		if err != nil {
			log.Fatalf("Error sending bid: %v", err)
		}

		// Receive the response from the server
		serverResponse, err := stream.Recv()
		if err != nil {
			log.Fatalf("Error receiving bid: %v", err)
		}

		switch resp := serverResponse.Response.(type) {
		case *AuctionHouse.GeneralResponse_BidResponse:
			if resp.BidResponse.Ack {
				response = "[YOUR BID WAS ACCEPTED]\n"
			} else {
				response = "[YOUR BID WAS TOO LOW]\n"
			}
		case *AuctionHouse.GeneralResponse_ResultResponse:
			response = fmt.Sprintf("The auction has ended! Winning bid is: %d from Bidder: %s\n", resp.ResultResponse.Result, resp.ResultResponse.WinnerName)
		}

		// // Print the response
		// if serverResponse.GetBidResponse().Ack { // If the bid was accepted
		// 	response = "[YOUR BID WAS ACCEPTED]\n"
		// } else { // If the bid was too low
		// 	response = "[YOUR BID WAS TOO LOW]\n"
		// }
	}
	// Return general response here
	return response
}

// Method with own go routine that constantly looks for new messages from the server
func (f *Frontend) SendResult() string {
	log.Printf("SendResult called")
	var response string

	// Request the current highest bid from the server
	for _, client := range f.clients {

		result, err := client.SendResult(context.Background(), &AuctionHouse.Empty{})
		if err != nil {
			log.Fatalf("Error starting result stream: %v", err)
		}

		var resultInfo = result.GetResultResponse()

		response = fmt.Sprintf("Auction status: %s, Current highest bid: %d, Bidder: %s\n", resultInfo.Status, resultInfo.Result, resultInfo.WinnerName)

		// Print the result
		log.Printf("Result received from server:")

	}

	// Return general response here
	return response
}

func (f *Frontend) ListenForWinner() {
	log.Printf("Listen for winner called")
	for _, client := range f.clients {
		for {
			stream, err := client.SendBid(context.Background())
			if err != nil {
				log.Fatalf("Error starting bid stream: %v", err)
			}

			log.Printf("Listening for winner result:")
			serverResponse, err := stream.Recv()
			if err != nil {
				log.Printf("Error receiving message: %v", err)
				return
			}
			log.Printf("Received a winner result:")

			switch resp := serverResponse.Response.(type) {
			case *AuctionHouse.GeneralResponse_ResultResponse:
				fmt.Printf("The auction has ended! Winning bid is: %d from Bidder: %s\n", resp.ResultResponse.Result, resp.ResultResponse.WinnerName)
			default:
				log.Printf("Received unexpected response type")
			}
		}
	}
}
