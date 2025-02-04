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
	calculationService "github.com/yourusername/proto-buf-experiment/services/calculation/internal/service"
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
		name        string
		numbers     []float64
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Empty Input",
			numbers:     []float64{},
			expectError: true,
			errorMsg:    "no numbers provided",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := &pb.AddRequest{
				Numbers:    tc.numbers,
				RequestId: "integration-error-test-" + tc.name,
			}

			resp, err := client.Add(ctx, req)

			if tc.expectError {
				assert.Error(t, err)
				assert.Contains(t, resp.Error, tc.errorMsg)
				assert.Equal(t, float64(0), resp.Result)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
