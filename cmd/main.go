package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type Image struct {
	Width  int
	Height int
	Data   []byte
}

var req struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

func ResizeImage(filename string, newWidth, newHeight int) ([]byte, error) {
	// Open the file
	imgFile, err := os.ReadFile(filename)
	if err != nil {
		return ([]byte)(nil), err
	}

	// Convert byte data to io.Reader
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

func isImageFile(file *multipart.FileHeader) (bool, int, int, error) {
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

func main() {
	r := gin.Default()

	r.POST("/upload", func(c *gin.Context) {
		// Retrieve the file from the form-data
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Retrieve the JSON data from a form field, e.g., "json"
		jsonFormData := c.Request.FormValue("json")

		// Deserialize the JSON data into the struct
		err = json.Unmarshal([]byte(jsonFormData), &req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Failed to parse JSON data: %v", err)})
			return
		}

		// Check if the file is an image and get its dimensions
		isImage, _, _, err := isImageFile(file)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		if !isImage {
			c.JSON(400, gin.H{
				"error": "Invalid image file",
			})
			return
		}

		// Save the file
		err = c.SaveUploadedFile(file, file.Filename)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Resize the image
		resizedImgData, err := ResizeImage(file.Filename, req.Width, req.Height)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Determine the correct content type (you might want to return this from your ResizeImage function)
		var fileExtension string
		if strings.HasSuffix(file.Filename, ".png") {
			c.Writer.Header().Set("Content-Type", "image/png")
			fileExtension = ".png"
		} else {
			c.Writer.Header().Set("Content-Type", "image/jpeg")
			fileExtension = ".jpeg"
		}

		// Set the headers for a file download response
		c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"resized-image%s\"", fileExtension))
		c.Writer.Header().Set("Content-Length", fmt.Sprint(len(resizedImgData)))

		// Write the image data to the response
		c.Writer.Write(resizedImgData)
	})

	r.Run(":8080")
}
