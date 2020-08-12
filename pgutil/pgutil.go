// Package pgutil provides utilities for Postgres
package pgutil

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

// ConnectionOptions represents the configurable options of a connection to a
// Postgres database
type ConnectionOptions struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type Opts func(*ConnectionOptions)

// Supported SSL modes.
type sslMode string

const (
	SSLBlank      sslMode = ""
	SSLDisable            = "disable"
	SSLAllow              = "allow"
	SSLPrefer             = "prefer"
	SSLRequire            = "require"
	SSLVerifyCa           = "verify-ca"
	SSLVerifyFull         = "verify-full"
)

// WithSSL sets the sslmode parameter for postgresql. See the
// postgresql documentation at
// https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNSTRING
func WithSSL(requestedSslMode sslMode) Opts {
	return func(c *ConnectionOptions) {
		c.SSLMode = string(requestedSslMode)
	}
}

// NewFromURL returns a ConnectionOptions from a given URL. This uses
// a format like `postgres://myuser:mypass@localhost/somedatabase`,
// which is commonly found on hosting platforms.
func NewFromURL(rawurl string, opts ...Opts) (ConnectionOptions, error) {
	var c ConnectionOptions
	parsed, err := url.Parse(rawurl)
	if err != nil {
		return c, errors.Wrap(err, "url parse")
	}

	c = ConnectionOptions{
		Host:    parsed.Host,
		User:    parsed.User.Username(),
		DBName:  strings.TrimPrefix(parsed.Path, "/"),
		SSLMode: string(SSLRequire),
	}

	if pass, ok := parsed.User.Password(); ok {
		c.Password = pass
	}

	// Split the URL host/port into parts
	hostComponents := strings.Split(parsed.Host, ":")
	switch len(hostComponents) {
	case 1:
		c.Host = hostComponents[0]
		c.Port = "5432"
	case 2:
		c.Host = hostComponents[0]
		c.Port = hostComponents[1]
	default:
		return c, errors.Errorf("Could not parse %s as host:port", parsed.Host)
	}

	for _, opt := range opts {
		opt(&c)
	}

	return c, nil
}

// String implements the Stringer interface so that a pgutil.ConnectionOptions
// can be converted into a value key/value connection string
func (c ConnectionOptions) String() string {
	s := fmt.Sprintf(
		"host=%s port=%s dbname=%s sslmode=%s user=%s",
		c.Host, c.Port, c.DBName, c.SSLMode, c.User,
	)

	if c.Password != "" {
		s = fmt.Sprintf("%s password=%s", s, c.Password)
	}
	return s
}
