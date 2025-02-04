package webhandler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	v1 "github.com/yourusername/proto-buf-experiment/gen/go/calculator/v1"
)

type WebHandler struct {
	calculationClient v1.AdditionServiceClient
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

func NewWebHandler(calculationClient v1.AdditionServiceClient) *WebHandler {
	return &WebHandler{
		calculationClient: calculationClient,
	}
}

func (h *WebHandler) AddHandler(w http.ResponseWriter, r *http.Request) {
	// Set content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Decode request body
	var req AddRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Prepare gRPC request with optional constraints
	addRequest := &v1.AddRequest{
		Numbers:   req.Numbers,
		RequestId: uuid.New().String(),
	}

	// Add optional constraints if provided
	if req.MinValue != nil || req.MaxValue != nil || req.MaxNumbers != nil {
		addRequest.Constraints = &v1.AddRequest_Constraints{
			MinValue:   req.MinValue,
			MaxValue:   req.MaxValue,
			MaxNumbers: req.MaxNumbers,
		}
	}

	// Call gRPC service
	resp, err := h.calculationClient.Add(r.Context(), addRequest)
	if err != nil {
		http.Error(w, "Error performing addition", http.StatusInternalServerError)
		return
	}

	// Prepare response
	response := AddResponse{
		Result:    resp.Result,
		RequestID: resp.RequestId,
	}

	// Handle potential error
	if resp.Error != nil {
		response.Error = &ErrorInfo{
			Code:     resp.Error.Code,
			Message:  resp.Error.Message,
			Severity: resp.Error.Severity.String(),
		}
	}

	// Add calculation metadata if available
	if resp.CalculationMetadata != nil {
		response.CalculationMetadata = &CalcMetadata{
			CalculationTime:   resp.CalculationMetadata.CalculationTime.AsTime().String(),
			NumbersProcessed:  resp.CalculationMetadata.NumbersProcessed,
			CalculationMethod: resp.CalculationMetadata.CalculationMethod,
		}
	}

	// Send JSON response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
