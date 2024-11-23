package main

import (
	"bufio"
	"fmt"
	"log"
	Frontend "main/Frontend"
	AuctionHouse "main/Handin5"
	"os"
	"strconv"
	"strings"
)

var clientId string
var bidfromclient int32 = 0
var haveBid bool = false

type UserBidRequest struct {
	BidAmount  int32
	BidderName string
}

func main() {
	frontend := Frontend.NewFrontend()
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("Error getting hostname: %v", err)
	}

	clientId = hostname

	var userCommand string

	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Write 'bid <amount>' to make a bid, or 'result' to get the current result. \n")
	for {
		userCommand, _ = reader.ReadString('\n')
		userCommand = strings.TrimSpace(userCommand)

		if userCommand == "Result" { // Request the current highest bid from the server
			fmt.Print(frontend.SendResult())
		} else if strings.HasPrefix(userCommand, "Bid") { // send Bid to the server
			if !haveBid {
				haveBid = true
				// Start listening for winner announcement
				go frontend.ListenForWinner()
			}
			// Parse the bid amount from the command
			parts := strings.Split(userCommand, " ")
			if len(parts) != 2 {
				log.Println("Invalid command. Use: Bid <amount>\n")
				continue
			}

			bidAmount, err := strconv.Atoi(parts[1])
			if err != nil {
				log.Fatalf("Error: This is not an integer: %v", err)
			}

			// Create an instance of UserBidRequest
			bidRequest := &AuctionHouse.UserBidRequest{
				BidAmount:  int32(bidAmount),
				BidderName: clientId,
			}
			fmt.Print(frontend.SendBid(bidRequest))
		}
	}
}
