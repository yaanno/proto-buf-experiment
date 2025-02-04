package calculation

import (
	"context"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	v1 "github.com/yourusername/proto-buf-experiment/gen/go/calculator/v1"
	"github.com/yourusername/proto-buf-experiment/services/calculation/service"
)

func TestAdditionService_Add(t *testing.T) {
	additionService := service.NewAdditionService()

	testCases := []struct {
		name           string
		request        *v1.AddRequest
		expectedResult float64
		expectedError  bool
		errorSeverity  v1.AddResponse_ErrorInfo_Severity
	}{
		{
			name: "Basic addition",
			request: &v1.AddRequest{
				Numbers:   []float64{1.0, 2.0, 3.0},
				RequestId: "test-request-1",
			},
			expectedResult: 6.0,
			expectedError:  false,
		},
		{
			name: "Empty input",
			request: &v1.AddRequest{
				Numbers:   []float64{},
				RequestId: "test-request-2",
			},
			expectedResult: 0,
			expectedError:  true,
			errorSeverity:  v1.AddResponse_ErrorInfo_SEVERITY_WARNING,
		},
		{
			name: "Constraints - Max Numbers",
			request: &v1.AddRequest{
				Numbers:   []float64{1.0, 2.0, 3.0, 4.0},
				RequestId: "test-request-3",
				Constraints: &v1.AddRequest_Constraints{
					MaxNumbers: intPtr(3),
				},
			},
			expectedResult: 0,
			expectedError:  true,
			errorSeverity:  v1.AddResponse_ErrorInfo_SEVERITY_WARNING,
		},
		{
			name: "Constraints - Min Value",
			request: &v1.AddRequest{
				Numbers:   []float64{1.0, 2.0, -3.0},
				RequestId: "test-request-4",
				Constraints: &v1.AddRequest_Constraints{
					MinValue: floatPtr(0.0),
				},
			},
			expectedResult: 0,
			expectedError:  true,
			errorSeverity:  v1.AddResponse_ErrorInfo_SEVERITY_ERROR,
		},
		{
			name: "Constraints - Max Value",
			request: &v1.AddRequest{
				Numbers:   []float64{1.0, 2.0, 100.0},
				RequestId: "test-request-5",
				Constraints: &v1.AddRequest_Constraints{
					MaxValue: floatPtr(10.0),
				},
			},
			expectedResult: 0,
			expectedError:  true,
			errorSeverity:  v1.AddResponse_ErrorInfo_SEVERITY_ERROR,
		},
		{
			name:     "Basic Addition",
			numbers:  []float64{1.0, 2.0, 3.0},
			expected: 6.0,
			hasError: false,
		},
		{
			name:     "Negative Numbers",
			numbers:  []float64{-1.0, 2.0, -3.0},
			expected: -2.0,
			hasError: false,
		},
		{
			name:     "Empty Input",
			numbers:  []float64{},
			expected: 0,
			hasError: true,
			errorMsg: "no numbers provided",
		},
		{
			name:     "Single Number",
			numbers:  []float64{5.5},
			expected: 5.5,
			hasError: false,
		},
		{
			name:     "Floating Point Precision",
			numbers:  []float64{0.1, 0.2},
			expected: 0.3,
			hasError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := additionService.Add(context.Background(), tc.request)

			if tc.expectedError {
				require.Error(t, err)
				require.NotNil(t, resp)
				require.NotNil(t, resp.Error)
				assert.Equal(t, tc.errorSeverity, resp.Error.Severity)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				assert.Equal(t, tc.expectedResult, resp.Result)
				assert.NotEmpty(t, resp.RequestId)
				assert.NotNil(t, resp.CalculationMetadata)
				assert.Equal(t, int32(len(tc.request.Numbers)), resp.CalculationMetadata.NumbersProcessed)
			}

			// Verify calculation metadata
			if resp.CalculationMetadata != nil {
				assert.NotNil(t, resp.CalculationMetadata.CalculationTime)
				assert.Equal(t, "simple_addition", resp.CalculationMetadata.CalculationMethod)
			}
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

func TestAdditionService_Overflow(t *testing.T) {
	additionService := service.NewAdditionService()

	largeNumbers := make([]float64, 1000)
	for i := range largeNumbers {
		largeNumbers[i] = math.MaxFloat64
	}

	req := &v1.AddRequest{
		Numbers:   largeNumbers,
		RequestId: "overflow-test",
	}

	resp, err := additionService.Add(context.Background(), req)

	require.Error(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Error)
	assert.Equal(t, v1.AddResponse_ErrorInfo_SEVERITY_CRITICAL, resp.Error.Severity)
	assert.Equal(t, "OVERFLOW", resp.Error.Code)
}

func TestAdditionService_RequestIDGeneration(t *testing.T) {
	calculationService := service.NewAdditionService()

	req := &v1.AddRequest{
		Numbers: []float64{1.0, 2.0},
		// Intentionally left empty to test auto-generation
		RequestId: "",
	}

	resp, err := calculationService.Add(context.Background(), req)

	assert.NoError(t, err)
	assert.NotEmpty(t, resp.RequestId, "Request ID should be auto-generated")
}

func BenchmarkAdditionService_Add(b *testing.B) {
	calculationService := service.NewAdditionService()
	req := &v1.AddRequest{
		Numbers:   []float64{1.0, 2.0, 3.0, 4.0, 5.0},
		RequestId: "bench-request-id",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = calculationService.Add(context.Background(), req)
	}
}
