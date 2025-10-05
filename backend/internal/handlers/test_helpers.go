package handlers

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
)

// PathValueContextKey is the context key used to store path parameters.
// This simulates the internal context key used by Go 1.22+ http.ServeMux.
type PathValueContextKey struct{}

// createRequestWithPathParams creates a new HTTP request with path parameters properly set for testing.
// This simulates the behavior of Go 1.22+ http.ServeMux with pattern matching.
func createRequestWithPathParams(method, path string, body *bytes.Buffer, pathParams map[string]string) *http.Request {
	var req *http.Request
	if body != nil {
		req = httptest.NewRequest(method, path, body)
	} else {
		req = httptest.NewRequest(method, path, nil)
	}

	// Set content type if body is provided
	if body != nil && body.Len() > 0 {
		req.Header.Set("Content-Type", "application/json")
	}

	// Create a context with path parameters
	ctx := req.Context()
	if pathParams != nil {
		ctx = context.WithValue(ctx, PathValueContextKey{}, pathParams)
	}

	return req.WithContext(ctx)
}

// extractPathParam extracts a path parameter from the request context.
// This is used by handlers to simulate r.PathValue() behavior in tests.
func extractPathParam(r *http.Request, key string) string {
	// First try the test context (for tests)
	if pathParams, ok := r.Context().Value(PathValueContextKey{}).(map[string]string); ok {
		if value, exists := pathParams[key]; exists {
			return value
		}
	}
	// Fall back to the actual PathValue (for production)
	return r.PathValue(key)
}
