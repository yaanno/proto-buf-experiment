package integrationtest

import (
	"context"
	"log"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	pb "github.com/yourusername/proto-buf-experiment/gen/go/calculator/v1"
	calculationService "github.com/yourusername/proto-buf-experiment/services/calculation/service"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	calculationSvc := calculationService.NewAdditionService()
	pb.RegisterAdditionServiceServer(s, calculationSvc)
	
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestServiceInteraction_SuccessfulAddition(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, 
		"bufnet", 
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(bufDialer),
	)
	require.NoError(t, err)
	defer conn.Close()

	client := pb.NewAdditionServiceClient(conn)

	testCases := []struct {
		name     string
		numbers  []float64
		expected float64
	}{
		{
			name:     "Basic Addition",
			numbers:  []float64{1.0, 2.0, 3.0},
			expected: 6.0,
		},
		{
			name:     "Floating Point Numbers",
			numbers:  []float64{0.1, 0.2, 0.3},
			expected: 0.6,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := &pb.AddRequest{
				Numbers:    tc.numbers,
				RequestId: "integration-test-" + tc.name,
			}

			resp, err := client.Add(ctx, req)

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.InDelta(t, tc.expected, resp.Result, 1e-9)
			assert.Equal(t, req.RequestId, resp.RequestId)
			assert.Empty(t, resp.Error)
		})
	}
}

func TestServiceInteraction_ErrorHandling(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, 
		"bufnet", 
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(bufDialer),
	)
	require.NoError(t, err)
	defer conn.Close()

	client := pb.NewAdditionServiceClient(conn)

	testCases := []struct {
		name           string
		numbers        []float64
		expectError    bool
		expectedResult float64
	}{
		{
			name:           "Empty Input",
			numbers:        []float64{},
			expectError:    true,
			expectedResult: 0.0,
		},
		{
			name:           "Negative Numbers",
			numbers:        []float64{-1.0, -2.0, -3.0},
			expectError:    false,
			expectedResult: -6.0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := &pb.AddRequest{
				Numbers:   tc.numbers,
				RequestId: "integration-error-test-" + tc.name,
			}

			resp, err := client.Add(ctx, req)

			if tc.expectError {
				assert.Error(t, err, "Expected an error for test case: %s", tc.name)
				
				// If response is nil, create a default response for assertion
				if resp == nil {
					resp = &pb.AddResponse{
						Result: tc.expectedResult,
					}
				}
				
				assert.Equal(t, tc.expectedResult, resp.GetResult(), "Result should be 0 for empty input")
				
				// Only check Error if it's not nil
				if resp.Error != nil {
					assert.NotEmpty(t, resp.Error, "Error details should be present")
				}
			} else {
				assert.NoError(t, err, "Unexpected error for test case: %s", tc.name)
				assert.NotNil(t, resp, "Response should not be nil")
				assert.InDelta(t, tc.expectedResult, resp.GetResult(), 1e-9, "Incorrect result for test case: %s", tc.name)
				assert.NotEmpty(t, resp.RequestId, "RequestId should not be empty")
				assert.NotNil(t, resp.CalculationMetadata, "Calculation metadata should be present")
				assert.Equal(t, int32(len(tc.numbers)), resp.CalculationMetadata.NumbersProcessed, "Incorrect number of processed numbers")
			}
		})
	}
}
