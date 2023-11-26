package http

import (
	"encoding/json"
	"fmt"
	"image-processing/internal/adapters"
	"image-processing/internal/usecases"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Server struct {
	Router       *gin.Engine
	imageService *usecases.ImageService
}

func NewServer(imageService *usecases.ImageService) *Server {
	router := gin.Default()
	return &Server{
		imageService: imageService,
		Router:       router,
	}
}

func (s *Server) SetupRoutes() {
	// ... setup your routes ...
	s.Router.GET("/ping", s.PingHandler)
	s.Router.POST("/resize", s.resizeHandler)
}

type ResizeRequest struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

func (s *Server) PingHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func (s *Server) resizeHandler(c *gin.Context) {
	// ... handler implementation using imageService ...
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Retrieve the JSON data from a form field, e.g., "json"
	jsonFormData := c.Request.FormValue("json")

	// Deserialize the JSON data into the struct
	var req ResizeRequest
	err = json.Unmarshal([]byte(jsonFormData), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Failed to parse JSON data: %v", err)})
		return
	}

	// Check if the file is an image and get its dimensions
	isImage, _, _, err := adapters.IsImageFile(file)
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
	// var localFilePath string
	err = c.SaveUploadedFile(file, file.Filename)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Resize the image
	resizedImgData, err := s.imageService.ResizeImage(file.Filename, req.Width, req.Height)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Determine the correct content type
	// TODO: return this from ResizeImage function
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
	_, err = c.Writer.Write(resizedImgData)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
}
