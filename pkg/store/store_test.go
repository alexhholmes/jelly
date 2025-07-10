package store

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewS3Storage(t *testing.T) {
	tests := []struct {
		name     string
		bucket   string
		region   string
		expected string
	}{
		{
			name:     "US East region",
			bucket:   "test-bucket",
			region:   "us-east-1",
			expected: "https://test-bucket.s3.us-east-1.amazonaws.com",
		},
		{
			name:     "EU West region",
			bucket:   "my-photos",
			region:   "eu-west-1",
			expected: "https://my-photos.s3.eu-west-1.amazonaws.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := NewS3Storage(nil, tt.bucket, tt.region)
			
			assert.Equal(t, tt.bucket, storage.bucket)
			assert.Equal(t, tt.region, storage.region)
			assert.Equal(t, tt.expected, storage.baseURL)
		})
	}
}

func TestNewS3StorageWithCustomURL(t *testing.T) {
	tests := []struct {
		name     string
		bucket   string
		region   string
		baseURL  string
		expected string
	}{
		{
			name:     "CloudFront URL",
			bucket:   "test-bucket",
			region:   "us-east-1",
			baseURL:  "https://d123456.cloudfront.net",
			expected: "https://d123456.cloudfront.net",
		},
		{
			name:     "Custom domain with trailing slash",
			bucket:   "photos",
			region:   "us-west-2",
			baseURL:  "https://photos.example.com/",
			expected: "https://photos.example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := NewS3StorageWithCustomURL(nil, tt.bucket, tt.region, tt.baseURL)
			
			assert.Equal(t, tt.bucket, storage.bucket)
			assert.Equal(t, tt.region, storage.region)
			assert.Equal(t, tt.expected, storage.baseURL)
		})
	}
}

func TestStorage_Upload(t *testing.T) {
	tests := []struct {
		name        string
		key         string
		data        []byte
		contentType string
		setupMock   func(*MockStorage)
		expectedURL string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "successful upload",
			key:         "photos/test.jpg",
			data:        []byte("test image data"),
			contentType: "image/jpeg",
			setupMock: func(m *MockStorage) {
				m.On("Upload", context.Background(), "photos/test.jpg", []byte("test image data"), "image/jpeg").
					Return("https://example.com/photos/test.jpg", nil)
			},
			expectedURL: "https://example.com/photos/test.jpg",
			expectError: false,
		},
		{
			name:        "upload with empty content type",
			key:         "files/document.pdf",
			data:        []byte("pdf data"),
			contentType: "",
			setupMock: func(m *MockStorage) {
				m.On("Upload", context.Background(), "files/document.pdf", []byte("pdf data"), "").
					Return("https://example.com/files/document.pdf", nil)
			},
			expectedURL: "https://example.com/files/document.pdf",
			expectError: false,
		},
		{
			name:        "upload failure",
			key:         "test.jpg",
			data:        []byte("test data"),
			contentType: "image/jpeg",
			setupMock: func(m *MockStorage) {
				m.On("Upload", context.Background(), "test.jpg", []byte("test data"), "image/jpeg").
					Return("", errors.New("upload failed"))
			},
			expectError: true,
			errorMsg:    "upload failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := NewMockStorage(t)
			tt.setupMock(mockStorage)
			
			url, err := mockStorage.Upload(context.Background(), tt.key, tt.data, tt.contentType)
			
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Empty(t, url)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedURL, url)
			}
		})
	}
}

func TestStorage_Download(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		setupMock    func(*MockStorage)
		expectedData []byte
		expectError  bool
		errorMsg     string
	}{
		{
			name: "successful download",
			key:  "photos/test.jpg",
			setupMock: func(m *MockStorage) {
				m.On("Download", context.Background(), "photos/test.jpg").
					Return([]byte("test image data"), nil)
			},
			expectedData: []byte("test image data"),
			expectError:  false,
		},
		{
			name: "download failure",
			key:  "missing.jpg",
			setupMock: func(m *MockStorage) {
				m.On("Download", context.Background(), "missing.jpg").
					Return([]byte(nil), errors.New("file not found"))
			},
			expectError: true,
			errorMsg:    "file not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := NewMockStorage(t)
			tt.setupMock(mockStorage)
			
			data, err := mockStorage.Download(context.Background(), tt.key)
			
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Nil(t, data)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedData, data)
			}
		})
	}
}

