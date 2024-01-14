package utils

import (
	"bytes"
	"golang.org/x/image/draw"
	"image"
	"image/jpeg"
	"os"
	"testing"
)

var (
	file, _ = os.Open("./test.jpeg")
)

func TestThumbnailImage(t *testing.T) {
	img, _, err := image.Decode(file)
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}

	rgbaImg := image.NewRGBA(img.Bounds())
	draw.Draw(rgbaImg, rgbaImg.Bounds(), img, img.Bounds().Min, draw.Over)

	out := CropThumbnailImage(rgbaImg, 1280, 720)
	buf := bytes.NewBuffer(nil)
	err = jpeg.Encode(buf, out, nil)
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}

	err = os.WriteFile("./thumbnail.jpg", buf.Bytes(), 0666)
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}

	/*webpData, err := Bytes2WebpBytes(buf.Bytes(), jpeg.DefaultQuality)
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}
	err = os.WriteFile("./thumbnail.webp", webpData, 0666)
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}*/
}
