package httputil

import "net/http"

// Middleware is a chainable decorator for HTTP Handlers.
type Middleware func(http.Handler) http.Handler

// Chain is a helper function for composing middlewares. Requests will
// traverse them in the order they're declared. That is, the first middleware
// is treated as the outermost middleware.
//
// Chain is identical to the go-kit helper for Endpoint Middleware.
func Chain(outer Middleware, others ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(others) - 1; i >= 0; i-- { // reverse
			next = others[i](next)
		}
		return outer(next)
	}
}
