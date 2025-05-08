package http

import (
	"context"
	"errors"
	"fmt"
	"load-generation-system/internal/core"
	"load-generation-system/internal/metrics"
	"load-generation-system/pkg/utils"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

// httpClient represents an HTTP client with an underlying HTTP client instance.
type httpClient struct {
	client *http.Client
}

// roundTripper wraps the http.RoundTripper interface and allows customization of request/response handling.
type roundTripper struct {
	Transport http.RoundTripper
}

// Regular expression pattern for matching UUIDs in URL paths.
var uuidPattern = regexp.MustCompile(`[a-f0-9\-]{8,}`)

// Constants used for request timeout status.
const timeoutStatus = "Timeout"

func NewClient(
	minIdleConnTimeoutSec, maxIdleConnTimeoutSec int64,
) core.Client {
	// Generate a random idle connection timeout within the given range.
	randomIdleConnTimeout := time.Duration(
		utils.GenerateInt64(minIdleConnTimeoutSec, maxIdleConnTimeoutSec),
	) * time.Second

	// Create a dialer with connection timeout and keep-alive timeout settings.
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 120 * time.Second,
	}

	// Configure the HTTP transport with various timeouts, connection settings, and proxy configuration.
	transport := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          0, // No limit on idle connections.
		IdleConnTimeout:       randomIdleConnTimeout,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		DialContext:           dialer.DialContext,
	}

	// Create and return the HTTP client with the custom transport.
	client := &http.Client{
		Transport: &roundTripper{
			Transport: transport,
		},
	}

	return &httpClient{client: client}
}

// R creates and returns a new core.Request instance for making HTTP requests.
// This method initializes an empty HTTP request and sets up necessary parameters.
func (c *httpClient) R() core.Request {
	// Create a new HTTP request with no body.
	req, err := http.NewRequest("", "", http.NoBody)
	if err != nil {
		// Panic in case of an error while creating the request.
		panic(fmt.Sprintf("cannot create new request: %v", err))
	}
	// Return a new instance of httpRequest with initialized fields.
	return &httpRequest{
		req:         req,
		httpClient:  c,
		queryParams: make(map[string][]string),
	}
}

// GetClient returns the underlying HTTP client.
func (c *httpClient) GetClient() *http.Client {
	return c.client
}

// normalizePath sanitizes the URL path by replacing any UUIDs with a placeholder string '%s'.
// This function is used for anonymizing paths with dynamic UUIDs when tracking metrics.
func normalizePath(path string) string {
	// Replace UUID patterns in the path with '%s' to normalize the URL.
	path = uuidPattern.ReplaceAllString(path, "%s")

	return path
}

// RoundTrip implements the RoundTripper interface, which intercepts the HTTP request,
// tracks metrics, processes the response, and returns it. It also handles errors and timeouts.
//
// Parameters:
//   - req: The HTTP request to be sent.
//
// Returns:
//   - *http.Response: The HTTP response received.
//   - error: Any error encountered during the request/response cycle.
func (rt *roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// Retrieve or construct the metric path for tracking.
	metricPath, ok := req.Context().Value(metricPath).(string)
	if !ok {
		metricPath = fmt.Sprintf("%s://%s%s", req.URL.Scheme, req.URL.Host, normalizePath(req.URL.Path))
	}

	// Increment the TotalRequestsCounter metric for the request.
	metrics.TotalRequestsCounter.WithLabelValues(
		metricPath,
		req.Method,
	).Inc()

	// Start tracking the request duration.
	start := time.Now()

	// Perform the HTTP request using the custom transport.
	resp, err := rt.Transport.RoundTrip(req)
	duration := time.Since(start).Seconds()

	if err != nil {
		// Check if the error is related to timeout or deadline exceeded.
		isDeadline := errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled)
		if err, ok := err.(net.Error); (ok && err.Timeout()) || isDeadline {
			// Record the timeout status in the metrics.
			metrics.ProcessedRequestsCounter.WithLabelValues(
				metricPath,
				req.Method,
				timeoutStatus,
			).Inc()

			// Record the request duration for the timeout.
			metrics.RequestDurationSecondsHist.WithLabelValues(
				metricPath,
				req.Method,
				timeoutStatus,
			).Observe(duration)
		}

		// Return the error if encountered during the round trip.
		return nil, err
	}

	// Record the status code in the metrics.
	status := strconv.Itoa(resp.StatusCode)
	metrics.ProcessedRequestsCounter.WithLabelValues(
		metricPath,
		req.Method,
		status,
	).Inc()

	// Record the request duration with the status code.
	metrics.RequestDurationSecondsHist.WithLabelValues(
		metricPath,
		req.Method,
		status,
	).Observe(duration)

	// Return the response if the round trip was successful.
	return resp, nil
}
