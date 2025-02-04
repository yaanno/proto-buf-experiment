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
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

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
		name               string
		requestBody        interface{}
		mockServiceResp    *v1.AddResponse
		expectedStatusCode int
		expectedResponse   webhandler.AddResponse
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
					NumbersProcessed:  3,
					CalculationMethod: "simple_addition",
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: webhandler.AddResponse{
				Result:    6.0,
				RequestID: "test-request-id",
				CalculationMetadata: &webhandler.CalcMetadata{
					NumbersProcessed:  3,
					CalculationMethod: "simple_addition",
				},
			},
		},
		{
			name: "Addition with Constraints",
			requestBody: webhandler.AddRequest{
				Numbers:    []float64{1.0, 2.0, 3.0},
				MinValue:   floatPtr(0.0),
				MaxValue:   floatPtr(10.0),
				MaxNumbers: intPtr(3),
			},
			mockServiceResp: &v1.AddResponse{
				Result:    6.0,
				RequestId: "test-request-id",
				CalculationMetadata: &v1.AddResponse_CalculationMetadata{
					NumbersProcessed:  3,
					CalculationMethod: "simple_addition",
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: webhandler.AddResponse{
				Result:    6.0,
				RequestID: "test-request-id",
				CalculationMetadata: &webhandler.CalcMetadata{
					NumbersProcessed:  3,
					CalculationMethod: "simple_addition",
				},
			},
		},
		{
			name: "Error Response",
			requestBody: webhandler.AddRequest{
				Numbers: []float64{1.0, 2.0, 3.0},
			},
			mockServiceResp: &v1.AddResponse{
				Result:    0,
				RequestId: "error-request-id",
				Error: &v1.AddResponse_ErrorInfo{
					Code:     "INVALID_INPUT",
					Message:  "Invalid input provided",
					Severity: v1.AddResponse_ErrorInfo_SEVERITY_ERROR,
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: webhandler.AddResponse{
				Result:    0,
				RequestID: "error-request-id",
				Error: &webhandler.ErrorInfo{
					Code:     "INVALID_INPUT",
					Message:  "Invalid input provided",
					Severity: "SEVERITY_ERROR",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create mock client
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
			req, err := http.NewRequest("POST", "/add", bytes.NewBuffer(jsonBody))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()

			// Call handler
			handler.AddHandler(w, req)

			// Check response status code
			assert.Equal(t, tc.expectedStatusCode, w.Code)

			// Parse response body
			var response webhandler.AddResponse
			err = json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			// Compare response
			assert.Equal(t, tc.expectedResponse.Result, response.Result)
			assert.Equal(t, tc.expectedResponse.RequestID, response.RequestID)

			// Check error if present
			if tc.expectedResponse.Error != nil {
				require.NotNil(t, response.Error)
				assert.Equal(t, tc.expectedResponse.Error.Code, response.Error.Code)
				assert.Equal(t, tc.expectedResponse.Error.Message, response.Error.Message)
				assert.Equal(t, tc.expectedResponse.Error.Severity, response.Error.Severity)
			}

			// Check calculation metadata if present
			if tc.expectedResponse.CalculationMetadata != nil {
				require.NotNil(t, response.CalculationMetadata)
				assert.Equal(t, tc.expectedResponse.CalculationMetadata.NumbersProcessed, response.CalculationMetadata.NumbersProcessed)
				assert.Equal(t, tc.expectedResponse.CalculationMetadata.CalculationMethod, response.CalculationMetadata.CalculationMethod)
			}

			// Verify mock expectations
			mockClient.AssertExpectations(t)
		})
	}
}

// Helper functions for creating pointers
func intPtr(i int32) *int32 {
	return &i
}

func floatPtr(f float64) *float64 {
	return &f
}
