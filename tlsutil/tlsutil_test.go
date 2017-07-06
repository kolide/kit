package tlsutil

import (
	"crypto/tls"
	"testing"
)

func TestNewConfig(t *testing.T) {
	// default, should have Modern compatibility.
	cfg := NewConfig()
	if have, want := cfg.MinVersion, uint16(tls.VersionTLS12); have != want {
		t.Errorf("have %d, want %d", have, want)
	}

	// test WithProfile
	cfg = NewConfig(WithProfile(Old))
	if have, want := cfg.MinVersion, uint16(tls.VersionSSL30); have != want {
		t.Errorf("have %d, want %d", have, want)
	}
}
