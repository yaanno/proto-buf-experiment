package service

import (
	"context"

	"github.com/yourusername/proto-buf-experiment/gen/go/calculator/v1"
	internalService "github.com/yourusername/proto-buf-experiment/services/calculation/internal/service"
)

// AdditionService provides a public wrapper around the internal service implementation
type AdditionService struct {
	v1.UnimplementedAdditionServiceServer
	internalService *internalService.AdditionService
}

// NewAdditionService creates a new instance of the public AdditionService
func NewAdditionService() *AdditionService {
	return &AdditionService{
		internalService: internalService.NewAdditionService(),
	}
}

// Add delegates the addition operation to the internal service
func (s *AdditionService) Add(ctx context.Context, req *v1.AddRequest) (*v1.AddResponse, error) {
	return s.internalService.Add(ctx, req)
}
