package photo

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/google/uuid"

	"jelly/pkg/api/gen"
	"jelly/pkg/api/util"
	"jelly/pkg/config"
	"jelly/pkg/model"
	"jelly/pkg/pgdb"
)

// PhotoHandler implements photo upload endpoints.
type PhotoHandler struct {
	DB *pgdb.Client
}

// UploadPhoto handles photo upload with optional caption and tags and processing.
// POST /photo
func (h PhotoHandler) UploadPhoto(w http.ResponseWriter, r *http.Request) {
	logger := r.Context().Value(util.ContextLogger).(*slog.Logger)

	// Parse multipart form using configured max file size
	maxFileSize := config.GetPhotoMaxFileSizeBytes()
	err := r.ParseMultipartForm(maxFileSize)
	if errors.Is(err, multipart.ErrMessageTooLarge) {
		// If the error is due to file size, return specific error
		logger.Info("File size too large",
			"error", err, "max_size_mb",
			maxFileSize/(1024*1024), "file_size_mb", r.ContentLength/(1024*1024),
		)
		http.Error(w, util.ErrMsgFileTooLarge, http.StatusBadRequest)
		return
	} else if err != nil {
		logger.Info("Failed to parse form", "error", err, "max_size_mb")
		http.Error(w, util.ErrMsgFailedToParseForm, http.StatusBadRequest)
		return
	}

	// Get uploaded file
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		logger.Error("Failed to get uploaded file", "error", err)
		http.Error(w, util.ErrMsgFileRequired, http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read file data
	bytes, err := io.ReadAll(file)
	if err != nil {
		logger.Error("Failed to read file", "error", err)
		http.Error(w, util.ErrMsgFailedToReadFile, http.StatusInternalServerError)
		return
	}

	// STEP 1: Generate unique IDs and extract metadata from uploaded file
	// - Create unique ID for raw photo record
	// - Extract original filename, file size, MIME type
	// - Calculate MD5 and SHA256 hashes for integrity verification
	// - Extract image dimensions (width/height) from file headers
	// - Parse EXIF data if available (camera settings, GPS, etc.)
	// - Set upload timestamp

	rawMetadata := model.RawPhoto{
		ID: uuid.New().String(),
		// UserID:           "",
		OriginalFilename: fileHeader.Filename,
		StorageURL:       "",
		FileSize:         fileHeader.Size,
		MimeType:         http.DetectContentType(bytes),
		MD5Hash:          util.CalculateMD5(bytes),
		UploadedAt:       time.Now(),
	}

	// Check if valid image file type
	if rawMetadata.MimeType != "image/jpeg" && rawMetadata.MimeType != "image/png" {
		logger.Info("Unsupported file type", "mime_type", rawMetadata.MimeType)
		http.Error(w, util.ErrMsgUnsupportedFileType, http.StatusBadRequest)
		return
	}

	// Async upload the raw image

	// Process the image, this will also update the rawMetadata with dimensions
	// and EXIF data if available.
	processed, err := util.ProcessPhoto(bytes, rawMetadata)

	// Update rawMetadata with processed data and write to database

	// Deferred db and s3 cleanup if any of the upcoming operations fail
	defer func() {
		if err != nil {
			// Write db schedule deletion in 7 days
			updateMetadata := model.RawPhoto{
				ScheduleDeletion: util.GetTimePtr(time.Now().Add(7 * 24 * time.Hour)),
			}
			if dbErr := h.DB.UpdateRawPhoto(r.Context(), rawMetadata.ID, updateMetadata); dbErr != nil {
				logger.Error("Failed to update metadata for scheduled deletion",
					"error", dbErr,
					"raw_photo_id", rawMetadata.ID,
					"schedule_deletion", updateMetadata.ScheduleDeletion,
				)
				return
			}
		}
	}()

	// STEP 2: Upload raw photo to storage
	// - Upload original file bytes to object storage (S3/local)
	// - Generate storage URL/path for raw file
	// - Handle upload errors and retry logic

	// STEP 3: Save raw photo metadata to database
	// - Create RawPhoto record with all metadata
	// - Save to database with transaction
	// - Handle database errors

	// STEP 4: Process the photo for web display
	// - Resize image to standard web sizes (e.g., 1920x1080 max)
	// - Apply image optimization (compression, format conversion)
	// - Strip sensitive EXIF data for privacy
	// - Generate thumbnail (e.g., 300x300)
	// - Convert to web-friendly formats (JPEG, WebP)

	// STEP 5: Upload processed images to storage
	// - Upload processed/optimized image to storage
	// - Upload thumbnail to storage
	// - Generate public URLs for both processed and thumbnail
	// - Handle upload errors

	// STEP 6: Save processed photo metadata to database
	// - Create Photo record with reference to RawPhoto
	// - Include user caption and tags
	// - Store processed image URLs and metadata
	// - Update RawPhoto.ProcessedAt timestamp
	// - Save with transaction

	// STEP 7: Return response to client
	// - Return processed photo metadata (not raw)
	// - Include public URLs for display
	// - Include user-provided caption and tags

	// TODO: Implement the above steps

	// Get optional caption and tags
	// caption := r.FormValue("caption")
	// var tags []string
	// if tagValues := r.Form["tags"]; len(tagValues) > 0 {
	// 	tags = tagValues
	// }

	// TODO: Save photo metadata to database
	// err = h.DB.SavePhoto(photoModel)
	// if err != nil {
	//     logger.Error("Failed to save photo metadata", "error", err)
	//     http.Error(w, "Failed to save photo", http.StatusInternalServerError)
	//     return
	// }

	// Create response
	photo := gen.Photo{
		Id: "",
		// Url:        photoModel.OriginalURL,
		// Caption:    photoModel.Caption,
		// Tags:       &photoModel.Tags,
		// UploadedAt: photoModel.UploadedAt,
	}

	resp := gen.PhotoUploadResponse{
		Photo:   photo,
		Message: util.StringPtr("Photo uploaded successfully"),
	}

	logger.Info("Photo uploaded", "photo_id", photo.Id, "filename", fileHeader.Filename)

	util.WriteJSONResponse(w, logger, http.StatusOK, resp)
}
