package core

import (
	"context"
	"net/http"
)

// Client defines an interface for making HTTP requests. It provides methods to initiate a
// new request and retrieve the underlying HTTP client used to send those requests.
type Client interface {
	// R returns a new Request object that can be used to configure and send an HTTP request.
	// It allows setting headers, body, query parameters, and other request-specific data.
	R() Request

	// GetClient returns the underlying *http.Client used to make the requests.
	// This allows access to the actual HTTP client, which may be useful for configuring
	// connection settings like timeouts, retries, or transport configurations.
	GetClient() *http.Client
}

// Request defines an interface for configuring and sending an HTTP request.
// It provides methods to set authentication tokens, request bodies, headers, query parameters,
// and HTTP methods such as GET, POST, PATCH, DELETE, and PUT.
type Request interface {
	// SetAuthToken sets an authorization token for the request, typically used in the
	// Authorization header to authenticate the request.
	//
	// Parameters:
	//   - token: The token to be included in the request's Authorization header.
	//
	// Returns:
	//   - The updated Request object, allowing for method chaining.
	SetAuthToken(token string) Request

	// SetBody sets the body of the request, allowing you to send data with the request.
	// This can be used for POST, PUT, or PATCH requests that require data to be sent.
	//
	// Parameters:
	//   - body: The body data to send with the request.
	//
	// Returns:
	//   - The updated Request object, allowing for method chaining.
	SetBody(body any) Request

	// SetFormData sets form data for the request, typically used with POST requests.
	// The data will be sent as "application/x-www-form-urlencoded".
	//
	// Parameters:
	//   - formData: A map of form field names to values.
	//
	// Returns:
	//   - The updated Request object, allowing for method chaining.
	SetFormData(formData map[string]string) Request

	// SetHeader sets a header for the request.
	// Headers can include metadata such as content type, authorization, or custom headers.
	//
	// Parameters:
	//   - header: The name of the header to set.
	//   - value: The value for the header.
	//
	// Returns:
	//   - The updated Request object, allowing for method chaining.
	SetHeader(header, value string) Request

	// SetQueryParams sets query parameters for the request's URL.
	// These parameters will be appended to the URL when the request is sent.
	//
	// Parameters:
	//   - params: A map of query parameter names to values.
	//
	// Returns:
	//   - The updated Request object, allowing for method chaining.
	SetQueryParams(params map[string][]string) Request

	// SetPath sets the path for the URL of the request. This is typically used to format
	// the URL dynamically, inserting values into the path.
	//
	// Parameters:
	//   - format: A string format for the path.
	//   - a: Any additional arguments to be used for formatting the path.
	//
	// Returns:
	//   - The updated Request object, allowing for method chaining.
	SetPath(format string, a ...any) Request

	// Get sends a GET request to the server with the specified configuration.
	// It retrieves data from the server.
	//
	// Parameters:
	//   - ctx: The context to manage the request's lifecycle, such as timeouts or cancellations.
	//
	// Returns:
	//   - A Response object representing the server's response, or an error if the request fails.
	Get(ctx context.Context) (Response, error)

	// Post sends a POST request to the server with the specified configuration.
	// This is typically used to send data to the server.
	//
	// Parameters:
	//   - ctx: The context to manage the request's lifecycle, such as timeouts or cancellations.
	//
	// Returns:
	//   - A Response object representing the server's response, or an error if the request fails.
	Post(ctx context.Context) (Response, error)

	// Patch sends a PATCH request to the server with the specified configuration.
	// This is typically used to partially update data on the server.
	//
	// Parameters:
	//   - ctx: The context to manage the request's lifecycle, such as timeouts or cancellations.
	//
	// Returns:
	//   - A Response object representing the server's response, or an error if the request fails.
	Patch(ctx context.Context) (Response, error)

	// Delete sends a DELETE request to the server with the specified configuration.
	// This is typically used to remove data from the server.
	//
	// Parameters:
	//   - ctx: The context to manage the request's lifecycle, such as timeouts or cancellations.
	//
	// Returns:
	//   - A Response object representing the server's response, or an error if the request fails.
	Delete(ctx context.Context) (Response, error)

	// Put sends a PUT request to the server with the specified configuration.
	// This is typically used to update data on the server.
	//
	// Parameters:
	//   - ctx: The context to manage the request's lifecycle, such as timeouts or cancellations.
	//
	// Returns:
	//   - A Response object representing the server's response, or an error if the request fails.
	Put(ctx context.Context) (Response, error)
}

// Response defines an interface for the server's response to an HTTP request.
// It provides a method to access the raw body of the response.
type Response interface {
	// Body retrieves the raw body content of the server's response.
	// It returns the body as a slice of bytes, which can be processed further
	// as needed, such as deserialization into a specific data structure.
	Body() []byte
}
