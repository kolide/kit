// Package httputil provides utilities on top of the net/http package.
package httputil

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/kolide/kit/tlsutil"
)

// Option configures an HTTP Server.
type Option func(*http.Server)

// WithTLSConfig allows overriding the default TLS Config in the call to NewServer.
func WithTLSConfig(cfg *tls.Config) Option {
	return func(s *http.Server) {
		s.TLSConfig = cfg
	}
}

// NewServer creates an HTTP Server with pre-configured timeouts and a secure TLS Config.
func NewServer(addr string, h http.Handler, opts ...Option) *http.Server {
	srv := http.Server{
		Addr:              addr,
		Handler:           h,
		ReadTimeout:       25 * time.Second,
		WriteTimeout:      40 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       5 * time.Minute,
		MaxHeaderBytes:    1 << 18, // 0.25 MB (262144 bytes)
	}

	for _, opt := range opts {
		opt(&srv)
	}

	// set a strict TLS config by default.
	if srv.TLSConfig == nil {
		srv.TLSConfig = tlsutil.NewConfig()
	}

	return &srv
}
