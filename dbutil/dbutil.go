// Package dbutil provides utilities for managing connections to a SQL database.
package dbutil

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type dbConfig struct {
	logger      log.Logger
	maxAttempts int
}

// WithLogger configures a logger Option.
func WithLogger(logger log.Logger) Option {
	return func(c *dbConfig) {
		c.logger = logger
	}
}

// WithMaxAttempts configures the number of maximum attempts to make
func WithMaxAttempts(maxAttempts int) Option {
	return func(c *dbConfig) {
		c.maxAttempts = maxAttempts
	}
}

// Option provides optional configuration for managing DB connections.
type Option func(*dbConfig)

// OpenDB creates a sql.DB connection to the database driver.
// OpenDB uses a linear backoff timer when attempting to establish a connection,
// only returning after the connection is successful or the number of attempts exceeds
// the maxAttempts value(defaults to 15 attempts).
func OpenDB(driver, dsn string, opts ...Option) (*sql.DB, error) {
	config := &dbConfig{
		logger:      log.NewNopLogger(),
		maxAttempts: 15,
	}

	for _, opt := range opts {
		opt(config)
	}

	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, errors.Wrapf(err, "opening %s connection, dsn=%s", driver, dsn)
	}

	var dbError error
	for attempt := 0; attempt < config.maxAttempts; attempt++ {
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

// dbutil.OpenDBX is similar to dbutil.OpenDB, except it returns a *sqlx.DB from
// the popular github.com/jmoiron/sqlx package.
func OpenDBX(driver, dsn string, opts ...Option) (*sqlx.DB, error) {
	db, err := OpenDB(driver, dsn, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "opening database/sql database connection")
	}

	return sqlx.NewDb(db, driver), nil
}
