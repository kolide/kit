// Package pgutil provides utilities for Postgres
package pgutil

import (
	"fmt"
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
