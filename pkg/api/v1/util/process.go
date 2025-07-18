package util

import (
	"crypto/md5"
	"encoding/hex"
	"time"

	"jelly/pkg/model"
)

func CalculateMD5(bytes []byte) string {
	hash := md5.Sum(bytes)
	return hex.EncodeToString(hash[:])
}

func GetTimePointer(t time.Time) *time.Time {
	if t.IsZero() {
		return nil
	}
	return &t
}

func ProcessPhoto(photo []byte, metadata model.RawPhoto, upload func(photo []byte,
	name string)) ([]byte, error) {
	return photo, nil
}