func TestStorage_Delete(t *testing.T) {
	tests := []struct {
		name        string
		key         string
		setupMock   func(*MockStorage)
		expectError bool
		errorMsg    string
	}{
		{
			name: "successful delete",
			key:  "photos/test.jpg",
			setupMock: func(m *MockStorage) {
				m.On("Delete", context.Background(), "photos/test.jpg").
					Return(nil)
			},
			expectError: false,
		},
		{
			name: "delete failure",
			key:  "test.jpg",
			setupMock: func(m *MockStorage) {
				m.On("Delete", context.Background(), "test.jpg").
					Return(errors.New("delete failed"))
			},
			expectError: true,
			errorMsg:    "delete failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := NewMockStorage(t)
			tt.setupMock(mockStorage)
			
			err := mockStorage.Delete(context.Background(), tt.key)
			
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestStorage_Exists(t *testing.T) {
	tests := []struct {
		name        string
		key         string
		setupMock   func(*MockStorage)
		expected    bool
		expectError bool
		errorMsg    string
	}{
		{
			name: "object exists",
			key:  "photos/test.jpg",
			setupMock: func(m *MockStorage) {
				m.On("Exists", context.Background(), "photos/test.jpg").
					Return(true, nil)
			},
			expected:    true,
			expectError: false,
		},
		{
			name: "object does not exist",
			key:  "photos/missing.jpg",
			setupMock: func(m *MockStorage) {
				m.On("Exists", context.Background(), "photos/missing.jpg").
					Return(false, nil)
			},
			expected:    false,
			expectError: false,
		},
		{
			name: "exists check failure",
			key:  "test.jpg",
			setupMock: func(m *MockStorage) {
				m.On("Exists", context.Background(), "test.jpg").
					Return(false, errors.New("permission denied"))
			},
			expected:    false,
			expectError: true,
			errorMsg:    "permission denied",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := NewMockStorage(t)
			tt.setupMock(mockStorage)
			
			exists, err := mockStorage.Exists(context.Background(), tt.key)
			
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, exists)
			}
		})
	}
}

func TestStorage_GeneratePresignedURL(t *testing.T) {
	tests := []struct {
		name        string
		key         string
		expiration  time.Duration
		setupMock   func(*MockStorage)
		expectedURL string
		expectError bool
		errorMsg    string
	}{
		{
			name:       "successful presigned URL generation",
			key:        "photos/test.jpg",
			expiration: 15 * time.Minute,
			setupMock: func(m *MockStorage) {
				m.On("GeneratePresignedURL", context.Background(), "photos/test.jpg", 15*time.Minute).
					Return("https://presigned.example.com/photos/test.jpg?expires=123456", nil)
			},
			expectedURL: "https://presigned.example.com/photos/test.jpg?expires=123456",
			expectError: false,
		},
		{
			name:       "presigned URL generation failure",
			key:        "test.jpg",
			expiration: 5 * time.Minute,
			setupMock: func(m *MockStorage) {
				m.On("GeneratePresignedURL", context.Background(), "test.jpg", 5*time.Minute).
					Return("", errors.New("presign failed"))
			},
			expectError: true,
			errorMsg:    "presign failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := NewMockStorage(t)
			tt.setupMock(mockStorage)
			
			url, err := mockStorage.GeneratePresignedURL(context.Background(), tt.key, tt.expiration)
			
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Empty(t, url)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedURL, url)
			}
		})
	}
}

// Integration-style test that validates the Storage interface implementation
func TestS3Storage_ImplementsStorageInterface(t *testing.T) {
	var _ Storage = (*S3Storage)(nil)
}

