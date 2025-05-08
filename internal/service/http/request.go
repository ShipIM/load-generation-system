package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"load-generation-system/internal/core"
	"net/http"
	"net/url"
	"strings"
)

// contextKey defines a custom key type for storing values in the context.
type contextKey string

const (
	// metricPath is the key used to store and retrieve the metric path value in context.
	metricPath contextKey = "metricPath"
)

// httpRequest represents an HTTP request and provides methods to configure and send the request.
type httpRequest struct {
	req          *http.Request       // The underlying HTTP request.
	httpClient   *httpClient         // The HTTP client used to send the request.
	queryParams  map[string][]string // Query parameters to be added to the URL.
	pathTemplate string              // Template for the URL path.
	pathParams   []any               // Parameters to replace placeholders in the URL path template.
}

// SetAuthToken sets the Authorization header with the given token.
//
// Parameters:
//   - token: The authorization token to be set.
//
// Returns:
//   - *httpRequest: The current instance of the request.
func (r *httpRequest) SetAuthToken(token string) core.Request {
	r.req.Header.Set("Authorization", "Bearer "+token)
	return r
}

// SetBody sets the body of the HTTP request as JSON-encoded data.
//
// Parameters:
//   - body: The data to be encoded as JSON and sent as the request body.
//
// Returns:
//   - *httpRequest: The current instance of the request.
func (r *httpRequest) SetBody(body any) core.Request {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		panic(fmt.Sprintf("cannot marshal body to JSON: %v", err))
	}
	r.req.Body = io.NopCloser(bytes.NewReader(jsonBody))
	r.req.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewReader(jsonBody)), nil
	}
	r.req.ContentLength = int64(len(jsonBody))
	r.req.Header.Set("Content-Type", "application/json")

	return r
}

// SetFormData sets the body of the HTTP request as form data.
//
// Parameters:
//   - formData: A map of key-value pairs representing form data.
//
// Returns:
//   - *httpRequest: The current instance of the request.
func (r *httpRequest) SetFormData(formData map[string]string) core.Request {
	form := url.Values{}
	for key, value := range formData {
		form.Add(key, value)
	}
	r.req.Body = io.NopCloser(strings.NewReader(form.Encode()))
	r.req.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(strings.NewReader(form.Encode())), nil
	}
	r.req.ContentLength = int64(len(form.Encode()))
	r.req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return r
}

// SetHeader sets a custom header for the HTTP request.
//
// Parameters:
//   - header: The name of the header to be set.
//   - value: The value of the header.
//
// Returns:
//   - *httpRequest: The current instance of the request.
func (r *httpRequest) SetHeader(header, value string) core.Request {
	r.req.Header.Set(header, value)
	return r
}

// SetQueryParams sets the query parameters for the HTTP request.
//
// Parameters:
//   - params: A map of key-value pairs representing the query parameters.
//
// Returns:
//   - *httpRequest: The current instance of the request.
func (r *httpRequest) SetQueryParams(params map[string][]string) core.Request {
	r.queryParams = params
	return r
}

// SetPath sets the URL path template and path parameters to format the final URL.
//
// Parameters:
//   - template: The URL path template with placeholders for parameters.
//   - a: The parameters to be inserted into the template.
//
// Returns:
//   - *httpRequest: The current instance of the request.
func (r *httpRequest) SetPath(template string, a ...any) core.Request {
	r.pathTemplate = template
	r.pathParams = a
	return r
}

// Get sends a GET request and returns the response.
//
// Parameters:
//   - ctx: The context for the request.
//
// Returns:
//   - core.Response: The response received from the server.
//   - error: Any error that occurred during the request.
func (r *httpRequest) Get(ctx context.Context) (core.Response, error) {
	return r.doRequest(ctx, http.MethodGet)
}

// Post sends a POST request and returns the response.
//
// Parameters:
//   - ctx: The context for the request.
//
// Returns:
//   - core.Response: The response received from the server.
//   - error: Any error that occurred during the request.
func (r *httpRequest) Post(ctx context.Context) (core.Response, error) {
	return r.doRequest(ctx, http.MethodPost)
}

// Patch sends a PATCH request and returns the response.
//
// Parameters:
//   - ctx: The context for the request.
//
// Returns:
//   - core.Response: The response received from the server.
//   - error: Any error that occurred during the request.
func (r *httpRequest) Patch(ctx context.Context) (core.Response, error) {
	return r.doRequest(ctx, http.MethodPatch)
}

// Delete sends a DELETE request and returns the response.
//
// Parameters:
//   - ctx: The context for the request.
//
// Returns:
//   - core.Response: The response received from the server.
//   - error: Any error that occurred during the request.
func (r *httpRequest) Delete(ctx context.Context) (core.Response, error) {
	return r.doRequest(ctx, http.MethodDelete)
}

// Put sends a PUT request and returns the response.
//
// Parameters:
//   - ctx: The context for the request.
//
// Returns:
//   - core.Response: The response received from the server.
//   - error: Any error that occurred during the request.
func (r *httpRequest) Put(ctx context.Context) (core.Response, error) {
	return r.doRequest(ctx, http.MethodPut)
}

// doRequest performs the HTTP request with the specified method and returns the response.
//
// Parameters:
//   - ctx: The context for the request.
//   - method: The HTTP method to be used (GET, POST, etc.).
//
// Returns:
//   - core.Response: The response received from the server.
//   - error: Any error that occurred during the request.
func (r *httpRequest) doRequest(ctx context.Context, method string) (core.Response, error) {
	r.req.Method = method

	// Add the metricPath to the request context.
	pathCtx := context.WithValue(ctx, metricPath, r.pathTemplate)
	r.req = r.req.WithContext(pathCtx)

	// Format the URL using the path template and path parameters.
	urlString := fmt.Sprintf(r.pathTemplate, r.pathParams...)

	// Parse the formatted URL string.
	parsedURL, err := url.Parse(urlString)
	if err != nil {
		return nil, err
	}
	r.req.URL = parsedURL

	// Add query parameters to the URL if they are present.
	if len(r.queryParams) > 0 {
		query := r.req.URL.Query()
		for key, values := range r.queryParams {
			for _, value := range values {
				query.Add(key, value)
			}
		}
		r.req.URL.RawQuery = query.Encode()
	}

	// Perform the HTTP request using the httpClient.
	resp, err := r.httpClient.client.Do(r.req)
	if err != nil {
		if urlErr, ok := err.(*url.Error); ok {
			// Unwrap the error if it's a URL-related error.
			unwrappedErr := urlErr.Unwrap()
			return nil, unwrappedErr
		}
		// Return any other error encountered during the request.
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body.
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Ensure that the status code is within the acceptable range.
	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return nil, core.ErrUnacceptableCode
	}

	// Return the response body as a core.Response.
	return &httpResponse{body: bodyBytes}, nil
}
