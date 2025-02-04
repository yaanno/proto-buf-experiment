package main

import (
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/yourusername/proto-buf-experiment/gen/go/calculator/v1"
	"github.com/yourusername/proto-buf-experiment/pkg/logging"
	"github.com/yourusername/proto-buf-experiment/services/calculation/internal/service"
)

const (
	port = ":50051"
)

func main() {
	// Create logger
	logger := logging.NewLogger(logging.LogConfig{
		ServiceName: "calculation-service",
		Debug:       os.Getenv("DEBUG") == "true",
		WriteToFile: true,
	})

	// Log service startup
	logger.Info().Msg("Starting calculation service")

	// Create a listener on TCP port
	lis, err := net.Listen("tcp", port)
	if err != nil {
		logger.Error().
			Err(err).
			Str("port", port).
			Msg("Failed to create listener")
		os.Exit(1)
	}

	// Create a gRPC server object with logging interceptor
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(logging.UnaryServerInterceptor(logger)),
	)

	// Attach the AdditionService implementation
	calculationService := service.NewAdditionService()
	pb.RegisterAdditionServiceServer(grpcServer, calculationService)

	// Register reflection service on gRPC server
	reflection.Register(grpcServer)

	// Log service start
	logger.Info().
		Str("port", port).
		Msg("Calculation service listening")

	// Start gRPC server
	if err := grpcServer.Serve(lis); err != nil {
		logger.Error().
			Err(err).
			Msg("Failed to serve gRPC server")
		os.Exit(1)
	}
}
