package storage

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockS3Client struct {
	mock.Mock
}

func (m *mockS3Client) PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	args := m.Called(ctx, params)
	return &s3.PutObjectOutput{}, args.Error(1)
}

func (m *mockS3Client) GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	args := m.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return &s3.GetObjectOutput{
		Body: io.NopCloser(strings.NewReader(args.String(0))),
	}, nil
}

func (m *mockS3Client) ListObjectsV2(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*s3.ListObjectsV2Output), args.Error(1)
}

func (m *mockS3Client) DeleteObject(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
	args := m.Called(ctx, params)
	return &s3.DeleteObjectOutput{}, args.Error(1)
}

func (m *mockS3Client) HeadObject(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
	args := m.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return &s3.HeadObjectOutput{}, nil
}

func TestS3Storage_Write(t *testing.T) {
	mockClient := new(mockS3Client)
	storage := &S3Storage{
		Storage: Storage{Type: "s3", Path: "testbucket"},
		client:  mockClient,
		bucket:  "testbucket",
	}

	mockClient.On("PutObject", mock.Anything, mock.Anything).Return(nil, nil)

	// Test direct content write
	err := storage.Write("test.txt", "hello world")
	assert.NoError(t, err)

	// Test streaming write
	err = storage.Write("test2.txt", "")
	assert.NoError(t, err)

	mockClient.AssertExpectations(t)
}

func TestS3Storage_Read(t *testing.T) {
	mockClient := new(mockS3Client)
	storage := &S3Storage{
		Storage: Storage{Type: "s3", Path: "testbucket"},
		client:  mockClient,
		bucket:  "testbucket",
	}

	mockClient.On("GetObject", mock.Anything, mock.Anything).Return("test content", nil)

	reader, err := storage.Read("test.txt")
	assert.NoError(t, err)

	content, err := io.ReadAll(reader)
	assert.NoError(t, err)
	assert.Equal(t, "test content", string(content))

	mockClient.AssertExpectations(t)
}

func TestS3Storage_ListDir(t *testing.T) {
	mockClient := new(mockS3Client)
	storage := &S3Storage{
		Storage: Storage{Type: "s3", Path: "testbucket"},
		client:  mockClient,
		bucket:  "testbucket",
	}

	mockResponse := &s3.ListObjectsV2Output{
		Contents: []types.Object{
			{Key: aws.String("folder/file1.txt"), Size: aws.Int64(100)},
			{Key: aws.String("folder/file2.txt"), Size: aws.Int64(200)},
		},
	}

	mockClient.On("ListObjectsV2", mock.Anything, mock.Anything).Return(mockResponse, nil)

	entries, err := storage.ListDir("folder")
	assert.NoError(t, err)
	assert.Len(t, entries, 2)
	assert.Equal(t, "folder/file1.txt", entries[0].Name())

	mockClient.AssertExpectations(t)
}

func TestS3Storage_Delete(t *testing.T) {
	mockClient := new(mockS3Client)
	storage := &S3Storage{
		Storage: Storage{Type: "s3", Path: "testbucket"},
		client:  mockClient,
		bucket:  "testbucket",
	}

	mockClient.On("DeleteObject", mock.Anything, mock.Anything).Return(nil, nil)

	err := storage.Delete("test.txt")
	assert.NoError(t, err)

	mockClient.AssertExpectations(t)
}

func TestS3Storage_Exists(t *testing.T) {
	mockClient := new(mockS3Client)
	storage := &S3Storage{
		Storage: Storage{Type: "s3", Path: "testbucket"},
		client:  mockClient,
		bucket:  "testbucket",
	}

	mockClient.On("HeadObject", mock.Anything, mock.Anything).Return(nil, nil)

	exists, err := storage.Exists("test.txt")
	assert.NoError(t, err)
	assert.True(t, exists)

	mockClient.AssertExpectations(t)
}
