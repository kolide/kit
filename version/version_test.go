package version

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestVersion(t *testing.T) {
	t.Parallel()

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

func Test_VersionNum(t *testing.T) {
	t.Parallel()

	var tests = map[string]struct {
		semver             string
		expectedVersionNum int
	}{
		"empty version": {
			semver:             "",
			expectedVersionNum: 0,
		},
		"unset version": {
			semver:             "unknown",
			expectedVersionNum: 0,
		},
		"basic version": {
			semver:             "1.2.3",
			expectedVersionNum: 1002003,
		},
		"4_part_version": {
			semver:             "1.2.3.4",
			expectedVersionNum: 1002003,
		},
		"max version": {
			semver:             "999.999.999",
			expectedVersionNum: 999999999,
		},
		"semver with leading v": {
			semver:             "v1.1.2",
			expectedVersionNum: 1001002,
		},
		"semver with leading zeros": {
			semver:             "01.01.002",
			expectedVersionNum: 1001002,
		},
		"semver with trailing branch info": {
			semver:             "1.10.3-1-g98Paoe",
			expectedVersionNum: 1010003,
		},
		"semver with leading v and trailing branch info": {
			semver:             "v1.10.3-1-g98Paoe",
			expectedVersionNum: 1010003,
		},
		"zero version": {
			semver:             "0.0.0",
			expectedVersionNum: 0,
		},
	}

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			version = tt.semver
			require.Equal(t, tt.expectedVersionNum, VersionNum())
		})
	}
}

func Test_SemverFromVersionNum(t *testing.T) {
	t.Parallel()

	var tests = map[string]struct {
		versionNum     int
		expectedSemver string
	}{
		"zero version": {
			versionNum:     0,
			expectedSemver: "0.0.0",
		},
		"1.10.3": {
			versionNum:     1010003,
			expectedSemver: "1.10.3",
		},
		"max version": {
			versionNum:     999999999,
			expectedSemver: "999.999.999",
		},
		"1.112.43": {
			versionNum:     1112043,
			expectedSemver: "1.112.43",
		},
	}

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.expectedSemver, SemverFromVersionNum(tt.versionNum))
		})
	}
}

func Test_VersionNumComparisons(t *testing.T) {
	t.Parallel()

	var tests = map[string]struct {
		lesserVersion  string
		greaterVersion string
	}{
		"empty version": {
			lesserVersion:  "",
			greaterVersion: "0.0.1",
		},
		"basic versions": {
			lesserVersion:  "1.2.3",
			greaterVersion: "1.2.4",
		},
		"max versions": {
			lesserVersion:  "999.999.998",
			greaterVersion: "999.999.999",
		},
		"large minor versions, no collisions": {
			lesserVersion:  "v1.999.999",
			greaterVersion: "v2.0.0",
		},
	}

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			version = tt.lesserVersion
			lesserParsed := VersionNum()
			version = tt.greaterVersion
			greaterParsed := VersionNum()
			require.True(t, lesserParsed < greaterParsed,
				fmt.Sprintf("expected %s to parse as lesser than %s. got lesser %d >= greater %d",
					tt.lesserVersion,
					tt.greaterVersion,
					lesserParsed,
					greaterParsed,
				),
			)
		})
	}
}

func Test_VersionNumIsReversible(t *testing.T) {
	t.Parallel()

	var tests = map[string]struct {
		testedVersion string
	}{
		"zero version": {
			testedVersion: "0.0.0",
		},
		"basic version": {
			testedVersion: "1.2.3",
		},
		"max version": {
			testedVersion: "999.999.999",
		},
		"random version": {
			testedVersion: "107.61.10",
		},
		"random version 2": {
			testedVersion: "0.118.919",
		},
	}

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			version = tt.testedVersion
			require.Equal(t, tt.testedVersion, SemverFromVersionNum(VersionNum()))
		})
	}
}
