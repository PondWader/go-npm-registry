package storage

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type fsStorageDriver struct {
	BaseDir string
}

func NewFsDriver(baseDir string) (StorageDriver, error) {
	baseDir = filepath.Clean(baseDir)
	if err := os.MkdirAll(baseDir, 0o700); err != nil {
		return nil, err
	}
	return fsStorageDriver{BaseDir: baseDir}, nil
}

// Resolves a path to the actual file system path, checking that the path is contained within the base directory
func (fs fsStorageDriver) ResolvePath(path string) (string, error) {
	resolved := filepath.Join(fs.BaseDir, path)
	// Check that the path is in the base dir for security
	if strings.HasPrefix(resolved, "..") || !strings.HasPrefix(resolved, fs.BaseDir+"/") {
		return "", errors.New("path is invalid")
	}
	return resolved, nil
}

func (fs fsStorageDriver) Write(path string, data []byte) error {
	resolved, err := fs.ResolvePath(path)
	if err != nil {
		return err
	}
	if err = os.MkdirAll(filepath.Dir(resolved), 0o700); err != nil {
		return err
	}
	return os.WriteFile(resolved, data, 0o700)
}

func (fs fsStorageDriver) Read(path string) ([]byte, error) {
	resolved, err := fs.ResolvePath(path)
	if err != nil {
		return nil, err
	}
	return os.ReadFile(resolved)
}

func (fs fsStorageDriver) NewReader(path string) (io.ReadCloser, error) {
	resolved, err := fs.ResolvePath(path)
	if err != nil {
		return nil, err
	}
	file, err := os.Open(resolved)
	if err != nil {
		return nil, err
	}
	return file, nil
}
