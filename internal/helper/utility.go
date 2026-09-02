package helper

import (
	"bytes"
	"errors"
	"image"

	"github.com/rs/xid"
)

// FetchNewID returns a new unique identifier
func FetchNewID() string {
	guid := xid.New()

	return guid.String()
}

// ToImage will convert an byte array to an Image
func ToImage(data []byte) (image.Image, error) {
	if img, _, err := image.Decode(bytes.NewReader(data)); err == nil {
		return img, nil
	}

	return nil, errors.New("Failed to Convert to Image")
}
