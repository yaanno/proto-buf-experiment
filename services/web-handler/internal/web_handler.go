package webhandler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/rs/zerolog"

	v1 "github.com/yourusername/proto-buf-experiment/gen/go/calculator/v1"
	"github.com/yourusername/proto-buf-experiment/pkg/logging"
)

type WebHandler struct {
	calculationClient v1.AdditionServiceClient
	logger            zerolog.Logger
}

type AddRequest struct {
	Numbers    []float64 `json:"numbers"`
	MinValue   *float64  `json:"min_value,omitempty"`
	MaxValue   *float64  `json:"max_value,omitempty"`
	MaxNumbers *int32    `json:"max_numbers,omitempty"`
}

type AddResponse struct {
	Result              float64       `json:"result"`
	Error               *ErrorInfo    `json:"error,omitempty"`
	RequestID           string        `json:"request_id"`
	CalculationMetadata *CalcMetadata `json:"calculation_metadata,omitempty"`
}

type ErrorInfo struct {
	Code     string `json:"code"`
	Message  string `json:"message"`
	Severity string `json:"severity"`
}

type CalcMetadata struct {
	CalculationTime   string `json:"calculation_time"`
	NumbersProcessed  int32  `json:"numbers_processed"`
	CalculationMethod string `json:"calculation_method"`
}

func NewWebHandler(
	calculationClient v1.AdditionServiceClient,
	logger logging.Logger,
) *WebHandler {
	return &WebHandler{
		calculationClient: calculationClient,
		logger:            logger.Logger,
	}
}

func (h *WebHandler) AddHandler(w http.ResponseWriter, r *http.Request) {
	// Set content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Decode request body
	var addRequest AddRequest
	if err := json.NewDecoder(r.Body).Decode(&addRequest); err != nil {
		h.logger.Error().
			Err(err).
			Msg("Failed to decode request body")

		response := AddResponse{
			RequestID: "error-request-id",
			Error: &ErrorInfo{
				Code:     "BAD_REQUEST",
				Message:  "Invalid request body",
				Severity: "ERROR",
			},
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Prepare gRPC request
	grpcRequest := &v1.AddRequest{
		Numbers: addRequest.Numbers,
	}

	// Optional: add validation parameters
	if addRequest.MinValue != nil || addRequest.MaxValue != nil || addRequest.MaxNumbers != nil {
		grpcRequest.Constraints = &v1.AddRequest_Constraints{
			MinValue:   addRequest.MinValue,
			MaxValue:   addRequest.MaxValue,
			MaxNumbers: addRequest.MaxNumbers,
		}
	}

	// Perform calculation
	start := time.Now()
	response, err := h.calculationClient.Add(context.Background(), grpcRequest)

	// Log calculation details
	duration := time.Since(start)
	logFields := map[string]interface{}{
		"request_id":     response.RequestId,
		"numbers_count":  len(addRequest.Numbers),
		"duration_ms":    duration.Milliseconds(),
		"calculation_ok": err == nil,
	}

	// Handle gRPC error or service-level error
	if err != nil || response.Error != nil {
		var errorCode, errorMessage string
		var severity string

		if err != nil {
			errorCode = "GRPC_ERROR"
			errorMessage = err.Error()
			severity = "ERROR"
		} else {
			errorCode = response.Error.Code
			errorMessage = response.Error.Message
			severity = response.Error.Severity.String()
		}

		h.logger.Error().
			Str("error_code", errorCode).
			Str("error_message", errorMessage).
			Fields(logFields).
			Msg("Calculation failed")

		httpResponse := AddResponse{
			RequestID: response.RequestId,
			Error: &ErrorInfo{
				Code:     errorCode,
				Message:  errorMessage,
				Severity: severity,
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(httpResponse)
		return
	}

	// Log successful calculation
	h.logger.Info().
		Fields(logFields).
		Msg("Calculation completed successfully")

	// Prepare HTTP response
	httpResponse := AddResponse{
		Result:    response.Result,
		RequestID: response.RequestId,
	}

	// Add calculation metadata if available
	if response.CalculationMetadata != nil {
		httpResponse.CalculationMetadata = &CalcMetadata{
			CalculationTime:   response.CalculationMetadata.CalculationTime.AsTime().String(),
			NumbersProcessed:  response.CalculationMetadata.NumbersProcessed,
			CalculationMethod: response.CalculationMetadata.CalculationMethod,
		}
	}

	// Send response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(httpResponse)
}
