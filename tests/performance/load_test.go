package performancetest

import (
	"context"
	"log"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	v1 "github.com/yourusername/proto-buf-experiment/gen/go/calculator/v1"
)

func TestClientConnection(t *testing.T) {
	// Connection timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Establish gRPC connection
	conn, err := grpc.DialContext(ctx,
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(), // Block until connection is established or timeout
	)
	require.NoError(t, err, "Failed to connect to gRPC server")
	defer conn.Close()

	// Create client
	client := v1.NewAdditionServiceClient(conn)

	// Prepare a simple addition request
	req := &v1.AddRequest{
		Numbers:   []float64{1.0, 2.0, 3.0},
		RequestId: "connection-test-request",
	}

	// Attempt to make a request
	resp, err := client.Add(ctx, req)
	
	// Assert no error and valid response
	assert.NoError(t, err, "Failed to make addition request")
	assert.NotNil(t, resp, "Response should not be nil")
	
	// Log the result for debugging
	if resp != nil {
		log.Printf("Connection Test Result: %v", resp.GetResult())
	}
}

func TestPerformance_ConcurrentAdditions(t *testing.T) {
	// Performance test configurations
	concurrentRequests := 100
	timeout := 5 * time.Second

	// Establish gRPC connection
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx,
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(), // Block until connection is established or timeout
	)
	assert.NoError(t, err)
	defer conn.Close()

	client := v1.NewAdditionServiceClient(conn)

	// Performance test
	var wg sync.WaitGroup
	results := make(chan *v1.AddResponse, concurrentRequests)
	errors := make(chan error, concurrentRequests)

	startTime := time.Now()

	for i := 0; i < concurrentRequests; i++ {
		wg.Add(1)
		go func(requestNum int) {
			defer wg.Done()

			req := &v1.AddRequest{
				Numbers:   []float64{float64(requestNum), 1.0, 2.0},
				RequestId: "perf-test-" + strconv.Itoa(requestNum),
			}

			resp, err := client.Add(ctx, req)
			if err != nil {
				errors <- err
				return
			}
			results <- resp
		}(i)
	}

	// Wait for all requests to complete
	wg.Wait()
	close(results)
	close(errors)

	// Check results
	elapsedTime := time.Since(startTime)
	successfulRequests := len(results)
	errorRequests := len(errors)

	log.Printf("Performance Test Results:")
	log.Printf("Total Requests: %d", concurrentRequests)
	log.Printf("Successful Requests: %d", successfulRequests)
	log.Printf("Failed Requests: %d", errorRequests)
	log.Printf("Total Execution Time: %v", elapsedTime)
	log.Printf("Requests per Second: %.2f", float64(concurrentRequests)/elapsedTime.Seconds())

	// Assertions
	assert.Equal(t, concurrentRequests, successfulRequests+errorRequests)
	assert.LessOrEqual(t, elapsedTime.Seconds(), timeout.Seconds())
}

func BenchmarkAdditionService_ConcurrentLoad(b *testing.B) {
	conn, err := grpc.Dial(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	assert.NoError(b, err)
	defer conn.Close()

	client := v1.NewAdditionServiceClient(conn)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req := &v1.AddRequest{
				Numbers:   []float64{1.0, 2.0, 3.0},
				RequestId: "bench-request",
			}

			_, err := client.Add(context.Background(), req)
			assert.NoError(b, err)
		}
	})
}
