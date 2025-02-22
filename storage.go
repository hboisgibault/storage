// Package storage implement different storage engines
package storage

import (
	"fmt"
	"io"
	"io/fs"
)

// Abstract interface to represent a storage system
type StorageInterface interface {
	// Make directory
	MakeDir(dir string, path string) error

	// Write file
	Write(key string, content string) (io.WriteCloser, error)

	// Read file
	Read(key string) (io.ReadCloser, error)

	// Read directory
	ListDir(key string) ([]fs.DirEntry, error)

	// Delete file
	Delete(key string) error

	// Check if file exists
	Exists(key string) (bool, error)
}

type Storage struct {
	Type string
	Path string
}

// Storage factory
func CreateStorage(storageType string, path string) (StorageInterface, error) {
	switch storageType {
	case "local":
		return newLocalStorage(path), nil
	default:
		return nil, fmt.Errorf("Wrong storage type passed")
	}
}
