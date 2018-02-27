package version

import (
	"testing"
	"time"
)

func TestVersion(t *testing.T) {
	now := time.Now().String()
	var tests = []struct {
		version   string
		buildDate string
		buildID   string
	}{
		{
			version:   "test",
			buildDate: now,
			buildID:   "42",
		},
		{
			version:   "test_without_build_id",
			buildDate: now,
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			version = tt.version
			buildDate = tt.buildDate
			buildID = tt.buildID

			info := Version()

			if have, want := info.Version, tt.version; have != want {
				t.Errorf("have %s, want %s", have, want)
			}

			if have, want := info.BuildDate, tt.buildDate; have != want {
				t.Errorf("have %s, want %s", have, want)
			}

			if have, want := info.BuildUser, "unknown"; have != want {
				t.Errorf("have %s, want %s", have, want)
			}

			if have, want := info.BuildID, tt.buildID; have != want {
				t.Errorf("have %s, want %s", have, want)
			}
		})
	}

}
