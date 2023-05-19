package pgutil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConversion(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		in   string
		opts []Opts
		out  string
		err  bool
	}{
		{
			in:  "postgres://myuser:mypass@localhost/somedatabase",
			out: "host=localhost port=5432 dbname=somedatabase sslmode=require user=myuser password=mypass",
		},
		{
			in:  "postgres://myuser@localhost:1234/somedatabase",
			out: "host=localhost port=1234 dbname=somedatabase sslmode=require user=myuser",
		},
		{
			in:   "postgres://myuser:mypass@localhost/somedatabase",
			opts: []Opts{WithSSL(SSLBlank)},
			out:  "host=localhost port=5432 dbname=somedatabase sslmode= user=myuser password=mypass",
		},
		{
			in:   "postgres://myuser:mypass@localhost/somedatabase",
			opts: []Opts{WithSSL(SSLDisable)},
			out:  "host=localhost port=5432 dbname=somedatabase sslmode=disable user=myuser password=mypass",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run("", func(t *testing.T) {
			t.Parallel()

			c, err := NewFromURL(tt.in, tt.opts...)
			if tt.err {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.out, c.String())
		})
	}

}
