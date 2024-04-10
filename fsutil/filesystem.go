// Package fsutil provides filesystem-related functions.
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

	"github.com/kolide/kit/env"
)

const (
	// DirMode is the default permission used when creating directories
	DirMode = 0755
	// FileMode is the default permission used when creating files
	FileMode = 0644
)

// Gopath will return the current GOPATH as set by environment variables and
// will fall back to ~/go if a GOPATH is not set.
func Gopath() string {
	home := env.String("HOME", "~/")
	return env.String("GOPATH", filepath.Join(home, "go"))
}

// CopyDir is a utility to assist with copying a directory from src to dest.
// Note that directory permissions are not maintained, but the permissions of
// the files in those directories are.
func CopyDir(src, dest string) error {
	dir, err := os.Open(src)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dest, DirMode); err != nil {
		return err
	}

	files, err := dir.Readdir(-1)
	if err != nil {
		return err
	}
	for _, file := range files {
		srcptr := filepath.Join(src, file.Name())
		dstptr := filepath.Join(dest, file.Name())
		if file.IsDir() {
			if err := CopyDir(srcptr, dstptr); err != nil {
				return err
			}
		} else {
			if err := CopyFile(srcptr, dstptr); err != nil {
				return err
			}
		}
	}
	return nil
}

// CopyFile is a utility to assist with copying a file from src to dest.
// Note that file permissions are maintained.
func CopyFile(src, dest string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destfile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destfile.Close()

	_, err = io.Copy(destfile, source)
	if err != nil {
		return err
	}
	sourceinfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	return os.Chmod(dest, sourceinfo.Mode())
}

// UntarBundle will untar a source tar.gz archive to the supplied
// destination. Note that this calls `filepath.Dir(destination)`,
// which has the effect of stripping the last component from
// destination.
func UntarBundle(destination string, source string) error {
	f, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("opening source: %w", err)
	}
	defer f.Close()

	gzr, err := gzip.NewReader(f)
	if err != nil {
		return fmt.Errorf("creating gzip reader from %s: %w", source, err)
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("reading tar file: %w", err)
		}

		if err := sanitizeExtractPath(filepath.Dir(destination), header.Name); err != nil {
			return fmt.Errorf("checking filename: %w", err)
		}

		destPath := filepath.Join(filepath.Dir(destination), header.Name)
		info := header.FileInfo()
		if info.IsDir() {
			if err = os.MkdirAll(destPath, info.Mode()); err != nil {
				return fmt.Errorf("creating directory %s for tar file: %w", destPath, err)
			}
			continue
		}

		if err := writeBundleFile(destPath, info.Mode(), tr); err != nil {
			return fmt.Errorf("writing file: %w", err)
		}
	}
	return nil
}

// UntarBundleWithRequiredFilePermission performs the same operation as UntarBundle,
// but enforces `requiredFilePerm` for all files in the bundle.
func UntarBundleWithRequiredFilePermission(destination string, source string, requiredFilePerm fs.FileMode) error {
	f, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("opening source: %w", err)
	}
	defer f.Close()

	gzr, err := gzip.NewReader(f)
	if err != nil {
		return fmt.Errorf("creating gzip reader from %s: %w", source, err)
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("reading tar file: %w", err)
		}

		if err := sanitizeExtractPath(filepath.Dir(destination), header.Name); err != nil {
			return fmt.Errorf("checking filename: %w", err)
		}

		destPath := filepath.Join(filepath.Dir(destination), header.Name)
		info := header.FileInfo()
		if info.IsDir() {
			if err = os.MkdirAll(destPath, info.Mode()); err != nil {
				return fmt.Errorf("creating directory %s for tar file: %w", destPath, err)
			}
			continue
		}

		if err := writeBundleFile(destPath, requiredFilePerm, tr); err != nil {
			return fmt.Errorf("writing file: %w", err)
		}
	}
	return nil
}

func writeBundleFile(destPath string, perm fs.FileMode, srcReader io.Reader) error {
	file, err := os.OpenFile(destPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, perm)
	if err != nil {
		return fmt.Errorf("opening %s: %w", destPath, err)
	}
	defer file.Close()
	if _, err := io.Copy(file, srcReader); err != nil {
		return fmt.Errorf("copying to %s: %w", destPath, err)
	}

	return nil
}

// sanitizeExtractPath checks that the supplied extraction path is nor
// vulnerable to zip slip attacks. See https://snyk.io/research/zip-slip-vulnerability
func sanitizeExtractPath(filePath string, destination string) error {
	destpath := filepath.Join(destination, filePath)
	if !strings.HasPrefix(destpath, filepath.Clean(destination)+string(os.PathSeparator)) {
		return fmt.Errorf("%s: illegal file path", filePath)
	}
	return nil
}
