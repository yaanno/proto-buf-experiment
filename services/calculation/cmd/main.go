package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/yourusername/proto-buf-experiment/gen/go/calculator/v1"
	"github.com/yourusername/proto-buf-experiment/services/calculation/internal/service"
)

const (
	port = ":50051"
)

func main() {
	// Create a listener on TCP port
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create a gRPC server object
	grpcServer := grpc.NewServer()

	// Attach the AdditionService implementation
	calculationService := service.NewAdditionService()
	pb.RegisterAdditionServiceServer(grpcServer, calculationService)

	// Register reflection service on gRPC server
	reflection.Register(grpcServer)

	log.Printf("Calculation service listening on %v", port)

	// Start gRPC server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
