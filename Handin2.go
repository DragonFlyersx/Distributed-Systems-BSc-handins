package main

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup

func main() {

	// Implement the TCP/IP Handshake using threads. This is not realistic (since the protocol should run across a network)
	//but your implementation needs to show that you have a good understanding of the protocol.
	client := Client{
		sequenceNumber: 100,
		// ackChannel:     make(chan int, 1),
		// seqChannel:     make(chan int, 1),

	}

	server := Server{
		sequenceNumber: 300,
		ackChannel:     make(chan int, 64000),
		seqChannel:     make(chan int, 64000),
	}

	wg.Add(2)

	go StartServer(&server)

	go StartClient(&client, &server)

	wg.Wait()

}

type Server struct {
	sequenceNumber int
	ackChannel     chan int
	seqChannel     chan int
}

type Client struct {
	sequenceNumber int
	// ackChannel     chan int
	// seqChannel     chan int
}

func StartClient(client *Client, server *Server) {
	defer wg.Done()
	fmt.Println("Client Established")

	// send client sequence number
	server.seqChannel <- client.sequenceNumber
	fmt.Printf("%s %d\n", "The client has sent the sequence number: ", client.sequenceNumber)

	// recieve back it's own sequence number incremented
	var recievedACK = <-server.ackChannel
	// receive the servers sequence number
	var recievedSEQ = <-server.seqChannel
	if recievedACK == client.sequenceNumber+1 {
		fmt.Printf("%s %d %s\n", "The client has received the correct sequence number: ", recievedACK, " from server")

	} else {
		// Handshake failed
		fmt.Printf("%s %d %s\n", "The client has received the wrong sequence number: ", recievedACK, " from server")
		fmt.Printf("%s %d %s\n", "The client will loop back and restart the handshake")
		// If client receives wrong sequence number, loop back and start handshake over again
		//quit <- true
	}

	// send servers sequence number back incremented
	server.ackChannel <- recievedSEQ + 1
	// send the clients own incremented sequence number back
	server.seqChannel <- recievedACK
	//fmt.Printf("%s %d %s %d\n", "Client has finalized the handshake using the client sequence number: ", recievedACK, " and the server sequence number ", recievedSEQ)
}

func StartServer(server *Server) {
	defer wg.Done()
	fmt.Println("Server Established")

	// Receive
	var recieved = <-server.seqChannel // Receive sequence number from client

	// Send
	server.ackChannel <- recieved + 1          // Sends client sequence number back incremented by 1
	server.seqChannel <- server.sequenceNumber // Sends servers sequence number back

	fmt.Println("Server recieved SYN")

	// Recieve
	var recievedACK = <-server.ackChannel
	recieved = <-server.seqChannel

	if recievedACK == server.sequenceNumber+1 {
		fmt.Println("Server has finalized the handshake using the client sequence number: ", recieved, " and the server sequence number ", recievedACK)
		fmt.Println("Connection Establish")
	} else {
		fmt.Println("Connection Failed")
	}
	// Send receive number + 1 and it's own sequence number
}

// 1 channel til server seq og en channel til client
// først sender client sit id
//når serveren modtager incrementer den client id med 1 og sender server id tilbage
// client incrementer server og sender tilbage
