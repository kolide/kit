// Package tlsutil provides utilities on top of the standard library TLS package.
package tlsutil

import (
	"crypto/tls"
	"fmt"
)

// Profile represents a collection of TLS CipherSuites and their compatibility with Web Browsers.
// The different profile types are defined on the Mozilla wiki: https://wiki.mozilla.org/Security/Server_Side_TLS
type Profile int

const (
	// Modern CipherSuites only.
	// This configuration is compatible with Firefox 27, Chrome 30, IE 11 on Windows 7,
	// Edge, Opera 17, Safari 9, Android 5.0, and Java 8.
	Modern Profile = iota

	// Intermediate supports a wider range of CipherSuites than Modern and
	// is compatible with Firefox 1, Chrome 1, IE 7, Opera 5 and Safari 1.
	Intermediate

	// Old provides backwards compatibility for legacy clients.
	// Should only be used as a last resort.
	Old
)

func (p Profile) String() string {
	switch p {
	case Modern:
		return "modern"
	case Intermediate:
		return "intermediate"
	case Old:
		return "old"
	default:
		panic("unknown TLS profile constant: " + fmt.Sprintf("%d", p))
	}
}

// Option is a TLS Config option. Options can be provided to the NewConfig function
// when creating a TLS Config.
type Option func(*tls.Config)

// WithProfile overrides the default Profile when creating a new *tls.Config.
func WithProfile(p Profile) Option {
	return func(config *tls.Config) {
		setProfile(config, p)
	}
}

// WithCertificates builds the tls.Config.NameToCertificate from the CommonName and
// SubjectAlternateName fields of the provided certificate.
//
// WithCertificates is useful for creating a TLS Config for servers which require SNI,
// for example reverse proxies.
func WithCertificates(certs []tls.Certificate) Option {
	return func(config *tls.Config) {
		config.Certificates = append(config.Certificates, certs...)
		config.BuildNameToCertificate()
	}
}

// NewConfig returns a configured *tls.Config. By default, the TLS Config is set to
// MinVersion of TLS 1.2 and a Modern Profile.
//
// Use one of the available Options to modify the default config.
func NewConfig(opts ...Option) *tls.Config {
	cfg := tls.Config{PreferServerCipherSuites: true}

	for _, opt := range opts {
		opt(&cfg)
	}

	// if a Profile was not specified, default to Modern.
	if cfg.MinVersion == 0 {
		setProfile(&cfg, Modern)
	}

	return &cfg
}

func setProfile(cfg *tls.Config, profile Profile) {
	switch profile {
	case Modern:
		cfg.MinVersion = tls.VersionTLS12
		cfg.CurvePreferences = append(cfg.CurvePreferences,
			tls.CurveP256,
			tls.CurveP384,
			tls.CurveP521,
			tls.X25519,
		)
		cfg.CipherSuites = append(cfg.CipherSuites,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		)
	case Intermediate:
		cfg.MinVersion = tls.VersionTLS10
		cfg.CurvePreferences = append(cfg.CurvePreferences,
			tls.CurveP256,
			tls.CurveP384,
			tls.CurveP521,
			tls.X25519,
		)
		cfg.CipherSuites = append(cfg.CipherSuites,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA,
			tls.TLS_RSA_WITH_RC4_128_SHA,
			tls.TLS_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
		)
	case Old:
		cfg.MinVersion = tls.VersionSSL30
		cfg.CurvePreferences = append(cfg.CurvePreferences,
			tls.CurveP256,
			tls.CurveP384,
			tls.CurveP521,
			tls.X25519,
		)
		cfg.CipherSuites = append(cfg.CipherSuites,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA,
			tls.TLS_RSA_WITH_RC4_128_SHA,
			tls.TLS_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
		)
	default:
		panic("invalid tls profile " + profile.String())
	}
}
