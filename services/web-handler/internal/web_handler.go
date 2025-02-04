package webhandler

import (
	"encoding/json"
	"net/http"

	v1 "github.com/yourusername/proto-buf-experiment/gen/go/calculator/v1"
)

type WebHandler struct {
	calculationClient v1.AdditionServiceClient
}

type AddRequest struct {
	Numbers []float64 `json:"numbers"`
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

	// Prepare gRPC request
	addRequest := &v1.AddRequest{
		Numbers:    req.Numbers,
		RequestId: "web-request-id",
	}

	// Call gRPC service
	resp, err := h.calculationClient.Add(r.Context(), addRequest)
	if err != nil {
		http.Error(w, "Error performing addition", http.StatusInternalServerError)
		return
	}

	// Prepare response
	response := map[string]interface{}{
		"result":     resp.Result,
		"request_id": resp.RequestId,
		"error":      resp.Error,
	}

	// Send JSON response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
