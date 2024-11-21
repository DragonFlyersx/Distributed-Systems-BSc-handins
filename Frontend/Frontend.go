package main

import (
	"context"
	"log"
	"sync"

	AuctionHouse "main/Handin5"

	"google.golang.org/grpc"
)

type Frontend struct {
	clients []AuctionHouse.AuctionServiceClient
}

func NewFrontend(addresses []string) *Frontend {
	var clients []AuctionHouse.AuctionServiceClient
	for _, address := range addresses {
		conn, err := grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Failed to connect to server: %v", err)
		}
		client := AuctionHouse.NewAuctionServiceClient(conn)
		clients = append(clients, client)
	}
	return &Frontend{clients: clients}
}

func (f *Frontend) SendBid(bidRequest *AuctionHouse.UserBidRequest) (*AuctionHouse.GeneralResponse, error) {
	var wg sync.WaitGroup
	responses := make(chan *AuctionHouse.GeneralResponse, len(f.clients))
	errors := make(chan error, len(f.clients))

	for _, client := range f.clients {
		wg.Add(1)
		go func(client AuctionHouse.AuctionServiceClient) {
			defer wg.Done()
			stream, err := client.SendBid(context.Background())
			if err != nil {
				errors <- err
				return
			}
			if err := stream.Send(bidRequest); err != nil {
				errors <- err
				return
			}
			response, err := stream.Recv()
			if err != nil {
				errors <- err
				return
			}
			responses <- response
		}(client)
	}

	wg.Wait()
	close(responses)
	close(errors)

	// Determine the correct response based on majority voting
	responseCount := make(map[*AuctionHouse.GeneralResponse]int)
	for response := range responses {
		responseCount[response]++
	}

	var correctResponse *AuctionHouse.GeneralResponse
	maxCount := 0
	for response, count := range responseCount {
		if count > maxCount {
			maxCount = count
			correctResponse = response
		}
	}

	if correctResponse == nil {
		return nil, <-errors
	}

	return correctResponse, nil
}
