package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/yourusername/proto-buf-experiment/gen/go/calculator/v1"
)

type WebHandler struct {
	calculationClient pb.AdditionServiceClient
}

type AddRequest struct {
	Numbers []float64 `json:"numbers"`
}

type AddResponse struct {
	Result    float64 `json:"result"`
	Error     string  `json:"error,omitempty"`
	RequestID string  `json:"request_id"`
}

func NewWebHandler(calculationClient pb.AdditionServiceClient) *WebHandler {
	return &WebHandler{
		calculationClient: calculationClient,
	}
}

func (h *WebHandler) AddHandler(w http.ResponseWriter, r *http.Request) {
	// Set content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Parse incoming request
	var req AddRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Generate request ID
	requestID := uuid.New().String()

	// Create gRPC request
	grpcReq := &pb.AddRequest{
		RequestId: requestID,
		Numbers:   req.Numbers,
	}

	// Set context with timeout
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Call gRPC service
	grpcResp, err := h.calculationClient.Add(ctx, grpcReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("gRPC error: %v", err), http.StatusInternalServerError)
		return
	}

	// Prepare response
	resp := AddResponse{
		Result:    grpcResp.Result,
		Error:     grpcResp.Error,
		RequestID: grpcResp.RequestId,
	}

	// Send response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func main() {
	// Set up gRPC connection to Calculation Service
	conn, err := grpc.Dial(
		"localhost:50051", 
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to connect to calculation service: %v", err)
	}
	defer conn.Close()

	// Create gRPC client
	calculationClient := pb.NewAdditionServiceClient(conn)

	// Create web handler
	handler := NewWebHandler(calculationClient)

	// Set up HTTP router
	r := mux.NewRouter()
	r.HandleFunc("/add", handler.AddHandler).Methods("POST")

	// Configure server
	srv := &http.Server{
		Handler:      r,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Web handler service listening on :8080")
	log.Fatal(srv.ListenAndServe())
}
