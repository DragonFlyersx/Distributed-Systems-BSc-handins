package main

import (
	"log"
	"main/Handin3"

	"context"

	"google.golang.org/grpc"
)

func c() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := Handin3.NewChittyChatClient(conn)
	client.PublishMessage(context.Background(), &Handin3.ChatMessage{Message: "Hello from the client!"})
	log.Println("Client connected successfully!")
}
