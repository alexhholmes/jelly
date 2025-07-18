package pgdb

import (
	"time"

	"github.com/google/uuid"

	"jelly/pkg/model"
)

func (c *Client) CreatePhoto(photo model.Photo) error {
	return nil
}

func (c *Client) GetPhotoByID(photoID uuid.UUID) (model.Photo, error) {
	return model.Photo{}, nil
}

func (c *Client) UpdatePhoto(photo model.Photo) error {
	return nil
}

func (c *Client) DeletePhoto(photoID uuid.UUID, deletionDuration time.Duration) error {
	return nil
}
