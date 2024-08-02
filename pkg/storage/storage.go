package storage

type StorageDriver interface {
	Write(path string, data []byte)
	Read(path string)
}
