package service

import (
	pb "github.com/yourusername/proto-buf-experiment/gen/go/calculator/v1"
	"github.com/yourusername/proto-buf-experiment/pkg/logging"
	internal "github.com/yourusername/proto-buf-experiment/services/web-handler/internal"
)

// NewWebHandler creates a new web handler using the internal implementation
func NewWebHandler(calculationClient pb.AdditionServiceClient, logger logging.Logger) *internal.WebHandler {
	return internal.NewWebHandler(calculationClient, logger)
}

// AddRequest represents the request structure for addition operations
type AddRequest struct {
	Numbers    []float64 `json:"numbers"`
	MinValue   *float64  `json:"min_value,omitempty"`
	MaxValue   *float64  `json:"max_value,omitempty"`
	MaxNumbers *int32    `json:"max_numbers,omitempty"`
}

// AddResponse represents the response structure for addition operations
type AddResponse struct {
	Result              float64       `json:"result"`
	Error               *ErrorInfo    `json:"error,omitempty"`
	RequestID           string        `json:"request_id"`
	CalculationMetadata *CalcMetadata `json:"calculation_metadata,omitempty"`
}

// ErrorInfo provides detailed error information
type ErrorInfo struct {
	Code     string `json:"code"`
	Message  string `json:"message"`
	Severity string `json:"severity"`
}

// CalcMetadata provides metadata about the calculation
type CalcMetadata struct {
	CalculationTime   string `json:"calculation_time"`
	NumbersProcessed  int32  `json:"numbers_processed"`
	CalculationMethod string `json:"calculation_method"`
}
