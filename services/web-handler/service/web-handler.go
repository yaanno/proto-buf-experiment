package service

import (
	pb "github.com/yourusername/proto-buf-experiment/gen/go/calculator/v1"
	internal "github.com/yourusername/proto-buf-experiment/services/web-handler/internal"
)

// NewWebHandler creates a new web handler using the internal implementation
func NewWebHandler(calculationClient pb.AdditionServiceClient) *internal.WebHandler {
	return internal.NewWebHandler(calculationClient)
}
