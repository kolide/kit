/*
Package version provides utilities for displaying version information about a Go application.

To use this package, a program would set the package variables at build time, using the
-ldflags go build flag.

Example:

	go build -ldflags "-X github.com/kolide/kit/version.version=1.0.0"

Available values and defaults to use with ldflags:

	version   = "unknown"
	branch    = "unknown"
	revision  = "unknown"
	goVersion = "unknown"
	buildDate = "unknown"
	buildUser = "unknown"
	appName   = "unknown"
*/
package version

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

// These values are private which ensures they can only be set with the build flags.
var (
	version   = "unknown"
	branch    = "unknown"
	revision  = "unknown"
	goVersion = runtime.Version()
	buildDate = "unknown"
	buildUser = "unknown"
	appName   = "unknown"
)

// Info is a structure with version build information about the current application.
type Info struct {
	Version   string `json:"version"`
	Branch    string `json:"branch"`
	Revision  string `json:"revision"`
	GoVersion string `json:"go_version"`
	BuildDate string `json:"build_date"`
	BuildUser string `json:"build_user"`
}

// semverRegexp is used to standardize various version strings by pulling out only
// the major.minor.patch submatch
var semverRegexp = regexp.MustCompile(`^v?([0-9]+\.[0-9]+\.[0-9]+).*$`)

// Version returns a structure with the current version information.
func Version() Info {
	return Info{
		Version:   version,
		Branch:    branch,
		Revision:  revision,
		GoVersion: goVersion,
		BuildDate: buildDate,
		BuildUser: buildUser,
	}
}

// Print outputs the application name and version string.
func Print() {
	v := Version()
	fmt.Printf("%s version %s\n", appName, v.Version)
}

// PrintFull prints the application name and detailed version information.
func PrintFull() {
	v := Version()
	fmt.Printf("%s - version %s\n", appName, v.Version)
	fmt.Printf("  branch: \t%s\n", v.Branch)
	fmt.Printf("  revision: \t%s\n", v.Revision)
	fmt.Printf("  build date: \t%s\n", v.BuildDate)
	fmt.Printf("  build user: \t%s\n", v.BuildUser)
	fmt.Printf("  go version: \t%s\n", v.GoVersion)
}

// Handler returns an HTTP Handler which returns JSON formatted version information.
func Handler() http.Handler {
	v := Version()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		enc.Encode(v)
	})
}

// versionMultiplier consts are used to create a comparable and reversible integer version from a semver string
// by using 1 million, 1 thousand, and 1 for each part we avoid collisions as long as all parts are less than 1 thousand
const (
	majorVersionMultiplier int = 1000000
	minorVersionMultiplier int = 1000
	patchVersionMultiplier int = 1
)

// VersionNumFromSemver parses the semver version string to look for only the major.minor.patch portion,
// splits that into 3 parts (disregarding any extra portions), and applies a multiplier to each
// part to generate a total int value representing the semver.
// note that this will generate a sortable, and reversible integer as long as all parts remain less than 1000
// This is currently intended for use in generating comparable versions to set within windows registry entries,
// allowing for an easy "upgrade-only" detection configuration within intune.
// Zero is returned for any case where the version cannot be reliably translated.
// VersionNumFromSemver should be used where the build time version value cannot be controlled- to restrict
// translations to the internally set version, use VersionNum
func VersionNumFromSemver(semver string) int {
	semverMatch := semverRegexp.FindStringSubmatch(semver)
	// expect the leftmost match as semverMatch[0] and the semver substring as semverMatch[1]
	if semverMatch == nil || len(semverMatch) != 2 {
		return 0
	}

	parts := strings.Split(semverMatch[1], ".")
	if len(parts) < 3 {
		return 0
	}

	versionNum := 0
	for i, part := range parts[:3] {
		partNum, err := strconv.Atoi(part)
		if err != nil {
			return 0
		}

		switch i {
		case 0:
			versionNum += (partNum * majorVersionMultiplier)
		case 1:
			versionNum += (partNum * minorVersionMultiplier)
		case 2:
			versionNum += (partNum * patchVersionMultiplier)
		}
	}

	return versionNum
}

// VersionNum returns an integer representing the version value set at build time.
// see VersionNumFromSemver for additional details regarding the general translation process
// and limitations. this will return 0 if version is unset/unknown
func VersionNum() int {
	return VersionNumFromSemver(version)
}

// SemverFromVersionNum provides the inverse functionality of VersionNum, allowing us
// to collect and report the integer version in a readable semver format
func SemverFromVersionNum(versionNum int) string {
	if versionNum == 0 {
		return "0.0.0"
	}

	major := versionNum / majorVersionMultiplier
	remaining := versionNum - (major * majorVersionMultiplier)
	minor := remaining / minorVersionMultiplier
	remaining -= (minor * minorVersionMultiplier)
	// not strictly needed because patchVersionMultiplier is 1 but here because it feels correct
	patch := remaining * patchVersionMultiplier

	return fmt.Sprintf("%d.%d.%d", major, minor, patch)
}
