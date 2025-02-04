package main

import (
	"log"
	"net/http"

	"google.golang.org/grpc"

	pb "github.com/yourusername/proto-buf-experiment/gen/go/calculator/v1"
	webhandler "github.com/yourusername/proto-buf-experiment/services/web-handler/internal"
)

func main() {
	// Establish gRPC connection
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Create gRPC client
	calculationClient := pb.NewAdditionServiceClient(conn)

	// Create web handler
	handler := webhandler.NewWebHandler(calculationClient)

	// Setup routes
	http.HandleFunc("/add", handler.AddHandler)

	// Start server
	log.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
