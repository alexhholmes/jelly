package photo

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"jelly/pkg/api/gen"
	"jelly/pkg/api/util"
	"jelly/pkg/config"
	"jelly/pkg/pgdb"
)

// PhotoHandler implements photo upload endpoints.
type PhotoHandler struct {
	DB *pgdb.Client
}

// UploadPhoto handles photo upload with optional caption and tags.
// POST /photo
func (h PhotoHandler) UploadPhoto(w http.ResponseWriter, r *http.Request) {
	logger := r.Context().Value(util.ContextLogger).(*slog.Logger)

	// Parse multipart form using configured max file size
	maxFileSize := config.GetPhotoMaxFileSizeBytes()
	err := r.ParseMultipartForm(maxFileSize)
	if err != nil {
		logger.Error("Failed to parse multipart form", "error", err)
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

	// Read file data (in real implementation, would save to storage)
	_, err = io.ReadAll(file)
	if err != nil {
		logger.Error("Failed to read file", "error", err)
		http.Error(w, util.ErrMsgFailedToReadFile, http.StatusInternalServerError)
		return
	}

	// Get optional caption
	caption := r.FormValue("caption")

	// Get optional tags
	var tags []string
	if tagValues := r.Form["tags"]; len(tagValues) > 0 {
		tags = tagValues
	}

	// Generate mock photo ID and URL
	photoID := fmt.Sprintf("photo_%d", time.Now().Unix())
	photoURL := fmt.Sprintf("https://example.com/photos/%s", fileHeader.Filename)

	// Create response
	photo := gen.Photo{
		Id:         photoID,
		Url:        photoURL,
		Caption:    &caption,
		Tags:       &tags,
		UploadedAt: time.Now(),
	}

	resp := gen.PhotoUploadResponse{
		Photo:   photo,
		Message: util.StringPtr("Photo uploaded successfully"),
	}

	logger.Info("Photo uploaded", "photo_id", photoID, "filename", fileHeader.Filename)

	util.WriteJSONResponse(w, logger, http.StatusOK, resp)
}
