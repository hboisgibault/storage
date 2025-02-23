package storage

import (
	"context"
	"io"
	"io/fs"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
)

type S3Storage struct {
	Storage
	client S3Client
	bucket string
}

type S3Client interface {
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
	ListObjectsV2(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error)
	DeleteObject(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error)
	HeadObject(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error)
}

func newS3Storage(bucket string) (StorageInterface, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(os.Getenv("AWS_S3_REGION")))
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg)
	return &S3Storage{
		Storage: Storage{
			Type: "s3",
			Path: bucket,
		},
		client: client,
		bucket: bucket,
	}, nil
}

func (s *S3Storage) MakeDir(dir string, path string) error {
	// S3 doesn't need explicit directory creation
	return nil
}

func (s *S3Storage) Write(key string, content string) error {
	if content != "" {
		_, err := s.client.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket: aws.String(s.bucket),
			Key:    aws.String(key),
			Body:   strings.NewReader(content),
		})
		return err
	}

	return nil
}

func (s *S3Storage) Read(key string) (io.ReadCloser, error) {
	output, err := s.client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	return output.Body, nil
}

func (s *S3Storage) ListDir(key string) ([]fs.DirEntry, error) {
	// Ensure the prefix ends with "/"
	prefix := key
	if !strings.HasSuffix(prefix, "/") {
		prefix = prefix + "/"
	}

	output, err := s.client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket:    aws.String(s.bucket),
		Prefix:    aws.String(prefix),
		Delimiter: aws.String("/"),
	})
	if err != nil {
		return nil, err
	}

	entries := make([]fs.DirEntry, 0)
	for _, obj := range output.Contents {
		if *obj.Size > 0 {
			entries = append(entries, newS3DirEntry(*obj.Key, *obj.Size))
		}
	}
	return entries, nil
}

func (s *S3Storage) Delete(key string) error {
	_, err := s.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	return err
}

func (s *S3Storage) Exists(key string) (bool, error) {
	_, err := s.client.HeadObject(context.TODO(), &s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return false, nil
	}
	return true, nil
}

type s3DirEntry struct {
	name string
	size int64
}

func newS3DirEntry(name string, size int64) *s3DirEntry {
	return &s3DirEntry{name: name, size: size}
}

func (e *s3DirEntry) Name() string               { return e.name }
func (e *s3DirEntry) IsDir() bool                { return strings.HasSuffix(e.name, "/") }
func (e *s3DirEntry) Type() fs.FileMode          { return fs.FileMode(0644) }
func (e *s3DirEntry) Info() (fs.FileInfo, error) { return nil, fs.ErrNotExist }
