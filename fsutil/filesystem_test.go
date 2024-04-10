package fsutil

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUntarBundle(t *testing.T) {
	t.Parallel()

	// Create tarball contents
	originalDir := t.TempDir()
	topLevelFile := filepath.Join(originalDir, "testfile.txt")
	require.NoError(t, os.WriteFile(topLevelFile, []byte("test1"), 0655))
	internalDir := filepath.Join(originalDir, "some", "path", "to")
	require.NoError(t, os.MkdirAll(internalDir, 0755))
	nestedFile := filepath.Join(internalDir, "anotherfile.txt")
	require.NoError(t, os.WriteFile(nestedFile, []byte("test2"), 0755))

	// Create test tarball
	tarballDir := t.TempDir()
	tarballFile := filepath.Join(tarballDir, "test.gz")
	createTar(t, tarballFile, originalDir)

	// Confirm we can untar the tarball successfully
	newDir := t.TempDir()
	require.NoError(t, UntarBundle(filepath.Join(newDir, "anything"), tarballFile))

	// Confirm the tarball has the contents we expect
	require.FileExists(t, filepath.Join(newDir, filepath.Base(topLevelFile)))
	require.FileExists(t, filepath.Join(newDir, "some", "path", "to", filepath.Base(nestedFile)))
}

func TestUntarBundleWithRequiredFilePermission(t *testing.T) {
	t.Parallel()

	// Create tarball contents
	originalDir := t.TempDir()
	topLevelFile := filepath.Join(originalDir, "testfile.txt")
	require.NoError(t, os.WriteFile(topLevelFile, []byte("test1"), 0655))
	internalDir := filepath.Join(originalDir, "some", "path", "to")
	require.NoError(t, os.MkdirAll(internalDir, 0755))
	nestedFile := filepath.Join(internalDir, "anotherfile.txt")
	require.NoError(t, os.WriteFile(nestedFile, []byte("test2"), 0744))

	// Create test tarball
	tarballDir := t.TempDir()
	tarballFile := filepath.Join(tarballDir, "test.gz")
	createTar(t, tarballFile, originalDir)

	// Confirm we can untar the tarball successfully
	newDir := t.TempDir()
	var requiredFileMode fs.FileMode = 0755
	require.NoError(t, UntarBundleWithRequiredFilePermission(filepath.Join(newDir, "anything"), tarballFile, requiredFileMode))

	// Confirm the tarball has the contents we expect
	newTopLevelFile := filepath.Join(newDir, filepath.Base(topLevelFile))
	require.FileExists(t, newTopLevelFile)
	newNestedFile := filepath.Join(newDir, "some", "path", "to", filepath.Base(nestedFile))
	require.FileExists(t, newNestedFile)

	// Require that both files have the required permission 0755
	topLevelFileInfo, err := os.Stat(newTopLevelFile)
	require.NoError(t, err)
	require.Equal(t, requiredFileMode, topLevelFileInfo.Mode())
	nestedFileInfo, err := os.Stat(newNestedFile)
	require.NoError(t, err)
	require.Equal(t, requiredFileMode, nestedFileInfo.Mode())
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

	require.NoError(t, filepath.Walk(sourceDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		srcInfo, err := os.Lstat(path)
		if os.IsNotExist(err) {
			return fmt.Errorf("error adding %s to tarball: %w", path, err)
		}

		hdr, err := tar.FileInfoHeader(srcInfo, path)
		if err != nil {
			return fmt.Errorf("error creating tar header: %w", err)
		}
		hdr.Name = strings.TrimPrefix(path, sourceDir+"/")

		if err := tw.WriteHeader(hdr); err != nil {
			return fmt.Errorf("error writing tar header: %w", err)
		}

		if !srcInfo.Mode().IsRegular() {
			// Don't open/copy over directories
			return nil
		}

		srcFile, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("error opening file to add to tarball: %w", err)
		}
		defer srcFile.Close()

		if _, err := io.Copy(tw, srcFile); err != nil {
			return fmt.Errorf("error copying file %s to tarball: %w", path, err)
		}

		return nil
	}))
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
