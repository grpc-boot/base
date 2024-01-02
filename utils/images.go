package utils

import (
	"bytes"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"os"

	"github.com/anthonynsimon/bild/transform"
	"github.com/chai2010/webp"
)

func CropThumbnailImage(img image.Image, width, height int) *image.RGBA {
	return transform.Resize(img, width, height, transform.Linear)
}

func SubImage(img *image.RGBA, width, height int) image.Image {
	var (
		bounds       = img.Bounds()
		oldW, oldH   = bounds.Max.X - bounds.Min.X, bounds.Max.Y - bounds.Min.Y
		xDiff, yDiff = 0, 0
	)

	if oldW == width && oldH == height {
		return img
	}

	if oldW > width && oldH > height {
		for oldW-2*xDiff > width {
			xDiff++
		}

		for oldH-2*yDiff > width {
			yDiff++
		}

		return img.SubImage(image.Rect(bounds.Min.X+xDiff, bounds.Min.Y+yDiff, bounds.Min.X+xDiff+width, bounds.Min.Y+yDiff+height))
	}

	return img
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

func Bytes2WebpBytes(input []byte, quality float32) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewBuffer(input))
	if err != nil {
		return nil, err
	}

	webpBytes, err := webp.EncodeRGBA(img, quality)
	if err != nil {
		return nil, err
	}

	return webpBytes, nil
}

func Image2Webp(in, out string, quality float32) error {
	fileBytes, err := os.ReadFile(in)
	if err != nil {
		return err
	}

	webpBytes, err := Bytes2WebpBytes(fileBytes, quality)
	if err != nil {
		return err
	}

	return os.WriteFile(out, webpBytes, 0666)
}
