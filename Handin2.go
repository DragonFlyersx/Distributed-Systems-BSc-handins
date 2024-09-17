package main

import (
	"fmt"
)

const clientAmount = 1

func main() {
	// Implement the TCP/IP Handshake using threads. This is not realistic (since the protocol should run across a network)
	//but your implementation needs to show that you have a good understanding of the protocol.
	clients := make([]Client, clientAmount)
	for index := range clients {
		clients[index] = Client{
			sequenceNumber: index * 100,
			ackChannel:     make(chan int, 1),
			seqChannel:     make(chan int, 1),
		}
	}
	server := Server{
		sequenceNumber: 300,
		ackChannel:     make(chan int, 64000),
		seqChannel:     make(chan int, 64000),
	}

	// hi Send
	//ack sender tilbage recieve
	// hi sender tilbagge send
	// connection estabilished
}

type Server struct {
	sequenceNumber int
	ackChannel     chan int
	seqChannel     chan int
}

type Client struct {
	sequenceNumber int
	ackChannel     chan int
	seqChannel     chan int
}

func SYN(client *Client) {
	// send client sequence number
	client.seqChannel <- client.sequenceNumber
}
func ACK(client *Client) {
	// recieve
	var recieved = <-client.seqChannel

	// send
	client.seqChannel <- recieved
	client.ackChannel <- recieved + 1
}

func SYNACK(server *Server) {
	for {
		// Receive
		var recieved = <-server.seqChannel

		// Send receive number + 1 and it's own sequence number
		server.seqChannel <- recieved + 1
		server.ackChannel <- server.sequenceNumber
		fmt.Println("Server recieved SYN")
	}
}
