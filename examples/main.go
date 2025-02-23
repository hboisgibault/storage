package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/hboisgibault/storage"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	store, err := storage.CreateStorage("s3", os.Getenv("AWS_S3_BUCKET"))
	if err != nil {
		log.Fatal(err)
	}

	dir := "output"

	dirEntries, err := store.ListDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range dirEntries {
		fmt.Println("File:", file.Name())
	}

	exists, err := store.Exists(dir + "/test.txt")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("File exists:", exists)

	store.Write(dir+"/test.txt", "Hello World")

	reader, err := store.Read(dir + "/test.txt")
	buf := new(strings.Builder)
	io.Copy(buf, reader)
	fmt.Println("File content:", buf.String())
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	store.Delete(dir + "/test.txt")
}
