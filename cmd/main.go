package main

import (
	"image-processing/internal/adapters"
	"image-processing/internal/delivery/http"
	"image-processing/internal/usecases"
)

func main() {
	processor := adapters.NewBasicImageProcessor()
	service := usecases.NewImageService(processor)
	server := http.NewServer(service)
	server.SetupRoutes()

	if err := server.Router.Run(":8080"); err != nil {
		// Handle the error as appropriate
		errorMessage := err.Error()
		panic(errorMessage)
	}
}
