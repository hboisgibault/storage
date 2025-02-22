package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func setupTest() func() {
	dir := "testdir"
	filename := "test.txt"
	fullPath := filepath.Join(dir, filename)
	subdir := filepath.Join(dir, "subdir")
	content := []byte("random text")

	os.Mkdir(dir, os.ModePerm)
	os.Mkdir(subdir, os.ModePerm)
	os.Create(filepath.Join(subdir, "1.txt"))
	os.Create(filepath.Join(subdir, "2.txt"))
	os.Create(filepath.Join(subdir, "3.txt"))
	os.Create(filepath.Join(subdir, "4.txt"))
	os.WriteFile(fullPath, content, 0644)

	return func() {
		os.RemoveAll(dir)
	}
}

func TestFactory(t *testing.T) {
	defer setupTest()()
	storage, error := CreateStorage("local", "testdir")
	storageClass := reflect.TypeOf(storage).String()

	if storageClass != "*storage.LocalStorage" {
		t.Fatalf(`Storage type = %q, %v, want match for %#q`, storageClass, error, "*storage.LocalStorage")
	}
}

func TestExists(t *testing.T) {
	defer setupTest()()
	storage, error := CreateStorage("local", "testdir")

	filename := "test.txt"
	exists, error := storage.Exists(filename)

	filenameisnot := "random.txt"
	isnot, error := storage.Exists(filenameisnot)

	if exists == false {
		t.Fatalf(`File does not exist = %q, %v`, filename, error)
	}

	if isnot == true {
		t.Fatalf(`File exists = %q, %v`, filenameisnot, error)
	}
}

func TestRead(t *testing.T) {
	defer setupTest()()
	storage, error := CreateStorage("local", "testdir")

	filename := "test.txt"
	content := "random text"
	reader, error := storage.Read(filename)

	buf := new(strings.Builder)
	io.Copy(buf, reader)

	if buf.String() != content {
		t.Fatalf(`Storage Read = %q, %v, want match for %#q`, buf.String(), error, content)
	}
}

func TestListDir(t *testing.T) {
	defer setupTest()()
	storage, error := CreateStorage("local", "testdir")

	dir := "subdir"
	files, error := storage.ListDir(dir)
	if error != nil {
		t.Fatalf(`File read fail, %v`, error)
	}

	for i, item := range files {
		filename := fmt.Sprintf("%d.txt", i+1)

		if item.Name() != filename {
			t.Fatalf(`File read = %q, %v, want match for %#q`, item.Name(), error, filename)
		}
	}
}

func TestWrite(t *testing.T) {
	defer setupTest()()
	storage, error := CreateStorage("local", "testdir")

	filename := "write.txt"
	content := "random text"
	storage.Write(filename, content)

	fullPath := filepath.Join("testdir", filename)
	file, _ := os.Open(fullPath)

	buf := new(strings.Builder)
	io.Copy(buf, file)

	if buf.String() != content {
		t.Fatalf(`Storage Write = %q, %v, want match for %#q`, buf.String(), error, content)
	}
}

func TestDelete(t *testing.T) {
	defer setupTest()()
	storage, _ := CreateStorage("local", "testdir")

	filename := "test.txt"
	delete := storage.Delete(filename)
	exists, _ := storage.Exists(filename)

	if exists == true {
		t.Fatalf(`File delete = %q, %v, fail`, filename, delete)
	}
}
