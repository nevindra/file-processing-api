package adapters

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"strings"
)

type BasicImageProcessor struct{}

func NewBasicImageProcessor() *BasicImageProcessor {
	return &BasicImageProcessor{}
}

func (p *BasicImageProcessor) ResizeImage(filename string, newWidth, newHeight int) ([]byte, error) {
	// ... implementation ...
	imgFile, err := os.ReadFile(filename)
	if err != nil {
		return ([]byte)(nil), err
	}

	// Convert byte data to a reader
	imgDataIn := bytes.NewReader(imgFile)

	// Decode the image
	img, format, err := image.Decode(imgDataIn)
	if err != nil {
		return ([]byte)(nil), err
	}

	// Create a new blank image with the new dimensions
	newImg := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	xRatio := float64(img.Bounds().Dx()) / float64(newWidth)
	yRatio := float64(img.Bounds().Dy()) / float64(newHeight)

	// Fill the new image with the resized data
	for newY := 0; newY < newHeight; newY++ {
		for newX := 0; newX < newWidth; newX++ {
			srcX := int(float64(newX) * xRatio)
			srcY := int(float64(newY) * yRatio)
			newImg.Set(newX, newY, img.At(srcX, srcY))
		}
	}

	// Prepare a buffer to receive the encoded image
	var buf bytes.Buffer

	// Encode the new image in the same format as the original
	switch strings.ToLower(format) {
	case "jpeg":
		err = jpeg.Encode(&buf, newImg, nil)
	case "png":
		err = png.Encode(&buf, newImg)
	default:
		return nil, fmt.Errorf("unsupported image format: %s", format)
	}
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
