package service

import (
	"context"
	"fmt"
	"math"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

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
	// Validate request ID
	requestID := req.RequestId
	if requestID == "" {
		requestID = uuid.New().String()
	}

	// Validate constraints if provided
	if req.Constraints != nil {
		// Check max number of numbers
		if req.Constraints.MaxNumbers != nil && len(req.Numbers) > int(*req.Constraints.MaxNumbers) {
			return &pb.AddResponse{
				RequestId: requestID,
				Error: &pb.AddResponse_ErrorInfo{
					Code:     "CONSTRAINT_VIOLATION",
					Message:  fmt.Sprintf("Too many numbers. Maximum allowed: %d", *req.Constraints.MaxNumbers),
					Severity: pb.AddResponse_ErrorInfo_SEVERITY_WARNING,
				},
			}, fmt.Errorf("too many numbers")
		}

		// Validate min and max values
		for _, num := range req.Numbers {
			if req.Constraints.MinValue != nil && num < *req.Constraints.MinValue {
				return &pb.AddResponse{
					RequestId: requestID,
					Error: &pb.AddResponse_ErrorInfo{
						Code:     "VALUE_TOO_LOW",
						Message:  fmt.Sprintf("Number %f is below minimum %f", num, *req.Constraints.MinValue),
						Severity: pb.AddResponse_ErrorInfo_SEVERITY_ERROR,
					},
				}, fmt.Errorf("number below minimum")
			}

			if req.Constraints.MaxValue != nil && num > *req.Constraints.MaxValue {
				return &pb.AddResponse{
					RequestId: requestID,
					Error: &pb.AddResponse_ErrorInfo{
						Code:     "VALUE_TOO_HIGH",
						Message:  fmt.Sprintf("Number %f is above maximum %f", num, *req.Constraints.MaxValue),
						Severity: pb.AddResponse_ErrorInfo_SEVERITY_ERROR,
					},
				}, fmt.Errorf("number above maximum")
			}
		}
	}

	// Validate request
	if len(req.Numbers) == 0 {
		return &pb.AddResponse{
			RequestId: requestID,
			Error: &pb.AddResponse_ErrorInfo{
				Code:     "NO_NUMBERS",
				Message:  "No numbers provided for addition",
				Severity: pb.AddResponse_ErrorInfo_SEVERITY_WARNING,
			},
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
			RequestId: requestID,
			Error: &pb.AddResponse_ErrorInfo{
				Code:     "OVERFLOW",
				Message:  "Calculation resulted in infinity",
				Severity: pb.AddResponse_ErrorInfo_SEVERITY_CRITICAL,
			},
		}, fmt.Errorf("calculation overflow")
	}

	// Prepare response with calculation metadata
	return &pb.AddResponse{
		Result:    result,
		RequestId: requestID,
		CalculationMetadata: &pb.AddResponse_CalculationMetadata{
			CalculationTime:   timestamppb.Now(),
			NumbersProcessed:  int32(len(req.Numbers)),
			CalculationMethod: "simple_addition",
		},
	}, nil
}
