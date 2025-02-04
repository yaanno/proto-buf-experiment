package webhandler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
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

	// Generate request ID
	requestID := uuid.New().String()

	// Create request scoped logger
	requestLogger := h.logger.With().
		Str("request_id", requestID).
		Str("method", r.Method).
		Str("path", r.URL.Path).
		Logger()

	// Log incoming request
	requestLogger.Info().Msg("Received add request")

	// Decode request body
	var addRequest AddRequest
	if err := json.NewDecoder(r.Body).Decode(&addRequest); err != nil {
		requestLogger.Error().
			Err(err).
			Msg("Failed to decode request body")

		response := AddResponse{
			RequestID: requestID,
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
		Numbers:   addRequest.Numbers,
		RequestId: requestID,
	}

	// Perform calculation
	start := time.Now()
	response, err := h.calculationClient.Add(context.Background(), grpcRequest)

	// Log calculation details
	duration := time.Since(start)
	logFields := map[string]interface{}{
		"request_id":     requestID,
		"numbers_count":  len(addRequest.Numbers),
		"duration_ms":    duration.Milliseconds(),
		"calculation_ok": err == nil,
	}

	if err != nil {
		requestLogger.Error().
			Err(err).
			Fields(logFields).
			Msg("Calculation failed")

		httpResponse := AddResponse{
			RequestID: requestID,
			Error: &ErrorInfo{
				Code:     "CALCULATION_ERROR",
				Message:  err.Error(),
				Severity: "ERROR",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(httpResponse)
		return
	}

	// Log successful calculation
	requestLogger.Info().
		Fields(logFields).
		Msg("Calculation completed successfully")

	// Prepare HTTP response
	httpResponse := AddResponse{
		Result:    response.Result,
		RequestID: requestID,
		CalculationMetadata: &CalcMetadata{
			CalculationTime:   duration.String(),
			NumbersProcessed:  int32(len(addRequest.Numbers)),
			CalculationMethod: "addition",
		},
	}

	// Send response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(httpResponse)
}
