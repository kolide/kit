package version

import (
	"testing"
	"time"
)

func TestVersion(t *testing.T) {
	now := time.Now().String()
	version = "test"
	buildDate = now

	info := Version()

	if have, want := info.Version, version; have != want {
		t.Errorf("have %s, want %s", have, want)
	}

	if have, want := info.BuildDate, now; have != want {
		t.Errorf("have %s, want %s", have, want)
	}

	if have, want := info.BuildUser, "unknown"; have != want {
		t.Errorf("have %s, want %s", have, want)
	}
}
