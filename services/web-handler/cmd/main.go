package main

import (
	"net/http"
	"os"

	"google.golang.org/grpc"

	pb "github.com/yourusername/proto-buf-experiment/gen/go/calculator/v1"
	"github.com/yourusername/proto-buf-experiment/pkg/logging"
	webhandler "github.com/yourusername/proto-buf-experiment/services/web-handler/internal"
)

func main() {
	// Create logger
	logger := logging.NewLogger(logging.LogConfig{
		ServiceName: "web-handler-service",
		Debug:       os.Getenv("DEBUG") == "true",
		WriteToFile: true,
	})

	// Log service startup
	logger.Info().Msg("Starting web handler service")

	// Establish gRPC connection
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		logger.Error().
			Err(err).
			Msg("Failed to connect to calculation service")
		os.Exit(1)
	}
	defer conn.Close()

	// Create gRPC client
	calculationClient := pb.NewAdditionServiceClient(conn)

	// Create web handler
	handler := webhandler.NewWebHandler(calculationClient, logger)

	// Setup routes
	http.HandleFunc("/add", handler.AddHandler)

	// Log server start
	logger.Info().
		Str("port", "8080").
		Msg("Web handler service listening")

	// Start server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		logger.Error().
			Err(err).
			Msg("Failed to start web server")
		os.Exit(1)
	}
}
