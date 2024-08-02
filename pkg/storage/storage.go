package storage

import (
	"errors"
	"io"
)

type StorageDriver interface {
	Write(path string, data []byte) error
	Read(path string) ([]byte, error)
	NewReader(path string) (io.ReadCloser, error)
}

func New(driver string, opts map[string]string) (StorageDriver, error) {
	if driver == "fs" {
		baseDir, ok := opts["base-dir"]
		if !ok {
			return fsStorageDriver{}, errors.New(`missing "base-path" option in fs storage driver options`)
		}
		return NewFsDriver(baseDir)
	}

	return fsStorageDriver{}, errors.New(`storage driver "` + driver + `" doesn't exist`)
}
