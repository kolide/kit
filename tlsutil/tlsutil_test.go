package tlsutil

import (
	"crypto/tls"
	"testing"
)

func TestNewConfig(t *testing.T) {
	t.Parallel()

	// default, should have Modern compatibility.
	cfg := NewConfig()
	if have, want := cfg.MinVersion, uint16(tls.VersionTLS12); have != want {
		t.Errorf("have %d, want %d", have, want)
	}

	// test WithProfile
	cfg = NewConfig(WithProfile(Old))
	if have, want := cfg.MinVersion, uint16(tls.VersionTLS10); have != want {
		t.Errorf("have %d, want %d", have, want)
	}
}
