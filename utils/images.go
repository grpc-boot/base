package utils

import (
	"bytes"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"

	"github.com/anthonynsimon/bild/transform"
)

func CropThumbnailImage(img image.Image, width, height int) *image.RGBA {
	return transform.Resize(img, width, height, transform.Linear)
}

func Bytes2JpgBytes(input []byte, quality int) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewBuffer(input))
	if err != nil {
		return nil, err
	}

	var buf = bytes.NewBuffer(nil)
	err = jpeg.Encode(buf, img, &jpeg.Options{Quality: quality})
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
