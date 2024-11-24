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
	streams map[AuctionHouse.AuctionServiceClient]grpc.BidiStreamingClient[AuctionHouse.UserBidRequest, AuctionHouse.GeneralResponse]
}

func NewFrontend() *Frontend {
	NodeOneAddress := "localhost:50051" // Address to the server
	//NodeTwoAddress := "localhost:50052"   // Address to the server
	//	NodeThreeAddress := "localhost:50053" // Address to the server
	// nodeAddresses := []string{NodeOneAddress, NodeTwoAddress, NodeThreeAddress}
	nodeAddresses := []string{NodeOneAddress}

	var clients []AuctionHouse.AuctionServiceClient
	// streams := make(map[AuctionHouse.AuctionServiceClient]grpc.BidiStreamingClient[AuctionHouse.UserBidRequest, AuctionHouse.GeneralResponse])
	streams := make(map[AuctionHouse.AuctionServiceClient]AuctionHouse.AuctionService_SendBidClient)

	for _, address := range nodeAddresses {
		conn, err := grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Failed to connect to server: %v", err)
		}
		client := AuctionHouse.NewAuctionServiceClient(conn)
		clients = append(clients, client)

		// Ensures that we reuse the same stream for both bidding and listening.
		stream, err := client.SendBid(context.Background())
		if err != nil {
			log.Fatalf("Error creating stream: %v", err)
		}
		streams[client] = stream
	}
	return &Frontend{clients: clients, streams: streams}
}

func (f *Frontend) SendBid(bidRequest *AuctionHouse.UserBidRequest) string {
	log.Printf("SendBid called")
	var response string
	var serverResponse *AuctionHouse.GeneralResponse
	var serverResponseSlice = make([]*AuctionHouse.GeneralResponse, 0)
	var ackCounter int = 0

	for _, client := range f.clients {
		// Send the bid to the server
		stream := f.streams[client]

		err := stream.Send(&AuctionHouse.UserBidRequest{
			BidAmount:  bidRequest.BidAmount,  // The amount being bid
			BidderName: bidRequest.BidderName, // The name of the bidder
		})
		if err != nil {
			log.Fatalf("Error sending bid: %v", err)
		}

		// Receive the response from the server
		serverResponse, err = stream.Recv()
		if err != nil {
			log.Fatalf("Error receiving bid: %v", err)
		}
		serverResponseSlice = append(serverResponseSlice, serverResponse)
	}

	// Checks responses from all servers
	for _, sliceResponse := range serverResponseSlice {
		switch resp := sliceResponse.Response.(type) {
		case *AuctionHouse.GeneralResponse_BidResponse:
			if resp.BidResponse.Ack {
				// Increment the ackCounter if the bid was accepted
				ackCounter++
			} else {
				// Decrement the ackCounter if the bid was too low to be accepted
				ackCounter--
			}
		case *AuctionHouse.GeneralResponse_ResultResponse:
			response = fmt.Sprintf("The auction has ended! Winning bid is: %d from Bidder: %s\n", resp.ResultResponse.Result, resp.ResultResponse.WinnerName)
		}
	}
	// Check the ackCounter to determine the response
	if ackCounter > 0 {
		response = "[YOUR BID WAS ACCEPTED]\n"
	} else if ackCounter < 0 {
		response = "[YOUR BID WAS TOO LOW]\n"
	} else if ackCounter == 0 {
		// Used when only 2 servers are left, and gives conflicting responses
		response = "[RECEIVED CONFLICTING ACK FROM SERVERS]\n"
	}

	// Clear slice
	serverResponseSlice = make([]*AuctionHouse.GeneralResponse, 0)

	// Return general response here
	return response
}

// Method with own go routine that constantly looks for new messages from the server
func (f *Frontend) SendResult() string {
	log.Printf("SendResult called")
	var response string
	var serverResponseSlice = make([]*AuctionHouse.GeneralResponse, 0)

	// Request the current highest bid from the server
	for _, client := range f.clients {
		result, err := client.SendResult(context.Background(), &AuctionHouse.Empty{})
		if err != nil {
			log.Fatalf("Error starting result stream: %v", err)
		}

		// var resultInfo = result.GetResultResponse()
		serverResponseSlice = append(serverResponseSlice, result)
		// Print the result
		log.Printf("Result received from server:")
	}

	// Determine the majority response
	responseCount := make(map[string]int)
	var majorityResponse *AuctionHouse.ResultResponse

	for _, sliceResponse := range serverResponseSlice {
		switch resp := sliceResponse.Response.(type) {
		case *AuctionHouse.GeneralResponse_ResultResponse:
			resultKey := fmt.Sprintf("%d-%s-%s", resp.ResultResponse.Result, resp.ResultResponse.Status, resp.ResultResponse.WinnerName)
			responseCount[resultKey]++
			if majorityResponse == nil || responseCount[resultKey] > responseCount[fmt.Sprintf("%d-%s-%s", majorityResponse.Result, majorityResponse.Status, majorityResponse.WinnerName)] {
				majorityResponse = resp.ResultResponse
			}
		}
	}

	if majorityResponse != nil {
		response = fmt.Sprintf("Auction status: %s, Current highest bid: %d, Bidder: %s\n", majorityResponse.Status, majorityResponse.Result, majorityResponse.WinnerName)
	} else {
		response = "[RECEIVED CONFLICTING RESPONSES FROM SERVERS]\n"
	}

	return response
}

func (f *Frontend) ListenForWinner() {
	var serverResponse *AuctionHouse.GeneralResponse
	var serverResponseSlice = make([]*AuctionHouse.GeneralResponse, 0)

	log.Printf("Listen for winner called")
	for _, client := range f.clients {
		stream := f.streams[client]
		for {
			serverResponse, err := stream.Recv()
			if err != nil {
				log.Printf("Error receiving message: %v", err)
				return
			}
			log.Printf("Received a winner result:")

			switch resp := serverResponse.Response.(type) {
			case *AuctionHouse.GeneralResponse_ResultResponse:
				fmt.Printf("The auction has ended! Winning bid is: %d from Bidder: %s\n",
					resp.ResultResponse.Result,
					resp.ResultResponse.WinnerName)
			}
		}
	}
}
