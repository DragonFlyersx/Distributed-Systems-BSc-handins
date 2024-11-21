package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	AuctionHouse "main/Handin5"
	"strconv"
	"strings"

	"os"

	"google.golang.org/grpc"
)

var clientId string = ""
var bidfromclient int32 = 0

type UserBidRequest struct {
	BidAmount  int32
	BidderName string
}

// first call to bid registers bidder
// then all should be input new bid
// should be able to query the system for the result / highest bid

func (bidRequest UserBidRequest) SendBid(client AuctionHouse.AuctionServiceClient, BidAmount int32) error {
	bidRequest.BidAmount = BidAmount
	bidRequest.BidderName = clientId

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
		return err
	}

	// Receive the response from the server
	response, err := stream.Recv()
	if err != nil {
		log.Fatalf("Error receiving bid: %v", err)
	}

	// Print the response
	log.Printf("Your bid has been placed: %v", response)

	return nil
}

// Method with own go routine that constantly looks for new messages from the server
func SendResult(client AuctionHouse.AuctionServiceClient) {
	// Request the current highest bid from the server
	result, err := client.SendResult(context.Background(), &AuctionHouse.Empty{})
	if err != nil {
		log.Fatalf("Error starting result stream: %v", err)
	}

	var resultInfo = result.GetResultResponse()
	// Print the result
	log.Printf("Result received from server:")
	log.Printf("Auction status: %s, Current highest bid: %d, Bidder: %s", resultInfo.Status, resultInfo.Result, resultInfo.WinnerName)
}

func main() {
	// Set up a connection to the server
	// needs to connect to 3 server Nodes
	NodeOneAddress := "localhost:50051" // Address to the server
	//NodeTwoAddress := "localhost:50052"   // Address to the server
	//NodeThreeAddress := "localhost:50053" // Address to the server
	/*
		clientId, err := os.Hostname() //
		if err != nil {
			//log.Fatalf("Error getting hostname: %v", err)
			log.Printf("Error getting hostname: %v", err)
			clientId = "Unknown"
		}*/

	nodeAddresses := []string{NodeOneAddress}

	connections := make(map[string]*grpc.ClientConn)
	clients := make(map[string]AuctionHouse.AuctionServiceClient)

	for _, address := range nodeAddresses {
		conn, err := grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect SERVER: %v", err)
		}
		log.Printf("Connected to SERVER: %v", address)
		connections[address] = conn

		// Create a client for each connection
		client := AuctionHouse.NewAuctionServiceClient(conn)
		clients[address] = client
	}

	defer func() {
		for address, conn := range connections {
			log.Printf("Closing connection to SERVER: %v", address)
			conn.Close()
		}
	}()

	var userCommand string

	reader := bufio.NewReader(os.Stdin)

	for _, client := range clients {
		SendResult(client)
	}

	fmt.Printf("Write 'bid <amount>' to make a bid, or 'result' to get the current result. \n")
	for {
		userCommand, _ = reader.ReadString('\n')
		userCommand = strings.TrimSpace(userCommand)

		if userCommand == "Result" { // Request the current highest bid from the server
			for _, client := range clients {
				SendResult(client)
			}

		} else if strings.HasPrefix(userCommand, "Bid") { // send Bid to the server
			// Parse the bid amount from the command
			parts := strings.Split(userCommand, " ")
			if len(parts) != 2 {
				log.Println("Invalid command. Use: Bid <amount>")
				continue
			}

			bidAmount, err := strconv.Atoi(parts[1])
			if err != nil {
				log.Fatalf("Error: This is not an integer: %v", err)
			}

			// Create an instance of UserBidRequest
			var bidRequest UserBidRequest
			for _, client := range clients {
				err = bidRequest.SendBid(client, int32(bidAmount))
				if err != nil {
					log.Fatalf("Error sending bid: %v", err)
				}
			}
		}
	}
}
