// Package dbutil provides utilities for managing connections to a SQL database.
package dbutil

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
)

type dbConfig struct {
	logger log.Logger
}

// WithLogger configures a logger Option.
func WithLogger(logger log.Logger) Option {
	return func(c *dbConfig) {
		c.logger = logger
	}
}

// Option provides optional configuration for managing DB connections.
type Option func(*dbConfig)

// OpenDB creates a sql.DB connection to the database driver.
// OpenDB uses an exponential backoff timer when attempting to establish a connection,
// only returning after the connection is successful or the number of attempts exceeds
// the maxAttempts value(defaults to 15 attempts).
func OpenDB(driver, dsn string, opts ...Option) (*sql.DB, error) {
	config := &dbConfig{
		logger: log.NewNopLogger(),
	}
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, errors.Wrapf(err, "opening %s connection, dsn=%s", driver, dsn)
	}

	const maxAttempts = 15
	var dbError error
	for attempt := 0; attempt < maxAttempts; attempt++ {
		dbError = db.Ping()
		if dbError == nil {
			// we're connected!
			break
		}
		interval := time.Duration(attempt) * time.Second
		level.Info(config.logger).Log(driver, fmt.Sprintf(
			"could not connect to db: %v, sleeping %v", dbError, interval))
		time.Sleep(interval)
	}
	if dbError != nil {
		return nil, dbError
	}

	return db, nil
}
