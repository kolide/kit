package fsutil

import (
	"embed"
	"io/fs"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

//go:embed test-data/fscopy
var embeddedSourceData embed.FS

func TestCopyFSToDisk(t *testing.T) {
	t.Parallel()

	subdir, err := fs.Sub(embeddedSourceData, "test-data/fscopy")
	require.NoError(t, err)

	destDir := t.TempDir()

	require.NoError(t, CopyFSToDisk(subdir, destDir, CommonFileMode))

	var tests = []struct {
		path string
	}{
		{path: "top.txt"},
		{path: path.Join("level1", "level2", "deep.txt")},
		{path: path.Join("level1", "l1.txt")},
		{path: path.Join("levelA", "levelB", "b.txt")},
		{path: path.Join("levelA", "a.txt")},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.path, func(t *testing.T) {
			t.Parallel()

			expected := "hello\n"
			if runtime.GOOS == "windows" {
				expected = "hello\r\n"
			}

			contents, err := os.ReadFile(path.Join(destDir, tt.path))
			require.NoError(t, err)
			require.Equal(t, expected, string(contents))
		})
	}
}
