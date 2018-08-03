package ulid

import (
	"crypto/rand"

	"github.com/oklog/ulid"
)

// New returns a Universally Unique Lexicographically Sortable Identifier via
// github.com/oklog/ulid
func New() string {
	return ulid.MustNew(ulid.Now(), rand.Reader).String()
}
