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

import "os"

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
