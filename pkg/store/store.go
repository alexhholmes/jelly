// Package store provides an abstract interface for a storage backend
// implemented using Amazon S3.
package store

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// Storage defines the interface for object storage operations
type Storage interface {
	// Upload uploads data to storage and returns the public URL
	Upload(ctx context.Context, key string, data []byte, contentType string) (string, error)

	// Download retrieves data from storage
	Download(ctx context.Context, key string) ([]byte, error)

	// Delete removes an object from storage
	Delete(ctx context.Context, key string) error

	// Exists checks if an object exists in storage
	Exists(ctx context.Context, key string) (bool, error)

	// GeneratePresignedURL creates a presigned URL for temporary access
	GeneratePresignedURL(ctx context.Context, key string, expiration time.Duration) (string, error)
}

// S3Storage implements Storage interface using Amazon S3
type S3Storage struct {
	client  *s3.Client
	bucket  string
	region  string
	baseURL string // Base URL for public access (e.g., CloudFront domain)
}

// NewS3Storage creates a new S3Storage instance
func NewS3Storage(client *s3.Client, bucket, region string) *S3Storage {
	baseURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com", bucket, region)
	return &S3Storage{
		client:  client,
		bucket:  bucket,
		region:  region,
		baseURL: baseURL,
	}
}

// NewS3StorageWithCustomURL creates a new S3Storage instance with custom base URL (e.g., CloudFront)
func NewS3StorageWithCustomURL(client *s3.Client, bucket, region, baseURL string) *S3Storage {
	return &S3Storage{
		client:  client,
		bucket:  bucket,
		region:  region,
		baseURL: strings.TrimSuffix(baseURL, "/"),
	}
}

// Upload uploads data to S3 and returns the public URL
func (s *S3Storage) Upload(ctx context.Context, key string, data []byte, contentType string) (string, error) {
	if len(data) == 0 {
		return "", fmt.Errorf("data cannot be empty")
	}

	if key == "" {
		return "", fmt.Errorf("key cannot be empty")
	}

	// Set default content type if not provided
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	input := &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
		ACL:         types.ObjectCannedACLPublicRead, // Make publicly readable
	}

	_, err := s.client.PutObject(ctx, input)
	if err != nil {
		return "", fmt.Errorf("failed to upload to S3: %w", err)
	}

	// Return public URL
	url := fmt.Sprintf("%s/%s", s.baseURL, key)
	return url, nil
}

// Download retrieves data from S3
func (s *S3Storage) Download(ctx context.Context, key string) ([]byte, error) {
	if key == "" {
		return nil, fmt.Errorf("key cannot be empty")
	}

	input := &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}

	result, err := s.client.GetObject(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to download from S3: %w", err)
	}
	defer result.Body.Close()

	data, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read object body: %w", err)
	}

	return data, nil
}

// Delete removes an object from S3
func (s *S3Storage) Delete(ctx context.Context, key string) error {
	if key == "" {
		return fmt.Errorf("key cannot be empty")
	}

	input := &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}

	_, err := s.client.DeleteObject(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete from S3: %w", err)
	}

	return nil
}

// Exists checks if an object exists in S3
func (s *S3Storage) Exists(ctx context.Context, key string) (bool, error) {
	if key == "" {
		return false, fmt.Errorf("key cannot be empty")
	}

	input := &s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}

	_, err := s.client.HeadObject(ctx, input)
	if err != nil {
		// Check if error is because object doesn't exist
		var noSuchKey *types.NoSuchKey
		var notFound *types.NotFound
		if errors.As(err, &noSuchKey) || errors.As(err, &notFound) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check object existence: %w", err)
	}

	return true, nil
}

// GeneratePresignedURL creates a presigned URL for temporary access
func (s *S3Storage) GeneratePresignedURL(ctx context.Context, key string, expiration time.Duration) (string, error) {
	if key == "" {
		return "", fmt.Errorf("key cannot be empty")
	}

	if expiration <= 0 {
		return "", fmt.Errorf("expiration must be positive")
	}

	presignClient := s3.NewPresignClient(s.client)

	input := &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}

	result, err := presignClient.PresignGetObject(ctx, input, func(opts *s3.PresignOptions) {
		opts.Expires = expiration
	})

	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return result.URL, nil
}
