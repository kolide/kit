package httputil

import (
	"net/http"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
)

type AcmOpt func(*autocert.Manager) error

func WithLetsEncryptStaging() AcmOpt {
	return func(m *autocert.Manager) error {
		m.Client.DirectoryURL = "https://acme-staging.api.letsencrypt.org/directory"
		return nil
	}
}

func WithEmail(e string) AcmOpt {
	return func(m *autocert.Manager) error {
		m.Email = e
		return nil
	}
}

func WithRenewBefore(t time.Duration) AcmOpt {
	return func(m *autocert.Manager) error {
		m.RenewBefore = t
		return nil
	}
}

func WithHttpClient(c *http.Client) AcmOpt {
	return func(m *autocert.Manager) error {
		m.Client.HTTPClient = c
		return nil
	}
}

func NewAutocertManager(cache autocert.Cache, allowedHosts []string, opts ...AcmOpt) (*autocert.Manager, error) {
	m := &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(allowedHosts...),
		Cache:      cache,
		Client:     &acme.Client{},
	}

	for _, opt := range opts {
		if err := opt(m); err != nil {
			return nil, errors.Wrap(err, "applying option to autocert manager")
		}
	}

	return m, nil
}
