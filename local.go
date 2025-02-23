package storage

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

type LocalStorage struct {
	Storage
}

func newLocalStorage(basePath string) StorageInterface {
	return &LocalStorage{
		Storage: Storage{
			Type: "local",
			Path: basePath,
		},
	}
}

func (l *LocalStorage) MakeDir(dir string, path string) error {
	fullPath := filepath.Join(l.Storage.Path, path, dir)
	return os.MkdirAll(fullPath, 0755)
}

func (l *LocalStorage) Write(key string, content string) error {
	fullPath := filepath.Join(l.Storage.Path, key)
	file, err := os.OpenFile(fullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	if content != "" {
		_, err := file.WriteString(content)
		file.Close()
		if err != nil {
			return err
		}
		return nil
	}

	return nil
}

func (l *LocalStorage) Read(key string) (io.ReadCloser, error) {
	fullPath := filepath.Join(l.Storage.Path, key)
	return os.Open(fullPath)
}

func (l *LocalStorage) ListDir(key string) ([]fs.DirEntry, error) {
	fullPath := filepath.Join(l.Storage.Path, key)
	return os.ReadDir(fullPath)
}

func (l *LocalStorage) Delete(key string) error {
	fullPath := filepath.Join(l.Storage.Path, key)
	return os.Remove(fullPath)
}

func (l *LocalStorage) Exists(key string) (bool, error) {
	fullPath := filepath.Join(l.Storage.Path, key)
	_, err := os.Stat(fullPath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
