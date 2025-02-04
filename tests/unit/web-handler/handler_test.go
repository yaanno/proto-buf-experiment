package webhandlertest

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	v1 "github.com/yourusername/proto-buf-experiment/gen/go/calculator/v1"
	"github.com/yourusername/proto-buf-experiment/pkg/logging"

	webhandler "github.com/yourusername/proto-buf-experiment/services/web-handler/service"
)

// MockAdditionServiceClient is a mock implementation of the AdditionServiceClient
type MockAdditionServiceClient struct {
	mock.Mock
}

func (m *MockAdditionServiceClient) Add(ctx context.Context, in *v1.AddRequest, opts ...grpc.CallOption) (*v1.AddResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*v1.AddResponse), args.Error(1)
}

func TestAddHandler(t *testing.T) {
	testCases := []struct {
		name            string
		requestBody     webhandler.AddRequest
		mockServiceResp *v1.AddResponse
		expectedStatus  int
	}{
		{
			name: "Successful Addition",
			requestBody: webhandler.AddRequest{
				Numbers: []float64{1.0, 2.0, 3.0},
			},
			mockServiceResp: &v1.AddResponse{
				Result:    6.0,
				RequestId: "test-request-id",
				CalculationMetadata: &v1.AddResponse_CalculationMetadata{
					CalculationTime:   timestamppb.New(time.Now()),
					NumbersProcessed:  3,
					CalculationMethod: "simple_addition",
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Addition with Constraints",
			requestBody: webhandler.AddRequest{
				Numbers:    []float64{1.0, 2.0, 3.0},
				MinValue:   floatPtr(0.0),
				MaxValue:   floatPtr(10.0),
				MaxNumbers: int32Ptr(5),
			},
			mockServiceResp: &v1.AddResponse{
				Result:    6.0,
				RequestId: "test-request-id",
				CalculationMetadata: &v1.AddResponse_CalculationMetadata{
					CalculationTime:   timestamppb.New(time.Now()),
					NumbersProcessed:  3,
					CalculationMethod: "simple_addition",
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Error Response",
			requestBody: webhandler.AddRequest{
				Numbers: []float64{1.0, 2.0, 3.0},
			},
			mockServiceResp: &v1.AddResponse{
				Error: &v1.AddResponse_ErrorInfo{
					Code:     "INVALID_INPUT",
					Message:  "Invalid input",
					Severity: v1.AddResponse_ErrorInfo_SEVERITY_ERROR,
				},
				RequestId: "error-request-id",
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create mock gRPC client
			mockClient := new(MockAdditionServiceClient)
			mockClient.On("Add", mock.Anything, mock.Anything, mock.Anything).Return(tc.mockServiceResp, nil)

			// log config
			logConfig := logging.LogConfig{
				ServiceName: "web-handler",
				Debug:       true,
				WriteToFile: true,
			}

			// Create logger
			logger := logging.NewLogger(logConfig)

			// Create web handler
			handler := webhandler.NewWebHandler(mockClient, logger)

			// Prepare request body
			jsonBody, err := json.Marshal(tc.requestBody)
			require.NoError(t, err)

			// Create HTTP request
			req := httptest.NewRequest(http.MethodPost, "/add", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			// Create ResponseRecorder
			w := httptest.NewRecorder()

			// Call handler
			handler.AddHandler(w, req)

			// Get response
			resp := w.Result()
			body, err := ioutil.ReadAll(resp.Body)
			require.NoError(t, err)

			// Check status code
			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			// Parse response body
			var addResp webhandler.AddResponse
			err = json.Unmarshal(body, &addResp)
			require.NoError(t, err)

			// Validate request ID
			if tc.mockServiceResp.RequestId != "" {
				assert.Equal(t, tc.mockServiceResp.RequestId, addResp.RequestID)
			}

			// Validate result for successful cases
			if tc.expectedStatus == http.StatusOK {
				assert.Equal(t, tc.mockServiceResp.Result, addResp.Result)

				// Validate calculation metadata
				require.NotNil(t, addResp.CalculationMetadata)
				assert.Equal(t, tc.mockServiceResp.CalculationMetadata.CalculationMethod, addResp.CalculationMetadata.CalculationMethod)
				assert.Equal(t, tc.mockServiceResp.CalculationMetadata.NumbersProcessed, addResp.CalculationMetadata.NumbersProcessed)
			}

			// Validate error response
			if tc.expectedStatus == http.StatusInternalServerError {
				require.NotNil(t, addResp.Error)
				assert.Equal(t, tc.mockServiceResp.Error.Code, addResp.Error.Code)
				assert.Equal(t, tc.mockServiceResp.Error.Message, addResp.Error.Message)
			}
		})
	}
}

// Helper functions for creating pointers
func int32Ptr(i int32) *int32 {
	return &i
}

func floatPtr(f float64) *float64 {
	return &f
}
