package webhandlertest

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"

	pb "github.com/yourusername/proto-buf-experiment/gen/go/calculator/v1"
	webhandler "github.com/yourusername/proto-buf-experiment/services/web-handler/service"
)

// MockCalculationClient is a mock type for the AdditionServiceClient
type MockCalculationClient struct {
	mock.Mock
}

func (m *MockCalculationClient) Add(ctx context.Context, in *pb.AddRequest, opts ...grpc.CallOption) (*pb.AddResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*pb.AddResponse), args.Error(1)
}

func TestAddHandler_SuccessfulAddition(t *testing.T) {
	// Create mock gRPC client
	mockClient := new(MockCalculationClient)

	// Prepare test data
	numbers := []float64{5.5, 3.7}
	expectedResult := 9.2
	expectedRequestID := "test-request-id"

	// Setup mock expectation
	mockClient.On("Add", mock.Anything, mock.Anything, mock.Anything).
		Return(&pb.AddResponse{
			Result:    expectedResult,
			RequestId: expectedRequestID,
			Error:     "",
		}, nil)

	// Create web handler
	handler := webhandler.NewWebHandler(mockClient)

	// Prepare request body
	reqBody, _ := json.Marshal(map[string][]float64{
		"numbers": numbers,
	})

	// Create HTTP request
	req := httptest.NewRequest(http.MethodPost, "/add", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Call handler
	handler.AddHandler(w, req)

	// Check response
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Parse response body
	var respBody map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&respBody)

	// Verify response
	assert.Equal(t, expectedResult, respBody["result"])
	assert.Equal(t, expectedRequestID, respBody["request_id"])
	assert.Empty(t, respBody["error"])

	// Verify mock expectations
	mockClient.AssertExpectations(t)
}

func TestAddHandler_ErrorHandling(t *testing.T) {
	testCases := []struct {
		name               string
		requestBody        []float64
		mockServerResponse *pb.AddResponse
		expectedStatusCode int
	}{
		{
			name:        "Empty Input",
			requestBody: []float64{},
			mockServerResponse: &pb.AddResponse{
				Result:    0,
				Error:     "no numbers provided",
				RequestId: "error-request-id",
			},
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create mock gRPC client
			mockClient := new(MockCalculationClient)

			// Setup mock expectation
			mockClient.On("Add", mock.Anything, mock.Anything, mock.Anything).
				Return(tc.mockServerResponse, nil)

			// Create web handler
			handler := webhandler.NewWebHandler(mockClient)

			// Prepare request body
			reqBody, _ := json.Marshal(map[string][]float64{
				"numbers": tc.requestBody,
			})

			// Create HTTP request
			req := httptest.NewRequest(http.MethodPost, "/add", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()

			// Call handler
			handler.AddHandler(w, req)

			// Check response
			resp := w.Result()
			assert.Equal(t, tc.expectedStatusCode, resp.StatusCode)

			// Parse response body
			var respBody map[string]interface{}
			json.NewDecoder(resp.Body).Decode(&respBody)

			// Verify response
			assert.Equal(t, tc.mockServerResponse.Result, respBody["result"])
			assert.Equal(t, tc.mockServerResponse.RequestId, respBody["request_id"])
			assert.Equal(t, tc.mockServerResponse.Error, respBody["error"])

			// Verify mock expectations
			mockClient.AssertExpectations(t)
		})
	}
}
