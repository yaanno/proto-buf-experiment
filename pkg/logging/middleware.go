package logging

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// UnaryServerInterceptor creates a logging interceptor for gRPC
func UnaryServerInterceptor(logger Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Start timing
		start := time.Now()

		// Extract request ID from metadata if exists
		var requestID string
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if ids := md.Get("request-id"); len(ids) > 0 {
				requestID = ids[0]
			}
		}

		// Create logger with request context
		logCtx := logger.
			WithRequestID(requestID).
			With().
			Str("method", info.FullMethod).
			Logger()

		// Log request
		logCtx.Info().Msg("Received gRPC request")

		// Call the actual handler
		resp, err := handler(ctx, req)

		// Calculate duration
		duration := time.Since(start)

		// Log response
		if err != nil {
			logCtx.Error().
				Err(err).
				Dur("duration", duration).
				Msg("gRPC request failed")
		} else {
			logCtx.Info().
				Dur("duration", duration).
				Msg("gRPC request completed")
		}

		return resp, err
	}
}
