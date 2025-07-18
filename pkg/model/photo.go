package model

import (
	"time"

	"jelly/pkg/api/v1/gen"
)

// RawPhoto represents the original unprocessed photo data
type RawPhoto struct {
	ID               string     `json:"id" db:"id"`
	UserID           string     `json:"user_id" db:"user_id"`
	OriginalFilename string     `json:"original_filename" db:"original_filename"`
	StorageURL       string     `json:"storage_url" db:"storage_url"`
	FileSize         int64      `json:"file_size" db:"file_size"`
	MimeType         string     `json:"mime_type" db:"mime_type"`
	MD5Hash          string     `json:"md5_hash" db:"md5_hash"`
	Width            *int       `json:"width,omitempty" db:"width"`
	Height           *int       `json:"height,omitempty" db:"height"`
	ExifData         *string    `json:"exif_data,omitempty" db:"exif_data"`
	UploadedAt       time.Time  `json:"uploaded_at" db:"uploaded_at"`
	ProcessedAt      *time.Time `json:"processed_at,omitempty" db:"processed_at"`
	ScheduleDeletion *time.Time `json:"schedule_deletion,omitempty" db:"schedule_deletion"`
}

func (rp *RawPhoto) ToRawPhotoDetails() gen.RawPhotoDetails {
	return gen.RawPhotoDetails{
		Id:               rp.ID,
		UserId:           rp.UserID,
		OriginalFilename: rp.OriginalFilename,
		StorageUrl:       rp.StorageURL,
		FileSize:         rp.FileSize,
		MimeType:         rp.MimeType,
		Md5Hash:          rp.MD5Hash,
		Width:            rp.Width,
		Height:           rp.Height,
		ProcessedAt:      rp.ProcessedAt,
		UploadedAt:       rp.UploadedAt,
	}
}

// Photo represents a photo in the database
type Photo struct {
	ID               string     `json:"id" db:"id"`
	RawPhotoID       string     `json:"raw_photo_id" db:"raw_photo_id"`
	UserID           string     `json:"user_id" db:"user_id"`
	Filename         string     `json:"filename" db:"filename"`
	OriginalURL      string     `json:"original_url" db:"original_url"`
	ThumbnailURL     string     `json:"thumbnail_url" db:"thumbnail_url"`
	Caption          *string    `json:"caption,omitempty" db:"caption"`
	Tags             []string   `json:"tags,omitempty" db:"tags"`
	FileSize         int64      `json:"file_size" db:"file_size"`
	MimeType         string     `json:"mime_type" db:"mime_type"`
	Width            *int       `json:"width,omitempty" db:"width"`
	Height           *int       `json:"height,omitempty" db:"height"`
	UploadedAt       time.Time  `json:"uploaded_at" db:"uploaded_at"`
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`
	ScheduleDeletion *time.Time `json:"schedule_deletion,omitempty" db:"schedule_deletion"`
}

func (p *Photo) ToPhotoDetails() gen.PhotoDetails {
	return gen.PhotoDetails{
		Id:               p.ID,
		UserId:           p.UserID,
		Caption:          p.Caption,
		FileSize:         p.FileSize,
		Filename:         p.Filename,
		Height:           p.Height,
		Width:            p.Width,
		MimeType:         p.MimeType,
		OriginalUrl:      p.OriginalURL,
		ThumbnailUrl:     p.ThumbnailURL,
		RawPhotoId:       p.RawPhotoID,
		ScheduleDeletion: p.ScheduleDeletion,
		Tags:             &p.Tags,
		UploadedAt:       p.UploadedAt,
		UpdatedAt:        p.UpdatedAt,
	}
}
