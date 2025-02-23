# Storage Module

A flexible Golang storage abstraction layer supporting multiple storage backends (local filesystem and Amazon S3).

## Features

- Abstract storage interface for different storage implementations
- Supported storage backends:
  - Local filesystem
  - Amazon S3
- Common operations:
  - Read/Write files
  - List directories
  - Check file existence
  - Delete files
  - Create directories

## Installation

```bash
go get github.com/hboisgibault/storage
```

## Usage

### Basic Example

```go
package main

import (
    "github.com/hboisgibault/storage"
)

func main() {
    // Create S3 storage
    store, err := storage.CreateStorage("s3", "your-bucket-name")
    if err != nil {
        panic(err)
    }

    // Write a file
    err = store.Write("test.txt", "Hello World")
    if err != nil {
        panic(err)
    }

    // Read a file
    reader, err := store.Read("test.txt")
    if err != nil {
        panic(err)
    }
    defer reader.Close()

    // List directory contents
    entries, err := store.ListDir("some/directory")
    if err != nil {
        panic(err)
    }
    for _, entry := range entries {
        fmt.Println(entry.Name())
    }
}
```

### Configuration

#### Local Storage
No configuration needed. Just specify the path when creating the storage:

```go
store, err := storage.CreateStorage("local", "/path/to/storage")
```

#### S3 Storage
Create a `.env` file with your AWS credentials:

```env
AWS_S3_ACCESS_KEY_ID=your_access_key
AWS_S3_SECRET_ACCESS_KEY=your_secret_key
AWS_S3_REGION=your_region
AWS_S3_BUCKET=your_bucket
```

Then initialize the S3 storage:

```go
store, err := storage.CreateStorage("s3", os.Getenv("AWS_S3_BUCKET"))
```

## Interface

```go
type StorageInterface interface {
    MakeDir(dir string, path string) error
    Write(key string, content string) error
    Read(key string) (io.ReadCloser, error)
    ListDir(key string) ([]fs.DirEntry, error)
    Delete(key string) error
    Exists(key string) (bool, error)
}
```

## Development

### Running Tests

```bash
go test ./...
```

### Adding New Storage Backends

1. Implement the `StorageInterface`
2. Add your implementation to the `CreateStorage` factory function
3. Add tests for your implementation

## License

This project is licensed under the MIT License - see the LICENSE file for details.
