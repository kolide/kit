/*
Package env provides utility functions for loading environment variables with defaults.

A common use of the env package is for combining flag with environment variables in a Go program.

Example:

	func main() {
		var (
			flProject = flag.String("http.addr", env.String("HTTP_ADDRESS", ":https"), "HTTP server address")
		)
		flag.Parse()
	}
*/
package env

import (
	"fmt"
	"os"
	"time"
)

// String returns the environment variable value specified by the key parameter,
// otherwise returning a default value if set.
func String(key, def string) string {
	if env, ok := os.LookupEnv(key); ok {
		return env
	}
	return def
}

// Bool returns the environment variable value specified by the key parameter,
// otherwise returning a default value if set.
func Bool(key string, def bool) bool {
	if env := os.Getenv(key); env == "true" || env == "TRUE" || env == "1" {
		return true
	}
	return def
}

// Duration returns the environment variable value specified by the key parameter,
// otherwise returning a default value if set.
// If the time.Duration value cannot be parsed, Duration will exit the program
// with an error status.
func Duration(key string, def time.Duration) time.Duration {
	if env, ok := os.LookupEnv(key); ok {
		t, err := time.ParseDuration(env)
		if err != nil {
			fmt.Fprintf(os.Stderr, "env: parse time.Duration from flag: %s\n", err)
			os.Exit(1)
		}
		return t
	}
	return def
}
