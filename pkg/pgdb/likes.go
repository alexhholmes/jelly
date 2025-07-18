package pgdb

func (c *Client) LikePhoto(userID, photoID string) error {
	return nil
}

func (c *Client) UnlikePhoto(userID, photoID string) error {
	return nil
}

func (c *Client) GetPhotoLikes(photoID string) ([]string, error) {
	var likes []string
	query := `SELECT user_id FROM photo_likes WHERE photo_id = $1`

	err := c.db.Select(&likes, query, photoID)
	if err != nil {
		// TODO
		return nil, err
	}

	return likes, nil
}

func (c *Client) CountPhotoLikes(photoID string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM photo_likes WHERE photo_id = $1`

	err := c.db.Get(&count, query, photoID)
	if err != nil {
		// TODO
		return 0, err
	}

	return count, nil
}

func (c *Client) IsPhotoLikedByUser(userID, photoID string) (bool, error) {
	return false, nil
}

func (c *Client) GetUserLikedPhotos(userID string) ([]string, error) {
	return nil, nil
}
