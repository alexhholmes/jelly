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

	// rawMetadata is for the original unprocessed photo
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
			updateMetadata := model.RawPhoto{
				// TODO add configuration to set deletion schedule
				ScheduleDeletion: util.GetTimePointer(time.Now().Add(7 * 24 * time.Hour)),
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

	// Get optional caption and tags
	caption := r.FormValue("caption")
	var tags []string
	if tagValues := r.Form["tags"]; len(tagValues) > 0 {
		tags = tagValues
	}

	// TODO: Save photo metadata to database
	// err = h.DB.SavePhoto(photoModel)
	// if err != nil {
	//     logger.Error("Failed to save photo metadata", "error", err)
	//     http.Error(w, "Failed to save photo", http.StatusInternalServerError)
	//     return
	// }

	// Create response
	photo := gen.Photo{
		Id:         "",
		Url:        "",
		Caption:    &caption,
		Tags:       &tags,
		UploadedAt: rawMetadata.UploadedAt,
	}

	resp := gen.PhotoUploadResponse{
		Photo:   photo,
		Message: util.StringPtr("Photo uploaded successfully"),
	}

	logger.Info("Photo uploaded", "photo_id", photo.Id, "filename", fileHeader.Filename)

	util.WriteJSONResponse(w, logger, http.StatusOK, resp)
}
