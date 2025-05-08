package interceptors

import (
	"context"
	"load-generation-system/internal/metrics"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// GRPCInterceptor is a gRPC interceptor that collects and reports metrics for gRPC requests.
// It tracks total request counts, processed request counts, and request durations.
func GRPCInterceptor(
	ctx context.Context, // The context for the gRPC call
	method string, // The name of the gRPC method being called
	req, reply any, // The request and response objects for the gRPC call
	cc *grpc.ClientConn, // The gRPC client connection
	invoker grpc.UnaryInvoker, // The actual invoker function to call the gRPC method
	opts ...grpc.CallOption, // Additional options for the RPC call
) error {
	// Extract the target (server address) of the client connection
	path := cc.Target()

	// Increment the TotalRequestsCounter metric for every incoming request
	metrics.TotalRequestsCounter.WithLabelValues(
		path,   // Target server address
		method, // Method being called
	).Inc()

	// Record the start time for measuring the request duration
	start := time.Now()

	// Invoke the actual gRPC method
	err := invoker(ctx, method, req, reply, cc, opts...)
	// Calculate the time it took to process the request
	duration := time.Since(start).Seconds()

	// Get the status code (OK, CANCELLED, etc.) from the error, if any
	statusCode := status.Code(err).String()

	// Increment the ProcessedRequestsCounter for the processed request, categorized by status code
	metrics.ProcessedRequestsCounter.WithLabelValues(
		path,       // Target server address
		method,     // Method being called
		statusCode, // The status code of the response
	).Inc()

	// Record the duration of the request in the RequestDurationSecondsHist histogram
	metrics.RequestDurationSecondsHist.WithLabelValues(
		path,       // Target server address
		method,     // Method being called
		statusCode, // The status code of the response
	).Observe(duration)

	// Return the error if any occurred during the invocation
	return err
}
