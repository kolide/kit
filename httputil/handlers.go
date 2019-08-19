package httputil

import (
	"crypto/subtle"
	"net/http"
)

// BasicAuthMiddleware is http middleware to authenticate based on a
// predefined map of usernames and passwords.
func BasicAuthMiddleware(basicauthPairs map[string][]byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok || username == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// username and password must match
		expectedPassword, ok := basicauthPairs[username]
		if !ok || subtle.ConstantTimeCompare([]byte(password), expectedPassword) != 1 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// handoff to the next handler
		next.ServeHTTP(w, r)
	})
}

// RedirectToSecureHandler is a simple handler to redirect to the secure URL.
func RedirectToSecureHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Connection", "close")
		url := r.URL
		url.Scheme = "https"
		url.Host = r.Host
		http.Redirect(w, r, url.String(), http.StatusMovedPermanently)
	})
}
