package calculationtest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	v1 "github.com/yourusername/proto-buf-experiment/gen/go/calculator/v1"
	"github.com/yourusername/proto-buf-experiment/services/calculation/service"
)

func TestAdditionService_Add(t *testing.T) {
	testCases := []struct {
		name     string
		numbers  []float64
		expected float64
		hasError bool
		errorMsg string
	}{
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

	calculationService := service.NewAdditionService()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := &v1.AddRequest{
				Numbers:   tc.numbers,
				RequestId: "test-request-id",
			}

			resp, err := calculationService.Add(context.Background(), req)

			if tc.hasError {
				assert.Error(t, err)
				assert.Contains(t, resp.Error, tc.errorMsg)
				assert.Equal(t, float64(0), resp.Result)
			} else {
				assert.NoError(t, err)
				assert.InDelta(t, tc.expected, resp.Result, 1e-9, "Result should match expected value")
				assert.Equal(t, "test-request-id", resp.RequestId)
				assert.Empty(t, resp.Error)
			}
		})
	}
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
