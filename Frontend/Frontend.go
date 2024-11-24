package frontend

import (
	"context"
	"fmt"
	"log"
	AuctionHouse "main/Handin5"
	"sync"
	"time"

	"google.golang.org/grpc"
)

// maybe context cancellation for the streams
var clientId string = ""
var bidfromclient int32 = 0
var mu sync.Mutex

type Frontend struct {
	clients []AuctionHouse.AuctionServiceClient
	streams map[AuctionHouse.AuctionServiceClient]grpc.BidiStreamingClient[AuctionHouse.UserBidRequest, AuctionHouse.GeneralResponse]
}

func NewFrontend() *Frontend {
	NodeOneAddress := "localhost:50051"   // Address to the server
	NodeTwoAddress := "localhost:50052"   // Address to the server
	NodeThreeAddress := "localhost:50053" // Address to the server
	nodeAddresses := []string{NodeOneAddress, NodeTwoAddress, NodeThreeAddress}
	// nodeAddresses := []string{NodeOneAddress}

	var clients []AuctionHouse.AuctionServiceClient
	// streams := make(map[AuctionHouse.AuctionServiceClient]grpc.BidiStreamingClient[AuctionHouse.UserBidRequest, AuctionHouse.GeneralResponse])
	streams := make(map[AuctionHouse.AuctionServiceClient]AuctionHouse.AuctionService_SendBidClient)

	for _, address := range nodeAddresses {
		conn, err := grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			//log.Fatalf("Failed to connect to server: %v", err)
		}
		client := AuctionHouse.NewAuctionServiceClient(conn)
		clients = append(clients, client)

		// Ensures that we reuse the same stream for both bidding and listening.
		stream, err := client.SendBid(context.Background())
		if err != nil {
			//log.Fatalf("Error creating stream: %v", err)
		}
		// lock until all streams have been inserted
		mu.Lock()
		streams[client] = stream
		mu.Unlock()
	}
	return &Frontend{clients: clients, streams: streams}
}

func (f *Frontend) SendBid(bidRequest *AuctionHouse.UserBidRequest) string {
	log.Printf("SendBid called")
	var response string
	var serverResponseSlice = make([]*AuctionHouse.GeneralResponse, 0)
	var ackCounter int = 0
	var wg sync.WaitGroup

	for _, client := range f.clients {
		wg.Add(1)
		go func(client AuctionHouse.AuctionServiceClient) {
			defer wg.Done()
			log.Printf("Creating new stream for client: %v", client)
			stream, err := client.SendBid(context.Background())
			if err != nil {
				log.Printf("Error creating stream: %v", err)
				return
			}

			err = stream.Send(&AuctionHouse.UserBidRequest{
				BidAmount:  bidRequest.BidAmount,
				BidderName: bidRequest.BidderName,
			})
			if err != nil {
				log.Printf("Error sending bid: %v", err)
				return
			}

			// Set a timeout for receiving the response
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			ch := make(chan *AuctionHouse.GeneralResponse, 1)
			go func() {
				serverResponse, err := stream.Recv()
				if err != nil {
					log.Printf("Error receiving bid: %v", err)
					ch <- nil
				} else {
					ch <- serverResponse
				}
			}()

			select {
			case res := <-ch:
				if res != nil {
					mu.Lock()
					serverResponseSlice = append(serverResponseSlice, res)
					mu.Unlock()
					log.Printf("Appending response to slice")
				}
			case <-ctx.Done():
				log.Printf("Timeout waiting for response from server: %v", client)
			}
		}(client)
	}

	wg.Wait()

	for _, sliceResponse := range serverResponseSlice {
		switch resp := sliceResponse.Response.(type) {
		case *AuctionHouse.GeneralResponse_BidResponse:
			if resp.BidResponse.Ack {
				log.Printf("Received ACK from server")
				ackCounter++
			} else {
				log.Printf("Received NACK from server")
				ackCounter--
			}
		case *AuctionHouse.GeneralResponse_ResultResponse:
			response = fmt.Sprintf("The auction has ended! Winning bid is: %d from Bidder: %s\n", resp.ResultResponse.Result, resp.ResultResponse.WinnerName)
		}
	}

	log.Printf("Size of slice is: %d", len(serverResponseSlice))

	if ackCounter > 0 {
		response = "[YOUR BID WAS ACCEPTED]\n"
	} else if ackCounter < 0 {
		response = "[YOUR BID WAS TOO LOW]\n"
	} else if ackCounter == 0 {
		response = "[RECEIVED CONFLICTING ACK FROM SERVERS]\n"
	}

	return response
}

// Method with own go routine that constantly looks for new messages from the server
func (f *Frontend) SendResult() string {
	log.Printf("SendResult called")
	var response string
	var serverResponseSlice = make([]*AuctionHouse.GeneralResponse, 0)
	responseChan := make(chan *AuctionHouse.GeneralResponse, len(f.clients))

	// Request the current highest bid from the server
	for _, client := range f.clients {
		go func(client AuctionHouse.AuctionServiceClient) {
			result, err := client.SendResult(context.Background(), &AuctionHouse.Empty{})
			if err != nil {
				log.Printf("Error starting result stream: %v", err)
			}
			responseChan <- result
		}(client)
	}

	// Collect responses from all servers
	for i := 0; i < len(f.clients); i++ {
		serverResponse := <-responseChan
		serverResponseSlice = append(serverResponseSlice, serverResponse)
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
	for {
		var serverResponseSlice = make([]*AuctionHouse.GeneralResponse, 0)

		// Collect responses from all servers
		for _, client := range f.clients {
			stream := f.streams[client]
			serverResponse, err := stream.Recv()
			if err != nil {
				log.Printf("Error receiving message from server: %v", err)
				continue
			}
			serverResponseSlice = append(serverResponseSlice, serverResponse)
		}

		// Count responses and find majority
		responseCount := make(map[string]int)
		resultMap := make(map[string]*AuctionHouse.ResultResponse)
		var majorityResult *AuctionHouse.ResultResponse
		var maxCount int

		// Process each response
		for _, response := range serverResponseSlice {
			if resp, ok := response.Response.(*AuctionHouse.GeneralResponse_ResultResponse); ok {
				resultKey := fmt.Sprintf("%d-%s-%s",
					resp.ResultResponse.Result,
					resp.ResultResponse.Status,
					resp.ResultResponse.WinnerName)

				responseCount[resultKey]++
				resultMap[resultKey] = resp.ResultResponse

				// Update majority if this response has more occurrences
				if responseCount[resultKey] > maxCount {
					maxCount = responseCount[resultKey]
					majorityResult = resp.ResultResponse
				}
			}
		}

		// Only announce if we have a clear majority
		if majorityResult != nil && maxCount > len(f.clients)/2 {
			fmt.Printf("[MAJORITY AGREEMENT] The auction has ended! "+
				"Winning bid: %d from Bidder: %s (Status: %s)\n",
				majorityResult.Result,
				majorityResult.WinnerName,
				majorityResult.Status)
		} else {
			log.Printf("[WARNING] No majority agreement on winner")
		}
	}
}
