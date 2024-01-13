package utils

import (
	"bytes"
	"image"
	"os"

	"github.com/chai2010/webp"
)

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
