package service

import (
	"context"
	"fmt"
	"math"

	"github.com/google/uuid"
	pb "github.com/yourusername/proto-buf-experiment/gen/go/calculator/v1"
)

// AdditionService implements the AdditionService interface
type AdditionService struct {
	pb.UnimplementedAdditionServiceServer
}

// NewAdditionService creates a new instance of AdditionService
func NewAdditionService() *AdditionService {
	return &AdditionService{}
}

// Add performs addition of numbers in the request
func (s *AdditionService) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
	// Validate request
	if len(req.Numbers) == 0 {
		return &pb.AddResponse{
			Result:     0,
			Error:      "no numbers provided for addition",
			RequestId:  req.RequestId,
		}, fmt.Errorf("no numbers provided")
	}

	// Perform addition
	var result float64
	for _, num := range req.Numbers {
		result += num
	}

	// Check for overflow
	if math.IsInf(result, 0) {
		return &pb.AddResponse{
			Result:     0,
			Error:      "calculation resulted in infinity",
			RequestId:  req.RequestId,
		}, fmt.Errorf("calculation overflow")
	}

	// If no request ID was provided, generate one
	requestID := req.RequestId
	if requestID == "" {
		requestID = uuid.New().String()
	}

	return &pb.AddResponse{
		Result:     result,
		RequestId:  requestID,
		Error:      "",
	}, nil
}
