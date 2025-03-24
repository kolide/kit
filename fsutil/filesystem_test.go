package fsutil

import (
	"archive/tar"
	"compress/gzip"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUntarBundle(t *testing.T) {
	t.Parallel()

	// Create tarball contents
	originalDir := t.TempDir()
	topLevelFile := filepath.Join(originalDir, "testfile.txt")
	var topLevelFileMode fs.FileMode = 0655
	require.NoError(t, os.WriteFile(topLevelFile, []byte("test1"), topLevelFileMode))
	internalDir := filepath.Join(originalDir, "some", "path", "to")
	var nestedFileMode fs.FileMode = 0755
	require.NoError(t, os.MkdirAll(internalDir, nestedFileMode))
	nestedFile := filepath.Join(internalDir, "anotherfile.txt")
	require.NoError(t, os.WriteFile(nestedFile, []byte("test2"), nestedFileMode))

	// Create test tarball
	tarballDir := t.TempDir()
	tarballFile := filepath.Join(tarballDir, "test.gz")
	createTar(t, tarballFile, originalDir)

	// Confirm we can untar the tarball successfully
	newDir := t.TempDir()
	require.NoError(t, UntarBundle(filepath.Join(newDir, "anything"), tarballFile))

	// Confirm the tarball has the contents we expect
	newTopLevelFile := filepath.Join(newDir, filepath.Base(topLevelFile))
	require.FileExists(t, newTopLevelFile)
	newNestedFile := filepath.Join(newDir, "some", "path", "to", filepath.Base(nestedFile))
	require.FileExists(t, newNestedFile)

	// Confirm each file retained its original permissions
	topLevelFileInfo, err := os.Stat(newTopLevelFile)
	require.NoError(t, err)
	require.Equal(t, topLevelFileMode.String(), topLevelFileInfo.Mode().String())
	nestedFileInfo, err := os.Stat(newNestedFile)
	require.NoError(t, err)
	require.Equal(t, nestedFileMode.String(), nestedFileInfo.Mode().String())
}

// createTar is a helper to create a test tar
func createTar(t *testing.T, createLocation string, sourceDir string) {
	tarballFile, err := os.Create(createLocation)
	require.NoError(t, err)
	defer tarballFile.Close()

	gzw := gzip.NewWriter(tarballFile)
	defer gzw.Close()

	tw := tar.NewWriter(gzw)
	defer tw.Close()

	require.NoError(t, tw.AddFS(os.DirFS(sourceDir)))
}

func TestSanitizeExtractPath(t *testing.T) {
	t.Parallel()

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
