package fs

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSanitizeExtractPath(t *testing.T) {
	var tests = []struct {
		filepath    string
		destination string
		expectError bool
	}{
		{
			filepath:    "file",
			destination: "/tmp",
			expectError: false,
		},
		{
			filepath:    "subdir/../subdir/file",
			destination: "/tmp",
			expectError: false,
		},

		{
			filepath:    "../../../file",
			destination: "/tmp",
			expectError: true,
		},
		{
			filepath:    "./././file",
			destination: "/tmp",
			expectError: false,
		},
	}

	for _, tt := range tests {
		if tt.expectError {
			require.Error(t, sanitizeExtractPath(tt.filepath, tt.destination), tt.filepath)
		} else {
			require.NoError(t, sanitizeExtractPath(tt.filepath, tt.destination), tt.filepath)
		}

	}

}
