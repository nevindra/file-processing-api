package adapters

import (
	"image"
	"mime/multipart"
	"net/http"
	"strings"
)

func IsImageFile(file *multipart.FileHeader) (bool, int, int, error) {
	// Open the file
	openedFile, err := file.Open()
	if err != nil {
		return false, 0, 0, err
	}
	defer openedFile.Close()

	// Read a small chunk to determine file type
	buffer := make([]byte, 512)
	_, err = openedFile.Read(buffer)
	if err != nil {
		return false, 0, 0, err
	}

	// Determine the file type
	fileType := http.DetectContentType(buffer)
	if !strings.HasPrefix(fileType, "image/") {
		return false, 0, 0, nil
	}

	// Reset the read pointer to start of the file after reading
	openedFile.Seek(0, 0)

	// Get image dimensions
	config, _, err := image.DecodeConfig(openedFile)
	if err != nil {
		return false, 0, 0, err
	}

	return true, config.Width, config.Height, nil
}