// Test edge cases and boundary conditions
func TestStorage_EdgeCases(t *testing.T) {
	t.Run("upload with large data", func(t *testing.T) {
		mockStorage := NewMockStorage(t)
		largeData := make([]byte, 10*1024*1024) // 10MB
		
		mockStorage.On("Upload", context.Background(), "large-file.bin", largeData, "application/octet-stream").
			Return("https://example.com/large-file.bin", nil)
		
		url, err := mockStorage.Upload(context.Background(), "large-file.bin", largeData, "application/octet-stream")
		
		assert.NoError(t, err)
		assert.Equal(t, "https://example.com/large-file.bin", url)
	})
	
	t.Run("download empty file", func(t *testing.T) {
		mockStorage := NewMockStorage(t)
		
		mockStorage.On("Download", context.Background(), "empty.txt").
			Return([]byte{}, nil)
		
		data, err := mockStorage.Download(context.Background(), "empty.txt")
		
		assert.NoError(t, err)
		assert.Empty(t, data)
	})
	
	t.Run("presigned URL with very short expiration", func(t *testing.T) {
		mockStorage := NewMockStorage(t)
		shortExpiration := 1 * time.Second
		
		mockStorage.On("GeneratePresignedURL", context.Background(), "temp.jpg", shortExpiration).
			Return("https://short-lived.example.com/temp.jpg", nil)
		
		url, err := mockStorage.GeneratePresignedURL(context.Background(), "temp.jpg", shortExpiration)
		
		assert.NoError(t, err)
		assert.Contains(t, url, "temp.jpg")
	})
}

// Benchmark tests for Storage interface
func BenchmarkStorage_Upload(b *testing.B) {
	mockStorage := NewMockStorage(b)
	data := make([]byte, 1024) // 1KB test data
	
	// Setup mock to handle many calls
	mockStorage.On("Upload", context.Background(), "bench-test.jpg", data, "image/jpeg").
		Return("https://example.com/bench-test.jpg", nil).
		Times(b.N)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = mockStorage.Upload(context.Background(), "bench-test.jpg", data, "image/jpeg")
	}
}

func BenchmarkStorage_Download(b *testing.B) {
	mockStorage := NewMockStorage(b)
	testData := []byte("benchmark test data")
	
	// Setup mock to handle many calls
	mockStorage.On("Download", context.Background(), "bench-test.jpg").
		Return(testData, nil).
		Times(b.N)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = mockStorage.Download(context.Background(), "bench-test.jpg")
	}
}

// Test multiple operations in sequence
func TestStorage_MultipleOperations(t *testing.T) {
	mockStorage := NewMockStorage(t)
	key := "multi-op-test.jpg"
	data := []byte("test data for multiple operations")
	
	// Setup mock expectations in order
	mockStorage.On("Upload", context.Background(), key, data, "image/jpeg").
		Return("https://example.com/"+key, nil).
		Once()
	
	mockStorage.On("Exists", context.Background(), key).
		Return(true, nil).
		Once()
	
	mockStorage.On("Download", context.Background(), key).
		Return(data, nil).
		Once()
	
	mockStorage.On("GeneratePresignedURL", context.Background(), key, 30*time.Minute).
		Return("https://presigned.example.com/"+key, nil).
		Once()
	
	mockStorage.On("Delete", context.Background(), key).
		Return(nil).
		Once()
	
	// Execute operations in sequence
	url, err := mockStorage.Upload(context.Background(), key, data, "image/jpeg")
	assert.NoError(t, err)
	assert.Equal(t, "https://example.com/"+key, url)
	
	exists, err := mockStorage.Exists(context.Background(), key)
	assert.NoError(t, err)
	assert.True(t, exists)
	
	downloadedData, err := mockStorage.Download(context.Background(), key)
	assert.NoError(t, err)
	assert.Equal(t, data, downloadedData)
	
	presignedURL, err := mockStorage.GeneratePresignedURL(context.Background(), key, 30*time.Minute)
	assert.NoError(t, err)
	assert.Contains(t, presignedURL, key)
	
	err = mockStorage.Delete(context.Background(), key)
	assert.NoError(t, err)
}

// Test concurrent operations using goroutines
func TestStorage_ConcurrentOperations(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping concurrent test in short mode")
	}
	
	mockStorage := NewMockStorage(t)
	
	// Setup mock for concurrent uploads
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("concurrent-test-%d.jpg", i)
		data := []byte(fmt.Sprintf("test data %d", i))
		
		mockStorage.On("Upload", context.Background(), key, data, "image/jpeg").
			Return("https://example.com/"+key, nil).
			Once()
	}
	
	// Execute concurrent uploads
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(index int) {
			defer func() { done <- true }()
			
			key := fmt.Sprintf("concurrent-test-%d.jpg", index)
			data := []byte(fmt.Sprintf("test data %d", index))
			
			url, err := mockStorage.Upload(context.Background(), key, data, "image/jpeg")
			assert.NoError(t, err)
			assert.Equal(t, "https://example.com/"+key, url)
		}(i)
	}
	
	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}