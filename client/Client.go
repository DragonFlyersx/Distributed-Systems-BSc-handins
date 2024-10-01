package main

import (
	"context"
	"log"
	"os"
	"time"

	"bob/Exercise05"

	"google.golang.org/grpc"
)

func getTimeFromService(address string) (string, time.Duration, error) {
	// Start tracking time for round-trip measurement
	start := time.Now()

	// Connect to the gRPC server
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return "", 0, err
	}
	defer conn.Close()

	// Create a client
	client := Exercise05.NewTimeServiceClient(conn)

	// Call the GetTime method
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.GetCurrentTime(ctx, &Exercise05.TimeRequest{})
	if err != nil {
		return "", 0, err
	}

	// Measure the round-trip time
	rtt := time.Since(start)

	return resp.CurrentTime, rtt, nil
}

func main() {
	// List of service addresses
	services := []string{
		"10.26.18.67:50051", // Replace with actual IPs and ports
	}

	// Open a log file to record the times
	logFile, err := os.OpenFile("time_log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()

	logger := log.New(logFile, "", log.LstdFlags)

	for _, service := range services {
		timeValue, rtt, err := getTimeFromService(service)
		if err != nil {
			log.Printf("Error getting time from %s: %v", service, err)
			continue
		}

		logger.Printf("Time from %s: %s | RTT: %v\n", service, timeValue, rtt)
		log.Printf("Time from %s: %s | RTT: %v\n", service, timeValue, rtt)
	}
}
