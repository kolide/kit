package fsutil

import (
	"fmt"
	"io/fs"
	"os"
	"path"
)

// CopyToDisk copies an embedded FS to a given directory. Because go's embed does not preserve the file mode, you must
// also pass a function that will return the desired file mode for each file.
func CopyFSToDisk(src fs.FS, destDir string, modeSetter func(fs.FileInfo) os.FileMode) error {
	if err := fs.WalkDir(src, ".", genCopyToDiskFunc(src, destDir, modeSetter)); err != nil {
		return fmt.Errorf("walking directory: %w", err)
	}

	return nil
}

// CommonFileMode is a function that returns the common file permissions: 0644 for all files, and 0755 for
// all directories. It is provided as a helper for CopyFSToDisk's common use case
func CommonFileMode(fi fs.FileInfo) fs.FileMode {
	if fi.IsDir() {
		return 0755
	} else {
		return 0644
	}
}

// genCopyToDiskFunc returns fs.WalkDirFunc function that will
// copy files to disk in a given location.
func genCopyToDiskFunc(srcFS fs.FS, destDir string, modeSetter func(fs.FileInfo) os.FileMode) fs.WalkDirFunc {
	return func(filepath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		fileinfo, err := d.Info()
		if err != nil {
			return fmt.Errorf("getting file info: %w", err)
		}

		fullpath := path.Join(destDir, filepath)

		// If it's a directory, make it under destdir
		if d.IsDir() {
			if err := os.MkdirAll(fullpath, modeSetter(fileinfo)); err != nil {
				return fmt.Errorf("making directory %s: %w", fullpath, err)
			}
			return nil
		}

		data, err := fs.ReadFile(srcFS, filepath)
		if err != nil {
			return fmt.Errorf("reading file from FS %s: %w", filepath, err)
		}

		if err := os.WriteFile(fullpath, data, modeSetter(fileinfo)); err != nil {
			return fmt.Errorf("writing %s: %w", filepath, err)
		}

		return nil
	}
}
