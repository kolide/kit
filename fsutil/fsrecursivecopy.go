package fsutil

import (
	"fmt"
	"io/fs"
	"os"
	"path"
)

func CopyFSToDisk(src fs.FS, destDir string) error {
	if err := fs.WalkDir(src, ".", genCopyToDiskFunc(src, destDir)); err != nil {
		return fmt.Errorf("walking directory: %w", err)
	}

	return nil

}

// genCopyToDiskFunc returns fs.WalkDirFunc function that will
// copy files to disk in a given location.
func genCopyToDiskFunc(srcFS fs.FS, destDir string) fs.WalkDirFunc {

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
			if err := os.MkdirAll(fullpath, 0755); err != nil {
				return fmt.Errorf("making directory %s: %w", fullpath, err)
			}
			return nil
		}

		data, err := fs.ReadFile(srcFS, filepath)
		if err != nil {
			return fmt.Errorf("reading file from FS %s: %w", filepath, err)
		}

		if err := os.WriteFile(fullpath, data, fileinfo.Mode()); err != nil {
			return fmt.Errorf("writing %s: %w", filepath, err)
		}

		return nil
	}
}
