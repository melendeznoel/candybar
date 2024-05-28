package media

import (
	"image"
	"time"
)

// Figure is the basic contract
type Figure struct {
	ID           string
	URL          string
	Name         string
	Type         string
	Data         []byte
	DataLength   int64
	LastModified time.Time
	Relatives    []Relative
	image        image.Image
	Expires      string
}

// ImageCompareResult is the image compate response contract
type ImageCompareResult struct {
	Images []Figure
	ID     string
}

// Relative is the basic contract
type Relative struct {
	Rank       int64
	ID         string
	Percentage int64
}
